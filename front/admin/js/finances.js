const SOURCES_PAIEMENT = ["abonnement", "vente", "publicite", "formation", "atelier", "evenement"];
const STATUTS_PAIEMENT = ["reussi", "en_attente", "echoue", "rembourse"];
const CATEGORIES_CHARGE = ["salaires", "loyers", "infrastructure", "marketing", "autre"];

function dateLocale(date) {
    const mois = String(date.getMonth() + 1).padStart(2, "0");
    const jour = String(date.getDate()).padStart(2, "0");
    return `${date.getFullYear()}-${mois}-${jour}`;
}

function periodeRaccourci(raccourci) {
    const aujourdhui = new Date();
    const debuts = {
        "12_mois": new Date(aujourdhui.getFullYear(), aujourdhui.getMonth() - 11, 1),
        mois_courant: new Date(aujourdhui.getFullYear(), aujourdhui.getMonth(), 1),
        annee_courante: new Date(aujourdhui.getFullYear(), 0, 1),
    };
    return { debut: dateLocale(debuts[raccourci]), fin: dateLocale(aujourdhui) };
}

function chargeVide() {
    return { libelle: "", categorie: "salaires", montant: null, date_charge: dateLocale(new Date()) };
}

demarrerApp({
    data() {
        return {
            periode: periodeRaccourci("12_mois"),
            raccourcis: ["12_mois", "mois_courant", "annee_courante"],
            onglets: ["synthese", "transactions", "abonnements", "publicite", "charges", "parametres"],
            onglet: "synthese",
            sources: SOURCES_PAIEMENT,
            statutsPaiement: STATUTS_PAIEMENT,
            categoriesCharge: CATEGORIES_CHARGE,
            resume: null,
            paiements: [],
            filtresPaiements: { source: "", statut: "" },
            pagePaiements: 1,
            totalPaiements: 0,
            parPagePaiements: 20,
            abonnements: [],
            campagnes: [],
            charges: [],
            nouvelleCharge: chargeVide(),
            parametres: null,
            messageParametres: "",
            erreur: "",
        };
    },

    computed: {
        nombrePagesPaiements() {
            return Math.max(1, Math.ceil(this.totalPaiements / this.parPagePaiements));
        },
    },

    mounted() {
        this.chargerResume();
    },

    methods: {
        formaterMontant(montant) {
            return montant.toLocaleString(langueCourante(), { style: "currency", currency: "EUR" });
        },

        formaterTaux(taux) {
            return (taux / 100).toLocaleString(langueCourante(), { style: "percent", maximumFractionDigits: 2 });
        },

        formaterDate(date) {
            return new Date(date).toLocaleDateString(langueCourante(), {
                day: "numeric", month: "short", year: "numeric", timeZone: "UTC",
            });
        },

        formaterMois(mois) {
            return new Date(`${mois}-01T00:00:00Z`).toLocaleDateString(langueCourante(), {
                month: "long", year: "numeric", timeZone: "UTC",
            });
        },

        formaterPeriode(debut, fin) {
            if (!fin) {
                return t("admin.finances.depuis", { date: this.formaterDate(debut) });
            }
            return `${this.formaterDate(debut)} → ${this.formaterDate(fin)}`;
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

        async chargerResume() {
            const reponse = await appelerApi(`/admin/finances/resume?${this.parametresPeriode()}`);
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.resume = await reponse.json();
        },

        appliquerPeriode() {
            this.erreur = "";
            this.pagePaiements = 1;
            this.chargerResume();
            this.chargerOnglet();
        },

        choisirRaccourci(raccourci) {
            this.periode = periodeRaccourci(raccourci);
            this.appliquerPeriode();
        },

        changerOnglet(onglet) {
            this.onglet = onglet;
            this.erreur = "";
            this.chargerOnglet();
        },

        chargerOnglet() {
            const chargements = {
                transactions: this.chargerPaiements,
                abonnements: this.chargerAbonnements,
                publicite: this.chargerCampagnes,
                charges: this.chargerCharges,
                parametres: this.chargerParametres,
            };
            if (chargements[this.onglet]) {
                chargements[this.onglet]();
            }
        },

        async chargerPaiements() {
            const supplement = { page: this.pagePaiements };
            for (const [nom, valeur] of Object.entries(this.filtresPaiements)) {
                if (valeur !== "") {
                    supplement[nom] = valeur;
                }
            }
            const reponse = await appelerApi(`/admin/finances/paiements?${this.parametresPeriode(supplement)}`);
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            const resultat = await reponse.json();
            this.paiements = resultat.paiements;
            this.totalPaiements = resultat.total;
            this.parPagePaiements = resultat.par_page;
        },

        filtrerPaiements() {
            this.pagePaiements = 1;
            this.chargerPaiements();
        },

        changerPagePaiements(page) {
            this.pagePaiements = page;
            this.chargerPaiements();
        },

        async chargerAbonnements() {
            const reponse = await appelerApi("/admin/finances/abonnements");
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.abonnements = await reponse.json();
        },

        async chargerCampagnes() {
            const reponse = await appelerApi("/admin/finances/campagnes");
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.campagnes = await reponse.json();
        },

        async chargerCharges() {
            const reponse = await appelerApi(`/admin/finances/charges?${this.parametresPeriode()}`);
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.charges = await reponse.json();
        },

        async ajouterCharge() {
            this.erreur = "";
            const reponse = await appelerApi("/admin/finances/charges", {
                method: "POST",
                body: JSON.stringify(this.nouvelleCharge),
            });
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.nouvelleCharge = chargeVide();
            this.chargerCharges();
            this.chargerResume();
        },

        async supprimerCharge(charge) {
            if (!confirm(t("admin.finances.confirmer_suppression_charge", { libelle: charge.libelle }))) {
                return;
            }
            this.erreur = "";
            const reponse = await appelerApi(`/admin/finances/charges/${charge.id}`, { method: "DELETE" });
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.chargerCharges();
            this.chargerResume();
        },

        async chargerParametres() {
            const reponse = await appelerApi("/admin/finances/parametres");
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.parametres = await reponse.json();
        },

        async enregistrerParametres() {
            this.erreur = "";
            this.messageParametres = "";
            const reponse = await appelerApi("/admin/finances/parametres", {
                method: "PUT",
                body: JSON.stringify({
                    taux_commission: this.parametres.taux_commission,
                    prix_premium_mensuel: this.parametres.prix_premium_mensuel,
                }),
            });
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.parametres = await reponse.json();
            this.messageParametres = t("admin.finances.parametres_enregistres");
        },
    },
});
