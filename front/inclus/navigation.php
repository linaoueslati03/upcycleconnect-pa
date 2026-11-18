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

    <div class="pied-menu">
        <select id="selecteur-langue" hidden></select>
        <button id="cloche-notifications" type="button" aria-expanded="false" aria-controls="panneau-notifications" hidden>
            <svg viewBox="0 0 24 24" width="21" height="21" fill="none" aria-hidden="true">
                <path d="M6 9a6 6 0 1 1 12 0c0 4 1.5 6 2.5 7h-17C4.5 15 6 13 6 9Z" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round"/>
                <path d="M10 19.5a2.2 2.2 0 0 0 4 0" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/>
            </svg>
            <span class="compteur-notifications" hidden></span>
        </button>
        <button id="logout-button" onclick="deconnecter()">
            <img src="../assets/img/logout-svgrepo-com.svg" alt="" data-t-alt="commun.deconnexion">
        </button>
    </div>

    <section id="panneau-notifications" class="panneau-notifications" hidden>
        <header>
            <h2></h2>
            <button type="button" class="tout-lu lien-action"></button>
        </header>
        <ul></ul>
    </section>

</aside>
