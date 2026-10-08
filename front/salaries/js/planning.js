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

// Nombre de semaines ISO de l'année (52 ou 53) : le 28 décembre est toujours dans la dernière.
function nombreSemaines(annee) {
    return numeroSemaine(new Date(annee, 11, 28));
}

// Année ISO d'une date : celle du jeudi de sa semaine (le 1er janvier peut appartenir à la
// dernière semaine de l'année précédente).
function anneeISO(date) {
    const jeudi = lundiDeLaSemaine(date);
    jeudi.setDate(jeudi.getDate() + 3);
    return jeudi.getFullYear();
}

const aujourdHui = new Date();

demarrerApp({
    data() {
        const annees = [];
        for (let a = aujourdHui.getFullYear() - 2; a <= aujourdHui.getFullYear() + 2; a++) {
            annees.push(a);
        }

        return {
            jours: ["jour.1", "jour.2", "jour.3", "jour.4", "jour.5", "jour.6", "jour.7"], // clés de traduction
            annees: annees,
            semaineSelectionnee: numeroSemaine(aujourdHui),
            anneeSelectionnee: anneeISO(aujourdHui),
            entrees: [], // événements et ateliers du salarié connecté
            erreur: "",
        };
    },

    computed: {
        semaines() {
            const resultat = [];
            for (let i = 1; i <= nombreSemaines(this.anneeSelectionnee); i++) {
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
            this.erreur = t("salaries.planning.erreur_chargement");
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
            return new Date(date).toLocaleTimeString(langueCourante(), { hour: "2-digit", minute: "2-digit" });
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
                this.semaineSelectionnee = nombreSemaines(this.anneeSelectionnee);
            }
        },

        semaineSuivante() {
            if (this.semaineSelectionnee < nombreSemaines(this.anneeSelectionnee)) {
                this.semaineSelectionnee++;
            } else {
                this.anneeSelectionnee++;
                this.semaineSelectionnee = 1;
            }
        },
    },
});
