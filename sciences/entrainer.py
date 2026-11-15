import argparse
import json
from datetime import datetime
from pathlib import Path

import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt
import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import (
    ConfusionMatrixDisplay,
    accuracy_score,
    classification_report,
    confusion_matrix,
    f1_score,
)
from sklearn.model_selection import GridSearchCV, StratifiedKFold
from sklearn.pipeline import make_pipeline
from sklearn.preprocessing import StandardScaler
from sklearn.tree import DecisionTreeClassifier, plot_tree

from extraire_dataset import VARIABLES
from profils import SITES, TYPES_DEMANDE

DOSSIER = Path(__file__).resolve().parent
DOSSIER_ML = DOSSIER.parent / "api" / "internal" / "ml"
GRAINE = 42
NOMBRE_CAS_TEST = 60

CATEGORIELLES = ["role", "derniere_demande", "categorie_preferee", "site_principal"]
NUMERIQUES = [v for v in VARIABLES if v not in CATEGORIELLES]
COMPTEURS = ["nb_formations", "nb_ateliers", "nb_evenements", "nb_depots", "nb_achats_annonces"]

NOMS_MODELES = {
    "arbre": "Arbre de décision",
    "foret": "Random forest",
    "logistique": "Régression logistique",
    "reference_majorite": "Référence : classe la plus fréquente",
    "reference_derniere": "Référence : même chose que la dernière demande",
}

NOMS_COURTS = {
    "arbre": "Arbre de décision",
    "foret": "Random forest",
    "logistique": "Régression\nlogistique",
    "reference_majorite": "Référence :\nclasse la plus\nfréquente",
    "reference_derniere": "Référence :\nmême chose que la\ndernière demande",
}


def lire_arguments():
    parser = argparse.ArgumentParser(description="Entraîne et compare les modèles de prédiction de la prochaine demande.")
    parser.add_argument("--dataset", type=Path, default=DOSSIER / "donnees" / "dataset.csv",
                        help="dataset produit par extraire_dataset.py (défaut : donnees/dataset.csv)")
    parser.add_argument("--coupures-test", type=int, default=2,
                        help="nombre de dernières coupures gardées pour le test (défaut : 2)")
    parser.add_argument("--sortie", type=Path, default=DOSSIER / "modeles",
                        help="dossier des métriques et figures (défaut : modeles)")
    parser.add_argument("--export", type=Path, default=DOSSIER_ML,
                        help="dossier où écrire arbre.json et cas_test.json pour l'API Go (défaut : api/internal/ml)")
    return parser.parse_args()


def valeurs_possibles(dataset):
    categories = sorted(set(dataset["categorie_preferee"]) | {"aucune"})
    return {
        "role": ["particulier", "professionnel"],
        "derniere_demande": TYPES_DEMANDE + ["aucune"],
        "categorie_preferee": categories,
        "site_principal": list(SITES) + ["autre", "inconnu"],
    }


def encoder(lignes, valeurs):
    colonnes = {nom: lignes[nom].astype(float) for nom in NUMERIQUES}
    total = lignes[COMPTEURS].sum(axis=1)
    for compteur in COMPTEURS:
        colonnes[f"part_{compteur[3:]}"] = (lignes[compteur] / total.where(total > 0, 1)).where(total > 0, 0.0)
    for variable, possibles in valeurs.items():
        for valeur in possibles:
            colonnes[f"{variable}={valeur}"] = (lignes[variable] == valeur).astype(float)
    return pd.DataFrame(colonnes, index=lignes.index)


def references(entrainement, test):
    majorite = entrainement["cible"].mode()[0]
    predictions_majorite = pd.Series(majorite, index=test.index)
    predictions_derniere = test["derniere_demande"].where(test["derniere_demande"].isin(TYPES_DEMANDE), majorite)
    return {"reference_majorite": predictions_majorite, "reference_derniere": predictions_derniere}


