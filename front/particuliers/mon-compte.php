<?php
require "inclus/espace.php";
$titrePage = "Mon compte";
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
                <h1>Mon compte</h1>
            </div>

            <div class="deux-colonnes">

                <div class="carte">
                    <h2>Profil</h2>
                    <form @submit.prevent="enregistrerProfil">
                        <div class="form-group">
                            <label for="nom">Nom</label>
                            <input id="nom" type="text" v-model="profil.nom" required>
                        </div>
                        <div class="form-group">
                            <label for="prenom">Prénom</label>
                            <input id="prenom" type="text" v-model="profil.prenom" required>
                        </div>
                        <div class="form-group">
                            <label for="email">Email</label>
                            <input id="email" type="email" v-model="profil.email" required>
                        </div>
                        <div class="form-group">
                            <label for="langue">Langue d'affichage</label>
                            <select id="langue" v-model="profil.langue_preferee_id">
                                <option v-for="l in langues" :key="l.id" :value="l.id">{{ l.libelle }}</option>
                            </select>
                        </div>

                        <p v-if="erreurProfil" class="erreur">{{ erreurProfil }}</p>
                        <p v-if="messageProfil" class="succes">{{ messageProfil }}</p>
                        <button type="submit" class="button-submit">Enregistrer</button>
                    </form>
                </div>

                <div>
                    <div class="carte section">
                        <h2>Mot de passe</h2>
                        <form @submit.prevent="changerMotDePasse">
                            <div class="form-group">
                                <label for="ancien">Mot de passe actuel</label>
                                <input id="ancien" type="password" v-model="motsDePasse.ancien_mot_de_passe"
                                    autocomplete="current-password" required>
                            </div>
                            <div class="form-group">
                                <label for="nouveau">Nouveau mot de passe (8 caractères minimum)</label>
                                <input id="nouveau" type="password" v-model="motsDePasse.nouveau_mot_de_passe"
                                    autocomplete="new-password" minlength="8" required>
                            </div>

                            <p v-if="erreurMotDePasse" class="erreur">{{ erreurMotDePasse }}</p>
                            <p v-if="messageMotDePasse" class="succes">{{ messageMotDePasse }}</p>
                            <button type="submit" class="button-submit">Changer le mot de passe</button>
                        </form>
                    </div>

                    <div class="carte">
                        <h2>Tutoriel</h2>
                        <p>Revoir la présentation des fonctionnalités de votre espace.</p>
                        <a class="btn-primary" href="tableau-de-bord.php?tutoriel=1">Revoir le tutoriel</a>
                    </div>
                </div>

            </div>

        </main>
    </div>

    <script src="js/mon-compte.js"></script>
</body>

</html>
