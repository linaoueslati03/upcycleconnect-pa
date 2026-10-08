// Fonctions et libellés communs aux pages de l'espace Particulier.

// Les libellés affichés sont des clés de traduction : t("annonce." + statut), t("offre." + type)…
const STATUTS_ANNONCE = ["en_ligne", "reservee", "cedee"];

const TYPES_OFFRE = ["formation", "atelier", "evenement"];

// Étapes d'un dépôt en conteneur, dans l'ordre imposé par l'API (libellé : t("depot." + statut)).
const ETAPES_DEPOT = ["demande", "validee", "deposee", "recuperee"];

// "2026-10-20T14:00:00Z" → "20 octobre 2026"
function formaterDate(date) {
    return new Date(date).toLocaleDateString(langueCourante(), { day: "numeric", month: "long", year: "numeric" });
}

// "2026-10-20T14:00:00Z" → "mardi 20 octobre 2026 à 16:00" (heure locale, langue de l'utilisateur)
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

// Lit le message d'erreur renvoyé par l'API ({"erreur": "..."}).
async function messageErreur(reponse) {
    const resultat = await reponse.json();
    return resultat.erreur || t("commun.erreur");
}
