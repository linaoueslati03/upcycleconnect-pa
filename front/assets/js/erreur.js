// Page d'erreur : bouton « Retour à l'accueil » vers l'espace du rôle connecté, bouton « Page précédente »
// quand il y a un historique, et traduction des textes si la langue choisie n'est pas le français.
// Sans Vue ni API : la page doit s'afficher même quand le reste du site ne fonctionne pas.

(function () {
    const carte = document.querySelector(".carte-erreur");
    const code = carte.dataset.code;
    const role = lireStockage("role");
    const langue = lireStockage("langue") || LANGUE_PAR_DEFAUT;

    function lireStockage(cle) {
        try {
            return localStorage.getItem(cle);
        } catch (erreur) {
            return null;
        }
    }

    const accueil = document.getElementById("retour-accueil");
    if (lireStockage("token") && ESPACES[role]) {
        accueil.href = ESPACES[role];
    }

    const precedente = document.getElementById("page-precedente");
    if (window.history.length > 1) {
        precedente.hidden = false;
        precedente.addEventListener("click", () => window.history.back());
    }

    async function traductions() {
        try {
            const enCache = JSON.parse(lireStockage(`traductions_${langue}`));
            if (enCache && enCache.erreur) {
                return enCache.erreur;
            }
        } catch (erreur) {
        }
        try {
            const reponse = await fetch(`${LOCALES_URL}/${langue}.json`);
            return reponse.ok ? (await reponse.json()).erreur : null;
        } catch (erreur) {
            return null;
        }
    }

    if (langue === LANGUE_PAR_DEFAUT) {
        return;
    }
    traductions().then((textes) => {
        if (!textes) {
            return;
        }
        document.documentElement.lang = langue;
        const propres = textes.codes[code] || textes.codes.autre;
        const valeurs = { ...textes, ...propres };
        document.querySelectorAll("[data-t]").forEach((element) => {
            const texte = valeurs[element.dataset.t];
            if (texte) {
                element.textContent = texte.replace("{code}", carte.querySelector(".code-erreur").textContent);
            }
        });
    });
})();
