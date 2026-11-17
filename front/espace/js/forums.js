demarrerApp({
    data() {
        return {
            sujets: [],
            sujet: null,
            messages: [],
            nouveauSujet: { titre: "", contenu: "" },
            reponse: "",
            signalementOuvert: null,
            signalementVisible: false,
            nouveauSujetOuvert: false,
            reponseOuverte: false,
            motif: "",
            message: "",
            erreur: "",
        };
    },

    mounted() {
        this.chargerSujets();
    },

    methods: {
        formaterDate,
        formaterDateHeure,

        async chargerSujets() {
            const reponse = await appelerApi("/forums/sujets");
            if (!reponse.ok) {
                this.erreur = t("forums.erreur_chargement");
                return;
            }
            this.sujets = await reponse.json();
        },

        async ouvrirSujet(sujet) {
            this.sujet = sujet;
            this.message = "";
            this.erreur = "";
            const reponse = await appelerApi(`/forums/sujets/${sujet.id}/messages`);
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.messages = await reponse.json();
        },

        ouvrirNouveauSujet() {
            this.nouveauSujet = { titre: "", contenu: "" };
            this.erreur = "";
            this.nouveauSujetOuvert = true;
        },

        ouvrirReponse() {
            this.reponse = "";
            this.erreur = "";
            this.reponseOuverte = true;
        },

        ouvrirSignalement(messageID) {
            this.signalementOuvert = messageID;
            this.motif = "";
            this.erreur = "";
            this.message = "";
            this.signalementVisible = true;
        },

        async creerSujet() {
            this.erreur = "";
            const reponse = await appelerApi("/forums/sujets", {
                method: "POST",
                body: JSON.stringify(this.nouveauSujet),
            });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.nouveauSujetOuvert = false;
            this.chargerSujets();
        },

        async repondre() {
            this.erreur = "";
            const reponse = await appelerApi(`/forums/sujets/${this.sujet.id}/messages`, {
                method: "POST",
                body: JSON.stringify({ contenu: this.reponse }),
            });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.reponse = "";
            this.reponseOuverte = false;
            this.ouvrirSujet(this.sujet);
        },

        async signaler(messageID) {
            const reponse = await appelerApi(`/forums/messages/${messageID}/signalement`, {
                method: "POST",
                body: JSON.stringify({ motif: this.motif }),
            });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.signalementVisible = false;
            this.motif = "";
            this.message = t("forums.signalement_envoye");
        },
    },
});
