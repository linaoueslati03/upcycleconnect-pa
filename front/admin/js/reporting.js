const COULEURS = {
    vert: "#2E7D32",
    vertClair: "#81C784",
    bleu: "#080341",
    bleuClair: "#8E8BB8",
    rouge: "#c62828",
};
const TOTAUX = ["acteurs", "particuliers", "professionnels", "nouveaux_periode", "actifs_periode"];
const TYPES_PRESTATION = ["formation", "atelier", "evenement"];
const COULEUR_TYPE = { formation: COULEURS.vert, atelier: COULEURS.bleuClair, evenement: COULEURS.bleu };
const CLASSES_PREDICTION = ["formation", "atelier", "evenement", "depot", "achat_annonce", "premium"];
const COULEUR_CLASSE = { ...COULEUR_TYPE, depot: "#8D6E63", achat_annonce: "#F9A825", premium: COULEURS.rouge };
const ROLES_ACTEUR = ["particulier", "professionnel"];
const VARIABLES_ENTIERES = ["anciennete_jours", "upcycling_score", "jours_depuis_derniere_demande"];
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
            onglets: ["acteurs", "prestations", "predictions"],
            onglet: "acteurs",
            totaux: TOTAUX,
            graphiquesEtat: ["role", "statut", "abonnement", "score"],
            graphiquesPeriode: ["nouveaux", "sites"],
            acteurs: null,
            typesPrestation: TYPES_PRESTATION,
            typePrestation: "",
            graphiquesPrestations: ["mois", "reinscription", "top"],
            prestations: null,
            classesPrediction: CLASSES_PREDICTION,
            rolesActeur: ROLES_ACTEUR,
            filtrePrediction: { classe: "", role: "" },
            pagePredictions: 1,
            predictions: null,
            detail: null,
            erreur: "",
        };
    },

    computed: {
        nombrePagesPredictions() {
            if (!this.predictions) {
                return 1;
            }
            return Math.max(1, Math.ceil(this.predictions.total / this.predictions.par_page));
        },
    },

    mounted() {
        Chart.defaults.font.family = "Inter, sans-serif";
        Chart.defaults.color = "#4a5160";
        const parametres = new URLSearchParams(window.location.search);
        if (this.onglets.includes(parametres.get("onglet"))) {
            this.onglet = parametres.get("onglet");
        }
        this.chargerOnglet();
        const acteur = Number(parametres.get("acteur"));
        if (this.onglet === "predictions" && acteur > 0) {
            this.expliquer(acteur);
        }
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

        formaterPourcentage(valeur) {
            if (valeur === null || valeur === undefined) {
                return "—";
            }
            return (valeur / 100).toLocaleString(langueCourante(), { style: "percent", maximumFractionDigits: 1 });
        },

        formaterMontant(montant) {
            return montant.toLocaleString(langueCourante(), { style: "currency", currency: "EUR" });
        },

        parametresPeriode(supplement = {}) {
            return new URLSearchParams({ debut: this.periode.debut, fin: this.periode.fin, ...supplement });
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

        async chargerPrestations() {
            this.erreur = "";
            const supplement = this.typePrestation ? { type: this.typePrestation } : {};
            const reponse = await appelerApi(`/admin/reporting/prestations?${this.parametresPeriode(supplement)}`);
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.prestations = await reponse.json();
            await this.$nextTick();
            this.dessinerPrestations();
        },

        dessinerPrestations() {
            const p = this.prestations;
            const types = p.type ? [p.type] : TYPES_PRESTATION;
            const nomType = (type) => t(`admin.reporting.prestations.type.${type}`);

            dessiner("mois", {
                type: "bar",
                data: {
                    labels: p.par_mois.map((m) => this.formaterMois(m.mois)),
                    datasets: types.map((type) => ({
                        label: nomType(type),
                        data: p.par_mois.map((m) => m[type]),
                        backgroundColor: COULEUR_TYPE[type],
                    })),
                },
                options: {
                    plugins: { legend: { position: "bottom" } },
                    scales: { x: { stacked: true }, y: { stacked: true, beginAtZero: true, ticks: { precision: 0 } } },
                },
            });

            const reinscription = barresSimples(
                p.par_type.map((s) => nomType(s.type)),
                p.par_type.map((s) => s.taux_reinscription ?? 0),
                p.par_type.map((s) => COULEUR_TYPE[s.type]),
            );
            reinscription.options.scales.y = {
                beginAtZero: true,
                ticks: { callback: (valeur) => this.formaterPourcentage(valeur) },
            };
            reinscription.options.plugins.tooltip = {
                callbacks: { label: (contexte) => this.formaterPourcentage(contexte.parsed.y) },
            };
            dessiner("reinscription", reinscription);

            const top = p.prestations.slice(0, 10);
            dessiner("top", {
                type: "bar",
                data: {
                    labels: top.map((s) => s.titre),
                    datasets: [{
                        label: t("admin.reporting.prestations.colonne.inscriptions"),
                        data: top.map((s) => s.inscriptions),
                        backgroundColor: top.map((s) => COULEUR_TYPE[s.type]),
                    }],
                },
                options: {
                    indexAxis: "y",
                    plugins: {
                        legend: { display: false },
                        tooltip: { callbacks: { afterLabel: (contexte) => nomType(top[contexte.dataIndex].type) } },
                    },
                    scales: { x: { beginAtZero: true, ticks: { precision: 0 } } },
                },
            });
        },

        couleurClasse(classe) {
            return COULEUR_CLASSE[classe];
        },

        nomClasse(classe) {
            return t(`admin.reporting.predictions.classe.${classe}`);
        },

        formaterProbabilite(probabilite) {
            return this.formaterPourcentage(probabilite * 100);
        },

        formaterDateHeure(version) {
            return new Date(version).toLocaleDateString(langueCourante(), { day: "numeric", month: "long", year: "numeric" });
        },

        async chargerPredictions(page = 1) {
            this.erreur = "";
            this.pagePredictions = page;
            const parametres = new URLSearchParams({ page });
            if (this.filtrePrediction.classe) {
                parametres.set("classe", this.filtrePrediction.classe);
            }
            if (this.filtrePrediction.role) {
                parametres.set("role", this.filtrePrediction.role);
            }
            const reponse = await appelerApi(`/admin/predictions?${parametres}`);
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.predictions = await reponse.json();
            await this.$nextTick();
            dessiner("repartition", barresSimples(
                this.predictions.repartition.map((c) => this.nomClasse(c.libelle)),
                this.predictions.repartition.map((c) => c.nombre),
                this.predictions.repartition.map((c) => COULEUR_CLASSE[c.libelle]),
            ));
        },

        filtrerPredictions() {
            this.chargerPredictions(1);
        },

        async expliquer(id) {
            this.erreur = "";
            const reponse = await appelerApi(`/admin/predictions/${id}`);
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.detail = await reponse.json();
            await this.$nextTick();
            document.getElementById("detail-prediction")?.scrollIntoView({ behavior: "smooth", block: "start" });
        },

        nombreDemandes(variables) {
            return variables.nb_formations + variables.nb_ateliers + variables.nb_evenements
                + variables.nb_depots + variables.nb_achats_annonces;
        },

        libelleValeur(variable, valeur) {
            if (variable === "role") {
                return t(`admin.reporting.predictions.role.${valeur}`);
            }
            if (variable === "derniere_demande") {
                return valeur === "aucune" ? t("admin.reporting.predictions.aucune_demande") : this.nomClasse(valeur);
            }
            return valeur;
        },

        phraseEtape(etape) {
            const base = "admin.reporting.predictions.chemin";
            const [variable, valeur] = etape.colonne.split("=");
            const sens = etape.inferieur_ou_egal ? "max" : "plus";

            if (variable === "role") {
                const role = etape.inferieur_ou_egal ? ROLES_ACTEUR.find((r) => r !== valeur) : valeur;
                return t(`${base}.indicateur.oui`, {
                    variable: t(`${base}.variable.role`),
                    valeur: this.libelleValeur("role", role),
                });
            }
            if (valeur !== undefined) {
                return t(`${base}.indicateur.${etape.inferieur_ou_egal ? "non" : "oui"}`, {
                    variable: t(`${base}.variable.${variable}`),
                    valeur: this.libelleValeur(variable, valeur),
                });
            }
            if (variable === "premium") {
                return t(`${base}.premium.${etape.inferieur_ou_egal ? "non" : "oui"}`);
            }
            if (variable.startsWith("part_")) {
                return t(`${base}.part.${sens}`, {
                    variable: t(`${base}.variable.${variable}`),
                    seuil: this.formaterProbabilite(etape.seuil),
                    valeur: this.formaterProbabilite(etape.valeur),
                });
            }
            if (variable === "montant_depense") {
                return t(`${base}.montant.${sens}`, {
                    variable: t(`${base}.variable.${variable}`),
                    seuil: this.formaterMontant(etape.seuil),
                    valeur: this.formaterMontant(etape.valeur),
                });
            }
            const entier = Math.floor(etape.seuil);
            const nombre = etape.inferieur_ou_egal ? entier : entier + 1;
            const cle = variable.startsWith("nb_") || VARIABLES_ENTIERES.includes(variable) ? sens : `${sens}_strict`;
            return t(`${base}.nombre.${cle}`, {
                variable: t(`${base}.variable.${variable}`),
                seuil: this.formaterNombre(cle.endsWith("_strict") ? etape.seuil : nombre),
                valeur: this.formaterNombre(etape.valeur),
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
                prestations: this.chargerPrestations,
                predictions: this.chargerPredictions,
            };
            if (chargements[this.onglet]) {
                chargements[this.onglet]();
            }
        },
    },
});
