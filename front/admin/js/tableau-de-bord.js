const COULEURS_TABLEAU = {
    vert: "#2E7D32",
    vertClair: "#81C784",
    nuit: "#222b3c",
    bleuClair: "#8E8BB8",
    beige: "#c9c1b1",
    orange: "#F9A825",
    marron: "#8D6E63",
    rouge: "#c62828",
};
const COULEUR_SOURCE = {
    abonnement: COULEURS_TABLEAU.vert,
    vente: COULEURS_TABLEAU.vertClair,
    publicite: COULEURS_TABLEAU.orange,
    formation: COULEURS_TABLEAU.nuit,
    atelier: COULEURS_TABLEAU.bleuClair,
    evenement: COULEURS_TABLEAU.marron,
};
// Pour les charges, une hausse est une mauvaise nouvelle : la couleur de l'évolution est inversée.
const KPI_FINANCES = [
    { cle: "revenus" },
    { cle: "charges", hausseNegative: true },
    { cle: "resultat" },
];
const KPI_ACTIVITE = ["nouveaux_acteurs", "acteurs_actifs", "annonces", "depots", "inscriptions"];
const A_TRAITER = [
    { cle: "offres_a_valider", lien: "validations.php" },
    { cle: "annonces_a_valider", lien: "annonces.php" },
    { cle: "sujets_a_moderer", lien: "moderation.php" },
    { cle: "signalements", lien: "moderation.php" },
    { cle: "conteneurs_maintenance", lien: "conteneurs.php" },
];
const graphiquesTableau = new Map();

function dessinerGraphique(nom, config) {
    graphiquesTableau.get(nom)?.destroy();
    const canvas = document.getElementById(`graphique-${nom}`);
    if (!canvas) {
        return;
    }
    config.options = { locale: langueCourante(), responsive: true, maintainAspectRatio: false, ...config.options };
    graphiquesTableau.set(nom, new Chart(canvas, config));
}

