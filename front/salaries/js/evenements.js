const { createApp } = Vue;

createApp({
    data() {
        return {
            evenements: [],
            filtres: FILTRES_STATUT,
            statutSelectionne: "tous",
            erreur: "",
        };
    },

    computed: {
        evenementsFiltres() {
            if (this.statutSelectionne === "tous") {
                return this.evenements;
            }
            return this.evenements.filter(e => e.statut === this.statutSelectionne);
        },
    },

    mounted() {
        this.chargerEvenements();
    },

    methods: {
        formaterDate,
        formaterStatut,

        async chargerEvenements() {
            const reponse = await appelerApi("/evenements");
            if (!reponse.ok) {
                this.erreur = "Erreur lors du chargement des événements";
                return;
            }
            this.evenements = await reponse.json();
        },

        modifierEvenement(id) {
            window.location.href = `creer-evenement.php?id=${id}`;
        },

        // L'API vérifie que le salarié connecté est responsable (sinon 403, affiché ici)
        async valider(id) {
            const reponse = await appelerApi(`/evenements/${id}/valider`, { method: "POST" });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.chargerEvenements();
        },

        async supprimer(id) {
            if (!confirm("Supprimer définitivement ?")) {
                return;
            }
            const reponse = await appelerApi(`/evenements/${id}`, { method: "DELETE" });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.chargerEvenements();
        },
    },
}).mount("#app");
