<?php
require "inclus/espace.php";
$titrePage = "nav.prestations";
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
                <h1>{{ t('admin.prestations.titre') }}</h1>
                <button class="btn-primary" @click="reinitialiserFormulaire">+ {{ t('admin.prestations.nouvelle') }}</button>
            </div>

            <p v-if="erreurListe" class="erreur">{{ erreurListe }}</p>

            <div class="mise-en-page-admin">

                <section class="carte">
                    <table class="tableau">
                        <thead>
                            <tr>
                                <th>{{ t('commun.titre') }}</th>
                                <th>{{ t('commun.categorie') }}</th>
                                <th>{{ t('commun.tarif') }}</th>
                                <th>{{ t('commun.statut') }}</th>
                                <th>{{ t('commun.actions') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="p in prestations" :key="p.id">
                                <td><strong>{{ p.titre }}</strong></td>
                                <td>{{ p.categorie }}</td>
                                <td>{{ formaterTarif(p.tarif) }}</td>
                                <td>{{ t('statut.' + p.statut) }}</td>
                                <td>
                                    <button class="lien-action" @click="modifier(p)">{{ t('commun.modifier') }}</button>
                                    <a class="lien-action" :href="lienPDF(p.id)">PDF</a>
                                    <button class="lien-action lien-supprimer" @click="supprimer(p.id)">{{ t('commun.supprimer') }}</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </section>

                <div class="carte carte-formulaire">
                    <h2>{{ formulaire.id ? t('admin.prestations.modifier') : t('admin.prestations.nouvelle') }}</h2>

                    <form @submit.prevent="enregistrer">
                        <div class="form-group">
                            <label for="titre">{{ t('commun.titre') }}</label>
                            <input id="titre" type="text" v-model="formulaire.titre" required>
                        </div>

                        <div class="form-group">
                            <label for="categorie">{{ t('commun.categorie') }}</label>
                            <input id="categorie" type="text" v-model="formulaire.categorie" required>
                        </div>

                        <div class="form-group">
                            <label for="tarif">{{ t('commun.tarif') }} (€)</label>
                            <input id="tarif" type="number" min="0" step="0.01" v-model.number="formulaire.tarif" required>
                        </div>

                        <div class="form-group">
                            <label for="statut">{{ t('commun.statut') }}</label>
                            <select id="statut" v-model="formulaire.statut">
                                <option value="brouillon">{{ t('statut.brouillon') }}</option>
                                <option value="publie">{{ t('statut.publie') }}</option>
                            </select>
                        </div>

                        <p v-if="erreur" class="erreur">{{ erreur }}</p>

                        <div class="form-buttons">
                            <button type="submit" class="button-submit">{{ t('commun.enregistrer') }}</button>
                            <button type="button" class="button-draft" @click="reinitialiserFormulaire">{{ t('commun.annuler') }}</button>
                        </div>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/prestations.js"></script>
</body>

</html>
