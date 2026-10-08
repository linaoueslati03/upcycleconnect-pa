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
│   ├── admin/           # back-office administrateur
│   ├── salaries/        # espace salariés
│   ├── particuliers/    # espace particuliers
│   ├── inclus/          # en-tête et menu communs aux espaces
│   └── assets/          # css, js, images partagés
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
nouveau fichier `db/migrations/009_....sql`, puis on l'applique sur la base existante :

```bash
docker compose exec postgres sh /db/migrer.sh
```

La table `schema_migrations` liste les migrations déjà appliquées. Pour repartir d'une base vide :
`docker compose down -v && docker compose up -d` (efface toutes les données locales).

## Organisation Git

- `main` : version stable, modifiée uniquement par Pull Request
- `develop` : branche d'intégration
- `feature/<sujet>`, `fix/<sujet>`, `chore/<sujet>` : une branche par tâche, fusionnée dans `develop` par Pull Request

Messages de commit au format [Conventional Commits](https://www.conventionalcommits.org/fr/) :
`feat(api): ...`, `fix(front): ...`, `chore: ...`.
