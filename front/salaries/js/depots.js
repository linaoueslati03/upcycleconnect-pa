const { createApp } = Vue;

createApp({
    data() {
        return {
            depots: [],
            statutSelectionne: "demande",
            erreur: "",
            filtres: [
                { nom: "À valider", valeur: "demande" },
                { nom: "Validés", valeur: "validee" },
                { nom: "Déposés", valeur: "deposee" },
                { nom: "Tous", valeur: "" },
            ],
            libellesStatut: { demande: "Demande", validee: "Validée", deposee: "Objet déposé", recuperee: "Récupéré" },
            // Le salarié valide la demande puis confirme le dépôt ; la récupération est
            // enregistrée par le professionnel qui scanne le code-barre (espace Pro).
            statutSuivant: { demande: "validee", validee: "deposee" },
            libellesAction: { demande: "Valider la demande", validee: "Confirmer le dépôt" },
        };
    },

    mounted() {
        this.chargerDepots();
    },

    methods: {
        formaterDate,

        async chargerDepots() {
            const chemin = this.statutSelectionne ? `/depots?statut=${this.statutSelectionne}` : "/depots";
            const reponse = await appelerApi(chemin);
            if (!reponse.ok) {
                this.erreur = "Erreur lors du chargement des dépôts";
                return;
            }
            this.erreur = "";
            this.depots = await reponse.json();
        },

        async avancer(depot) {
            const reponse = await appelerApi(`/depots/${depot.id}/statut`, {
                method: "PUT",
                body: JSON.stringify({ statut: this.statutSuivant[depot.statut] }),
            });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.chargerDepots();
        },
    },
}).mount("#app");
