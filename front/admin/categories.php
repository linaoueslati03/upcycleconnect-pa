<?php
require "inclus/espace.php";
$titrePage = "nav.categories";
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
                <h1>{{ $t('admin.categories.titre') }}</h1>
            </div>

            <p class="discret">{{ $t('admin.categories.explication') }}</p>
            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <section v-for="(groupe, type) in groupes" :key="type" class="carte section-categories">
                <div class="entete-carte">
                    <h2>{{ $t('admin.categories.' + type) }}</h2>
                    <button class="btn-primary" @click="ouvrirAjout(type)">+ {{ $t('admin.categories.ajouter') }}</button>
                </div>

                <table class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('commun.nom') }}</th>
                            <th>{{ $t('commun.statut') }}</th>
                            <th>{{ $t('commun.actions') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="c in groupe.liste" :key="c.id">
                            <td>{{ c.libelle }}</td>
                            <td :class="{ 'statut-suspendu': !c.actif }">
                                {{ c.actif ? $t('admin.categories.active') : $t('admin.categories.inactive') }}
                            </td>
                            <td>
                                <button class="lien-action" @click="commencerEdition(type, c)">{{ $t('admin.categories.renommer') }}</button>
                                <button class="lien-action" @click="enregistrer(type, c, c.libelle, !c.actif)">
                                    {{ c.actif ? $t('admin.categories.desactiver') : $t('admin.categories.activer') }}
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </section>

            <fenetre-modale :ouverte="edition.ouvert" :titre="$t('admin.categories.renommer') + (edition.categorie ? ' · ' + edition.categorie.libelle : '')" @fermer="annulerEdition">
                <form @submit.prevent="enregistrer(edition.type, edition.categorie, edition.libelle, edition.categorie.actif)">
                    <div class="form-group">
                        <label for="libelle-categorie">{{ $t('commun.nom') }}</label>
                        <input id="libelle-categorie" type="text" v-model="edition.libelle" required>
                    </div>
                    <p v-if="erreur" class="erreur">{{ erreur }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                        <button type="button" class="button-draft" @click="annulerEdition">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

            <fenetre-modale :ouverte="ajoutOuvert" :titre="$t('admin.categories.ajouter') + ' · ' + $t('admin.categories.' + ajoutType)" @fermer="ajoutOuvert = false">
                <form @submit.prevent="ajouter(ajoutType)">
                    <div class="form-group">
                        <label for="nouvelle-categorie">{{ $t('admin.categories.nouvelle') }}</label>
                        <input id="nouvelle-categorie" type="text" v-model="groupes[ajoutType].nouvelle" required>
                    </div>
                    <p v-if="erreur" class="erreur">{{ erreur }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('admin.categories.ajouter') }}</button>
                        <button type="button" class="button-draft" @click="ajoutOuvert = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

            <section class="carte section-categories">
                <h2>{{ $t('admin.categories.impact.titre') }}</h2>
                <p class="discret">{{ $t('admin.categories.impact.explication') }}</p>
                <table class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('commun.nom') }}</th>
                            <th class="montant">{{ $t('admin.categories.impact.facteur') }}</th>
                            <th>{{ $t('admin.categories.impact.source') }}</th>
                            <th>{{ $t('commun.actions') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="c in groupes.materiaux.liste" :key="c.id">
                            <td>{{ c.libelle }}</td>
                            <td class="montant">{{ c.co2_evite_par_kg !== null ? c.co2_evite_par_kg.toLocaleString(langueCourante()) : '—' }}</td>
                            <td class="source-facteur">{{ c.source_facteur || $t('admin.categories.impact.sans_source') }}</td>
                            <td><button class="lien-action" @click="ouvrirImpact(c)">{{ $t('commun.modifier') }}</button></td>
                        </tr>
                    </tbody>
                </table>
            </section>

            <fenetre-modale :ouverte="impact.ouvert" :titre="$t('admin.categories.impact.titre') + (impact.categorie ? ' · ' + impact.categorie.libelle : '')" @fermer="impact.ouvert = false">
                <form @submit.prevent="enregistrerImpact">
                    <div class="form-group">
                        <label for="facteur-co2">{{ $t('admin.categories.impact.facteur') }}</label>
                        <input id="facteur-co2" type="number" v-model.number="impact.facteur" min="0" max="100" step="0.001" required>
                    </div>
                    <div class="form-group">
                        <label for="source-co2">{{ $t('admin.categories.impact.source') }}</label>
                        <textarea id="source-co2" v-model="impact.source" rows="3" required></textarea>
                    </div>
                    <p v-if="impact.erreur" class="erreur">{{ impact.erreur }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                        <button type="button" class="button-draft" @click="impact.ouvert = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

        </main>
    </div>

    <script src="<?= ressource("js/categories.js") ?>"></script>
</body>

</html>
