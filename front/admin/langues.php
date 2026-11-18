<?php
require "inclus/espace.php";
$titrePage = "nav.langues";
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
                <h1>{{ $t('admin.langues.titre') }}</h1>
                <button class="btn-primary" @click="ouvrirAjout">+ {{ $t('admin.langues.ajouter') }}</button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="succes" class="succes">{{ succes }}</p>

            <div class="mise-en-page-admin">

                <section class="carte">
                    <table class="tableau">
                        <thead>
                            <tr>
                                <th>{{ $t('admin.langues.code') }}</th>
                                <th>{{ $t('commun.nom') }}</th>
                                <th>{{ $t('admin.langues.fichier') }}</th>
                                <th>{{ $t('commun.statut') }}</th>
                                <th>{{ $t('commun.actions') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="l in langues" :key="l.id">
                                <td>{{ l.code }}</td>
                                <td>{{ l.libelle }}</td>
                                <td :class="{ 'statut-suspendu': !fichiers[l.code] }">
                                    {{ fichiers[l.code] ? $t('admin.langues.fichier_present') : $t('admin.langues.fichier_manquant') }}
                                </td>
                                <td>{{ l.actif ? $t('admin.langues.active') : $t('admin.langues.inactive') }}</td>
                                <td>
                                    <template v-if="l.code !== langueParDefaut">
                                        <label class="lien-action">
                                            {{ $t('admin.langues.importer') }}
                                            <input type="file" accept=".json,application/json" hidden @change="importer(l, $event)">
                                        </label>
                                        <button class="lien-action" @click="basculer(l)">
                                            {{ l.actif ? $t('admin.langues.desactiver') : $t('admin.langues.activer') }}
                                        </button>
                                        <button class="lien-action lien-supprimer" @click="demanderSuppression(l)">
                                            {{ $t('commun.supprimer') }}
                                        </button>
                                    </template>
                                    <span v-else class="discret">{{ $t('admin.langues.reference') }}</span>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </section>

                <div class="carte carte-formulaire">
                    <h2>{{ $t('admin.langues.procedure_titre') }}</h2>
                    <ol>
                        <li>{{ $t('admin.langues.etape_1') }}</li>
                        <li>
                            {{ $t('admin.langues.etape_2') }}
                            <a href="../assets/locales/fr.json" download="modele-traduction.json">{{ $t('admin.langues.telecharger_modele') }}</a>
                        </li>
                        <li>{{ $t('admin.langues.etape_3') }}</li>
                    </ol>
                </div>

            </div>

            <fenetre-modale :ouverte="suppression.ouverte" @fermer="suppression.ouverte = false"
                :titre="$t('admin.langues.supprimer_titre', { langue: suppression.langue ? suppression.langue.libelle : '' })">
                <template v-if="suppression.langue">
                    <p>{{ $t('admin.langues.supprimer_intro', { langue: suppression.langue.libelle, code: suppression.langue.code }) }}</p>
                    <ul class="details-offre">
                        <li>{{ $t('admin.langues.supprimer_ligne') }}</li>
                        <li>{{ $t('admin.langues.supprimer_fichier', { fichier: suppression.langue.code + '.json' }) }}</li>
                        <li>{{ $t('admin.langues.supprimer_utilisateurs', { nombre: suppression.langue.nb_utilisateurs }) }}</li>
                    </ul>
                    <p class="erreur">{{ $t('admin.langues.supprimer_irreversible') }}</p>
                    <p v-if="erreurSuppression" class="erreur">{{ erreurSuppression }}</p>
                    <div class="form-buttons">
                        <button type="button" class="button-submit bouton-refus" :disabled="suppression.enCours" @click="supprimer">
                            {{ $t('admin.langues.supprimer_confirmer') }}
                        </button>
                        <button type="button" class="button-draft" @click="suppression.ouverte = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </template>
            </fenetre-modale>

            <fenetre-modale :ouverte="ajoutOuvert" :titre="$t('admin.langues.ajouter')" @fermer="ajoutOuvert = false">
                    <form @submit.prevent="ajouter">
                        <div class="form-group">
                            <label for="code">{{ $t('admin.langues.code') }}</label>
                            <input id="code" type="text" v-model="formulaire.code" maxlength="5" placeholder="es" required>
                        </div>
                        <div class="form-group">
                            <label for="libelle">{{ $t('commun.nom') }}</label>
                            <input id="libelle" type="text" v-model="formulaire.libelle" placeholder="Español" required>
                        </div>

                        <p v-if="erreurFormulaire" class="erreur">{{ erreurFormulaire }}</p>

                        <div class="form-buttons">
                            <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                            <button type="button" class="button-draft" @click="ajoutOuvert = false">{{ $t('commun.annuler') }}</button>
                        </div>
                    </form>
            </fenetre-modale>

        </main>
    </div>

    <script src="<?= ressource("js/langues.js") ?>"></script>
</body>

</html>
