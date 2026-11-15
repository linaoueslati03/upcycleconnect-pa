<?php
require "inclus/espace.php";
$titrePage = "nav.reporting";
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
                <h1>{{ $t('admin.reporting.titre') }}</h1>
            </div>

            <form class="filtres-utilisateurs periode-finances" @submit.prevent="appliquerPeriode">
                <div class="form-group">
                    <label for="debut">{{ $t('commun.date_debut') }}</label>
                    <input id="debut" type="date" v-model="periode.debut" required>
                </div>
                <div class="form-group">
                    <label for="fin">{{ $t('commun.date_fin') }}</label>
                    <input id="fin" type="date" v-model="periode.fin" required>
                </div>
                <button type="submit" class="button-submit">{{ $t('admin.finances.appliquer') }}</button>
                <button v-for="raccourci in raccourcis" :key="raccourci" type="button" class="pill"
                    @click="choisirRaccourci(raccourci)">
                    {{ $t('admin.finances.raccourci.' + raccourci) }}
                </button>
            </form>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="filtres">
                <button v-for="o in onglets" :key="o" class="pill" :class="{ active: onglet === o }" @click="changerOnglet(o)">
                    {{ $t('admin.reporting.onglet.' + o) }}
                </button>
            </div>

            <template v-if="onglet === 'acteurs' && acteurs">
                <div class="grille-finances grille-totaux-reporting">
                    <div v-for="cle in totaux" :key="cle" class="carte">
                        <span class="discret">{{ $t('admin.reporting.totaux.' + cle) }}</span>
                        <p class="chiffre-finances">{{ formaterNombre(acteurs.totaux[cle]) }}</p>
                    </div>
                </div>

                <h2 class="titre-bloc-reporting">{{ $t('admin.reporting.etat_actuel') }}</h2>
                <p class="discret">{{ $t('admin.reporting.etat_actuel_detail') }}</p>
                <div class="grille-graphiques">
                    <section v-for="g in graphiquesEtat" :key="g" class="carte">
                        <h3>{{ $t('admin.reporting.graphique.' + g + '.titre') }}</h3>
                        <p class="discret">{{ $t('admin.reporting.graphique.' + g + '.detail') }}</p>
                        <div class="zone-graphique">
                            <canvas :id="'graphique-' + g" role="img" :aria-label="$t('admin.reporting.graphique.' + g + '.titre')"></canvas>
                        </div>
                    </section>
                </div>

                <h2 class="titre-bloc-reporting">{{ $t('admin.reporting.sur_periode') }}</h2>
                <p class="discret">{{ $t('admin.reporting.sur_periode_detail', { debut: formaterDate(acteurs.debut), fin: formaterDate(acteurs.fin) }) }}</p>
                <div class="grille-graphiques grille-graphiques-larges">
                    <section v-for="g in graphiquesPeriode" :key="g" class="carte">
                        <h3>{{ $t('admin.reporting.graphique.' + g + '.titre') }}</h3>
                        <p class="discret">{{ $t('admin.reporting.graphique.' + g + '.detail') }}</p>
                        <div class="zone-graphique">
                            <canvas :id="'graphique-' + g" role="img" :aria-label="$t('admin.reporting.graphique.' + g + '.titre')"></canvas>
                        </div>
                    </section>
                </div>
            </template>

        </main>
    </div>

    <script src="https://unpkg.com/chart.js@4.4.7/dist/chart.umd.js"></script>
    <script src="js/periode.js"></script>
    <script src="js/reporting.js"></script>
</body>

</html>
