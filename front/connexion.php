
<!DOCTYPE html>
<html lang="fr">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title data-titre="connexion.titre">UpcycleConnect</title>

    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>

    <link
        rel="stylesheet"
        href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Merriweather:wght@400;700&family=Poppins:wght@500;600;700&display=swap"
    >

    <link rel="stylesheet" href="assets/css/style.css">
    <link rel="stylesheet" href="assets/css/connexion.css">

    <script src="assets/js/config.js"></script>
    <script src="assets/js/api.js"></script>
    <script src="https://unpkg.com/vue@3.5.13/dist/vue.global.prod.js"></script>
</head>

<body class="page-connexion">

    <div id="app" class="carte-connexion">

        <div class="panneau-logo">

            <img
                src="assets/img/LOGO-pa.png"
                alt="Logo UpcycleConnect"
            >

            <h2>UpcycleConnect</h2>

        </div>


        <div class="panneau-formulaire">

            <div class="choix-langue">

                <label for="langue">
                    {{ t('commun.langue') }}
                </label>

                <select
                    id="langue"
                    :value="langueCourante()"
                    @change="changerLangue($event.target.value)"
                >

                    <option
                        v-for="l in langues"
                        :key="l.code"
                        :value="l.code"
                    >
                        {{ l.libelle }}
                    </option>

                </select>

            </div>


            <h1>{{ t('connexion.titre') }}</h1>


            <form @submit.prevent="seConnecter">

                <div class="form-group">

                    <input
                        type="email"
                        id="email"
                        v-model="email"
                        :placeholder="t('commun.email')"
                        autocomplete="username"
                        required
                    >

                </div>


                <div class="form-group">

                    <input
                        type="password"
                        id="mot-de-passe"
                        v-model="motDePasse"
                        :placeholder="t('commun.mot_de_passe')"
                        autocomplete="current-password"
                        required
                    >

                </div>


                <p
                    v-if="erreur"
                    class="erreur"
                >
                    {{ erreur }}
                </p>


                <button
                    type="submit"
                    class="button-submit"
                >
                    {{ t('connexion.se_connecter') }}
                </button>

            </form>


            <p class="mot-de-passe-oublie">

                <a href="mot-de-passe-oublie.php">
                    Vous avez oublié votre mot de passe ?
                </a>

            </p>

        </div>

    </div>


    <script src="assets/js/connexion.js"></script>

</body>

</html>