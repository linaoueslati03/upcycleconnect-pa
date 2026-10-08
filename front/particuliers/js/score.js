const { createApp } = Vue;

createApp({
    data() {
        return {
            score: 0,
            historique: [],
            bareme: [],
            erreur: "",
        };
    },

    async mounted() {
        const [reponseScore, reponseHistorique, reponseBareme] = await Promise.all([
            appelerApi("/moi/score"),
            appelerApi("/moi/score/historique"),
            appelerApi("/score/bareme"),
        ]);
        if (!reponseScore.ok || !reponseHistorique.ok || !reponseBareme.ok) {
            this.erreur = "Erreur lors du chargement du score";
            return;
        }
        this.score = (await reponseScore.json()).upcycling_score;
        this.historique = await reponseHistorique.json();
        this.bareme = await reponseBareme.json();
    },

    methods: {
        formaterDate,
    },
}).mount("#app");
