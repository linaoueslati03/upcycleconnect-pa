<?php
require "inclus/espace.php";
$titrePage = "nav.prestations";
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
                <h1>{{ t('nav.prestations') }} (Factures & Contrats)</h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="mise-en-page-admin">
                <section class="carte">
                    <table class="tableau">
                        <thead>
                            <tr>
                                <th>{{ t('commun.titre') }}</th>
                                <th>{{ t('commun.categorie') }}</th>
                                <th>{{ t('commun.tarif') }}</th>
                                <th>{{ t('commun.statut') }}</th>
                                <th>{{ t('commun.actions') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="p in prestations" :key="p.id">
                                <td><strong>{{ p.titre }}</strong></td>
                                <td>{{ p.categorie }}</td>
                                <td>{{ formaterTarif(p.tarif) }}</td>
                                <td>{{ t('statut.' + p.statut) }}</td>
                                <td>
                                    <a class="lien-action" :href="lienPDF(p.id)">Télécharger PDF</a>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                    <p v-if="prestations.length === 0" style="margin-top: 15px;">Aucun contrat ou prestation pour le moment.</p>
                </section>
            </div>
        </main>
    </div>
    <script src="js/contrats.js"></script>
</body>
</html>
