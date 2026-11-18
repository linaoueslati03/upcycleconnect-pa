demarrerApp({
    data() {
        return {
            annonces: [],
            categories: [],
            filtres: { type: "", localisation: "" },
            selection: null,
            erreur: "",
            erreurAchat: "",
            idsALaUne: new Set(),
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
        const aLaUne = await appelerApi("/a-la-une");
        if (aLaUne.ok) {
            this.idsALaUne = new Set((await aLaUne.json()).annonces);
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
            const annonces = await reponse.json();
            // Les annonces mises en avant par une campagne publicitaire passent en premier
            annonces.forEach((a) => { a.a_la_une = this.idsALaUne.has(a.id); });
            this.annonces = [...annonces.filter((a) => a.a_la_une), ...annonces.filter((a) => !a.a_la_une)];
            this.selection = null;
        },
    },
});
