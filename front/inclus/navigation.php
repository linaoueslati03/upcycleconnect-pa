<?php
// Menu latéral commun aux espaces connectés.
// $liens (défini dans inclus/espace.php de chaque espace) : fichier => [clé du libellé, icône]
// Les libellés sont traduits au chargement de la page (attribut data-t, voir api.js).
?>
<aside>

    <a id="toggle-button" href="<?= array_key_first($liens) ?>" aria-label="UpcycleConnect">
        <img src="../assets/img/LOGO-pa.png" alt="Logo">
    </a>

    <ul>
        <?php foreach ($liens as $fichier => [$libelle, $icone]): ?>
            <li>
                <a class="nav-button" href="<?= $fichier ?>">
                    <img src="../assets/img/<?= $icone ?>" alt="" aria-hidden="true">
                    <span data-t="<?= $libelle ?>"></span>
                </a>
            </li>
        <?php endforeach; ?>
    </ul>

    <button id="logout-button" onclick="deconnecter()">
        <img src="../assets/img/logout-svgrepo-com.svg" alt="" data-t-alt="commun.deconnexion">
    </button>

</aside>
