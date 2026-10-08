demarrerApp({
    data() {
        return {
            email: "",
            motDePasse: "",
            erreur: "",
            langues: [],
        };
    },

    // Langues proposées : celles actives en base (ajouter une langue = une ligne SQL)
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

            const reponse = await appelerApi("/login", {
                method: "POST",
                body: JSON.stringify({ email: this.email, mot_de_passe: this.motDePasse }),
            });
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
            localStorage.setItem("langue", resultat.langue);
            window.location.href = espace;
        },
    },
});
