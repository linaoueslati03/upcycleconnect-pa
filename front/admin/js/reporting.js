const COULEURS = {
    vert: "#2E7D32",
    vertClair: "#81C784",
    bleu: "#080341",
    bleuClair: "#8E8BB8",
    rouge: "#c62828",
};
const TOTAUX = ["acteurs", "particuliers", "professionnels", "nouveaux_periode", "actifs_periode"];
const graphiques = new Map();

function dessiner(nom, config) {
    const ancien = graphiques.get(nom);
    if (ancien) {
        ancien.destroy();
    }
    const canvas = document.getElementById(`graphique-${nom}`);
    if (!canvas) {
        return;
    }
    config.options = { locale: langueCourante(), responsive: true, maintainAspectRatio: false, ...config.options };
    graphiques.set(nom, new Chart(canvas, config));
}

function viderGraphiques() {
    graphiques.forEach((graphique) => graphique.destroy());
    graphiques.clear();
}

function barresSimples(libelles, valeurs, couleurs) {
    return {
        type: "bar",
        data: { labels: libelles, datasets: [{ data: valeurs, backgroundColor: couleurs }] },
        options: {
            plugins: { legend: { display: false } },
            scales: { y: { beginAtZero: true, ticks: { precision: 0 } } },
        },
    };
}

demarrerApp({
    data() {
        return {
            periode: periodeRaccourci("12_mois"),
            raccourcis: ["12_mois", "mois_courant", "annee_courante"],
            onglets: ["acteurs"],
            onglet: "acteurs",
            totaux: TOTAUX,
            graphiquesEtat: ["role", "statut", "abonnement", "score"],
            graphiquesPeriode: ["nouveaux", "sites"],
            acteurs: null,
            erreur: "",
        };
    },

    mounted() {
        Chart.defaults.font.family = "Inter, sans-serif";
        Chart.defaults.color = "#4a5160";
        this.chargerOnglet();
    },

    methods: {
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
                month: "short", year: "numeric", timeZone: "UTC",
            });
        },

        parametresPeriode() {
            return new URLSearchParams({ debut: this.periode.debut, fin: this.periode.fin });
        },

        async lireErreur(reponse) {
            try {
                return (await reponse.json()).erreur || t("commun.erreur");
            } catch (erreur) {
                return t("commun.erreur");
            }
        },

        async chargerActeurs() {
            const reponse = await appelerApi(`/admin/reporting/acteurs?${this.parametresPeriode()}`);
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.acteurs = await reponse.json();
            await this.$nextTick();
            this.dessinerActeurs();
        },

        dessinerActeurs() {
            const a = this.acteurs;
            const libelles = (liste, prefixe) => liste.map((c) => t(`admin.reporting.${prefixe}.${c.libelle}`));
            const nombres = (liste) => liste.map((c) => c.nombre);

            dessiner("role", {
                type: "doughnut",
                data: {
                    labels: libelles(a.par_role, "role"),
                    datasets: [{ data: nombres(a.par_role), backgroundColor: [COULEURS.vert, COULEURS.bleu] }],
                },
                options: { plugins: { legend: { position: "bottom" } } },
            });

            dessiner("statut", barresSimples(
                libelles(a.par_statut, "statut"), nombres(a.par_statut), [COULEURS.vert, COULEURS.rouge],
            ));

            dessiner("abonnement", barresSimples(
                libelles(a.par_abonnement, "abonnement"), nombres(a.par_abonnement), [COULEURS.bleuClair, COULEURS.vert],
            ));

            const score = barresSimples(
                a.distribution_score.map((c) => c.libelle), nombres(a.distribution_score), COULEURS.vert,
            );
            score.options.scales.x = { title: { display: true, text: t("admin.reporting.points") } };
            dessiner("score", score);

            dessiner("nouveaux", {
                type: "line",
                data: {
                    labels: a.nouveaux_par_mois.map((m) => this.formaterMois(m.mois)),
                    datasets: [
                        {
                            label: t("admin.reporting.role.particulier"),
                            data: a.nouveaux_par_mois.map((m) => m.particuliers),
                            borderColor: COULEURS.vert,
                            backgroundColor: COULEURS.vert,
                            tension: 0.3,
                        },
                        {
                            label: t("admin.reporting.role.professionnel"),
                            data: a.nouveaux_par_mois.map((m) => m.professionnels),
                            borderColor: COULEURS.bleu,
                            backgroundColor: COULEURS.bleu,
                            tension: 0.3,
                        },
                    ],
                },
                options: {
                    plugins: { legend: { position: "bottom" } },
                    scales: { y: { beginAtZero: true, ticks: { precision: 0 } } },
                },
            });

            dessiner("sites", {
                type: "bar",
                data: {
                    labels: a.activite_par_site.map((s) => (s.site === "autre" ? t("admin.reporting.site_autre") : s.site)),
                    datasets: [
                        {
                            label: t("admin.reporting.demandes"),
                            data: a.activite_par_site.map((s) => s.demandes),
                            backgroundColor: COULEURS.vert,
                        },
                        {
                            label: t("admin.reporting.acteurs_distincts"),
                            data: a.activite_par_site.map((s) => s.acteurs),
                            backgroundColor: COULEURS.bleuClair,
                        },
                    ],
                },
                options: {
                    indexAxis: "y",
                    plugins: { legend: { position: "bottom" } },
                    scales: { x: { beginAtZero: true, ticks: { precision: 0 } } },
                },
            });
        },

        appliquerPeriode() {
            this.erreur = "";
            this.chargerOnglet();
        },

        choisirRaccourci(raccourci) {
            this.periode = periodeRaccourci(raccourci);
            this.appliquerPeriode();
        },

        changerOnglet(onglet) {
            this.onglet = onglet;
            this.erreur = "";
            viderGraphiques();
            this.chargerOnglet();
        },

        chargerOnglet() {
            const chargements = {
                acteurs: this.chargerActeurs,
            };
            if (chargements[this.onglet]) {
                chargements[this.onglet]();
            }
        },
    },
});
