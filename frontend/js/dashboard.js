// Vérifier l'authentification
if (!requireAuth()) {
    throw new Error("Utilisateur non authentifié.");
}


// ==============================
// UTILISATEUR
// ==============================

const user =
    JSON.parse(
        localStorage.getItem("user")
    );


if (user) {

    const welcome =
        document.getElementById("welcome");

    if (welcome) {

        welcome.textContent =
            `Bonjour ${user.name || ""} 👋`;
    }
}


// ==============================
// LOGOUT
// ==============================

const logoutBtn =
    document.getElementById("logoutBtn");


if (logoutBtn) {

    logoutBtn.addEventListener(
        "click",
        logout
    );
}



// ==============================
// CHARGER LES STATISTIQUES
// ==============================

async function loadDashboard() {

    try {

        // Professeurs
        const professors =
            await apiFetch("/professors");


        // Cours
        const courses =
            await apiFetch("/courses");


        const professorList =
            Array.isArray(professors)
                ? professors
                : professors.data || [];


        const courseList =
            Array.isArray(courses)
                ? courses
                : courses.data || [];


        // Afficher les nombres
        const professorCount =
            document.getElementById(
                "professorCount"
            );


        const courseCount =
            document.getElementById(
                "courseCount"
            );


        if (professorCount) {

            professorCount.textContent =
                professorList.length;
        }


        if (courseCount) {

            courseCount.textContent =
                courseList.length;
        }


    } catch (error) {

        console.error(
            "Erreur dashboard :",
            error
        );
    }
}


loadDashboard();