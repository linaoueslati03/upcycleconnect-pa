<?php
require "inclus/espace.php";
$titrePage = "nav.reporting";
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
                <h1>{{ $t('admin.reporting.titre') }}</h1>
            </div>

            <form v-if="onglet !== 'predictions'" class="filtres-utilisateurs periode-finances" @submit.prevent="appliquerPeriode">
                <div class="form-group">
                    <label for="debut">{{ $t('commun.date_debut') }}</label>
                    <input id="debut" type="date" v-model="periode.debut" required>
                </div>
                <div class="form-group">
                    <label for="fin">{{ $t('commun.date_fin') }}</label>
                    <input id="fin" type="date" v-model="periode.fin" required>
                </div>
                <button type="submit" class="button-submit">{{ $t('admin.finances.appliquer') }}</button>
                <button v-for="raccourci in raccourcis" :key="raccourci" type="button" class="pill"
                    @click="choisirRaccourci(raccourci)">
                    {{ $t('admin.finances.raccourci.' + raccourci) }}
                </button>
            </form>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <div class="filtres">
                <button v-for="o in onglets" :key="o" class="pill" :class="{ active: onglet === o }" @click="changerOnglet(o)">
                    {{ $t('admin.reporting.onglet.' + o) }}
                </button>
            </div>

            <template v-if="onglet === 'acteurs' && acteurs">
                <div class="grille-finances grille-totaux-reporting">
                    <div v-for="cle in totaux" :key="cle" class="carte">
                        <span class="discret">{{ $t('admin.reporting.totaux.' + cle) }}</span>
                        <p class="chiffre-finances">{{ formaterNombre(acteurs.totaux[cle]) }}</p>
                    </div>
                </div>

                <h2 class="titre-bloc-reporting">{{ $t('admin.reporting.etat_actuel') }}</h2>
                <p class="discret">{{ $t('admin.reporting.etat_actuel_detail') }}</p>
                <div class="grille-graphiques">
                    <section v-for="g in graphiquesEtat" :key="g" class="carte">
                        <h3>{{ $t('admin.reporting.graphique.' + g + '.titre') }}</h3>
                        <p class="discret">{{ $t('admin.reporting.graphique.' + g + '.detail') }}</p>
                        <div class="zone-graphique">
                            <canvas :id="'graphique-' + g" role="img" :aria-label="$t('admin.reporting.graphique.' + g + '.titre')"></canvas>
                        </div>
                    </section>
                </div>

                <h2 class="titre-bloc-reporting">{{ $t('admin.reporting.sur_periode') }}</h2>
                <p class="discret">{{ $t('admin.reporting.sur_periode_detail', { debut: formaterDate(acteurs.debut), fin: formaterDate(acteurs.fin) }) }}</p>
                <div class="grille-graphiques grille-graphiques-larges">
                    <section v-for="g in graphiquesPeriode" :key="g" class="carte">
                        <h3>{{ $t('admin.reporting.graphique.' + g + '.titre') }}</h3>
                        <p class="discret">{{ $t('admin.reporting.graphique.' + g + '.detail') }}</p>
                        <div class="zone-graphique">
                            <canvas :id="'graphique-' + g" role="img" :aria-label="$t('admin.reporting.graphique.' + g + '.titre')"></canvas>
                        </div>
                    </section>
                </div>
            </template>

            <template v-if="onglet === 'prestations'">
                <div class="filtres-utilisateurs periode-finances">
                    <div class="form-group">
                        <label for="type-prestation">{{ $t('admin.reporting.prestations.filtre_type') }}</label>
                        <select id="type-prestation" v-model="typePrestation" @change="chargerPrestations">
                            <option value="">{{ $t('admin.reporting.prestations.tous_types') }}</option>
                            <option v-for="t in typesPrestation" :key="t" :value="t">{{ $t('admin.reporting.prestations.type.' + t) }}</option>
                        </select>
                    </div>
                </div>

                <template v-if="prestations">
                    <div class="grille-finances grille-totaux-reporting">
                        <div class="carte">
                            <span class="discret">{{ $t('admin.reporting.prestations.totaux.inscriptions') }}</span>
                            <p class="chiffre-finances">{{ formaterNombre(prestations.totaux.inscriptions) }}</p>
                        </div>
                        <div class="carte">
                            <span class="discret">{{ $t('admin.reporting.prestations.totaux.participants') }}</span>
                            <p class="chiffre-finances">{{ formaterNombre(prestations.totaux.participants) }}</p>
                        </div>
                        <div class="carte">
                            <span class="discret">{{ $t('admin.reporting.prestations.totaux.reinscription') }}</span>
                            <p class="chiffre-finances">{{ formaterPourcentage(prestations.totaux.taux_reinscription) }}</p>
                        </div>
                        <div class="carte">
                            <span class="discret">{{ $t('admin.reporting.prestations.totaux.remplissage') }}</span>
                            <p class="chiffre-finances">{{ formaterPourcentage(prestations.totaux.taux_remplissage) }}</p>
                        </div>
                        <div class="carte">
                            <span class="discret">{{ $t('admin.reporting.prestations.totaux.chiffre_affaires') }}</span>
                            <p class="chiffre-finances">{{ formaterMontant(prestations.totaux.chiffre_affaires) }}</p>
                        </div>
                    </div>
                    <p class="discret">{{ $t('admin.reporting.sur_periode_detail', { debut: formaterDate(prestations.debut), fin: formaterDate(prestations.fin) }) }}
                        {{ $t('admin.reporting.prestations.definition') }}</p>

                    <div class="grille-graphiques">
                        <section v-for="g in graphiquesPrestations" :key="g" class="carte" :class="{ 'graphique-large': g === 'top' }">
                            <h3>{{ $t('admin.reporting.prestations.graphique.' + g + '.titre') }}</h3>
                            <p class="discret">{{ $t('admin.reporting.prestations.graphique.' + g + '.detail') }}</p>
                            <div class="zone-graphique">
                                <canvas :id="'graphique-' + g" role="img" :aria-label="$t('admin.reporting.prestations.graphique.' + g + '.titre')"></canvas>
                            </div>
                        </section>
                    </div>

                    <section class="carte section-categories">
                        <h2>{{ $t('admin.reporting.prestations.tableau') }}</h2>
                        <p class="discret">{{ $t('admin.reporting.prestations.tableau_detail') }}</p>
                        <p v-if="prestations.prestations.length === 0">{{ $t('admin.reporting.prestations.aucune') }}</p>
                        <table v-else class="tableau">
                            <thead>
                                <tr>
                                    <th>{{ $t('admin.reporting.prestations.colonne.prestation') }}</th>
                                    <th>{{ $t('admin.reporting.prestations.colonne.type') }}</th>
                                    <th class="montant">{{ $t('admin.reporting.prestations.colonne.sessions') }}</th>
                                    <th class="montant">{{ $t('admin.reporting.prestations.colonne.inscriptions') }}</th>
                                    <th class="montant">{{ $t('admin.reporting.prestations.colonne.participants') }}</th>
                                    <th class="montant">{{ $t('admin.reporting.prestations.colonne.reinscription') }}</th>
                                    <th class="montant">{{ $t('admin.reporting.prestations.colonne.remplissage') }}</th>
                                    <th class="montant">{{ $t('admin.reporting.prestations.colonne.chiffre_affaires') }}</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="p in prestations.prestations" :key="p.type + p.titre">
                                    <td>{{ p.titre }}</td>
                                    <td>{{ $t('admin.reporting.prestations.type.' + p.type) }}</td>
                                    <td class="montant">{{ formaterNombre(p.sessions) }}</td>
                                    <td class="montant"><strong>{{ formaterNombre(p.inscriptions) }}</strong></td>
                                    <td class="montant">{{ formaterNombre(p.participants) }}</td>
                                    <td class="montant">{{ formaterPourcentage(p.taux_reinscription) }}</td>
                                    <td class="montant">{{ formaterPourcentage(p.taux_remplissage) }}</td>
                                    <td class="montant">{{ formaterMontant(p.chiffre_affaires) }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </section>
                </template>
            </template>

            <template v-if="onglet === 'predictions'">
                <p class="discret">{{ $t('admin.reporting.predictions.calcul_actuel') }}</p>

                <div v-if="predictions" class="carte bandeau-modele">
                    <strong>{{ $t('admin.reporting.predictions.modele_titre', { date: formaterDateHeure(predictions.modele.version) }) }}</strong>
                    <p>{{ $t('admin.reporting.predictions.modele_detail', { exactitude: formaterProbabilite(predictions.modele.exactitude_test) }) }}</p>
                </div>

                <section v-if="detail" id="detail-prediction" class="carte detail-prediction">
                    <div class="entete-detail">
                        <h2>{{ detail.utilisateur.prenom }} {{ detail.utilisateur.nom }}</h2>
                        <button class="lien-action" @click="detail = null">{{ $t('admin.utilisateurs.fermer') }}</button>
                    </div>
                    <p class="discret">
                        {{ $t('admin.reporting.predictions.role.' + detail.utilisateur.role) }} · {{ detail.utilisateur.email }} ·
                        {{ $t('admin.reporting.predictions.resume_acteur', {
                            jours: formaterNombre(detail.variables.anciennete_jours),
                            score: formaterNombre(detail.variables.upcycling_score),
                            demandes: formaterNombre(nombreDemandes(detail.variables)) }) }}
                    </p>
                    <p class="prediction-principale">
                        {{ $t('admin.reporting.predictions.prochaine') }}
                        <strong>{{ nomClasse(detail.prediction.classe) }}</strong>
                        ({{ formaterProbabilite(detail.prediction.probabilites[0].probabilite) }})
                    </p>
                    <div class="mise-en-page-detail">
                        <div>
                            <h3>{{ $t('admin.reporting.predictions.probabilites') }}</h3>
                            <div v-for="p in detail.prediction.probabilites" :key="p.classe" class="barre-probabilite">
                                <span>{{ nomClasse(p.classe) }}</span>
                                <div class="barre"><div :style="{ width: (p.probabilite * 100) + '%', background: couleurClasse(p.classe) }"></div></div>
                                <span class="montant">{{ formaterProbabilite(p.probabilite) }}</span>
                            </div>
                        </div>
                        <div>
                            <h3>{{ $t('admin.reporting.predictions.pourquoi') }}</h3>
                            <ol class="chemin-prediction">
                                <li v-for="(etape, i) in detail.prediction.chemin" :key="i">{{ phraseEtape(etape) }}</li>
                            </ol>
                            <p class="discret">{{ $t('admin.reporting.predictions.effectif', { nombre: formaterNombre(detail.prediction.effectif) }) }}</p>
                        </div>
                    </div>
                </section>

                <div class="filtres-utilisateurs periode-finances">
                    <div class="form-group">
                        <label for="classe-prediction">{{ $t('admin.reporting.predictions.filtre_classe') }}</label>
                        <select id="classe-prediction" v-model="filtrePrediction.classe" @change="filtrerPredictions">
                            <option value="">{{ $t('admin.reporting.predictions.toutes_classes') }}</option>
                            <option v-for="c in classesPrediction" :key="c" :value="c">{{ nomClasse(c) }}</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="role-prediction">{{ $t('commun.role') }}</label>
                        <select id="role-prediction" v-model="filtrePrediction.role" @change="filtrerPredictions">
                            <option value="">{{ $t('admin.reporting.predictions.tous_roles') }}</option>
                            <option v-for="r in rolesActeur" :key="r" :value="r">{{ $t('admin.reporting.predictions.role.' + r) }}</option>
                        </select>
                    </div>
                </div>

                <template v-if="predictions">
                    <div class="grille-graphiques">
                        <section class="carte graphique-large">
                            <h3>{{ $t('admin.reporting.predictions.graphique_titre') }}</h3>
                            <p class="discret">{{ $t('admin.reporting.predictions.graphique_detail') }}</p>
                            <div class="zone-graphique">
                                <canvas id="graphique-repartition" role="img" :aria-label="$t('admin.reporting.predictions.graphique_titre')"></canvas>
                            </div>
                        </section>
                    </div>

                    <section class="carte section-categories">
                        <h2>{{ $t('admin.reporting.predictions.liste_titre') }}</h2>
                        <p class="discret">{{ $t('admin.reporting.predictions.liste_detail') }}</p>
                        <p v-if="predictions.predictions.length === 0">{{ $t('admin.reporting.predictions.aucun') }}</p>
                        <table v-else class="tableau">
                            <thead>
                                <tr>
                                    <th>{{ $t('admin.reporting.predictions.colonne.acteur') }}</th>
                                    <th>{{ $t('commun.role') }}</th>
                                    <th>{{ $t('admin.reporting.predictions.colonne.prochaine') }}</th>
                                    <th>{{ $t('admin.reporting.predictions.colonne.confiance') }}</th>
                                    <th class="montant">{{ $t('admin.reporting.predictions.colonne.derniere') }}</th>
                                    <th>{{ $t('commun.actions') }}</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="p in predictions.predictions" :key="p.id">
                                    <td>{{ p.prenom }} {{ p.nom }}<br><span class="discret">{{ p.email }}</span></td>
                                    <td>{{ $t('admin.reporting.predictions.role.' + p.role) }}</td>
                                    <td>{{ nomClasse(p.classe) }}</td>
                                    <td>
                                        <div class="barre-confiance">
                                            <div class="barre"><div :style="{ width: (p.probabilite * 100) + '%', background: couleurClasse(p.classe) }"></div></div>
                                            <span>{{ formaterProbabilite(p.probabilite) }}</span>
                                        </div>
                                    </td>
                                    <td class="montant">{{ p.jours_depuis_derniere_demande >= 999 ? '—' : $t('admin.reporting.predictions.il_y_a', { jours: formaterNombre(p.jours_depuis_derniere_demande) }) }}</td>
                                    <td><button class="lien-action" @click="expliquer(p.id)">{{ $t('admin.reporting.predictions.expliquer') }}</button></td>
                                </tr>
                            </tbody>
                        </table>

                        <div class="pagination">
                            <button class="lien-action" :disabled="pagePredictions <= 1" @click="chargerPredictions(pagePredictions - 1)">← {{ $t('commun.precedent') }}</button>
                            <span>{{ $t('admin.reporting.predictions.page', { page: pagePredictions, pages: nombrePagesPredictions, total: predictions.total }) }}</span>
                            <button class="lien-action" :disabled="pagePredictions >= nombrePagesPredictions" @click="chargerPredictions(pagePredictions + 1)">{{ $t('commun.suivant') }} →</button>
                        </div>
                    </section>
                </template>
            </template>

        </main>
    </div>

    <script src="https://unpkg.com/chart.js@4.4.7/dist/chart.umd.js"></script>
    <script src="js/periode.js"></script>
    <script src="js/reporting.js"></script>
</body>

</html>
