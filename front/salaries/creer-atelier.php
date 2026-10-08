<?php
require "inclus/espace.php";
$titrePage = "salaries.ateliers.creer";
?>
<!DOCTYPE html>
<html lang="fr">

<head>
    <?php include "../inclus/entete.php"; ?>
</head>

<body>
    <?php include "../inclus/navigation.php"; ?>

    <div id="app">
        <main id="main-content">

            <h1>{{ id ? t('salaries.ateliers.modifier') : t('salaries.ateliers.creer') }}</h1>

            <form @submit.prevent>

                <div class="form-group">
                    <label for="titre">{{ t('commun.titre') }}</label>
                    <input type="text" id="titre" v-model="formulaire.titre" required>
                </div>

                <div class="form-group">
                    <label for="description">{{ t('commun.description') }}</label>
                    <textarea id="description" v-model="formulaire.description"></textarea>
                </div>

                <div class="form-group">
                    <label for="date-debut">{{ t('commun.date_debut') }}</label>
                    <input type="datetime-local" id="date-debut" v-model="formulaire.date_debut" required>
                </div>

                <div class="form-group">
                    <label for="date-fin">{{ t('commun.date_fin') }}</label>
                    <input type="datetime-local" id="date-fin" v-model="formulaire.date_fin">
                </div>

                <div class="form-group">
                    <label for="lieu">{{ t('commun.lieu') }}</label>
                    <input type="text" id="lieu" v-model="formulaire.lieu">
                </div>

                <p v-if="erreur" class="erreur">{{ erreur }}</p>

                <div class="form-buttons">
                    <button type="button" class="button-draft" @click="envoyerAtelier('brouillon')">
                        {{ t('commun.enregistrer_brouillon') }}
                    </button>
                    <button type="button" class="button-submit" @click="envoyerAtelier('en_attente')">
                        {{ t('commun.soumettre_validation') }}
                    </button>
                </div>

            </form>

            <div v-if="confirmation" id="confirmation">
                {{ t('salaries.ateliers.enregistre') }} <a href="ateliers.php">{{ t('salaries.ateliers.retour') }}</a>
            </div>

        </main>
    </div>

    <script src="js/creer-atelier.js"></script>
</body>

</html>
