
const STATUTS_ANNONCE = ["en_ligne", "reservee", "cedee"];

const TYPES_OFFRE = ["formation", "atelier", "evenement"];

const ETAPES_DEPOT = ["demande", "validee", "deposee", "recuperee"];

// 00 mois année
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

async function messageErreur(reponse) {
    const resultat = await reponse.json();
    return resultat.erreur || t("commun.erreur");
}
