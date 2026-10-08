demarrerApp({
    data() {
        return {
            ateliers: [],
            filtres: FILTRES_STATUT,
            statutSelectionne: "tous",
            erreur: "",
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

        modifierAtelier(id) {
            window.location.href = `creer-atelier.php?id=${id}`;
        },

        // L'API vérifie que le salarié connecté est responsable (sinon 403, affiché ici)
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
