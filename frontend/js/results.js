document.addEventListener("DOMContentLoaded", function() {

    var professorSelect = document.getElementById("professorSelect");
    var courseSelect = document.getElementById("courseSelect");
    var academicYearSelect = document.getElementById("academicYearSelect");
    var periodSelect = document.getElementById("periodSelect");
    var loadResultsBtn = document.getElementById("loadResultsBtn");
    var resultsContainer = document.getElementById("resultsContainer");
    var filterMessage = document.getElementById("filterMessage");

    var professors = [];
    var courses = [];


    // =========================================================
    // CHARGER LES PROFESSEURS
    // =========================================================

    function loadProfessors() {

        professorSelect.innerHTML =
            '<option value="">Chargement...</option>';

        apiFetch("/professors")
            .then(function(response) {

                console.log("Réponse /professors :", response);

                if (response.data) {
                    professors = response.data;
                } else if (Array.isArray(response)) {
                    professors = response;
                } else {
                    professors = [];
                }

                professorSelect.innerHTML =
                    '<option value="">Sélectionner un professeur</option>';

                for (var i = 0; i < professors.length; i++) {

                    var professor = professors[i];

                    var option = document.createElement("option");

                    option.value = professor.id;

                    option.textContent =
                        professor.name ||
                        professor.full_name ||
                        professor.nom ||
                        "Professeur #" + professor.id;

                    professorSelect.appendChild(option);
                }

            })
            .catch(function(error) {

                console.error("Erreur professeurs :", error);

                professorSelect.innerHTML =
                    '<option value="">Erreur de chargement</option>';

                showMessage(
                    "Impossible de charger les professeurs.",
                    "danger"
                );
            });
    }


    // =========================================================
    // CHARGER LES COURS
    // =========================================================

    function loadCourses() {

        courseSelect.innerHTML =
            '<option value="">Chargement...</option>';

        courseSelect.disabled = true;

        apiFetch("/courses")
            .then(function(response) {

                console.log("Réponse /courses :", response);

                if (response.data) {
                    courses = response.data;
                } else if (response.courses) {
                    courses = response.courses;
                } else if (Array.isArray(response)) {
                    courses = response;
                } else {
                    courses = [];
                }

                courseSelect.innerHTML =
                    '<option value="">Sélectionner un cours</option>';

                for (var i = 0; i < courses.length; i++) {

                    var course = courses[i];

                    var option = document.createElement("option");

                    option.value = course.id;

                    var courseName = course.name || "Cours";
                    var courseCode = course.code || "";

                    if (courseCode !== "") {
                        option.textContent =
                            courseCode + " - " + courseName;
                    } else {
                        option.textContent = courseName;
                    }

                    courseSelect.appendChild(option);
                }

                courseSelect.disabled = false;

            })
            .catch(function(error) {

                console.error("Erreur cours :", error);

                courseSelect.innerHTML =
                    '<option value="">Erreur de chargement</option>';

                courseSelect.disabled = true;

                showMessage(
                    "Impossible de charger les cours.",
                    "danger"
                );
            });
    }


    // =========================================================
    // VÉRIFIER LES FILTRES
    // =========================================================

    function checkFilters() {

        var professorId = professorSelect.value;
        var courseId = courseSelect.value;
        var academicYear = academicYearSelect.value;
        var period = periodSelect.value;

        if (
            professorId !== "" &&
            courseId !== "" &&
            academicYear !== "" &&
            period !== ""
        ) {

            loadResultsBtn.disabled = false;

            filterMessage.textContent =
                "Tous les filtres sont sélectionnés.";

        } else {

            loadResultsBtn.disabled = true;

            filterMessage.textContent =
                "Sélectionnez un professeur, un cours, une année académique et une période.";
        }
    }


    // =========================================================
    // CHARGER LES RÉSULTATS
    // =========================================================

    function loadResults() {

        var professorId = professorSelect.value;
        var courseId = courseSelect.value;
        var academicYear = academicYearSelect.value;
        var period = periodSelect.value;


        if (professorId === "") {

            showMessage(
                "Veuillez sélectionner un professeur.",
                "warning"
            );

            return;
        }


        if (courseId === "") {

            showMessage(
                "Veuillez sélectionner un cours.",
                "warning"
            );

            return;
        }


        if (academicYear === "") {

            showMessage(
                "Veuillez sélectionner une année académique.",
                "warning"
            );

            return;
        }


        if (period === "") {

            showMessage(
                "Veuillez sélectionner une période.",
                "warning"
            );

            return;
        }


        resultsContainer.innerHTML =
            '<div class="text-center p-4">' +
            '<div class="spinner-border" role="status"></div>' +
            '<p class="mt-2">Chargement des résultats...</p>' +
            '</div>';


        /*
         * Le backend attend :
         *
         * /results/professors/:professor_id
         * ?course_id=...
         * &academic_year=...
         * &period=...
         */

        var url =
            "/results/professors/" +
            professorId +
            "?course_id=" +
            encodeURIComponent(courseId) +
            "&academic_year=" +
            encodeURIComponent(academicYear) +
            "&period=" +
            encodeURIComponent(period);


        console.log("URL résultats :", url);


        apiFetch(url)
            .then(function(response) {

                console.log(
                    "Réponse résultats :",
                    response
                );

                displayResults(response);

            })
            .catch(function(error) {

                console.error(
                    "Erreur résultat :",
                    error
                );

                resultsContainer.innerHTML = "";

                showMessage(
                    error.message ||
                    "Impossible de charger les résultats.",
                    "danger"
                );
            });
    }


    // =========================================================
    // AFFICHER LES RÉSULTATS
    // =========================================================

    function displayResults(result) {

        resultsContainer.innerHTML = "";


        if (!result) {

            showMessage(
                "Aucun résultat disponible.",
                "info"
            );

            return;
        }


        var professorId = result.professor_id;
        var courseId = result.course_id;
        var academicYear = result.academic_year;
        var period = result.period;

        var totalReviews = result.total_reviews;
        var globalAverage = result.global_average;


        // -----------------------------------------------------
        // RÉSUMÉ
        // -----------------------------------------------------

        var summary = document.createElement("div");

        summary.className = "row mb-4";


        summary.innerHTML =

            '<div class="col-md-3 mb-3">' +
            '<div class="card shadow-sm h-100">' +
            '<div class="card-body">' +
            '<h6 class="text-muted">Professeur</h6>' +
            '<h4>' +
            getProfessorName(professorId) +
            '</h4>' +
            '</div>' +
            '</div>' +
            '</div>' +

            '<div class="col-md-3 mb-3">' +
            '<div class="card shadow-sm h-100">' +
            '<div class="card-body">' +
            '<h6 class="text-muted">Cours</h6>' +
            '<h4>' +
            getCourseName(courseId) +
            '</h4>' +
            '</div>' +
            '</div>' +
            '</div>' +

            '<div class="col-md-3 mb-3">' +
            '<div class="card shadow-sm h-100">' +
            '<div class="card-body">' +
            '<h6 class="text-muted">Évaluations</h6>' +
            '<h4>' +
            formatNumber(totalReviews) +
            '</h4>' +
            '</div>' +
            '</div>' +
            '</div>' +

            '<div class="col-md-3 mb-3">' +
            '<div class="card shadow-sm h-100">' +
            '<div class="card-body">' +
            '<h6 class="text-muted">Moyenne globale</h6>' +
            '<h4>' +
            formatAverage(globalAverage) +
            ' / 5' +
            '</h4>' +
            '</div>' +
            '</div>' +
            '</div>';


        resultsContainer.appendChild(summary);


        // -----------------------------------------------------
        // ANNÉE ET PÉRIODE
        // -----------------------------------------------------

        var periodInfo = document.createElement("div");

        periodInfo.className =
            "alert alert-light border";


        periodInfo.innerHTML =
            "<strong>Année académique :</strong> " +
            (academicYear || "N/A") +
            " &nbsp;&nbsp; " +
            "<strong>Période :</strong> " +
            formatPeriod(period);


        resultsContainer.appendChild(periodInfo);


        // -----------------------------------------------------
        // CRITÈRES
        // -----------------------------------------------------

        var criteria = result.criteria;


        if (!criteria || criteria.length === 0) {

            var noCriteria =
                document.createElement("div");

            noCriteria.className =
                "alert alert-info";

            noCriteria.textContent =
                "Aucun résultat par critère n'est disponible.";

            resultsContainer.appendChild(noCriteria);

            return;
        }


        var title =
            document.createElement("h5");

        title.className =
            "mt-4 mb-3";

        title.textContent =
            "Résultats par critère";

        resultsContainer.appendChild(title);


        // -----------------------------------------------------
        // TABLEAU
        // -----------------------------------------------------

        var table =
            document.createElement("table");

        table.className =
            "table table-bordered table-hover align-middle";


        var thead =
            document.createElement("thead");

        thead.innerHTML =
            "<tr>" +
            "<th>Critère</th>" +
            "<th>Moyenne</th>" +
            "<th>Nombre de réponses</th>" +
            "</tr>";

        table.appendChild(thead);


        var tbody =
            document.createElement("tbody");


        for (var i = 0; i < criteria.length; i++) {

            var criterion = criteria[i];

            var row =
                document.createElement("tr");


            var criterionId =
                criterion.criterion_id;

            var average =
                criterion.average;

            var responses =
                criterion.responses;


            row.innerHTML =
                "<td>" +
                getCriterionName(criterionId) +
                "</td>" +

                "<td>" +
                formatAverage(average) +
                " / 5" +
                "</td>" +

                "<td>" +
                formatNumber(responses) +
                "</td>";


            tbody.appendChild(row);
        }


        table.appendChild(tbody);

        resultsContainer.appendChild(table);
    }


    // =========================================================
    // NOM DU PROFESSEUR
    // =========================================================

    function getProfessorName(id) {

        for (var i = 0; i < professors.length; i++) {

            if (
                Number(professors[i].id) ===
                Number(id)
            ) {

                return (
                    professors[i].name ||
                    professors[i].full_name ||
                    professors[i].nom ||
                    "Professeur #" + id
                );
            }
        }

        return "Professeur #" + id;
    }


    // =========================================================
    // NOM DU COURS
    // =========================================================

    function getCourseName(id) {

        for (var i = 0; i < courses.length; i++) {

            if (
                Number(courses[i].id) ===
                Number(id)
            ) {

                var code =
                    courses[i].code || "";

                var name =
                    courses[i].name || "Cours";


                if (code !== "") {

                    return code + " - " + name;
                }

                return name;
            }
        }

        return "Cours #" + id;
    }


    // =========================================================
    // NOM DU CRITÈRE
    // =========================================================

    function getCriterionName(id) {

        return "Critère #" + id;
    }


    // =========================================================
    // FORMATTER UNE MOYENNE
    // =========================================================

    function formatAverage(value) {

        if (
            value === null ||
            value === undefined ||
            value === ""
        ) {

            return "0.00";
        }


        var number = Number(value);


        if (isNaN(number)) {

            return "0.00";
        }


        return number.toFixed(2);
    }


    // =========================================================
    // FORMATTER UN NOMBRE
    // =========================================================

    function formatNumber(value) {

        if (
            value === null ||
            value === undefined ||
            value === ""
        ) {

            return "0";
        }


        var number = Number(value);


        if (isNaN(number)) {

            return "0";
        }


        return number.toString();
    }


    // =========================================================
    // FORMATTER LA PÉRIODE
    // =========================================================

    function formatPeriod(period) {

        if (period === "S1") {

            return "Semestre 1";
        }


        if (period === "S2") {

            return "Semestre 2";
        }


        return period || "N/A";
    }


    // =========================================================
    // AFFICHER UN MESSAGE
    // =========================================================

    function showMessage(message, type) {

        if (!filterMessage) {
            return;
        }


        filterMessage.className =
            "alert alert-" + type;

        filterMessage.textContent =
            message;
    }


    // =========================================================
    // ÉVÉNEMENTS
    // =========================================================

    professorSelect.addEventListener(
        "change",
        function() {

            checkFilters();
        }
    );


    courseSelect.addEventListener(
        "change",
        function() {

            checkFilters();
        }
    );


    academicYearSelect.addEventListener(
        "change",
        function() {

            checkFilters();
        }
    );


    periodSelect.addEventListener(
        "change",
        function() {

            checkFilters();
        }
    );


    loadResultsBtn.addEventListener(
        "click",
        function() {

            loadResults();
        }
    );


    // =========================================================
    // INITIALISATION
    // =========================================================

    loadProfessors();

    loadCourses();

    checkFilters();

});