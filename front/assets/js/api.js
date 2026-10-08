// Fonctions partagées par toutes les pages pour appeler l'API Go.
// Le token reçu à la connexion est gardé dans le localStorage du navigateur.

const PAGE_CONNEXION = "/connexion.php";

// Page d'accueil de chaque espace, selon le rôle renvoyé par POST /api/login.
const ESPACES = {
    administrateur: "/admin/utilisateurs.php",
    salarie: "/salaries/planning.php",
};

function oublierConnexion() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
}

// Appelle l'API en ajoutant le token de connexion, et renvoie la réponse de fetch.
// Si la session a expiré (401), renvoie vers la page de connexion.
async function appelerApi(chemin, options = {}) {
    const entetes = { "Content-Type": "application/json" };
    const token = localStorage.getItem("token");
    if (token) {
        entetes.Authorization = token;
    }

    const reponse = await fetch(API_BASE_URL + chemin, { ...options, headers: entetes });

    if (reponse.status === 401 && token) {
        oublierConnexion();
        window.location.href = PAGE_CONNEXION;
    }
    return reponse;
}

// Renvoie vers la connexion si l'utilisateur n'est pas connecté avec l'un des rôles donnés.
// C'est un confort d'affichage : les droits sont toujours revérifiés par l'API.
function exigerConnexion(...roles) {
    if (!localStorage.getItem("token") || !roles.includes(localStorage.getItem("role"))) {
        window.location.href = PAGE_CONNEXION;
    }
}

async function deconnecter() {
    await appelerApi("/logout", { method: "POST" });
    oublierConnexion();
    window.location.href = PAGE_CONNEXION;
}
