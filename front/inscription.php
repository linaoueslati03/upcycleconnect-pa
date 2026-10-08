<!DOCTYPE html>
<html lang="fr">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title data-titre="connexion.creer_compte">UpcycleConnect</title>
    <!-- Polices de la charte graphique (Google Fonts) -->
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Merriweather:wght@400;700&family=Poppins:wght@500;600;700&display=swap">
    <link rel="stylesheet" href="assets/css/style.css">
    <link rel="stylesheet" href="assets/css/connexion.css">
    <script src="assets/js/config.js"></script>
    <script src="assets/js/api.js"></script>
    <script src="https://unpkg.com/vue@3.5.13/dist/vue.global.prod.js"></script>
</head>

<body class="page-connexion">

    <div id="app" class="carte-connexion">
        <img src="assets/img/LOGO-pa.png" alt="UpcycleConnect">
        <h1>{{ t('connexion.creer_compte') }}</h1>

        <form @submit.prevent="creerCompte">
            <div class="form-group">
                <label for="role">{{ t('inscription.je_suis') }}</label>
                <select id="role" v-model="compte.role">
                    <option value="particulier">{{ t('inscription.particulier') }}</option>
                    <option value="professionnel">{{ t('inscription.professionnel') }}</option>
                </select>
            </div>
            <div class="form-group">
                <label for="nom">{{ t('commun.nom') }}</label>
                <input id="nom" type="text" v-model="compte.nom" autocomplete="family-name" required>
            </div>
            <div class="form-group">
                <label for="prenom">{{ t('commun.prenom') }}</label>
                <input id="prenom" type="text" v-model="compte.prenom" autocomplete="given-name" required>
            </div>
            <div class="form-group">
                <label for="email">{{ t('commun.email') }}</label>
                <input id="email" type="email" v-model="compte.email" autocomplete="email" required>
            </div>
            <div class="form-group">
                <label for="mot-de-passe">{{ t('commun.mot_de_passe_min') }}</label>
                <input id="mot-de-passe" type="password" v-model="compte.mot_de_passe" autocomplete="new-password"
                    minlength="8" required>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <button type="submit" class="button-submit">{{ t('inscription.creer') }}</button>
        </form>

        <p><a href="connexion.php">{{ t('inscription.deja_compte') }}</a></p>
    </div>

    <script src="assets/js/inscription.js"></script>
</body>

</html>
