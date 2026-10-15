demarrerApp({
    data() {
        return {
            entrees: [],
            erreur: "",
            participation: "",
            aVenirSeulement: true,
            recherche: "",
            filtres: [
                { nom: "commun.tous", valeur: "" },
                { nom: "admin.plannings.anime", valeur: "anime" },
                { nom: "admin.plannings.inscrit", valeur: "inscrit" },
            ],
        };
    },

    computed: {
        entreesFiltrees() {
            const maintenant = new Date();
            const motCle = this.recherche.toLowerCase();
            return this.entrees.filter((e) =>
                (!this.participation || e.participation === this.participation) &&
                (!this.aVenirSeulement || new Date(e.date_debut) >= maintenant) &&
                e.personne.toLowerCase().includes(motCle)
            );
        },
    },

    async mounted() {
        const reponse = await appelerApi("/admin/plannings");
        if (!reponse.ok) {
            this.erreur = t("admin.plannings.erreur_chargement");
            return;
        }
        this.entrees = await reponse.json();
    },

    methods: {
        formaterDateHeure(date) {
            return new Date(date).toLocaleString(langueCourante(), {
                day: "numeric", month: "long", year: "numeric", hour: "2-digit", minute: "2-digit",
            });
        },
    },
});
