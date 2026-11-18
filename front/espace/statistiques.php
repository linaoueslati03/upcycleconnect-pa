<?php
require "inclus/espace.php";
$titrePage = "nav.statistiques";
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
                <h1>{{ $t('statistiques.titre') }}</h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <!-- Freemium : présentation de l'offre -->
            <section v-if="estFreemium" class="carte offre-premium">
                <span class="badge-premium">Premium</span>
                <h2>{{ $t('statistiques.verrou.titre') }}</h2>
                <p>{{ $t('statistiques.verrou.texte') }}</p>
                <ul class="avantages-premium">
                    <li v-for="a in ['statistiques', 'stocks', 'alertes']" :key="a">{{ $t('compte.abonnement.avantage.' + a) }}</li>
                </ul>
                <button class="btn-primary" :disabled="achatEnCours" @click="passerPremium">
                    {{ $t('compte.abonnement.passer', { prix: formaterPrix(prixPremium) }) }}
                </button>
            </section>

            <template v-if="stats">
                <div class="grille-finances grille-kpi">
                    <div class="carte carte-kpi carte-kpi-disponibles">
                        <span class="discret">{{ $t('statistiques.kpi.disponibles') }}</span>
                        <p class="chiffre-finances">{{ formaterNombre(stats.objets_disponibles.length) }}</p>
                        <a class="lien-action" href="#objets-disponibles">{{ $t('statistiques.kpi.voir_disponibles') }}</a>
                    </div>
                    <div class="carte carte-kpi">
                        <span class="discret">{{ $t('statistiques.kpi.recuperes') }}</span>
                        <p class="chiffre-finances">{{ formaterNombre(stats.recuperations.nombre) }}</p>
                    </div>
                    <div class="carte carte-kpi">
                        <span class="discret">{{ $t('statistiques.kpi.poids') }}</span>
                        <p class="chiffre-finances">{{ formaterKg(stats.recuperations.poids_kg) }}</p>
                    </div>
                    <div class="carte carte-kpi">
                        <span class="discret">{{ $t('statistiques.kpi.co2') }}</span>
                        <p class="chiffre-finances">≈ {{ formaterKg(stats.recuperations.co2_kg) }}</p>
                        <span class="discret">{{ $t('statistiques.kpi.co2_detail') }}</span>
                    </div>
                    <div class="carte carte-kpi">
                        <span class="discret">{{ $t('statistiques.kpi.achats') }}</span>
                        <p class="chiffre-finances">{{ formaterNombre(stats.achats.nombre) }}</p>
                        <span class="discret">{{ $t('statistiques.kpi.depenses', { montant: formaterMontant(stats.achats.depenses) }) }}</span>
                    </div>
                    <div class="carte carte-kpi">
                        <span class="discret">{{ $t('statistiques.kpi.campagnes') }}</span>
                        <p class="chiffre-finances">{{ formaterNombre(stats.campagnes_actives) }}</p>
                        <a class="lien-action" href="contrats.php">{{ $t('statistiques.kpi.voir_campagnes') }}</a>
                    </div>
                </div>

                <section id="objets-disponibles" class="carte section-categories">
                    <div class="entete-carte">
                        <h2>{{ $t('statistiques.stocks.titre') }}</h2>
                        <a class="btn-primary" href="depots.php">{{ $t('statistiques.stocks.recuperer') }}</a>
                    </div>
                    <p class="discret">{{ $t('statistiques.stocks.texte') }}</p>
                    <p v-if="stats.objets_disponibles.length === 0" class="discret">{{ $t('statistiques.stocks.aucun') }}</p>
                    <template v-else>
                        <div class="filtres">
                            <span v-for="site in stocksParSite" :key="site.site" class="pill">{{ site.site }} <span class="nombre-pill">{{ site.total }}</span></span>
                        </div>
                        <table class="tableau">
                            <thead>
                                <tr>
                                    <th>{{ $t('statistiques.stocks.objet') }}</th>
                                    <th>{{ $t('commun.site') }}</th>
                                    <th>{{ $t('depots.materiau') }}</th>
                                    <th class="montant">{{ $t('statistiques.methode.poids') }}</th>
                                    <th>{{ $t('statistiques.stocks.depose_le') }}</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="o in stats.objets_disponibles" :key="o.id">
                                    <td><strong>{{ o.description_objet }}</strong></td>
                                    <td>{{ o.site }}</td>
                                    <td>{{ o.materiau || '—' }}</td>
                                    <td class="montant">{{ o.poids_kg ? formaterKg(o.poids_kg) : '—' }}</td>
                                    <td class="date">{{ formaterDate(o.date_depot) }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </template>
                </section>

                <div class="grille-graphiques">
                    <section class="carte">
                        <h3>{{ $t('statistiques.graphique.mois') }}</h3>
                        <p class="discret">{{ $t('statistiques.graphique.mois_detail') }}</p>
                        <div class="zone-graphique"><canvas id="graphique-mois" role="img" :aria-label="$t('statistiques.graphique.mois')"></canvas></div>
                    </section>
                    <section class="carte">
                        <h3>{{ $t('statistiques.graphique.materiaux') }}</h3>
                        <p class="discret">{{ $t('statistiques.graphique.materiaux_detail') }}</p>
                        <p v-if="stats.recuperations.nombre === 0" class="discret">{{ $t('statistiques.aucune_recuperation') }}</p>
                        <div v-else class="zone-graphique"><canvas id="graphique-materiaux" role="img" :aria-label="$t('statistiques.graphique.materiaux')"></canvas></div>
                    </section>
                </div>

                <section class="carte section-categories">
                    <h2>{{ $t('statistiques.methode.titre') }}</h2>
                    <p class="discret">{{ $t('statistiques.methode.texte') }}</p>
                    <p v-if="stats.recuperations.poids_estimes" class="discret">
                        {{ $t('statistiques.methode.poids_estimes', { nombre: stats.recuperations.poids_estimes }) }}
                    </p>
                    <table class="tableau">
                        <thead>
                            <tr>
                                <th>{{ $t('depots.materiau') }}</th>
                                <th class="montant">{{ $t('statistiques.methode.objets') }}</th>
                                <th class="montant">{{ $t('statistiques.methode.poids') }}</th>
                                <th class="montant">{{ $t('statistiques.methode.facteur') }}</th>
                                <th class="montant">{{ $t('statistiques.methode.co2') }}</th>
                                <th>{{ $t('statistiques.methode.source') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="m in stats.par_materiau" :key="m.code">
                                <td>{{ nomMateriau(m) }}</td>
                                <td class="montant">{{ formaterNombre(m.nombre) }}</td>
                                <td class="montant">{{ formaterKg(m.poids_kg) }}</td>
                                <td class="montant">{{ m.facteur !== null ? formaterNombre(m.facteur, 2) : '—' }}</td>
                                <td class="montant">{{ formaterKg(m.co2_kg) }}</td>
                                <td class="source-facteur">{{ m.source || $t('statistiques.methode.sans_source') }}</td>
                            </tr>
                        </tbody>
                    </table>
                </section>


                <section class="carte section-categories">
                    <div class="entete-carte">
                        <h2>{{ $t('statistiques.alertes.titre') }}</h2>
                        <button class="btn-primary" @click="ouvrirAlerte">+ {{ $t('statistiques.alertes.nouvelle') }}</button>
                    </div>
                    <p class="discret">{{ $t('statistiques.alertes.texte') }}</p>
                    <p v-if="messageAlerte" class="succes">{{ messageAlerte }}</p>
                    <p v-if="alertes.length === 0" class="discret">{{ $t('statistiques.alertes.aucune') }}</p>
                    <ul v-else class="liste-alertes">
                        <li v-for="a in alertes" :key="a.id">
                            <span>{{ descriptionAlerte(a) }}</span>
                            <button class="lien-action lien-supprimer" @click="supprimerAlerte(a)">{{ $t('commun.supprimer') }}</button>
                        </li>
                    </ul>
                </section>
            </template>

            <fenetre-modale :ouverte="alerteOuverte" :titre="$t('statistiques.alertes.nouvelle')" @fermer="alerteOuverte = false">
                <form @submit.prevent="enregistrerAlerte">
                    <p class="discret">{{ $t('statistiques.alertes.aide') }}</p>
                    <div class="form-group">
                        <label for="alerte-materiau">{{ $t('depots.materiau') }}</label>
                        <select id="alerte-materiau" v-model="formulaireAlerte.categorie_id">
                            <option :value="null">{{ $t('statistiques.alertes.tous_materiaux') }}</option>
                            <option v-for="c in categories" :key="c.id" :value="c.id">{{ nomMateriau(c) }}</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="alerte-site">{{ $t('commun.site') }}</label>
                        <select id="alerte-site" v-model="formulaireAlerte.site">
                            <option value="">{{ $t('statistiques.alertes.tous_sites') }}</option>
                            <option v-for="s in sites" :key="s" :value="s">{{ s }}</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="alerte-mot-cle">{{ $t('statistiques.alertes.mot_cle') }}</label>
                        <input id="alerte-mot-cle" type="text" v-model="formulaireAlerte.mot_cle" maxlength="100"
                            :placeholder="$t('statistiques.alertes.mot_cle_exemple')">
                    </div>
                    <p v-if="erreurAlerte" class="erreur">{{ erreurAlerte }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                        <button type="button" class="button-draft" @click="alerteOuverte = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>
        </main>
    </div>

    <script src="https://unpkg.com/chart.js@4.4.7/dist/chart.umd.js"></script>
    <script src="<?= ressource("js/statistiques.js") ?>"></script>
</body>

</html>
