<?php
require "inclus/espace.php";
$titrePage = "Forums";
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
                <h1>Forums</h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="deux-colonnes">

                <section>
                    <div v-if="sujet" class="carte section">
                        <button class="button-draft" @click="sujet = null">← Tous les sujets</button>
                        <h2>{{ sujet.titre }}</h2>

                        <ul class="liste-simple">
                            <li v-for="m in messages" :key="m.id">
                                <strong>{{ m.auteur }}</strong>
                                <span class="discret"> · {{ formaterDateHeure(m.created_at) }}</span>
                                <p>{{ m.contenu }}</p>

                                <form v-if="signalementOuvert === m.id" @submit.prevent="signaler(m.id)">
                                    <div class="form-group">
                                        <label :for="'motif-' + m.id">Motif du signalement</label>
                                        <input :id="'motif-' + m.id" type="text" v-model="motif" required>
                                    </div>
                                    <button type="submit" class="button-submit">Envoyer le signalement</button>
                                </form>
                                <button v-else class="button-draft" @click="signalementOuvert = m.id">Signaler</button>
                            </li>
                        </ul>

                        <p v-if="message" class="succes">{{ message }}</p>

                        <form @submit.prevent="repondre">
                            <div class="form-group">
                                <label for="reponse">Votre réponse</label>
                                <textarea id="reponse" v-model="reponse" required></textarea>
                            </div>
                            <button type="submit" class="button-submit">Répondre</button>
                        </form>
                    </div>

                    <div v-else>
                        <p v-if="sujets.length === 0" class="discret">Aucun sujet pour le moment.</p>
                        <div v-for="s in sujets" :key="s.id" class="carte carte-cliquable section" @click="ouvrirSujet(s)">
                            <span v-if="s.statut === 'en_attente'" class="badge badge-attente">En attente de modération</span>
                            <h3>{{ s.titre }}</h3>
                            <p class="discret">Par {{ s.auteur }} · {{ formaterDate(s.created_at) }}</p>
                        </div>
                    </div>
                </section>

                <div class="carte">
                    <h2>Nouveau sujet</h2>
                    <form @submit.prevent="creerSujet">
                        <div class="form-group">
                            <label for="titre">Titre</label>
                            <input id="titre" type="text" v-model="nouveauSujet.titre" required>
                        </div>
                        <div class="form-group">
                            <label for="contenu">Message</label>
                            <textarea id="contenu" v-model="nouveauSujet.contenu" required></textarea>
                        </div>
                        <button type="submit" class="button-submit">Publier</button>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/forums.js"></script>
</body>

</html>
