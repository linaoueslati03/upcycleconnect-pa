<?php
require "inclus/espace.php";
$titrePage = "nav.mon_compte";
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

            <div class="deux-colonnes">

                <div class="carte">
                    <h2>{{ $t('compte.profil') }}</h2>
                    <form @submit.prevent="enregistrerProfil">
                        <div class="form-group">
                            <label for="nom">{{ $t('commun.nom') }}</label>
                            <input id="nom" type="text" v-model="profil.nom" required>
                        </div>
                        <div class="form-group">
                            <label for="prenom">{{ $t('commun.prenom') }}</label>
                            <input id="prenom" type="text" v-model="profil.prenom" required>
                        </div>
                        <div class="form-group">
                            <label for="email">{{ $t('commun.email') }}</label>
                            <input id="email" type="email" v-model="profil.email" required>
                        </div>
                        <div class="form-group">
                            <label for="langue">{{ $t('compte.langue_affichage') }}</label>
                            <select id="langue" v-model="profil.langue_preferee_id">
                                <option v-for="l in langues" :key="l.id" :value="l.id">{{ l.libelle }}</option>
                            </select>
                        </div>

                        <p v-if="erreurProfil" class="erreur">{{ erreurProfil }}</p>
                        <p v-if="messageProfil" class="succes">{{ messageProfil }}</p>
                        <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                    </form>
                </div>

                <div>
                    <div class="carte section">
                        <h2>{{ $t('commun.mot_de_passe') }}</h2>
                        <form @submit.prevent="changerMotDePasse">
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
                            <p v-if="messageMotDePasse" class="succes">{{ messageMotDePasse }}</p>
                            <button type="submit" class="button-submit">{{ $t('compte.changer_mot_de_passe') }}</button>
                        </form>
                    </div>

                    <div v-if="roleConnecte() === 'particulier'" class="carte">
                        <h2>{{ $t('compte.tutoriel') }}</h2>
                        <p>{{ $t('compte.tutoriel_explication') }}</p>
                        <a class="btn-primary" href="tableau-de-bord.php?tutoriel=1">{{ $t('compte.revoir_tutoriel') }}</a>
                    </div>

                    <div v-if="roleConnecte() === 'professionnel'" class="carte section">
                        <h2>Profil de l'entreprise</h2>
                        <div class="form-group">
                            <label for="siret">SIRET</label>
                            <input id="siret" type="text" v-model="entreprise.siret">
                        </div>
                        <div class="form-group">
                            <label for="adresse">Adresse de l'entreprise</label>
                            <input id="adresse" type="text" v-model="entreprise.adresse">
                        </div>
                    </div>

                    <div v-if="roleConnecte() === 'professionnel'" class="carte">
                        <h2>Abonnement</h2>
                        <p>Votre abonnement actuel : <strong>{{ entreprise.abonnement }}</strong></p>
                        <button v-if="entreprise.abonnement !== 'Premium'" class="btn-primary" @click="passerPremium">Passer Premium</button>
                    </div>
                </div>

            </div>

        </main>
    </div>

    <script src="js/mon-compte.js"></script>
</body>

</html>
