// ==============================
// AUTHENTIFICATION
// ==============================

if (!requireAuth()) {
    throw new Error("Utilisateur non authentifié.");
}


// ==============================
// DÉCONNEXION
// ==============================

const logoutBtn = document.getElementById("logoutBtn");

if (logoutBtn) {
    logoutBtn.addEventListener("click", logout);
}


// ==============================
// VARIABLES
// ==============================

let criteria = [];

const params = new URLSearchParams(window.location.search);

const professorFromUrl = params.get("professor");


// ==============================
// CHARGER LA PAGE
// ==============================

async function init() {

    try {

        await loadProfessors();
        await loadCourses();
        await loadCriteria();

        if (professorFromUrl) {

            document.getElementById("professor").value =
                professorFromUrl;
        }

    } catch (error) {

        console.error(error);

        showMessage(
            error.message || "Erreur lors du chargement.",
            "danger"
        );
    }
}


// ==============================
// PROFESSEURS
// ==============================

async function loadProfessors() {

    const select = document.getElementById("professor");

    const response = await apiFetch("/professors");

    console.log("Professeurs :", response);


    // Selon la réponse de l'API
    const professors =
        response.data ||
        response.professors ||
        response;


    if (!Array.isArray(professors)) {

        throw new Error(
            "Format incorrect pour les professeurs."
        );
    }


    select.innerHTML = `
        <option value="">
            Sélectionner un professeur
        </option>
    `;


    professors.forEach(professor => {

        const option = document.createElement("option");

        option.value = professor.id;

        option.textContent =
            professor.name ||
            professor.full_name ||
            `Professeur ${professor.id}`;

        select.appendChild(option);
    });
}


// ==============================
// COURS
// ==============================

async function loadCourses() {

    const select = document.getElementById("course");

    console.log("Chargement des cours...");


    const response = await apiFetch("/courses");

    console.log("Réponse API /courses :", response);


    const courses =
        response.data ||
        response.courses ||
        response;


    if (!Array.isArray(courses)) {

        throw new Error(
            "Format incorrect pour les cours."
        );
    }


    select.innerHTML = `
        <option value="">
            Sélectionner un cours
        </option>
    `;


    if (courses.length === 0) {

        select.innerHTML += `
            <option value="">
                Aucun cours disponible
            </option>
        `;

        return;
    }


    courses.forEach(course => {

        const option = document.createElement("option");

        option.value = course.id;


        if (course.code && course.name) {

            option.textContent =
                `${course.code} - ${course.name}`;

        } else {

            option.textContent =
                course.name ||
                course.code ||
                `Cours ${course.id}`;
        }


        select.appendChild(option);
    });


    console.log(
        `${courses.length} cours affichés.`
    );
}


// ==============================
// CRITÈRES
// ==============================

async function loadCriteria() {

    const container =
        document.getElementById("criteriaContainer");


    const response =
        await apiFetch("/criteria");


    console.log("Critères :", response);


    criteria =
        response.data ||
        response.criteria ||
        response;


    if (!Array.isArray(criteria)) {

        throw new Error(
            "Format incorrect pour les critères."
        );
    }


    // Garder uniquement les critères actifs
    criteria = criteria.filter(
        criterion => criterion.active !== false
    );


    if (criteria.length === 0) {

        container.innerHTML = `
            <div class="alert alert-warning">
                Aucun critère disponible.
            </div>
        `;

        return;
    }


    container.innerHTML = "";


    criteria.forEach(criterion => {

        const maxScore =
            Number(criterion.max_score) || 5;


        let ratings = "";


        for (let score = 1; score <= maxScore; score++) {

            ratings += `
                <label class="me-3">

                    <input
                        type="radio"
                        name="criterion-${criterion.id}"
                        value="${score}"
                        required
                    >

                    ${score}

                </label>
            `;
        }


        container.innerHTML += `

            <div class="border rounded p-3 mb-3">

                <h5 class="fw-bold">
                    ${criterion.name}
                </h5>

                <p class="text-secondary mb-3">
                    ${criterion.description || ""}
                </p>

                <div>
                    ${ratings}
                </div>

            </div>

        `;
    });
}


// ==============================
// MESSAGE
// ==============================

function showMessage(text, type) {

    const message =
        document.getElementById(
            "evaluationMessage"
        );


    message.className =
        `alert alert-${type}`;


    message.textContent = text;
}


// ==============================
// ENVOYER L'ÉVALUATION
// ==============================

const evaluationForm =
    document.getElementById(
        "evaluationForm"
    );


evaluationForm.addEventListener(
    "submit",
    async function(event) {

        event.preventDefault();


        const professorId =
            document.getElementById(
                "professor"
            ).value;


        const courseId =
            document.getElementById(
                "course"
            ).value;


        const comment =
            document.getElementById(
                "comment"
            ).value.trim();


        // ------------------------------
        // Vérifier professeur
        // ------------------------------

        if (!professorId) {

            showMessage(
                "Veuillez sélectionner un professeur.",
                "danger"
            );

            return;
        }


        // ------------------------------
        // Vérifier cours
        // ------------------------------

        if (!courseId) {

            showMessage(
                "Veuillez sélectionner un cours.",
                "danger"
            );

            return;
        }


        // ------------------------------
        // Récupérer les notes
        // ------------------------------

        const answers = [];


        for (const criterion of criteria) {

            const selected =
                document.querySelector(
                    `input[name="criterion-${criterion.id}"]:checked`
                );


            if (!selected) {

                showMessage(
                    `Veuillez noter : ${criterion.name}`,
                    "danger"
                );

                return;
            }


            answers.push({

                criterion_id: Number(criterion.id),

                score: Number(selected.value)

            });
        }


        // ------------------------------
        // Envoyer
        // ------------------------------

        try {

            showMessage(
                "Envoi de l'évaluation...",
                "info"
            );


            await apiFetch(
                "/evaluations", {

                    method: "POST",

                    body: JSON.stringify({

                        professor_id: Number(professorId),

                        course_id: Number(courseId),

                        comment: comment,

                        answers: answers

                    })
                }
            );


            showMessage(
                "Évaluation envoyée avec succès.",
                "success"
            );


            evaluationForm.reset();


        } catch (error) {

            console.error(
                "Erreur :",
                error
            );


            showMessage(
                error.message ||
                "Une erreur est survenue.",
                "danger"
            );
        }
    }
);


// ==============================
// DÉMARRER
// ==============================

init();