def entrainer_modeles(x, y):
    validation = StratifiedKFold(n_splits=5, shuffle=True, random_state=GRAINE)
    recherches = {
        "arbre": GridSearchCV(
            DecisionTreeClassifier(random_state=GRAINE),
            {"max_depth": [3, 4, 5, 6, 7, 8, 10], "min_samples_leaf": [5, 10, 20, 40]},
            scoring="f1_macro", cv=validation,
        ),
        "foret": GridSearchCV(
            RandomForestClassifier(n_estimators=300, random_state=GRAINE, n_jobs=-1),
            {"max_depth": [6, 10, None], "min_samples_leaf": [1, 5, 10]},
            scoring="f1_macro", cv=validation,
        ),
        "logistique": GridSearchCV(
            make_pipeline(StandardScaler(), LogisticRegression(max_iter=3000)),
            {"logisticregression__C": [0.1, 1.0, 10.0]},
            scoring="f1_macro", cv=validation,
        ),
    }
    modeles = {}
    for nom, recherche in recherches.items():
        recherche.fit(x, y)
        modeles[nom] = recherche
        print(f"  {NOMS_MODELES[nom]} : réglage retenu {recherche.best_params_}")
    return modeles


def scores(y_vrai, y_predit):
    return {
        "exactitude": round(accuracy_score(y_vrai, y_predit), 4),
        "f1_macro": round(f1_score(y_vrai, y_predit, average="macro", labels=TYPES_DEMANDE, zero_division=0), 4),
        "par_classe": {
            classe: {k: round(v, 4) for k, v in valeurs.items()}
            for classe, valeurs in classification_report(
                y_vrai, y_predit, labels=TYPES_DEMANDE, output_dict=True, zero_division=0
            ).items()
            if classe in TYPES_DEMANDE
        },
    }


def importances(modele, colonnes, nombre=15):
    serie = pd.Series(modele.feature_importances_, index=colonnes).sort_values(ascending=False)
    return {nom: round(valeur, 4) for nom, valeur in serie.head(nombre).items() if valeur > 0}


def figure_comparaison(resultats, chemin):
    noms = list(resultats)
    exactitudes = [resultats[n]["exactitude"] * 100 for n in noms]
    f1 = [resultats[n]["f1_macro"] * 100 for n in noms]
    positions = range(len(noms))
    fig, axe = plt.subplots(figsize=(10, 5))
    axe.bar([p - 0.2 for p in positions], exactitudes, width=0.4, label="Exactitude", color="#2E7D32")
    axe.bar([p + 0.2 for p in positions], f1, width=0.4, label="F1 macro", color="#8E8BB8")
    for p, (e, f) in enumerate(zip(exactitudes, f1)):
        axe.text(p - 0.2, e + 1, f"{e:.0f} %", ha="center", fontsize=8)
        axe.text(p + 0.2, f + 1, f"{f:.0f} %", ha="center", fontsize=8)
    axe.set_xticks(list(positions), [NOMS_COURTS[n] for n in noms], fontsize=9)
    axe.set_ylabel("Score sur le jeu de test (%)")
    axe.set_ylim(0, 100)
    axe.set_title("Comparaison des modèles et des références")
    axe.legend()
    fig.tight_layout()
    fig.savefig(chemin, dpi=150)
    plt.close(fig)


def figure_confusion(y_vrai, y_predit, chemin):
    fig, axe = plt.subplots(figsize=(7, 6))
    ConfusionMatrixDisplay(confusion_matrix(y_vrai, y_predit, labels=TYPES_DEMANDE), display_labels=TYPES_DEMANDE).plot(
        ax=axe, cmap="Greens", colorbar=False, xticks_rotation=35)
    axe.set_xlabel("Classe prédite")
    axe.set_ylabel("Vraie classe")
    axe.set_title("Arbre de décision : matrice de confusion (test)")
    fig.tight_layout()
    fig.savefig(chemin, dpi=150)
    plt.close(fig)


def figure_importances(valeurs, chemin):
    serie = pd.Series(valeurs).sort_values()
    fig, axe = plt.subplots(figsize=(8, 6))
    axe.barh(serie.index, serie.values * 100, color="#2E7D32")
    axe.set_xlabel("Importance (%)")
    axe.set_title("Arbre de décision : variables les plus importantes")
    fig.tight_layout()
    fig.savefig(chemin, dpi=150)
    plt.close(fig)


def figure_arbre(arbre, colonnes, chemin):
    fig, axe = plt.subplots(figsize=(22, 10))
    plot_tree(arbre, max_depth=3, feature_names=list(colonnes), class_names=list(arbre.classes_),
              filled=True, impurity=False, proportion=True, rounded=True, fontsize=7, ax=axe)
    axe.set_title("Arbre de décision : les 3 premiers niveaux")
    fig.tight_layout()
    fig.savefig(chemin, dpi=120)
    plt.close(fig)


