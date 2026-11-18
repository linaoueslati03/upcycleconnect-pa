demarrerApp({
    data() {
        return {
            prenom: "",
            resume: null,
            abonnement: "",
            campagnesActives: 0,
            projetsALaUne: [],
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

        // Projets de professionnels mis en avant par une campagne publicitaire active
        async chargerALaUne() {
            const reponse = await appelerApi("/a-la-une");
            if (reponse.ok) {
                this.projetsALaUne = (await reponse.json()).projets;
            }
        },

        async charger() {
            const reponseCompte = await appelerApi("/moi");
            if (!reponseCompte.ok) {
                this.erreur = t("accueil.erreur_chargement");
                return;
            }
            const compte = await reponseCompte.json();
            this.prenom = compte.prenom;
            this.abonnement = compte.abonnement;
            this.chargerALaUne();

            if (roleConnecte() === "professionnel") {
                const campagnes = await appelerApi("/pro/campagnes");
                if (campagnes.ok) {
                    this.campagnesActives = (await campagnes.json()).filter((c) => c.statut === "active").length;
                }
                return;
            }
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
