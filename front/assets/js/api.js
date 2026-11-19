function roleConnecte() {
    return localStorage.getItem("role");
}

function afficherMenuDuRole() {
    const role = roleConnecte();
    const elements = document.querySelectorAll("aside li[data-roles]:not([data-roles=''])");
    elements.forEach((element) => {
        element.hidden = !element.dataset.roles.split(" ").includes(role);
    });

    // Adresses sans « .php » (réécriture d'URL) : on compare les noms de page sans l'extension
    const pageCourante = location.pathname.split("/").pop().replace(/\.php$/, "");
    document.querySelectorAll("aside a.nav-button").forEach((lien) => {
        const actif = lien.getAttribute("href").replace(/\.php$/, "") === pageCourante;
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

// Après /login (page de connexion ou fin d'inscription) : garde la session dans le navigateur
function enregistrerConnexion(resultat) {
    localStorage.setItem("token", resultat.token);
    localStorage.setItem("user_id", resultat.id);
    localStorage.setItem("role", resultat.role);
    if (resultat.langue) {
        localStorage.setItem("langue", resultat.langue);
    }
    if (resultat.doit_changer_mot_de_passe) {
        localStorage.setItem("changer_mot_de_passe", "1");
    }
}

function pageMonCompte() {
    return roleConnecte() === "administrateur" ? "/admin/mon-compte" : "/espace/mon-compte";
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
    app.config.globalProperties.lienTutoriel = () => (typeof lienTutoriel === "function" ? lienTutoriel() : "");
    app.mount("#app");
    document.body.classList.add("app-prete");
    if (typeof lancerTutorielSiBesoin === "function") {
        lancerTutorielSiBesoin();
    }
    if (typeof clocheNotifications !== "undefined") {
        clocheNotifications.demarrer();
    }
}

// Recherche au fil de la frappe : n'appelle fonction qu'une fois que l'utilisateur
// a cessé de taper pendant « delai » millisecondes (évite une requête par touche).
function differer(fonction, delai = 300) {
    let minuteur;
    return (...parametres) => {
        clearTimeout(minuteur);
        minuteur = setTimeout(() => fonction(...parametres), delai);
    };
}

// PDF archivés : l'API exige le jeton dans l'en-tête, un simple lien <a href> ne suffit donc pas.
// On récupère le fichier avec le jeton, puis on l'ouvre dans un onglet ou on le télécharge.
// L'onglet est ouvert tout de suite (pendant le clic) pour ne pas être bloqué par le navigateur.
async function recupererPdf(chemin, { nomFichier = "document.pdf", telecharger = false } = {}) {
    const onglet = telecharger ? null : window.open("", "_blank");
    try {
        const reponse = await appelerApi(chemin);
        if (!reponse.ok) {
            let message = t("commun.erreur");
            try {
                message = (await reponse.json()).erreur || message;
            } catch (erreur) {
            }
            throw new Error(message);
        }
        const adresse = URL.createObjectURL(new Blob([await reponse.blob()], { type: "application/pdf" }));
        if (onglet) {
            onglet.location.href = adresse;
        } else {
            const lien = Object.assign(document.createElement("a"), { href: adresse, download: nomFichier });
            document.body.append(lien);
            lien.click();
            lien.remove();
        }
        setTimeout(() => URL.revokeObjectURL(adresse), 60000);
    } catch (erreur) {
        onglet?.close();
        throw erreur;
    }
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
    if (!localStorage.getItem("token")) {
        window.location.href = PAGE_CONNEXION;
        return;
    }
    // Connecté, mais la page est réservée à un autre rôle : page « Accès refusé »
    if (!roles.includes(roleConnecte())) {
        window.location.replace("/erreur?code=403");
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