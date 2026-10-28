<?php
require "inclus/espace.php";
$titrePage = "salaries.formations.creer";
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
                <h1>
                    {{ modification
                        ? $t('salaries.formations.modifier')
                        : $t('salaries.formations.creer') }}
                </h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="succes" class="succes">{{ succes }}</p>

            <form @submit.prevent="enregistrer">
                <div class="form-group">
                    <label for="titre">{{ $t('commun.titre') }}</label>
                    <input id="titre" v-model.trim="formation.titre" required maxlength="150">
                </div>

                <div class="form-group">
                    <label for="description">{{ $t('commun.description') }}</label>
                    <textarea id="description" v-model="formation.description"></textarea>
                </div>

                <div class="form-group">
                    <label for="date_debut">{{ $t('commun.date_debut') }}</label>
                    <input id="date_debut" type="datetime-local" v-model="formation.date_debut" required>
                </div>

                <div class="form-group">
                    <label for="date_fin">{{ $t('commun.date_fin') }}</label>
                    <input id="date_fin" type="datetime-local" v-model="formation.date_fin">
                </div>

                <div class="form-group">
                    <label for="lieu">{{ $t('commun.lieu') }}</label>
                    <input id="lieu" v-model.trim="formation.lieu" maxlength="150">
                </div>

                <div class="form-group">
                    <label for="nb_places">{{ $t('salaries.formations.places') }}</label>
                    <input id="nb_places" type="number" min="1" v-model.number="formation.nb_places">
                </div>

                <div class="form-group">
                    <label for="tarif">{{ $t('commun.tarif') }} (€)</label>
                    <input id="tarif" type="number" min="0" step="0.01" v-model.number="formation.tarif" required>
                </div>

                <div class="form-group">
                    <label for="materiel_necessaire">
                        {{ $t('salaries.formations.materiel') }}
                    </label>
                    <textarea id="materiel_necessaire" v-model="formation.materiel_necessaire"></textarea>
                </div>

                <div class="form-group">
                    <label for="statut">{{ $t('commun.statut') }}</label>
                    <select id="statut" v-model="formation.statut">
                        <option value="brouillon">{{ formaterStatut('brouillon') }}</option>
                        <option value="en_attente">{{ formaterStatut('en_attente') }}</option>
                    </select>
                </div>

                <div class="boutons-formulaire">
                    <button type="submit" class="button-submit" :disabled="enregistrement">
                        {{ enregistrement
                            ? $t('commun.chargement')
                            : $t('commun.enregistrer') }}
                    </button>

                    <a href="formations.php" class="lien-action">
                        {{ $t('commun.annuler') }}
                    </a>
                </div>
            </form>
        </main>
    </div>

    <script src="js/creer-formation.js"></script>
</body>

</html>