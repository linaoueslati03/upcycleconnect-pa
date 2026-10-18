const ARTICLE_VIDE = { id: null, titre: "", categorie: "", contenu: "" };

demarrerApp({
    data() {
        return {
            conseils: [],
            recherche: "",
            categorie: "",
            selection: null,
            erreur: "",

            articles: [],
            formulaireVisible: false,
            formulaire: { ...ARTICLE_VIDE },
            erreurArticle: "",
        };
    },

    computed: {
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

    mounted() {
        this.chargerConseils();
        if (roleConnecte() === "salarie") {
            this.chargerArticles();
        }
    },

    methods: {
        formaterDate,
        formaterStatut,

        async chargerConseils() {
            const reponse = await appelerApi("/conseils");
            if (!reponse.ok) {
                this.erreur = t("conseils.erreur_chargement");
                return;
            }
            this.conseils = await reponse.json();
        },

        async chargerArticles() {
            const reponse = await appelerApi("/salaries/conseils");
            if (!reponse.ok) {
                this.erreurArticle = t("salaries.articles.erreur_chargement");
                return;
            }
            this.articles = await reponse.json();
        },

        nouvelArticle() {
            this.formulaire = { ...ARTICLE_VIDE };
            this.formulaireVisible = true;
            this.erreurArticle = "";
        },

        modifier(article) {
            this.formulaire = {
                id: article.id,
                titre: article.titre,
                categorie: article.categorie || "",
                contenu: article.contenu,
            };
            this.formulaireVisible = true;
            this.erreurArticle = "";
        },

        async enregistrer(statut) {
            this.erreurArticle = "";
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
                this.erreurArticle = (await reponse.json()).erreur;
                return;
            }

            this.formulaireVisible = false;
            this.chargerArticles();
            this.chargerConseils();
        },
    },
});
