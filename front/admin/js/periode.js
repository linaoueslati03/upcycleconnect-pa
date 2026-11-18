function dateLocale(date) {
    const mois = String(date.getMonth() + 1).padStart(2, "0");
    const jour = String(date.getDate()).padStart(2, "0");
    return `${date.getFullYear()}-${mois}-${jour}`;
}

function periodeRaccourci(raccourci) {
    const aujourdhui = new Date();
    const debuts = {
        "30_jours": new Date(aujourdhui.getFullYear(), aujourdhui.getMonth(), aujourdhui.getDate() - 29),
        "12_mois": new Date(aujourdhui.getFullYear(), aujourdhui.getMonth() - 11, 1),
        mois_courant: new Date(aujourdhui.getFullYear(), aujourdhui.getMonth(), 1),
        annee_courante: new Date(aujourdhui.getFullYear(), 0, 1),
    };
    return { debut: dateLocale(debuts[raccourci]), fin: dateLocale(aujourdhui) };
}
