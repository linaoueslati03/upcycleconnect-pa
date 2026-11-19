// Tutoriel de première connexion (demandé par le sujet : « affichable à la première connexion,
// bloquer l'utilisateur et utiliser des overlays »).
// L'écran est assombri, une découpe met en lumière l'élément réel dont parle chaque étape,
// et une bulle donne l'explication. Tant que la visite n'est pas terminée, le reste de la page
// ne réagit pas aux clics. À la fin, l'API enregistre que le tutoriel a été vu.
//
// Une étape = { cle } pour une bulle centrée, ou { cle, cible } pour mettre un élément en lumière.
// Les textes sont dans les fichiers de langue : tutoriel.<rôle>.<cle>.titre / .texte

// Les liens du menu sont écrits sans « .php » (réécriture d'URL)
const lienMenu = (page) => `aside a.nav-button[href="${page.replace(/\.php$/, "")}"]`;

const ETAPES_TUTORIEL = {
    particulier: [
        { cle: "bienvenue" },
        { cle: "annonces", cible: lienMenu("mes-annonces.php") },
        { cle: "depots", cible: lienMenu("depots.php") },
        { cle: "catalogue", cible: lienMenu("catalogue.php") },
        { cle: "planning", cible: lienMenu("planning.php") },
        { cle: "score", cible: lienMenu("score.php") },
        { cle: "communaute", cible: lienMenu("forums.php") },
        { cle: "compte", cible: lienMenu("mon-compte.php") },
    ],
    professionnel: [
        { cle: "bienvenue" },
        { cle: "annonces", cible: lienMenu("annonces.php") },
        { cle: "depots", cible: lienMenu("depots.php") },
        { cle: "projets", cible: lienMenu("projets.php") },
        { cle: "contrats", cible: lienMenu("contrats.php") },
        { cle: "statistiques", cible: lienMenu("statistiques.php") },
        { cle: "compte", cible: lienMenu("mon-compte.php") },
    ],
    salarie: [
        { cle: "bienvenue" },
        { cle: "planning", cible: lienMenu("planning.php") },
        { cle: "formations", cible: lienMenu("formations.php") },
        { cle: "ateliers", cible: lienMenu("ateliers.php") },
        { cle: "evenements", cible: lienMenu("evenements.php") },
        { cle: "conseils", cible: lienMenu("conseils.php") },
        { cle: "depots", cible: lienMenu("depots.php") },
        { cle: "compte", cible: lienMenu("mon-compte.php") },
    ],
    administrateur: [
        { cle: "bienvenue" },
        { cle: "tableau_de_bord", cible: lienMenu("tableau-de-bord.php") },
        { cle: "utilisateurs", cible: lienMenu("utilisateurs.php") },
        { cle: "validations", cible: lienMenu("validations.php") },
        { cle: "prestations", cible: lienMenu("prestations.php") },
        { cle: "conteneurs", cible: lienMenu("conteneurs.php") },
        { cle: "finances", cible: lienMenu("finances.php") },
        { cle: "reporting", cible: lienMenu("reporting.php") },
        { cle: "notifications", cible: lienMenu("notifications.php") },
        { cle: "langues", cible: lienMenu("langues.php") },
    ],
};

const MARGE_DECOUPE = 6;
const ECART_BULLE = 16;

class VisiteGuidee {
    constructor(role, etapes) {
        this.role = role;
        this.etapes = etapes;
        this.index = 0;
        this.repositionner = () => this.placer();
        this.clavier = (evenement) => this.gererClavier(evenement);
    }

