<?php
// Routeur du serveur de développement PHP, qui ne lit pas le .htaccess :
//   php -S localhost:8000 -t front front/routeur.php
// Mêmes règles que front/.htaccess (Apache) : adresses sans .php, anciennes adresses en .php redirigées,
// dossiers inclus/ cachés, pages d'erreur personnalisées.

$racine = __DIR__;
$chemin = rawurldecode(parse_url($_SERVER["REQUEST_URI"], PHP_URL_PATH) ?? "/");
$requete = $_SERVER["QUERY_STRING"] ?? "";

function pageErreur(int $codeErreur): bool
{
    require __DIR__ . "/erreur.php";
    return true;
}

function rediriger(string $adresse, int $code, string $requete): bool
{
    header("Location: " . $adresse . ($requete !== "" ? "?" . $requete : ""), true, $code);
    return true;
}

// Fait exécuter une page comme si elle avait été appelée directement (chemins relatifs, nom de la page)
function executer(string $fichier, string $racine): bool
{
    $adresse = substr($fichier, strlen($racine));
    $_SERVER["SCRIPT_FILENAME"] = $fichier;
    $_SERVER["SCRIPT_NAME"] = $adresse;
    $_SERVER["PHP_SELF"] = $adresse;
    chdir(dirname($fichier));
    require $fichier;
    return true;
}

if (str_contains($chemin, "..") || str_contains($chemin, "\0")) {
    return pageErreur(400);
}

// 1. Morceaux de page et routeur : jamais accessibles par une adresse
if (preg_match('#(^|/)inclus(/|$)#', $chemin) || $chemin === "/routeur.php") {
    return pageErreur(404);
}

// 2. /espace/ et /admin/ seuls : vers la page d'accueil
if (preg_match('#^/(espace|admin)/?$#', $chemin)) {
    return rediriger("/", 302, "");
}

// 3. Anciennes adresses en .php → adresse sans .php
if (preg_match('#^(.*/)index\.php$#', $chemin, $m)) {
    return rediriger($m[1], 301, $requete);
}
if ($chemin !== "/erreur.php" && preg_match('#^(.+)\.php$#', $chemin, $m)) {
    return rediriger($m[1], 301, $requete);
}

$fichier = $racine . $chemin;

// Fichiers statiques (CSS, JS, images, traductions) : servis tels quels par PHP
if (is_file($fichier) && !str_ends_with($fichier, ".php")) {
    return false;
}

if (is_dir($fichier)) {
    if (is_file(rtrim($fichier, "/") . "/index.php")) {
        return executer(rtrim($fichier, "/") . "/index.php", $racine);
    }
    return pageErreur(403);
}

// 4. Barre oblique finale en trop (/espace/annonces/) → /espace/annonces (sinon les chemins relatifs cassent)
if (strlen($chemin) > 1 && str_ends_with($chemin, "/")) {
    return rediriger(rtrim($chemin, "/"), 301, $requete);
}

// 5. Adresse sans .php → page .php correspondante
$page = $racine . $chemin . ".php";
if (is_file($page)) {
    return executer($page, $racine);
}
if ($chemin === "/erreur.php") {
    return executer($fichier, $racine);
}

return pageErreur(404);
