demarrerApp({
    data() {
        return {
            codeObjet: "",
            objetTrouve: null,
            erreur: "",
            succes: ""
        };
    },
    methods: {
        async verifierCode() {
            this.erreur = "";
            this.succes = "";
            this.objetTrouve = null;

            let id = this.codeObjet.replace(/\D/g, "");
            if (!id) id = this.codeObjet;

            const reponse = await appelerApi(`/depots/${id}`);
            if (!reponse.ok) {
                this.erreur = "Objet introuvable ou code invalide.";
                return;
            }

            const depot = await reponse.json();
            
            if (depot.statut === "recupere") {
                this.erreur = "Cet objet a déjà été récupéré.";
                return;
            }

            this.objetTrouve = depot;
        },

        async confirmerRecuperation() {
            this.erreur = "";
            const reponse = await appelerApi(`/depots/${this.objetTrouve.id}/statut`, {
                method: "PUT",
                body: JSON.stringify({ statut: "recupere" })
            });

            if (!reponse.ok) {
                this.erreur = "Erreur lors de la récupération.";
                return;
            }

            this.succes = "Récupération confirmée avec succès ! Le conteneur est désormais vide.";
            
            setTimeout(() => {
                this.objetTrouve = null;
                this.codeObjet = "";
                this.succes = "";
            }, 3000);
        }
    }
});
