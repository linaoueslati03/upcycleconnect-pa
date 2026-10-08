<?php
require "inclus/espace.php";
$titrePage = "Articles";
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
                <h1>Articles de conseils</h1>
                <button class="btn-primary" @click="nouvelArticle">Créer un article</button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div v-if="formulaireVisible" class="form-article">
                <h2>{{ formulaire.id ? "Modifier l'article" : "Nouvel article" }}</h2>

                <div class="form-group">
                    <label for="titre">Titre</label>
                    <input id="titre" type="text" placeholder="Titre de l'article" v-model="formulaire.titre">
                </div>

                <div class="form-group">
                    <label for="categorie">Catégorie</label>
                    <select id="categorie" v-model="formulaire.categorie">
                        <option value="">Sélectionner une catégorie</option>
                        <option value="actualites">Actualités</option>
                        <option value="conseils">Conseils</option>
                        <option value="evenements">Événements</option>
                        <option value="upcycling">Upcycling</option>
                    </select>
                </div>

                <div class="form-group">
                    <label for="texte">Article</label>
                    <textarea id="texte" placeholder="Écrivez votre article..." v-model="formulaire.contenu"></textarea>
                </div>

                <div class="form-buttons">
                    <button type="button" class="button-draft" @click="enregistrer('brouillon')">Enregistrer en brouillon</button>
                    <button type="button" class="button-submit" @click="enregistrer('publie')">Publier</button>
                </div>
            </div>

            <div class="evenements">
                <div v-for="a in articles" :key="a.id" class="evenement-row">
                    <span class="evenement-titre">{{ a.titre }}</span>
                    <span class="evenement-date">{{ formaterDate(a.updated_at) }}</span>
                    <span class="badge-statut" :class="'statut-' + a.statut">{{ formaterStatut(a.statut) }}</span>
                    <button class="btn-modifier" @click="modifier(a)">Modifier</button>
                </div>
            </div>

        </main>
    </div>

    <script src="js/articles.js"></script>
</body>

</html>
