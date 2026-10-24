const STATUTS_ANNONCE = ["en_ligne", "reservee", "cedee"];

const TYPES_OFFRE = ["formation", "atelier", "evenement"];

const ETAPES_DEPOT = ["demande", "validee", "deposee", "recuperee"];

const FILTRES_STATUT = [
    { nom: "commun.tous", valeur: "tous" },
    { nom: "statut.brouillon", valeur: "brouillon" },
    { nom: "statut.en_attente", valeur: "en_attente" },
    { nom: "statut.publie", valeur: "publie" },
    { nom: "statut.refuse", valeur: "refuse" },
    { nom: "statut.annule", valeur: "annule" },
    { nom: "statut.termine", valeur: "termine" },
];

function formaterDate(date) {
    return new Date(date).toLocaleDateString(langueCourante(), { day: "numeric", month: "long", year: "numeric" });
}

function formaterDateHeure(date) {
    return new Date(date).toLocaleString(langueCourante(), {
        weekday: "long", day: "numeric", month: "long", year: "numeric", hour: "2-digit", minute: "2-digit",
    });
}

function formaterPrix(prix) {
    if (!prix) {
        return t("commun.gratuit");
    }
    return prix.toLocaleString(langueCourante(), { style: "currency", currency: "EUR" });
}

function formaterStatut(statut) {
    return t("statut." + statut);
}

async function messageErreur(reponse) {
    const resultat = await reponse.json();
    return resultat.erreur || t("commun.erreur");
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
