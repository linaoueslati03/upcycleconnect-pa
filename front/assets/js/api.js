// Fonctions partagées par toutes les pages pour appeler l'API Go.
// Le token reçu à la connexion est gardé dans le localStorage du navigateur.

const PAGE_CONNEXION = "/connexion.php";

// Page d'accueil de chaque espace, selon le rôle renvoyé par POST /api/login.
const ESPACES = {
    administrateur: "/admin/utilisateurs.php",
    salarie: "/salaries/planning.php",
    particulier: "/particuliers/tableau-de-bord.php",
};

function oublierConnexion() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
}

// ---------- Multilingue
// Les textes de l'interface sont en base (table traductions) et chargés au démarrage
// de chaque page : ajouter une langue ne demande aucune modification du code.

let traductions = {};

function langueCourante() {
    return localStorage.getItem("langue") || "fr";
}

function changerLangue(code) {
    localStorage.setItem("langue", code);
    window.location.reload();
}

// Texte de la clé dans la langue courante. {nom} est remplacé par valeurs.nom.
// Si la clé n'existe pas, on affiche la clé elle-même (repérable pendant le développement).
function t(cle, valeurs = {}) {
    let texte = traductions[cle] || cle;
    for (const [nom, valeur] of Object.entries(valeurs)) {
        texte = texte.replace(`{${nom}}`, valeur);
    }
    return texte;
}

// Comme t(), mais pour une donnée venant de la base (ex. code de catégorie) : si aucune
// traduction n'existe encore pour cette valeur, on affiche la valeur elle-même.
function traduireOu(cle, valeurParDefaut) {
    return traductions[cle] || valeurParDefaut;
}

// Charge les traductions puis traduit les textes écrits en PHP hors de Vue
// (menu, titre de l'onglet), repérés par l'attribut data-t="clé".
async function chargerTraductions() {
    try {
        const reponse = await fetch(`${API_BASE_URL}/traductions?langue=${langueCourante()}`);
        if (reponse.ok) {
            traductions = await reponse.json();
        }
    } catch (erreur) {
        // API injoignable : les clés s'affichent, la page reste utilisable
    }

    document.documentElement.lang = langueCourante();
    document.querySelectorAll("[data-t]").forEach((element) => {
        element.textContent = t(element.dataset.t);
    });
    document.querySelectorAll("[data-t-alt]").forEach((element) => {
        element.alt = t(element.dataset.tAlt);
    });
    const titre = document.querySelector("title[data-titre]");
    if (titre) {
        document.title = `${t(titre.dataset.titre)} — UpcycleConnect`;
    }
}

// Démarre l'application Vue de la page une fois les traductions chargées :
// t() est alors utilisable dans tous les templates ({{ t("cle") }}).
async function demarrerApp(options) {
    await chargerTraductions();
    const app = Vue.createApp(options);
    app.config.globalProperties.t = t;
    app.config.globalProperties.traduireOu = traduireOu;
    app.config.globalProperties.langueCourante = langueCourante;
    app.mount("#app");
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
