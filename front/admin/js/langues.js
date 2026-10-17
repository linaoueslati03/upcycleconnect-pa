demarrerApp({
    data() {
        return {
            langues: [],
            fichiers: {},
            langueParDefaut: LANGUE_PAR_DEFAUT,
            formulaire: { code: "", libelle: "" },
            erreur: "",
            succes: "",
            erreurFormulaire: "",
        };
    },

    mounted() {
        this.chargerLangues();
    },

    methods: {
        async chargerLangues() {
            const reponse = await appelerApi("/admin/langues");
            if (!reponse.ok) {
                this.erreur = t("admin.langues.erreur_chargement");
                return;
            }
            this.langues = await reponse.json();
            this.verifierFichiers();
        },

        async verifierFichiers() {
            for (const langue of this.langues) {
                try {
                    const reponse = await fetch(`${LOCALES_URL}/${langue.code}.json`, { method: "HEAD", cache: "no-store" });
                    this.fichiers[langue.code] = reponse.ok;
                } catch (erreur) {
                    this.fichiers[langue.code] = false;
                }
            }
        },

        async ajouter() {
            this.erreurFormulaire = "";
            const reponse = await appelerApi("/admin/langues", {
                method: "POST",
                body: JSON.stringify(this.formulaire),
            });
            if (!reponse.ok) {
                this.erreurFormulaire = (await reponse.json()).erreur;
                return;
            }
            this.formulaire = { code: "", libelle: "" };
            this.chargerLangues();
        },

        async importer(langue, evenement) {
            this.erreur = "";
            this.succes = "";
            const fichier = evenement.target.files[0];
            evenement.target.value = "";
            if (!fichier) {
                return;
            }

            const reponse = await appelerApi(`/admin/langues/${langue.id}/traduction`, {
                method: "PUT",
                body: await fichier.text(),
            });
            const resultat = await reponse.json();
            if (!reponse.ok) {
                this.erreur = resultat.erreur;
                return;
            }
            this.succes = t("admin.langues.import_reussi", { langue: langue.libelle, nombre: resultat.textes });
            this.verifierFichiers();
        },

        async basculer(langue) {
            this.erreur = "";
            const reponse = await appelerApi(`/admin/langues/${langue.id}`, {
                method: "PUT",
                body: JSON.stringify({ libelle: langue.libelle, actif: !langue.actif }),
            });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.chargerLangues();
        },
    },
});
