const { createApp } = Vue;

// Numéro de semaine ISO (la semaine 1 est celle qui contient le premier jeudi de l'année).
function numeroSemaine(date) {
    const copie = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
    const jour = copie.getUTCDay() || 7;
    copie.setUTCDate(copie.getUTCDate() + 4 - jour);
    const debutAnnee = new Date(Date.UTC(copie.getUTCFullYear(), 0, 1));
    return Math.ceil(((copie - debutAnnee) / 86400000 + 1) / 7);
}

// Lundi de la semaine qui contient la date donnée.
function lundiDeLaSemaine(date) {
    const copie = new Date(date);
    const jour = copie.getDay() || 7;
    copie.setDate(copie.getDate() - jour + 1);
    copie.setHours(0, 0, 0, 0);
    return copie;
}

// Lundi de la semaine n°semaine de l'année (le 4 janvier est toujours en semaine 1).
function dateDepuisSemaine(semaine, annee) {
    const debut = lundiDeLaSemaine(new Date(annee, 0, 4));
    debut.setDate(debut.getDate() + (semaine - 1) * 7);
    return debut;
}

const aujourdHui = new Date();

createApp({
    data() {
        const annees = [];
        for (let a = aujourdHui.getFullYear() - 2; a <= aujourdHui.getFullYear() + 2; a++) {
            annees.push(a);
        }

        return {
            jours: ["Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi", "Dimanche"],
            annees: annees,
            semaineSelectionnee: numeroSemaine(aujourdHui),
            anneeSelectionnee: aujourdHui.getFullYear(),
            entrees: [], // événements et ateliers du salarié connecté
            erreur: "",
        };
    },

    computed: {
        semaines() {
            const resultat = [];
            for (let i = 1; i <= 53; i++) {
                resultat.push({ numero: i, date: dateDepuisSemaine(i, this.anneeSelectionnee) });
            }
            return resultat;
        },

        debutSemaine() {
            return dateDepuisSemaine(this.semaineSelectionnee, this.anneeSelectionnee);
        },
    },

    async mounted() {
        const reponse = await appelerApi("/salaries/planning");
        if (!reponse.ok) {
            this.erreur = "Erreur lors du chargement du planning";
            return;
        }
        this.entrees = await reponse.json();
    },

    methods: {
        formaterStatut,

        // Entrées du jour n°index (0 = lundi) de la semaine affichée
        entreesDuJour(index) {
            const jour = new Date(this.debutSemaine);
            jour.setDate(jour.getDate() + index);
            return this.entrees.filter((e) => new Date(e.date_debut).toDateString() === jour.toDateString());
        },

        // "2026-10-20T14:00:00Z" → "16:00" (heure locale)
        formaterHeure(date) {
            return new Date(date).toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" });
        },

        // Date → "20/10"
        formaterJourMois(date) {
            const jour = String(date.getDate()).padStart(2, "0");
            const mois = String(date.getMonth() + 1).padStart(2, "0");
            return `${jour}/${mois}`;
        },

        semainePrecedente() {
            if (this.semaineSelectionnee > 1) {
                this.semaineSelectionnee--;
            } else {
                this.anneeSelectionnee--;
                this.semaineSelectionnee = 52;
            }
        },

        semaineSuivante() {
            if (this.semaineSelectionnee < 53) {
                this.semaineSelectionnee++;
            } else {
                this.anneeSelectionnee++;
                this.semaineSelectionnee = 1;
            }
        },
    },
}).mount("#app");
