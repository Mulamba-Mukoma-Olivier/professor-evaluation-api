// ==============================
// CONNEXION
// ==============================

const loginForm = document.getElementById("loginForm");

if (loginForm) {

    loginForm.addEventListener("submit", async (event) => {

        event.preventDefault();

        const email =
            document.getElementById("email").value.trim();

        const password =
            document.getElementById("password").value;

        const errorMessage =
            document.getElementById("errorMessage");

        errorMessage.textContent =
            "Connexion en cours...";

        try {

            const data = await apiFetch("/auth/login", {
                method: "POST",

                body: JSON.stringify({
                    email: email,
                    password: password
                })
            });


            // Selon la réponse du backend
            const token =
                data.token ||
                data.access_token;


            if (!token) {
                throw new Error(
                    "Le serveur n'a pas retourné de token."
                );
            }


            // Stocker le token
            localStorage.setItem(
                "token",
                token
            );


            // Stocker les informations utilisateur
            if (data.user) {

                localStorage.setItem(
                    "user",
                    JSON.stringify(data.user)
                );
            }


            // Aller au dashboard
            window.location.href =
                "dashboard.html";

        } catch (error) {

            errorMessage.textContent =
                error.message;
        }
    });
}



// ==============================
// INSCRIPTION
// ==============================

const registerForm =
    document.getElementById("registerForm");


if (registerForm) {

    registerForm.addEventListener(
        "submit",
        async (event) => {

            event.preventDefault();


            const matricule =
                document
                    .getElementById("matricule")
                    .value
                    .trim();


            const name =
                document
                    .getElementById("name")
                    .value
                    .trim();


            const email =
                document
                    .getElementById("email")
                    .value
                    .trim();


            const password =
                document
                    .getElementById("password")
                    .value;


            const message =
                document.getElementById(
                    "registerMessage"
                );


            message.className =
                "message";

            message.textContent =
                "Création du compte...";


            try {

                await apiFetch(
                    "/auth/register",
                    {
                        method: "POST",

                        body: JSON.stringify({

                            matricule: matricule,

                            name: name,

                            email: email,

                            password: password,

                            role: "student"
                        })
                    }
                );


                message.className =
                    "message success";

                message.textContent =
                    "Compte créé avec succès.";


                // Redirection vers login
                setTimeout(() => {

                    window.location.href =
                        "login.html";

                }, 1200);


            } catch (error) {

                message.className =
                    "message error";

                message.textContent =
                    error.message;
            }
        }
    );
}