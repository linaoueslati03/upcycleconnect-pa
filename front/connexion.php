<!DOCTYPE html>
<html lang="fr">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title data-titre="connexion.titre">UpcycleConnect</title>
    <link rel="stylesheet" href="assets/css/style.css">
    <link rel="stylesheet" href="assets/css/connexion.css">
    <script src="assets/js/config.js"></script>
    <script src="assets/js/api.js"></script>
    <script src="https://unpkg.com/vue@3.5.13/dist/vue.global.prod.js"></script>
</head>

<body class="page-connexion">

    <div id="app" class="carte-connexion">
        <img src="assets/img/LOGO-pa.png" alt="UpcycleConnect">
        <h1>{{ t('connexion.titre') }}</h1>

        <form @submit.prevent="seConnecter">
            <div class="form-group">
                <label for="email">{{ t('commun.email') }}</label>
                <input type="email" id="email" v-model="email" autocomplete="username" required>
            </div>

            <div class="form-group">
                <label for="mot-de-passe">{{ t('commun.mot_de_passe') }}</label>
                <input type="password" id="mot-de-passe" v-model="motDePasse" autocomplete="current-password" required>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <button type="submit" class="button-submit">{{ t('connexion.se_connecter') }}</button>
        </form>

        <p><a href="inscription.php">{{ t('connexion.creer_compte') }}</a></p>

        <div class="form-group">
            <label for="langue">{{ t('commun.langue') }}</label>
            <select id="langue" :value="langueCourante()" @change="changerLangue($event.target.value)">
                <option v-for="l in langues" :key="l.code" :value="l.code">{{ l.libelle }}</option>
            </select>
        </div>
    </div>

    <script src="assets/js/connexion.js"></script>
</body>

</html>
