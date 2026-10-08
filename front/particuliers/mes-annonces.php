<?php
require "inclus/espace.php";
$titrePage = "nav.mes_annonces";
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
                <h1>{{ t('nav.mes_annonces') }}</h1>
                <button class="btn-primary" @click="reinitialiserFormulaire">+ {{ t('annonces.nouvelle') }}</button>
            </div>

            <p v-if="erreurListe" class="erreur">{{ erreurListe }}</p>

            <div class="deux-colonnes">

                <section>
                    <p v-if="annonces.length === 0" class="discret">{{ t('annonces.aucune_mienne') }}</p>
                    <div v-for="a in annonces" :key="a.id" class="carte section">
                        <span class="badge" :class="{ 'badge-attente': a.statut !== 'en_ligne' }">{{ t('annonce.' + a.statut) }}</span>
                        <h3>{{ a.titre }}</h3>
                        <p class="discret">
                            {{ t('annonce.type.' + a.type) }}<span v-if="a.type === 'vente'"> · {{ formaterPrix(a.prix) }}</span> · {{ a.localisation }}
                        </p>
                        <div class="form-buttons">
                            <button class="button-draft" @click="modifier(a)">{{ t('commun.modifier') }}</button>
                            <button v-if="['en_ligne', 'reservee'].includes(a.statut)" class="button-submit" @click="marquerCedee(a)">{{ t('annonces.marquer_cedee') }}</button>
                            <button class="button-draft" @click="supprimer(a.id)">{{ t('commun.supprimer') }}</button>
                        </div>
                    </div>
                </section>

                <div class="carte">
                    <h2>{{ formulaire.id ? t('annonces.modifier') : t('annonces.deposer') }}</h2>

                    <form @submit.prevent="enregistrer">
                        <div class="form-group">
                            <label for="titre">{{ t('commun.titre') }}</label>
                            <input id="titre" type="text" v-model="formulaire.titre" required>
                        </div>

                        <div class="form-group">
                            <label for="description">{{ t('commun.description') }}</label>
                            <textarea id="description" v-model="formulaire.description"></textarea>
                        </div>

                        <div class="form-group">
                            <label for="type">{{ t('commun.type') }}</label>
                            <select id="type" v-model="formulaire.type">
                                <option value="don">{{ t('annonce.type.don') }}</option>
                                <option value="vente">{{ t('annonce.type.vente') }}</option>
                            </select>
                        </div>

                        <div v-if="formulaire.type === 'vente'" class="form-group">
                            <label for="prix">{{ t('commun.prix') }} (€)</label>
                            <input id="prix" type="number" min="0.01" step="0.01" v-model.number="formulaire.prix" required>
                        </div>

                        <div class="form-group">
                            <label for="categorie">{{ t('annonces.categorie_materiau') }}</label>
                            <select id="categorie" v-model="formulaire.categorie_id">
                                <option :value="null">{{ t('commun.non_precisee') }}</option>
                                <option v-for="c in categories" :key="c.id" :value="c.id">{{ traduireOu('categorie.' + c.code, c.code) }}</option>
                            </select>
                        </div>

                        <div class="form-group">
                            <label for="localisation">{{ t('commun.localisation') }}</label>
                            <input id="localisation" type="text" v-model="formulaire.localisation">
                        </div>

                        <div v-if="formulaire.id && statutsAnnonce.includes(formulaire.statut)" class="form-group">
                            <label for="statut">{{ t('commun.statut') }}</label>
                            <select id="statut" v-model="formulaire.statut">
                                <option v-for="statut in statutsAnnonce" :key="statut" :value="statut">{{ t('annonce.' + statut) }}</option>
                            </select>
                        </div>

                        <p v-if="erreur" class="erreur">{{ erreur }}</p>
                        <p v-if="message" class="succes">{{ message }}</p>

                        <div class="form-buttons">
                            <button type="submit" class="button-submit">{{ t('commun.enregistrer') }}</button>
                            <button type="button" class="button-draft" @click="reinitialiserFormulaire">{{ t('commun.annuler') }}</button>
                        </div>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/mes-annonces.js"></script>
</body>

</html>
