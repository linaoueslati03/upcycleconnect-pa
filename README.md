# UpcycleConnect

Plateforme web d'upcycling : mise en relation de particuliers, professionnels/artisans et salariés
autour du réemploi d'objets et de matériaux.

## Stack

| Couche | Technologie |
|---|---|
| API (toute la logique métier) | Go (`net/http`, `database/sql`) |
| Base de données | PostgreSQL |
| Front | PHP + JavaScript (Vue 3 chargé par CDN, sans build) |
| Conteneurs | Docker Compose |

## Arborescence

```
.
├── api/                 # API Go (module github.com/linaoueslati03/upcycleconnect-pa/api)
│   ├── cmd/api/         # point d'entrée : main.go
│   └── internal/
│       ├── api/         # handlers HTTP, un fichier par domaine + routes.go
│       ├── database/    # connexion PostgreSQL
│       └── middleware/  # CORS
├── db/
│   ├── migrations/      # schéma SQL, un fichier numéroté par évolution
│   ├── seeds/           # données de démonstration
│   ├── migrer.sh        # applique les migrations pas encore passées
│   └── initialiser.sh   # migrations + données, à la création de la base
├── front/
│   ├── espace/          # application des particuliers, professionnels et salariés
│   │   └── inclus/espace.php   # liste des pages et rôles autorisés pour chacune
│   ├── admin/           # back-office administrateur
│   ├── inclus/          # en-tête et menu communs
│   └── assets/          # css, js, traductions (locales/), images
├── docker-compose.yml
└── .env.example
```

## Lancer le projet

Prérequis : Docker Desktop et PHP 8.

```bash
cp .env.example .env
docker compose up -d --build
php -S localhost:8000 -t front front/routeur.php
```

| Service | URL |
|---|---|
| Front | http://localhost:8000 (page de connexion) |
| API | http://localhost:8081/api/sante |
| Adminer (base de données) | http://localhost:8080 |

## Base de données

À la première création du volume, `db/initialiser.sh` applique toutes les migrations puis charge
les données de démonstration (mot de passe des comptes : `Test1234!`) :

| Compte | Rôle |
|---|---|
| admin@upcycleconnect.fr | Administrateur |
| salarie@upcycleconnect.fr | Salarié (responsable) |
| particulier@upcycleconnect.fr | Particulier |
| pro@upcycleconnect.fr | Professionnel |

Pour faire évoluer le schéma, on ne modifie jamais une migration déjà appliquée : on ajoute un
nouveau fichier `db/migrations/013_....sql`, puis on l'applique sur la base existante :

```bash
docker compose exec postgres sh /db/migrer.sh
```

La table `schema_migrations` liste les migrations déjà appliquées. Pour repartir d'une base vide :
`docker compose down -v && docker compose up -d` (efface toutes les données locales).

### Base vide (déploiement)

Pour un vrai déploiement, sans aucune donnée de test, `db/creer_base_vide.sh` crée une base à côté de la base de travail :

```bash
sh db/creer_base_vide.sh                        # base « upcycleconnect_vide »
sh db/creer_base_vide.sh ma_base --remplacer    # recrée une base existante
```

