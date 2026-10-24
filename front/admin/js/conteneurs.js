const STATUTS_CONTENEUR = ["vide", "occupe", "maintenance"];

demarrerApp({
    data() {
        return {
            conteneurs: [],
            statuts: STATUTS_CONTENEUR,
            filtreStatut: "",
            filtres: [
                { nom: "commun.tous", valeur: "" },
                { nom: "admin.conteneurs.statut.vide", valeur: "vide" },
                { nom: "admin.conteneurs.statut.occupe", valeur: "occupe" },
                { nom: "admin.conteneurs.statut.maintenance", valeur: "maintenance" },
            ],
            nouveauSite: "",
            changement: { conteneur: null, statut: "", motif: "" },
            historique: { conteneur: null, lignes: [] },
            erreur: "",
        };
    },

    computed: {
        conteneursFiltres() {
            if (!this.filtreStatut) {
                return this.conteneurs;
            }
            return this.conteneurs.filter((c) => c.statut === this.filtreStatut);
        },
    },

    mounted() {
        this.charger();
    },

    methods: {
        formaterDateHeure(date) {
            return new Date(date).toLocaleString(langueCourante(), {
                day: "numeric", month: "long", year: "numeric", hour: "2-digit", minute: "2-digit",
            });
        },

        compter(statut) {
            if (!statut) {
                return this.conteneurs.length;
            }
            return this.conteneurs.filter((c) => c.statut === statut).length;
        },

        async charger() {
            const reponse = await appelerApi("/admin/conteneurs");
            if (!reponse.ok) {
                this.erreur = t("admin.conteneurs.erreur_chargement");
                return;
            }
            this.conteneurs = await reponse.json();
        },

        async ajouter() {
            this.erreur = "";
            const reponse = await appelerApi("/admin/conteneurs", {
                method: "POST",
                body: JSON.stringify({ site: this.nouveauSite }),
            });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.nouveauSite = "";
            this.charger();
        },

        ouvrirChangement(conteneur) {
            const autreStatut = STATUTS_CONTENEUR.find((s) => s !== conteneur.statut);
            this.changement = { conteneur, statut: autreStatut, motif: "" };
        },

        fermerChangement() {
            this.changement = { conteneur: null, statut: "", motif: "" };
        },

        async enregistrerStatut() {
            this.erreur = "";
            const conteneur = this.changement.conteneur;
            const reponse = await appelerApi(`/admin/conteneurs/${conteneur.id}/statut`, {
                method: "PUT",
                body: JSON.stringify({ statut: this.changement.statut, motif: this.changement.motif }),
            });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.fermerChangement();
            await this.charger();
            if (this.historique.conteneur && this.historique.conteneur.id === conteneur.id) {
                this.afficherHistorique(conteneur);
            }
        },

        async afficherHistorique(conteneur) {
            this.erreur = "";
            const reponse = await appelerApi(`/admin/conteneurs/${conteneur.id}/historique`);
            if (!reponse.ok) {
                this.erreur = t("admin.conteneurs.erreur_chargement");
                return;
            }
            this.historique = { conteneur, lignes: await reponse.json() };
        },
    },
});
