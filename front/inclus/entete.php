<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title data-titre="<?= htmlspecialchars($titrePage) ?>">UpcycleConnect</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Merriweather:wght@400;700&family=Poppins:wght@500;600;700&display=swap">
<link rel="stylesheet" href="../assets/css/style.css">
<?php foreach ($feuillesStyle ?? [] as $feuille): ?>
<link rel="stylesheet" href="<?= $feuille ?>">
<?php endforeach; ?>
<script src="../assets/js/config.js"></script>
<script src="https://unpkg.com/vue@3.5.13/dist/vue.global.prod.js"></script>
<script src="https://unpkg.com/vue-i18n@11.4.12/dist/vue-i18n.global.prod.js"></script>
<script src="../assets/js/api.js"></script>
<script>exigerConnexion(...<?= json_encode($rolesAutorises) ?>);</script>
<?php foreach ($scriptsEspace ?? [] as $script): ?>
<script src="<?= $script ?>"></script>
<?php endforeach; ?>
