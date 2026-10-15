const PAGE_CONNEXION = "/connexion.php";

const ESPACES = {
    administrateur: "/admin/utilisateurs.php",
    salarie: "/salaries/planning.php",
    particulier: "/particuliers/tableau-de-bord.php",
    professionnel: "/professionnels/tableau-de-bord.php",
};

function oublierConnexion() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
}

let traductions = {};

function langueCourante() {
    return localStorage.getItem("langue") || "fr";
}

function changerLangue(code) {
    localStorage.setItem("langue", code);
    window.location.reload();
}

function t(cle, valeurs = {}) {
    let texte = traductions[cle] || cle;
    for (const [nom, valeur] of Object.entries(valeurs)) {
        texte = texte.replace(`{${nom}}`, valeur);
    }
    return texte;
}

function traduireOu(cle, valeurParDefaut) {
    return traductions[cle] || valeurParDefaut;
}

async function chargerTraductions() {
    try {
        const reponse = await fetch(`${API_BASE_URL}/traductions?langue=${langueCourante()}`);
        if (reponse.ok) {
            traductions = await reponse.json();
        }
    } catch (erreur) {
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

async function demarrerApp(options) {
    await chargerTraductions();
    const app = Vue.createApp(options);
    app.config.globalProperties.t = t;
    app.config.globalProperties.traduireOu = traduireOu;
    app.config.globalProperties.langueCourante = langueCourante;
    app.mount("#app");
}

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

function exigerConnexion(...roles) {
    if (!localStorage.getItem("token") || !roles.includes(localStorage.getItem("role"))) {
        window.location.href = PAGE_CONNEXION;
    }
}

async function deconnecter() {
    try {
        await appelerApi("/logout", { method: "POST" });
    } catch (erreur) {
    }
    oublierConnexion();
    window.location.href = PAGE_CONNEXION;
}
