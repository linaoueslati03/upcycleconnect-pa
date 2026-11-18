// Statistiques et impact des professionnels (Premium, contrôlé par l'API) :
// objets récupérés, poids réemployé, CO2 évité estimé (poids × facteur sourcé du matériau),
// achats, stocks disponibles en temps réel par site, alertes de collecte.

const COULEURS_MATERIAU = {
    bois: "#8D6E63",
    metal: "#607D8B",
    textile: "#2E7D32",
    verre: "#81C784",
    plastique: "#F9A825",
};
const alerteVide = () => ({ categorie_id: null, site: "", mot_cle: "" });
const graphiquesStats = new Map();

function dessinerStat(nom, config) {
    graphiquesStats.get(nom)?.destroy();
    const canvas = document.getElementById(`graphique-${nom}`);
    if (!canvas) {
        return;
    }
    config.options = { locale: langueCourante(), responsive: true, maintainAspectRatio: false, ...config.options };
    graphiquesStats.set(nom, new Chart(canvas, config));
}

demarrerApp({
    data() {
        return {
            stats: null,
            estFreemium: false,
            prixPremium: 20,
            achatEnCours: false,
            categories: [],
            sites: [],
            alertes: [],
            alerteOuverte: false,
            formulaireAlerte: alerteVide(),
            erreurAlerte: "",
            messageAlerte: "",
            erreur: "",
        };
    },

    computed: {
        stocksParSite() {
            const sites = new Map();
            for (const d of this.stats.disponibles) {
                if (!sites.has(d.site)) {
                    sites.set(d.site, { site: d.site, total: 0, materiaux: [] });
                }
                const site = sites.get(d.site);
                site.total += d.nombre;
                site.materiaux.push(d);
            }
            return [...sites.values()];
        },
    },

    async mounted() {
        Chart.defaults.font.family = "Inter, sans-serif";
        Chart.defaults.color = "#4a5160";
        const [categories, conteneurs] = await Promise.all([appelerApi("/categories"), appelerApi("/conteneurs")]);
        if (categories.ok) {
            this.categories = (await categories.json()).filter((c) => c.actif);
        }
        if (conteneurs.ok) {
            this.sites = [...new Set((await conteneurs.json()).map((c) => c.site))].sort();
        }
        this.charger();
    },

    methods: {
        formaterPrix,
        formaterDate,

        formaterMontant(montant) {
            return Number(montant).toLocaleString(langueCourante(), { style: "currency", currency: "EUR" });
        },

        formaterNombre(nombre, decimales = 0) {
            return Number(nombre).toLocaleString(langueCourante(), { maximumFractionDigits: decimales });
        },

        formaterKg(kg) {
            return `${this.formaterNombre(kg, kg < 10 ? 1 : 0)} kg`;
        },

        nomMateriau(m) {
            return t("categorie." + m.code, m.libelle);
        },

        async charger() {
            const reponse = await appelerApi("/pro/statistiques");
            if (reponse.status === 403) {
                this.estFreemium = true;
                const abonnement = await appelerApi("/pro/abonnement");
                if (abonnement.ok) {
                    this.prixPremium = (await abonnement.json()).prix_mensuel;
                }
                return;
            }
            if (!reponse.ok) {
                this.erreur = t("commun.erreur_chargement");
                return;
            }
            this.stats = await reponse.json();
            this.chargerAlertes();
            this.$nextTick(() => this.dessiner());
        },

        async passerPremium() {
            this.achatEnCours = true;
            const erreur = await demarrerAchat("abonnement", 0);
            if (erreur) {
                this.erreur = erreur;
                this.achatEnCours = false;
            }
        },

        dessiner() {
            const mois = this.stats.par_mois;
            const kg = (v) => this.formaterKg(v);
            dessinerStat("mois", {
                data: {
                    labels: mois.map((m) => new Date(`${m.mois}-01T00:00:00Z`).toLocaleDateString(langueCourante(), {
                        month: "short", year: "2-digit", timeZone: "UTC",
                    })),
                    datasets: [
                        {
                            type: "line", label: t("statistiques.graphique.co2"), data: mois.map((m) => m.co2_kg),
                            borderColor: "#222b3c", backgroundColor: "#222b3c", tension: 0.3, yAxisID: "co2", order: 0,
                        },
                        {
                            type: "bar", label: t("statistiques.graphique.objets"), data: mois.map((m) => m.nombre),
                            backgroundColor: "#2E7D32", borderRadius: 4, yAxisID: "objets", order: 1,
                        },
                    ],
                },
                options: {
                    interaction: { mode: "index", intersect: false },
                    plugins: {
                        legend: { position: "bottom" },
                        tooltip: { callbacks: { label: (c) => `${c.dataset.label} : ${c.dataset.yAxisID === "co2" ? kg(c.parsed.y) : c.parsed.y}` } },
                    },
                    scales: {
                        objets: { position: "left", beginAtZero: true, ticks: { precision: 0 } },
                        co2: { position: "right", beginAtZero: true, grid: { drawOnChartArea: false }, ticks: { callback: (v) => kg(v) } },
                    },
                },
            });

            const materiaux = this.stats.par_materiau.filter((m) => m.co2_kg > 0);
            dessinerStat("materiaux", {
                type: "doughnut",
                data: {
                    labels: materiaux.map((m) => this.nomMateriau(m)),
                    datasets: [{ data: materiaux.map((m) => m.co2_kg), backgroundColor: materiaux.map((m) => COULEURS_MATERIAU[m.code] || "#9e9e9e"), borderWidth: 2 }],
                },
                options: {
                    cutout: "60%",
                    plugins: { legend: { position: "bottom" }, tooltip: { callbacks: { label: (c) => `${c.label} : ≈ ${kg(c.parsed)}` } } },
                },
            });
        },

        async chargerAlertes() {
            const reponse = await appelerApi("/pro/alertes");
            if (reponse.ok) {
                this.alertes = await reponse.json();
            }
        },

        descriptionAlerte(a) {
            const morceaux = [];
            const categorie = this.categories.find((c) => c.id === a.categorie_id);
            morceaux.push(categorie ? this.nomMateriau(categorie) : t("statistiques.alertes.tous_materiaux"));
            morceaux.push(a.site || t("statistiques.alertes.tous_sites"));
            if (a.mot_cle) {
                morceaux.push(t("statistiques.alertes.contenant", { mot: a.mot_cle }));
            }
            return morceaux.join(" · ");
        },

        ouvrirAlerte() {
            this.formulaireAlerte = alerteVide();
            this.erreurAlerte = "";
            this.messageAlerte = "";
            this.alerteOuverte = true;
        },

        async enregistrerAlerte() {
            this.erreurAlerte = "";
            const reponse = await appelerApi("/pro/alertes", { method: "POST", body: JSON.stringify(this.formulaireAlerte) });
            if (!reponse.ok) {
                this.erreurAlerte = await messageErreur(reponse);
                return;
            }
            this.alerteOuverte = false;
            this.messageAlerte = t("statistiques.alertes.enregistree");
            this.chargerAlertes();
        },

        async supprimerAlerte(alerte) {
            const reponse = await appelerApi(`/pro/alertes/${alerte.id}`, { method: "DELETE" });
            if (reponse.ok) {
                this.alertes = this.alertes.filter((a) => a.id !== alerte.id);
            }
        },
    },
});
