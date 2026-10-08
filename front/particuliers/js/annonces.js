const { createApp } = Vue;

createApp({
    data() {
        return {
            annonces: [],
            categories: [],
            filtres: { type: "", localisation: "" },
            selection: null,
            erreur: "",
        };
    },

    async mounted() {
        const reponse = await appelerApi("/categories");
        if (reponse.ok) {
            this.categories = await reponse.json();
        }
        this.chargerAnnonces();
    },

    methods: {
        formaterDate,
        formaterPrix,

        nomCategorie(id) {
            const categorie = this.categories.find((c) => c.id === id);
            return categorie ? categorie.code : "Non précisée";
        },

        async chargerAnnonces() {
            // Les filtres sont envoyés à l'API : c'est la base qui filtre, pas le navigateur
            const parametres = new URLSearchParams();
            if (this.filtres.type) {
                parametres.set("type", this.filtres.type);
            }
            if (this.filtres.localisation) {
                parametres.set("localisation", this.filtres.localisation);
            }

            const reponse = await appelerApi(`/annonces?${parametres}`);
            if (!reponse.ok) {
                this.erreur = "Erreur lors du chargement des annonces";
                return;
            }
            this.erreur = "";
            this.annonces = await reponse.json();
            this.selection = null;
        },
    },
}).mount("#app");
