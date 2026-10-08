<?php
// Menu latéral commun aux espaces connectés.
// $liens (défini dans inclus/espace.php de chaque espace) : fichier => [libellé, icône]
?>
<aside>

    <a id="toggle-button" href="<?= array_key_first($liens) ?>" aria-label="Accueil">
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

    <button id="logout-button" onclick="deconnecter()">
        <img src="../assets/img/logout-svgrepo-com.svg" alt="Déconnexion">
    </button>

</aside>
