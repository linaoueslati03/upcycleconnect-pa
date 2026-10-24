<?php
require "inclus/espace.php";
$titrePage = "nav.finances";
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
                <h1>{{ $t('admin.finances.titre') }}</h1>
            </div>

            <form class="filtres-utilisateurs periode-finances" @submit.prevent="appliquerPeriode">
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

            <div v-if="resume" class="grille-finances">
                <div class="carte">
                    <span class="discret">{{ $t('admin.finances.revenus') }}</span>
                    <p class="chiffre-finances">{{ formaterMontant(resume.total_revenus) }}</p>
                </div>
                <div class="carte">
                    <span class="discret">{{ $t('admin.finances.charges') }}</span>
                    <p class="chiffre-finances">{{ formaterMontant(resume.total_charges) }}</p>
                </div>
                <div class="carte">
                    <span class="discret">{{ $t('admin.finances.resultat') }}</span>
                    <p class="chiffre-finances" :class="{ negatif: resume.resultat < 0 }">{{ formaterMontant(resume.resultat) }}</p>
                </div>
                <div class="carte">
                    <span class="discret">{{ $t('admin.finances.recurrent') }}</span>
                    <p class="chiffre-finances">{{ formaterMontant(resume.revenu_mensuel_recurrent) }}</p>
                    <span class="discret">{{ $t('admin.finances.recurrent_detail', { abonnes: resume.abonnes_actifs, campagnes: resume.campagnes_actives }) }}</span>
                </div>
            </div>

            <div class="filtres">
                <button v-for="o in onglets" :key="o" class="pill" :class="{ active: onglet === o }" @click="changerOnglet(o)">
                    {{ $t('admin.finances.onglet.' + o) }}
                </button>
            </div>

            <template v-if="onglet === 'synthese' && resume">
                <section class="carte section-categories">
                    <h2>{{ $t('admin.finances.par_source') }}</h2>
                    <p class="discret">{{ $t('admin.finances.explication_commission', { volume: formaterMontant(resume.volume_ventes) }) }}</p>
                    <table class="tableau">
                        <thead>
                            <tr>
                                <th>{{ $t('admin.finances.source_colonne') }}</th>
                                <th class="montant">{{ $t('admin.finances.nb_paiements') }}</th>
                                <th class="montant">{{ $t('admin.finances.encaisse') }}</th>
                                <th class="montant">{{ $t('admin.finances.revenu_plateforme') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="s in resume.revenus_par_source" :key="s.source">
                                <td>{{ $t('admin.finances.source.' + s.source) }}</td>
                                <td class="montant">{{ s.nb_paiements }}</td>
                                <td class="montant">{{ formaterMontant(s.montant_encaisse) }}</td>
                                <td class="montant"><strong>{{ formaterMontant(s.revenu) }}</strong></td>
                            </tr>
                        </tbody>
                    </table>
                </section>

                <section class="carte section-categories">
                    <h2>{{ $t('admin.finances.bilan_mensuel') }}</h2>
                    <table class="tableau">
                        <thead>
                            <tr>
                                <th>{{ $t('admin.finances.mois') }}</th>
                                <th class="montant">{{ $t('admin.finances.revenus') }}</th>
                                <th class="montant">{{ $t('admin.finances.charges') }}</th>
                                <th class="montant">{{ $t('admin.finances.resultat') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="m in resume.par_mois" :key="m.mois">
                                <td>{{ formaterMois(m.mois) }}</td>
                                <td class="montant">{{ formaterMontant(m.revenus) }}</td>
                                <td class="montant">{{ formaterMontant(m.charges) }}</td>
                                <td class="montant" :class="{ negatif: m.resultat < 0 }">{{ formaterMontant(m.resultat) }}</td>
                            </tr>
                        </tbody>
                    </table>
                </section>
            </template>

            <section v-if="onglet === 'transactions'" class="carte section-categories">
                <div class="filtres-utilisateurs">
                    <select v-model="filtresPaiements.source" @change="filtrerPaiements" :aria-label="$t('admin.finances.source_colonne')">
                        <option value="">{{ $t('admin.finances.toutes_sources') }}</option>
                        <option v-for="s in sources" :key="s" :value="s">{{ $t('admin.finances.source.' + s) }}</option>
                    </select>
                    <select v-model="filtresPaiements.statut" @change="filtrerPaiements" :aria-label="$t('commun.statut')">
                        <option value="">{{ $t('admin.finances.tous_statuts') }}</option>
                        <option v-for="s in statutsPaiement" :key="s" :value="s">{{ $t('admin.finances.statut_paiement.' + s) }}</option>
                    </select>
                </div>

                <p v-if="paiements.length === 0">{{ $t('admin.finances.aucun_paiement') }}</p>
                <table v-else class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('admin.finances.date') }}</th>
                            <th>{{ $t('admin.finances.payeur') }}</th>
                            <th>{{ $t('admin.finances.source_colonne') }}</th>
                            <th class="montant">{{ $t('admin.finances.montant') }}</th>
                            <th class="montant">{{ $t('admin.finances.commission') }}</th>
                            <th>{{ $t('commun.statut') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="p in paiements" :key="p.id">
                            <td class="date">{{ formaterDate(p.date_paiement) }}</td>
                            <td>{{ p.utilisateur }}</td>
                            <td>{{ $t('admin.finances.source.' + p.source) }}</td>
                            <td class="montant">{{ formaterMontant(p.montant) }}</td>
                            <td class="montant">
                                <template v-if="p.montant_commission !== null">
                                    {{ formaterMontant(p.montant_commission) }} <span class="discret">({{ formaterTaux(p.taux_commission) }})</span>
                                </template>
                                <span v-else class="discret">—</span>
                            </td>
                            <td>
                                <span class="badge-statut" :class="'paiement-' + p.statut">
                                    {{ $t('admin.finances.statut_paiement.' + p.statut) }}
                                </span>
                            </td>
                        </tr>
                    </tbody>
                </table>

                <div class="pagination">
                    <button class="lien-action" :disabled="pagePaiements <= 1" @click="changerPagePaiements(pagePaiements - 1)">← {{ $t('commun.precedent') }}</button>
                    <span>{{ $t('admin.finances.page', { page: pagePaiements, pages: nombrePagesPaiements, total: totalPaiements }) }}</span>
                    <button class="lien-action" :disabled="pagePaiements >= nombrePagesPaiements" @click="changerPagePaiements(pagePaiements + 1)">{{ $t('commun.suivant') }} →</button>
                </div>
            </section>

            <section v-if="onglet === 'abonnements'" class="carte section-categories">
                <h2>{{ $t('admin.finances.abonnements_premium') }}</h2>
                <p v-if="abonnements.length === 0">{{ $t('admin.finances.aucun_abonnement') }}</p>
                <table v-else class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('admin.finances.professionnel') }}</th>
                            <th class="montant">{{ $t('admin.finances.mensualite') }}</th>
                            <th>{{ $t('admin.finances.periode') }}</th>
                            <th>{{ $t('commun.statut') }}</th>
                            <th class="montant">{{ $t('admin.finances.total_paye') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="a in abonnements" :key="a.id">
                            <td>{{ a.utilisateur }}<br><span class="discret">{{ a.email }}</span></td>
                            <td class="montant">{{ formaterMontant(a.montant_mensuel) }}</td>
                            <td>{{ formaterPeriode(a.date_debut, a.date_fin) }}</td>
                            <td>
                                <span class="badge-statut" :class="'contrat-' + a.statut">
                                    {{ $t('admin.finances.statut_abonnement.' + a.statut) }}
                                </span>
                            </td>
                            <td class="montant">{{ formaterMontant(a.total_paye) }}</td>
                        </tr>
                    </tbody>
                </table>
            </section>

            <section v-if="onglet === 'publicite'" class="carte section-categories">
                <h2>{{ $t('admin.finances.campagnes') }}</h2>
                <p v-if="campagnes.length === 0">{{ $t('admin.finances.aucune_campagne') }}</p>
                <table v-else class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('admin.finances.annonceur') }}</th>
                            <th>{{ $t('admin.finances.mis_en_avant') }}</th>
                            <th class="montant">{{ $t('admin.finances.mensualite') }}</th>
                            <th>{{ $t('admin.finances.periode') }}</th>
                            <th>{{ $t('commun.statut') }}</th>
                            <th class="montant">{{ $t('admin.finances.total_paye') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="c in campagnes" :key="c.id">
                            <td>{{ c.annonceur }}</td>
                            <td>
                                {{ c.cible_titre || '—' }}<br>
                                <span class="discret">{{ $t('admin.finances.cible.' + c.cible_type) }}</span>
                            </td>
                            <td class="montant">{{ formaterMontant(c.montant_mensuel) }}</td>
                            <td>{{ formaterPeriode(c.date_debut, c.date_fin) }}</td>
                            <td>
                                <span class="badge-statut" :class="'contrat-' + c.statut">
                                    {{ $t('admin.finances.statut_campagne.' + c.statut) }}
                                </span>
                            </td>
                            <td class="montant">{{ formaterMontant(c.total_paye) }}</td>
                        </tr>
                    </tbody>
                </table>
            </section>

            <div v-if="onglet === 'charges'" class="mise-en-page-admin">
                <section class="carte">
                    <h2>{{ $t('admin.finances.charges_periode') }}</h2>
                    <p v-if="charges.length === 0">{{ $t('admin.finances.aucune_charge') }}</p>
                    <table v-else class="tableau">
                        <thead>
                            <tr>
                                <th>{{ $t('admin.finances.date') }}</th>
                                <th>{{ $t('admin.finances.libelle') }}</th>
                                <th>{{ $t('commun.categorie') }}</th>
                                <th class="montant">{{ $t('admin.finances.montant') }}</th>
                                <th>{{ $t('commun.actions') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="c in charges" :key="c.id">
                                <td class="date">{{ formaterDate(c.date_charge) }}</td>
                                <td>{{ c.libelle }}</td>
                                <td>{{ $t('admin.finances.categorie_charge.' + c.categorie) }}</td>
                                <td class="montant">{{ formaterMontant(c.montant) }}</td>
                                <td><button class="lien-supprimer" @click="supprimerCharge(c)">{{ $t('commun.supprimer') }}</button></td>
                            </tr>
                        </tbody>
                    </table>
                </section>

                <section class="carte carte-formulaire">
                    <h2>{{ $t('admin.finances.ajouter_charge') }}</h2>
                    <form @submit.prevent="ajouterCharge">
                        <div class="form-group">
                            <label for="libelle">{{ $t('admin.finances.libelle') }}</label>
                            <input id="libelle" type="text" v-model="nouvelleCharge.libelle" maxlength="150" required>
                        </div>
                        <div class="form-group">
                            <label for="categorie">{{ $t('commun.categorie') }}</label>
                            <select id="categorie" v-model="nouvelleCharge.categorie" required>
                                <option v-for="c in categoriesCharge" :key="c" :value="c">{{ $t('admin.finances.categorie_charge.' + c) }}</option>
                            </select>
                        </div>
                        <div class="form-group">
                            <label for="montant">{{ $t('admin.finances.montant_euros') }}</label>
                            <input id="montant" type="number" min="0.01" step="0.01" v-model.number="nouvelleCharge.montant" required>
                        </div>
                        <div class="form-group">
                            <label for="date-charge">{{ $t('admin.finances.date') }}</label>
                            <input id="date-charge" type="date" v-model="nouvelleCharge.date_charge" required>
                        </div>
                        <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                    </form>
                </section>
            </div>

            <section v-if="onglet === 'parametres' && parametres" class="carte section-categories">
                <h2>{{ $t('admin.finances.onglet.parametres') }}</h2>
                <p class="discret">{{ $t('admin.finances.explication_parametres') }}</p>
                <p v-if="messageParametres" class="succes">{{ messageParametres }}</p>
                <form class="formulaire-parametres" @submit.prevent="enregistrerParametres">
                    <div class="form-group">
                        <label for="taux">{{ $t('admin.finances.taux_commission') }}</label>
                        <input id="taux" type="number" min="5" max="10" step="0.5" v-model.number="parametres.taux_commission" required>
                    </div>
                    <div class="form-group">
                        <label for="prix-premium">{{ $t('admin.finances.prix_premium') }}</label>
                        <input id="prix-premium" type="number" min="15" max="30" step="0.5" v-model.number="parametres.prix_premium_mensuel" required>
                    </div>
                    <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                </form>
                <p class="discret">{{ $t('admin.finances.derniere_modification', { date: formaterDate(parametres.date_modification) }) }}</p>
            </section>

        </main>
    </div>

    <script src="js/finances.js"></script>
</body>

</html>
