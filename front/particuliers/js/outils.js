// Fonctions et libellés communs aux pages de l'espace Particulier.

const LIBELLES_ANNONCE = { en_ligne: "En ligne", reservee: "Réservée", cedee: "Cédée" };

const LIBELLES_OFFRE = { formation: "Formation", atelier: "Atelier", evenement: "Événement" };

// Étapes d'un dépôt en conteneur, dans l'ordre imposé par l'API.
const ETAPES_DEPOT = [
    { statut: "demande", libelle: "Demande envoyée" },
    { statut: "validee", libelle: "Validée" },
    { statut: "deposee", libelle: "Objet déposé" },
    { statut: "recuperee", libelle: "Récupéré par un professionnel" },
];

// "2026-10-20T14:00:00Z" → "20 octobre 2026"
function formaterDate(date) {
    return new Date(date).toLocaleDateString("fr-FR", { day: "numeric", month: "long", year: "numeric" });
}

// "2026-10-20T14:00:00Z" → "mardi 20 octobre 2026 à 16:00" (heure locale du navigateur)
function formaterDateHeure(date) {
    return new Date(date).toLocaleString("fr-FR", {
        weekday: "long", day: "numeric", month: "long", year: "numeric", hour: "2-digit", minute: "2-digit",
    });
}

function formaterPrix(prix) {
    if (!prix) {
        return "Gratuit";
    }
    return prix.toLocaleString("fr-FR", { style: "currency", currency: "EUR" });
}

// Lit le message d'erreur renvoyé par l'API ({"erreur": "..."}).
async function messageErreur(reponse) {
    const resultat = await reponse.json();
    return resultat.erreur || "Une erreur est survenue";
}
