<?php
require "inclus/espace.php";
$titrePage = "Annonces";
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
                <h1>Annonces</h1>
                <a class="btn-primary" href="mes-annonces.php">Déposer une annonce</a>
            </div>

            <form class="barre-filtres" @submit.prevent="chargerAnnonces">
                <div class="form-group">
                    <label for="type">Type</label>
                    <select id="type" v-model="filtres.type">
                        <option value="">Tous</option>
                        <option value="don">Don</option>
                        <option value="vente">Vente</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="localisation">Localisation</label>
                    <input id="localisation" type="text" v-model="filtres.localisation" placeholder="Ex. Paris">
                </div>
                <button type="submit" class="button-submit">Rechercher</button>
            </form>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="!erreur && annonces.length === 0" class="discret">Aucune annonce ne correspond à votre recherche.</p>

            <div class="deux-colonnes">
                <div class="grille">
                    <div v-for="a in annonces" :key="a.id" class="carte carte-cliquable"
                        :class="{ 'carte-selectionnee': selection && selection.id === a.id }" @click="selection = a">
                        <span class="badge">{{ a.type === "don" ? "Don" : "Vente" }}</span>
                        <h3>{{ a.titre }}</h3>
                        <p class="discret">{{ a.localisation }} · {{ formaterDate(a.date_creation) }}</p>
                        <p v-if="a.type === 'vente'"><strong>{{ formaterPrix(a.prix) }}</strong></p>
                    </div>
                </div>

                <div v-if="selection" class="carte">
                    <span class="badge">{{ selection.type === "don" ? "Don" : "Vente" }}</span>
                    <h2>{{ selection.titre }}</h2>
                    <p>{{ selection.description || "Pas de description." }}</p>
                    <ul class="liste-simple">
                        <li><strong>Catégorie :</strong> {{ nomCategorie(selection.categorie_id) }}</li>
                        <li><strong>Localisation :</strong> {{ selection.localisation || "Non précisée" }}</li>
                        <li v-if="selection.type === 'vente'"><strong>Prix :</strong> {{ formaterPrix(selection.prix) }}</li>
                        <li><strong>Publiée le :</strong> {{ formaterDate(selection.date_creation) }}</li>
                    </ul>
                </div>
            </div>

        </main>
    </div>

    <script src="js/annonces.js"></script>
</body>

</html>
