// Création d'un compte particulier ou professionnel, puis connexion automatique et arrivée dans son espace.
// Depuis la page d'accueil, /inscription?role=professionnel arrive avec le bon profil déjà choisi.

const ROLES_INSCRIPTION = ["particulier", "professionnel"];

demarrerApp({
    data() {
        const roleDemande = new URLSearchParams(window.location.search).get("role");
        return {
            compte: {
                role: ROLES_INSCRIPTION.includes(roleDemande) ? roleDemande : "particulier",
                nom: "",
                prenom: "",
                email: "",
                mot_de_passe: "",
                langue_preferee_id: null,
            },
            confirmation: "",
            erreur: "",
            envoi: false,
        };
    },

    async mounted() {
        // Le compte garde la langue dans laquelle la personne s'est inscrite
        try {
            const reponse = await appelerApi("/langues");
            if (reponse.ok) {
                const langue = (await reponse.json()).find((l) => l.code === langueCourante());
                this.compte.langue_preferee_id = langue ? langue.id : null;
            }
        } catch (erreur) {
        }
    },

    methods: {
        async creerCompte() {
            this.erreur = "";
            if (this.compte.mot_de_passe !== this.confirmation) {
                this.erreur = t("inscription.mots_de_passe_differents");
                return;
            }
            this.envoi = true;
            oublierConnexion();

            try {
                const reponse = await appelerApi("/comptes", { method: "POST", body: JSON.stringify(this.compte) });
                if (!reponse.ok) {
                    this.erreur = (await reponse.json()).erreur;
                    return;
                }

                const connexion = await appelerApi("/login", {
                    method: "POST",
                    body: JSON.stringify({ email: this.compte.email, mot_de_passe: this.compte.mot_de_passe }),
                });
                const resultat = await connexion.json();
                if (!connexion.ok) {
                    window.location.href = PAGE_CONNEXION;
                    return;
                }
                const espace = ESPACES[resultat.role];
                if (!espace) {
                    this.erreur = t("inscription.compte_cree_espace_indisponible");
                    return;
                }

                enregistrerConnexion(resultat);
                window.location.href = espace;
            } catch (erreur) {
                this.erreur = t("commun.serveur_injoignable");
            } finally {
                this.envoi = false;
            }
        },
    },
});
