<?php
$titrePage = "Articles";
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

            <div class="page-header">
                <h1>Articles</h1>
                <button class="btn-primary" @click="formulaireVisible = true">Créer un article</button>
            </div>

            <div v-if="formulaireVisible" class="form-article">

                <div class="form-group">
                    <label for="titre">Titre</label>
                    <input id="titre" type="text" placeholder="Titre de l'article" v-model="titre">
                </div>

                <div class="form-group">
                    <label for="categorie">Catégorie</label>
                    <select id="categorie" v-model="categorie">
                        <option value="">Sélectionner une catégorie</option>
                        <option value="actualites">Actualités</option>
                        <option value="conseils">Conseils</option>
                        <option value="evenements">Événements</option>
                        <option value="upcycling">Upcycling</option>
                    </select>
                </div>

                <div class="form-group">
                    <label for="texte">Article</label>
                    <textarea id="texte" placeholder="Écrivez votre article..." v-model="texte"></textarea>
                </div>

                <div>
                    <button type="button" class="button-submit" @click="publier">Publier</button>
                </div>

            </div>

        </main>
    </div>

    <script src="js/articles.js"></script>
</body>

</html>
