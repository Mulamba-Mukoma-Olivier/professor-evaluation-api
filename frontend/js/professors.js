// Vérifier l'authentification
if (!requireAuth()) {
    throw new Error("Utilisateur non authentifié.");
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


// Liste globale
let professors = [];


// ==============================
// CHARGER LES PROFESSEURS
// ==============================

async function loadProfessors() {

    const container =
        document.getElementById(
            "professorsList"
        );


    try {

        const data =
            await apiFetch("/professors");


        professors =
            Array.isArray(data)
                ? data
                : data.data || [];


        displayProfessors(
            professors
        );


    } catch (error) {

        container.innerHTML = `
            <p class="message error">
                ${error.message}
            </p>
        `;
    }
}



// ==============================
// AFFICHER LES PROFESSEURS
// ==============================

function displayProfessors(list) {

    const container =
        document.getElementById(
            "professorsList"
        );


    if (!list.length) {

        container.innerHTML = `
            <p class="muted">
                Aucun professeur trouvé.
            </p>
        `;

        return;
    }


    container.innerHTML =
        list.map(professor => {

            const id =
                professor.id;


            const name =
                professor.name ||
                professor.full_name ||
                "Professeur";


            const department =
                professor.department ||
                "Département non spécifié";


            return `

                <article class="professor-card">

                    <div class="professor-avatar">
                        👨‍🏫
                    </div>

                    <h3>
                        ${name}
                    </h3>

                    <p>
                        ${department}
                    </p>

                    <a
                        href="evaluation.html?professor=${id}"
                        class="btn btn-primary"
                    >
                        Évaluer
                    </a>

                </article>

            `;

        }).join("");
}



// ==============================
// RECHERCHE
// ==============================

const searchInput =
    document.getElementById(
        "searchProfessor"
    );


if (searchInput) {

    searchInput.addEventListener(
        "input",
        (event) => {

            const search =
                event.target.value
                    .toLowerCase()
                    .trim();


            const filtered =
                professors.filter(
                    professor => {

                        const name =
                            (
                                professor.name ||
                                professor.full_name ||
                                ""
                            )
                            .toLowerCase();


                        return name.includes(
                            search
                        );
                    }
                );


            displayProfessors(
                filtered
            );
        }
    );
}


// Lancer le chargement
loadProfessors();