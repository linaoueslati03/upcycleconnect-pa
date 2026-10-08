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
            if (!connexion.ok) {
                window.location.href = PAGE_CONNEXION;
                return;
            }
            const espace = ESPACES[resultat.role];
            if (!espace) {
                // Compte créé, mais l'espace de ce rôle (professionnel) n'existe pas encore
                this.erreur = t("inscription.compte_cree_espace_indisponible");
                return;
            }

            localStorage.setItem("token", resultat.token);
            localStorage.setItem("role", resultat.role);
            // Langue enregistrée dans le compte ; sinon on garde celle choisie à l'écran
            if (resultat.langue) {
                localStorage.setItem("langue", resultat.langue);
            }
            window.location.href = espace;
        },
    },
});
