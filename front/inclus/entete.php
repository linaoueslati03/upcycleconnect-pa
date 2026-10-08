<?php
// En-tête commun aux pages des espaces connectés (admin, salariés, particuliers).
// Variables définies avant l'inclusion :
//   $titrePage     clé de traduction du titre de l'onglet
//   $roleRequis    rôle qui a accès à l'espace (sinon renvoi vers la connexion)
//   $feuillesStyle feuilles CSS propres à l'espace (facultatif)
//   $scriptsEspace scripts JS communs à l'espace (facultatif)
?>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title data-titre="<?= htmlspecialchars($titrePage) ?>">UpcycleConnect</title>
<link rel="stylesheet" href="../assets/css/style.css">
<?php foreach ($feuillesStyle ?? [] as $feuille): ?>
<link rel="stylesheet" href="<?= $feuille ?>">
<?php endforeach; ?>
<script src="../assets/js/config.js"></script>
<script src="../assets/js/api.js"></script>
<script>exigerConnexion(<?= json_encode($roleRequis) ?>);</script>
<script src="https://unpkg.com/vue@3.5.13/dist/vue.global.prod.js"></script>
<?php foreach ($scriptsEspace ?? [] as $script): ?>
<script src="<?= $script ?>"></script>
<?php endforeach; ?>
