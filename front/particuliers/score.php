<?php
require "inclus/espace.php";
$titrePage = "Upcycling Score";
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
                <h1>Mon Upcycling Score</h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="deux-colonnes">

                <section class="carte">
                    <h2>Historique</h2>
                    <p v-if="historique.length === 0" class="discret">Aucun point gagné pour le moment.</p>
                    <ul class="liste-simple">
                        <li v-for="h in historique" :key="h.id">
                            <strong :class="h.delta >= 0 ? 'succes' : 'erreur'">{{ h.delta >= 0 ? "+" : "" }}{{ h.delta }}</strong>
                            {{ h.motif }}
                            <span class="discret"> · {{ formaterDate(h.date) }}</span>
                        </li>
                    </ul>
                </section>

                <div>
                    <div class="carte section">
                        <span class="discret">Score actuel</span>
                        <p class="chiffre-cle">{{ score }} points</p>
                    </div>

                    <div class="carte">
                        <h3>Comment gagner des points ?</h3>
                        <ul class="liste-simple">
                            <li v-for="regle in bareme" :key="regle.action">
                                {{ regle.action }} : <strong>+{{ regle.points }}</strong>
                            </li>
                        </ul>
                    </div>
                </div>

            </div>

        </main>
    </div>

    <script src="js/score.js"></script>
</body>

</html>
