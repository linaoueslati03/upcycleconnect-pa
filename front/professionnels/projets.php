<?php
require "inclus/espace.php";
$titrePage = "nav.mes_projets";
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
                <h1>{{ t('projets.titre') }}</h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="deux-colonnes">

                <section>
                    <p v-if="projets.length === 0" class="discret">{{ t('projets.aucun') }}</p>
                    <div class="grille section">
                        <div v-for="p in projets" :key="p.id" class="carte carte-cliquable"
                            :class="{ 'carte-selectionnee': projet && projet.id === p.id }" @click="ouvrirProjet(p.id)">
                            <span class="badge" :class="{ 'badge-attente': !p.partage_public }">
                                {{ p.partage_public ? t('projets.partage') : t('projets.prive') }}
                            </span>
                            <h3>{{ p.titre }}</h3>
                            <p class="discret">{{ t('projets.cree_le', { date: formaterDate(p.date_creation) }) }}</p>
                        </div>
                    </div>

                    <div v-if="projet" class="carte">
                        <h2>{{ projet.titre }}</h2>
                        <p>{{ projet.description }}</p>

                        <div class="form-buttons section">
                            <button class="button-submit" @click="basculerPartage">
                                {{ projet.partage_public ? t('projets.ne_plus_partager') : t('projets.partager') }}
                            </button>
                            <button class="button-draft" @click="supprimerProjet">{{ t('projets.supprimer') }}</button>
                        </div>

                        <h3>{{ t('projets.etapes') }}</h3>
                        <p v-if="projet.etapes.length === 0" class="discret">{{ t('projets.aucune_etape') }}</p>
                        <ol>
                            <li v-for="e in projet.etapes" :key="e.id">
                                {{ e.description }}
                                <span class="discret"> · {{ formaterDate(e.date) }}</span>
                                <a v-if="e.photo_url" :href="e.photo_url" target="_blank" rel="noopener"> ({{ t('projets.photo') }})</a>
                            </li>
                        </ol>

                        <form @submit.prevent="ajouterEtape">
                            <div class="form-group">
                                <label for="etape">{{ t('projets.nouvelle_etape') }}</label>
                                <textarea id="etape" v-model="nouvelleEtape.description" required></textarea>
                            </div>
                            <div class="form-group">
                                <label for="photo">{{ t('projets.lien_photo') }}</label>
                                <input id="photo" type="url" v-model="nouvelleEtape.photo_url">
                            </div>
                            <button type="submit" class="button-submit">{{ t('projets.ajouter_etape') }}</button>
                        </form>
                    </div>
                </section>

                <div class="carte">
                    <h2>{{ t('projets.nouveau') }}</h2>
                    <form @submit.prevent="creerProjet">
                        <div class="form-group">
                            <label for="titre">{{ t('commun.titre') }}</label>
                            <input id="titre" type="text" v-model="formulaire.titre" required>
                        </div>
                        <div class="form-group">
                            <label for="description">{{ t('commun.description') }}</label>
                            <textarea id="description" v-model="formulaire.description"></textarea>
                        </div>
                        <div class="form-group">
                            <label><input type="checkbox" v-model="formulaire.partage_public"> {{ t('projets.partager') }}</label>
                        </div>
                        <button type="submit" class="button-submit">{{ t('projets.creer') }}</button>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/projets.js"></script>
</body>

</html>
