const { createApp } = Vue;

createApp({
    data() {
        return {
            projets: [],
            projet: null, // projet ouvert, avec ses étapes
            formulaire: { titre: "", description: "", partage_public: false },
            nouvelleEtape: { description: "", photo_url: "" },
            erreur: "",
        };
    },

    mounted() {
        this.chargerProjets();
    },

    methods: {
        formaterDate,

        async chargerProjets() {
            const reponse = await appelerApi("/moi/projets");
            if (!reponse.ok) {
                this.erreur = "Erreur lors du chargement de vos projets";
                return;
            }
            this.projets = await reponse.json();
        },

        async ouvrirProjet(id) {
            const reponse = await appelerApi(`/projets/${id}`);
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.projet = await reponse.json();
        },

        async creerProjet() {
            const reponse = await appelerApi("/projets", { method: "POST", body: JSON.stringify(this.formulaire) });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            const projet = await reponse.json();
            this.formulaire = { titre: "", description: "", partage_public: false };
            await this.chargerProjets();
            this.ouvrirProjet(projet.id);
        },

        async basculerPartage() {
            const reponse = await appelerApi(`/projets/${this.projet.id}`, {
                method: "PUT",
                body: JSON.stringify({
                    titre: this.projet.titre,
                    description: this.projet.description,
                    partage_public: !this.projet.partage_public,
                }),
            });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            await this.chargerProjets();
            this.ouvrirProjet(this.projet.id);
        },

        async ajouterEtape() {
            const reponse = await appelerApi(`/projets/${this.projet.id}/etapes`, {
                method: "POST",
                body: JSON.stringify({
                    description: this.nouvelleEtape.description,
                    photo_url: this.nouvelleEtape.photo_url || null,
                    ordre: this.projet.etapes.length + 1,
                }),
            });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.nouvelleEtape = { description: "", photo_url: "" };
            this.ouvrirProjet(this.projet.id);
        },

        async supprimerProjet() {
            if (!confirm("Supprimer ce projet et ses étapes ?")) {
                return;
            }
            const reponse = await appelerApi(`/projets/${this.projet.id}`, { method: "DELETE" });
            if (!reponse.ok) {
                this.erreur = await messageErreur(reponse);
                return;
            }
            this.projet = null;
            this.chargerProjets();
        },
    },
}).mount("#app");
