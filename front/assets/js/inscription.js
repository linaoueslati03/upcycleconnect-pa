demarrerApp({
    data() {
        return {
            compte: { role: "particulier", nom: "", prenom: "", email: "", mot_de_passe: "" },
            erreur: "",
        };
    },

    methods: {
        async creerCompte() {
            this.erreur = "";
            oublierConnexion();

            const reponse = await appelerApi("/comptes", { method: "POST", body: JSON.stringify(this.compte) });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }

            // Compte créé : on connecte directement l'utilisateur
            const connexion = await appelerApi("/login", {
                method: "POST",
                body: JSON.stringify({ email: this.compte.email, mot_de_passe: this.compte.mot_de_passe }),
            });
            const resultat = await connexion.json();
            const espace = ESPACES[resultat.role];
            if (!connexion.ok || !espace) {
                window.location.href = PAGE_CONNEXION;
                return;
            }

            localStorage.setItem("token", resultat.token);
            localStorage.setItem("role", resultat.role);
            localStorage.setItem("langue", resultat.langue);
            window.location.href = espace;
        },
    },
});
