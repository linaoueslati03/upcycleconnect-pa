<?php
// Ajoute à chaque fichier CSS/JS local sa date de modification (?v=...) :
// dès qu'un fichier change, le navigateur télécharge la nouvelle version au lieu de garder l'ancienne en cache.
if (!function_exists("ressource")) {
    function ressource(string $chemin): string
    {
        $fichier = realpath(dirname($_SERVER["SCRIPT_FILENAME"]) . "/" . $chemin);
        return htmlspecialchars($chemin . ($fichier ? "?v=" . filemtime($fichier) : ""));
    }
}
