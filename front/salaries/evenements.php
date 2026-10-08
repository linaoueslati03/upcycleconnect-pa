<?php
require "inclus/espace.php";
$titrePage = "Mes événements";
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
                <h1>Mes événements</h1>
                <a class="btn-primary" href="creer-evenement.php">+ Créer un événement</a>
            </div>

            <div class="filtres">
                <button v-for="filtre in filtres" :key="filtre.valeur"
                    class="pill" :class="{ active: statutSelectionne === filtre.valeur }"
                    @click="statutSelectionne = filtre.valeur">
                    {{ filtre.nom }}
                </button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="evenements">
                <div v-for="evenement in evenementsFiltres" :key="evenement.id" class="evenement-row">
                    <span class="evenement-titre">{{ evenement.titre }}</span>
                    <span class="evenement-date">{{ formaterDate(evenement.date_debut) }}</span>
                    <span class="badge-statut" :class="'statut-' + evenement.statut">
                        {{ formaterStatut(evenement.statut) }}
                    </span>
                    <button class="btn-modifier" @click="modifierEvenement(evenement.id)">Modifier</button>
                    <button v-if="evenement.statut === 'en_attente'" class="btn-modifier" @click="valider(evenement.id)">Valider</button>
                    <button class="btn-modifier" @click="supprimer(evenement.id)">Supprimer</button>
                </div>
            </div>

        </main>
    </div>

    <script src="js/evenements.js"></script>
</body>

</html>
