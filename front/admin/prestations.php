<!DOCTYPE html>
<html lang="fr">
    <head>
        <meta charset="UTF-8">
        <title>UpcycleConnect - Prestations</title>
        <link rel="stylesheet" href="../assets/css/style.css">
        <link rel="stylesheet" href="../assets/css/prestations.css">
    </head>
    <body>
        <aside>
            <button id="toggle-button">
                <img src="../assets/img/LOGO-pa.png" alt="Logo UpcycleConnect">
            </button>

            <ul>
                <li>
                    <a class="nav-button" href="utilisateurs.php">
                        <img src="../assets/img/people-svgrepo-com.svg" alt="Utilisateurs">
                        <span>Utilisateurs</span>
                    </a>
                </li>
                <li>
                    <a class="nav-button" href="prestations.php">
                        <img src="../assets/img/book-svgrepo-com.svg" alt="Prestations">
                        <span>Prestations</span>
                    </a>
                </li>
            </ul>

            <button id="logout-button" onclick="deconnecter()">
                <img src="../assets/img/logout-svgrepo-com.svg" alt="Déconnexion">
            </button>
        </aside>

        <main class="main-content">
            <header class="page-header">
                <h1>GESTION DES PRESTATIONS</h1>
            </header>

            <div style="display: flex; gap: 40px; align-items: flex-start;">
                
                <section class="content-left" style="flex: 7;">
                    <div class="table-container">
                        <table class="presta-table">
                            <thead>
                                <tr>
                                    <th>Titre</th>
                                    <th>Catégorie</th>
                                    <th>Tarif</th>
                                    <th>Statut</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody id="prestations-list">
                                <!-- JS -->
                            </tbody>
                        </table>
                    </div>
                </section>

                <div class="content-right" style="flex: 3;">
                    <div class="table-container">
                        <h3 style="color: #2E7D32; margin-top: 0;" id="form-title">NOUVELLE OFFRE</h3>
                        <form id="prestation-form">
                            <input type="hidden" id="presta-id">
                            
                            <label>Titre</label>
                            <input type="text" id="presta-titre" required style="width: 100%; margin-bottom: 10px;">
                            
                            <label>Catégorie</label>
                            <input type="text" id="presta-categorie" required style="width: 100%; margin-bottom: 10px;">
                            
                            <label>Tarif (€)</label>
                            <input type="number" id="presta-tarif" required style="width: 100%; margin-bottom: 10px;">
                            
                            <label>Statut</label>
                            <select id="presta-statut" style="width: 100%; margin-bottom: 20px;">
                                <option value="Brouillon">Brouillon</option>
                                <option value="Publiée">Publiée</option>
                            </select>
                            
                            <button type="submit" class="btn-add" style="width: 100%; margin-bottom: 10px;">Enregistrer</button>
                            <button type="button" onclick="reinitialiserFormulaire()" style="width: 100%;">Annuler</button>
                        </form>
                    </div>
                </div>

            </div>
        </main>

        <script src="../assets/js/config.js"></script>
        <script src="../assets/js/api.js"></script>
        <script>exigerConnexion("administrateur");</script>
        <script src="js/prestations.js"></script>
    </body>
</html>