<?php
require "inclus/espace.php";
$titrePage = "nav.tableau_de_bord";
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
                <h1>{{ $t('admin.tableau_de_bord.titre') }}</h1>
            </div>

            <form class="filtres-utilisateurs periode-finances" @submit.prevent="charger">
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
                    :class="{ active: raccourciActif === raccourci }" @click="choisirRaccourci(raccourci)">
                    {{ $t('admin.finances.raccourci.' + raccourci) }}
                </button>
            </form>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <template v-if="tableau">
                <p class="discret comparaison-periode">
                    {{ $t('admin.tableau_de_bord.comparaison', { debut: formaterDate(tableau.debut_precedent), fin: formaterDate(tableau.fin_precedent) }) }}
                </p>

                <h2 class="titre-bloc-reporting">{{ $t('admin.tableau_de_bord.bloc_finances') }}</h2>
                <div class="grille-finances grille-kpi">
                    <div v-for="kpi in kpiFinances" :key="kpi.cle" class="carte carte-kpi">
                        <span class="discret">{{ $t('admin.tableau_de_bord.kpi.' + kpi.cle) }}</span>
                        <p class="chiffre-finances" :class="{ negatif: kpi.valeur < 0 }">{{ formaterMontant(kpi.valeur) }}</p>
                        <span v-if="kpi.evolution" class="evolution" :class="kpi.evolution.classe" :title="kpi.evolution.description" :aria-label="kpi.evolution.description">
                            <span aria-hidden="true">{{ kpi.evolution.fleche }} {{ kpi.evolution.texte }}</span>
                        </span>
                        <span v-else class="discret">{{ kpi.detail }}</span>
                    </div>
                </div>

                <h2 class="titre-bloc-reporting">{{ $t('admin.tableau_de_bord.bloc_activite') }}</h2>
                <div class="grille-finances grille-kpi grille-kpi-activite">
                    <div v-for="kpi in kpiActivite" :key="kpi.cle" class="carte carte-kpi">
                        <span class="discret">{{ $t('admin.tableau_de_bord.kpi.' + kpi.cle) }}</span>
                        <p class="chiffre-finances">{{ formaterNombre(kpi.valeur) }}</p>
                        <span v-if="kpi.evolution" class="evolution" :class="kpi.evolution.classe" :title="kpi.evolution.description" :aria-label="kpi.evolution.description">
                            <span aria-hidden="true">{{ kpi.evolution.fleche }} {{ kpi.evolution.texte }}</span>
                        </span>
                    </div>
                </div>

                <div class="grille-tableau-de-bord">
                    <section class="carte graphique-principal">
                        <h3>{{ $t('admin.tableau_de_bord.graphique.finances.titre') }}</h3>
                        <p class="discret">{{ $t('admin.tableau_de_bord.graphique.series_detail', { debut: formaterMoisDate(tableau.debut_series), fin: formaterMoisDate(tableau.fin) }) }}</p>
                        <div class="zone-graphique">
                            <canvas id="graphique-finances" role="img" :aria-label="$t('admin.tableau_de_bord.graphique.finances.titre')"></canvas>
                        </div>
                    </section>

                    <section class="carte a-traiter">
                        <h3>{{ $t('admin.tableau_de_bord.a_traiter.titre') }}</h3>
                        <p v-if="totalATraiter === 0" class="discret">{{ $t('admin.tableau_de_bord.a_traiter.rien') }}</p>
                        <ul v-else>
                            <li v-for="element in elementsATraiter" :key="element.cle">
                                <a :href="element.lien" :class="{ vide: element.nombre === 0 }">
                                    <span>{{ $t('admin.tableau_de_bord.a_traiter.' + element.cle) }}</span>
                                    <strong class="compteur-a-traiter">{{ element.nombre }}</strong>
                                </a>
                            </li>
                        </ul>
                    </section>

                    <section class="carte">
                        <h3>{{ $t('admin.tableau_de_bord.graphique.activite.titre') }}</h3>
                        <p class="discret">{{ $t('admin.tableau_de_bord.graphique.series_detail', { debut: formaterMoisDate(tableau.debut_series), fin: formaterMoisDate(tableau.fin) }) }}</p>
                        <div class="zone-graphique">
                            <canvas id="graphique-activite" role="img" :aria-label="$t('admin.tableau_de_bord.graphique.activite.titre')"></canvas>
                        </div>
                    </section>

                    <section class="carte">
                        <h3>{{ $t('admin.tableau_de_bord.graphique.sources.titre') }}</h3>
                        <p class="discret">{{ $t('admin.tableau_de_bord.graphique.sources.detail') }}</p>
                        <p v-if="tableau.actuel.revenus === 0" class="discret">{{ $t('admin.tableau_de_bord.aucun_revenu') }}</p>
                        <div v-else class="zone-graphique">
                            <canvas id="graphique-sources" role="img" :aria-label="$t('admin.tableau_de_bord.graphique.sources.titre')"></canvas>
                        </div>
                    </section>
                </div>

                <p class="liens-detail">
                    <a href="finances.php" class="lien-action">{{ $t('admin.tableau_de_bord.voir_finances') }}</a>
                    <a href="reporting.php" class="lien-action">{{ $t('admin.tableau_de_bord.voir_reporting') }}</a>
                </p>
            </template>

        </main>
    </div>

    <script src="https://unpkg.com/chart.js@4.4.7/dist/chart.umd.js"></script>
    <script src="<?= ressource("js/periode.js") ?>"></script>
    <script src="<?= ressource("js/tableau-de-bord.js") ?>"></script>
</body>

</html>
