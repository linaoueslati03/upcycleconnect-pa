TYPES_DEMANDE = ["formation", "atelier", "evenement", "depot", "achat_annonce", "premium"]

SITES = {
    "Paris 10e": 1.0,
    "Paris 11e": 1.6,
    "Paris 13e": 1.4,
    "Paris 16e": 0.7,
    "Bourg-la-Reine": 0.5,
    "Ivry": 0.8,
    "Montreuil": 1.5,
}

PROFILS = {
    "apprenant": {
        "role": "particulier",
        "part": 0.25,
        "demandes_par_mois": 0.9,
        "preferences": {"formation": 6, "atelier": 3, "evenement": 1, "depot": 1, "achat_annonce": 0.3, "premium": 0},
        "annonces_par_mois": 0.1,
        "part_vente": 0.3,
        "categories": {"textile": 3, "bois": 2, "verre": 1, "metal": 1, "plastique": 1},
    },
    "donateur": {
        "role": "particulier",
        "part": 0.25,
        "demandes_par_mois": 0.8,
        "preferences": {"formation": 0.5, "atelier": 1, "evenement": 1, "depot": 7, "achat_annonce": 0.3, "premium": 0},
        "annonces_par_mois": 0.5,
        "part_vente": 0.2,
        "categories": {"plastique": 3, "textile": 3, "verre": 2, "bois": 1, "metal": 1},
    },
    "sociable": {
        "role": "particulier",
        "part": 0.20,
        "demandes_par_mois": 1.0,
        "preferences": {"formation": 1, "atelier": 3, "evenement": 6, "depot": 1, "achat_annonce": 0.2, "premium": 0},
        "annonces_par_mois": 0.15,
        "part_vente": 0.3,
        "categories": {"textile": 2, "verre": 2, "bois": 1, "plastique": 1, "metal": 1},
    },
    "artisan_acheteur": {
        "role": "professionnel",
        "part": 0.18,
        "demandes_par_mois": 1.3,
        "preferences": {"formation": 0.5, "atelier": 0.3, "evenement": 0.5, "depot": 0.3, "achat_annonce": 7, "premium": 1.5},
        "annonces_par_mois": 0.2,
        "part_vente": 0.8,
        "categories": {"bois": 4, "metal": 3, "verre": 1, "textile": 1, "plastique": 0.5},
    },
    "pro_curieux": {
        "role": "professionnel",
        "part": 0.12,
        "demandes_par_mois": 0.7,
        "preferences": {"formation": 4, "atelier": 1, "evenement": 3, "depot": 0.3, "achat_annonce": 1.5, "premium": 0.8},
        "annonces_par_mois": 0.1,
        "part_vente": 0.6,
        "categories": {"bois": 2, "textile": 2, "metal": 1, "verre": 1, "plastique": 1},
    },
}

SUITES = {
    "formation": {"formation": 1.8, "atelier": 1.5},
    "atelier": {"atelier": 1.6, "evenement": 1.4},
    "evenement": {"evenement": 1.5, "atelier": 1.3},
    "depot": {"depot": 2.0},
    "achat_annonce": {"achat_annonce": 1.8, "premium": 1.5},
    "premium": {"achat_annonce": 1.6},
}

PART_HORS_PROFIL = 0.15
PART_SITE_HABITUEL = 0.75


def tirer_poids(rng, poids):
    cles = [c for c, p in poids.items() if p > 0]
    return rng.choices(cles, weights=[poids[c] for c in cles])[0]


def tirer_demande(rng, profil, derniere, types_permis):
    if rng.random() < PART_HORS_PROFIL:
        candidats = [t for t in types_permis if PROFILS[profil]["preferences"][t] > 0]
        return rng.choice(candidats)
    poids = {}
    bonus = SUITES.get(derniere, {})
    for t in types_permis:
        poids[t] = PROFILS[profil]["preferences"][t] * bonus.get(t, 1.0)
    return tirer_poids(rng, poids)
