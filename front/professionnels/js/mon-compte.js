demarrerApp({
    data() {
        return {
            profil: { nom: "", siret: "", adresse: "", abonnement: "Freemium" },
            succesProfil: false,
            erreurProfil: ""
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
                this.profil.nom = donnees.nom || "";
                this.profil.adresse = donnees.adresse || "";
                this.profil.siret = donnees.siret || "Non renseigné";
                this.profil.abonnement = donnees.abonnement || "Freemium";
            }
        },
        async enregistrerProfil() {
            this.erreurProfil = "";
            this.succesProfil = false;
            const reponse = await appelerApi("/moi", {
                method: "PUT",
                body: JSON.stringify({ nom: this.profil.nom, adresse: this.profil.adresse })
            });

            if (!reponse.ok) {
                const res = await reponse.json();
                this.erreurProfil = res.erreur || "Erreur de sauvegarde.";
                return;
            }
            this.succesProfil = true;
            setTimeout(() => { this.succesProfil = false; }, 3000);
        },
        passerPremium() {
            alert("Redirection vers la passerelle de paiement Stripe...");
            this.profil.abonnement = "Premium";
        }
    }
});
