const { createApp } = Vue;

createApp({
    data() {
        return {
            profil: { nom: "", prenom: "", email: "", langue_preferee_id: null },
            langues: [],
            motsDePasse: { ancien_mot_de_passe: "", nouveau_mot_de_passe: "" },
            erreurProfil: "",
            messageProfil: "",
            erreurMotDePasse: "",
            messageMotDePasse: "",
        };
    },

    async mounted() {
        // Les langues proposées viennent de la base : en ajouter une ne demande aucun changement de code
        const [reponseCompte, reponseLangues] = await Promise.all([appelerApi("/moi"), appelerApi("/langues")]);
        if (!reponseCompte.ok || !reponseLangues.ok) {
            this.erreurProfil = "Erreur lors du chargement de votre compte";
            return;
        }
        const compte = await reponseCompte.json();
        this.profil = {
            nom: compte.nom,
            prenom: compte.prenom,
            email: compte.email,
            langue_preferee_id: compte.langue_preferee_id,
        };
        this.langues = await reponseLangues.json();
    },

    methods: {
        async enregistrerProfil() {
            this.erreurProfil = "";
            this.messageProfil = "";
            const reponse = await appelerApi("/moi", { method: "PUT", body: JSON.stringify(this.profil) });
            if (!reponse.ok) {
                this.erreurProfil = await messageErreur(reponse);
                return;
            }
            this.messageProfil = "Profil enregistré.";
        },

        async changerMotDePasse() {
            this.erreurMotDePasse = "";
            this.messageMotDePasse = "";
            const reponse = await appelerApi("/moi/mot-de-passe", {
                method: "PUT",
                body: JSON.stringify(this.motsDePasse),
            });
            if (!reponse.ok) {
                this.erreurMotDePasse = await messageErreur(reponse);
                return;
            }
            this.motsDePasse = { ancien_mot_de_passe: "", nouveau_mot_de_passe: "" };
            this.messageMotDePasse = "Mot de passe modifié.";
        },
    },
}).mount("#app");