demarrerApp({
    data() {
        return {
            periode: periodeRaccourci("30_jours"),
            raccourcis: ["30_jours", "mois_courant", "annee_courante", "12_mois"],
            raccourciActif: "30_jours",
            tableau: null,
            erreur: "",
            numeroRequete: 0,
        };
    },

    computed: {
        kpiFinances() {
            const { actuel, precedent } = this.tableau;
            return [
                ...KPI_FINANCES.map(({ cle, hausseNegative }) => ({
                    cle,
                    valeur: actuel[cle],
                    evolution: this.evolution(actuel[cle], precedent[cle], hausseNegative),
                })),
                {
                    cle: "recurrent",
                    valeur: this.tableau.revenu_mensuel_recurrent,
                    detail: this.$t("admin.finances.recurrent_detail", {
                        abonnes: this.tableau.abonnes_actifs, campagnes: this.tableau.campagnes_actives,
                    }),
                },
            ];
        },

        kpiActivite() {
            const { actuel, precedent } = this.tableau;
            return KPI_ACTIVITE.map((cle) => ({
                cle,
                valeur: actuel[cle],
                evolution: this.evolution(actuel[cle], precedent[cle], false),
            }));
        },

        elementsATraiter() {
            return A_TRAITER.map((element) => ({ ...element, nombre: this.tableau.a_traiter[element.cle] }));
        },

        totalATraiter() {
            return this.elementsATraiter.reduce((total, element) => total + element.nombre, 0);
        },
    },

    watch: {
        periode: {
            deep: true,
            handler() {
                const choisi = this.raccourcis.find((r) => {
                    const p = periodeRaccourci(r);
                    return p.debut === this.periode.debut && p.fin === this.periode.fin;
                });
                this.raccourciActif = choisi || "";
            },
        },
    },

    mounted() {
        Chart.defaults.font.family = "Inter, sans-serif";
        Chart.defaults.color = "#4a5160";
        this.charger();
    },

    methods: {
        formaterMontant(montant) {
            return montant.toLocaleString(langueCourante(), { style: "currency", currency: "EUR" });
        },

        formaterNombre(nombre) {
            return nombre.toLocaleString(langueCourante());
        },

        formaterDate(date) {
            return new Date(`${date}T00:00:00Z`).toLocaleDateString(langueCourante(), {
                day: "numeric", month: "short", year: "numeric", timeZone: "UTC",
            });
        },

        formaterMois(mois) {
            return new Date(`${mois}-01T00:00:00Z`).toLocaleDateString(langueCourante(), {
                month: "short", year: "2-digit", timeZone: "UTC",
            });
        },

        formaterMoisDate(date) {
            return new Date(`${date}T00:00:00Z`).toLocaleDateString(langueCourante(), {
                month: "long", year: "numeric", timeZone: "UTC",
            });
        },

        // Évolution par rapport à la période précédente de même durée
        evolution(actuel, precedent, hausseNegative) {
            if (!precedent) {
                return null;
            }
            const variation = ((actuel - precedent) / Math.abs(precedent)) * 100;
            const arrondi = Math.round(variation);
            if (arrondi === 0) {
                const texte = this.$t("admin.tableau_de_bord.stable");
                return { fleche: "=", texte, description: texte, classe: "stable" };
            }
            const hausse = arrondi > 0;
            const texte = `${hausse ? "+" : ""}${arrondi.toLocaleString(langueCourante())} %`;
            return {
                fleche: hausse ? "▲" : "▼",
                texte,
                description: this.$t("admin.tableau_de_bord.vs_precedent", { variation: texte }),
                classe: hausse !== Boolean(hausseNegative) ? "positive" : "negative",
            };
        },

        choisirRaccourci(raccourci) {
            this.periode = periodeRaccourci(raccourci);
            this.charger();
        },

        async charger() {
            this.erreur = "";
            const numero = ++this.numeroRequete;
            const parametres = new URLSearchParams(this.periode);
            try {
                const reponse = await appelerApi(`/admin/tableau-de-bord?${parametres}`);
                const donnees = await reponse.json();
                if (numero !== this.numeroRequete) {
                    return;
                }
                if (!reponse.ok) {
                    this.erreur = donnees.erreur || this.$t("commun.erreur_chargement");
                    return;
                }
                this.tableau = donnees;
                this.$nextTick(() => this.dessiner());
            } catch (erreur) {
                if (numero === this.numeroRequete) {
                    this.erreur = this.$t("commun.erreur_chargement");
                }
            }
        },

        dessiner() {
            const t = this.tableau;
            const euros = (valeur) => this.formaterMontant(valeur);

            dessinerGraphique("finances", {
                data: {
                    labels: t.finances_par_mois.map((m) => this.formaterMois(m.mois)),
                    datasets: [
                        {
                            type: "line", label: this.$t("admin.tableau_de_bord.kpi.resultat"),
                            data: t.finances_par_mois.map((m) => m.resultat),
                            borderColor: COULEURS_TABLEAU.nuit, backgroundColor: COULEURS_TABLEAU.nuit,
                            tension: 0.3, pointRadius: 3, order: 0,
                        },
                        {
                            type: "bar", label: this.$t("admin.tableau_de_bord.kpi.revenus"),
                            data: t.finances_par_mois.map((m) => m.revenus),
                            backgroundColor: COULEURS_TABLEAU.vert, borderRadius: 4, order: 1,
                        },
                        {
                            type: "bar", label: this.$t("admin.tableau_de_bord.kpi.charges"),
                            data: t.finances_par_mois.map((m) => m.charges),
                            backgroundColor: COULEURS_TABLEAU.beige, borderRadius: 4, order: 1,
                        },
                    ],
                },
                options: {
                    interaction: { mode: "index", intersect: false },
                    plugins: {
                        legend: { position: "bottom" },
                        tooltip: { callbacks: { label: (c) => `${c.dataset.label} : ${euros(c.parsed.y)}` } },
                    },
                    scales: {
                        y: {
                            ticks: {
                                callback: (valeur) => valeur.toLocaleString(langueCourante(), {
                                    style: "currency", currency: "EUR", maximumFractionDigits: 0,
                                }),
                            },
                        },
                    },
                },
            });

            const series = [
                ["nouveaux_acteurs", COULEURS_TABLEAU.vert],
                ["annonces", COULEURS_TABLEAU.orange],
                ["depots", COULEURS_TABLEAU.marron],
                ["inscriptions", COULEURS_TABLEAU.nuit],
            ];
            dessinerGraphique("activite", {
                type: "line",
                data: {
                    labels: t.activite_par_mois.map((m) => this.formaterMois(m.mois)),
                    datasets: series.map(([cle, couleur]) => ({
                        label: this.$t(`admin.tableau_de_bord.kpi.${cle}`),
                        data: t.activite_par_mois.map((m) => m[cle]),
                        borderColor: couleur, backgroundColor: couleur, tension: 0.3, pointRadius: 2,
                    })),
                },
                options: {
                    interaction: { mode: "index", intersect: false },
                    plugins: { legend: { position: "bottom" } },
                    scales: { y: { beginAtZero: true, ticks: { precision: 0 } } },
                },
            });

            const sources = t.revenus_par_source.filter((s) => s.revenu > 0);
            dessinerGraphique("sources", {
                type: "doughnut",
                data: {
                    labels: sources.map((s) => this.$t(`admin.finances.source.${s.source}`)),
                    datasets: [{
                        data: sources.map((s) => s.revenu),
                        backgroundColor: sources.map((s) => COULEUR_SOURCE[s.source]),
                        borderWidth: 2,
                    }],
                },
                options: {
                    cutout: "62%",
                    plugins: {
                        legend: { position: "bottom" },
                        tooltip: {
                            callbacks: {
                                label: (c) => `${c.label} : ${euros(c.parsed)} (${Math.round((c.parsed / t.actuel.revenus) * 100)} %)`,
                            },
                        },
                    },
                },
            });
        },
    },
});
