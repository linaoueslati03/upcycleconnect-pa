demarrerApp({
    data() {
        return {
            ateliers: [],
            filtres: FILTRES_STATUT,
            statutSelectionne: "tous",
            erreur: "",
            formulaireOuvert: false,
            editionId: null,
            formulaire: { titre: "", description: "", date_debut: "", date_fin: "", lieu: "" },
            erreurFormulaire: "",
            confirmation: false,
        };
    },

    computed: {
        ateliersFiltres() {
            if (this.statutSelectionne === "tous") {
                return this.ateliers;
            }
            return this.ateliers.filter(a => a.statut === this.statutSelectionne);
        },
    },

    mounted() {
        this.chargerAteliers();
        const aModifier = new URLSearchParams(window.location.search).get("modifier");
        if (aModifier) {
            this.modifierAtelier(aModifier);
        }
    },

    methods: {
        formaterDate,
        formaterStatut,

        async chargerAteliers() {
            const reponse = await appelerApi("/ateliers");
            if (!reponse.ok) {
                this.erreur = t("salaries.ateliers.erreur_chargement");
                return;
            }
            this.ateliers = await reponse.json();
        },

        ouvrirCreation() {
            this.editionId = null;
            this.formulaire = { titre: "", description: "", date_debut: "", date_fin: "", lieu: "" };
            this.erreurFormulaire = "";
            this.confirmation = false;
            this.formulaireOuvert = true;
        },

        async modifierAtelier(id) {
            this.erreur = "";
            this.confirmation = false;
            const reponse = await appelerApi(`/ateliers/${id}`);
            if (!reponse.ok) {
                this.erreur = t("salaries.ateliers.introuvable");
                return;
            }
            const atelier = await reponse.json();
            this.editionId = atelier.id ?? id;
            this.formulaire = {
                titre: atelier.titre || "",
                description: atelier.description || "",
                lieu: atelier.lieu || "",
                date_debut: versChampDate(atelier.date_debut),
                date_fin: versChampDate(atelier.date_fin),
            };
            this.erreurFormulaire = "";
            this.formulaireOuvert = true;
        },

        async envoyerAtelier(statut) {
            this.erreurFormulaire = "";
            this.confirmation = false;

            if (!this.formulaire.titre || !this.formulaire.date_debut) {
                this.erreurFormulaire = t("salaries.titre_date_obligatoires");
                return;
            }

            const id = this.editionId;
            const reponse = await appelerApi(id ? `/ateliers/${id}` : "/ateliers", {
                method: id ? "PUT" : "POST",
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
                this.erreurFormulaire = (await reponse.json()).erreur;
                return;
            }

            this.formulaireOuvert = false;
            this.confirmation = true;
            this.chargerAteliers();
        },

        async valider(id) {
            const reponse = await appelerApi(`/ateliers/${id}/valider`, { method: "POST" });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.chargerAteliers();
        },

        async supprimer(id) {
            if (!confirm(t("commun.confirmer_suppression"))) {
                return;
            }
            const reponse = await appelerApi(`/ateliers/${id}`, { method: "DELETE" });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.chargerAteliers();
        },
    },
});