Elle contient toutes les migrations et les référentiels (rôles, langues, matériaux avec leurs facteurs CO2, paramètres financiers), un conteneur par site de collecte et **un seul compte administrateur**. Le script demande son e-mail, son nom et son mot de passe au clavier (même règle que l'application, jamais écrit dans un fichier). Pour que l'API l'utilise : `DB_NAME=upcycleconnect_vide` dans `.env`, puis `docker compose up -d api`. Retirer la ligne pour revenir à la base de travail.

## Page d'accueil et inscription

`/` (`front/index.php`) est la page d'accueil publique : bandeau plein écran avec le logo et le slogan, le concept en 3 étapes, des chiffres, les 6 dernières annonces en ligne, les 3 prochains rendez-vous publiés, « Pour qui ? » et un appel à l'inscription. Les données viennent d'une seule route publique, `GET /api/accueil` (`api/internal/api/accueil.go`), qui ne renvoie aucune donnée personnelle. Un visiteur connecté voit « Mon espace » à la place des boutons de connexion.

L'inscription (`/inscription`) crée un compte particulier ou professionnel puis connecte directement la personne. Le profil peut être choisi depuis l'accueil (`/inscription?role=professionnel`), le mot de passe est saisi deux fois et le compte garde la langue choisie. L'e-mail est enregistré en minuscules et sans espaces, et la connexion ne tient pas compte des majuscules. La page de connexion propose « Pas encore de compte ? Créer un compte ».

## Adresses et pages d'erreur

Les adresses s'écrivent sans `.php` : `/connexion`, `/espace/annonces`, `/admin/tableau-de-bord?onglet=…`. Les règles sont dans `front/.htaccess` (Apache, conteneur `web`, avec `mod_rewrite` activé dans `docker-compose.yml`) et reprises dans `front/routeur.php` pour `php -S` ; une modification de l'un doit être reportée dans l'autre.

- une ancienne adresse en `.php` redirige définitivement (301) vers la nouvelle, paramètres conservés (anciens liens de notifications, favoris) ;
- les dossiers `inclus/` (morceaux de page) et `routeur.php` renvoient 404, la liste des fichiers d'un dossier renvoie 403 ;
- une page d'erreur unique, `front/erreur.php`, avec le fond et la carte de la page de connexion (403, 404, 500, 503, et un message générique pour les autres codes). Elle n'utilise ni Vue ni l'API, pour s'afficher même quand le reste est en panne ; textes dans `erreur.*` des fichiers de langue ;
- une erreur PHP dans une page affiche la page 500 au lieu d'une page blanche ou d'un message technique (`front/inclus/erreurs.php`, le détail va dans le journal du serveur) ;
- un utilisateur connecté qui ouvre une page réservée à un autre rôle arrive sur la page 403 (`exigerConnexion` dans `api.js`).

Pour essayer : `/nimporte-quoi` (404), `/assets/` (403), `/erreur?code=500`, ou `/admin/utilisateurs` connecté en particulier (403).

## Une application, un affichage par rôle

Particuliers, professionnels et salariés utilisent la même application (`front/espace/`). Le fichier
`front/espace/inclus/espace.php` liste chaque page avec les rôles qui y ont accès : le menu n'affiche
que les pages du rôle connecté, et une page refuse un rôle qui n'est pas dans sa liste. Dans une page
partagée, une partie propre à un rôle s'écrit `<div v-if="roleConnecte() === 'professionnel'">`.

Ces contrôles ne règlent que l'affichage : les droits sont vérifiés par l'API Go sur chaque route.

Pour ajouter une page : créer `front/espace/ma-page.php` et `front/espace/js/ma-page.js`, puis l'ajouter
dans la liste de `espace.php` avec son libellé, son icône et ses rôles.

## Multilingue

Le front utilise la bibliothèque [vue-i18n](https://vue-i18n.intlify.dev/) (v11, chargée par CDN).
Les textes sont dans un fichier JSON par langue : `front/assets/locales/fr.json`, `en.json`. La liste
des langues proposées reste en base (table `langues`, liée à la langue préférée de chaque utilisateur).

Dans un template Vue, un texte s'écrit `{{ $t('nav.mes_annonces') }}`, avec des paramètres
`{{ $t('accueil.bonjour', { prenom }) }}` ou une valeur par défaut `{{ $t('categorie.' + code, code) }}`.
Dans le JavaScript, on utilise `t(...)` de la même façon.

Pour ajouter une langue, sans modifier le code :

1. copier `front/assets/locales/fr.json` en `es.json` et traduire les textes ;
2. déclarer la langue en base :

```sql
INSERT INTO langues (code, libelle) VALUES ('es', 'Español');
```

Une clé absente d'une langue s'affiche en français (`fallbackLocale`).

## Achats et paiement

Quatre achats passent par le même parcours : payer une **formation** (Catalogue → « Payer et s'inscrire »), acheter un **objet en vente** (Annonces → « Acheter »), et pour les professionnels l'**abonnement Premium** (Mon compte → « Passer Premium ») et une **campagne publicitaire** (Contrats & Facturation → « Nouvelle campagne », 100 à 500 €/mois, l'élément apparaît « À la une »). Chaque achat produit une facture PDF ; l'abonnement et la campagne produisent aussi un contrat PDF. Le parcours est dans `api/internal/api/commandes.go` :

1. `POST /api/commandes` crée une commande en attente : l'objet est réservé, la place de formation est retenue pendant 30 minutes ;
2. paiement confirmé → `confirmerCommande()` : inscription payée ou objet cédé à l'acheteur, paiement enregistré avec `enregistrerPaiement()` (commission comprise pour les ventes), facture PDF archivée, vendeur ou formateur notifié ;
3. paiement abandonné ou délai dépassé → `annulerCommande()` : l'objet est remis en vente.

Tant que Stripe n'est pas branché, `PAIEMENT_FOURNISSEUR=simulation` (valeur par défaut) : la page `/espace/paiement` affiche un bandeau « Mode démonstration » et un bouton « Payer » qui confirme la commande sans aucun débit.

### Impact CO2 (statistiques des professionnels)

CO2 évité = poids de chaque objet récupéré × facteur de son matériau. Le particulier indique le matériau et un poids approximatif à la demande de dépôt, le salarié peut corriger le poids à la validation. Les facteurs (kg CO2e évités par kg réemployé) sont modifiables par l'admin (page Catégories) et sourcés : ADEME Base Carbone n° 20907 (bois/mobilier), 26729 (acier), 26741 (verre), 20836 (plastique), et Refashion, bilan environnemental 2026 (textile). C'est une estimation, présentée comme telle aux professionnels.

### Brancher Stripe

Tout le reste est déjà en place, seules ces étapes manquent :

1. Créer un compte Stripe (mode test) et mettre dans `.env` : `STRIPE_SECRET_KEY=sk_test_…`, `STRIPE_WEBHOOK_SECRET=whsec_…` et `PAIEMENT_FOURNISSEUR=stripe` (déclarer les deux clés dans `docker-compose.yml`, comme OneSignal).
2. Ajouter la librairie : `go get github.com/stripe/stripe-go/v82` (dans `api/`).
3. Dans `urlPaiement()` (`commandes.go`), cas `"stripe"` : créer une session **Stripe Checkout** (mode `payment`, montant `c.Montant` en centimes, libellé `c.Libelle`, `metadata.commande_id`, `success_url` et `cancel_url` vers `/espace/paiement?commande=ID`, expiration égale à `c.DateExpiration`), enregistrer l'identifiant de session dans `commandes.reference_externe`, renvoyer `session.URL`.
4. Créer la route `POST /api/stripe/webhook` : vérifier la signature avec `STRIPE_WEBHOOK_SECRET`, puis
   - `checkout.session.completed` → `s.confirmerCommande(commandeID, idDuPaiementStripe)` ;
   - `checkout.session.expired` → `s.annulerCommande(commandeID, "expiree")`.

   `confirmerCommande` ne crée jamais de second paiement si Stripe renvoie le même événement.
5. Abonnement et campagnes : pour un prélèvement **mensuel** automatique, utiliser le mode `subscription` de Checkout et l'événement `invoice.paid` pour enregistrer chaque mois un paiement (`enregistrerPaiement`) ; aujourd'hui seul le premier mois est payé.
6. Tester en local avec la Stripe CLI : `stripe listen --forward-to localhost:8081/api/stripe/webhook`, et la carte de test `4242 4242 4242 4242` (date future, CVC au choix).

La page `/espace/paiement` sert alors de page de retour : elle affiche le résultat de la commande (le bouton de démonstration disparaît, car `payer-simulation` refuse les commandes Stripe).

## Organisation Git

- `main` : version stable, modifiée uniquement par Pull Request
- `develop` : branche d'intégration
- `feature/<sujet>`, `fix/<sujet>`, `chore/<sujet>` : une branche par tâche, fusionnée dans `develop` par Pull Request

Messages de commit au format [Conventional Commits](https://www.conventionalcommits.org/fr/) :
`feat(api): ...`, `fix(front): ...`, `chore: ...`.
