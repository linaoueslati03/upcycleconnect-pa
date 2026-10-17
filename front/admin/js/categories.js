const EDITION_VIDE = { type: null, id: null, libelle: "" };

demarrerApp({
    data() {
        return {
            groupes: {
                prestations: { liste: [], nouvelle: "", lecture: "/admin/categories/prestations", ecriture: "/admin/categories/prestations" },
                materiaux: { liste: [], nouvelle: "", lecture: "/categories", ecriture: "/admin/categories/materiaux" },
            },
            edition: { ...EDITION_VIDE },
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
            this.charger(type);
        },

        commencerEdition(type, categorie) {
            this.edition = { type, id: categorie.id, libelle: categorie.libelle };
        },

        annulerEdition() {
            this.edition = { ...EDITION_VIDE };
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
