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
                <h1>{{ t('nav.conseils') }}</h1>
            </div>

            <div class="barre-filtres">
                <div class="form-group">
                    <label for="recherche">{{ t('commun.rechercher') }}</label>
                    <input id="recherche" type="search" v-model="recherche" :placeholder="t('conseils.mot_cle')">
                </div>
                <div class="form-group">
                    <label for="categorie">{{ t('commun.categorie') }}</label>
                    <select id="categorie" v-model="categorie">
                        <option value="">{{ t('commun.toutes') }}</option>
                        <option v-for="c in categories" :key="c" :value="c">{{ traduireOu('categorie_article.' + c, c) }}</option>
                    </select>
                </div>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="!erreur && conseilsFiltres.length === 0" class="discret">{{ t('conseils.aucun') }}</p>

            <div class="deux-colonnes">
                <div class="grille">
                    <div v-for="c in conseilsFiltres" :key="c.id" class="carte carte-cliquable"
                        :class="{ 'carte-selectionnee': selection && selection.id === c.id }" @click="selection = c">
                        <span v-if="c.categorie" class="badge">{{ traduireOu('categorie_article.' + c.categorie, c.categorie) }}</span>
                        <h3>{{ c.titre }}</h3>
                        <p class="discret">{{ formaterDate(c.created_at) }}</p>
                    </div>
                </div>

                <article v-if="selection" class="carte">
                    <span v-if="selection.categorie" class="badge">{{ traduireOu('categorie_article.' + selection.categorie, selection.categorie) }}</span>
                    <h2>{{ selection.titre }}</h2>
                    <p class="discret">{{ t('conseils.publie_le', { date: formaterDate(selection.created_at) }) }}</p>
                    <p>{{ selection.contenu }}</p>
                </article>
            </div>

        </main>
    </div>

    <script src="js/conseils.js"></script>
</body>

</html>
