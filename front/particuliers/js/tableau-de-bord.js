const { createApp } = Vue;

createApp({
    data() {
        return {
            prenom: "",
            resume: null,
            erreur: "",
            tutorielVisible: false,
            indexEtape: 0,
            etapesTutoriel: [
                { titre: "Bienvenue sur UpcycleConnect !", texte: "Voici un rapide tour des fonctionnalités de votre espace." },
                { titre: "Vos annonces", texte: "Donnez ou vendez vos objets et matériaux depuis « Mes annonces », et parcourez celles des autres membres." },
                { titre: "Le dépôt en conteneur", texte: "Pas le temps de gérer une vente ? Déposez votre objet dans un conteneur UpcycleConnect : vous recevez un code d'ouverture une fois la demande validée." },
                { titre: "Le catalogue", texte: "Inscrivez-vous aux formations, ateliers et événements : ils apparaissent ensuite dans « Mon planning »." },
                { titre: "Votre Upcycling Score", texte: "Chaque dépôt, don ou inscription vous rapporte des points. Suivez leur évolution dans « Upcycling Score »." },
                { titre: "La communauté", texte: "Partagez vos projets d'upcycling et échangez avec les autres membres sur les forums." },
            ],
        };
    },

    computed: {
        etapeCourante() {
            return this.etapesTutoriel[this.indexEtape];
        },
        derniereEtape() {
            return this.indexEtape === this.etapesTutoriel.length - 1;
        },
    },

    mounted() {
        this.charger();
    },

    methods: {
        formaterDate,

        libelleDepot(statut) {
            const etape = ETAPES_DEPOT.find((e) => e.statut === statut);
            return etape ? etape.libelle : "Aucun";
        },

        async charger() {
            const [reponseCompte, reponseResume] = await Promise.all([
                appelerApi("/moi"),
                appelerApi("/moi/dashboard"),
            ]);
            if (!reponseCompte.ok || !reponseResume.ok) {
                this.erreur = "Erreur lors du chargement du tableau de bord";
                return;
            }

            const compte = await reponseCompte.json();
            this.prenom = compte.prenom;
            this.resume = await reponseResume.json();

            // ?tutoriel=1 : lien « Revoir le tutoriel » de la page Mon compte
            const revoir = new URLSearchParams(window.location.search).has("tutoriel");
            this.tutorielVisible = !compte.tutoriel_vu || revoir;
        },

        async terminerTutoriel() {
            await appelerApi("/moi/tutoriel", { method: "PUT" });
            this.tutorielVisible = false;
            this.indexEtape = 0;
        },
    },
}).mount("#app");
