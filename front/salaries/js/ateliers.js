const { createApp } = Vue;

createApp({
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
                this.erreur = "Erreur lors du chargement des ateliers";
                return;
            }
            this.ateliers = await reponse.json();
        },

        modifierAtelier(id) {
            // Pas encore de page ni de route de modification : à faire.
            console.log("Modifier l'atelier :", id);
        },
    },
}).mount("#app");
