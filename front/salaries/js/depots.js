demarrerApp({
    data() {
        return {
            depots: [],
            statutSelectionne: "demande",
            erreur: "",
            filtres: [
                { nom: "salaries.depots.a_valider", valeur: "demande" },
                { nom: "salaries.depots.valides", valeur: "validee" },
                { nom: "salaries.depots.deposes", valeur: "deposee" },
                { nom: "commun.tous", valeur: "" },
            ],
            // Le salarié valide la demande puis confirme le dépôt ; la récupération est
            // enregistrée par le professionnel qui scanne le code-barre (espace Pro).
            statutSuivant: { demande: "validee", validee: "deposee" },
            libellesAction: { demande: "salaries.depots.valider_demande", validee: "salaries.depots.confirmer_depot" },
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
                this.erreur = t("salaries.depots.erreur_chargement");
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
});
