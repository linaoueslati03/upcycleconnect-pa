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
php -S localhost:8000 -t front
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

Deux achats sont possibles : payer une **formation** (Catalogue → « Payer et s'inscrire ») et acheter un **objet en vente** (Annonces → « Acheter »). Le parcours est dans `api/internal/api/commandes.go` :

1. `POST /api/commandes` crée une commande en attente : l'objet est réservé, la place de formation est retenue pendant 30 minutes ;
2. paiement confirmé → `confirmerCommande()` : inscription payée ou objet cédé à l'acheteur, paiement enregistré avec `enregistrerPaiement()` (commission comprise pour les ventes), facture PDF archivée, vendeur ou formateur notifié ;
3. paiement abandonné ou délai dépassé → `annulerCommande()` : l'objet est remis en vente.

Tant que Stripe n'est pas branché, `PAIEMENT_FOURNISSEUR=simulation` (valeur par défaut) : la page `/espace/paiement.php` affiche un bandeau « Mode démonstration » et un bouton « Payer » qui confirme la commande sans aucun débit.

### Brancher Stripe

Tout le reste est déjà en place, seules ces étapes manquent :

1. Créer un compte Stripe (mode test) et mettre dans `.env` : `STRIPE_SECRET_KEY=sk_test_…`, `STRIPE_WEBHOOK_SECRET=whsec_…` et `PAIEMENT_FOURNISSEUR=stripe` (déclarer les deux clés dans `docker-compose.yml`, comme OneSignal).
2. Ajouter la librairie : `go get github.com/stripe/stripe-go/v82` (dans `api/`).
3. Dans `urlPaiement()` (`commandes.go`), cas `"stripe"` : créer une session **Stripe Checkout** (mode `payment`, montant `c.Montant` en centimes, libellé `c.Libelle`, `metadata.commande_id`, `success_url` et `cancel_url` vers `/espace/paiement.php?commande=ID`, expiration égale à `c.DateExpiration`), enregistrer l'identifiant de session dans `commandes.reference_externe`, renvoyer `session.URL`.
4. Créer la route `POST /api/stripe/webhook` : vérifier la signature avec `STRIPE_WEBHOOK_SECRET`, puis
   - `checkout.session.completed` → `s.confirmerCommande(commandeID, idDuPaiementStripe)` ;
   - `checkout.session.expired` → `s.annulerCommande(commandeID, "expiree")`.

   `confirmerCommande` ne crée jamais de second paiement si Stripe renvoie le même événement.
5. Tester en local avec la Stripe CLI : `stripe listen --forward-to localhost:8081/api/stripe/webhook`, et la carte de test `4242 4242 4242 4242` (date future, CVC au choix).

La page `/espace/paiement.php` sert alors de page de retour : elle affiche le résultat de la commande (le bouton de démonstration disparaît, car `payer-simulation` refuse les commandes Stripe).

## Organisation Git

- `main` : version stable, modifiée uniquement par Pull Request
- `develop` : branche d'intégration
- `feature/<sujet>`, `fix/<sujet>`, `chore/<sujet>` : une branche par tâche, fusionnée dans `develop` par Pull Request

Messages de commit au format [Conventional Commits](https://www.conventionalcommits.org/fr/) :
`feat(api): ...`, `fix(front): ...`, `chore: ...`.
