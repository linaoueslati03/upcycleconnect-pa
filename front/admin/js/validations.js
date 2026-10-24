demarrerApp({
    data() {
        return {
            offres: [],
            filtreType: "",
            filtres: [
                { nom: "commun.tous", valeur: "" },
                { nom: "admin.validations.formations", valeur: "formation" },
                { nom: "admin.validations.ateliers", valeur: "atelier" },
                { nom: "admin.validations.evenements", valeur: "evenement" },
            ],
            refus: { cle: "", motif: "" },
            erreur: "",
            succes: "",
        };
    },

    computed: {
        offresFiltrees() {
            if (!this.filtreType) {
                return this.offres;
            }
            return this.offres.filter((o) => o.type === this.filtreType);
        },
    },

    mounted() {
        this.charger();
    },

    methods: {
        formaterDate(date) {
            return new Date(date).toLocaleDateString(langueCourante(), { day: "numeric", month: "long", year: "numeric" });
        },

        formaterDateHeure(date) {
            return new Date(date).toLocaleString(langueCourante(), {
                day: "numeric", month: "long", year: "numeric", hour: "2-digit", minute: "2-digit",
            });
        },

        compter(type) {
            if (!type) {
                return this.offres.length;
            }
            return this.offres.filter((o) => o.type === type).length;
        },

        async charger() {
            const reponse = await appelerApi("/admin/validations");
            if (!reponse.ok) {
                this.erreur = t("admin.validations.erreur_chargement");
                return;
            }
            this.offres = await reponse.json();
        },

        commencerRefus(offre) {
            this.refus = { cle: offre.type + offre.id, motif: "" };
        },

        annulerRefus() {
            this.refus = { cle: "", motif: "" };
        },

        async decider(offre, decision) {
            this.erreur = "";
            this.succes = "";
            const reponse = await appelerApi(`/validations/${offre.type}/${offre.id}`, {
                method: "PUT",
                body: JSON.stringify({ decision, motif: decision === "refuser" ? this.refus.motif : "" }),
            });
            if (!reponse.ok) {
                this.erreur = (await reponse.json()).erreur;
                return;
            }
            const message = decision === "valider" ? "admin.validations.validee" : "admin.validations.refusee";
            this.succes = t(message, { titre: offre.titre });
            this.annulerRefus();
            this.charger();
        },
    },
});
