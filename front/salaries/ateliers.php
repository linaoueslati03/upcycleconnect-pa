<?php
require "inclus/espace.php";
$titrePage = "Mes ateliers";
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
                <h1>Mes ateliers</h1>
                <a class="btn-primary" href="creer-atelier.php">+ Créer un atelier</a>
            </div>

            <div class="filtres">
                <button v-for="filtre in filtres" :key="filtre.valeur"
                    class="pill" :class="{ active: statutSelectionne === filtre.valeur }"
                    @click="statutSelectionne = filtre.valeur">
                    {{ filtre.nom }}
                </button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="ateliers">
                <div v-for="atelier in ateliersFiltres" :key="atelier.id" class="atelier-row">
                    <span class="atelier-titre">{{ atelier.titre }}</span>
                    <span class="atelier-date">{{ formaterDate(atelier.date_debut) }}</span>
                    <span class="badge-statut" :class="'statut-' + atelier.statut">
                        {{ formaterStatut(atelier.statut) }}
                    </span>
                    <button class="btn-modifier" @click="modifierAtelier(atelier.id)">Modifier</button>
                    <button v-if="atelier.statut === 'en_attente'" class="btn-modifier" @click="valider(atelier.id)">Valider</button>
                    <button class="btn-modifier" @click="supprimer(atelier.id)">Supprimer</button>
                </div>
            </div>

        </main>
    </div>

    <script src="js/ateliers.js"></script>
</body>

</html>
