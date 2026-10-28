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
                <a class="btn-primary" href="creer-formation.php">
                    + {{ $t('salaries.formations.creer') }}
                </a>
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
        </main>
    </div>

    <script src="js/formations.js"></script>
</body>

</html>