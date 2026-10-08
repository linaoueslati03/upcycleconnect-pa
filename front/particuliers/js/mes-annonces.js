const ANNONCE_VIDE = {
    id: null, titre: "", description: "", type: "don", prix: null,
    categorie_id: null, localisation: "", statut: "en_ligne",
};

demarrerApp({
    data() {
        return {
            annonces: [],
            categories: [],
            statutsAnnonce: STATUTS_ANNONCE,
            formulaire: { ...ANNONCE_VIDE },
            erreur: "",
            erreurListe: "",
        };
    },

    async mounted() {
        const reponse = await appelerApi("/categories");
        if (reponse.ok) {
            this.categories = await reponse.json();
        }
        this.chargerAnnonces();
    },

    methods: {
        formaterPrix,

        async chargerAnnonces() {
            const reponse = await appelerApi("/annonces?mine=true");
            if (!reponse.ok) {
                this.erreurListe = t("annonces.erreur_chargement_miennes");
                return;
            }
            this.annonces = await reponse.json();
        },

        reinitialiserFormulaire() {
            this.formulaire = { ...ANNONCE_VIDE };
            this.erreur = "";
        },

        modifier(annonce) {
            this.formulaire = { ...annonce };
            this.erreur = "";
        },

        // L'API attend l'annonce complète (PUT remplace toutes les valeurs)
        corpsAnnonce(annonce) {
            return JSON.stringify({
                titre: annonce.titre,
                description: annonce.description,
                type: annonce.type,
                prix: annonce.type === "vente" ? annonce.prix : null,
                categorie_id: annonce.categorie_id,
                localisation: annonce.localisation,
                statut: annonce.statut,
            });
        },

        async enregistrer() {
            const id = this.formulaire.id;
            const reponse = await appelerApi(id ? `/annonces/${id}` : "/annonces", {
                method: id ? "PUT" : "POST",
                body: this.corpsAnnonce(this.formulaire),
            });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.reinitialiserFormulaire();
            this.chargerAnnonces();
        },

        async marquerCedee(annonce) {
            const reponse = await appelerApi(`/annonces/${annonce.id}`, {
                method: "PUT",
                body: this.corpsAnnonce({ ...annonce, statut: "cedee" }),
            });
            if (!reponse.ok) {
                this.erreurListe = await messageErreur(reponse);
                return;
            }
            this.chargerAnnonces();
        },

        async supprimer(id) {
            if (!confirm(t("annonces.confirmer_suppression"))) {
                return;
            }
            const reponse = await appelerApi(`/annonces/${id}`, { method: "DELETE" });
            if (!reponse.ok) {
                this.erreurListe = await messageErreur(reponse);
                return;
            }
            this.chargerAnnonces();
        },
    },
});
