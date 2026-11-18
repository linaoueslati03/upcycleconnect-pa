demarrerApp({
    data() {
        return {
            offres: [],
            inscriptions: [],
            typesOffre: TYPES_OFFRE,
            typeFiltre: "",
            erreur: "",
            message: "",
            achatEnCours: null,
        };
    },

    mounted() {
        this.chargerCatalogue();
        this.chargerInscriptions();
    },

    methods: {
        formaterDateHeure,
        formaterPrix,

        async chargerCatalogue() {
            const chemin = this.typeFiltre ? `/catalogue?type=${this.typeFiltre}` : "/catalogue";
            const reponse = await appelerApi(chemin);
            if (!reponse.ok) {
                this.erreur = t("catalogue.erreur_chargement");
                return;
            }
            this.offres = await reponse.json();
        },

        async chargerInscriptions() {
            const reponse = await appelerApi("/moi/planning");
            if (reponse.ok) {
                this.inscriptions = await reponse.json();
            }
        },

        estInscrit(offre) {
            return this.inscriptions.some((i) => i.type === offre.type && i.offre_id === offre.id);
        },

        // Formation payante : passage par la page de paiement, l'inscription est faite une fois le paiement confirmé
        async payerEtInscrire(offre) {
            this.erreur = "";
            this.message = "";
            this.achatEnCours = offre.type + offre.id;
            const erreur = await demarrerAchat(offre.type, offre.id);
            if (erreur) {
                this.erreur = erreur;
                this.achatEnCours = null;
            }
        },

        async inscrire(offre) {
            this.erreur = "";
            this.message = "";

            const reponse = await appelerApi("/inscriptions", {
                method: "POST",
                body: JSON.stringify({ type_offre: offre.type, offre_id: offre.id }),
            });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }

            this.message = t("catalogue.inscription_confirmee", { titre: offre.titre });
            this.chargerCatalogue();
            this.chargerInscriptions();
        },
    },
});
