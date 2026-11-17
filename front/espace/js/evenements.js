demarrerApp({
    data() {
        return {
            evenements: [],
            filtres: FILTRES_STATUT,
            statutSelectionne: "tous",
            erreur: "",
            formulaireOuvert: false,
            editionId: null,
            formulaire: { titre: "", description: "", date_debut: "", date_fin: "", lieu: "", site: "" },
            erreurFormulaire: "",
            confirmation: false,
        };
    },

    computed: {
        evenementsFiltres() {
            if (this.statutSelectionne === "tous") {
                return this.evenements;
            }
            return this.evenements.filter(e => e.statut === this.statutSelectionne);
        },
    },

    mounted() {
        this.chargerEvenements();
        const aModifier = new URLSearchParams(window.location.search).get("modifier");
        if (aModifier) {
            this.modifierEvenement(aModifier);
        }
    },

    methods: {
        formaterDate,
        formaterStatut,

        async chargerEvenements() {
            const reponse = await appelerApi("/evenements");
            if (!reponse.ok) {
                this.erreur = t("salaries.evenements.erreur_chargement");
                return;
            }
            this.evenements = await reponse.json();
        },

        ouvrirCreation() {
            this.editionId = null;
            this.formulaire = { titre: "", description: "", date_debut: "", date_fin: "", lieu: "", site: "" };
            this.erreurFormulaire = "";
            this.confirmation = false;
            this.formulaireOuvert = true;
        },

        async modifierEvenement(id) {
            this.erreur = "";
            this.confirmation = false;
            const reponse = await appelerApi(`/evenements/${id}`);
            if (!reponse.ok) {
                this.erreur = t("salaries.evenements.introuvable");
                return;
            }
            const evenement = await reponse.json();
            this.editionId = evenement.id ?? id;
            this.formulaire = {
                titre: evenement.titre || "",
                description: evenement.description || "",
                lieu: evenement.lieu || "",
                site: evenement.site || "",
                date_debut: versChampDate(evenement.date_debut),
                date_fin: versChampDate(evenement.date_fin),
            };
            this.erreurFormulaire = "";
            this.formulaireOuvert = true;
        },

        async envoyerEvenement(statut) {
            this.erreurFormulaire = "";
            this.confirmation = false;

            if (!this.formulaire.titre || !this.formulaire.date_debut) {
                this.erreurFormulaire = t("salaries.titre_date_obligatoires");
                return;
            }

            const id = this.editionId;
            const reponse = await appelerApi(id ? `/evenements/${id}` : "/evenements", {
                method: id ? "PUT" : "POST",
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
                this.erreurFormulaire = (await reponse.json()).erreur;
                return;
            }

            this.formulaireOuvert = false;
            this.confirmation = true;
            this.chargerEvenements();
        },

        async valider(id) {
            const reponse = await appelerApi(`/evenements/${id}/valider`, { method: "POST" });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.chargerEvenements();
        },

        async supprimer(id) {
            if (!confirm(t("commun.confirmer_suppression"))) {
                return;
            }
            const reponse = await appelerApi(`/evenements/${id}`, { method: "DELETE" });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.chargerEvenements();
        },
    },
});
