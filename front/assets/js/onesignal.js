window.OneSignalDeferred = window.OneSignalDeferred || [];

window.OneSignalDeferred.push(async function (OneSignal) {
    const appId = window.ONESIGNAL_APP_ID;

    if (!appId) {
        console.warn("OneSignal : ONESIGNAL_APP_ID n'est pas configuré.");
        return;
    }

    await OneSignal.init({ appId });

    const userId = localStorage.getItem("user_id");
    if (userId) {
        await OneSignal.login(String(userId));
    }

    const bouton = document.getElementById("activer-notifications");
    if (bouton) {
        bouton.addEventListener("click", async () => {
            bouton.disabled = true;
            try {
                await OneSignal.Notifications.requestPermission();
            } catch (erreur) {
                console.error("Impossible d'activer les notifications :", erreur);
            } finally {
                bouton.disabled = false;
            }
        });
    }
});

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