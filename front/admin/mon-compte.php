<?php
require "inclus/espace.php";
$titrePage = "nav.mon_compte";
$scriptsEspace = ["../espace/js/outils.js"];
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
                <h1>{{ $t('nav.mon_compte') }}</h1>
            </div>

            <p v-if="motDePasseObligatoire" class="erreur">{{ $t('compte.changement_obligatoire') }}</p>

            <div class="mise-en-page-admin">

                <div class="carte carte-formulaire">
                    <h2>{{ $t('compte.profil') }}</h2>
                    <dl class="fiche-profil">
                        <dt>{{ $t('commun.nom') }}</dt><dd>{{ profil.nom || '—' }}</dd>
                        <dt>{{ $t('commun.prenom') }}</dt><dd>{{ profil.prenom || '—' }}</dd>
                        <dt>{{ $t('commun.email') }}</dt><dd>{{ profil.email || '—' }}</dd>
                        <dt>{{ $t('compte.langue_affichage') }}</dt><dd>{{ libelleLangue || '—' }}</dd>
                    </dl>
                    <p v-if="erreurProfil && !profilOuvert" class="erreur">{{ erreurProfil }}</p>
                    <p v-if="messageProfil" class="succes">{{ messageProfil }}</p>
                    <button class="btn-primary" @click="ouvrirProfil">{{ $t('commun.modifier') }}</button>
                </div>

                <div class="carte carte-formulaire">
                    <h2>{{ $t('commun.mot_de_passe') }}</h2>
                    <p class="discret">{{ $t('commun.regle_mot_de_passe') }}</p>
                    <p v-if="messageMotDePasse" class="succes">{{ messageMotDePasse }}</p>
                    <button class="btn-primary" @click="ouvrirMotDePasse">{{ $t('compte.changer_mot_de_passe') }}</button>
                </div>

                <div class="carte carte-formulaire">
                    <h2>{{ $t('compte.tutoriel') }}</h2>
                    <p>{{ $t('compte.tutoriel_explication') }}</p>
                    <a class="btn-primary" :href="lienTutoriel()">{{ $t('compte.revoir_tutoriel') }}</a>
                </div>

                <div class="carte carte-formulaire">
                    <h2>Notifications</h2>
                    <p>Active les notifications pour recevoir les informations importantes d'UpcycleConnect.</p>
                    <button id="activer-notifications" type="button" class="button-submit">
                        Activer les notifications
                    </button>
                </div>

            </div>

            <fenetre-modale :ouverte="profilOuvert" :titre="$t('compte.profil')" @fermer="profilOuvert = false">
                <form @submit.prevent="enregistrerProfil">
                    <div class="form-group">
                        <label for="nom">{{ $t('commun.nom') }}</label>
                        <input id="nom" type="text" v-model="brouillonProfil.nom" required>
                    </div>
                    <div class="form-group">
                        <label for="prenom">{{ $t('commun.prenom') }}</label>
                        <input id="prenom" type="text" v-model="brouillonProfil.prenom" required>
                    </div>
                    <div class="form-group">
                        <label for="email">{{ $t('commun.email') }}</label>
                        <input id="email" type="email" v-model="brouillonProfil.email" required>
                    </div>
                    <div class="form-group">
                        <label for="langue">{{ $t('compte.langue_affichage') }}</label>
                        <select id="langue" v-model="brouillonProfil.langue_preferee_id">
                            <option v-for="l in langues" :key="l.id" :value="l.id">{{ l.libelle }}</option>
                        </select>
                    </div>
                    <p v-if="erreurProfil" class="erreur">{{ erreurProfil }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                        <button type="button" class="button-draft" @click="profilOuvert = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

            <fenetre-modale :ouverte="motDePasseOuvert" :titre="$t('compte.changer_mot_de_passe')" @fermer="motDePasseOuvert = false">
                <form @submit.prevent="changerMotDePasse">
                    <p v-if="motDePasseObligatoire" class="erreur">{{ $t('compte.changement_obligatoire') }}</p>
                    <div class="form-group">
                        <label for="ancien">{{ $t('compte.mot_de_passe_actuel') }}</label>
                        <input id="ancien" type="password" v-model="motsDePasse.ancien_mot_de_passe"
                            autocomplete="current-password" required>
                    </div>
                    <div class="form-group">
                        <label for="nouveau">{{ $t('compte.nouveau_mot_de_passe') }}</label>
                        <input id="nouveau" type="password" v-model="motsDePasse.nouveau_mot_de_passe"
                            autocomplete="new-password" minlength="12" maxlength="72" required
                            aria-describedby="regle-mot-de-passe">
                        <small id="regle-mot-de-passe">{{ $t('commun.regle_mot_de_passe') }}</small>
                    </div>
                    <p v-if="erreurMotDePasse" class="erreur">{{ erreurMotDePasse }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('compte.changer_mot_de_passe') }}</button>
                        <button type="button" class="button-draft" @click="motDePasseOuvert = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

        </main>
    </div>

    <script src="<?= ressource("../espace/js/mon-compte.js") ?>"></script>
</body>

</html>
