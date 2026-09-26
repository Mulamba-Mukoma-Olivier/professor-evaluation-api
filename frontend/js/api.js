const API_URL = "http://localhost:8080";

async function apiFetch(endpoint, options = {}) {
    const token = localStorage.getItem("token");

    const headers = {
        "Content-Type": "application/json",
        ...options.headers
    };

    if (token) {
        headers.Authorization = `Bearer ${token}`;
    }

    const response = await fetch(`${API_URL}${endpoint}`, {
        ...options,
        headers
    });

    // Token expiré ou invalide
    if (response.status === 401) {
        localStorage.removeItem("token");
        localStorage.removeItem("user");

        window.location.href = "login.html";

        throw new Error("Session expirée.");
    }

    const contentType = response.headers.get("content-type");

    let data;

    if (contentType && contentType.includes("application/json")) {
        data = await response.json();
    } else {
        data = await response.text();
    }

    if (!response.ok) {
        throw new Error(
            data?.message ||
            data?.error ||
            "Une erreur est survenue."
        );
    }

    return data;
}


// Vérifier que l'utilisateur est connecté
function requireAuth() {
    const token = localStorage.getItem("token");

    if (!token) {
        window.location.href = "login.html";
        return false;
    }

    return true;
}


// Déconnexion
function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("user");

    window.location.href = "login.html";
}