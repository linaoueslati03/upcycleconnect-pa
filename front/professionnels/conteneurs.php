<?php
require "inclus/espace.php";
$titrePage = "nav.depot_conteneur";
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
                <h1>{{ t('nav.depot_conteneur') }} (Récupération)</h1>
            </div>

            <div class="mise-en-page-admin">
                <section class="carte">
                    <h2>Scanner un objet en conteneur</h2>
                    <p>Entrez le code de l'objet pour vérifier sa disponibilité et le récupérer.</p>
                    
                    <form @submit.prevent="verifierCode">
                        <div class="form-group">
                            <label for="code_objet">Code de l'objet (Scan ou manuel)</label>
                            <input id="code_objet" type="text" v-model="codeObjet" required placeholder="Ex: OBJ-12345">
                        </div>
                        <button type="submit" class="btn-primary">Vérifier le code</button>
                    </form>

                    <p v-if="erreur" class="erreur" style="margin-top: 15px;">{{ erreur }}</p>
                </section>

                <section v-if="objetTrouve" class="carte" style="margin-top: 20px; background-color: #f9f9f9; border-left: 4px solid #4CAF50;">
                    <h2>Objet trouvé</h2>
                    <p><strong>Type :</strong> {{ objetTrouve.type }}</p>
                    <p><strong>Description :</strong> {{ objetTrouve.description }}</p>
                    <p><strong>Conteneur / Box :</strong> {{ objetTrouve.conteneur }}</p>
                    
                    <button @click="confirmerRecuperation" class="btn-primary" style="margin-top: 15px; background-color: #4CAF50;">Confirmer la récupération</button>
                    <p v-if="succes" style="color: green; margin-top: 10px; font-weight: bold;">{{ succes }}</p>
                </section>
            </div>
        </main>
    </div>
    <script src="js/conteneurs.js"></script>
</body>
</html>
