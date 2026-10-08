<?php
require "inclus/espace.php";
$titrePage = "nav.forums";
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
                <h1>{{ t('nav.forums') }}</h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="deux-colonnes">

                <section>
                    <div v-if="sujet" class="carte section">
                        <button class="button-draft" @click="sujet = null">← {{ t('forums.tous_les_sujets') }}</button>
                        <h2>{{ sujet.titre }}</h2>

                        <ul class="liste-simple">
                            <li v-for="m in messages" :key="m.id">
                                <strong>{{ m.auteur }}</strong>
                                <span class="discret"> · {{ formaterDateHeure(m.created_at) }}</span>
                                <p>{{ m.contenu }}</p>

                                <form v-if="signalementOuvert === m.id" @submit.prevent="signaler(m.id)">
                                    <div class="form-group">
                                        <label :for="'motif-' + m.id">{{ t('forums.motif') }}</label>
                                        <input :id="'motif-' + m.id" type="text" v-model="motif" required>
                                    </div>
                                    <button type="submit" class="button-submit">{{ t('forums.envoyer_signalement') }}</button>
                                </form>
                                <button v-else class="button-draft" @click="signalementOuvert = m.id">{{ t('forums.signaler') }}</button>
                            </li>
                        </ul>

                        <p v-if="message" class="succes">{{ message }}</p>

                        <form @submit.prevent="repondre">
                            <div class="form-group">
                                <label for="reponse">{{ t('forums.votre_reponse') }}</label>
                                <textarea id="reponse" v-model="reponse" required></textarea>
                            </div>
                            <button type="submit" class="button-submit">{{ t('forums.repondre') }}</button>
                        </form>
                    </div>

                    <div v-else>
                        <p v-if="sujets.length === 0" class="discret">{{ t('forums.aucun') }}</p>
                        <div v-for="s in sujets" :key="s.id" class="carte carte-cliquable section" @click="ouvrirSujet(s)">
                            <span v-if="s.statut === 'en_attente'" class="badge badge-attente">{{ t('forums.en_attente') }}</span>
                            <h3>{{ s.titre }}</h3>
                            <p class="discret">{{ t('forums.par', { auteur: s.auteur }) }} · {{ formaterDate(s.created_at) }}</p>
                        </div>
                    </div>
                </section>

                <div class="carte">
                    <h2>{{ t('forums.nouveau') }}</h2>
                    <form @submit.prevent="creerSujet">
                        <div class="form-group">
                            <label for="titre">{{ t('commun.titre') }}</label>
                            <input id="titre" type="text" v-model="nouveauSujet.titre" required>
                        </div>
                        <div class="form-group">
                            <label for="contenu">{{ t('forums.message') }}</label>
                            <textarea id="contenu" v-model="nouveauSujet.contenu" required></textarea>
                        </div>
                        <button type="submit" class="button-submit">{{ t('commun.publier') }}</button>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/forums.js"></script>
</body>

</html>
