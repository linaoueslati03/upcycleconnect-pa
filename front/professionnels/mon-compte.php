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
                <h1>{{ t('nav.mon_compte') }}</h1>
            </div>

            <div class="mise-en-page-admin">
                <div class="carte carte-formulaire">
                    <h2>Profil de l'entreprise</h2>
                    <form @submit.prevent="enregistrerProfil">
                        <div class="form-group">
                            <label for="nom">{{ t('commun.nom') }} (Entreprise)</label>
                            <input id="nom" type="text" v-model="profil.nom" required>
                        </div>

                        <div class="form-group">
                            <label for="siret">SIRET</label>
                            <input id="siret" type="text" v-model="profil.siret">
                        </div>
                        
                        <div class="form-group">
                            <label for="adresse">Adresse de l'entreprise</label>
                            <input id="adresse" type="text" v-model="profil.adresse">
                        </div>

                        <div class="form-buttons">
                            <button type="submit" class="button-submit">{{ t('commun.enregistrer') }}</button>
                        </div>
                        <p v-if="succesProfil" style="color: green; margin-top: 10px;">Profil mis à jour avec succès.</p>
                        <p v-if="erreurProfil" class="erreur">{{ erreurProfil }}</p>
                    </form>
                </div>

                <div class="carte" style="margin-top: 20px;">
                    <h2>Abonnement</h2>
                    <p>Votre abonnement actuel : <strong>{{ profil.abonnement || 'Freemium' }}</strong></p>
                    <button class="btn-primary" @click="passerPremium" v-if="profil.abonnement !== 'Premium'">Passer Premium</button>
                </div>
            </div>
        </main>
    </div>
    <script src="js/mon-compte.js"></script>
</body>
</html>