def encoder_comme_go(modele, variables):
    encodage = modele["encodage"]
    vecteur = [float(variables[nom]) for nom in encodage["numeriques"]]
    total = sum(float(variables[compteur]) for compteur in encodage["parts"])
    vecteur += [float(variables[compteur]) / total if total > 0 else 0.0 for compteur in encodage["parts"]]
    for variable, possibles in encodage["categorielles"].items():
        vecteur += [1.0 if variables[variable] == valeur else 0.0 for valeur in possibles]
    return vecteur


def parcourir(modele, vecteur):
    noeud = modele["noeuds"][0]
    while "probabilites" not in noeud:
        suivant = noeud["gauche"] if vecteur[noeud["variable"]] <= noeud["seuil"] else noeud["droite"]
        noeud = modele["noeuds"][suivant]
    return noeud["probabilites"]


def classe_predite(modele, probabilites):
    return modele["classes"][max(range(len(probabilites)), key=lambda i: (probabilites[i], -i))]


def construire_modele(arbre, colonnes, valeurs, version, resultat_test):
    t = arbre.tree_
    noeuds = []
    for i in range(t.node_count):
        if t.children_left[i] == -1:
            valeurs_feuille = t.value[i][0]
            noeuds.append({
                "probabilites": [float(v) for v in valeurs_feuille / valeurs_feuille.sum()],
                "effectif": int(t.n_node_samples[i]),
            })
        else:
            noeuds.append({
                "variable": int(t.feature[i]),
                "seuil": float(t.threshold[i]),
                "gauche": int(t.children_left[i]),
                "droite": int(t.children_right[i]),
            })
    return {
        "version": version,
        "exactitude_test": resultat_test["exactitude"],
        "f1_macro_test": resultat_test["f1_macro"],
        "profondeur": int(arbre.get_depth()),
        "classes": [str(c) for c in arbre.classes_],
        "encodage": {"numeriques": NUMERIQUES, "parts": COMPTEURS, "categorielles": valeurs},
        "colonnes": list(colonnes),
        "noeuds": noeuds,
    }


def valeur_brute(valeur):
    return valeur.item() if hasattr(valeur, "item") else valeur


def exporter_arbre(arbre, colonnes, valeurs, version, resultat_test, dataset, dossier):
    modele = construire_modele(arbre, colonnes, valeurs, version, resultat_test)
    attendues = arbre.predict_proba(encoder(dataset, valeurs))
    ecarts = 0
    for (_, ligne), attendu in zip(dataset.iterrows(), attendues):
        vecteur = encoder_comme_go(modele, ligne)
        if len(vecteur) != len(modele["colonnes"]):
            raise SystemExit("Export refusé : l'encodage ne produit pas le bon nombre de colonnes.")
        obtenu = parcourir(modele, vecteur)
        if any(abs(a - b) > 1e-12 for a, b in zip(obtenu, attendu)) or \
                classe_predite(modele, obtenu) != arbre.classes_[attendu.argmax()]:
            ecarts += 1
    if ecarts:
        raise SystemExit(f"Export refusé : le parcours du JSON diffère de scikit-learn sur {ecarts} lignes.")

    cas = []
    for _, ligne in dataset.sample(NOMBRE_CAS_TEST, random_state=GRAINE).iterrows():
        probabilites = parcourir(modele, encoder_comme_go(modele, ligne))
        cas.append({
            "variables": {nom: valeur_brute(ligne[nom]) for nom in VARIABLES},
            "classe": classe_predite(modele, probabilites),
            "probabilites": probabilites,
        })

    dossier.mkdir(parents=True, exist_ok=True)
    with open(dossier / "arbre.json", "w", encoding="utf-8") as fichier:
        json.dump(modele, fichier, ensure_ascii=False, indent=1)
    with open(dossier / "cas_test.json", "w", encoding="utf-8") as fichier:
        json.dump(cas, fichier, ensure_ascii=False, indent=1)
    return len(modele["noeuds"])


