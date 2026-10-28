demarrerApp({
    data() {
        return {
            id: new URLSearchParams(window.location.search).get("id"),
            formation: {
                titre: "",
                description: "",
                date_debut: "",
                date_fin: "",
                lieu: "",
                nb_places: null,
                tarif: 0,
                materiel_necessaire: "",
                statut: "brouillon",
            },
            erreur: "",
            succes: "",
            enregistrement: false,
        };
    },

    computed: {
        modification() {
            return Boolean(this.id);
        },
    },

    mounted() {
        if (this.modification) {
            this.chargerFormation();
        }
    },

    methods: {
        formaterStatut,

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

        async chargerFormation() {
            this.erreur = "";

            try {
                const reponse = await appelerApi(`/formations/${this.id}`);

                if (!reponse.ok) {
                    this.erreur = (await reponse.json()).erreur;
                    return;
                }

                const f = await reponse.json();

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
            } catch (e) {
                this.erreur = t("salaries.formations.erreur_chargement");
            }
        },

        async enregistrer() {
            this.erreur = "";
            this.succes = "";

            if (
                this.formation.date_fin &&
                this.formation.date_fin < this.formation.date_debut
            ) {
                this.erreur = "La date de fin doit être postérieure à la date de début.";
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
                    this.modification
                        ? `/formations/${this.id}`
                        : "/formations",
                    {
                        method: this.modification ? "PUT" : "POST",
                        body: JSON.stringify(donnees),
                    }
                );

                if (!reponse.ok) {
                    this.erreur = (await reponse.json()).erreur
                        || "Impossible d'enregistrer la formation.";
                    return;
                }

                window.location.href = "formations.php";
            } catch (e) {
                this.erreur = "Une erreur est survenue lors de l'enregistrement.";
            } finally {
                this.enregistrement = false;
            }
        },
    },
});