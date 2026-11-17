<?php
require "inclus/espace.php";
$titrePage = "salaries.evenements.titre";
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
                <h1>{{ $t('salaries.evenements.titre') }}</h1>
                <button class="btn-primary" @click="ouvrirCreation">+ {{ $t('salaries.evenements.creer') }}</button>
            </div>

            <div class="filtres">
                <button v-for="filtre in filtres" :key="filtre.valeur"
                    class="pill" :class="{ active: statutSelectionne === filtre.valeur }"
                    @click="statutSelectionne = filtre.valeur">
                    {{ $t(filtre.nom) }}
                </button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="confirmation" class="succes">{{ $t('salaries.evenements.enregistre') }}</p>

            <div class="evenements">
                <div v-for="evenement in evenementsFiltres" :key="evenement.id" class="evenement-row">
                    <span class="evenement-titre">{{ evenement.titre }}</span>
                    <span class="evenement-date">{{ formaterDate(evenement.date_debut) }}</span>
                    <span class="badge-statut" :class="'statut-' + evenement.statut">
                        {{ formaterStatut(evenement.statut) }}
                    </span>
                    <button class="btn-modifier" @click="modifierEvenement(evenement.id)">{{ $t('commun.modifier') }}</button>
                    <button v-if="evenement.statut === 'en_attente'" class="btn-modifier" @click="valider(evenement.id)">{{ $t('commun.valider') }}</button>
                    <button class="btn-modifier" @click="supprimer(evenement.id)">{{ $t('commun.supprimer') }}</button>
                </div>
            </div>

            <fenetre-modale :ouverte="formulaireOuvert" :titre="editionId ? $t('salaries.evenements.modifier') : $t('salaries.evenements.creer')" @fermer="formulaireOuvert = false">
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

                <div class="form-group">
                    <label for="site">{{ $t('commun.site') }}</label>
                    <select id="site" v-model="formulaire.site">
                        <option value="">{{ $t('salaries.evenements.choisir_site') }}</option>
                        <option value="Paris 10e">Paris 10e</option>
                        <option value="Paris 11e">Paris 11e</option>
                        <option value="Paris 13e">Paris 13e</option>
                        <option value="Montreuil">Montreuil</option>
                        <option value="Suisse">{{ $t('site.suisse') }}</option>
                    </select>
                </div>

                <p v-if="erreurFormulaire" class="erreur">{{ erreurFormulaire }}</p>

                <div class="form-buttons">
                    <button type="button" class="button-draft" @click="envoyerEvenement('brouillon')">
                        {{ $t('commun.enregistrer_brouillon') }}
                    </button>
                    <button type="button" class="button-submit" @click="envoyerEvenement('en_attente')">
                        {{ $t('commun.soumettre_validation') }}
                    </button>
                </div>

                </form>
            </fenetre-modale>

        </main>
    </div>

    <script src="<?= ressource("js/evenements.js") ?>"></script>
</body>

</html>
