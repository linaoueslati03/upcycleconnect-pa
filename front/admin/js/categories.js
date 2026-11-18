const EDITION_VIDE = { ouvert: false, type: null, categorie: null, libelle: "" };

demarrerApp({
    data() {
        return {
            groupes: {
                prestations: { liste: [], nouvelle: "", lecture: "/admin/categories/prestations", ecriture: "/admin/categories/prestations" },
                materiaux: { liste: [], nouvelle: "", lecture: "/categories", ecriture: "/admin/categories/materiaux" },
            },
            edition: { ...EDITION_VIDE },
            ajoutOuvert: false,
            ajoutType: "prestations",
            impact: { ouvert: false, categorie: null, facteur: null, source: "", erreur: "" },
            erreur: "",
        };
    },

    mounted() {
        this.charger("prestations");
        this.charger("materiaux");
    },

    methods: {
        async charger(type) {
            const groupe = this.groupes[type];
            const reponse = await appelerApi(groupe.lecture);
            if (!reponse.ok) {
                this.erreur = t("admin.categories.erreur_chargement");
                return;
            }
            groupe.liste = await reponse.json();
        },

        ouvrirAjout(type) {
            this.ajoutType = type;
            this.groupes[type].nouvelle = "";
            this.erreur = "";
            this.ajoutOuvert = true;
        },

        async ajouter(type) {
            this.erreur = "";
            const groupe = this.groupes[type];
            const reponse = await appelerApi(groupe.ecriture, {
                method: "POST",
                body: JSON.stringify({ libelle: groupe.nouvelle }),
            });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            groupe.nouvelle = "";
            this.ajoutOuvert = false;
            this.charger(type);
        },

        // Facteur de CO2 évité par kg réemployé (statistiques des professionnels), avec sa source
        ouvrirImpact(categorie) {
            this.impact = {
                ouvert: true, categorie, facteur: categorie.co2_evite_par_kg, source: categorie.source_facteur || "", erreur: "",
            };
        },

        async enregistrerImpact() {
            this.impact.erreur = "";
            const reponse = await appelerApi(`/admin/categories/materiaux/${this.impact.categorie.id}/impact`, {
                method: "PUT",
                body: JSON.stringify({ co2_evite_par_kg: this.impact.facteur, source_facteur: this.impact.source }),
            });
            if (!reponse.ok) {
                this.impact.erreur = (await reponse.json()).erreur;
                return;
            }
            this.impact.ouvert = false;
            this.charger("materiaux");
        },

        commencerEdition(type, categorie) {
            this.erreur = "";
            this.edition = { ouvert: true, type, categorie, libelle: categorie.libelle };
        },

        annulerEdition() {
            this.edition.ouvert = false;
        },

        async enregistrer(type, categorie, libelle, actif) {
            this.erreur = "";
            const reponse = await appelerApi(`${this.groupes[type].ecriture}/${categorie.id}`, {
                method: "PUT",
                body: JSON.stringify({ libelle, actif }),
            });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.annulerEdition();
            this.charger(type);
        },
    },
});
