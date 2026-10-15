<?php
const PARTICULIER = "particulier";
const PROFESSIONNEL = "professionnel";
const SALARIE = "salarie";

$pages = [
    "tableau-de-bord.php" => ["nav.accueil", "people-svgrepo-com.svg", [PARTICULIER, PROFESSIONNEL]],
    "planning.php"        => ["nav.mon_planning", "party-horn-svgrepo-com.svg", [PARTICULIER, SALARIE]],
    "annonces.php"        => ["nav.annonces", "book-svgrepo-com.svg", [PARTICULIER, PROFESSIONNEL]],
    "mes-annonces.php"    => ["nav.mes_annonces", "pencil-square-svgrepo-com.svg", [PARTICULIER]],
    "depots.php"          => ["nav.depots", "book-svgrepo-com.svg", [PARTICULIER, SALARIE, PROFESSIONNEL]],
    "catalogue.php"       => ["nav.catalogue", "party-horn-svgrepo-com.svg", [PARTICULIER]],
    "evenements.php"      => ["nav.evenements", "party-horn-svgrepo-com.svg", [SALARIE]],
    "ateliers.php"        => ["nav.ateliers", "pencil-square-svgrepo-com.svg", [SALARIE]],
    "projets.php"         => ["nav.mes_projets", "pencil-square-svgrepo-com.svg", [PARTICULIER, PROFESSIONNEL]],
    "conseils.php"        => ["nav.conseils", "news-svgrepo-com.svg", [PARTICULIER, SALARIE]],
    "forums.php"          => ["nav.forums", "message-square-shield-svgrepo-com.svg", [PARTICULIER]],
    "score.php"           => ["nav.score", "news-svgrepo-com.svg", [PARTICULIER]],
    "contrats.php"        => ["nav.contrats", "news-svgrepo-com.svg", [PROFESSIONNEL]],
    "statistiques.php"    => ["nav.statistiques", "party-horn-svgrepo-com.svg", [PROFESSIONNEL]],
    "mon-compte.php"      => ["nav.mon_compte", "people-svgrepo-com.svg", [PARTICULIER, PROFESSIONNEL, SALARIE]],
];

$pagesHorsMenu = [
    "creer-evenement.php" => [SALARIE],
    "creer-atelier.php"   => [SALARIE],
];

$pageCourante = basename($_SERVER["SCRIPT_NAME"]);
$rolesAutorises = $pages[$pageCourante][2] ?? $pagesHorsMenu[$pageCourante];

$liens = $pages;
$feuillesStyle = ["../assets/css/espace.css"];
$scriptsEspace = ["js/outils.js"];
