<?php
require "inclus/espace.php";
$titrePage = "admin.annonces.titre";
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
                <h1>{{ t('admin.annonces.titre') }}</h1>
            </div>

            <div class="filtres">
                <button v-for="filtre in filtres" :key="filtre.valeur"
                    class="pill" :class="{ active: statutSelectionne === filtre.valeur }"
                    @click="statutSelectionne = filtre.valeur; chargerAnnonces()">
                    {{ t(filtre.nom) }}
                </button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="!erreur && annonces.length === 0">{{ t('admin.annonces.aucune') }}</p>

            <section v-if="annonces.length > 0" class="carte">
                <table class="tableau">
                    <thead>
                        <tr>
                            <th>{{ t('commun.titre') }}</th>
                            <th>{{ t('admin.annonces.auteur') }}</th>
                            <th>{{ t('commun.type') }}</th>
                            <th>{{ t('commun.localisation') }}</th>
                            <th>{{ t('commun.statut') }}</th>
                            <th>{{ t('commun.actions') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="a in annonces" :key="a.id">
                            <td>
                                <strong>{{ a.titre }}</strong><br>
                                <span>{{ a.description }}</span>
                            </td>
                            <td>{{ a.auteur }}</td>
                            <td>
                                {{ t('annonce.type.' + a.type) }}
                                <span v-if="a.type === 'vente'"> · {{ formaterPrix(a.prix) }}</span>
                            </td>
                            <td>{{ a.localisation }}</td>
                            <td>{{ t('annonce.' + a.statut) }}</td>
                            <td>
                                <template v-if="a.statut === 'en_attente'">
                                    <button class="lien-action" @click="decider(a.id, 'valider')">{{ t('commun.valider') }}</button>
                                    <button class="lien-action lien-supprimer" @click="decider(a.id, 'refuser')">{{ t('commun.refuser') }}</button>
                                </template>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </section>

        </main>
    </div>

    <script src="js/annonces.js"></script>
</body>

</html>
