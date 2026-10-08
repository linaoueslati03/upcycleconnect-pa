const { createApp } = Vue;

createApp({
    data() {
        return {
            formulaireVisible: false,
            titre: "",
            categorie: "",
            texte: "",
        };
    },

    methods: {
        publier() {
            // Pas encore de route POST /api/conseils côté API : à brancher quand elle existera.
            console.log({
                titre: this.titre,
                categorie: this.categorie,
                texte: this.texte,
            });
        },
    },
}).mount("#app");
