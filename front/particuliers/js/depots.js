const { createApp } = Vue;

createApp({
    data() {
        return {
            depots: [],
            conteneurs: [],
            etapes: ETAPES_DEPOT,
            formulaire: { description_objet: "", conteneur_id: null },
            erreur: "",
            erreurListe: "",
            confirmation: false,
        };
    },

    async mounted() {
        const reponse = await appelerApi("/conteneurs");
        if (reponse.ok) {
            this.conteneurs = await reponse.json();
        }
        this.chargerDepots();
    },

    methods: {
        formaterDate,

        nomConteneur(id) {
            const conteneur = this.conteneurs.find((c) => c.id === id);
            return conteneur ? conteneur.site : "—";
        },

        indexStatut(statut) {
            return ETAPES_DEPOT.findIndex((e) => e.statut === statut);
        },

        async chargerDepots() {
            const reponse = await appelerApi("/moi/depots");
            if (!reponse.ok) {
                this.erreurListe = "Erreur lors du chargement de vos dépôts";
                return;
            }
            this.depots = await reponse.json();
        },

        async demanderDepot() {
            this.erreur = "";
            this.confirmation = false;

            const reponse = await appelerApi("/depots", {
                method: "POST",
                body: JSON.stringify(this.formulaire),
            });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }

            this.confirmation = true;
            this.formulaire = { description_objet: "", conteneur_id: null };
            this.chargerDepots();
        },
    },
}).mount("#app");
