<?php require_once __DIR__ . "/inclus/ressource.php"; ?>
<!DOCTYPE html>
<html lang="fr">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title data-titre="vitrine.titre_page">UpcycleConnect</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Merriweather:wght@400;700&family=Poppins:wght@500;600;700&display=swap">
    <!-- Drapeaux des langues : bibliothèque flag-icons (licence MIT) -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/flag-icons/7.2.3/css/flag-icons.min.css">
    <link rel="stylesheet" href="<?= ressource("assets/css/style.css") ?>">
    <link rel="stylesheet" href="<?= ressource("assets/css/accueil.css") ?>">
    <script src="<?= ressource("assets/js/config.js") ?>"></script>
    <script src="https://unpkg.com/vue@3.5.13/dist/vue.global.prod.js"></script>
    <script src="https://unpkg.com/vue-i18n@11.4.12/dist/vue-i18n.global.prod.js"></script>
    <script src="<?= ressource("assets/js/api.js") ?>"></script>
</head>

<body class="page-accueil">

    <div id="app">

        <!-- Bandeau plein écran : paysage, logo et slogan (sur le modèle de la maquette) -->
        <header class="bandeau">

            <div class="barre-haut">
                <nav class="liens-secondaires" :aria-label="$t('vitrine.nav.compte')">
                    <!-- Choix de la langue : liste déroulante avec drapeaux (bibliothèque flag-icons), sans rechargement -->
                    <div class="choix-langue-vitrine" ref="choixLangue" @keydown.esc="menuLangue = false">
                        <button type="button" class="langue-actuelle" :aria-label="$t('commun.langue')"
                            aria-haspopup="true" :aria-expanded="menuLangue" @click="menuLangue = !menuLangue">
                            <span v-if="pays(langue)" :class="'fi fi-' + pays(langue)"></span>
                            <span>{{ langue.toUpperCase() }}</span>
                            <span class="chevron-langue" :class="{ ouvert: menuLangue }"></span>
                        </button>
                        <transition name="liste-langues">
                            <ul v-if="menuLangue" class="liste-langues">
                                <li v-for="l in langues" :key="l.code">
                                    <button type="button" :class="{ actif: l.code === langue }"
                                        :aria-current="l.code === langue" @click="choisirLangueVitrine(l.code)">
                                        <span v-if="pays(l.code)" :class="'fi fi-' + pays(l.code)"></span>
                                        <span v-else class="code-sans-drapeau">{{ l.code.toUpperCase() }}</span>
                                        <span>{{ l.libelle }}</span>
                                    </button>
                                </li>
                            </ul>
                        </transition>
                    </div>
                    <template v-if="espace">
                        <a class="bouton-vitrine bouton-plein" :href="espace">{{ $t('vitrine.mon_espace') }}</a>
                    </template>
                    <template v-else>
                        <a class="bouton-vitrine bouton-contour" href="/connexion">{{ $t('connexion.se_connecter') }}</a>
                        <a class="bouton-vitrine bouton-plein inscription-haut" href="/inscription">{{ $t('connexion.creer_compte') }}</a>
                    </template>
                </nav>
            </div>

            <div class="barre-principale">
                <div class="marque">
                    <img src="assets/img/LOGO-pa.png" alt="" class="marque-logo">
                    <span>UpcycleConnect</span>
                </div>
                <nav class="liens-principaux" :aria-label="$t('vitrine.nav.sections')">
                    <a href="#concept">{{ $t('vitrine.nav.concept') }}</a>
                    <a href="#annonces">{{ $t('vitrine.nav.annonces') }}</a>
                    <a href="#agenda">{{ $t('vitrine.nav.agenda') }}</a>
                    <a href="#pour-qui">{{ $t('vitrine.nav.pour_qui') }}</a>
                    <a href="#rejoindre">{{ $t('vitrine.nav.rejoindre') }}</a>
                </nav>
            </div>

            <div class="bandeau-centre">
                <div class="bandeau-marque">
                    <img src="assets/img/LOGO-pa.png" alt="" class="bandeau-logo">
                    <span>UpcycleConnect</span>
                </div>
                <h1>{{ $t('vitrine.slogan') }}</h1>
                <p class="bandeau-texte">{{ $t('vitrine.sous_titre') }}</p>
                <div class="bandeau-actions">
                    <template v-if="espace">
                        <a class="bouton-vitrine bouton-plein bouton-grand" :href="espace">{{ $t('vitrine.mon_espace') }}</a>
                    </template>
                    <template v-else>
                        <a class="bouton-vitrine bouton-plein bouton-grand" href="/inscription?role=particulier">{{ $t('vitrine.cta_particulier') }}</a>
                        <a class="bouton-vitrine bouton-contour bouton-grand" href="/inscription?role=professionnel">{{ $t('vitrine.cta_professionnel') }}</a>
                    </template>
                </div>
            </div>

            <a class="defiler" href="#concept">{{ $t('vitrine.decouvrir') }}</a>
        </header>

        <main>

            <!-- Le concept en 3 étapes -->
            <section id="concept" class="bloc-vitrine">
                <div class="contenu">
                    <h2 class="titre-section">{{ $t('vitrine.concept.titre') }}</h2>
                    <p class="sous-titre-section">{{ $t('vitrine.concept.texte') }}</p>
                    <ol class="etapes">
                        <li v-for="numero in 3" :key="numero" class="etape">
                            <span class="etape-numero">{{ String(numero).padStart(2, '0') }}</span>
                            <h3>{{ $t('vitrine.concept.etapes.' + numero + '.titre') }}</h3>
                            <p>{{ $t('vitrine.concept.etapes.' + numero + '.texte') }}</p>
                        </li>
                    </ol>
                </div>
            </section>

            <!-- Chiffres clés (base de données) -->
            <section v-if="donnees" class="bande-chiffres" :aria-label="$t('vitrine.chiffres.titre')">
                <div class="contenu chiffres">
                    <div v-for="cle in clesChiffres" :key="cle" class="chiffre">
                        <strong>{{ nombre(donnees.chiffres[cle]) }}</strong>
                        <span>{{ $t('vitrine.chiffres.' + cle) }}</span>
                    </div>
                </div>
            </section>

            <!-- Dernières annonces -->
            <section id="annonces" class="bloc-vitrine bloc-clair">
                <div class="contenu">
                    <h2 class="titre-section">{{ $t('vitrine.annonces.titre') }}</h2>
                    <p class="sous-titre-section">{{ $t('vitrine.annonces.sous_titre') }}</p>

                    <p v-if="chargement" class="message-vitrine">{{ $t('commun.chargement') }}</p>
                    <p v-else-if="erreur" class="message-vitrine">{{ $t('vitrine.erreur_chargement') }}</p>
                    <p v-else-if="!donnees.annonces.length" class="message-vitrine">{{ $t('vitrine.annonces.aucune') }}</p>

                    <ul v-else class="grille-annonces">
                        <li v-for="annonce in donnees.annonces" :key="annonce.id" class="carte-annonce">
                            <div class="carte-annonce-haut">
                                <span :class="['pastille', 'pastille-' + annonce.type]">{{ $t('vitrine.annonces.' + annonce.type) }}</span>
                                <span class="prix">{{ prix(annonce) }}</span>
                            </div>
                            <h3>{{ annonce.titre }}</h3>
                            <p class="discret">{{ categorie(annonce) }}<template v-if="annonce.localisation"> · {{ annonce.localisation }}</template></p>
                            <p class="date-annonce">{{ $t('vitrine.annonces.publiee_le', { date: date(annonce.date_creation) }) }}</p>
                        </li>
                    </ul>

                    <div class="centre">
                        <a class="bouton-vitrine bouton-plein" :href="lienAnnonces">{{ $t('vitrine.annonces.voir_tout') }}</a>
                    </div>
                </div>
            </section>

            <!-- Prochains rendez-vous -->
            <section id="agenda" class="bloc-vitrine">
                <div class="contenu">
                    <h2 class="titre-section">{{ $t('vitrine.agenda.titre') }}</h2>
                    <p class="sous-titre-section">{{ $t('vitrine.agenda.sous_titre') }}</p>

                    <p v-if="chargement" class="message-vitrine">{{ $t('commun.chargement') }}</p>
                    <p v-else-if="erreur" class="message-vitrine">{{ $t('vitrine.erreur_chargement') }}</p>
                    <p v-else-if="!donnees.rendez_vous.length" class="message-vitrine">{{ $t('vitrine.agenda.aucun') }}</p>

                    <ul v-else class="liste-agenda">
                        <li v-for="rdv in donnees.rendez_vous" :key="rdv.type + rdv.id" class="rendez-vous">
                            <div class="rendez-vous-date">
                                <strong>{{ jour(rdv.date_debut) }}</strong>
                                <span>{{ mois(rdv.date_debut) }}</span>
                            </div>
                            <div class="rendez-vous-detail">
                                <span class="pastille pastille-offre">{{ $t('vitrine.agenda.types.' + rdv.type) }}</span>
                                <h3>{{ rdv.titre }}</h3>
                                <p class="discret">{{ heure(rdv.date_debut) }}<template v-if="rdv.lieu"> · {{ rdv.lieu }}</template></p>
                            </div>
                        </li>
                    </ul>

                    <div class="centre">
                        <a class="bouton-vitrine bouton-contour-sombre" :href="lienCatalogue">{{ $t('vitrine.agenda.voir_tout') }}</a>
                    </div>
                </div>
            </section>

            <!-- Pour qui ? -->
            <section id="pour-qui" class="bloc-vitrine bloc-clair">
                <div class="contenu">
                    <h2 class="titre-section">{{ $t('vitrine.pour_qui.titre') }}</h2>
                    <div class="grille-profils">
                        <article v-for="profil in profils" :key="profil" class="carte-profil">
                            <h3>{{ $t('vitrine.pour_qui.' + profil + '.titre') }}</h3>
                            <ul>
                                <li v-for="numero in 4" :key="numero">{{ $t('vitrine.pour_qui.' + profil + '.points.' + numero) }}</li>
                            </ul>
                            <a v-if="!espace" class="lien-profil"
                                :href="'/inscription?role=' + profil">{{ $t('vitrine.pour_qui.' + profil + '.action') }}</a>
                        </article>
                    </div>
                </div>
            </section>

            <!-- Appel final -->
            <section id="rejoindre" class="bande-rejoindre">
                <div class="contenu">
                    <h2>{{ $t('vitrine.rejoindre.titre') }}</h2>
                    <p>{{ $t('vitrine.rejoindre.texte') }}</p>
                    <a v-if="espace" class="bouton-vitrine bouton-plein bouton-grand" :href="espace">{{ $t('vitrine.mon_espace') }}</a>
                    <a v-else class="bouton-vitrine bouton-plein bouton-grand" href="/inscription">{{ $t('vitrine.rejoindre.bouton') }}</a>
                </div>
            </section>

        </main>

        <footer class="pied-vitrine">
            <div class="contenu pied-contenu">
                <div class="marque">
                    <img src="assets/img/LOGO-pa.png" alt="" class="marque-logo">
                    <span>UpcycleConnect</span>
                </div>
                <p>{{ $t('vitrine.pied.texte') }}</p>
                <p class="droits">© {{ annee }} UpcycleConnect · AlterNext</p>
            </div>
        </footer>

    </div>

    <script src="<?= ressource("assets/js/accueil.js") ?>"></script>
</body>

</html>
