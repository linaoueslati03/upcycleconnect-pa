<?php
require "inclus/espace.php";
$titrePage = "Mes projets";
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
                <h1>Mes projets d'upcycling</h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="deux-colonnes">

                <section>
                    <p v-if="projets.length === 0" class="discret">Vous n'avez pas encore de projet.</p>
                    <div class="grille section">
                        <div v-for="p in projets" :key="p.id" class="carte carte-cliquable"
                            :class="{ 'carte-selectionnee': projet && projet.id === p.id }" @click="ouvrirProjet(p.id)">
                            <span class="badge" :class="{ 'badge-attente': !p.partage_public }">
                                {{ p.partage_public ? "Partagé" : "Privé" }}
                            </span>
                            <h3>{{ p.titre }}</h3>
                            <p class="discret">Créé le {{ formaterDate(p.date_creation) }}</p>
                        </div>
                    </div>

                    <div v-if="projet" class="carte">
                        <h2>{{ projet.titre }}</h2>
                        <p>{{ projet.description }}</p>

                        <div class="form-buttons section">
                            <button class="button-submit" @click="basculerPartage">
                                {{ projet.partage_public ? "Ne plus partager" : "Partager avec la communauté" }}
                            </button>
                            <button class="button-draft" @click="supprimerProjet">Supprimer le projet</button>
                        </div>

                        <h3>Étapes de transformation</h3>
                        <p v-if="projet.etapes.length === 0" class="discret">Aucune étape documentée.</p>
                        <ol>
                            <li v-for="e in projet.etapes" :key="e.id">
                                {{ e.description }}
                                <span class="discret"> · {{ formaterDate(e.date) }}</span>
                                <a v-if="e.photo_url" :href="e.photo_url" target="_blank" rel="noopener"> (photo)</a>
                            </li>
                        </ol>

                        <form @submit.prevent="ajouterEtape">
                            <div class="form-group">
                                <label for="etape">Nouvelle étape</label>
                                <textarea id="etape" v-model="nouvelleEtape.description" required></textarea>
                            </div>
                            <div class="form-group">
                                <label for="photo">Lien vers une photo (facultatif)</label>
                                <input id="photo" type="url" v-model="nouvelleEtape.photo_url">
                            </div>
                            <button type="submit" class="button-submit">Ajouter l'étape</button>
                        </form>
                    </div>
                </section>

                <div class="carte">
                    <h2>Nouveau projet</h2>
                    <form @submit.prevent="creerProjet">
                        <div class="form-group">
                            <label for="titre">Titre</label>
                            <input id="titre" type="text" v-model="formulaire.titre" required>
                        </div>
                        <div class="form-group">
                            <label for="description">Description</label>
                            <textarea id="description" v-model="formulaire.description"></textarea>
                        </div>
                        <div class="form-group">
                            <label><input type="checkbox" v-model="formulaire.partage_public"> Partager avec la communauté</label>
                        </div>
                        <button type="submit" class="button-submit">Créer le projet</button>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/projets.js"></script>
</body>

</html>
