<?php
require "inclus/espace.php";
$titrePage = "salaries.formations.titre";
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
                <h1>{{ $t('salaries.formations.titre') }}</h1>
                <button class="btn-primary" @click="ouvrirCreation">
                    + {{ $t('salaries.formations.creer') }}
                </button>
            </div>

            <div class="filtres">
                <button
                    v-for="filtre in filtres"
                    :key="filtre.valeur"
                    class="pill"
                    :class="{ active: statutSelectionne === filtre.valeur }"
                    @click="statutSelectionne = filtre.valeur">
                    {{ $t(filtre.nom) }}
                </button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="succes" class="succes">{{ succes }}</p>

            <p v-if="chargement">
                {{ $t('commun.chargement') }}
            </p>

            <p v-else-if="formationsFiltres.length === 0">
                {{ $t('salaries.formations.aucune') }}
            </p>

            <div v-else class="ateliers">
                <div
                    v-for="formation in formationsFiltres"
                    :key="formation.id"
                    class="atelier-row">

                    <div>
                        <span class="atelier-titre">
                            {{ formation.titre }}
                        </span>

                        <span class="atelier-date">
                            {{ formaterDate(formation.date_debut) }}
                        </span>

                        <span class="atelier-date">
                            {{ formation.nb_places ?? '—' }}
                            {{ $t('salaries.formations.places') }}
                        </span>

                        <span class="atelier-date">
                            {{ Number(formation.tarif).toLocaleString(langueCourante()) }} €
                        </span>
                    </div>

                    <span
                        class="badge-statut"
                        :class="'statut-' + formation.statut">
                        {{ formaterStatut(formation.statut) }}
                    </span>

                    <button
                        class="btn-modifier"
                        @click="modifierFormation(formation.id)">
                        {{ $t('commun.modifier') }}
                    </button>

                    <button
                        v-if="formation.statut === 'en_attente'"
                        class="btn-modifier"
                        @click="valider(formation.id)">
                        {{ $t('commun.valider') }}
                    </button>

                    <button
                        class="btn-modifier"
                        @click="supprimer(formation.id)">
                        {{ $t('commun.supprimer') }}
                    </button>
                </div>
            </div>
            <fenetre-modale :ouverte="formulaireOuvert" :titre="editionId ? $t('salaries.formations.modifier') : $t('salaries.formations.creer')" @fermer="formulaireOuvert = false">
                <form @submit.prevent="enregistrer">
                <div class="form-group">
                    <label for="titre-formation">{{ $t('commun.titre') }}</label>
                    <input id="titre-formation" v-model.trim="formation.titre" required maxlength="150">
                </div>

                <div class="form-group">
                    <label for="description">{{ $t('commun.description') }}</label>
                    <textarea id="description" v-model="formation.description"></textarea>
                </div>

                <div class="form-group">
                    <label for="date_debut">{{ $t('commun.date_debut') }}</label>
                    <input id="date_debut" type="datetime-local" v-model="formation.date_debut" required>
                </div>

                <div class="form-group">
                    <label for="date_fin">{{ $t('commun.date_fin') }}</label>
                    <input id="date_fin" type="datetime-local" v-model="formation.date_fin">
                </div>

                <div class="form-group">
                    <label for="lieu">{{ $t('commun.lieu') }}</label>
                    <input id="lieu" v-model.trim="formation.lieu" maxlength="150">
                </div>

                <div class="form-group">
                    <label for="nb_places">{{ $t('salaries.formations.places') }}</label>
                    <input id="nb_places" type="number" min="1" v-model.number="formation.nb_places">
                </div>

                <div class="form-group">
                    <label for="tarif">{{ $t('commun.tarif') }} (€)</label>
                    <input id="tarif" type="number" min="0" step="0.01" v-model.number="formation.tarif" required>
                </div>

                <div class="form-group">
                    <label for="materiel_necessaire">
                        {{ $t('salaries.formations.materiel') }}
                    </label>
                    <textarea id="materiel_necessaire" v-model="formation.materiel_necessaire"></textarea>
                </div>

                <div class="form-group">
                    <label for="statut">{{ $t('commun.statut') }}</label>
                    <select id="statut" v-model="formation.statut">
                        <option value="brouillon">{{ formaterStatut('brouillon') }}</option>
                        <option value="en_attente">{{ formaterStatut('en_attente') }}</option>
                    </select>
                </div>

                <p v-if="erreurFormulaire" class="erreur">{{ erreurFormulaire }}</p>

                <div class="form-buttons">
                    <button type="submit" class="button-submit" :disabled="enregistrement">
                        {{ enregistrement
                            ? $t('commun.chargement')
                            : $t('commun.enregistrer') }}
                    </button>
                    <button type="button" class="button-draft" @click="formulaireOuvert = false">
                        {{ $t('commun.annuler') }}
                    </button>
                </div>
                </form>
            </fenetre-modale>
        </main>
    </div>

    <script src="<?= ressource("js/formations.js") ?>"></script>
</body>

</html>