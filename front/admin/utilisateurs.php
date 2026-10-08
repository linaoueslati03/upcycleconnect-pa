<!DOCTYPE html>
<html lang="fr">

<head>
    <meta charset="utf-8">

    <title>Gestion des utilisateurs</title>

    <link rel="stylesheet" href="../pages/style.css">

    <script src="https://cdn.tailwindcss.com"></script>
    <script src="../js/config.js"></script>
    <script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
</head>

<body class="p-8">

    <div id="app">

        <h1 class="text-2xl font-display font-semibold text-brand-navy mb-6">Gestion des utilisateurs</h1>

        <button @click="ouvrirFormulaireCreation"
            class="mb-4 px-4 py-2 rounded-lg bg-brand-green text-white font-semibold">
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
            <tbody>
                <tr v-for="u in utilisateurs" :key="u.id" class="border-t">
                    <td class="p-3">{{ u.id }}</td>
                    <td class="p-3">{{ u.nom }}</td>
                    <td class="p-3">{{ u.prenom }}</td>
                    <td class="p-3">{{ u.email }}</td>
                    <td class="p-3">{{ nomsRoles[u.role_id] ?? u.role_id }}</td>
                    <td class="p-3">{{ u.statut }}</td>
                    <td class="p-3">{{ u.upcycling_score }}</td>
                    <td class="p-3">
                        <button @click="ouvrirFormulaireEdition(u)" class="text-brand-navy underline mr-2">Modifier</button>
                        <button @click="supprimerUtilisateur(u.id)" class="text-red-600 underline">Supprimer</button>
                    </td>
                </tr>
            </tbody>
        </table>

        <form v-if="formulaireVisible" @submit.prevent="soumettreFormulaire"
            class="mt-6 max-w-md flex flex-col gap-3 bg-white p-6 rounded-lg shadow">

            <h2 class="text-lg font-semibold text-brand-navy">
                {{ formulaire.id ? "Modifier l'utilisateur" : "Nouvel utilisateur" }}
            </h2>

            <label class="text-sm font-medium">Nom
                <input type="text" v-model="formulaire.nom" class="w-full border rounded px-3 py-2" required>
            </label>

            <label class="text-sm font-medium">Prénom
                <input type="text" v-model="formulaire.prenom" class="w-full border rounded px-3 py-2" required>
            </label>

            <label class="text-sm font-medium">Email
                <input type="email" v-model="formulaire.email" class="w-full border rounded px-3 py-2" required>
            </label>

            <label class="text-sm font-medium">Rôle
                <select v-model="formulaire.role_id" class="w-full border rounded px-3 py-2">
                    <option :value="1">Particulier</option>
                    <option :value="2">Professionnel</option>
                    <option :value="3">Salarié</option>
                    <option :value="4">Administrateur</option>
                </select>
            </label>

            <label v-if="!formulaire.id" class="text-sm font-medium">Mot de passe
                <input type="password" v-model="formulaire.mot_de_passe" class="w-full border rounded px-3 py-2">
            </label>

            <p v-if="erreur" class="text-red-600 text-sm">{{ erreur }}</p>

            <div class="flex gap-2">
                <button type="submit" class="px-4 py-2 rounded-lg bg-brand-green text-white font-semibold">Enregistrer</button>
                <button type="button" @click="fermerFormulaire" class="px-4 py-2 rounded-lg bg-gray-200">Annuler</button>
            </div>
        </form>

    </div>

    <script src="js/utilisateurs.js"></script>
</body>

</html>
