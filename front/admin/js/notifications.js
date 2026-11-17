const CIBLES = ["tous", "particuliers", "professionnels", "utilisateur", "prediction"];
const CLASSES = ["formation", "atelier", "evenement", "depot", "achat_annonce", "premium"];
const ROLES_ACTEURS = [1, 2];

function formulaireVide() {
    return {
        cible: "tous", destinataire_id: null, classe_predite: "formation",
        titre_fr: "", message_fr: "", titre_en: "", message_en: "",
    };
}

demarrerApp({
    data() {
        return {
            cibles: CIBLES,
            classes: CLASSES,
            formulaire: formulaireVide(),
            formulaireOuvert: false,
            destinataire: null,
            recherche: "",
            resultatsRecherche: [],
            delaiRecherche: null,
            apercu: null,
            envoiEnCours: false,
            resultat: null,
            erreur: "",
            historique: { notifications: [], total: 0, par_page: 20 },
            page: 1,
        };
    },

    computed: {
        nombrePages() {
            return Math.max(1, Math.ceil(this.historique.total / this.historique.par_page));
        },

        messageResultat() {
            if (!this.resultat) {
                return "";
            }
            if (this.resultat.statut === "envoyee") {
                return t("admin.notifications.envoyee", { nombre: this.resultat.nb_destinataires });
            }
            return t("admin.notifications.echec", { erreur: this.resultat.erreur });
        },
    },

    mounted() {
        this.actualiserApercu();
        this.chargerHistorique(1);
    },

    methods: {
        formaterDate(date) {
            return new Date(date).toLocaleString(langueCourante(), {
                day: "numeric", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit",
            });
        },

        libelleCible(notification) {
            const libelle = t(`admin.notifications.cibles.${notification.cible}`);
            if (notification.cible === "utilisateur") {
                return `${libelle} : ${notification.destinataire || "—"}`;
            }
            if (notification.cible === "prediction") {
                return `${libelle} : ${t(`admin.reporting.predictions.classe.${notification.classe_predite}`)}`;
            }
            return libelle;
        },

        async lireErreur(reponse) {
            try {
                return (await reponse.json()).erreur || t("commun.erreur");
            } catch (erreur) {
                return t("commun.erreur");
            }
        },

        parametresCible() {
            const parametres = new URLSearchParams({ cible: this.formulaire.cible });
            if (this.formulaire.cible === "utilisateur" && this.formulaire.destinataire_id) {
                parametres.set("destinataire_id", this.formulaire.destinataire_id);
            }
            if (this.formulaire.cible === "prediction") {
                parametres.set("classe_predite", this.formulaire.classe_predite);
            }
            return parametres;
        },

        async actualiserApercu() {
            this.apercu = null;
            if (this.formulaire.cible === "utilisateur" && !this.formulaire.destinataire_id) {
                return;
            }
            const reponse = await appelerApi(`/admin/notifications/destinataires?${this.parametresCible()}`);
            if (reponse.ok) {
                this.apercu = (await reponse.json()).nombre;
            }
        },

        ouvrirFormulaire() {
            this.resultat = null;
            this.erreur = "";
            this.formulaireOuvert = true;
        },

        changerCible() {
            this.resultat = null;
            this.erreur = "";
            this.retirerDestinataire();
        },

        rechercherDestinataires() {
            clearTimeout(this.delaiRecherche);
            this.delaiRecherche = setTimeout(async () => {
                const texte = this.recherche.trim();
                if (texte.length < 2) {
                    this.resultatsRecherche = [];
                    return;
                }
                const reponse = await appelerApi(`/utilisateurs?${new URLSearchParams({ recherche: texte, statut: "actif" })}`);
                if (reponse.ok) {
                    const { utilisateurs } = await reponse.json();
                    this.resultatsRecherche = utilisateurs.filter((u) => ROLES_ACTEURS.includes(u.role_id)).slice(0, 8);
                }
            }, 250);
        },

        choisirDestinataire(utilisateur) {
            this.destinataire = utilisateur;
            this.formulaire.destinataire_id = utilisateur.id;
            this.resultatsRecherche = [];
            this.actualiserApercu();
        },

        retirerDestinataire() {
            this.destinataire = null;
            this.formulaire.destinataire_id = null;
            this.recherche = "";
            this.resultatsRecherche = [];
            this.actualiserApercu();
        },

        async envoyer() {
            this.erreur = "";
            this.resultat = null;
            if (!confirm(t("admin.notifications.confirmer", { nombre: this.apercu }))) {
                return;
            }
            this.envoiEnCours = true;
            try {
                const corps = { ...this.formulaire };
                if (corps.cible !== "utilisateur") {
                    corps.destinataire_id = null;
                }
                const reponse = await appelerApi("/admin/notifications", { method: "POST", body: JSON.stringify(corps) });
                if (!reponse.ok) {
                    this.erreur = await this.lireErreur(reponse);
                    return;
                }
                this.resultat = await reponse.json();
                if (this.resultat.statut === "envoyee") {
                    const { cible, classe_predite: classe } = this.formulaire;
                    this.formulaire = { ...formulaireVide(), cible, classe_predite: classe };
                    this.retirerDestinataire();
                    this.formulaireOuvert = false;
                }
                this.chargerHistorique(1);
            } finally {
                this.envoiEnCours = false;
            }
        },

        async chargerHistorique(page) {
            const reponse = await appelerApi(`/admin/notifications?page=${page}`);
            if (!reponse.ok) {
                this.erreur = await this.lireErreur(reponse);
                return;
            }
            this.historique = await reponse.json();
            this.page = page;
        },
    },
});
