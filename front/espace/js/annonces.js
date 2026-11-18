demarrerApp({
    data() {
        return {
            annonces: [],
            categories: [],
            filtres: { type: "", localisation: "" },
            selection: null,
            erreur: "",
            erreurAchat: "",
            achatEnCours: false,
        };
    },

    watch: {
        selection() {
            this.erreurAchat = "";
        },
        "filtres.type"() {
            this.chargerAnnonces();
        },
        "filtres.localisation"() {
            this.rechercherEnDirect();
        },
    },

    created() {
        this.rechercherEnDirect = differer(() => this.chargerAnnonces());
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

        estMonObjet(annonce) {
            return annonce.utilisateur_id === Number(localStorage.getItem("user_id"));
        },

        // L'objet est réservé le temps du paiement, puis cédé à l'acheteur une fois payé
        async acheter(annonce) {
            this.erreurAchat = "";
            this.achatEnCours = true;
            const erreur = await demarrerAchat("vente", annonce.id);
            if (erreur) {
                this.erreurAchat = erreur;
                this.achatEnCours = false;
                this.chargerAnnonces();
            }
        },

        nomCategorie(id) {
            const categorie = this.categories.find((c) => c.id === id);
            return categorie ? t("categorie." + categorie.code, categorie.libelle) : t("commun.non_precisee");
        },

        async chargerAnnonces() {
            const numero = (this.numeroRequete = (this.numeroRequete || 0) + 1);
            const parametres = new URLSearchParams();
            if (this.filtres.type) {
                parametres.set("type", this.filtres.type);
            }
            if (this.filtres.localisation) {
                parametres.set("localisation", this.filtres.localisation);
            }

            const reponse = await appelerApi(`/annonces?${parametres}`);
            if (numero !== this.numeroRequete) {
                return;
            }
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
