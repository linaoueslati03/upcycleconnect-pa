// Dépôts en conteneur, selon le rôle :
//  - particulier : demande de dépôt (objet, matériau, poids approximatif) et suivi avec les codes ;
//  - salarié : validation des demandes (avec correction du poids) puis confirmation du dépôt ;
//  - professionnel : récupération d'un objet déposé grâce à son code-barre.

const formulaireDepotVide = () => ({ description_objet: "", conteneur_id: null, categorie_id: null, poids_kg: null });

demarrerApp({
    data() {
        return {
            depots: [],
            conteneurs: [],
            categories: [],
            etapes: ETAPES_DEPOT,
            formulaire: formulaireDepotVide(),
            erreur: "",
            erreurListe: "",
            confirmation: false,
            demandeOuverte: false,
            scannerOuvert: false,
            poidsCorriges: {},

            statutSelectionne: "demande",
            filtres: [
                { nom: "salaries.depots.a_valider", valeur: "demande" },
                { nom: "salaries.depots.valides", valeur: "validee" },
                { nom: "salaries.depots.deposes", valeur: "deposee" },
                { nom: "commun.tous", valeur: "" },
            ],
            statutSuivant: { demande: "validee", validee: "deposee" },
            libellesAction: { demande: "salaries.depots.valider_demande", validee: "salaries.depots.confirmer_depot" },

            codeObjet: "",
            objetTrouve: null,
            succes: "",
            recuperationEnCours: false,
        };
    },

    computed: {
        categoriesActives() {
            return this.categories.filter((c) => c.actif);
        },
    },

    async mounted() {
        const role = roleConnecte();
        if (role === "particulier" || role === "salarie") {
            const reponse = await appelerApi("/categories");
            if (reponse.ok) {
                this.categories = await reponse.json();
            }
        }
        if (role === "particulier") {
            const reponse = await appelerApi("/conteneurs");
            if (reponse.ok) {
                this.conteneurs = await reponse.json();
            }
            this.chargerDepots();
        } else if (role === "salarie") {
            this.chargerDepotsAValider();
        }
    },

    methods: {
        formaterDate,

        nomConteneur(id) {
            const conteneur = this.conteneurs.find((c) => c.id === id);
            return conteneur ? conteneur.site : "—";
        },

        nomMateriau(id) {
            const categorie = this.categories.find((c) => c.id === id);
            return categorie ? t("categorie." + categorie.code, categorie.libelle) : "";
        },

        formaterPoids(poids) {
            return `${Number(poids).toLocaleString(langueCourante(), { maximumFractionDigits: 1 })} kg`;
        },

        indexStatut(statut) {
            return ETAPES_DEPOT.indexOf(statut);
        },

        async chargerDepots() {
            const reponse = await appelerApi("/moi/depots");
            if (!reponse.ok) {
                this.erreurListe = t("depots.erreur_chargement");
                return;
            }
            this.depots = await reponse.json();
        },

        ouvrirDemande() {
            this.formulaire = formulaireDepotVide();
            this.erreur = "";
            this.confirmation = false;
            this.demandeOuverte = true;
        },

        async demanderDepot() {
            this.erreur = "";
            this.confirmation = false;
            const corps = { ...this.formulaire, poids_kg: this.formulaire.poids_kg ? Number(this.formulaire.poids_kg) : null };
            const reponse = await appelerApi("/depots", { method: "POST", body: JSON.stringify(corps) });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.confirmation = true;
            this.demandeOuverte = false;
            this.formulaire = formulaireDepotVide();
            this.chargerDepots();
        },

        async chargerDepotsAValider() {
            const chemin = this.statutSelectionne ? `/depots?statut=${this.statutSelectionne}` : "/depots";
            const reponse = await appelerApi(chemin);
            if (!reponse.ok) {
                this.erreur = t("salaries.depots.erreur_chargement");
                return;
            }
            this.erreur = "";
            this.depots = await reponse.json();
            this.poidsCorriges = Object.fromEntries(this.depots.map((d) => [d.id, d.poids_kg ?? ""]));
        },

        async avancer(depot) {
            const corps = { statut: this.statutSuivant[depot.statut] };
            // À la validation, le salarié peut corriger le poids annoncé
            const poids = this.poidsCorriges[depot.id];
            if (depot.statut === "demande" && poids !== "" && poids !== null && Number(poids) !== depot.poids_kg) {
                corps.poids_kg = Number(poids);
            }
            const reponse = await appelerApi(`/depots/${depot.id}/statut`, { method: "PUT", body: JSON.stringify(corps) });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.chargerDepotsAValider();
        },

        ouvrirScanner() {
            this.codeObjet = "";
            this.objetTrouve = null;
            this.erreur = "";
            this.succes = "";
            this.scannerOuvert = true;
        },

        // Le code-barre est celui du bon de dépôt ; seuls les objets déposés dans un conteneur sont trouvés
        async verifierCode() {
            this.erreur = "";
            this.succes = "";
            this.objetTrouve = null;
            const code = this.codeObjet.trim();
            if (!code) {
                return;
            }
            const reponse = await appelerApi(`/pro/depots/code-barre/${encodeURIComponent(code)}`);
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.objetTrouve = await reponse.json();
        },

        async confirmerRecuperation() {
            this.erreur = "";
            this.recuperationEnCours = true;
            const reponse = await appelerApi(`/depots/${this.objetTrouve.id}/statut`, {
                method: "PUT",
                body: JSON.stringify({ statut: "recuperee" }),
            });
            this.recuperationEnCours = false;
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.succes = t("depots.pro.recuperation_confirmee", { objet: this.objetTrouve.description_objet });
            this.objetTrouve = null;
            this.codeObjet = "";
        },
    },
});
