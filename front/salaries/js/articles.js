const ARTICLE_VIDE = { id: null, titre: "", categorie: "", contenu: "" };

demarrerApp({
    data() {
        return {
            articles: [],
            formulaireVisible: false,
            formulaire: { ...ARTICLE_VIDE },
            erreur: "",
        };
    },

    mounted() {
        this.chargerArticles();
    },

    methods: {
        formaterDate,
        formaterStatut,

        // Articles du salarié connecté, brouillons compris
        async chargerArticles() {
            const reponse = await appelerApi("/salaries/conseils");
            if (!reponse.ok) {
                this.erreur = t("salaries.articles.erreur_chargement");
                return;
            }
            this.articles = await reponse.json();
        },

        nouvelArticle() {
            this.formulaire = { ...ARTICLE_VIDE };
            this.formulaireVisible = true;
            this.erreur = "";
        },

        modifier(article) {
            this.formulaire = {
                id: article.id,
                titre: article.titre,
                categorie: article.categorie || "",
                contenu: article.contenu,
            };
            this.formulaireVisible = true;
            this.erreur = "";
        },

        // statut = "brouillon" ou "publie" selon le bouton cliqué
        async enregistrer(statut) {
            this.erreur = "";
            const id = this.formulaire.id;
            const reponse = await appelerApi(id ? `/conseils/${id}` : "/conseils", {
                method: id ? "PUT" : "POST",
                body: JSON.stringify({
                    titre: this.formulaire.titre,
                    contenu: this.formulaire.contenu,
                    categorie: this.formulaire.categorie || null,
                    statut: statut,
                }),
            });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }

            this.formulaireVisible = false;
            this.chargerArticles();
        },
    },
});
