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

            localStorage.setItem("token", resultat.token);
            localStorage.setItem("role", resultat.role);
            if (resultat.langue) {
                localStorage.setItem("langue", resultat.langue);
            }
            if (resultat.doit_changer_mot_de_passe) {
                localStorage.setItem("changer_mot_de_passe", "1");
                window.location.href = pageMonCompte();
                return;
            }
            window.location.href = espace;
        },
    },
});
