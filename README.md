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

## Organisation Git

- `main` : version stable, modifiée uniquement par Pull Request
- `develop` : branche d'intégration
- `feature/<sujet>`, `fix/<sujet>`, `chore/<sujet>` : une branche par tâche, fusionnée dans `develop` par Pull Request

Messages de commit au format [Conventional Commits](https://www.conventionalcommits.org/fr/) :
`feat(api): ...`, `fix(front): ...`, `chore: ...`.
