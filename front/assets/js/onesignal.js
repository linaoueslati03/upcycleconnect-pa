// Notifications push OneSignal (SDK web v16).
// Au chargement : initialisation, puis le navigateur est rattaché au compte connecté (External ID = id utilisateur).
// « notificationsPush » donne l'état réel des notifications à la page Mon compte et permet de les activer
// ou de les désactiver.

window.OneSignalDeferred = window.OneSignalDeferred || [];

let resoudreOneSignal;
// Promesse résolue avec l'objet OneSignal une fois prêt, ou avec null s'il n'est pas utilisable
window.oneSignalPret = new Promise((resoudre) => {
    resoudreOneSignal = resoudre;
});

if (!window.ONESIGNAL_APP_ID) {
    console.warn("OneSignal : ONESIGNAL_APP_ID n'est pas configuré.");
    resoudreOneSignal(null);
} else {
    window.OneSignalDeferred.push(async function (OneSignal) {
        try {
            await OneSignal.init({ appId: window.ONESIGNAL_APP_ID });
            const userId = localStorage.getItem("user_id");
            if (userId) {
                await OneSignal.login(String(userId));
            }
            resoudreOneSignal(OneSignal);
        } catch (erreur) {
            console.error("OneSignal : initialisation impossible :", erreur);
            window.oneSignalErreur = String(erreur && erreur.message ? erreur.message : erreur);
            resoudreOneSignal(null);
        }
    });
}

const notificationsPush = {
    // Attend OneSignal, au plus « delai » ms (le script peut être bloqué ou hors ligne)
    async sdk(delai = 8000) {
        return Promise.race([window.oneSignalPret, new Promise((resoudre) => setTimeout(() => resoudre(null), delai))]);
    },

    // « activees », « desactivees », « bloquees » (refus du navigateur), « erreur » (OneSignal n'a pas démarré)
    // ou « indisponibles » (navigateur non compatible, OneSignal non configuré)
    async etat() {
        const OneSignal = await this.sdk();
        if (!OneSignal && window.oneSignalErreur) {
            return "erreur";
        }
        if (!OneSignal || !OneSignal.Notifications.isPushSupported()) {
            return "indisponibles";
        }
        if (OneSignal.Notifications.permissionNative === "denied") {
            return "bloquees";
        }
        return OneSignal.Notifications.permission && OneSignal.User.PushSubscription.optedIn ? "activees" : "desactivees";
    },

    async activer() {
        const OneSignal = await this.sdk();
        if (!OneSignal) {
            return "indisponibles";
        }
        if (!OneSignal.Notifications.permission) {
            await OneSignal.Notifications.requestPermission();
        }
        if (OneSignal.Notifications.permission && !OneSignal.User.PushSubscription.optedIn) {
            await OneSignal.User.PushSubscription.optIn();
        }
        return this.etat();
    },

    async desactiver() {
        const OneSignal = await this.sdk();
        if (OneSignal) {
            await OneSignal.User.PushSubscription.optOut();
        }
        return this.etat();
    },

    // Prévient la page quand l'état change (autorisation donnée ou retirée ailleurs, abonnement mis à jour)
    async surChangement(rappel) {
        const OneSignal = await this.sdk();
        if (!OneSignal) {
            return;
        }
        OneSignal.Notifications.addEventListener("permissionChange", () => this.etat().then(rappel));
        OneSignal.User.PushSubscription.addEventListener("change", () => this.etat().then(rappel));
    },
};

function deconnecterOneSignal() {
    if (!window.OneSignalDeferred) return;

    window.OneSignalDeferred.push(async function (OneSignal) {
        try {
            await OneSignal.logout();
        } catch (erreur) {
            console.error("Erreur de déconnexion OneSignal :", erreur);
        }
    });
}
