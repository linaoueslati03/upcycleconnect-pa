const jours = [
    "lundi",
    "mardi",
    "mercredi",
    "jeudi",
    "vendredi",
    "samedi",
    "dimanche"
];

let dateActuelle = new Date();

function getNumeroSemaine(date) {
    const copie = new Date(Date.UTC(
        date.getFullYear(),
        date.getMonth(),
        date.getDate()
    ));

    const jour = copie.getUTCDay() || 7;

    copie.setUTCDate(
        copie.getUTCDate() + 4 - jour
    );

    const debutAnnee = new Date(
        Date.UTC(copie.getUTCFullYear(), 0, 1)
    );

    return Math.ceil(
        (((copie - debutAnnee) / 86400000) + 1) / 7
    );
}

function getDebutSemaine(date) {
    const copie = new Date(date);
    const jour = copie.getDay() || 7;

    copie.setDate(
        copie.getDate() - jour + 1
    );

    copie.setHours(0, 0, 0, 0);

    return copie;
}

function getDateDepuisSemaine(semaine, annee) {
    const date = new Date(annee, 0, 4);
    const debut = getDebutSemaine(date);

    debut.setDate(
        debut.getDate() + (semaine - 1) * 7
    );

    return debut;
}

function remplirSelects() {
    const selectSemaine = document.getElementById("semaine");
    const selectAnnee = document.getElementById("annee");

    selectSemaine.innerHTML = "";
    selectAnnee.innerHTML = "";

    for (let i = 1; i <= 53; i++) {
        const option = document.createElement("option");

        const dateSemaine = getDateDepuisSemaine(
            i,
            dateActuelle.getFullYear()
        );

        const jour = String(dateSemaine.getDate()).padStart(2, "0");
        const mois = String(dateSemaine.getMonth() + 1).padStart(2, "0");

        option.value = i;
        option.textContent = `Semaine du ${jour}/${mois}`;

        selectSemaine.appendChild(option);
    }

    const anneeActuelle = dateActuelle.getFullYear();

    for (
        let annee = anneeActuelle - 2;
        annee <= anneeActuelle + 2;
        annee++
    ) {
        const option = document.createElement("option");

        option.value = annee;
        option.textContent = annee;

        selectAnnee.appendChild(option);
    }
}

function afficherPlanning() {
    const semaine = getNumeroSemaine(dateActuelle);
    const annee = dateActuelle.getFullYear();

    const debut = getDebutSemaine(dateActuelle);

    const jour = String(debut.getDate()).padStart(2, "0");
    const mois = String(debut.getMonth() + 1).padStart(2, "0");

    const texteSemaine = `Semaine du ${jour}/${mois}`;

    document.getElementById("semaine").value = semaine;
    document.getElementById("annee").value = annee;

    document.getElementById("titre-semaine").textContent =
        texteSemaine;

    jours.forEach(function(jour) {
        document.getElementById(jour).innerHTML = "";
    });
}

document
    .getElementById("semaine-precedente")
    .addEventListener("click", function() {
        dateActuelle.setDate(
            dateActuelle.getDate() - 7
        );

        afficherPlanning();
    });

document
    .getElementById("semaine-suivante")
    .addEventListener("click", function() {
        dateActuelle.setDate(
            dateActuelle.getDate() + 7
        );

        afficherPlanning();
    });

document
    .getElementById("semaine")
    .addEventListener("change", function() {
        const semaine = parseInt(this.value);

        const annee = parseInt(
            document.getElementById("annee").value
        );

        dateActuelle = getDateDepuisSemaine(
            semaine,
            annee
        );

        afficherPlanning();
    });

document
    .getElementById("annee")
    .addEventListener("change", function() {
        const annee = parseInt(this.value);

        const semaine = parseInt(
            document.getElementById("semaine").value
        );

        dateActuelle = getDateDepuisSemaine(
            semaine,
            annee
        );

        remplirSelects();
        afficherPlanning();
    });

remplirSelects();
afficherPlanning();