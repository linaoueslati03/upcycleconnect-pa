demarrerApp({
    data() {
        return {
            depots: [],
            conteneurs: [],
            etapes: ETAPES_DEPOT,
            formulaire: { description_objet: "", conteneur_id: null },
            erreur: "",
            erreurListe: "",
            confirmation: false,

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
        };
    },

    async mounted() {
        const role = roleConnecte();
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

        async demanderDepot() {
            this.erreur = "";
            this.confirmation = false;

            const reponse = await appelerApi("/depots", {
                method: "POST",
                body: JSON.stringify(this.formulaire),
            });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }

            this.confirmation = true;
            this.formulaire = { description_objet: "", conteneur_id: null };
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
        },

        async avancer(depot) {
            const reponse = await appelerApi(`/depots/${depot.id}/statut`, {
                method: "PUT",
                body: JSON.stringify({ statut: this.statutSuivant[depot.statut] }),
            });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.chargerDepotsAValider();
        },

        async verifierCode() {
            this.erreur = "";
            this.succes = "";
            this.objetTrouve = null;

            let id = this.codeObjet.replace(/\D/g, "");
            if (!id) id = this.codeObjet;

            const reponse = await appelerApi(`/depots/${id}`);
            if (!reponse.ok) {
                this.erreur = "Objet introuvable ou code invalide.";
                return;
            }

            const depot = await reponse.json();

            if (depot.statut === "recupere") {
                this.erreur = "Cet objet a déjà été récupéré.";
                return;
            }

            this.objetTrouve = depot;
        },

        async confirmerRecuperation() {
            this.erreur = "";
            const reponse = await appelerApi(`/depots/${this.objetTrouve.id}/statut`, {
                method: "PUT",
                body: JSON.stringify({ statut: "recupere" })
            });

            if (!reponse.ok) {
                this.erreur = "Erreur lors de la récupération.";
                return;
            }

            this.succes = "Récupération confirmée avec succès ! Le conteneur est désormais vide.";

            setTimeout(() => {
                this.objetTrouve = null;
                this.codeObjet = "";
                this.succes = "";
            }, 3000);
        },
    },
});