    demarrer() {
        this.fond = document.createElement("div");
        this.fond.className = "tuto-fond";
        // Bloque tous les clics sur la page tant que la visite n'est pas terminée
        this.fond.addEventListener("mousedown", (e) => e.preventDefault());

        this.decoupe = document.createElement("div");
        this.decoupe.className = "tuto-decoupe";

        this.bulle = document.createElement("div");
        this.bulle.className = "tuto-bulle";
        this.bulle.setAttribute("role", "dialog");
        this.bulle.setAttribute("aria-modal", "true");
        this.bulle.setAttribute("aria-labelledby", "tuto-titre");
        this.bulle.innerHTML = `
            <p class="tuto-compteur"></p>
            <h2 id="tuto-titre"></h2>
            <p class="tuto-texte" aria-live="polite"></p>
            <div class="tuto-points" aria-hidden="true"></div>
            <div class="tuto-boutons">
                <button type="button" class="button-draft tuto-precedent"></button>
                <button type="button" class="button-submit tuto-suivant"></button>
            </div>`;
        this.bulle.querySelector(".tuto-precedent").addEventListener("click", () => this.aller(this.index - 1));
        this.bulle.querySelector(".tuto-suivant").addEventListener("click", () => this.suivant());
        this.bulle.querySelector(".tuto-points").append(...this.etapes.map(() => document.createElement("span")));

        document.body.append(this.fond, this.decoupe, this.bulle);
        document.body.classList.add("modale-ouverte", "tutoriel-actif");
        window.addEventListener("resize", this.repositionner);
        document.addEventListener("keydown", this.clavier, true);
        this.aller(0);
    }

    suivant() {
        if (this.index === this.etapes.length - 1) {
            this.terminer();
        } else {
            this.aller(this.index + 1);
        }
    }

    aller(index) {
        if (index < 0 || index >= this.etapes.length) {
            return;
        }
        this.index = index;
        const etape = this.etapes[index];
        const base = `tutoriel.${this.role}.${etape.cle}`;
        this.bulle.querySelector(".tuto-compteur").textContent =
            t("tutoriel.etape_sur", { numero: index + 1, total: this.etapes.length });
        this.bulle.querySelector("h2").textContent = t(`${base}.titre`);
        this.bulle.querySelector(".tuto-texte").textContent = t(`${base}.texte`);
        this.bulle.querySelectorAll(".tuto-points span").forEach((point, i) => point.classList.toggle("actif", i === index));

        const precedent = this.bulle.querySelector(".tuto-precedent");
        precedent.textContent = t("commun.precedent");
        precedent.hidden = index === 0;
        this.bulle.querySelector(".tuto-suivant").textContent =
            index === this.etapes.length - 1 ? t("commun.terminer") : t("commun.suivant");

        this.cible = etape.cible ? document.querySelector(etape.cible) : null;
        if (this.cible && this.cible.offsetParent === null) {
            this.cible = null;
        }
        this.cible?.scrollIntoView({ block: "nearest", inline: "center", behavior: "instant" });

        this.bulle.classList.remove("tuto-apparition");
        void this.bulle.offsetWidth;
        this.bulle.classList.add("tuto-apparition");
        this.placer();
        this.bulle.querySelector(".tuto-suivant").focus({ preventScroll: true });
    }

