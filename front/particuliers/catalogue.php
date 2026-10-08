<?php
require "inclus/espace.php";
$titrePage = "Catalogue";
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
                <h1>Formations, ateliers et événements</h1>
            </div>

            <div class="barre-filtres">
                <div class="form-group">
                    <label for="type">Type</label>
                    <select id="type" v-model="typeFiltre" @change="chargerCatalogue">
                        <option value="">Tous</option>
                        <option v-for="(libelle, type) in libellesOffre" :key="type" :value="type">{{ libelle }}</option>
                    </select>
                </div>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="message" class="succes">{{ message }}</p>
            <p v-if="!erreur && offres.length === 0" class="discret">Aucune offre publiée pour le moment.</p>

            <div class="grille">
                <div v-for="o in offres" :key="o.type + o.id" class="carte">
                    <span class="badge">{{ libellesOffre[o.type] }}</span>
                    <h3>{{ o.titre }}</h3>
                    <p>{{ o.description }}</p>
                    <ul class="liste-simple">
                        <li>{{ formaterDateHeure(o.date_debut) }}</li>
                        <li v-if="o.lieu">{{ o.lieu }}</li>
                        <li>{{ formaterPrix(o.tarif) }}</li>
                        <li v-if="o.nb_places !== null">{{ o.nb_places }} places restantes</li>
                    </ul>
                    <p v-if="estInscrit(o)" class="succes">Vous êtes inscrit·e</p>
                    <button v-else class="button-submit" @click="inscrire(o)">S'inscrire</button>
                </div>
            </div>

        </main>
    </div>

    <script src="js/catalogue.js"></script>
</body>

</html>
