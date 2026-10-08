const { createApp } = Vue;

const FORMULAIRE_VIDE = { id: null, titre: "", categorie: "", tarif: 0, statut: "Brouillon" };

createApp({
    data() {
        return {
            prestations: [],
            formulaire: { ...FORMULAIRE_VIDE },
            erreur: "",
            erreurListe: "",
        };
    },

    mounted() {
        this.chargerPrestations();
    },

    methods: {
        async chargerPrestations() {
            const reponse = await appelerApi("/prestations");
            if (!reponse.ok) {
                this.erreurListe = "Erreur lors du chargement des prestations";
                return;
            }
            this.prestations = await reponse.json();
        },

        formaterTarif(tarif) {
            return tarif.toLocaleString("fr-FR", { style: "currency", currency: "EUR" });
        },

        // Le récapitulatif PDF est généré par l'API Go (route publique GET /prestations/{id}/pdf)
        lienPDF(id) {
            return `${API_BASE_URL}/prestations/${id}/pdf`;
        },

        reinitialiserFormulaire() {
            this.formulaire = { ...FORMULAIRE_VIDE };
            this.erreur = "";
        },

        modifier(prestation) {
            this.formulaire = { ...prestation };
            this.erreur = "";
        },

        async enregistrer() {
            const { id, ...donnees } = this.formulaire;
            const chemin = id ? `/prestations/${id}` : "/prestations";
            const methode = id ? "PUT" : "POST";

            const reponse = await appelerApi(chemin, { method: methode, body: JSON.stringify(donnees) });
            if (!reponse.ok) {
                const resultat = await reponse.json();
                this.erreur = resultat.erreur;
                return;
            }

            this.reinitialiserFormulaire();
            this.chargerPrestations();
        },

        async supprimer(id) {
            if (!confirm("Supprimer cette prestation ?")) {
                return;
            }

            const reponse = await appelerApi(`/prestations/${id}`, { method: "DELETE" });
            if (!reponse.ok) {
                const resultat = await reponse.json();
                this.erreurListe = resultat.erreur;
                return;
            }

            this.chargerPrestations();
        },
    },
}).mount("#app");
