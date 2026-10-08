<?php
require "inclus/espace.php";
$titrePage = "Conseils";
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
                <h1>Conseils</h1>
            </div>

            <div class="barre-filtres">
                <div class="form-group">
                    <label for="recherche">Rechercher</label>
                    <input id="recherche" type="search" v-model="recherche" placeholder="Mot-clé">
                </div>
                <div class="form-group">
                    <label for="categorie">Catégorie</label>
                    <select id="categorie" v-model="categorie">
                        <option value="">Toutes</option>
                        <option v-for="c in categories" :key="c" :value="c">{{ c }}</option>
                    </select>
                </div>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="!erreur && conseilsFiltres.length === 0" class="discret">Aucun conseil trouvé.</p>

            <div class="deux-colonnes">
                <div class="grille">
                    <div v-for="c in conseilsFiltres" :key="c.id" class="carte carte-cliquable"
                        :class="{ 'carte-selectionnee': selection && selection.id === c.id }" @click="selection = c">
                        <span v-if="c.categorie" class="badge">{{ c.categorie }}</span>
                        <h3>{{ c.titre }}</h3>
                        <p class="discret">{{ formaterDate(c.created_at) }}</p>
                    </div>
                </div>

                <article v-if="selection" class="carte">
                    <span v-if="selection.categorie" class="badge">{{ selection.categorie }}</span>
                    <h2>{{ selection.titre }}</h2>
                    <p class="discret">Publié le {{ formaterDate(selection.created_at) }}</p>
                    <p>{{ selection.contenu }}</p>
                </article>
            </div>

        </main>
    </div>

    <script src="js/conseils.js"></script>
</body>

</html>
