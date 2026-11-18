demarrerApp({
    data() {
        return {
            profil: { nom: "", prenom: "", email: "", langue_preferee_id: null },
            brouillonProfil: { nom: "", prenom: "", email: "", langue_preferee_id: null },
            profilOuvert: false,
            motDePasseOuvert: false,
            langues: [],
            motsDePasse: { ancien_mot_de_passe: "", nouveau_mot_de_passe: "" },
            erreurProfil: "",
            messageProfil: "",
            erreurMotDePasse: "",
            messageMotDePasse: "",
            motDePasseObligatoire: doitChangerMotDePasse(),
            entreprise: { siret: "", adresse: "" },
            brouillonEntreprise: { siret: "", adresse: "" },
            entrepriseOuverte: false,
            erreurEntreprise: "",
            messageEntreprise: "",
            abonnement: null,
            achatEnCours: false,
            resiliationOuverte: false,
            messageAbonnement: "",
            erreurAbonnement: "",
            documents: [],
            etatPush: "verification",
            pushEnCours: false,
            messagePush: "",
            erreurPush: false,
            erreurDocuments: "",
        };
    },

    async mounted() {
        this.chargerDocuments();
        this.chargerEtatPush();
        const [reponseCompte, reponseLangues] = await Promise.all([appelerApi("/moi"), appelerApi("/langues")]);
        if (!reponseCompte.ok || !reponseLangues.ok) {
            this.erreurProfil = t("compte.erreur_chargement");
            return;
        }
        const compte = await reponseCompte.json();
        this.motDePasseObligatoire = compte.doit_changer_mot_de_passe;
        this.profil = {
            nom: compte.nom,
            prenom: compte.prenom,
            email: compte.email,
            langue_preferee_id: compte.langue_preferee_id,
        };
        this.entreprise = {
            siret: compte.siret || "",
            adresse: compte.adresse_entreprise || "",
        };
        if (roleConnecte() === "professionnel") {
            this.chargerAbonnement();
        }
        this.langues = await reponseLangues.json();
        if (this.motDePasseObligatoire) {
            this.ouvrirMotDePasse();
        }
    },

    computed: {
        libelleLangue() {
            return this.langues.find((l) => l.id === this.profil.langue_preferee_id)?.libelle;
        },
    },

    methods: {
        // Notifications push : l'état affiché est celui du navigateur et de OneSignal, pas une simple préférence
        async chargerEtatPush() {
            if (typeof notificationsPush === "undefined") {
                this.etatPush = "indisponibles";
                return;
            }
            this.etatPush = await notificationsPush.etat();
            notificationsPush.surChangement((etat) => {
                this.etatPush = etat;
            });
        },

        async activerPush() {
            this.pushEnCours = true;
            this.messagePush = "";
            try {
                this.etatPush = await notificationsPush.activer();
                this.erreurPush = this.etatPush !== "activees";
                this.messagePush = t(this.erreurPush ? "compte.notifications.non_activees" : "compte.notifications.activees_message");
            } catch (erreur) {
                this.erreurPush = true;
                this.messagePush = t("compte.notifications.non_activees");
            }
            this.pushEnCours = false;
        },

        async desactiverPush() {
            this.pushEnCours = true;
            this.messagePush = "";
            try {
                this.etatPush = await notificationsPush.desactiver();
                this.erreurPush = false;
                this.messagePush = t("compte.notifications.desactivees_message");
            } catch (erreur) {
                this.erreurPush = true;
                this.messagePush = t("commun.erreur");
            }
            this.pushEnCours = false;
        },

        async testerPush() {
            this.pushEnCours = true;
            this.messagePush = "";
            try {
                const reponse = await appelerApi("/moi/notifications/test", { method: "POST" });
                const resultat = reponse.ok ? await reponse.json() : null;
                this.erreurPush = !resultat || resultat.push !== "envoyee";
                this.messagePush = t(this.erreurPush ? "compte.notifications.test_echec" : "compte.notifications.test_envoye");
            } catch (erreur) {
                this.erreurPush = true;
                this.messagePush = t("compte.notifications.test_echec");
            }
            this.pushEnCours = false;
        },

        // PDF archivés à mon nom (contrats, attestations, documents administratifs)
        async chargerDocuments() {
            try {
                const reponse = await appelerApi("/moi/documents");
                if (!reponse.ok) {
                    this.erreurDocuments = t("commun.erreur_chargement");
                    return;
                }
                this.documents = await reponse.json();
            } catch (erreur) {
                this.erreurDocuments = t("commun.erreur_chargement");
            }
        },

        formaterDateDocument(date) {
            return new Date(date).toLocaleDateString(langueCourante(), { day: "numeric", month: "long", year: "numeric" });
        },

        async ouvrirDocument(d) {
            try {
                await recupererPdf(`/documents/${d.id}/pdf`);
            } catch (erreur) {
                this.erreurDocuments = erreur.message;
            }
        },

        async telechargerDocument(d) {
            try {
                await recupererPdf(`/documents/${d.id}/pdf`, { nomFichier: `document-${d.id}.pdf`, telecharger: true });
            } catch (erreur) {
                this.erreurDocuments = erreur.message;
            }
        },

        ouvrirProfil() {
            this.brouillonProfil = { ...this.profil };
            this.erreurProfil = "";
            this.messageProfil = "";
            this.profilOuvert = true;
        },

        ouvrirMotDePasse() {
            this.motsDePasse = { ancien_mot_de_passe: "", nouveau_mot_de_passe: "" };
            this.erreurMotDePasse = "";
            this.messageMotDePasse = "";
            this.motDePasseOuvert = true;
        },

        async enregistrerProfil() {
            this.erreurProfil = "";
            this.messageProfil = "";
            const reponse = await appelerApi("/moi", { method: "PUT", body: JSON.stringify(this.brouillonProfil) });
            if (!reponse.ok) {
                this.erreurProfil = await messageErreur(reponse);
                return;
            }
            this.profil = { ...this.brouillonProfil };
            this.profilOuvert = false;
            this.messageProfil = t("compte.profil_enregistre");

            const langue = this.langues.find((l) => l.id === this.profil.langue_preferee_id);
            if (langue && langue.code !== langueCourante()) {
                changerLangue(langue.code);
            }
        },

        async changerMotDePasse() {
            this.erreurMotDePasse = "";
            this.messageMotDePasse = "";
            const reponse = await appelerApi("/moi/mot-de-passe", {
                method: "PUT",
                body: JSON.stringify(this.motsDePasse),
            });
            if (!reponse.ok) {
                this.erreurMotDePasse = await messageErreur(reponse);
                return;
            }
            this.motsDePasse = { ancien_mot_de_passe: "", nouveau_mot_de_passe: "" };
            this.motDePasseOuvert = false;
            this.messageMotDePasse = t("compte.mot_de_passe_modifie");
            this.motDePasseObligatoire = false;
            localStorage.removeItem("changer_mot_de_passe");
        },

        formaterPrix,

        ouvrirEntreprise() {
            this.brouillonEntreprise = { siret: this.entreprise.siret, adresse: this.entreprise.adresse };
            this.erreurEntreprise = "";
            this.messageEntreprise = "";
            this.entrepriseOuverte = true;
        },

        async enregistrerEntreprise() {
            this.erreurEntreprise = "";
            const reponse = await appelerApi("/pro/entreprise", { method: "PUT", body: JSON.stringify(this.brouillonEntreprise) });
            if (!reponse.ok) {
                this.erreurEntreprise = await messageErreur(reponse);
                return;
            }
            const donnees = await reponse.json();
            this.entreprise = { siret: donnees.siret, adresse: donnees.adresse_entreprise };
            this.entrepriseOuverte = false;
            this.messageEntreprise = t("compte.entreprise.enregistree");
        },

        async chargerAbonnement() {
            const reponse = await appelerApi("/pro/abonnement");
            if (reponse.ok) {
                this.abonnement = await reponse.json();
            }
        },

        // Paiement par la page de paiement (simulée tant que Stripe n'est pas branché)
        async passerPremium() {
            this.erreurAbonnement = "";
            this.achatEnCours = true;
            const erreur = await demarrerAchat("abonnement", 0);
            if (erreur) {
                this.erreurAbonnement = erreur;
                this.achatEnCours = false;
            }
        },

        async resilier() {
            this.resiliationOuverte = false;
            this.erreurAbonnement = "";
            const reponse = await appelerApi("/pro/abonnement/resilier", { method: "POST" });
            if (!reponse.ok) {
                this.erreurAbonnement = await messageErreur(reponse);
                return;
            }
            this.messageAbonnement = t("compte.abonnement.resilie");
            this.chargerAbonnement();
        },
    },
});
