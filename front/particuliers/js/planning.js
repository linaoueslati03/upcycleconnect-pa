const { createApp } = Vue;

createApp({
    data() {
        return {
            entrees: [],
            libellesOffre: LIBELLES_OFFRE,
            erreur: "",
        };
    },

    async mounted() {
        // L'API fusionne formations, ateliers et événements et les trie par date
        const reponse = await appelerApi("/moi/planning");
        if (!reponse.ok) {
            this.erreur = "Erreur lors du chargement du planning";
            return;
        }
        this.entrees = await reponse.json();
    },

    methods: {
        formaterDateHeure,
    },
}).mount("#app");
