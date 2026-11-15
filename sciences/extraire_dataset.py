import argparse
import csv
import sys
from collections import Counter
from datetime import date, datetime, time, timedelta
from pathlib import Path

import psycopg

from config import parametres_connexion
from profils import TYPES_DEMANDE

DOSSIER = Path(__file__).resolve().parent

VARIABLES = [
    "role",
    "premium",
    "anciennete_jours",
    "upcycling_score",
    "nb_annonces",
    "nb_annonces_vente",
    "nb_formations",
    "nb_ateliers",
    "nb_evenements",
    "nb_depots",
    "nb_achats_annonces",
    "montant_depense",
    "jours_depuis_derniere_demande",
    "derniere_demande",
    "categorie_preferee",
    "site_principal",
]
COLONNES = ["utilisateur_id", "coupure"] + VARIABLES + ["cible"]
LIGNES_MINIMUM = 500

REQUETE = """
    WITH coupures AS (
        SELECT %(derniere)s::timestamp - make_interval(months => k) AS coupure
        FROM generate_series(0, %(mois)s - 1) AS k
    )
    SELECT v.utilisateur_id, c.coupure::date,
           v.role, v.premium, v.anciennete_jours, v.upcycling_score,
           v.nb_annonces, v.nb_annonces_vente, v.nb_formations, v.nb_ateliers, v.nb_evenements,
           v.nb_depots, v.nb_achats_annonces, v.montant_depense, v.jours_depuis_derniere_demande,
           v.derniere_demande, COALESCE(cm.code, 'aucune'), v.site_principal,
           suivante.type_demande
    FROM coupures c
    CROSS JOIN LATERAL variables_acteurs(c.coupure) v
    JOIN LATERAL (
        SELECT d.type_demande
        FROM vue_demandes_acteurs d
        WHERE d.utilisateur_id = v.utilisateur_id
          AND d.date_demande >= c.coupure
          AND d.date_demande < c.coupure + make_interval(days => %(horizon)s)
        ORDER BY d.date_demande, d.type_demande
        LIMIT 1
    ) suivante ON TRUE
    LEFT JOIN categories_materiaux cm ON cm.id = v.categorie_preferee
    ORDER BY c.coupure, v.utilisateur_id
"""


def lire_arguments():
    parser = argparse.ArgumentParser(description="Extrait le dataset d'entraînement du modèle de prédiction.")
    parser.add_argument("--mois", type=int, default=12, help="nombre de dates de coupure mensuelles (défaut : 12)")
    parser.add_argument("--horizon", type=int, default=60,
                        help="jours observés après chaque coupure pour trouver la cible (défaut : 60)")
    parser.add_argument("--date-reference", type=date.fromisoformat, default=date.today(),
                        help="fin des données, AAAA-MM-JJ (défaut : aujourd'hui)")
    parser.add_argument("--sortie", type=Path, default=DOSSIER / "donnees" / "dataset.csv",
                        help="fichier CSV produit (défaut : donnees/dataset.csv)")
    return parser.parse_args()


def afficher_resume(lignes):
    print(f"{len(lignes)} lignes")
    par_coupure = Counter(ligne[1] for ligne in lignes)
    print("Par coupure :", ", ".join(f"{c.isoformat()} {n}" for c, n in sorted(par_coupure.items())))
    classes = Counter(ligne[-1] for ligne in lignes)
    print("Classes :", ", ".join(f"{t} {classes[t]} ({100 * classes[t] / len(lignes):.0f} %)" for t in TYPES_DEMANDE))
    if len(lignes) < LIGNES_MINIMUM:
        print(f"Attention : moins de {LIGNES_MINIMUM} lignes, générez plus d'acteurs ou augmentez --mois.")
    absentes = [t for t in TYPES_DEMANDE if classes[t] == 0]
    if absentes:
        print("Attention : classes absentes du dataset :", ", ".join(absentes))


def main():
    args = lire_arguments()
    reference = datetime.combine(args.date_reference, time())
    derniere = reference - timedelta(days=args.horizon)
    with psycopg.connect(**parametres_connexion()) as connexion, connexion.cursor() as cur:
        cur.execute("SELECT to_regclass('vue_variables_acteurs') IS NOT NULL")
        if not cur.fetchone()[0]:
            sys.exit("La base n'a pas la migration 022_reporting.sql.\n"
                     "Appliquez-la depuis la racine du projet : docker compose exec postgres sh /db/migrer.sh")
        cur.execute(REQUETE, {"derniere": derniere, "mois": args.mois, "horizon": args.horizon})
        lignes = cur.fetchall()

    if not lignes:
        sys.exit("Aucune ligne : la base ne contient pas de demandes. Lancez d'abord generer_donnees.py.")

    args.sortie.parent.mkdir(parents=True, exist_ok=True)
    with open(args.sortie, "w", newline="", encoding="utf-8") as fichier:
        ecrivain = csv.writer(fichier)
        ecrivain.writerow(COLONNES)
        ecrivain.writerows(lignes)

    print(f"Dataset écrit dans {args.sortie}")
    afficher_resume(lignes)


if __name__ == "__main__":
    main()
