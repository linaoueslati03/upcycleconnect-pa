function formationVide() {
    return {
        titre: "",
        description: "",
        date_debut: "",
        date_fin: "",
        lieu: "",
        nb_places: null,
        tarif: 0,
        materiel_necessaire: "",
        statut: "brouillon",
    };
}

demarrerApp({
    data() {
        return {
            formations: [],
            filtres: FILTRES_STATUT,
            statutSelectionne: "tous",
            erreur: "",
            chargement: true,
            formulaireOuvert: false,
            editionId: null,
            formation: formationVide(),
            erreurFormulaire: "",
            succes: "",
            enregistrement: false,
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
        const aModifier = new URLSearchParams(window.location.search).get("modifier");
        if (aModifier) {
            this.modifierFormation(aModifier);
        }
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

        convertirDate(date) {
            if (!date) {
                return "";
            }

            const d = new Date(date);
            const decalage = d.getTimezoneOffset() * 60000;

            return new Date(d.getTime() - decalage)
                .toISOString()
                .slice(0, 16);
        },

        ouvrirCreation() {
            this.editionId = null;
            this.formation = formationVide();
            this.erreurFormulaire = "";
            this.succes = "";
            this.formulaireOuvert = true;
        },

        async modifierFormation(id) {
            this.erreur = "";
            this.succes = "";

            try {
                const reponse = await appelerApi(`/formations/${id}`);

                if (!reponse.ok) {
                    this.erreur = (await reponse.json()).erreur;
                    return;
                }

                const f = await reponse.json();

                this.editionId = id;
                this.formation = {
                    titre: f.titre ?? "",
                    description: f.description ?? "",
                    date_debut: this.convertirDate(f.date_debut),
                    date_fin: this.convertirDate(f.date_fin),
                    lieu: f.lieu ?? "",
                    nb_places: f.nb_places ?? null,
                    tarif: Number(f.tarif ?? 0),
                    materiel_necessaire: f.materiel_necessaire ?? "",
                    statut: f.statut ?? "brouillon",
                };
                this.erreurFormulaire = "";
                this.formulaireOuvert = true;
            } catch (e) {
                this.erreur = t("salaries.formations.erreur_chargement");
            }
        },

        async enregistrer() {
            this.erreurFormulaire = "";
            this.succes = "";

            if (
                this.formation.date_fin &&
                this.formation.date_fin < this.formation.date_debut
            ) {
                this.erreurFormulaire = t("salaries.formations.date_fin_invalide");
                return;
            }

            this.enregistrement = true;

            const donnees = {
                ...this.formation,
                date_debut: new Date(this.formation.date_debut).toISOString(),
                date_fin: this.formation.date_fin
                    ? new Date(this.formation.date_fin).toISOString()
                    : null,
                nb_places: this.formation.nb_places || null,
                tarif: Number(this.formation.tarif || 0),
            };

            try {
                const reponse = await appelerApi(
                    this.editionId
                        ? `/formations/${this.editionId}`
                        : "/formations",
                    {
                        method: this.editionId ? "PUT" : "POST",
                        body: JSON.stringify(donnees),
                    }
                );

                if (!reponse.ok) {
                    this.erreurFormulaire = (await reponse.json()).erreur
                        || t("salaries.formations.erreur_enregistrement");
                    return;
                }

                this.formulaireOuvert = false;
                await this.chargerFormations();
            } catch (e) {
                this.erreurFormulaire = t("commun.erreur");
            } finally {
                this.enregistrement = false;
            }
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