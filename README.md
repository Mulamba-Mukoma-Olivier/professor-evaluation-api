# Endpoints de l'API REST

## Déploiement de l'API sur Vercel

Le projet contient une API Go avec son point d'entrée dans `cmd/server/main.go`. Pour le déployer :

1. Poussez le dépôt vers GitHub, GitLab ou Bitbucket, puis importez-le dans Vercel.
2. Vérifiez que le répertoire du projet est la racine du dépôt et que le Framework Preset est `Go`.
3. Ajoutez les variables d'environnement suivantes dans les réglages du projet Vercel :

   - `JWT_SECRET` : une clé aléatoire longue et unique.
   - `DB_DRIVER` : `postgres`.
   - `DATABASE_URL` : l'URL de connexion fournie par un hébergeur PostgreSQL, avec TLS (par exemple `sslmode=require`).
   - `CORS_ALLOWED_ORIGINS` : l'origine exacte du frontend autorisé, par exemple `https://mon-frontend.vercel.app`.

4. Déployez le projet. L'API et le contrôle de santé seront disponibles à l'URL Vercel, avec le contrôle sur `/health`.

En local, si `DB_DRIVER` n'est pas défini, l'application continue d'utiliser le fichier SQLite `professor_evaluation.db`. Sur Vercel, utilisez PostgreSQL : le fichier local du service n'est pas un stockage persistant. Gardez les identifiants dans les variables Vercel et ne les ajoutez pas au dépôt.

## 1. Santé de l'API

* `GET /health`

## 2. Authentification

* `POST /auth/register`
* `POST /auth/login`

## 3. Professeurs

* `GET /professors`
* `GET /professors/active`
* `GET /professors/by-status`
* `GET /professors/:id`
* `POST /professors`
* `PUT /professors/:id`
* `DELETE /professors/:id`

## 4. Cours

* `GET /courses`
* `GET /courses/:id`
* `POST /courses`
* `PUT /courses/:id`
* `DELETE /courses/:id`

## 5. Critères d'évaluation

* `GET /criteria`
* `GET /criteria/active`
* `GET /criteria/:id`
* `POST /criteria`
* `PUT /criteria/:id`
* `DELETE /criteria/:id`

## 6. Évaluations

* `GET /evaluations`
* `GET /evaluations/:id`
* `POST /evaluations`

## 7. Résultats

* `GET /results/professors/:professor_id`

## 8. Éligibilité

* `GET /eligibility/:student_id`
* `POST /eligibility`
* `PUT /eligibility/:student_id`
* `DELETE /eligibility/:student_id`

**Total : 29 endpoints.**
