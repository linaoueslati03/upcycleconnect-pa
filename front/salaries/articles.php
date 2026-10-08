<?php
require "inclus/espace.php";
$titrePage = "nav.articles";
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
                <h1>{{ t('salaries.articles.titre') }}</h1>
                <button class="btn-primary" @click="nouvelArticle">{{ t('salaries.articles.creer') }}</button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div v-if="formulaireVisible" class="form-article">
                <h2>{{ formulaire.id ? t('salaries.articles.modifier') : t('salaries.articles.nouveau') }}</h2>

                <div class="form-group">
                    <label for="titre">{{ t('commun.titre') }}</label>
                    <input id="titre" type="text" :placeholder="t('salaries.articles.placeholder_titre')" v-model="formulaire.titre">
                </div>

                <div class="form-group">
                    <label for="categorie">{{ t('commun.categorie') }}</label>
                    <select id="categorie" v-model="formulaire.categorie">
                        <option value="">{{ t('commun.choisir_categorie') }}</option>
                        <option value="actualites">{{ t('categorie_article.actualites') }}</option>
                        <option value="conseils">{{ t('categorie_article.conseils') }}</option>
                        <option value="evenements">{{ t('categorie_article.evenements') }}</option>
                        <option value="upcycling">{{ t('categorie_article.upcycling') }}</option>
                    </select>
                </div>

                <div class="form-group">
                    <label for="texte">{{ t('salaries.articles.article') }}</label>
                    <textarea id="texte" :placeholder="t('salaries.articles.placeholder_texte')" v-model="formulaire.contenu"></textarea>
                </div>

                <div class="form-buttons">
                    <button type="button" class="button-draft" @click="enregistrer('brouillon')">{{ t('commun.enregistrer_brouillon') }}</button>
                    <button type="button" class="button-submit" @click="enregistrer('publie')">{{ t('commun.publier') }}</button>
                </div>
            </div>

            <div class="evenements">
                <div v-for="a in articles" :key="a.id" class="evenement-row">
                    <span class="evenement-titre">{{ a.titre }}</span>
                    <span class="evenement-date">{{ formaterDate(a.updated_at) }}</span>
                    <span class="badge-statut" :class="'statut-' + a.statut">{{ formaterStatut(a.statut) }}</span>
                    <button class="btn-modifier" @click="modifier(a)">{{ t('commun.modifier') }}</button>
                </div>
            </div>

        </main>
    </div>

    <script src="js/articles.js"></script>
</body>

</html>
