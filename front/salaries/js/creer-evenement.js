const { createApp } = Vue;

createApp({
    data() {
        return {
            formulaire: {
                titre: "",
                description: "",
                date_debut: "",
                date_fin: "",
                lieu: "",
                site: "",
            },
            confirmation: false,
            erreur: "",
        };
    },

    methods: {
        // statut = "brouillon" ou "en_attente" selon le bouton cliqué
        async envoyerEvenement(statut) {
            this.erreur = "";
            this.confirmation = false;

            if (!this.formulaire.titre || !this.formulaire.date_debut) {
                this.erreur = "Le titre et la date de début sont obligatoires.";
                return;
            }

            const reponse = await appelerApi("/evenements", {
                method: "POST",
                body: JSON.stringify({
                    titre: this.formulaire.titre,
                    description: this.formulaire.description,
                    date_debut: versDateAPI(this.formulaire.date_debut),
                    date_fin: versDateAPI(this.formulaire.date_fin),
                    lieu: this.formulaire.lieu,
                    site: this.formulaire.site,
                    statut: statut,
                }),
            });

            if (!reponse.ok) {
                const resultat = await reponse.json();
                this.erreur = resultat.erreur;
                return;
            }

            this.confirmation = true;
        },
    },
}).mount("#app");
