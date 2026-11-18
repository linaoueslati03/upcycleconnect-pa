// Cloche de notifications de l'espace (particuliers, professionnels, salariés).
// Les notifications sont enregistrées par l'API à chaque notifierUtilisateur() : elles s'affichent ici
// même sans HTTPS ni OneSignal. Le compteur est rafraîchi toutes les minutes.

const DELAI_RAFRAICHISSEMENT_NOTIFICATIONS = 60000;

const clocheNotifications = {
    bouton: null,
    compteur: null,
    panneau: null,
    notifications: [],
    nonLues: 0,

    async charger() {
        try {
            const reponse = await appelerApi("/moi/notifications");
            if (!reponse.ok) {
                return;
            }
            const donnees = await reponse.json();
            this.notifications = donnees.notifications;
            this.nonLues = donnees.non_lues;
            this.afficherCompteur();
            if (!this.panneau.hidden) {
                this.remplirPanneau();
            }
        } catch (erreur) {
        }
    },

    afficherCompteur() {
        this.compteur.hidden = this.nonLues === 0;
        this.compteur.textContent = this.nonLues > 99 ? "99+" : String(this.nonLues);
        const libelle = t("notifications_site.titre") + (this.nonLues ? ` (${t("notifications_site.non_lues", { nombre: this.nonLues })})` : "");
        this.bouton.setAttribute("aria-label", libelle);
        this.bouton.title = libelle;
    },

    remplirPanneau() {
        const liste = this.panneau.querySelector("ul");
        liste.replaceChildren();
        this.panneau.querySelector(".tout-lu").hidden = this.nonLues === 0;
        if (this.notifications.length === 0) {
            const vide = document.createElement("li");
            vide.className = "notification-vide";
            vide.textContent = t("notifications_site.aucune");
            liste.append(vide);
            return;
        }
        for (const n of this.notifications) {
            const element = document.createElement("li");
            element.className = n.lue ? "" : "non-lue";
            const lien = document.createElement(n.lien ? "a" : "div");
            lien.className = "notification";
            if (n.lien) {
                lien.href = n.lien;
            }
            const titre = document.createElement("strong");
            titre.textContent = n.titre;
            const message = document.createElement("span");
            message.textContent = n.message;
            const date = document.createElement("time");
            date.dateTime = n.date_creation;
            date.textContent = new Date(n.date_creation).toLocaleString(langueCourante(), {
                day: "numeric", month: "short", hour: "2-digit", minute: "2-digit",
            });
            lien.append(titre, message, date);
            lien.addEventListener("click", () => this.marquerLue(n));
            element.append(lien);
            liste.append(element);
        }
    },

    async marquerLue(notification) {
        if (notification.lue) {
            return;
        }
        notification.lue = true;
        this.nonLues = Math.max(0, this.nonLues - 1);
        this.afficherCompteur();
        try {
            await appelerApi(`/moi/notifications/${notification.id}/lue`, { method: "PUT" });
        } catch (erreur) {
        }
    },

    async toutMarquerLu() {
        this.notifications.forEach((n) => { n.lue = true; });
        this.nonLues = 0;
        this.afficherCompteur();
        this.remplirPanneau();
        try {
            await appelerApi("/moi/notifications/lues", { method: "PUT" });
        } catch (erreur) {
        }
    },

    basculer(ouvrir = this.panneau.hidden) {
        this.panneau.hidden = !ouvrir;
        this.bouton.setAttribute("aria-expanded", String(ouvrir));
        if (ouvrir) {
            this.remplirPanneau();
            this.charger();
        }
    },

    demarrer() {
        this.bouton = document.getElementById("cloche-notifications");
        if (!this.bouton || !["particulier", "professionnel", "salarie"].includes(roleConnecte())) {
            return;
        }
        this.compteur = this.bouton.querySelector(".compteur-notifications");
        this.panneau = document.getElementById("panneau-notifications");
        this.panneau.querySelector("h2").textContent = t("notifications_site.titre");
        const toutLu = this.panneau.querySelector(".tout-lu");
        toutLu.textContent = t("notifications_site.tout_lu");
        toutLu.addEventListener("click", () => this.toutMarquerLu());
        this.bouton.hidden = false;
        this.bouton.addEventListener("click", (e) => {
            e.stopPropagation();
            this.basculer();
        });
        // Fermeture au clic à côté ou avec Échap
        document.addEventListener("click", (e) => {
            if (!this.panneau.hidden && !this.panneau.contains(e.target)) {
                this.basculer(false);
            }
        });
        document.addEventListener("keydown", (e) => {
            if (e.key === "Escape" && !this.panneau.hidden) {
                this.basculer(false);
                this.bouton.focus();
            }
        });
        this.charger();
        setInterval(() => this.charger(), DELAI_RAFRAICHISSEMENT_NOTIFICATIONS);
    },
};
