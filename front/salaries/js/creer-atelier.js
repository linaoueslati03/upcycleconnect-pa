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
        const reponse = await appelerApi(`/ateliers/${this.id}`);
        if (!reponse.ok) {
            this.erreur = t("salaries.ateliers.introuvable");
            return;
        }
        const atelier = await reponse.json();
        this.formulaire = {
            titre: atelier.titre || "",
            description: atelier.description || "",
            lieu: atelier.lieu || "",
            date_debut: versChampDate(atelier.date_debut),
            date_fin: versChampDate(atelier.date_fin),
        };
    },

    methods: {
        // statut = "brouillon" ou "en_attente" selon le bouton cliqué
        async envoyerAtelier(statut) {
            this.erreur = "";
            this.confirmation = false;

            if (!this.formulaire.titre || !this.formulaire.date_debut) {
                this.erreur = t("salaries.titre_date_obligatoires");
                return;
            }

            const reponse = await appelerApi(this.id ? `/ateliers/${this.id}` : "/ateliers", {
                method: this.id ? "PUT" : "POST",
                body: JSON.stringify({
                    titre: this.formulaire.titre,
                    description: this.formulaire.description,
                    date_debut: versDateAPI(this.formulaire.date_debut),
                    date_fin: versDateAPI(this.formulaire.date_fin),
                    lieu: this.formulaire.lieu,
                    statut: statut,
                }),
            });

            if (!reponse.ok) {
                const resultat = await reponse.json();
                this.erreur = resultat.erreur;
                return;
            }

            // Le formulaire passe en modification : un nouveau clic ne recrée pas l'élément
            this.id = (await reponse.json()).id;
            this.confirmation = true;
        },
    },
});
