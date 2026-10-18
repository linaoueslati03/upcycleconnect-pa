<?php
require "inclus/espace.php";
$titrePage = "nav.conseils";
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
                <h1>{{ $t('nav.conseils') }}</h1>
                <button v-if="roleConnecte() === 'salarie'" class="btn-primary" @click="nouvelArticle">{{ $t('salaries.articles.creer') }}</button>
            </div>

            <section v-if="roleConnecte() === 'salarie'" class="section">
                <h2>{{ $t('salaries.articles.titre') }}</h2>
                <p v-if="erreurArticle" class="erreur">{{ erreurArticle }}</p>

                <div v-if="formulaireVisible" class="form-article">
                    <h2>{{ formulaire.id ? $t('salaries.articles.modifier') : $t('salaries.articles.nouveau') }}</h2>

                    <div class="form-group">
                        <label for="titre">{{ $t('commun.titre') }}</label>
                        <input id="titre" type="text" :placeholder="$t('salaries.articles.placeholder_titre')" v-model="formulaire.titre">
                    </div>

                    <div class="form-group">
                        <label for="categorie-article">{{ $t('commun.categorie') }}</label>
                        <select id="categorie-article" v-model="formulaire.categorie">
                            <option value="">{{ $t('commun.choisir_categorie') }}</option>
                            <option value="actualites">{{ $t('categorie_article.actualites') }}</option>
                            <option value="conseils">{{ $t('categorie_article.conseils') }}</option>
                            <option value="evenements">{{ $t('categorie_article.evenements') }}</option>
                            <option value="upcycling">{{ $t('categorie_article.upcycling') }}</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="texte">{{ $t('salaries.articles.article') }}</label>
                        <textarea id="texte" :placeholder="$t('salaries.articles.placeholder_texte')" v-model="formulaire.contenu"></textarea>
                    </div>

                    <div class="form-buttons">
                        <button type="button" class="button-draft" @click="enregistrer('brouillon')">{{ $t('commun.enregistrer_brouillon') }}</button>
                        <button type="button" class="button-submit" @click="enregistrer('publie')">{{ $t('commun.publier') }}</button>
                    </div>
                </div>

                <div class="evenements">
                    <div v-for="a in articles" :key="a.id" class="evenement-row">
                        <span class="evenement-titre">{{ a.titre }}</span>
                        <span class="evenement-date">{{ formaterDate(a.updated_at) }}</span>
                        <span class="badge-statut" :class="'statut-' + a.statut">{{ formaterStatut(a.statut) }}</span>
                        <button class="btn-modifier" @click="modifier(a)">{{ $t('commun.modifier') }}</button>
                    </div>
                </div>
            </section>

            <div class="barre-filtres">
                <div class="form-group">
                    <label for="recherche">{{ $t('commun.rechercher') }}</label>
                    <input id="recherche" type="search" v-model="recherche" :placeholder="$t('conseils.mot_cle')">
                </div>
                <div class="form-group">
                    <label for="categorie">{{ $t('commun.categorie') }}</label>
                    <select id="categorie" v-model="categorie">
                        <option value="">{{ $t('commun.toutes') }}</option>
                        <option v-for="c in categories" :key="c" :value="c">{{ $t('categorie_article.' + c, c) }}</option>
                    </select>
                </div>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="!erreur && conseilsFiltres.length === 0" class="discret">{{ $t('conseils.aucun') }}</p>

            <div class="deux-colonnes">
                <div class="grille">
                    <div v-for="c in conseilsFiltres" :key="c.id" class="carte carte-cliquable"
                        :class="{ 'carte-selectionnee': selection && selection.id === c.id }" @click="selection = c">
                        <span v-if="c.categorie" class="badge">{{ $t('categorie_article.' + c.categorie, c.categorie) }}</span>
                        <h3>{{ c.titre }}</h3>
                        <p class="discret">{{ formaterDate(c.created_at) }}</p>
                    </div>
                </div>

                <article v-if="selection" class="carte">
                    <span v-if="selection.categorie" class="badge">{{ $t('categorie_article.' + selection.categorie, selection.categorie) }}</span>
                    <h2>{{ selection.titre }}</h2>
                    <p class="discret">{{ $t('conseils.publie_le', { date: formaterDate(selection.created_at) }) }}</p>
                    <p>{{ selection.contenu }}</p>
                </article>
            </div>

        </main>
    </div>

    <script src="js/conseils.js"></script>
</body>

</html>
