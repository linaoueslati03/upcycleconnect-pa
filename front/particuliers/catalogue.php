<?php
require "inclus/espace.php";
$titrePage = "nav.catalogue";
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
                <h1>{{ t('catalogue.titre') }}</h1>
            </div>

            <div class="barre-filtres">
                <div class="form-group">
                    <label for="type">{{ t('commun.type') }}</label>
                    <select id="type" v-model="typeFiltre" @change="chargerCatalogue">
                        <option value="">{{ t('commun.tous') }}</option>
                        <option v-for="type in typesOffre" :key="type" :value="type">{{ t('offre.' + type) }}</option>
                    </select>
                </div>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="message" class="succes">{{ message }}</p>
            <p v-if="!erreur && offres.length === 0" class="discret">{{ t('catalogue.aucune') }}</p>

            <div class="grille">
                <div v-for="o in offres" :key="o.type + o.id" class="carte">
                    <span class="badge">{{ t('offre.' + o.type) }}</span>
                    <h3>{{ o.titre }}</h3>
                    <p>{{ o.description }}</p>
                    <ul class="liste-simple">
                        <li>{{ formaterDateHeure(o.date_debut) }}</li>
                        <li v-if="o.lieu">{{ o.lieu }}</li>
                        <li>{{ formaterPrix(o.tarif) }}</li>
                        <li v-if="o.nb_places !== null">{{ t('catalogue.places_restantes', { nombre: o.nb_places }) }}</li>
                    </ul>
                    <p v-if="estInscrit(o)" class="succes">{{ t('catalogue.inscrit') }}</p>
                    <button v-else class="button-submit" @click="inscrire(o)">{{ t('catalogue.sinscrire') }}</button>
                </div>
            </div>

        </main>
    </div>

    <script src="js/catalogue.js"></script>
</body>

</html>
