// Page de paiement. Tant que Stripe n'est pas branché, le paiement est simulé : le bouton « Payer »
// confirme la commande sans débit réel. Avec Stripe, cette page servira de page de retour
// (success_url / cancel_url) et affichera simplement le résultat de la commande.

demarrerApp({
    data() {
        return {
            commande: null,
            erreur: "",
            enCours: false,
            maintenant: Date.now(),
        };
    },

    computed: {
        pageRetour() {
            return { formation: "catalogue.php", vente: "annonces.php", abonnement: "statistiques.php", publicite: "contrats.php" }[this.commande.type];
        },

        tempsRestant() {
            const secondes = Math.max(0, Math.floor((new Date(this.commande.date_expiration) - this.maintenant) / 1000));
            return `${Math.floor(secondes / 60)} min ${String(secondes % 60).padStart(2, "0")} s`;
        },
    },

    mounted() {
        this.charger();
        this.minuteur = setInterval(() => {
            this.maintenant = Date.now();
            const enAttente = this.commande?.statut === "en_attente";
            // Expiration de la réservation, ou retour de Stripe en attente du webhook (vérifié toutes les 3 s)
            if (enAttente && (new Date(this.commande.date_expiration) <= this.maintenant
                || (this.commande.fournisseur !== "simulation" && Math.floor(this.maintenant / 1000) % 3 === 0))) {
                this.charger();
            }
        }, 1000);
    },

    unmounted() {
        clearInterval(this.minuteur);
    },

    methods: {
        formaterPrix,

        idCommande() {
            return Number(new URLSearchParams(window.location.search).get("commande"));
        },

        async appel(chemin, options) {
            this.erreur = "";
            try {
                const reponse = await appelerApi(chemin, options);
                if (!reponse.ok) {
                    const message = await messageErreur(reponse);
                    if (reponse.status === 409) {
                        // La commande a changé entre-temps (expirée, déjà payée…) : on affiche son état réel
                        await this.charger();
                    }
                    this.erreur = message;
                    return;
                }
                this.commande = await reponse.json();
            } catch (erreur) {
                this.erreur = t("commun.serveur_injoignable");
            }
        },

        charger() {
            if (!this.idCommande()) {
                this.erreur = t("paiement.introuvable");
                return;
            }
            return this.appel(`/commandes/${this.idCommande()}`);
        },

        async payer() {
            this.enCours = true;
            await this.appel(`/commandes/${this.commande.id}/payer-simulation`, { method: "POST" });
            this.enCours = false;
        },

        async annuler() {
            this.enCours = true;
            await this.appel(`/commandes/${this.commande.id}/annuler`, { method: "POST" });
            this.enCours = false;
        },

        async ouvrirFacture() {
            try {
                await recupererPdf(`/documents/${this.commande.facture_id}/pdf`);
            } catch (erreur) {
                this.erreur = erreur.message;
            }
        },
    },
});
