<?php
require "inclus/espace.php";
$titrePage = "nav.score";
?>
<!DOCTYPE html>
<html lang="fr">
<head>
    <?php include "../inclus/entete.php"; ?>
    <style>
        .graph-placeholder {
            background-color: #f1f8e9;
            border: 2px dashed #8bc34a;
            border-radius: 8px;
            height: 200px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #558b2f;
            font-weight: bold;
            margin-bottom: 20px;
        }
        .lock-banner {
            background-color: #ffebee;
            border-left: 4px solid #f44336;
            padding: 20px;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <?php include "../inclus/navigation.php"; ?>

    <div id="app">
        <main>
            <div class="page-header">
                <h1>Statistiques & Impact (Premium)</h1>
            </div>

            <div class="mise-en-page-admin">
                
                <section v-if="profil.abonnement !== 'Premium'" class="carte lock-banner">
                    <h2 style="color: #d32f2f;">Fonctionnalité Premium</h2>
                    <p>Votre compte actuel est en mode <strong>Freemium</strong>. Les statistiques d'impact écologique, l'état des stocks en temps réel et les alertes de collecte priorisées sont réservés aux abonnés Premium.</p>
                    <button class="btn-primary" @click="passerPremium" style="background-color: #d32f2f; margin-top: 10px;">Souscrire à l'offre Premium</button>
                </section>

                <div v-else>
                    <section class="carte">
                        <h2>Votre Impact Écologique</h2>
                        <p>Total de CO2 évité grâce à vos récupérations ce mois-ci : <strong>450 kg</strong></p>
                        <div class="graph-placeholder">
                            [Graphique d'évolution des émissions évitées]
                        </div>
                    </section>

                    <section class="carte" style="margin-top: 20px;">
                        <h2>Alertes priorisées de collecte</h2>
                        <p>Configurez vos critères pour recevoir une notification dès qu'un matériau spécifique est déposé.</p>
                        
                        <form @submit.prevent="enregistrerAlerte">
                            <div class="form-group">
                                <label for="materiau">Type de matériau recherché</label>
                                <select id="materiau" v-model="filtreAlerte.materiau" required>
                                    <option value="">-- Sélectionnez --</option>
                                    <option value="bois">Bois massif</option>
                                    <option value="metal">Métal / Acier</option>
                                    <option value="verre">Verre</option>
                                    <option value="textile">Textile</option>
                                </select>
                            </div>
                            <div class="form-group">
                                <label for="localisation">Secteur / Agence</label>
                                <input id="localisation" type="text" v-model="filtreAlerte.localisation" placeholder="Paris 11ème" required>
                            </div>
                            <button type="submit" class="btn-primary">M'avertir</button>
                            <p v-if="messageAlerte" style="color: green; margin-top: 10px; font-weight: bold;">{{ messageAlerte }}</p>
                        </form>
                    </section>
                </div>

            </div>
        </main>
    </div>
    <script src="js/statistiques.js"></script>
</body>
</html>
