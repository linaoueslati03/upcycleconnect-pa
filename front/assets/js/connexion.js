demarrerApp({
    data() {
        return {
            email: "",
            motDePasse: "",
            erreur: "",
            langues: [],
        };
    },

    async mounted() {
        const reponse = await appelerApi("/langues");
        if (reponse.ok) {
            this.langues = await reponse.json();
        }
    },

    methods: {
        changerLangue,

        async seConnecter() {
            this.erreur = "";
            oublierConnexion();

            let reponse;
            try {
                reponse = await appelerApi("/login", {
                    method: "POST",
                    body: JSON.stringify({ email: this.email, mot_de_passe: this.motDePasse }),
                });
            } catch (erreur) {
                this.erreur = t("commun.serveur_injoignable");
                return;
            }
            const resultat = await reponse.json();

            if (!reponse.ok) {
                this.erreur = resultat.erreur;
                return;
            }

            const espace = ESPACES[resultat.role];
            if (!espace) {
                this.erreur = t("connexion.espace_indisponible");
                return;
            }

            enregistrerConnexion(resultat);
            if (resultat.doit_changer_mot_de_passe) {
                window.location.href = pageMonCompte();
                return;
            }
            window.location.href = espace;
        },
    },
});
