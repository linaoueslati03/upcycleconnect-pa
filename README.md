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
├── api/                 # API Go
├── db/
│   └── init/            # scripts SQL exécutés à la création de la base
├── front/
│   ├── admin/           # back-office administrateur
│   ├── salaries/        # espace salariés
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
| Front | http://localhost:8000/admin/utilisateurs.php |
| API | http://localhost:8081/api/sante |
| Adminer (base de données) | http://localhost:8080 |

Les scripts de `db/init/` ne s'exécutent qu'à la première création du volume PostgreSQL.

## Organisation Git

- `main` : version stable, modifiée uniquement par Pull Request
- `develop` : branche d'intégration
- `feature/<sujet>`, `fix/<sujet>`, `chore/<sujet>` : une branche par tâche, fusionnée dans `develop` par Pull Request

Messages de commit au format [Conventional Commits](https://www.conventionalcommits.org/fr/) :
`feat(api): ...`, `fix(front): ...`, `chore: ...`.
