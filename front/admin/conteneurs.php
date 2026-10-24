<?php
require "inclus/espace.php";
$titrePage = "nav.conteneurs";
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
                <h1>{{ $t('admin.conteneurs.titre') }}</h1>
            </div>

            <p class="discret">{{ $t('admin.conteneurs.explication') }}</p>
            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="filtres">
                <button v-for="filtre in filtres" :key="filtre.valeur"
                    class="pill" :class="{ active: filtreStatut === filtre.valeur }"
                    @click="filtreStatut = filtre.valeur">
                    {{ $t(filtre.nom) }} ({{ compter(filtre.valeur) }})
                </button>
            </div>

            <section class="carte section-categories">
                <form class="ajout-categorie" @submit.prevent="ajouter">
                    <input type="text" v-model="nouveauSite" :placeholder="$t('admin.conteneurs.nouveau_site')"
                        :aria-label="$t('admin.conteneurs.nouveau_site')" required>
                    <button type="submit" class="button-submit">{{ $t('admin.conteneurs.ajouter') }}</button>
                </form>

                <p v-if="conteneursFiltres.length === 0">{{ $t('admin.conteneurs.aucun') }}</p>
                <table v-else class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('commun.site') }}</th>
                            <th>{{ $t('commun.statut') }}</th>
                            <th>{{ $t('admin.conteneurs.objet_en_cours') }}</th>
                            <th>{{ $t('commun.actions') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="c in conteneursFiltres" :key="c.id">
                            <td>{{ c.site }} <span class="discret">#{{ c.id }}</span></td>
                            <td>
                                <span class="badge-statut" :class="'conteneur-' + c.statut">
                                    {{ $t('admin.conteneurs.statut.' + c.statut) }}
                                </span>
                            </td>
                            <td>
                                <template v-if="c.depot_id">
                                    <strong>{{ c.description_objet }}</strong><br>
                                    <span class="discret">{{ c.deposant }} · {{ $t('depot.' + c.depot_statut) }}</span><br>
                                    {{ $t('depots.code_ouverture') }} : <span class="code">{{ c.code_ouverture }}</span><br>
                                    {{ $t('depots.code_barre') }} : <span class="code">{{ c.code_barre }}</span>
                                </template>
                                <span v-else class="discret">—</span>
                            </td>
                            <td>
                                <button class="lien-action" @click="ouvrirChangement(c)">{{ $t('admin.conteneurs.changer_statut') }}</button>
                                <button class="lien-action" @click="afficherHistorique(c)">{{ $t('admin.conteneurs.historique') }}</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </section>

            <section v-if="changement.conteneur" class="carte section-categories">
                <h2>{{ $t('admin.conteneurs.changer_statut_de', { site: changement.conteneur.site, id: changement.conteneur.id }) }}</h2>
                <p v-if="changement.conteneur.depot_id" class="erreur">{{ $t('admin.conteneurs.alerte_depot_en_cours') }}</p>

                <form @submit.prevent="enregistrerStatut">
                    <div class="form-group">
                        <label for="nouveau-statut">{{ $t('admin.conteneurs.nouveau_statut') }}</label>
                        <select id="nouveau-statut" v-model="changement.statut" required>
                            <option v-for="s in statuts" :key="s" :value="s" :disabled="s === changement.conteneur.statut">
                                {{ $t('admin.conteneurs.statut.' + s) }}
                            </option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="motif">{{ $t('admin.conteneurs.motif') }}</label>
                        <input id="motif" type="text" v-model="changement.motif" required
                            :placeholder="$t('admin.conteneurs.motif_exemple')">
                    </div>
                    <div class="boutons-formulaire">
                        <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                        <button type="button" class="lien-action" @click="fermerChangement">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </section>

            <section v-if="historique.conteneur" class="carte section-categories">
                <h2>{{ $t('admin.conteneurs.historique_de', { site: historique.conteneur.site, id: historique.conteneur.id }) }}</h2>
                <p v-if="historique.lignes.length === 0">{{ $t('admin.conteneurs.historique_vide') }}</p>
                <table v-else class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('admin.conteneurs.date') }}</th>
                            <th>{{ $t('admin.conteneurs.changement') }}</th>
                            <th>{{ $t('admin.conteneurs.par') }}</th>
                            <th>{{ $t('admin.conteneurs.motif') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="h in historique.lignes" :key="h.id">
                            <td>{{ formaterDateHeure(h.date_changement) }}</td>
                            <td>
                                <template v-if="h.ancien_statut">{{ $t('admin.conteneurs.statut.' + h.ancien_statut) }} → </template>
                                {{ $t('admin.conteneurs.statut.' + h.nouveau_statut) }}
                            </td>
                            <td>{{ h.auteur || '—' }}</td>
                            <td>
                                <template v-if="h.motif">{{ h.motif }}</template>
                                <template v-else-if="h.depot_id">{{ $t('admin.conteneurs.motif_depot', { id: h.depot_id }) }}</template>
                                <template v-else>—</template>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <button class="lien-action" @click="historique.conteneur = null">{{ $t('admin.conteneurs.fermer') }}</button>
            </section>

        </main>
    </div>

    <script src="js/conteneurs.js"></script>
</body>

</html>
