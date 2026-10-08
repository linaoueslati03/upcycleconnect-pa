// Fonctions communes aux pages de l'espace Salarié.

// Filtres affichés au-dessus des listes d'événements et d'ateliers.
// nom = clé de traduction du libellé
const FILTRES_STATUT = [
    { nom: "commun.tous", valeur: "tous" },
    { nom: "statut.brouillon", valeur: "brouillon" },
    { nom: "statut.en_attente", valeur: "en_attente" },
    { nom: "statut.publie", valeur: "publie" },
    { nom: "statut.annule", valeur: "annule" },
    { nom: "statut.termine", valeur: "termine" },
];

function formaterStatut(statut) {
    return t("statut." + statut);
}

// "2026-10-20T14:00:00Z" → "20 octobre 2026"
function formaterDate(date) {
    return new Date(date).toLocaleDateString(langueCourante(), {
        day: "numeric",
        month: "long",
        year: "numeric",
    });
}

// Inverse de versDateAPI : pré-remplit un <input type="datetime-local"> en heure locale.
// "2026-10-20T14:00:00Z" → "2026-10-20T16:00" (à Paris, en heure d'été)
function versChampDate(dateAPI) {
    if (!dateAPI) {
        return "";
    }
    const date = new Date(dateAPI);
    const decalage = date.getTimezoneOffset() * 60000;
    return new Date(date - decalage).toISOString().slice(0, 16);
}

// Un champ <input type="datetime-local"> donne "2026-10-20T14:00",
// alors que l'API Go attend une date complète ("2026-10-20T14:00:00.000Z").
function versDateAPI(valeurChamp) {
    if (!valeurChamp) {
        return null;
    }
    return new Date(valeurChamp).toISOString();
}
