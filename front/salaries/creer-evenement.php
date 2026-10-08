<?php
require "inclus/espace.php";
$titrePage = "salaries.evenements.creer";
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

            <h1>{{ id ? t('salaries.evenements.modifier') : t('salaries.evenements.creer') }}</h1>

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

                <div class="form-group">
                    <label for="site">{{ t('commun.site') }}</label>
                    <select id="site" v-model="formulaire.site">
                        <option value="">{{ t('salaries.evenements.choisir_site') }}</option>
                        <option value="paris-10">Paris 10e</option>
                        <option value="paris-11">Paris 11e</option>
                        <option value="paris-13">Paris 13e</option>
                        <option value="montreuil">Montreuil</option>
                        <option value="suisse">{{ t('site.suisse') }}</option>
                    </select>
                </div>

                <p v-if="erreur" class="erreur">{{ erreur }}</p>

                <div class="form-buttons">
                    <button type="button" class="button-draft" @click="envoyerEvenement('brouillon')">
                        {{ t('commun.enregistrer_brouillon') }}
                    </button>
                    <button type="button" class="button-submit" @click="envoyerEvenement('en_attente')">
                        {{ t('commun.soumettre_validation') }}
                    </button>
                </div>

            </form>

            <div v-if="confirmation" id="confirmation">
                {{ t('salaries.evenements.enregistre') }} <a href="evenements.php">{{ t('salaries.evenements.retour') }}</a>
            </div>

        </main>
    </div>

    <script src="js/creer-evenement.js"></script>
</body>

</html>
