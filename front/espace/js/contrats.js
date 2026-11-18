// Contrats & Facturation (professionnels) : abonnement, campagnes publicitaires (mise en avant d'un projet
// public ou d'une annonce en ligne, 100 à 500 €/mois), factures et contrats archivés, historique des paiements,
// et catalogue des prestations de la plateforme.

const TYPES_CONTRAT = ["facture", "contrat"];
const campagneVide = () => ({ cible_type: "projet", cible_id: null, montant: 200 });

demarrerApp({
    data() {
        return {
            abonnement: null,
            campagnes: [],
            documents: [],
            paiements: [],
            prestations: [],
            projets: [],
            annonces: [],
            campagneOuverte: false,
            formulaireCampagne: campagneVide(),
            erreurCampagne: "",
            campagneAArreter: null,
            achatEnCours: false,
            erreur: "",
            message: "",
        };
    },

    computed: {
        elementsDisponibles() {
            const occupes = new Set(this.campagnes.filter((c) => c.statut === "active").map((c) => `${c.cible_type}-${c.cible_id}`));
            const elements = this.formulaireCampagne.cible_type === "projet"
                ? this.projets.filter((p) => p.partage_public)
                : this.annonces.filter((a) => a.statut === "en_ligne");
            return elements.filter((e) => !occupes.has(`${this.formulaireCampagne.cible_type}-${e.id}`));
        },
    },

    mounted() {
        this.charger();
    },

    methods: {
        formaterPrix,
        formaterDate,

        async lire(chemin) {
            try {
                const reponse = await appelerApi(chemin);
                return reponse.ok ? await reponse.json() : null;
            } catch (erreur) {
                return null;
            }
        },

        async charger() {
            const [abonnement, campagnes, documents, paiements, prestations] = await Promise.all([
                this.lire("/pro/abonnement"), this.lire("/pro/campagnes"), this.lire("/moi/documents"),
                this.lire("/moi/paiements"), this.lire("/prestations"),
            ]);
            if (!campagnes || !documents || !paiements) {
                this.erreur = t("commun.erreur_chargement");
            }
            this.abonnement = abonnement;
            this.campagnes = campagnes || [];
            this.documents = (documents || []).filter((d) => TYPES_CONTRAT.includes(d.type_document));
            this.paiements = paiements || [];
            this.prestations = prestations || [];
        },

        lienPDF(id) {
            return `${API_BASE_URL}/prestations/${id}/pdf`;
        },

        async ouvrirDocument(document) {
            try {
                await recupererPdf(`/documents/${document.id}/pdf`);
            } catch (erreur) {
                this.erreur = erreur.message;
            }
        },

        async telechargerDocument(document) {
            try {
                await recupererPdf(`/documents/${document.id}/pdf`, { nomFichier: `document-${document.id}.pdf`, telecharger: true });
            } catch (erreur) {
                this.erreur = erreur.message;
            }
        },

        async ouvrirCampagne() {
            this.formulaireCampagne = campagneVide();
            this.erreurCampagne = "";
            this.message = "";
            const [projets, annonces] = await Promise.all([this.lire("/moi/projets"), this.lire("/annonces?mine=true")]);
            this.projets = projets || [];
            this.annonces = annonces || [];
            if (this.elementsDisponibles.length === 0 && this.annonces.some((a) => a.statut === "en_ligne")) {
                this.formulaireCampagne.cible_type = "annonce";
            }
            this.campagneOuverte = true;
        },

        // Paiement par la page de paiement (simulée tant que Stripe n'est pas branché)
        async lancerCampagne() {
            this.erreurCampagne = "";
            this.achatEnCours = true;
            const erreur = await demarrerAchat("publicite", 0, {
                cible_type: this.formulaireCampagne.cible_type,
                cible_id: this.formulaireCampagne.cible_id,
                montant: this.formulaireCampagne.montant,
            });
            if (erreur) {
                this.erreurCampagne = erreur;
                this.achatEnCours = false;
            }
        },

        async arreterCampagne() {
            const campagne = this.campagneAArreter;
            this.campagneAArreter = null;
            const reponse = await appelerApi(`/pro/campagnes/${campagne.id}/arreter`, { method: "POST" });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.message = t("contrats.campagnes.arretee");
            this.charger();
        },
    },
});
