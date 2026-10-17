const FORMULAIRE_VIDE = { id: null, titre: "", categorie_id: null, tarif: 0, statut: "brouillon" };

demarrerApp({
    data() {
        return {
            prestations: [],
            categories: [],
            categorieInitiale: null,
            formulaire: { ...FORMULAIRE_VIDE },
            erreur: "",
            erreurListe: "",
        };
    },

    computed: {
        categoriesProposees() {
            return this.categories.filter((c) => c.actif || c.id === this.categorieInitiale);
        },
    },

    mounted() {
        this.chargerPrestations();
        this.chargerCategories();
    },

    methods: {
        async chargerPrestations() {
            const reponse = await appelerApi("/prestations");
            if (!reponse.ok) {
                this.erreurListe = t("admin.prestations.erreur_chargement");
                return;
            }
            this.prestations = await reponse.json();
        },

        async chargerCategories() {
            const reponse = await appelerApi("/admin/categories/prestations");
            if (reponse.ok) {
                this.categories = await reponse.json();
            }
        },

        formaterTarif(tarif) {
            return tarif.toLocaleString(langueCourante(), { style: "currency", currency: "EUR" });
        },

        lienPDF(id) {
            return `${API_BASE_URL}/prestations/${id}/pdf`;
        },

        reinitialiserFormulaire() {
            this.formulaire = { ...FORMULAIRE_VIDE };
            this.categorieInitiale = null;
            this.erreur = "";
        },

        modifier(prestation) {
            this.formulaire = {
                id: prestation.id,
                titre: prestation.titre,
                categorie_id: prestation.categorie_id,
                tarif: prestation.tarif,
                statut: prestation.statut,
            };
            this.categorieInitiale = prestation.categorie_id;
            this.erreur = "";
        },

        async enregistrer() {
            const { id, ...donnees } = this.formulaire;
            const chemin = id ? `/prestations/${id}` : "/prestations";
            const methode = id ? "PUT" : "POST";

            const reponse = await appelerApi(chemin, { method: methode, body: JSON.stringify(donnees) });
            if (!reponse.ok) {
                const resultat = await reponse.json();
                this.erreur = resultat.erreur;
                return;
            }

            this.reinitialiserFormulaire();
            this.chargerPrestations();
        },

        async supprimer(id) {
            if (!confirm(t("admin.prestations.confirmer_suppression"))) {
                return;
            }

            const reponse = await appelerApi(`/prestations/${id}`, { method: "DELETE" });
            if (!reponse.ok) {
                const resultat = await reponse.json();
                this.erreurListe = resultat.erreur;
                return;
            }

            this.chargerPrestations();
        },
    },
});
