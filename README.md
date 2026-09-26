# Endpoints de l'API REST

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
