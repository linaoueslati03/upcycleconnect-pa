<?php
require "inclus/espace.php";
$titrePage = "nav.mon_planning";
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
                <h1>{{ t('nav.mon_planning') }}</h1>
                <a class="btn-primary" href="catalogue.php">{{ t('planning.voir_catalogue') }}</a>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="!erreur && entrees.length === 0" class="discret">
                {{ t('planning.aucune_inscription') }}
            </p>

            <div v-for="e in entrees" :key="e.inscription_id" class="carte section">
                <span class="badge">{{ t('offre.' + e.type) }}</span>
                <h3>{{ e.titre }}</h3>
                <p class="discret">
                    {{ formaterDateHeure(e.date_debut) }}<span v-if="e.lieu"> · {{ e.lieu }}</span>
                </p>
                <p v-if="e.description">{{ e.description }}</p>
            </div>

        </main>
    </div>

    <script src="js/planning.js"></script>
</body>

</html>
