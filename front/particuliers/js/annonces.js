demarrerApp({
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
            return categorie ? traduireOu("categorie." + categorie.code, categorie.code) : t("commun.non_precisee");
        },

        async chargerAnnonces() {
            const parametres = new URLSearchParams();
            if (this.filtres.type) {
                parametres.set("type", this.filtres.type);
            }
            if (this.filtres.localisation) {
                parametres.set("localisation", this.filtres.localisation);
            }

            const reponse = await appelerApi(`/annonces?${parametres}`);
            if (!reponse.ok) {
                this.erreur = t("annonces.erreur_chargement");
                return;
            }
            this.erreur = "";
            this.annonces = await reponse.json();
            this.selection = null;
        },
    },
});
