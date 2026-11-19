<?php
// Erreurs PHP : au lieu d'une page blanche ou d'un message technique (chemins du serveur, numéros de ligne),
// le visiteur voit la page « Erreur interne » (500). Le détail reste dans le journal du serveur.
// Inclus en tête de chaque page de l'application (front/espace/inclus/espace.php, front/admin/inclus/espace.php).

ini_set("display_errors", "0");
ini_set("log_errors", "1");

// La page est d'abord écrite en mémoire : en cas d'erreur, ce qui était déjà produit est jeté
// et remplacé par la page d'erreur, avec le bon code HTTP.
ob_start();

function afficherPageErreur(int $codeErreur): void
{
    while (ob_get_level() > 0) {
        ob_end_clean();
    }
    require __DIR__ . "/../erreur.php";
}

set_exception_handler(function (Throwable $exception): void {
    error_log("UpcycleConnect : " . $exception);
    afficherPageErreur(500);
});

register_shutdown_function(function (): void {
    $erreur = error_get_last();
    if ($erreur && in_array($erreur["type"], [E_ERROR, E_PARSE, E_CORE_ERROR, E_COMPILE_ERROR, E_USER_ERROR], true)) {
        afficherPageErreur(500);
    }
});
