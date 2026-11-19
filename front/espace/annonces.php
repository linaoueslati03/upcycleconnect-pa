<?php
require "inclus/espace.php";
$titrePage = "nav.annonces";
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
                <h1>{{ $t('nav.annonces') }}</h1>
                <a v-if="roleConnecte() === 'particulier'" class="btn-primary" href="mes-annonces">{{ $t('annonces.deposer') }}</a>
            </div>

            <form class="barre-filtres" @submit.prevent="chargerAnnonces">
                <div class="form-group">
                    <label for="type">{{ $t('commun.type') }}</label>
                    <select id="type" v-model="filtres.type">
                        <option value="">{{ $t('commun.tous') }}</option>
                        <option value="don">{{ $t('annonce.type.don') }}</option>
                        <option value="vente">{{ $t('annonce.type.vente') }}</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="localisation">{{ $t('commun.localisation') }}</label>
                    <input id="localisation" type="search" v-model="filtres.localisation" :placeholder="$t('annonces.exemple_localisation')">
                </div>
            </form>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="!erreur && annonces.length === 0" class="discret">{{ $t('annonces.aucune') }}</p>

            <div class="deux-colonnes">
                <div class="grille">
                    <div v-for="a in annonces" :key="a.id" class="carte carte-cliquable"
                        :class="{ 'carte-selectionnee': selection && selection.id === a.id, 'carte-a-la-une': a.a_la_une }" @click="selection = a">
                        <span v-if="a.a_la_une" class="badge-a-la-une">{{ $t('a_la_une.badge') }}</span>
                        <span class="badge">{{ $t('annonce.type.' + a.type) }}</span>
                        <h3>{{ a.titre }}</h3>
                        <p class="discret">{{ a.localisation }} · {{ formaterDate(a.date_creation) }}</p>
                        <p v-if="a.type === 'vente'"><strong>{{ formaterPrix(a.prix) }}</strong></p>
                    </div>
                </div>

                <div v-if="selection" class="carte">
                    <span class="badge">{{ $t('annonce.type.' + selection.type) }}</span>
                    <h2>{{ selection.titre }}</h2>
                    <p>{{ selection.description || $t('annonces.pas_de_description') }}</p>
                    <ul class="liste-simple">
                        <li><strong>{{ $t('commun.categorie') }} :</strong> {{ nomCategorie(selection.categorie_id) }}</li>
                        <li><strong>{{ $t('commun.localisation') }} :</strong> {{ selection.localisation || $t('commun.non_precisee') }}</li>
                        <li v-if="selection.type === 'vente'"><strong>{{ $t('commun.prix') }} :</strong> {{ formaterPrix(selection.prix) }}</li>
                        <li><strong>{{ $t('annonces.publiee_le') }} :</strong> {{ formaterDate(selection.date_creation) }}</li>
                    </ul>
                    <template v-if="selection.type === 'vente'">
                        <p v-if="estMonObjet(selection)" class="discret">{{ $t('annonces.votre_objet') }}</p>
                        <button v-else class="button-submit bouton-acheter" :disabled="achatEnCours" @click="acheter(selection)">
                            {{ $t('annonces.acheter', { prix: formaterPrix(selection.prix) }) }}
                        </button>
                        <p v-if="erreurAchat" class="erreur">{{ erreurAchat }}</p>
                    </template>
                </div>
            </div>

        </main>
    </div>

    <script src="<?= ressource("js/annonces.js") ?>"></script>
</body>

</html>
