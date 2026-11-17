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

    const pageCourante = location.pathname.split("/").pop();
    document.querySelectorAll("aside a.nav-button").forEach((lien) => {
        const actif = lien.getAttribute("href") === pageCourante;
        lien.classList.toggle("actif", actif);
        if (actif) {
            lien.setAttribute("aria-current", "page");
        }
    });

    const premierLien = document.querySelector("aside li:not([hidden]) a");
    const logo = document.getElementById("toggle-button");
    if (premierLien && logo) {
        logo.href = premierLien.getAttribute("href");
    }
}

function oublierConnexion() {
    localStorage.removeItem("token");
    localStorage.removeItem("user_id");
    localStorage.removeItem("role");
    localStorage.removeItem("changer_mot_de_passe");
}

function pageMonCompte() {
    return roleConnecte() === "administrateur" ? "/admin/mon-compte.php" : "/espace/mon-compte.php";
}

function doitChangerMotDePasse() {
    return localStorage.getItem("changer_mot_de_passe") === "1";
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

function traductionsEnCache(code) {
    try {
        return JSON.parse(localStorage.getItem(`traductions_${code}`));
    } catch (erreur) {
        return null;
    }
}

async function telechargerMessages(code) {
    try {
        const reponse = await fetch(`${LOCALES_URL}/${code}.json`, { cache: "no-cache" });
        if (!reponse.ok) {
            return false;
        }
        const messages = await reponse.json();
        i18n.global.setLocaleMessage(code, messages);
        try {
            localStorage.setItem(`traductions_${code}`, JSON.stringify(messages));
        } catch (erreur) {
        }
        return true;
    } catch (erreur) {
        return false;
    }
}

// Les traductions déjà connues sont utilisées tout de suite (la page s'affiche sans attendre le réseau),
// puis rafraîchies en arrière-plan : une traduction modifiée apparaît dès qu'elle est téléchargée.
async function chargerMessages(code) {
    const enCache = traductionsEnCache(code);
    if (enCache) {
        i18n.global.setLocaleMessage(code, enCache);
        telechargerMessages(code).then((misAJour) => {
            if (misAJour) {
                traduireElementsStatiques();
            }
        });
        return;
    }
    await telechargerMessages(code);
}

async function chargerTraductions() {
    const codes = new Set([LANGUE_PAR_DEFAUT, langueCourante()]);
    await Promise.all([...codes].map(chargerMessages));
    traduireElementsStatiques();
}

function traduireElementsStatiques() {
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

async function preparerSelecteurLangue() {
    const selecteur = document.getElementById("selecteur-langue");
    if (!selecteur) {
        return;
    }

    try {
        const reponse = await fetch(`${API_BASE_URL}/langues`);
        if (!reponse.ok) {
            return;
        }
        const langues = await reponse.json();
        selecteur.replaceChildren(...langues.map((langue) => new Option(langue.libelle, langue.code)));
        selecteur.value = langueCourante();
        selecteur.setAttribute("aria-label", t("commun.langue"));
        selecteur.addEventListener("change", () => choisirLangue(selecteur.value));
        selecteur.hidden = false;
    } catch (erreur) {
    }
}

async function choisirLangue(code) {
    if (localStorage.getItem("token")) {
        await appelerApi("/moi/langue", { method: "PUT", body: JSON.stringify({ code }) });
    }
    changerLangue(code);
}

async function demarrerApp(options) {
    afficherMenuDuRole();
    await chargerTraductions();
    preparerSelecteurLangue();
    const app = Vue.createApp(options);
    app.use(i18n);
    app.config.globalProperties.langueCourante = langueCourante;
    app.config.globalProperties.roleConnecte = roleConnecte;
    if (typeof FenetreModale !== "undefined") {
        app.component("fenetre-modale", FenetreModale);
    }
    app.mount("#app");
    document.body.classList.add("app-prete");
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
        return;
    }
    if (doitChangerMotDePasse() && window.location.pathname !== pageMonCompte()) {
        window.location.href = pageMonCompte();
    }
}

async function deconnecter() {
    try {
        await appelerApi("/logout", { method: "POST" });
    } catch (erreur) {
    }
    if (typeof deconnecterOneSignal === "function") {
        deconnecterOneSignal();
    }
    oublierConnexion();
    window.location.href = PAGE_CONNEXION;
}