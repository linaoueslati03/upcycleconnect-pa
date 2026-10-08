const { createApp } = Vue;

createApp({
    data() {
        return {
            evenements: [],
            filtres: FILTRES_STATUT,
            statutSelectionne: "tous",
            erreur: "",
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
    },

    methods: {
        formaterDate,
        formaterStatut,

        async chargerEvenements() {
            const reponse = await fetch(`${API_BASE_URL}/evenements`);
            if (!reponse.ok) {
                this.erreur = "Erreur lors du chargement des événements";
                return;
            }
            this.evenements = await reponse.json();
        },

        modifierEvenement(id) {
            // Pas encore de page de modification : à faire (la route PUT existe déjà).
            console.log("Modifier l'événement :", id);
        },
    },
}).mount("#app");
