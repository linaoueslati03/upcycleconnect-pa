<?php
$titrePage = "Mon planning";
?>
<!DOCTYPE html>
<html lang="fr">

<head>
    <?php include "inclus/entete.php"; ?>
</head>

<body>
    <?php include "inclus/navigation.php"; ?>

    <div id="app">
        <main>

            <div class="planning-header">
                <h1>Mon planning</h1>

                <div class="planning-filtres">
                    <div>
                        <label for="semaine">Semaine</label>
                        <select id="semaine" v-model="semaineSelectionnee">
                            <option v-for="semaine in semaines" :key="semaine.numero" :value="semaine.numero">
                                Semaine du {{ formaterJourMois(semaine.date) }}
                            </option>
                        </select>
                    </div>

                    <div>
                        <label for="annee">Année</label>
                        <select id="annee" v-model="anneeSelectionnee">
                            <option v-for="annee in annees" :key="annee" :value="annee">{{ annee }}</option>
                        </select>
                    </div>
                </div>
            </div>

            <div class="semaine-titre">
                <button @click="semainePrecedente">←</button>
                <h2>Semaine du {{ formaterJourMois(debutSemaine) }}</h2>
                <button @click="semaineSuivante">→</button>
            </div>

            <div class="planning">
                <div class="jour" v-for="jour in jours" :key="jour">
                    <h3>{{ jour }}</h3>
                    <div class="evenements"></div>
                </div>
            </div>

        </main>
    </div>

    <script src="js/planning.js"></script>
</body>

</html>
