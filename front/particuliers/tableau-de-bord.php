<?php
require "inclus/espace.php";
$titrePage = "Accueil";
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
                <h1>Bonjour {{ prenom }}</h1>
                <div class="form-buttons">
                    <a class="btn-primary" href="mes-annonces.php">Déposer une annonce</a>
                    <a class="btn-primary" href="depots.php">Demander un dépôt en conteneur</a>
                </div>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div v-if="resume" class="grille">
                <a class="carte carte-cliquable" href="mes-annonces.php">
                    <span class="discret">Annonces en ligne</span>
                    <p class="chiffre-cle">{{ resume.annonces_actives }}</p>
                </a>
                <a class="carte carte-cliquable" href="depots.php">
                    <span class="discret">Dernier dépôt en conteneur</span>
                    <p class="chiffre-cle">{{ libelleDepot(resume.dernier_depot_statut) }}</p>
                </a>
                <a class="carte carte-cliquable" href="score.php">
                    <span class="discret">Upcycling Score</span>
                    <p class="chiffre-cle">{{ resume.upcycling_score }}</p>
                </a>
                <a class="carte carte-cliquable" href="planning.php">
                    <span class="discret">Prochain rendez-vous</span>
                    <p class="chiffre-cle" v-if="resume.prochaine_offre_titre">{{ resume.prochaine_offre_titre }}</p>
                    <p class="discret" v-if="resume.prochaine_offre_date">{{ formaterDate(resume.prochaine_offre_date) }}</p>
                    <p class="chiffre-cle" v-else>—</p>
                </a>
            </div>

        </main>

        <!-- Tutoriel de première connexion : s'affiche tant que tutoriel_vu est faux -->
        <div v-if="tutorielVisible" class="fond-modal">
            <div class="modal" role="dialog" aria-modal="true" aria-labelledby="titre-tutoriel">
                <h2 id="titre-tutoriel">{{ etapeCourante.titre }}</h2>
                <p>{{ etapeCourante.texte }}</p>
                <p class="points-etapes">Étape {{ indexEtape + 1 }} sur {{ etapesTutoriel.length }}</p>

                <div class="form-buttons">
                    <button v-if="indexEtape > 0" class="button-draft" @click="indexEtape--">Précédent</button>
                    <button v-if="!derniereEtape" class="button-submit" @click="indexEtape++">Suivant</button>
                    <button v-else class="button-submit" @click="terminerTutoriel">Terminer</button>
                </div>
            </div>
        </div>
    </div>

    <script src="js/tableau-de-bord.js"></script>
</body>

</html>
