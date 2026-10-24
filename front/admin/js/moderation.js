demarrerApp({
    data() {
        return {
            signalements: [],
            sujets: [],
            erreur: "",
            statutSignalements: "a_traiter",
            statutSujets: "en_attente",
            filtresSignalements: [
                { nom: "admin.moderation.a_traiter", valeur: "a_traiter" },
                { nom: "admin.moderation.traites", valeur: "traite" },
                { nom: "admin.moderation.rejetes", valeur: "rejete" },
                { nom: "commun.tous", valeur: "tous" },
            ],
            filtresSujets: [
                { nom: "admin.moderation.en_attente", valeur: "en_attente" },
                { nom: "admin.moderation.valides", valeur: "valide" },
                { nom: "admin.moderation.rejetes", valeur: "rejete" },
                { nom: "commun.tous", valeur: "tous" },
            ],
        };
    },

    mounted() {
        this.chargerSignalements();
        this.chargerSujets();
    },

    methods: {
        formaterDate(date) {
            return new Date(date).toLocaleDateString(langueCourante(), { day: "numeric", month: "long", year: "numeric" });
        },

        async chargerSignalements() {
            const reponse = await appelerApi(`/moderation/signalements?statut=${this.statutSignalements}`);
            if (!reponse.ok) {
                this.erreur = t("admin.moderation.erreur_chargement");
                return;
            }
            this.signalements = await reponse.json();
        },

        async chargerSujets() {
            const reponse = await appelerApi(`/moderation/sujets?statut=${this.statutSujets}`);
            if (!reponse.ok) {
                this.erreur = t("admin.moderation.erreur_chargement");
                return;
            }
            this.sujets = await reponse.json();
        },

        async decider(chemin, statut) {
            this.erreur = "";
            const reponse = await appelerApi(chemin, { method: "PUT", body: JSON.stringify({ statut }) });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return false;
            }
            return true;
        },

        async changerStatutMessage(messageID, statut) {
            if (await this.decider(`/forums/messages/${messageID}/statut`, statut)) {
                this.chargerSignalements();
            }
        },

        async rejeterSignalement(signalementID) {
            if (await this.decider(`/moderation/signalements/${signalementID}`, "rejete")) {
                this.chargerSignalements();
            }
        },

        async changerStatutSujet(sujetID, statut) {
            if (await this.decider(`/forums/sujets/${sujetID}/statut`, statut)) {
                this.chargerSujets();
            }
        },
    },
});
