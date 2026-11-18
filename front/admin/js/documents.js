// Back-office : PDF archivés (contrats, attestations, fiches…) et codes de dépôt.
// Les PDF sont générés et archivés par l'API (table documents) ; cette page permet de les retrouver,
// de les ouvrir, de les télécharger et de créer un contrat ou un document administratif.

const ONGLETS_DOCUMENTS = ["documents", "codes"];
const STATUTS_CODES = ["validee", "deposee", "recuperee"];
const formulaireVide = () => ({ type_document: "contrat", titre: "", contenu: "" });

demarrerApp({
    data() {
        return {
            onglets: ONGLETS_DOCUMENTS,
            onglet: "documents",
            filtresDocuments: { recherche: "", type: "", page: 1 },
            documents: { documents: [], par_type: [], total: 0, par_page: 20 },
            filtresCodes: { recherche: "", statut: "", page: 1 },
            codes: { codes: [], total: 0, par_page: 20 },
            statutsCodes: STATUTS_CODES,
            chargement: false,
            numeroRequete: 0,
            erreur: "",
            succes: "",
            formulaireOuvert: false,
            formulaire: formulaireVide(),
            beneficiaire: null,
            rechercheBeneficiaire: "",
            propositions: [],
            rechercheEnCours: false,
            numeroRechercheBeneficiaire: 0,
            erreurFormulaire: "",
            enregistrement: false,
        };
    },

    computed: {
        totalTousTypes() {
            return this.documents.par_type.reduce((total, c) => total + c.nombre, 0);
        },
        pagesDocuments() {
            return Math.max(1, Math.ceil(this.documents.total / this.documents.par_page));
        },
        pagesCodes() {
            return Math.max(1, Math.ceil(this.codes.total / this.codes.par_page));
        },
    },

    watch: {
        "filtresDocuments.recherche"() {
            this.filtresDocuments.page = 1;
            this.chargerDocumentsEnDirect();
        },
        "filtresCodes.recherche"() {
            this.filtresCodes.page = 1;
            this.chargerCodesEnDirect();
        },
        "filtresCodes.statut"() {
            this.filtresCodes.page = 1;
            this.chargerCodes();
        },
        rechercheBeneficiaire() {
            this.chercherBeneficiaireEnDirect();
        },
    },

    created() {
        this.chargerDocumentsEnDirect = differer(() => this.chargerDocuments());
        this.chargerCodesEnDirect = differer(() => this.chargerCodes());
        this.chercherBeneficiaireEnDirect = differer(() => this.chercherBeneficiaire(), 250);
    },

    mounted() {
        const onglet = new URLSearchParams(window.location.search).get("onglet");
        if (this.onglets.includes(onglet)) {
            this.onglet = onglet;
        }
        this.chargerOnglet();
    },

    methods: {
        formaterDate(date) {
            return new Date(date).toLocaleDateString(langueCourante(), { day: "numeric", month: "short", year: "numeric" });
        },

        formaterTaille(octets) {
            return octets < 1024 * 1024
                ? `${Math.max(1, Math.round(octets / 1024)).toLocaleString(langueCourante())} ko`
                : `${(octets / 1024 / 1024).toLocaleString(langueCourante(), { maximumFractionDigits: 1 })} Mo`;
        },

        changerOnglet(onglet) {
            this.onglet = onglet;
            this.erreur = "";
            this.succes = "";
            const url = new URL(window.location.href);
            url.searchParams.set("onglet", onglet);
            history.replaceState(null, "", url);
            this.chargerOnglet();
        },

        chargerOnglet() {
            if (this.onglet === "documents") {
                this.chargerDocuments();
            } else {
                this.chargerCodes();
            }
        },

        // Une seule requête « en cours » compte : une réponse plus ancienne arrivée en retard est ignorée
        async charger(chemin, filtres) {
            const numero = ++this.numeroRequete;
            const parametres = new URLSearchParams(Object.entries(filtres).filter(([, valeur]) => valeur !== ""));
            this.chargement = true;
            try {
                const reponse = await appelerApi(`${chemin}?${parametres}`);
                const donnees = await reponse.json();
                if (numero !== this.numeroRequete) {
                    return null;
                }
                if (!reponse.ok) {
                    this.erreur = donnees.erreur || this.$t("commun.erreur_chargement");
                    return null;
                }
                this.erreur = "";
                return donnees;
            } catch (erreur) {
                if (numero === this.numeroRequete) {
                    this.erreur = this.$t("commun.erreur_chargement");
                }
                return null;
            } finally {
                if (numero === this.numeroRequete) {
                    this.chargement = false;
                }
            }
        },

        async chargerDocuments() {
            const donnees = await this.charger("/admin/documents", this.filtresDocuments);
            if (donnees) {
                this.documents = donnees;
            }
        },

        async chargerCodes() {
            const donnees = await this.charger("/admin/codes", this.filtresCodes);
            if (donnees) {
                this.codes = donnees;
            }
        },

        choisirType(type) {
            this.filtresDocuments.type = type;
            this.filtresDocuments.page = 1;
            this.chargerDocuments();
        },

        changerPageDocuments(sens) {
            this.filtresDocuments.page += sens;
            this.chargerDocuments();
        },

        changerPageCodes(sens) {
            this.filtresCodes.page += sens;
            this.chargerCodes();
        },

        nomFichier(document) {
            const titre = document.titre.normalize("NFD").replace(/[̀-ͯ]/g, "")
                .replace(/[^a-zA-Z0-9]+/g, "-").replace(/^-|-$/g, "").toLowerCase().slice(0, 60);
            return `${titre || "document"}-${document.id}.pdf`;
        },

        async ouvrir(document) {
            try {
                await recupererPdf(`/documents/${document.id}/pdf`);
            } catch (erreur) {
                this.erreur = erreur.message;
            }
        },

        async telecharger(document) {
            try {
                await recupererPdf(`/documents/${document.id}/pdf`, { nomFichier: this.nomFichier(document), telecharger: true });
            } catch (erreur) {
                this.erreur = erreur.message;
            }
        },

        ouvrirNouveau() {
            this.formulaire = formulaireVide();
            this.beneficiaire = null;
            this.rechercheBeneficiaire = "";
            this.propositions = [];
            this.erreurFormulaire = "";
            this.formulaireOuvert = true;
        },

        async chercherBeneficiaire() {
            const recherche = this.rechercheBeneficiaire.trim();
            const numero = ++this.numeroRechercheBeneficiaire;
            if (recherche.length < 2) {
                this.propositions = [];
                return;
            }
            this.rechercheEnCours = true;
            try {
                const reponse = await appelerApi(`/utilisateurs?${new URLSearchParams({ recherche })}`);
                const donnees = await reponse.json();
                if (numero === this.numeroRechercheBeneficiaire && reponse.ok) {
                    this.propositions = donnees.utilisateurs.slice(0, 6);
                }
            } catch (erreur) {
            } finally {
                if (numero === this.numeroRechercheBeneficiaire) {
                    this.rechercheEnCours = false;
                }
            }
        },

        choisirBeneficiaire(utilisateur) {
            this.beneficiaire = utilisateur;
            this.propositions = [];
        },

        changerBeneficiaire() {
            this.beneficiaire = null;
            this.$nextTick(() => document.getElementById("beneficiaire")?.focus());
        },

        async creer() {
            if (!this.beneficiaire) {
                this.erreurFormulaire = this.$t("admin.documents.beneficiaire_obligatoire");
                return;
            }
            this.erreurFormulaire = "";
            this.enregistrement = true;
            try {
                const reponse = await appelerApi("/admin/documents", {
                    method: "POST",
                    body: JSON.stringify({ ...this.formulaire, utilisateur_id: this.beneficiaire.id }),
                });
                const donnees = await reponse.json();
                if (!reponse.ok) {
                    this.erreurFormulaire = donnees.erreur || this.$t("commun.erreur");
                    return;
                }
                this.formulaireOuvert = false;
                this.succes = this.$t("admin.documents.cree", {
                    titre: this.formulaire.titre,
                    personne: `${this.beneficiaire.prenom} ${this.beneficiaire.nom}`,
                });
                this.filtresDocuments = { recherche: "", type: "", page: 1 };
                await this.$nextTick();
                this.chargerDocuments();
            } catch (erreur) {
                this.erreurFormulaire = this.$t("commun.serveur_injoignable");
            } finally {
                this.enregistrement = false;
            }
        },
    },
});
