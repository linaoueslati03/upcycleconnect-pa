<?php
require "inclus/espace.php";
$titrePage = "Créer un événement";
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

            <h1>{{ id ? "Modifier" : "Créer" }} un événement</h1>

            <form @submit.prevent>

                <div class="form-group">
                    <label for="titre">Titre</label>
                    <input type="text" id="titre" v-model="formulaire.titre" required>
                </div>

                <div class="form-group">
                    <label for="description">Description</label>
                    <textarea id="description" v-model="formulaire.description"></textarea>
                </div>

                <div class="form-group">
                    <label for="date-debut">Date de début</label>
                    <input type="datetime-local" id="date-debut" v-model="formulaire.date_debut" required>
                </div>

                <div class="form-group">
                    <label for="date-fin">Date de fin</label>
                    <input type="datetime-local" id="date-fin" v-model="formulaire.date_fin">
                </div>

                <div class="form-group">
                    <label for="lieu">Lieu</label>
                    <input type="text" id="lieu" v-model="formulaire.lieu">
                </div>

                <div class="form-group">
                    <label for="site">Site</label>
                    <select id="site" v-model="formulaire.site">
                        <option value="">Sélectionner un site</option>
                        <option value="paris-10">Paris 10e</option>
                        <option value="paris-11">Paris 11e</option>
                        <option value="paris-13">Paris 13e</option>
                        <option value="montreuil">Montreuil</option>
                        <option value="suisse">Suisse</option>
                    </select>
                </div>

                <p v-if="erreur" class="erreur">{{ erreur }}</p>

                <div class="form-buttons">
                    <button type="button" class="button-draft" @click="envoyerEvenement('brouillon')">
                        Enregistrer en brouillon
                    </button>
                    <button type="button" class="button-submit" @click="envoyerEvenement('en_attente')">
                        Soumettre à validation
                    </button>
                </div>

            </form>

            <div v-if="confirmation" id="confirmation">
                Événement enregistré ! <a href="evenements.php">Retour à mes événements</a>
            </div>

        </main>
    </div>

    <script src="js/creer-evenement.js"></script>
</body>

</html>
