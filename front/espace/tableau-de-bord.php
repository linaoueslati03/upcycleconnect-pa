<?php
require "inclus/espace.php";
$titrePage = "nav.accueil";
?>
<!DOCTYPE html>
<html lang="fr">

<head>
    <?php include "../inclus/entete.php"; ?>
</head>

<body>
    <?php include "../inclus/navigation.php"; ?>

    <div id="app">
        <main>

            <div class="page-header">
                <h1>{{ $t('accueil.bonjour', { prenom }) }}</h1>
                <div v-if="roleConnecte() === 'particulier'" class="form-buttons">
                    <a class="btn-primary" href="mes-annonces">{{ $t('annonces.deposer') }}</a>
                    <a class="btn-primary" href="depots">{{ $t('accueil.demander_depot') }}</a>
                </div>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div v-if="resume && roleConnecte() === 'particulier'" class="grille">
                <a class="carte carte-cliquable" href="mes-annonces">
                    <span class="discret">{{ $t('accueil.annonces_en_ligne') }}</span>
                    <p class="chiffre-cle">{{ resume.annonces_actives }}</p>
                </a>
                <a class="carte carte-cliquable" href="depots">
                    <span class="discret">{{ $t('accueil.dernier_depot') }}</span>
                    <p class="chiffre-cle">{{ libelleDepot(resume.dernier_depot_statut) }}</p>
                </a>
                <a class="carte carte-cliquable" href="score">
                    <span class="discret">{{ $t('nav.score') }}</span>
                    <p class="chiffre-cle">{{ resume.upcycling_score }}</p>
                </a>
                <a class="carte carte-cliquable" href="planning">
                    <span class="discret">{{ $t('accueil.prochain_rdv') }}</span>
                    <p class="chiffre-cle" v-if="resume.prochaine_offre_titre">{{ resume.prochaine_offre_titre }}</p>
                    <p class="discret" v-if="resume.prochaine_offre_date">{{ formaterDate(resume.prochaine_offre_date) }}</p>
                    <p class="chiffre-cle" v-else>—</p>
                </a>
            </div>

            <div v-if="roleConnecte() === 'professionnel'" class="grille">
                <a class="carte carte-cliquable" href="mon-compte">
                    <span class="discret">{{ $t('compte.abonnement.titre') }}</span>
                    <p class="chiffre-cle">
                        <span v-if="abonnement === 'premium'" class="badge-premium">Premium</span>
                        <template v-else>Freemium</template>
                    </p>
                </a>
                <a class="carte carte-cliquable" href="contrats">
                    <span class="discret">{{ $t('accueil.pro.campagnes_actives') }}</span>
                    <p class="chiffre-cle">{{ campagnesActives }}</p>
                </a>
                <a class="carte carte-cliquable" href="depots">
                    <span class="discret">{{ $t('accueil.pro.recuperer') }}</span>
                    <p class="discret">{{ $t('accueil.pro.recuperer_detail') }}</p>
                </a>
                <a class="carte carte-cliquable" href="statistiques">
                    <span class="discret">{{ $t('nav.statistiques') }}</span>
                    <p class="discret">{{ abonnement === 'premium' ? $t('accueil.pro.stats_premium') : $t('accueil.pro.stats_freemium') }}</p>
                </a>
            </div>

            <section v-if="projetsALaUne.length" class="a-la-une">
                <h2 class="titre-bloc-reporting">{{ $t('a_la_une.projets') }}</h2>
                <div class="grille">
                    <article v-for="p in projetsALaUne" :key="p.id" class="carte carte-a-la-une">
                        <span class="badge-a-la-une">{{ $t('a_la_une.badge') }}</span>
                        <h3>{{ p.titre }}</h3>
                        <p v-if="p.description">{{ p.description }}</p>
                        <p class="discret">{{ $t('a_la_une.par', { nom: p.professionnel }) }}</p>
                    </article>
                </div>
            </section>

        </main>
    </div>

    <script src="<?= ressource("js/tableau-de-bord.js") ?>"></script>
</body>

</html>
