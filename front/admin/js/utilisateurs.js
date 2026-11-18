const ROLE_PARTICULIER = 1;
const ROLE_PROFESSIONNEL = 2;
const FORMULAIRE_VIDE = { id: null, nom: "", prenom: "", email: "", role_id: 1, mot_de_passe: "", langue_preferee_id: null };

demarrerApp({
    data() {
        return {
            utilisateurs: [],
            filtres: { recherche: "", role_id: "", statut: "" },
            page: 1,
            total: 0,
            parPage: 20,
            motDePasseTemporaire: null,
            formulaire: { ...FORMULAIRE_VIDE },
            formulaireOuvert: false,
            prediction: null,
            erreur: "",
            erreurListe: "",
        };
    },

    computed: {
        nombrePages() {
            return Math.max(1, Math.ceil(this.total / this.parPage));
        },
    },

    watch: {
        "filtres.recherche"() {
            this.rechercherEnDirect();
        },
    },

    created() {
        this.rechercherEnDirect = differer(() => this.rechercher());
    },

    mounted() {
        this.chargerUtilisateurs();
    },

    methods: {
        async chargerUtilisateurs() {
            const numero = (this.numeroRequete = (this.numeroRequete || 0) + 1);
            const parametres = new URLSearchParams({ page: this.page });
            for (const [nom, valeur] of Object.entries(this.filtres)) {
                if (valeur !== "") {
                    parametres.set(nom, valeur);
                }
            }
            const reponse = await appelerApi(`/utilisateurs?${parametres}`);
            if (numero !== this.numeroRequete) {
                return;
            }
            if (!reponse.ok) {
                this.erreurListe = t("admin.utilisateurs.erreur_chargement");
                return;
            }
            const resultat = await reponse.json();
            this.utilisateurs = resultat.utilisateurs;
            this.total = resultat.total;
            this.parPage = resultat.par_page;
        },

        rechercher() {
            this.page = 1;
            this.chargerUtilisateurs();
        },

        changerPage(page) {
            this.page = page;
            this.chargerUtilisateurs();
        },

        async reinitialiserMotDePasse(utilisateur) {
            if (!confirm(t("admin.utilisateurs.confirmer_reinitialisation"))) {
                return;
            }
            this.erreurListe = "";
            const reponse = await appelerApi(`/utilisateurs/${utilisateur.id}/mot-de-passe`, { method: "PUT" });
            const resultat = await reponse.json();
            if (!reponse.ok) {
                this.erreurListe = resultat.erreur;
                return;
            }
            this.motDePasseTemporaire = {
                personne: `${utilisateur.prenom} ${utilisateur.nom}`,
                valeur: resultat.mot_de_passe_temporaire,
            };
        },

        reinitialiserFormulaire() {
            this.formulaire = { ...FORMULAIRE_VIDE };
            this.prediction = null;
            this.erreur = "";
        },

        async chargerPrediction(utilisateur) {
            this.prediction = null;
            if (![ROLE_PARTICULIER, ROLE_PROFESSIONNEL].includes(utilisateur.role_id)) {
                return;
            }
            const reponse = await appelerApi(`/admin/predictions/${utilisateur.id}`);
            if (!reponse.ok || this.formulaire.id !== utilisateur.id) {
                return;
            }
            const { prediction } = await reponse.json();
            this.prediction = {
                classe: prediction.classe,
                probabilite: prediction.probabilites.find((p) => p.classe === prediction.classe).probabilite,
            };
        },

        ouvrirNouveau() {
            this.reinitialiserFormulaire();
            this.formulaireOuvert = true;
        },

        modifier(utilisateur) {
            this.formulaire = {
                id: utilisateur.id,
                nom: utilisateur.nom,
                prenom: utilisateur.prenom,
                email: utilisateur.email,
                role_id: utilisateur.role_id,
                mot_de_passe: "",
                langue_preferee_id: utilisateur.langue_preferee_id,
            };
            this.erreur = "";
            this.formulaireOuvert = true;
            this.chargerPrediction(utilisateur);
        },

        async enregistrer() {
            const donnees = {
                nom: this.formulaire.nom,
                prenom: this.formulaire.prenom,
                email: this.formulaire.email,
                role_id: this.formulaire.role_id,
                langue_preferee_id: this.formulaire.langue_preferee_id,
            };

            let chemin = "/utilisateurs";
            let methode = "POST";
            if (this.formulaire.id) {
                chemin = `/utilisateurs/${this.formulaire.id}`;
                methode = "PUT";
            } else {
                donnees.mot_de_passe = this.formulaire.mot_de_passe;
            }

            const reponse = await appelerApi(chemin, { method: methode, body: JSON.stringify(donnees) });
            if (!reponse.ok) {
                const resultat = await reponse.json();
                this.erreur = resultat.erreur;
                return;
            }

            this.formulaireOuvert = false;
            this.chargerUtilisateurs();
        },

        async changerStatut(utilisateur) {
            const nouveauStatut = utilisateur.statut === "suspendu" ? "actif" : "suspendu";
            if (nouveauStatut === "suspendu" && !confirm(t("admin.utilisateurs.confirmer_suspension"))) {
                return;
            }

            const reponse = await appelerApi(`/utilisateurs/${utilisateur.id}/statut`, {
                method: "PUT",
                body: JSON.stringify({ statut: nouveauStatut }),
            });
            if (!reponse.ok) {
                this.erreurListe = (await reponse.json()).erreur;
                return;
            }
            this.chargerUtilisateurs();
        },

        async changerResponsable(utilisateur, estResponsable) {
            const reponse = await appelerApi(`/utilisateurs/${utilisateur.id}/responsable`, {
                method: "PUT",
                body: JSON.stringify({ est_responsable: estResponsable }),
            });
            if (!reponse.ok) {
                this.erreurListe = (await reponse.json()).erreur;
            }
            this.chargerUtilisateurs();
        },

        async supprimer(id) {
            if (!confirm(t("admin.utilisateurs.confirmer_suppression"))) {
                return;
            }

            const reponse = await appelerApi(`/utilisateurs/${id}`, { method: "DELETE" });
            if (!reponse.ok) {
                const resultat = await reponse.json();
                this.erreurListe = resultat.erreur;
                return;
            }

            this.chargerUtilisateurs();
        },
    },
});
