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

            <div class="planning-header">
                <h1>{{ $t('nav.mon_planning') }}</h1>
                <a v-if="roleConnecte() === 'particulier'" class="btn-primary" href="catalogue.php">{{ $t('planning.voir_catalogue') }}</a>

                <div class="planning-filtres">
                    <div>
                        <label for="semaine">{{ $t('salaries.planning.semaine') }}</label>
                        <select id="semaine" v-model="semaineSelectionnee">
                            <option v-for="semaine in semaines" :key="semaine.numero" :value="semaine.numero">
                                {{ $t('salaries.planning.semaine_du', { date: formaterJourMois(semaine.date) }) }}
                            </option>
                        </select>
                    </div>

                    <div>
                        <label for="annee">{{ $t('salaries.planning.annee') }}</label>
                        <select id="annee" v-model="anneeSelectionnee">
                            <option v-for="annee in annees" :key="annee" :value="annee">{{ annee }}</option>
                        </select>
                    </div>
                </div>
            </div>

            <div class="semaine-titre">
                <button @click="semainePrecedente">←</button>
                <h2>{{ $t('salaries.planning.semaine_du', { date: formaterJourMois(debutSemaine) }) }}</h2>
                <button @click="semaineSuivante">→</button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="!erreur && roleConnecte() === 'particulier' && entrees.length === 0" class="discret">
                {{ $t('planning.aucune_inscription') }}
            </p>

            <div class="planning">
                <div class="jour" v-for="(jour, index) in jours" :key="jour">
                    <h3>{{ $t(jour) }}</h3>
                    <div class="evenements">
                        <template v-for="e in entreesDuJour(index)" :key="e.type + (e.id || e.inscription_id)">
                            <a v-if="roleConnecte() === 'salarie'" class="evenement"
                                :href="(e.type === 'atelier' ? 'creer-atelier.php' : 'creer-evenement.php') + '?id=' + e.id">
                                <span class="evenement-titre">{{ e.titre }}</span>
                                <span class="evenement-heure">{{ formaterHeure(e.date_debut) }} · {{ formaterStatut(e.statut) }}</span>
                            </a>
                            <div v-else class="evenement">
                                <span class="evenement-titre">{{ e.titre }}</span>
                                <span class="evenement-heure">{{ formaterHeure(e.date_debut) }}<template v-if="e.lieu"> · {{ e.lieu }}</template></span>
                            </div>
                        </template>
                    </div>
                </div>
            </div>

        </main>
    </div>

    <script src="js/planning.js"></script>
</body>

</html>
