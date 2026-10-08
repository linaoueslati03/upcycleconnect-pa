<?php
// Liens du menu : fichier => [libellé, icône]
$liens = [
    "planning.php"   => ["Mon planning", "people-svgrepo-com.svg"],
    "evenements.php" => ["Événements", "party-horn-svgrepo-com.svg"],
    "ateliers.php"   => ["Ateliers", "pencil-square-svgrepo-com.svg"],
    "articles.php"   => ["Articles", "news-svgrepo-com.svg"],
];
?>
<aside>

    <a id="toggle-button" href="planning.php" aria-label="Accueil">
        <img src="../assets/img/LOGO-pa.png" alt="Logo">
    </a>

    <ul>
        <?php foreach ($liens as $fichier => [$libelle, $icone]): ?>
            <li>
                <a class="nav-button" href="<?= $fichier ?>">
                    <img src="../assets/img/<?= $icone ?>" alt="<?= $libelle ?>">
                    <span><?= $libelle ?></span>
                </a>
            </li>
        <?php endforeach; ?>
    </ul>

    <button id="logout-button">
        <img src="../assets/img/logout-svgrepo-com.svg" alt="Déconnexion">
    </button>

</aside>
