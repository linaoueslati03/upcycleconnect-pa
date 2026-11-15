# Bloc Sciences — module 1

Outils hors ligne du module Sciences : génération des données et du dataset, entraînement du modèle.
Rien dans ce dossier n'est déployé : le reporting et la prédiction tournent dans l'API Go.

## Installation

```bash
cd sciences
python3 -m venv .venv
.venv/bin/pip install -r requirements.txt
```

La connexion à la base est lue dans le `.env` à la racine (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`),
avec les mêmes valeurs par défaut que l'API (`localhost:5433`). Une variable d'environnement passe avant le `.env`.

La base doit avoir toutes les migrations, dont `022_reporting.sql`.

## Générer les données

```bash
.venv/bin/python generer_donnees.py                      # 600 acteurs, graine 42, historique jusqu'à aujourd'hui
.venv/bin/python generer_donnees.py --reinitialiser      # remplace les données déjà générées
.venv/bin/python generer_donnees.py --acteurs 15000 --reinitialiser   # base peuplée en masse (étape 5)
```

Options : `--acteurs`, `--graine`, `--date-reference AAAA-MM-JJ` (mêmes options = mêmes données), `--reinitialiser`.

Le script écrit tout dans une seule transaction : en cas d'erreur, rien n'est enregistré.
Les comptes générés ont une adresse en `@example.com`, domaine réservé aux adresses fictives (RFC 2606), avec le même
mot de passe que les comptes de test des seeds. Rien d'autre ne les distingue dans l'appli.
`--reinitialiser` supprime uniquement ces comptes et ce qui leur est rattaché : ne pas créer de vrai compte en `@example.com`.

### Ce qui est généré

- 6 salariés, qui créent et valident les offres ;
- des formations, ateliers et événements sur 18 mois, répartis sur les 7 sites, plus un mois et demi à venir ;
- les acteurs (70 % de particuliers, 30 % de professionnels) et leur historique : inscriptions, dépôts en conteneur,
  annonces, achats d'annonces, abonnements premium, paiements et points d'Upcycling Score.

Les règles de l'appli sont respectées : paiements au format d'`enregistrerPaiement()` (commission seulement sur les ventes,
7 % puis 8 %), un premium crée une ligne `abonnements` et met à jour `utilisateurs.abonnement`, dépôts dans l'ordre
demande → validée → déposée → récupérée, points de score identiques à l'API (+5 inscription, +15 annonce cédée,
+10 dépôt récupéré).

### Les profils (`profils.py`)

Chaque acteur suit l'un des 5 profils, qui fixe ses préférences entre les 6 types de demande
(formation, atelier, événement, dépôt, achat d'annonce, premium) :

| Profil | Rôle | Part | Demandes surtout |
|---|---|---|---|
| apprenant | particulier | 25 % | formations, ateliers |
| donateur | particulier | 25 % | dépôts en conteneur |
| sociable | particulier | 20 % | événements, ateliers |
| artisan_acheteur | professionnel | 18 % | achats d'annonces, premium |
| pro_curieux | professionnel | 12 % | formations, événements |

S'y ajoutent l'effet de la demande précédente (après une formation, une autre formation ou un atelier est plus probable)
et 15 % de demandes hors profil. Le profil n'est pas enregistré en base : le modèle doit le retrouver à partir du comportement.

## Extraire le dataset

```bash
.venv/bin/python extraire_dataset.py          # écrit donnees/dataset.csv
```

Options : `--mois 12` (nombre de dates de coupure), `--horizon 60` (jours observés après chaque coupure),
`--date-reference AAAA-MM-JJ`, `--sortie`.

Chaque ligne décrit un acteur à une date de coupure : ses variables sont calculées par `variables_acteurs(coupure)`
(migration 022), donc uniquement avec ce qui s'est passé avant cette date, et la `cible` est le type de sa première
demande dans les 60 jours suivants. Les acteurs sans demande dans cette fenêtre sont écartés. Les coupures sont
mensuelles, la dernière 60 jours avant la date de référence pour que chaque fenêtre soit complète.

Colonnes : `utilisateur_id` et `coupure` (suivi et découpage entraînement / test, pas des variables du modèle),
les 16 variables, puis `cible`. `categorie_preferee` est exportée sous forme de code (`bois`, `metal`…, `aucune`)
car les identifiants changent d'une base à l'autre.
