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

Le script génère aussi les charges de la plateforme (loyer, serveurs, licences, salaires, rémunération des animateurs selon le nombre de sessions, marketing, achat de conteneurs) sur toute la période, pour que les pages Finances et Tableau de bord aient un résultat réaliste. Ces charges n'ont pas d'auteur (`auteur_id` vide), contrairement à celles saisies par un administrateur. `--charges-seulement` les régénère seules, sans toucher au reste des données.

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

## Entraîner le modèle

```bash
.venv/bin/python entrainer.py        # écrit modeles/metriques.json et modeles/figures/*.png
```

Options : `--dataset`, `--coupures-test 2` (dernières coupures gardées pour le test), `--sortie`.

### Ce que le modèle prédit (la cible)

Pour un acteur, le type de la prochaine prestation qu'il demandera, parmi 6 : formation, atelier, événement,
dépôt en conteneur, achat d'annonce, abonnement premium. Ce sont toutes les façons de demander une prestation dans l'appli.
Dans le dataset, c'est la première demande dans les 60 jours après la coupure ; les acteurs qui ne demandent rien
pendant ces 60 jours sont écartés (on prédit ce que l'acteur demandera, pas s'il reviendra).

### Méthode

- Variables : les 16 colonnes du dataset. Les 4 variables texte (`role`, `derniere_demande`, `categorie_preferee`,
  `site_principal`) sont encodées en colonnes 0/1 (une par valeur possible) et 5 variables sont ajoutées :
  la part de chaque type dans l'historique de l'acteur (`part_formations` = `nb_formations` / total des 5 compteurs,
  0 si aucune demande). Soit 41 colonnes. L'API Go devra reproduire exactement cet encodage.
- Découpage par date : entraînement sur les coupures les plus anciennes, test sur les 2 dernières.
  Le modèle est jugé sur le futur, comme dans l'appli.
- Trois modèles comparés : arbre de décision, random forest, régression logistique. Leurs réglages sont choisis par
  validation croisée (5 plis) sur l'entraînement uniquement, en maximisant le F1 macro. Graine fixe : résultats reproductibles.
- Deux références sans apprentissage : toujours la classe la plus fréquente, et « la même chose que la dernière demande ».
- Pas de pondération des classes : pondérer pour faire ressortir premium (3 % des cas) faisait perdre 6 points
  d'exactitude pour des prédictions premium fausses 93 % du temps. Premium reste visible par sa probabilité.

### Résultats (test : 756 lignes, août et septembre 2026)

| Modèle | Exactitude | F1 macro |
|---|---|---|
| Arbre de décision (déployé) | 46,4 % | 37,0 % |
| Random forest | 47,3 % | 40,7 % |
| Régression logistique | 47,6 % | 39,3 % |
| Référence : classe la plus fréquente | 20,6 % | 5,7 % |
| Référence : même chose que la dernière demande | 41,0 % | 33,9 % |

L'arbre retenu (profondeur 6, 27 feuilles) fait presque aussi bien que les deux autres modèles tout en restant lisible
et léger à exécuter dans l'API Go : c'est pour cela qu'il est déployé. Les données étant générées avec une part de hasard
(chaque demande est tirée selon le profil, avec 15 % de demandes hors profil), aucun modèle ne peut dépasser
environ 50 %. Les variables les plus utiles sont le rôle, la part de dépôts et la part d'achats d'annonces.

Figures pour le rapport : `modeles/figures/` (comparaison des modèles, matrice de confusion, importance des variables,
aperçu de l'arbre).

## Export de l'arbre pour l'API Go

`entrainer.py` exporte aussi l'arbre retenu dans `api/internal/ml/` :

- `arbre.json` : la version (date d'entraînement), les scores de test, les 6 classes, la recette d'encodage
  (variables numériques, compteurs des parts, valeurs possibles des variables texte), la liste ordonnée des 41 colonnes
  et les nœuds. Un nœud question se lit « si colonne[`variable`] ≤ `seuil` aller à `gauche`, sinon à `droite` » ;
  une feuille donne les `probabilites` des 6 classes (même ordre que `classes`) et son `effectif` d'entraînement.
  On part du nœud 0 ; la classe prédite est la plus probable (la première en cas d'égalité).
- `cas_test.json` : 60 acteurs du dataset avec leurs 16 variables brutes, la classe et les probabilités attendues,
  pour vérifier en Go que l'encodage et le parcours donnent exactement le même résultat.

Avant d'écrire ces fichiers, le script refait le parcours du JSON sur tout le dataset et s'arrête si une seule prédiction
diffère de scikit-learn. Pour mettre un nouveau modèle en ligne : relancer `entrainer.py`, puis reconstruire l'API
(`docker compose up -d --build api`).
