const { createApp } = Vue;

createApp({
    data() {
        return {
            email: "",
            motDePasse: "",
            erreur: "",
        };
    },

    methods: {
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
                this.erreur = "Votre espace n'est pas encore disponible.";
                return;
            }

            localStorage.setItem("token", resultat.token);
            localStorage.setItem("role", resultat.role);
            window.location.href = espace;
        },
    },
}).mount("#app");
