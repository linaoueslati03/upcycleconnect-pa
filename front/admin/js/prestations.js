document.addEventListener('DOMContentLoaded', () => {
    const tbody = document.getElementById('prestations-list'); 

    fetch('http://localhost:8081/api/prestations')
        .then(response => {
            if (!response.ok) throw new Error("Erreur serveur");
            return response.json();
        })
        .then(data => {
            tbody.innerHTML = ''; 
            data.forEach(prest => {
                let row = `
                    <tr>
                        <td class="presta-title">${prest.titre}</td>
                        <td>${prest.categorie}</td>
                        <td>${prest.tarif} €</td>
                        <td>${prest.statut}</td>
                        <td>
                            <a href="#" class="orange-link" onclick="editer(${prest.id}, '${prest.titre.replace(/'/g, "\\'")}', '${prest.categorie}', ${prest.tarif}, '${prest.statut}')">Editer</a>
                            <a href="#" style="color: #d32f2f; margin-left: 15px; text-decoration: none; font-weight: bold;" onclick="supprimer(${prest.id})">Supprimer</a>
                        </td>
                    </tr>
                `;
                tbody.innerHTML += row;
            });
        })
        .catch(error => {
            console.error("Erreur API :", error);
            tbody.innerHTML = `<tr><td colspan="5" style="color:red;">Erreur de chargement</td></tr>`;
        });
});

document.getElementById('prestation-form').addEventListener('submit', function(e) {
    e.preventDefault();
    const id = document.getElementById('presta-id').value;
    const method = id ? 'PUT' : 'POST';
    const url = id ? `http://localhost:8081/api/prestations/${id}` : 'http://localhost:8081/api/prestations';

    const payload = {
        titre: document.getElementById('presta-titre').value,
        categorie: document.getElementById('presta-categorie').value,
        tarif: parseFloat(document.getElementById('presta-tarif').value),
        statut: document.getElementById('presta-statut').value
    };

    fetch(url, {
        method: method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
    }).then(() => {
        reinitialiserFormulaire();
        location.reload();
    });
});

function editer(id, titre, categorie, tarif, statut) {
    document.getElementById('form-title').innerText = "MODIFIER L'OFFRE";
    document.getElementById('presta-id').value = id;
    document.getElementById('presta-titre').value = titre;
    document.getElementById('presta-categorie').value = categorie;
    document.getElementById('presta-tarif').value = tarif;
    document.getElementById('presta-statut').value = statut;
}

function supprimer(id) {
    if(confirm("Supprimer cette prestation ?")) {
        fetch(`http://localhost:8081/api/prestations/${id}`, { method: 'DELETE' })
            .then(() => location.reload());
    }
}

function reinitialiserFormulaire() {
    document.getElementById('prestation-form').reset();
    document.getElementById('presta-id').value = '';
    document.getElementById('form-title').innerText = "NOUVELLE OFFRE";
}