### Endpoints de l'API

#### Health

* `GET /health`

#### Authentification

* `POST /auth/register`
* `POST /auth/login`

#### Professeurs

* `GET /professors`
* `GET /professors/active`
* `GET /professors/:id`
* `POST /professors`

#### Cours

* `GET /courses`
* `GET /courses/:id`
* `POST /courses`

#### Critères

* `GET /criteria`
* `GET /criteria/active`
* `GET /criteria/:id`
* `POST /criteria`

#### Évaluations

* `GET /evaluations`
* `GET /evaluations/:id`
* `POST /evaluations`
* `DELETE /evaluations/:id`

#### Résultats

* `GET /results/professors/:professor_id`

#### Éligibilité

* `GET /eligibility/:student_id`
