const NOMS_ROLES = { 1: "Particulier", 2: "Professionnel", 3: "Salarié", 4: "Administrateur" };

let utilisateursActuels = [];

async function chargerUtilisateurs() {
    const reponse = await fetch(`${API_BASE_URL}/utilisateurs`);
    utilisateursActuels = await reponse.json();

    const corps = document.getElementById("corps-tableau");
    corps.innerHTML = "";

    for (const u of utilisateursActuels) {
        const ligne = document.createElement("tr");
        ligne.className = "border-t";
        ligne.innerHTML = `
            <td class="p-3">${u.id}</td>
            <td class="p-3">${u.nom}</td>
            <td class="p-3">${u.prenom}</td>
            <td class="p-3">${u.email}</td>
            <td class="p-3">${NOMS_ROLES[u.role_id] ?? u.role_id}</td>
            <td class="p-3">${u.statut}</td>
            <td class="p-3">${u.upcycling_score}</td>
            <td class="p-3">
                <button class="text-brand-navy underline mr-2" onclick="ouvrirFormulaireEditionParId(${u.id})">Modifier</button>
                <button class="text-red-600 underline" onclick="supprimerUtilisateur(${u.id})">Supprimer</button>
            </td>
        `;
        corps.appendChild(ligne);
    }
}

function ouvrirFormulaireEditionParId(id) {
    const utilisateur = utilisateursActuels.find((u) => u.id === id);
    if (utilisateur) {
        ouvrirFormulaireEdition(utilisateur);
    }
}

function ouvrirFormulaireCreation() {
    document.getElementById("titre-formulaire").textContent = "Nouvel utilisateur";
    document.getElementById("champ-id").value = "";
    document.getElementById("champ-nom").value = "";
    document.getElementById("champ-prenom").value = "";
    document.getElementById("champ-email").value = "";
    document.getElementById("champ-role").value = "1";
    document.getElementById("champ-mot-de-passe").value = "";
    document.getElementById("conteneur-mot-de-passe").classList.remove("hidden");
    document.getElementById("erreur-formulaire").classList.add("hidden");
    document.getElementById("formulaire-utilisateur").classList.remove("hidden");
}

function ouvrirFormulaireEdition(u) {
    document.getElementById("titre-formulaire").textContent = "Modifier l'utilisateur";
    document.getElementById("champ-id").value = u.id;
    document.getElementById("champ-nom").value = u.nom;
    document.getElementById("champ-prenom").value = u.prenom;
    document.getElementById("champ-email").value = u.email;
    document.getElementById("champ-role").value = u.role_id;
    document.getElementById("champ-mot-de-passe").value = "";
    document.getElementById("conteneur-mot-de-passe").classList.add("hidden");
    document.getElementById("erreur-formulaire").classList.add("hidden");
    document.getElementById("formulaire-utilisateur").classList.remove("hidden");
}

function fermerFormulaire() {
    document.getElementById("formulaire-utilisateur").classList.add("hidden");
}

async function soumettreFormulaire(evenement) {
    evenement.preventDefault();

    const id = document.getElementById("champ-id").value;
    const donnees = {
        nom: document.getElementById("champ-nom").value,
        prenom: document.getElementById("champ-prenom").value,
        email: document.getElementById("champ-email").value,
        role_id: parseInt(document.getElementById("champ-role").value, 10),
    };

    let url = `${API_BASE_URL}/utilisateurs`;
    let methode = "POST";
    if (id) {
        url = `${API_BASE_URL}/utilisateurs/${id}`;
        methode = "PUT";
    } else {
        donnees.mot_de_passe = document.getElementById("champ-mot-de-passe").value;
    }

    const reponse = await fetch(url, {
        method: methode,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(donnees),
    });

    const resultat = await reponse.json();

    if (!reponse.ok) {
        const erreur = document.getElementById("erreur-formulaire");
        erreur.textContent = resultat.erreur;
        erreur.classList.remove("hidden");
        return;
    }

    fermerFormulaire();
    chargerUtilisateurs();
}

async function supprimerUtilisateur(id) {
    if (!confirm("Supprimer cet utilisateur ?")) {
        return;
    }
    await fetch(`${API_BASE_URL}/utilisateurs/${id}`, { method: "DELETE" });
    chargerUtilisateurs();
}

document.addEventListener("DOMContentLoaded", () => {
    chargerUtilisateurs();
    document.getElementById("bouton-nouveau").addEventListener("click", ouvrirFormulaireCreation);
    document.getElementById("bouton-annuler").addEventListener("click", fermerFormulaire);
    document.getElementById("formulaire-utilisateur").addEventListener("submit", soumettreFormulaire);
});
