<?php
require "inclus/espace.php";
$titrePage = "nav.depot_conteneur";
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
                <h1>{{ t('nav.depot_conteneur') }}</h1>
            </div>

            <div class="deux-colonnes">

                <section>
                    <h2>{{ t('depots.suivi') }}</h2>
                    <p v-if="erreurListe" class="erreur">{{ erreurListe }}</p>
                    <p v-if="depots.length === 0" class="discret">{{ t('depots.aucun') }}</p>

                    <div v-for="d in depots" :key="d.id" class="carte section">
                        <h3>{{ d.description_objet }}</h3>
                        <p class="discret">{{ t('depots.conteneur_demande_le', { conteneur: nomConteneur(d.conteneur_id), date: formaterDate(d.date_demande) }) }}</p>

                        <ol class="suivi-etapes">
                            <li v-for="(etape, index) in etapes" :key="etape"
                                :class="{ atteinte: index <= indexStatut(d.statut) }">
                                {{ t('depot.' + etape) }}
                            </li>
                        </ol>

                        <p v-if="d.code_ouverture">
                            {{ t('depots.code_ouverture') }} : <span class="code">{{ d.code_ouverture }}</span><br>
                            {{ t('depots.code_barre') }} : <span class="code">{{ d.code_barre }}</span>
                        </p>
                        <p v-else class="discret">{{ t('depots.codes_apres_validation') }}</p>
                    </div>
                </section>

                <div class="carte">
                    <h2>{{ t('depots.demander') }}</h2>
                    <p class="discret">{{ t('depots.explication') }}</p>

                    <form @submit.prevent="demanderDepot">
                        <div class="form-group">
                            <label for="description">{{ t('depots.description_objet') }}</label>
                            <textarea id="description" v-model="formulaire.description_objet" required
                                :placeholder="t('depots.placeholder_objet')"></textarea>
                        </div>

                        <div class="form-group">
                            <label for="conteneur">{{ t('depots.conteneur_souhaite') }}</label>
                            <select id="conteneur" v-model="formulaire.conteneur_id" required>
                                <option :value="null" disabled>{{ t('depots.choisir_conteneur') }}</option>
                                <option v-for="c in conteneurs" :key="c.id" :value="c.id">{{ c.site }}</option>
                            </select>
                        </div>

                        <p v-if="erreur" class="erreur">{{ erreur }}</p>
                        <p v-if="confirmation" class="succes">{{ t('depots.demande_envoyee') }}</p>

                        <button type="submit" class="button-submit">{{ t('depots.envoyer') }}</button>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/depots.js"></script>
</body>

</html>
