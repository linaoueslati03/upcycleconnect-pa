import os
from pathlib import Path

RACINE = Path(__file__).resolve().parent.parent


def lire_env():
    valeurs = {}
    fichier = RACINE / ".env"
    if fichier.exists():
        for ligne in fichier.read_text(encoding="utf-8").splitlines():
            ligne = ligne.strip()
            if not ligne or ligne.startswith("#") or "=" not in ligne:
                continue
            cle, valeur = ligne.split("=", 1)
            valeurs[cle.strip()] = valeur.strip().strip('"').strip("'")
    return valeurs


def parametres_connexion():
    env = lire_env()

    def variable(nom, defaut):
        return os.environ.get(nom) or env.get(nom) or defaut

    return {
        "host": variable("DB_HOST", "localhost"),
        "port": variable("DB_PORT", "5433"),
        "user": variable("DB_USER", "postgres"),
        "password": variable("DB_PASSWORD", "postgres"),
        "dbname": variable("DB_NAME", "upcycleconnect"),
    }
