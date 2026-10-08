<?php
require "inclus/espace.php";
$titrePage = "nav.score";
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
                <h1>{{ t('score.titre') }}</h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="deux-colonnes">

                <section class="carte">
                    <h2>{{ t('score.historique') }}</h2>
                    <p v-if="historique.length === 0" class="discret">{{ t('score.aucun_point') }}</p>
                    <ul class="liste-simple">
                        <li v-for="h in historique" :key="h.id">
                            <strong :class="h.delta >= 0 ? 'succes' : 'erreur'">{{ h.delta >= 0 ? "+" : "" }}{{ h.delta }}</strong>
                            {{ traduireOu('score.motif.' + h.motif, h.motif) }}
                            <span class="discret"> · {{ formaterDate(h.date) }}</span>
                        </li>
                    </ul>
                </section>

                <div>
                    <div class="carte section">
                        <span class="discret">{{ t('score.actuel') }}</span>
                        <p class="chiffre-cle">{{ t('score.points', { nombre: score }) }}</p>
                    </div>

                    <div class="carte">
                        <h3>{{ t('score.comment_gagner') }}</h3>
                        <ul class="liste-simple">
                            <li v-for="regle in bareme" :key="regle.cle">
                                {{ t(regle.cle) }} : <strong>+{{ regle.points }}</strong>
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
