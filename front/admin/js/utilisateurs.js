const { createApp } = Vue;

createApp({
    data() {
        return {
            nomsRoles: { 1: "Particulier", 2: "Professionnel", 3: "Salarié", 4: "Administrateur" },
            utilisateurs: [],
            formulaireVisible: false,
            erreur: "",
            formulaire: {
                id: null,
                nom: "",
                prenom: "",
                email: "",
                role_id: 1,
                mot_de_passe: "",
            },
        };
    },

    mounted() {
        this.chargerUtilisateurs();
    },

    methods: {
        deconnecter,

        async chargerUtilisateurs() {
            const reponse = await appelerApi("/utilisateurs");
            if (!reponse.ok) {
                this.erreur = "Erreur lors du chargement des utilisateurs";
                return;
            }
            this.utilisateurs = await reponse.json();
        },

        reinitialiserFormulaire() {
            this.formulaire = { id: null, nom: "", prenom: "", email: "", role_id: 1, mot_de_passe: "" };
            this.erreur = "";
        },

        ouvrirFormulaireCreation() {
            this.reinitialiserFormulaire();
            this.formulaireVisible = true;
        },

        ouvrirFormulaireEdition(utilisateur) {
            this.formulaire = {
                id: utilisateur.id,
                nom: utilisateur.nom,
                prenom: utilisateur.prenom,
                email: utilisateur.email,
                role_id: utilisateur.role_id,
                mot_de_passe: "",
            };
            this.erreur = "";
            this.formulaireVisible = true;
        },

        fermerFormulaire() {
            this.formulaireVisible = false;
        },

        async soumettreFormulaire() {
            const donnees = {
                nom: this.formulaire.nom,
                prenom: this.formulaire.prenom,
                email: this.formulaire.email,
                role_id: this.formulaire.role_id,
            };

            let chemin = "/utilisateurs";
            let methode = "POST";
            if (this.formulaire.id) {
                chemin = `/utilisateurs/${this.formulaire.id}`;
                methode = "PUT";
            } else {
                donnees.mot_de_passe = this.formulaire.mot_de_passe;
            }

            const reponse = await appelerApi(chemin, {
                method: methode,
                body: JSON.stringify(donnees),
            });

            const resultat = await reponse.json();

            if (!reponse.ok) {
                this.erreur = resultat.erreur;
                return;
            }

            this.formulaireVisible = false;
            this.chargerUtilisateurs();
        },

        async supprimerUtilisateur(id) {
            if (!confirm("Supprimer cet utilisateur ?")) {
                return;
            }
            await appelerApi(`/utilisateurs/${id}`, { method: "DELETE" });
            this.chargerUtilisateurs();
        },
    },
}).mount("#app");
