<?php
// Page d'erreur commune (403, 404, 500, 503…), même charte que la page de connexion.
// Appelée par Apache (ErrorDocument dans .htaccess), par front/routeur.php (php -S),
// par inclus/erreurs.php (erreur PHP) et par api.js (page d'un autre rôle → 403).
// Volontairement autonome : ni Vue ni API, pour s'afficher même quand le reste ne marche pas.
// Les chemins sont absolus : la page peut être servie à n'importe quelle adresse (/espace/inclus/…).

$erreurs = [
    400 => ["Requête invalide", "L'adresse demandée n'est pas valide."],
    403 => ["Accès refusé", "Cette page n'est pas accessible avec ton compte."],
    404 => ["Page introuvable", "La page que tu cherches n'existe pas ou a été déplacée."],
    500 => ["Erreur interne", "Une erreur inattendue s'est produite. Réessaie dans quelques instants."],
    503 => ["Service indisponible", "Le site est momentanément indisponible. Réessaie dans quelques instants."],
];

$code = (int) ($codeErreur ?? $_GET["code"] ?? 404);
if ($code < 400 || $code > 599) {
    $code = 404;
}
$cle = isset($erreurs[$code]) ? (string) $code : "autre";
[$titre, $message] = $erreurs[$code] ?? ["Une erreur est survenue", "La page n'a pas pu être affichée."];

if (!headers_sent()) {
    http_response_code($code);
    header("Content-Type: text/html; charset=utf-8");
    header("Cache-Control: no-store");
}

$version = fn(string $chemin) => $chemin . "?v=" . @filemtime(__DIR__ . $chemin);
?>
<!DOCTYPE html>
<html lang="fr">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="robots" content="noindex">
    <title><?= $code ?> · UpcycleConnect</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Poppins:wght@500;600;700&display=swap">
    <link rel="stylesheet" href="<?= $version("/assets/css/style.css") ?>">
    <link rel="stylesheet" href="<?= $version("/assets/css/connexion.css") ?>">
    <link rel="stylesheet" href="<?= $version("/assets/css/erreur.css") ?>">
    <script src="<?= $version("/assets/js/config.js") ?>"></script>
    <script src="<?= $version("/assets/js/erreur.js") ?>" defer></script>
</head>

<body class="page-connexion page-erreur">

    <main class="carte-connexion carte-erreur" data-code="<?= $cle ?>">

        <div class="panneau-logo">
            <img src="/assets/img/LOGO-pa.png" alt="Logo UpcycleConnect">
            <h2>UpcycleConnect</h2>
        </div>

        <div class="panneau-formulaire panneau-erreur">
            <p class="code-erreur" aria-hidden="true"><?= $code ?></p>
            <h1 data-t="titre"><?= htmlspecialchars($titre) ?></h1>
            <p class="message-erreur" data-t="message"><?= htmlspecialchars($message) ?></p>

            <div class="actions-erreur">
                <a class="bouton-erreur bouton-accueil" id="retour-accueil" href="/" data-t="accueil">Retour à l'accueil</a>
                <button class="bouton-erreur bouton-retour" id="page-precedente" type="button" data-t="precedente" hidden>Page précédente</button>
            </div>

            <p class="detail-erreur" data-t="detail">Erreur <?= $code ?></p>
        </div>

    </main>

</body>

</html>
