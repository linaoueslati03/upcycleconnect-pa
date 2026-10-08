<?php
require "inclus/espace.php";
$titrePage = "Prestations";
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
                <h1>Gestion des prestations</h1>
                <button class="btn-primary" @click="reinitialiserFormulaire">+ Nouvelle prestation</button>
            </div>

            <p v-if="erreurListe" class="erreur">{{ erreurListe }}</p>

            <div class="mise-en-page-admin">

                <section class="carte">
                    <table class="tableau">
                        <thead>
                            <tr>
                                <th>Titre</th>
                                <th>Catégorie</th>
                                <th>Tarif</th>
                                <th>Statut</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="p in prestations" :key="p.id">
                                <td><strong>{{ p.titre }}</strong></td>
                                <td>{{ p.categorie }}</td>
                                <td>{{ formaterTarif(p.tarif) }}</td>
                                <td>{{ p.statut === "publie" ? "Publiée" : "Brouillon" }}</td>
                                <td>
                                    <button class="lien-action" @click="modifier(p)">Modifier</button>
                                    <a class="lien-action" :href="lienPDF(p.id)">PDF</a>
                                    <button class="lien-action lien-supprimer" @click="supprimer(p.id)">Supprimer</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </section>

                <div class="carte carte-formulaire">
                    <h2>{{ formulaire.id ? "Modifier la prestation" : "Nouvelle prestation" }}</h2>

                    <form @submit.prevent="enregistrer">
                        <div class="form-group">
                            <label for="titre">Titre</label>
                            <input id="titre" type="text" v-model="formulaire.titre" required>
                        </div>

                        <div class="form-group">
                            <label for="categorie">Catégorie</label>
                            <input id="categorie" type="text" v-model="formulaire.categorie" required>
                        </div>

                        <div class="form-group">
                            <label for="tarif">Tarif (€)</label>
                            <input id="tarif" type="number" min="0" step="0.01" v-model.number="formulaire.tarif" required>
                        </div>

                        <div class="form-group">
                            <label for="statut">Statut</label>
                            <select id="statut" v-model="formulaire.statut">
                                <option value="brouillon">Brouillon</option>
                                <option value="publie">Publiée</option>
                            </select>
                        </div>

                        <p v-if="erreur" class="erreur">{{ erreur }}</p>

                        <div class="form-buttons">
                            <button type="submit" class="button-submit">Enregistrer</button>
                            <button type="button" class="button-draft" @click="reinitialiserFormulaire">Annuler</button>
                        </div>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/prestations.js"></script>
</body>

</html>
