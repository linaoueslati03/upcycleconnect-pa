<?php
require "inclus/espace.php";
$titrePage = "Mes annonces";
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
                <h1>Mes annonces</h1>
                <button class="btn-primary" @click="reinitialiserFormulaire">+ Nouvelle annonce</button>
            </div>

            <p v-if="erreurListe" class="erreur">{{ erreurListe }}</p>

            <div class="deux-colonnes">

                <section>
                    <p v-if="annonces.length === 0" class="discret">Vous n'avez pas encore publié d'annonce.</p>
                    <div v-for="a in annonces" :key="a.id" class="carte section">
                        <span class="badge" :class="{ 'badge-attente': a.statut !== 'en_ligne' }">{{ libelleStatut(a.statut) }}</span>
                        <h3>{{ a.titre }}</h3>
                        <p class="discret">
                            {{ a.type === "don" ? "Don" : "Vente · " + formaterPrix(a.prix) }} · {{ a.localisation }}
                        </p>
                        <div class="form-buttons">
                            <button class="button-draft" @click="modifier(a)">Modifier</button>
                            <button v-if="a.statut !== 'cedee'" class="button-submit" @click="marquerCedee(a)">Marquer comme cédée</button>
                            <button class="button-draft" @click="supprimer(a.id)">Supprimer</button>
                        </div>
                    </div>
                </section>

                <div class="carte">
                    <h2>{{ formulaire.id ? "Modifier l'annonce" : "Déposer une annonce" }}</h2>

                    <form @submit.prevent="enregistrer">
                        <div class="form-group">
                            <label for="titre">Titre</label>
                            <input id="titre" type="text" v-model="formulaire.titre" required>
                        </div>

                        <div class="form-group">
                            <label for="description">Description</label>
                            <textarea id="description" v-model="formulaire.description"></textarea>
                        </div>

                        <div class="form-group">
                            <label for="type">Type</label>
                            <select id="type" v-model="formulaire.type">
                                <option value="don">Don</option>
                                <option value="vente">Vente</option>
                            </select>
                        </div>

                        <div v-if="formulaire.type === 'vente'" class="form-group">
                            <label for="prix">Prix (€)</label>
                            <input id="prix" type="number" min="0.01" step="0.01" v-model.number="formulaire.prix" required>
                        </div>

                        <div class="form-group">
                            <label for="categorie">Catégorie de matériau</label>
                            <select id="categorie" v-model="formulaire.categorie_id">
                                <option :value="null">Non précisée</option>
                                <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.code }}</option>
                            </select>
                        </div>

                        <div class="form-group">
                            <label for="localisation">Localisation</label>
                            <input id="localisation" type="text" v-model="formulaire.localisation">
                        </div>

                        <div v-if="formulaire.id" class="form-group">
                            <label for="statut">Statut</label>
                            <select id="statut" v-model="formulaire.statut">
                                <option v-for="(libelle, statut) in libellesStatut" :key="statut" :value="statut">{{ libelle }}</option>
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

    <script src="js/mes-annonces.js"></script>
</body>

</html>
