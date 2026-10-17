<?php
require "inclus/espace.php";
$titrePage = "nav.categories";
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
                <h1>{{ $t('admin.categories.titre') }}</h1>
            </div>

            <p class="discret">{{ $t('admin.categories.explication') }}</p>
            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <section v-for="(groupe, type) in groupes" :key="type" class="carte section-categories">
                <h2>{{ $t('admin.categories.' + type) }}</h2>

                <form class="ajout-categorie" @submit.prevent="ajouter(type)">
                    <input type="text" v-model="groupe.nouvelle" :placeholder="$t('admin.categories.nouvelle')"
                        :aria-label="$t('admin.categories.nouvelle')" required>
                    <button type="submit" class="button-submit">{{ $t('admin.categories.ajouter') }}</button>
                </form>

                <table class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('commun.nom') }}</th>
                            <th>{{ $t('commun.statut') }}</th>
                            <th>{{ $t('commun.actions') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="c in groupe.liste" :key="c.id">
                            <td>
                                <input v-if="edition.type === type && edition.id === c.id" type="text" v-model="edition.libelle"
                                    :aria-label="$t('commun.nom')">
                                <template v-else>{{ c.libelle }}</template>
                            </td>
                            <td :class="{ 'statut-suspendu': !c.actif }">
                                {{ c.actif ? $t('admin.categories.active') : $t('admin.categories.inactive') }}
                            </td>
                            <td>
                                <template v-if="edition.type === type && edition.id === c.id">
                                    <button class="lien-action" @click="enregistrer(type, c, edition.libelle, c.actif)">{{ $t('commun.enregistrer') }}</button>
                                    <button class="lien-action" @click="annulerEdition">{{ $t('commun.annuler') }}</button>
                                </template>
                                <template v-else>
                                    <button class="lien-action" @click="commencerEdition(type, c)">{{ $t('admin.categories.renommer') }}</button>
                                    <button class="lien-action" @click="enregistrer(type, c, c.libelle, !c.actif)">
                                        {{ c.actif ? $t('admin.categories.desactiver') : $t('admin.categories.activer') }}
                                    </button>
                                </template>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </section>

        </main>
    </div>

    <script src="js/categories.js"></script>
</body>

</html>
