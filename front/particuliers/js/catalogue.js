demarrerApp({
    data() {
        return {
            offres: [],
            inscriptions: [],
            typesOffre: TYPES_OFFRE,
            typeFiltre: "",
            erreur: "",
            message: "",
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

        // Mon planning = la liste de mes inscriptions : sert à afficher « inscrit·e »
        async chargerInscriptions() {
            const reponse = await appelerApi("/moi/planning");
            if (reponse.ok) {
                this.inscriptions = await reponse.json();
            }
        },

        estInscrit(offre) {
            return this.inscriptions.some((i) => i.type === offre.type && i.offre_id === offre.id);
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
