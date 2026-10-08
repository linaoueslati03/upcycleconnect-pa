const { createApp } = Vue;

createApp({
    data() {
        return {
            conseils: [],
            recherche: "",
            categorie: "",
            selection: null,
            erreur: "",
        };
    },

    computed: {
        // Catégories présentes dans les articles publiés (sans doublon)
        categories() {
            return [...new Set(this.conseils.map((c) => c.categorie).filter(Boolean))];
        },

        conseilsFiltres() {
            const motCle = this.recherche.toLowerCase();
            return this.conseils.filter((c) =>
                (!this.categorie || c.categorie === this.categorie) &&
                (c.titre.toLowerCase().includes(motCle) || c.contenu.toLowerCase().includes(motCle))
            );
        },
    },

    async mounted() {
        const reponse = await appelerApi("/conseils");
        if (!reponse.ok) {
            this.erreur = "Erreur lors du chargement des conseils";
            return;
        }
        this.conseils = await reponse.json();
    },

    methods: {
        formaterDate,
    },
}).mount("#app");
