// Fonctions communes aux pages de l'espace Salarié.

const NOMS_STATUTS = {
    brouillon: "Brouillon",
    en_attente: "En attente",
    publie: "Publié",
    annule: "Annulé",
    termine: "Terminé",
};

// Filtres affichés au-dessus des listes d'événements et d'ateliers.
const FILTRES_STATUT = [
    { nom: "Tous", valeur: "tous" },
    { nom: "Brouillon", valeur: "brouillon" },
    { nom: "En attente", valeur: "en_attente" },
    { nom: "Publié", valeur: "publie" },
    { nom: "Annulé", valeur: "annule" },
    { nom: "Terminé", valeur: "termine" },
];

function formaterStatut(statut) {
    return NOMS_STATUTS[statut] || statut;
}

// "2026-10-20T14:00:00Z" → "20 octobre 2026"
function formaterDate(date) {
    return new Date(date).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
        year: "numeric",
    });
}

// Un champ <input type="datetime-local"> donne "2026-10-20T14:00",
// alors que l'API Go attend une date complète ("2026-10-20T14:00:00.000Z").
function versDateAPI(valeurChamp) {
    if (!valeurChamp) {
        return null;
    }
    return new Date(valeurChamp).toISOString();
}
