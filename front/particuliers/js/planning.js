demarrerApp({
    data() {
        return {
            entrees: [],
            erreur: "",
        };
    },

    async mounted() {
        // L'API fusionne formations, ateliers et événements et les trie par date
        const reponse = await appelerApi("/moi/planning");
        if (!reponse.ok) {
            this.erreur = t("salaries.planning.erreur_chargement");
            return;
        }
        this.entrees = await reponse.json();
    },

    methods: {
        formaterDateHeure,
    },
});
