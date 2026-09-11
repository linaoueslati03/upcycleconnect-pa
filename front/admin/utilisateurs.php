<!DOCTYPE html>
<html lang="fr">

<head>
    <meta charset="utf-8">

    <title>Gestion des utilisateurs</title>

    <link rel="stylesheet" href="../pages/style.css">

    <script src="https://cdn.tailwindcss.com"></script>
    <script src="../js/config.js"></script>
</head>

<body class="p-8">

    <h1 class="text-2xl font-display font-semibold text-brand-navy mb-6">Gestion des utilisateurs</h1>

    <button id="bouton-nouveau" class="mb-4 px-4 py-2 rounded-lg bg-brand-green text-white font-semibold">
        + Nouvel utilisateur
    </button>

    <table class="w-full bg-white rounded-lg overflow-hidden shadow">
        <thead class="bg-brand-navy text-white text-left">
            <tr>
                <th class="p-3">ID</th>
                <th class="p-3">Nom</th>
                <th class="p-3">Prénom</th>
                <th class="p-3">Email</th>
                <th class="p-3">Rôle</th>
                <th class="p-3">Statut</th>
                <th class="p-3">Score</th>
                <th class="p-3">Actions</th>
            </tr>
        </thead>
        <tbody id="corps-tableau"></tbody>
    </table>

    <form id="formulaire-utilisateur" class="hidden mt-6 max-w-md flex flex-col gap-3 bg-white p-6 rounded-lg shadow">
        <h2 id="titre-formulaire" class="text-lg font-semibold text-brand-navy">Nouvel utilisateur</h2>

        <input type="hidden" id="champ-id">

        <label class="text-sm font-medium">Nom
            <input type="text" id="champ-nom" class="w-full border rounded px-3 py-2" required>
        </label>

        <label class="text-sm font-medium">Prénom
            <input type="text" id="champ-prenom" class="w-full border rounded px-3 py-2" required>
        </label>

        <label class="text-sm font-medium">Email
            <input type="email" id="champ-email" class="w-full border rounded px-3 py-2" required>
        </label>

        <label class="text-sm font-medium">Rôle
            <select id="champ-role" class="w-full border rounded px-3 py-2">
                <option value="1">Particulier</option>
                <option value="2">Professionnel</option>
                <option value="3">Salarié</option>
                <option value="4">Administrateur</option>
            </select>
        </label>

        <label id="conteneur-mot-de-passe" class="text-sm font-medium">Mot de passe
            <input type="password" id="champ-mot-de-passe" class="w-full border rounded px-3 py-2">
        </label>

        <p id="erreur-formulaire" class="text-red-600 text-sm hidden"></p>

        <div class="flex gap-2">
            <button type="submit" class="px-4 py-2 rounded-lg bg-brand-green text-white font-semibold">Enregistrer</button>
            <button type="button" id="bouton-annuler" class="px-4 py-2 rounded-lg bg-gray-200">Annuler</button>
        </div>
    </form>

    <script src="js/utilisateurs.js"></script>
</body>

</html>
