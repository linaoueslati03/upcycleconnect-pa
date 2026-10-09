<?php
require "inclus/espace.php";
$titrePage = "nav.accueil";
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
                <h1>{{ t('nav.accueil') }} - Espace Professionnel</h1>
            </div>
            <div class="mise-en-page-admin">
                <section class="carte">
                    <h2>Résumé de votre activité</h2>
                    <p>Bienvenue sur votre espace.</p>
                </section>
            </div>
        </main>
    </div>
    <script src="js/tableau-de-bord.js"></script>
</body>
</html>
