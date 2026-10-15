demarrerApp({
    data() {
        return {
            prestations: [],
            erreur: ""
        };
    },
    mounted() {
        this.chargerPrestations();
    },
    methods: {
        async chargerPrestations() {
            const reponse = await appelerApi("/prestations");
            if (!reponse.ok) {
                this.erreur = t("commun.erreur_chargement") || "Erreur de chargement des prestations.";
                return;
            }
            this.prestations = await reponse.json();
        },
        formaterTarif(tarif) {
            return tarif.toLocaleString(langueCourante(), { style: "currency", currency: "EUR" });
        },
        lienPDF(id) {
            return `${API_BASE_URL}/prestations/${id}/pdf`;
        }
    }
});
