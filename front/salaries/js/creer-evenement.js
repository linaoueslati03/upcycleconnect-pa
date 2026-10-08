demarrerApp({
    data() {
        return {
            id: new URLSearchParams(window.location.search).get("id"),
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

    // En modification, on pré-remplit le formulaire avec les valeurs actuelles
    async mounted() {
        if (!this.id) {
            return;
        }
        const reponse = await appelerApi(`/evenements/${this.id}`);
        if (!reponse.ok) {
            this.erreur = t("salaries.evenements.introuvable");
            return;
        }
        const evenement = await reponse.json();
        this.formulaire = {
            titre: evenement.titre || "",
            description: evenement.description || "",
            lieu: evenement.lieu || "",
            site: evenement.site || "",
            date_debut: versChampDate(evenement.date_debut),
            date_fin: versChampDate(evenement.date_fin),
        };
    },

    methods: {
        // statut = "brouillon" ou "en_attente" selon le bouton cliqué
        async envoyerEvenement(statut) {
            this.erreur = "";
            this.confirmation = false;

            if (!this.formulaire.titre || !this.formulaire.date_debut) {
                this.erreur = t("salaries.titre_date_obligatoires");
                return;
            }

            const reponse = await appelerApi(this.id ? `/evenements/${this.id}` : "/evenements", {
                method: this.id ? "PUT" : "POST",
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
});
