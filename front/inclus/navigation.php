<aside>

    <a id="toggle-button" href="<?= array_key_first($liens) ?>" aria-label="UpcycleConnect">
        <img src="../assets/img/LOGO-pa.png" alt="Logo">
    </a>

    <ul>
        <?php foreach ($liens as $fichier => $lien): ?>
            <?php $roles = $lien[2] ?? []; ?>
            <li data-roles="<?= implode(" ", $roles) ?>" <?= $roles ? "hidden" : "" ?>>
                <a class="nav-button" href="<?= $fichier ?>">
                    <img src="../assets/img/<?= $lien[1] ?>" alt="" aria-hidden="true">
                    <span data-t="<?= $lien[0] ?>"></span>
                </a>
            </li>
        <?php endforeach; ?>
    </ul>

    <button id="logout-button" onclick="deconnecter()">
        <img src="../assets/img/logout-svgrepo-com.svg" alt="" data-t-alt="commun.deconnexion">
    </button>

</aside>
