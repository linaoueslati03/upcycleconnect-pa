const { createApp } = Vue;

const FORMULAIRE_VIDE = { id: null, nom: "", prenom: "", email: "", role_id: 1, mot_de_passe: "" };

createApp({
    data() {
        return {
            nomsRoles: { 1: "Particulier", 2: "Professionnel", 3: "Salarié", 4: "Administrateur" },
            utilisateurs: [],
            formulaire: { ...FORMULAIRE_VIDE },
            erreur: "",
            erreurListe: "",
        };
    },

    mounted() {
        this.chargerUtilisateurs();
    },

    methods: {
        async chargerUtilisateurs() {
            const reponse = await appelerApi("/utilisateurs");
            if (!reponse.ok) {
                this.erreurListe = "Erreur lors du chargement des utilisateurs";
                return;
            }
            this.utilisateurs = await reponse.json();
        },

        reinitialiserFormulaire() {
            this.formulaire = { ...FORMULAIRE_VIDE };
            this.erreur = "";
        },

        modifier(utilisateur) {
            this.formulaire = {
                id: utilisateur.id,
                nom: utilisateur.nom,
                prenom: utilisateur.prenom,
                email: utilisateur.email,
                role_id: utilisateur.role_id,
                mot_de_passe: "",
            };
            this.erreur = "";
        },

        async enregistrer() {
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

            const reponse = await appelerApi(chemin, { method: methode, body: JSON.stringify(donnees) });
            if (!reponse.ok) {
                const resultat = await reponse.json();
                this.erreur = resultat.erreur;
                return;
            }

            this.reinitialiserFormulaire();
            this.chargerUtilisateurs();
        },

        async supprimer(id) {
            if (!confirm("Supprimer cet utilisateur ?")) {
                return;
            }

            const reponse = await appelerApi(`/utilisateurs/${id}`, { method: "DELETE" });
            if (!reponse.ok) {
                const resultat = await reponse.json();
                this.erreurListe = resultat.erreur;
                return;
            }

            this.chargerUtilisateurs();
        },
    },
}).mount("#app");