    placer() {
        const largeurEcran = window.innerWidth;
        const hauteurEcran = window.innerHeight;
        const bulle = this.bulle.getBoundingClientRect();

        if (!this.cible) {
            // Étape sans élément : écran entièrement assombri, bulle au centre
            Object.assign(this.decoupe.style, {
                top: `${hauteurEcran / 2}px`, left: `${largeurEcran / 2}px`, width: "0px", height: "0px",
            });
            this.decoupe.classList.add("sans-cible");
            Object.assign(this.bulle.style, {
                top: `${Math.max(16, (hauteurEcran - bulle.height) / 2)}px`,
                left: `${Math.max(16, (largeurEcran - bulle.width) / 2)}px`,
            });
            this.bulle.dataset.cote = "centre";
            return;
        }

        this.decoupe.classList.remove("sans-cible");
        const zone = this.cible.getBoundingClientRect();
        Object.assign(this.decoupe.style, {
            top: `${zone.top - MARGE_DECOUPE}px`,
            left: `${zone.left - MARGE_DECOUPE}px`,
            width: `${zone.width + 2 * MARGE_DECOUPE}px`,
            height: `${zone.height + 2 * MARGE_DECOUPE}px`,
        });

        // Bulle à droite de l'élément si la place le permet (menu latéral), sinon dessous, sinon dessus
        let top;
        let left;
        let cote;
        if (zone.right + ECART_BULLE + bulle.width <= largeurEcran - 16) {
            cote = "droite";
            left = zone.right + MARGE_DECOUPE + ECART_BULLE;
            top = zone.top + zone.height / 2 - bulle.height / 2;
        } else if (zone.bottom + ECART_BULLE + bulle.height <= hauteurEcran - 16) {
            cote = "dessous";
            top = zone.bottom + MARGE_DECOUPE + ECART_BULLE;
            left = zone.left + zone.width / 2 - bulle.width / 2;
        } else {
            cote = "dessus";
            top = zone.top - MARGE_DECOUPE - ECART_BULLE - bulle.height;
            left = zone.left + zone.width / 2 - bulle.width / 2;
        }
        top = Math.min(Math.max(16, top), hauteurEcran - bulle.height - 16);
        left = Math.min(Math.max(16, left), largeurEcran - bulle.width - 16);
        Object.assign(this.bulle.style, { top: `${top}px`, left: `${left}px` });
        this.bulle.dataset.cote = cote;
    }

    gererClavier(evenement) {
        // Échap ne ferme pas : la visite est bloquante, on en sort avec « Terminer »
        if (evenement.key === "Escape") {
            evenement.preventDefault();
            evenement.stopPropagation();
        } else if (evenement.key === "ArrowRight") {
            evenement.preventDefault();
            this.suivant();
        } else if (evenement.key === "ArrowLeft") {
            evenement.preventDefault();
            this.aller(this.index - 1);
        } else if (evenement.key === "Tab") {
            // Le focus clavier reste dans la bulle
            const boutons = [...this.bulle.querySelectorAll("button:not([hidden])")];
            const position = boutons.indexOf(document.activeElement);
            const prochain = evenement.shiftKey ? position - 1 : position + 1;
            boutons[(prochain + boutons.length) % boutons.length].focus();
            evenement.preventDefault();
        }
    }

    async terminer() {
        window.removeEventListener("resize", this.repositionner);
        document.removeEventListener("keydown", this.clavier, true);
        document.body.classList.remove("modale-ouverte", "tutoriel-actif");
        this.fond.classList.add("tuto-sortie");
        this.decoupe.classList.add("tuto-sortie");
        this.bulle.classList.add("tuto-sortie");
        setTimeout(() => [this.fond, this.decoupe, this.bulle].forEach((element) => element.remove()), 250);

        const url = new URL(window.location.href);
        if (url.searchParams.has("tutoriel")) {
            url.searchParams.delete("tutoriel");
            history.replaceState(null, "", url);
        }
        try {
            await appelerApi("/moi/tutoriel", { method: "PUT" });
        } catch (erreur) {
        }
    }
}

function lienTutoriel() {
    return `${ESPACES[roleConnecte()]}?tutoriel=1`;
}

// Appelée par demarrerApp() une fois la page affichée.
// Ne démarre que sur la page d'arrivée de l'espace : à la première connexion, ou sur demande (?tutoriel=1).
async function lancerTutorielSiBesoin() {
    const role = roleConnecte();
    const etapes = ETAPES_TUTORIEL[role];
    if (!etapes || window.location.pathname !== ESPACES[role] || doitChangerMotDePasse()) {
        return;
    }
    if (!new URLSearchParams(window.location.search).has("tutoriel")) {
        try {
            const reponse = await appelerApi("/moi");
            if (!reponse.ok || (await reponse.json()).tutoriel_vu) {
                return;
            }
        } catch (erreur) {
            return;
        }
    }
    // Laisse la page finir de s'afficher avant de mesurer la position des éléments
    requestAnimationFrame(() => requestAnimationFrame(() => new VisiteGuidee(role, etapes).demarrer()));
}
