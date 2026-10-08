demarrerApp({
    data() {
        return {
            annonces: [],
            statutSelectionne: "en_attente",
            erreur: "",
            filtres: [
                { nom: "admin.annonces.a_valider", valeur: "en_attente" },
                { nom: "annonce.en_ligne", valeur: "en_ligne" },
                { nom: "admin.annonces.refusees", valeur: "refusee" },
                { nom: "commun.tous", valeur: "tous" },
            ],
        };
    },

    mounted() {
        this.chargerAnnonces();
    },

    methods: {
        formaterPrix(prix) {
            return prix.toLocaleString(langueCourante(), { style: "currency", currency: "EUR" });
        },

        async chargerAnnonces() {
            const reponse = await appelerApi(`/admin/annonces?statut=${this.statutSelectionne}`);
            if (!reponse.ok) {
                this.erreur = t("admin.annonces.erreur_chargement");
                return;
            }
            this.erreur = "";
            this.annonces = await reponse.json();
        },

        // decision = "valider" (l'annonce passe en ligne) ou "refuser"
        async decider(id, decision) {
            const reponse = await appelerApi(`/annonces/${id}/${decision}`, { method: "POST" });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            this.chargerAnnonces();
        },
    },
});
