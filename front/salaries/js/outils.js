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

function formaterDate(date) {
    return new Date(date).toLocaleDateString(langueCourante(), {
        day: "numeric",
        month: "long",
        year: "numeric",
    });
}

function versChampDate(dateAPI) {
    if (!dateAPI) {
        return "";
    }
    const date = new Date(dateAPI);
    const decalage = date.getTimezoneOffset() * 60000;
    return new Date(date - decalage).toISOString().slice(0, 16);
}

function versDateAPI(valeurChamp) {
    if (!valeurChamp) {
        return null;
    }
    return new Date(valeurChamp).toISOString();
}
