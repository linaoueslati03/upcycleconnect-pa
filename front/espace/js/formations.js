demarrerApp({
    data() {
        return {
            formations: [],
            filtres: FILTRES_STATUT,
            statutSelectionne: "tous",
            erreur: "",
            chargement: true,
        };
    },

    computed: {
        formationsFiltres() {
            if (this.statutSelectionne === "tous") {
                return this.formations;
            }

            return this.formations.filter(
                f => f.statut === this.statutSelectionne
            );
        },
    },

    mounted() {
        this.chargerFormations();
    },

    methods: {
        formaterDate,
        formaterStatut,

        async chargerFormations() {
            this.chargement = true;
            this.erreur = "";

            try {
                const reponse = await appelerApi("/formations");

                if (!reponse.ok) {
                    this.erreur = t("salaries.formations.erreur_chargement");
                    return;
                }

                this.formations = await reponse.json();
            } catch (e) {
                this.erreur = t("salaries.formations.erreur_chargement");
            } finally {
                this.chargement = false;
            }
        },

        modifierFormation(id) {
            window.location.href = `creer-formation.php?id=${id}`;
        },

        async valider(id) {
            this.erreur = "";

            try {
                const reponse = await appelerApi(`/formations/${id}/valider`, {
                    method: "POST",
                });

                if (!reponse.ok) {
                    this.erreur = (await reponse.json()).erreur;
                    return;
                }

                await this.chargerFormations();
            } catch (e) {
                this.erreur = t("salaries.formations.erreur_chargement");
            }
        },

        async supprimer(id) {
            if (!confirm(t("commun.confirmer_suppression"))) {
                return;
            }

            this.erreur = "";

            try {
                const reponse = await appelerApi(`/formations/${id}`, {
                    method: "DELETE",
                });

                if (!reponse.ok) {
                    this.erreur = (await reponse.json()).erreur;
                    return;
                }

                await this.chargerFormations();
            } catch (e) {
                this.erreur = t("salaries.formations.erreur_chargement");
            }
        },
    },
});