const PAGE_CONNEXION = "/connexion.php";

const ESPACES = {
    administrateur: "/admin/utilisateurs.php",
    salarie: "/espace/planning.php",
    particulier: "/espace/tableau-de-bord.php",
    professionnel: "/espace/tableau-de-bord.php",
};

function roleConnecte() {
    return localStorage.getItem("role");
}

function afficherMenuDuRole() {
    const role = roleConnecte();
    const elements = document.querySelectorAll("aside li[data-roles]:not([data-roles=''])");
    elements.forEach((element) => {
        element.hidden = !element.dataset.roles.split(" ").includes(role);
    });

    const premierLien = document.querySelector("aside li:not([hidden]) a");
    const logo = document.getElementById("toggle-button");
    if (premierLien && logo) {
        logo.href = premierLien.getAttribute("href");
    }
}

function oublierConnexion() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
}

function langueCourante() {
    return localStorage.getItem("langue") || LANGUE_PAR_DEFAUT;
}

function changerLangue(code) {
    localStorage.setItem("langue", code);
    window.location.reload();
}

const i18n = VueI18n.createI18n({
    legacy: false,
    locale: langueCourante(),
    fallbackLocale: LANGUE_PAR_DEFAUT,
    messages: {},
    missingWarn: false,
    fallbackWarn: false,
});

const { t } = i18n.global;

async function chargerMessages(code) {
    try {
        const reponse = await fetch(`${LOCALES_URL}/${code}.json`);
        if (reponse.ok) {
            i18n.global.setLocaleMessage(code, await reponse.json());
        }
    } catch (erreur) {
    }
}

async function chargerTraductions() {
    const codes = new Set([LANGUE_PAR_DEFAUT, langueCourante()]);
    await Promise.all([...codes].map(chargerMessages));

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
    afficherMenuDuRole();
    await chargerTraductions();
    const app = Vue.createApp(options);
    app.use(i18n);
    app.config.globalProperties.langueCourante = langueCourante;
    app.config.globalProperties.roleConnecte = roleConnecte;
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
    if (!localStorage.getItem("token") || !roles.includes(roleConnecte())) {
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
