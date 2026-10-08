<?php
require "inclus/espace.php";
$titrePage = "Utilisateurs";
?>
<!DOCTYPE html>
<html lang="fr">

<head>
    <?php include "../inclus/entete.php"; ?>
</head>

<body>
    <?php include "../inclus/navigation.php"; ?>

    <div id="app">
        <main>

            <div class="page-header">
                <h1>Gestion des utilisateurs</h1>
                <button class="btn-primary" @click="reinitialiserFormulaire">+ Nouvel utilisateur</button>
            </div>

            <p v-if="erreurListe" class="erreur">{{ erreurListe }}</p>

            <div class="mise-en-page-admin">

                <section class="carte">
                    <table class="tableau">
                        <thead>
                            <tr>
                                <th>Nom</th>
                                <th>Prénom</th>
                                <th>Email</th>
                                <th>Rôle</th>
                                <th>Statut</th>
                                <th>Score</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="u in utilisateurs" :key="u.id">
                                <td>{{ u.nom }}</td>
                                <td>{{ u.prenom }}</td>
                                <td>{{ u.email }}</td>
                                <td>{{ nomsRoles[u.role_id] ?? u.role_id }}</td>
                                <td>{{ u.statut }}</td>
                                <td>{{ u.upcycling_score }}</td>
                                <td>
                                    <button class="lien-action" @click="modifier(u)">Modifier</button>
                                    <button class="lien-action lien-supprimer" @click="supprimer(u.id)">Supprimer</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </section>

                <div class="carte carte-formulaire">
                    <h2>{{ formulaire.id ? "Modifier l'utilisateur" : "Nouvel utilisateur" }}</h2>

                    <form @submit.prevent="enregistrer">
                        <div class="form-group">
                            <label for="nom">Nom</label>
                            <input id="nom" type="text" v-model="formulaire.nom" required>
                        </div>

                        <div class="form-group">
                            <label for="prenom">Prénom</label>
                            <input id="prenom" type="text" v-model="formulaire.prenom" required>
                        </div>

                        <div class="form-group">
                            <label for="email">Email</label>
                            <input id="email" type="email" v-model="formulaire.email" required>
                        </div>

                        <div class="form-group">
                            <label for="role">Rôle</label>
                            <select id="role" v-model="formulaire.role_id">
                                <option v-for="(nom, id) in nomsRoles" :key="id" :value="Number(id)">{{ nom }}</option>
                            </select>
                        </div>

                        <div v-if="!formulaire.id" class="form-group">
                            <label for="mot-de-passe">Mot de passe</label>
                            <input id="mot-de-passe" type="password" v-model="formulaire.mot_de_passe" required>
                        </div>

                        <p v-if="erreur" class="erreur">{{ erreur }}</p>

                        <div class="form-buttons">
                            <button type="submit" class="button-submit">Enregistrer</button>
                            <button type="button" class="button-draft" @click="reinitialiserFormulaire">Annuler</button>
                        </div>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/utilisateurs.js"></script>
</body>

</html>
