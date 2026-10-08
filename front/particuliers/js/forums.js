const { createApp } = Vue;

createApp({
    data() {
        return {
            sujets: [],
            sujet: null, // sujet ouvert
            messages: [],
            nouveauSujet: { titre: "", contenu: "" },
            reponse: "",
            signalementOuvert: null, // id du message dont le formulaire de signalement est ouvert
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
                this.erreur = "Erreur lors du chargement des sujets";
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

        async creerSujet() {
            this.erreur = "";
            const reponse = await appelerApi("/forums/sujets", {
                method: "POST",
                body: JSON.stringify(this.nouveauSujet),
            });
            // 429 : l'API limite la fréquence des publications (anti-spam)
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.nouveauSujet = { titre: "", contenu: "" };
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
            this.signalementOuvert = null;
            this.motif = "";
            this.message = "Merci, le message a été signalé à l'équipe de modération.";
        },
    },
}).mount("#app");
