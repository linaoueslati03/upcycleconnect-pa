demarrerApp({
    data() {
        return {
            prenom: "",
            resume: null,
            erreur: "",
        };
    },

    mounted() {
        this.charger();
    },

    methods: {
        formaterDate,

        libelleDepot(statut) {
            return statut ? t("depot." + statut) : t("accueil.aucun_depot");
        },

        async charger() {
            const reponseCompte = await appelerApi("/moi");
            if (!reponseCompte.ok) {
                this.erreur = t("accueil.erreur_chargement");
                return;
            }
            const compte = await reponseCompte.json();
            this.prenom = compte.prenom;

            if (roleConnecte() !== "particulier") {
                return;
            }

            const reponseResume = await appelerApi("/moi/dashboard");
            if (!reponseResume.ok) {
                this.erreur = t("accueil.erreur_chargement");
                return;
            }
            this.resume = await reponseResume.json();
        },
    },
});
