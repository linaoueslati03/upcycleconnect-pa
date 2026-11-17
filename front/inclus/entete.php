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
?>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title data-titre="<?= htmlspecialchars($titrePage) ?>">UpcycleConnect</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Merriweather:wght@400;700&family=Poppins:wght@500;600;700&display=swap">
<link rel="stylesheet" href="<?= ressource("../assets/css/style.css") ?>">
<?php foreach ($feuillesStyle ?? [] as $feuille): ?>
<link rel="stylesheet" href="<?= ressource($feuille) ?>">
<?php endforeach; ?>
<script src="<?= ressource("../assets/js/config.js") ?>"></script>
<script src="https://unpkg.com/vue@3.5.13/dist/vue.global.prod.js"></script>
<script src="https://unpkg.com/vue-i18n@11.4.12/dist/vue-i18n.global.prod.js"></script>
<script src="<?= ressource("../assets/js/api.js") ?>"></script>
<script src="<?= ressource("../assets/js/fenetre-modale.js") ?>"></script>
<script src="<?= ressource("../assets/js/tutoriel.js") ?>"></script>
<script>exigerConnexion(...<?= json_encode($rolesAutorises) ?>);</script>
<?php foreach ($scriptsEspace ?? [] as $script): ?>
<script src="<?= ressource($script) ?>"></script>
<?php endforeach; ?>
<?php include __DIR__ . '/onesignal.php'; ?>
<?php /* Chrome prépare en arrière-plan la page d'un lien du menu dès que la souris s'y attarde :
         au clic, elle s'affiche presque instantanément, avec ses données déjà chargées. */ ?>
<script type="speculationrules">
{"prerender": [{"where": {"selector_matches": "aside a.nav-button"}, "eagerness": "moderate"}]}
</script>