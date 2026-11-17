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
                <h1>{{ $t('nav.forums') }}</h1>
                <button class="btn-primary" @click="ouvrirNouveauSujet">+ {{ $t('forums.nouveau') }}</button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div>

                <section>
                    <div v-if="sujet" class="carte section">
                        <button class="button-draft" @click="sujet = null">← {{ $t('forums.tous_les_sujets') }}</button>
                        <h2>{{ sujet.titre }}</h2>

                        <ul class="liste-simple">
                            <li v-for="m in messages" :key="m.id">
                                <strong>{{ m.auteur }}</strong>
                                <span class="discret"> · {{ formaterDateHeure(m.created_at) }}</span>
                                <p>{{ m.contenu }}</p>

                                <button class="button-draft" @click="ouvrirSignalement(m.id)">{{ $t('forums.signaler') }}</button>
                            </li>
                        </ul>

                        <p v-if="message" class="succes">{{ message }}</p>

                        <button class="btn-primary" @click="ouvrirReponse">{{ $t('forums.repondre') }}</button>
                    </div>

                    <div v-else>
                        <p v-if="sujets.length === 0" class="discret">{{ $t('forums.aucun') }}</p>
                        <div v-for="s in sujets" :key="s.id" class="carte carte-cliquable section" @click="ouvrirSujet(s)">
                            <span v-if="s.statut === 'en_attente'" class="badge badge-attente">{{ $t('forums.en_attente') }}</span>
                            <h3>{{ s.titre }}</h3>
                            <p class="discret">{{ $t('forums.par', { auteur: s.auteur }) }} · {{ formaterDate(s.created_at) }}</p>
                        </div>
                    </div>
                </section>

            </div>

            <fenetre-modale :ouverte="nouveauSujetOuvert" :titre="$t('forums.nouveau')" @fermer="nouveauSujetOuvert = false">
                    <form @submit.prevent="creerSujet">
                        <div class="form-group">
                            <label for="titre">{{ $t('commun.titre') }}</label>
                            <input id="titre" type="text" v-model="nouveauSujet.titre" required>
                        </div>
                        <div class="form-group">
                            <label for="contenu">{{ $t('forums.message') }}</label>
                            <textarea id="contenu" v-model="nouveauSujet.contenu" required></textarea>
                        </div>
                        <p v-if="erreur" class="erreur">{{ erreur }}</p>
                        <div class="form-buttons">
                            <button type="submit" class="button-submit">{{ $t('commun.publier') }}</button>
                            <button type="button" class="button-draft" @click="nouveauSujetOuvert = false">{{ $t('commun.annuler') }}</button>
                        </div>
                    </form>
            </fenetre-modale>

            <fenetre-modale :ouverte="reponseOuverte" :titre="$t('forums.repondre') + (sujet ? ' · ' + sujet.titre : '')" @fermer="reponseOuverte = false">
                <form @submit.prevent="repondre">
                    <div class="form-group">
                        <label for="reponse">{{ $t('forums.votre_reponse') }}</label>
                        <textarea id="reponse" v-model="reponse" required></textarea>
                    </div>
                    <p v-if="erreur" class="erreur">{{ erreur }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('forums.repondre') }}</button>
                        <button type="button" class="button-draft" @click="reponseOuverte = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

            <fenetre-modale :ouverte="signalementVisible" :titre="$t('forums.signaler')" @fermer="signalementVisible = false">
                <form @submit.prevent="signaler(signalementOuvert)">
                    <div class="form-group">
                        <label for="motif-signalement">{{ $t('forums.motif') }}</label>
                        <input id="motif-signalement" type="text" v-model="motif" required>
                    </div>
                    <p v-if="erreur" class="erreur">{{ erreur }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('forums.envoyer_signalement') }}</button>
                        <button type="button" class="button-draft" @click="signalementVisible = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

        </main>
    </div>

    <script src="<?= ressource("js/forums.js") ?>"></script>
</body>

</html>
