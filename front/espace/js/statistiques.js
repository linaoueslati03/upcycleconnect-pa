demarrerApp({
    data() {
        return {
            profil: { abonnement: "Freemium" },
            filtreAlerte: { materiau: "", localisation: "" },
            messageAlerte: ""
        };
    },
    mounted() {
        this.chargerProfil();
    },
    methods: {
        async chargerProfil() {
            const reponse = await appelerApi("/moi");
            if (reponse.ok) {
                const donnees = await reponse.json();
                this.profil.abonnement = donnees.abonnement || "Freemium";
            }
        },
        passerPremium() {
            alert("Redirection vers la passerelle de paiement Stripe..");
            this.profil.abonnement = "Premium";
        },
        enregistrerAlerte() {
            this.messageAlerte = "Alerte priorisée configurée avec succès !";
            setTimeout(() => {
                this.messageAlerte = "";
                this.filtreAlerte = { materiau: "", localisation: "" };
            }, 4000);
        }
    }
});
