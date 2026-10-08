<?php
$titrePage = "Créer un atelier";
?>
<!DOCTYPE html>
<html lang="fr">

<head>
    <?php include "inclus/entete.php"; ?>
</head>

<body>
    <?php include "inclus/navigation.php"; ?>

    <div id="app">
        <main id="main-content">

            <h1>Créer un atelier</h1>

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

                <p v-if="erreur" class="erreur">{{ erreur }}</p>

                <div class="form-buttons">
                    <button type="button" class="button-draft" @click="envoyerAtelier('brouillon')">
                        Enregistrer en brouillon
                    </button>
                    <button type="button" class="button-submit" @click="envoyerAtelier('en_attente')">
                        Soumettre à validation
                    </button>
                </div>

            </form>

            <div v-if="confirmation" id="confirmation">
                Atelier enregistré ! <a href="ateliers.php">Retour à mes ateliers</a>
            </div>

        </main>
    </div>

    <script src="js/creer-atelier.js"></script>
</body>

</html>
