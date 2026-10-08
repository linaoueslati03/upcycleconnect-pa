demarrerApp({
    data() {
        return {
            prenom: "",
            resume: null,
            erreur: "",
            tutorielVisible: false,
            indexEtape: 0,
            // Le contenu des étapes est en base (traductions tutoriel.N.titre / tutoriel.N.texte)
            etapesTutoriel: ["tutoriel.1", "tutoriel.2", "tutoriel.3", "tutoriel.4", "tutoriel.5", "tutoriel.6"],
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
            return statut ? t("depot." + statut) : t("accueil.aucun_depot");
        },

        async charger() {
            const [reponseCompte, reponseResume] = await Promise.all([
                appelerApi("/moi"),
                appelerApi("/moi/dashboard"),
            ]);
            if (!reponseCompte.ok || !reponseResume.ok) {
                this.erreur = t("accueil.erreur_chargement");
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
});
