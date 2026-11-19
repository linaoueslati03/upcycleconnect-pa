<?php require_once __DIR__ . "/inclus/ressource.php"; ?>
<!DOCTYPE html>
<html lang="fr">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title data-titre="connexion.creer_compte">UpcycleConnect</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Merriweather:wght@400;700&family=Poppins:wght@500;600;700&display=swap">
    <link rel="stylesheet" href="<?= ressource("assets/css/style.css") ?>">
    <link rel="stylesheet" href="<?= ressource("assets/css/connexion.css") ?>">
    <script src="<?= ressource("assets/js/config.js") ?>"></script>
    <script src="https://unpkg.com/vue@3.5.13/dist/vue.global.prod.js"></script>
    <script src="https://unpkg.com/vue-i18n@11.4.12/dist/vue-i18n.global.prod.js"></script>
    <script src="<?= ressource("assets/js/api.js") ?>"></script>
</head>

<body class="page-connexion">

    <div id="app" class="carte-connexion">

        <div class="panneau-logo">
            <img src="assets/img/LOGO-pa.png" alt="UpcycleConnect">
            <h2>UpcycleConnect</h2>
        </div>

        <div class="panneau-formulaire panneau-inscription">
            <h1>{{ $t('connexion.creer_compte') }}</h1>

            <form @submit.prevent="creerCompte">
                <div class="form-group">
                    <label for="role">{{ $t('inscription.je_suis') }}</label>
                    <select id="role" v-model="compte.role">
                        <option value="particulier">{{ $t('inscription.particulier') }}</option>
                        <option value="professionnel">{{ $t('inscription.professionnel') }}</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="nom">{{ $t('commun.nom') }}</label>
                    <input id="nom" type="text" v-model="compte.nom" autocomplete="family-name" required>
                </div>
                <div class="form-group">
                    <label for="prenom">{{ $t('commun.prenom') }}</label>
                    <input id="prenom" type="text" v-model="compte.prenom" autocomplete="given-name" required>
                </div>
                <div class="form-group">
                    <label for="email">{{ $t('commun.email') }}</label>
                    <input id="email" type="email" v-model="compte.email" autocomplete="email" required>
                </div>
                <div class="form-group">
                    <label for="mot-de-passe">{{ $t('commun.mot_de_passe_min') }}</label>
                    <input id="mot-de-passe" type="password" v-model="compte.mot_de_passe" autocomplete="new-password"
                        minlength="12" maxlength="72" required aria-describedby="regle-mot-de-passe">
                    <small id="regle-mot-de-passe">{{ $t('commun.regle_mot_de_passe') }}</small>
                </div>
                <div class="form-group">
                    <label for="confirmation-mot-de-passe">{{ $t('inscription.confirmer_mot_de_passe') }}</label>
                    <input id="confirmation-mot-de-passe" type="password" v-model="confirmation" autocomplete="new-password"
                        minlength="12" maxlength="72" required>
                </div>

                <p v-if="erreur" class="erreur" role="alert">{{ erreur }}</p>

                <button type="submit" class="button-submit" :disabled="envoi">{{ $t('inscription.creer') }}</button>
            </form>

            <p class="lien-secondaire"><a href="connexion">{{ $t('inscription.deja_compte') }}</a></p>
        </div>

    </div>

    <script src="<?= ressource("assets/js/inscription.js") ?>"></script>
</body>

</html>
