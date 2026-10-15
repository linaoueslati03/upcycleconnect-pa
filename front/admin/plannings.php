<?php
require "inclus/espace.php";
$titrePage = "nav.plannings";
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
                <h1>{{ $t('admin.plannings.titre') }}</h1>
            </div>

            <div class="filtres">
                <button v-for="filtre in filtres" :key="filtre.valeur"
                    class="pill" :class="{ active: participation === filtre.valeur }"
                    @click="participation = filtre.valeur">
                    {{ $t(filtre.nom) }}
                </button>
                <button class="pill" :class="{ active: aVenirSeulement }" @click="aVenirSeulement = !aVenirSeulement">
                    {{ $t('admin.plannings.a_venir') }}
                </button>
                <input type="search" v-model="recherche" :placeholder="$t('admin.plannings.rechercher_personne')"
                    :aria-label="$t('admin.plannings.rechercher_personne')">
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="!erreur && entreesFiltrees.length === 0">{{ $t('admin.plannings.aucune') }}</p>

            <section v-if="entreesFiltrees.length > 0" class="carte">
                <table class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('commun.date_debut') }}</th>
                            <th>{{ $t('admin.plannings.personne') }}</th>
                            <th>{{ $t('commun.role') }}</th>
                            <th>{{ $t('admin.plannings.participation') }}</th>
                            <th>{{ $t('commun.type') }}</th>
                            <th>{{ $t('commun.titre') }}</th>
                            <th>{{ $t('commun.lieu') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="e in entreesFiltrees" :key="e.participation + e.type + e.offre_id + '-' + e.utilisateur_id">
                            <td>{{ formaterDateHeure(e.date_debut) }}</td>
                            <td>{{ e.personne }}</td>
                            <td>{{ $t('role.' + e.role_id) }}</td>
                            <td>{{ $t('admin.plannings.' + e.participation) }}</td>
                            <td>{{ $t('offre.' + e.type) }}</td>
                            <td>{{ e.titre }}</td>
                            <td>{{ e.lieu || '—' }}</td>
                        </tr>
                    </tbody>
                </table>
            </section>

        </main>
    </div>

    <script src="js/plannings.js"></script>
</body>

</html>
