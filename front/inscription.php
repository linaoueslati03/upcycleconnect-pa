<!DOCTYPE html>
<html lang="fr">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Créer un compte — UpcycleConnect</title>
    <link rel="stylesheet" href="assets/css/style.css">
    <link rel="stylesheet" href="assets/css/connexion.css">
    <script src="assets/js/config.js"></script>
    <script src="assets/js/api.js"></script>
    <script src="https://unpkg.com/vue@3.5.13/dist/vue.global.prod.js"></script>
</head>

<body class="page-connexion">

    <div id="app" class="carte-connexion">
        <img src="assets/img/LOGO-pa.png" alt="UpcycleConnect">
        <h1>Créer un compte</h1>

        <form @submit.prevent="creerCompte">
            <div class="form-group">
                <label for="role">Je suis</label>
                <select id="role" v-model="compte.role">
                    <option value="particulier">Un particulier</option>
                    <option value="professionnel">Un professionnel / artisan</option>
                </select>
            </div>
            <div class="form-group">
                <label for="nom">Nom</label>
                <input id="nom" type="text" v-model="compte.nom" autocomplete="family-name" required>
            </div>
            <div class="form-group">
                <label for="prenom">Prénom</label>
                <input id="prenom" type="text" v-model="compte.prenom" autocomplete="given-name" required>
            </div>
            <div class="form-group">
                <label for="email">Email</label>
                <input id="email" type="email" v-model="compte.email" autocomplete="email" required>
            </div>
            <div class="form-group">
                <label for="mot-de-passe">Mot de passe (8 caractères minimum)</label>
                <input id="mot-de-passe" type="password" v-model="compte.mot_de_passe" autocomplete="new-password"
                    minlength="8" required>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <button type="submit" class="button-submit">Créer mon compte</button>
        </form>

        <p><a href="connexion.php">J'ai déjà un compte</a></p>
    </div>

    <script src="assets/js/inscription.js"></script>
</body>

</html>
