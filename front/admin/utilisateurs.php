<?php
require "inclus/espace.php";
$titrePage = "nav.utilisateurs";
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
                <h1>{{ t('admin.utilisateurs.titre') }}</h1>
                <button class="btn-primary" @click="reinitialiserFormulaire">+ {{ t('admin.utilisateurs.nouveau') }}</button>
            </div>

            <p v-if="erreurListe" class="erreur">{{ erreurListe }}</p>

            <div class="mise-en-page-admin">

                <section class="carte">
                    <table class="tableau">
                        <thead>
                            <tr>
                                <th>{{ t('commun.nom') }}</th>
                                <th>{{ t('commun.prenom') }}</th>
                                <th>{{ t('commun.email') }}</th>
                                <th>{{ t('commun.role') }}</th>
                                <th>{{ t('commun.statut') }}</th>
                                <th>{{ t('commun.score') }}</th>
                                <th>{{ t('commun.actions') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="u in utilisateurs" :key="u.id">
                                <td>{{ u.nom }}</td>
                                <td>{{ u.prenom }}</td>
                                <td>{{ u.email }}</td>
                                <td>{{ t('role.' + u.role_id) }}</td>
                                <td>{{ u.statut }}</td>
                                <td>{{ u.upcycling_score }}</td>
                                <td>
                                    <button class="lien-action" @click="modifier(u)">{{ t('commun.modifier') }}</button>
                                    <button class="lien-action lien-supprimer" @click="supprimer(u.id)">{{ t('commun.supprimer') }}</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </section>

                <div class="carte carte-formulaire">
                    <h2>{{ formulaire.id ? t('admin.utilisateurs.modifier') : t('admin.utilisateurs.nouveau') }}</h2>

                    <form @submit.prevent="enregistrer">
                        <div class="form-group">
                            <label for="nom">{{ t('commun.nom') }}</label>
                            <input id="nom" type="text" v-model="formulaire.nom" required>
                        </div>

                        <div class="form-group">
                            <label for="prenom">{{ t('commun.prenom') }}</label>
                            <input id="prenom" type="text" v-model="formulaire.prenom" required>
                        </div>

                        <div class="form-group">
                            <label for="email">{{ t('commun.email') }}</label>
                            <input id="email" type="email" v-model="formulaire.email" required>
                        </div>

                        <div class="form-group">
                            <label for="role">{{ t('commun.role') }}</label>
                            <select id="role" v-model="formulaire.role_id">
                                <option v-for="id in [1, 2, 3, 4]" :key="id" :value="id">{{ t('role.' + id) }}</option>
                            </select>
                        </div>

                        <div v-if="!formulaire.id" class="form-group">
                            <label for="mot-de-passe">{{ t('commun.mot_de_passe') }}</label>
                            <input id="mot-de-passe" type="password" v-model="formulaire.mot_de_passe" required>
                        </div>

                        <p v-if="erreur" class="erreur">{{ erreur }}</p>

                        <div class="form-buttons">
                            <button type="submit" class="button-submit">{{ t('commun.enregistrer') }}</button>
                            <button type="button" class="button-draft" @click="reinitialiserFormulaire">{{ t('commun.annuler') }}</button>
                        </div>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/utilisateurs.js"></script>
</body>

</html>