def main():
    args = lire_arguments()
    dataset = pd.read_csv(args.dataset)
    coupures = sorted(dataset["coupure"].unique())
    coupures_test = coupures[-args.coupures_test:]
    test = dataset[dataset["coupure"].isin(coupures_test)]
    entrainement = dataset[~dataset["coupure"].isin(coupures_test)]

    valeurs = valeurs_possibles(dataset)
    x_entrainement, x_test = encoder(entrainement, valeurs), encoder(test, valeurs)
    y_entrainement, y_test = entrainement["cible"], test["cible"]
    print(f"Dataset : {len(dataset)} lignes, {x_entrainement.shape[1]} colonnes après encodage")
    print(f"Entraînement : {len(entrainement)} lignes ({coupures[0]} → {coupures[-args.coupures_test - 1]})")
    print(f"Test : {len(test)} lignes ({coupures_test[0]} → {coupures_test[-1]})")

    print("Entraînement des modèles (validation croisée sur l'entraînement uniquement)…")
    modeles = entrainer_modeles(x_entrainement, y_entrainement)

    predictions = {nom: pd.Series(m.predict(x_test), index=test.index) for nom, m in modeles.items()}
    predictions.update(references(entrainement, test))
    resultats = {nom: scores(y_test, predictions[nom]) for nom in NOMS_MODELES}

    arbre = modeles["arbre"].best_estimator_
    foret = modeles["foret"].best_estimator_
    premium_particuliers = int(((predictions["arbre"] == "premium") & (test["role"] == "particulier")).sum())

    args.sortie.mkdir(parents=True, exist_ok=True)
    figures = args.sortie / "figures"
    figures.mkdir(exist_ok=True)
    importance_arbre = importances(arbre, x_entrainement.columns)
    figure_comparaison(resultats, figures / "comparaison_modeles.png")
    figure_confusion(y_test, predictions["arbre"], figures / "matrice_confusion_arbre.png")
    figure_importances(importance_arbre, figures / "importance_variables_arbre.png")
    figure_arbre(arbre, x_entrainement.columns, figures / "arbre_apercu.png")

    version = datetime.now().isoformat(timespec="seconds")
    metriques = {
        "date_entrainement": version,
        "dataset": {
            "lignes": len(dataset),
            "entrainement": {"lignes": len(entrainement), "coupures": [coupures[0], coupures[-args.coupures_test - 1]]},
            "test": {"lignes": len(test), "coupures": [coupures_test[0], coupures_test[-1]]},
            "repartition_test": y_test.value_counts().to_dict(),
        },
        "variables": {
            "numeriques": NUMERIQUES,
            "derivees": {f"part_{c[3:]}": f"{c} / somme des {len(COMPTEURS)} compteurs (0 si aucune demande)"
                         for c in COMPTEURS},
            "categorielles": valeurs,
        },
        "modele_retenu": "arbre",
        "resultats": {
            nom: {"nom": NOMS_MODELES[nom], **resultats[nom],
                  **({"reglage": modeles[nom].best_params_} if nom in modeles else {})}
            for nom in NOMS_MODELES
        },
        "arbre": {
            "profondeur": int(arbre.get_depth()),
            "feuilles": int(arbre.get_n_leaves()),
            "importance_variables": importance_arbre,
            "matrice_confusion": {"classes": TYPES_DEMANDE,
                                  "valeurs": confusion_matrix(y_test, predictions["arbre"], labels=TYPES_DEMANDE).tolist()},
            "premium_predit_pour_des_particuliers": premium_particuliers,
        },
        "foret": {"importance_variables": importances(foret, x_entrainement.columns)},
    }
    with open(args.sortie / "metriques.json", "w", encoding="utf-8") as fichier:
        json.dump(metriques, fichier, ensure_ascii=False, indent=2)

    print("\nRésultats sur le jeu de test :")
    print(f"  {'Modèle':<50} {'Exactitude':>10} {'F1 macro':>9}")
    for nom in NOMS_MODELES:
        print(f"  {NOMS_MODELES[nom]:<50} {resultats[nom]['exactitude']:>9.1%} {resultats[nom]['f1_macro']:>9.1%}")
    print(f"\nArbre retenu : profondeur {arbre.get_depth()}, {arbre.get_n_leaves()} feuilles")
    print(f"Premium prédit pour des particuliers : {premium_particuliers}")
    print(f"Métriques et figures écrites dans {args.sortie}")

    nb_noeuds = exporter_arbre(arbre, x_entrainement.columns, valeurs, version, resultats["arbre"], dataset, args.export)
    print(f"Arbre exporté ({nb_noeuds} nœuds), identique à scikit-learn sur les {len(dataset)} lignes du dataset")
    print(f"Fichiers pour l'API Go : {args.export / 'arbre.json'} et {args.export / 'cas_test.json'}")


if __name__ == "__main__":
    main()
