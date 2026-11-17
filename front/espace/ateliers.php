<?php
require "inclus/espace.php";
$titrePage = "salaries.ateliers.titre";
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
                <h1>{{ $t('salaries.ateliers.titre') }}</h1>
                <button class="btn-primary" @click="ouvrirCreation">+ {{ $t('salaries.ateliers.creer') }}</button>
            </div>

            <div class="filtres">
                <button v-for="filtre in filtres" :key="filtre.valeur"
                    class="pill" :class="{ active: statutSelectionne === filtre.valeur }"
                    @click="statutSelectionne = filtre.valeur">
                    {{ $t(filtre.nom) }}
                </button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="confirmation" class="succes">{{ $t('salaries.ateliers.enregistre') }}</p>

            <div class="ateliers">
                <div v-for="atelier in ateliersFiltres" :key="atelier.id" class="atelier-row">
                    <span class="atelier-titre">{{ atelier.titre }}</span>
                    <span class="atelier-date">{{ formaterDate(atelier.date_debut) }}</span>
                    <span class="badge-statut" :class="'statut-' + atelier.statut">
                        {{ formaterStatut(atelier.statut) }}
                    </span>
                    <button class="btn-modifier" @click="modifierAtelier(atelier.id)">{{ $t('commun.modifier') }}</button>
                    <button v-if="atelier.statut === 'en_attente'" class="btn-modifier" @click="valider(atelier.id)">{{ $t('commun.valider') }}</button>
                    <button class="btn-modifier" @click="supprimer(atelier.id)">{{ $t('commun.supprimer') }}</button>
                </div>
            </div>

            <fenetre-modale :ouverte="formulaireOuvert" :titre="editionId ? $t('salaries.ateliers.modifier') : $t('salaries.ateliers.creer')" @fermer="formulaireOuvert = false">
                <form @submit.prevent>

                <div class="form-group">
                    <label for="titre-offre">{{ $t('commun.titre') }}</label>
                    <input type="text" id="titre-offre" v-model="formulaire.titre" required>
                </div>

                <div class="form-group">
                    <label for="description">{{ $t('commun.description') }}</label>
                    <textarea id="description" v-model="formulaire.description"></textarea>
                </div>

                <div class="form-group">
                    <label for="date-debut">{{ $t('commun.date_debut') }}</label>
                    <input type="datetime-local" id="date-debut" v-model="formulaire.date_debut" required>
                </div>

                <div class="form-group">
                    <label for="date-fin">{{ $t('commun.date_fin') }}</label>
                    <input type="datetime-local" id="date-fin" v-model="formulaire.date_fin">
                </div>

                <div class="form-group">
                    <label for="lieu">{{ $t('commun.lieu') }}</label>
                    <input type="text" id="lieu" v-model="formulaire.lieu">
                </div>

                <p v-if="erreurFormulaire" class="erreur">{{ erreurFormulaire }}</p>

                <div class="form-buttons">
                    <button type="button" class="button-draft" @click="envoyerAtelier('brouillon')">
                        {{ $t('commun.enregistrer_brouillon') }}
                    </button>
                    <button type="button" class="button-submit" @click="envoyerAtelier('en_attente')">
                        {{ $t('commun.soumettre_validation') }}
                    </button>
                </div>

                </form>
            </fenetre-modale>

        </main>
    </div>

    <script src="<?= ressource("js/ateliers.js") ?>"></script>
</body>

</html>
