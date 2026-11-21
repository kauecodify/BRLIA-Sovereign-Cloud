// BRLIA - Portal JavaScript
// Minimalista, Profissional, Sem Emojis

const PORTAL = {
    name: "operations-center",
    title: "Centro de Operacoes",
    pages: "Alertas,Incidentes,Saude dos Servicos,Clusters,Backup e DR".split(","),
};

document.addEventListener("DOMContentLoaded", function () {
    renderNav();
    checkStatus();
    initNavigation();
});

function renderNav() {
    const container = document.getElementById("nav-sections");
    if (!container) return;

    const navTitle = document.createElement("div");
    navTitle.className = "nav-section-title";
    navTitle.textContent = PORTAL.title;
    container.appendChild(navTitle);

    const wrapper = document.createElement("div");
    PORTAL.pages.forEach(function (label) {
        const link = document.createElement("a");
        link.href = "#";
        link.className = "nav-item";
        link.innerHTML =
            '<span class="nav-icon">&#9679;</span><span>' + label + "</span>";
        link.addEventListener("click", function (event) {
            event.preventDefault();
            const pageTitle = document.getElementById("page-title");
            if (pageTitle) pageTitle.textContent = label;
        });
        wrapper.appendChild(link);
    });
    container.appendChild(wrapper);
}

function initNavigation() {
    const pageTitle = document.getElementById("page-title");
    document.querySelectorAll(".nav-item").forEach(function (item) {
        item.addEventListener("click", function () {
            document.querySelectorAll(".nav-item").forEach(function (el) {
                el.classList.remove("active");
            });
            item.classList.add("active");
            const label = item.textContent.trim();
            if (pageTitle && label) pageTitle.textContent = label;
        });
    });
}

function checkStatus() {
    const chip = document.getElementById("status-chip");
    const text = document.getElementById("platform-status");
    const raw = document.getElementById("raw-status");

    fetch("/api/health")
        .then(function (response) {
            if (!response.ok) throw new Error("HTTP " + response.status);
            return response.json();
        })
        .then(function (body) {
            if (chip) {
                chip.textContent = "operacional";
                chip.classList.add("ok");
            }
            if (text) {
                text.textContent =
                    "API de Controle respondeu com estado saudavel (servico: " +
                    (body.service || "desconhecido") +
                    ").";
            }
            if (raw) raw.textContent = JSON.stringify(body, null, 2);
        })
        .catch(function (error) {
            if (chip) {
                chip.textContent = "offline";
                chip.classList.add("down");
            }
            if (text) {
                text.textContent =
                    "Nao foi possivel consultar a API de Controle (" + error.message + ").";
            }
        });
}
