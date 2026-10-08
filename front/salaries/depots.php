<?php
require "inclus/espace.php";
$titrePage = "Dépôts en conteneur";
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
                <h1>Dépôts en conteneur</h1>
            </div>

            <div class="filtres">
                <button v-for="filtre in filtres" :key="filtre.valeur"
                    class="pill" :class="{ active: statutSelectionne === filtre.valeur }"
                    @click="statutSelectionne = filtre.valeur; chargerDepots()">
                    {{ filtre.nom }}
                </button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="!erreur && depots.length === 0">Aucun dépôt.</p>

            <div class="evenements">
                <div v-for="d in depots" :key="d.id" class="evenement-row">
                    <span class="evenement-titre">{{ d.description_objet }}</span>
                    <span class="evenement-date">{{ formaterDate(d.date_demande) }}</span>
                    <span class="badge-statut">{{ libellesStatut[d.statut] }}</span>
                    <span v-if="d.code_ouverture">Code : {{ d.code_ouverture }}</span>
                    <button v-if="statutSuivant[d.statut]" class="btn-modifier" @click="avancer(d)">
                        {{ libellesAction[d.statut] }}
                    </button>
                </div>
            </div>

        </main>
    </div>

    <script src="js/depots.js"></script>
</body>

</html>
