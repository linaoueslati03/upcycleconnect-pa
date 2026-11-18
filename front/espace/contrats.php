<?php
require "inclus/espace.php";
$titrePage = "nav.contrats";
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
                <h1>{{ $t('contrats.titre') }}</h1>
                <button class="btn-primary" @click="ouvrirCampagne">+ {{ $t('contrats.campagnes.nouvelle') }}</button>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="message" class="succes">{{ message }}</p>

            <!-- Abonnement -->
            <section v-if="abonnement" class="carte section-categories carte-abonnement" :class="{ premium: abonnement.abonnement === 'premium' }">
                <div class="entete-carte">
                    <h2>{{ $t('compte.abonnement.titre') }}</h2>
                    <a class="lien-action" href="mon-compte.php">{{ $t('contrats.gerer_abonnement') }}</a>
                </div>
                <p v-if="abonnement.abonnement === 'premium'">
                    <span class="badge-premium">Premium</span>
                    {{ $t('contrats.abonnement_premium', { montant: formaterPrix(abonnement.montant_mensuel || abonnement.prix_mensuel), date: abonnement.prochaine_echeance ? formaterDate(abonnement.prochaine_echeance) : '—' }) }}
                </p>
                <p v-else>{{ $t('contrats.abonnement_freemium') }}</p>
            </section>

            <!-- Campagnes publicitaires -->
            <section class="carte section-categories">
                <h2>{{ $t('contrats.campagnes.titre') }}</h2>
                <p class="discret">{{ $t('contrats.campagnes.explication') }}</p>
                <p v-if="campagnes.length === 0" class="discret">{{ $t('contrats.campagnes.aucune') }}</p>
                <table v-else class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('contrats.campagnes.element') }}</th>
                            <th class="montant">{{ $t('contrats.campagnes.budget') }}</th>
                            <th>{{ $t('commun.statut') }}</th>
                            <th>{{ $t('commun.date_debut') }}</th>
                            <th>{{ $t('commun.actions') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="c in campagnes" :key="c.id">
                            <td>
                                <span class="badge-statut badge-document">{{ $t('contrats.campagnes.cible.' + c.cible_type) }}</span>
                                {{ c.cible_titre || '—' }}
                            </td>
                            <td class="montant">{{ formaterPrix(c.montant_mensuel) }} / {{ $t('compte.abonnement.mois') }}</td>
                            <td><span class="badge-statut" :class="'contrat-' + c.statut">{{ $t('contrats.campagnes.statut.' + c.statut) }}</span></td>
                            <td class="date">{{ formaterDate(c.date_debut) }}</td>
                            <td>
                                <button v-if="c.statut === 'active'" class="lien-action lien-supprimer" @click="campagneAArreter = c">{{ $t('contrats.campagnes.arreter') }}</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </section>

            <!-- Factures et contrats archivés -->
            <section class="carte section-categories">
                <h2>{{ $t('contrats.documents.titre') }}</h2>
                <p v-if="documents.length === 0" class="discret">{{ $t('contrats.documents.aucun') }}</p>
                <ul v-else class="liste-documents">
                    <li v-for="d in documents" :key="d.id">
                        <img src="../assets/img/documents.svg" alt="" aria-hidden="true">
                        <div class="infos-document">
                            <strong>{{ d.titre }}</strong>
                            <span class="discret">{{ $t('admin.documents.type.' + d.type_document) }} · {{ formaterDate(d.date_creation) }}</span>
                        </div>
                        <div class="actions-document">
                            <button class="lien-action" @click="ouvrirDocument(d)">{{ $t('admin.documents.ouvrir') }}</button>
                            <button class="lien-action" @click="telechargerDocument(d)">{{ $t('admin.documents.telecharger') }}</button>
                        </div>
                    </li>
                </ul>
            </section>

            <!-- Paiements -->
            <section class="carte section-categories">
                <h2>{{ $t('contrats.paiements.titre') }}</h2>
                <p v-if="paiements.length === 0" class="discret">{{ $t('contrats.paiements.aucun') }}</p>
                <table v-else class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('admin.documents.colonne.date') }}</th>
                            <th>{{ $t('commun.type') }}</th>
                            <th>{{ $t('contrats.paiements.libelle') }}</th>
                            <th class="montant">{{ $t('paiement.montant') }}</th>
                            <th>{{ $t('contrats.paiements.facture') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="p in paiements" :key="p.id">
                            <td class="date">{{ formaterDate(p.date_paiement) }}</td>
                            <td>{{ $t('admin.finances.source.' + p.source) }}</td>
                            <td>{{ p.libelle || '—' }}</td>
                            <td class="montant">{{ formaterPrix(p.montant) }}</td>
                            <td>
                                <button v-if="p.facture_id" class="lien-action" @click="ouvrirDocument({ id: p.facture_id })">{{ $t('admin.documents.ouvrir') }}</button>
                                <span v-else class="discret">—</span>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </section>

            <!-- Catalogue des prestations de la plateforme -->
            <section class="carte section-categories">
                <h2>{{ $t('contrats.prestations.titre') }}</h2>
                <p v-if="prestations.length === 0" class="discret">{{ $t('contrats.prestations.aucune') }}</p>
                <table v-else class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('commun.titre') }}</th>
                            <th>{{ $t('commun.categorie') }}</th>
                            <th class="montant">{{ $t('commun.tarif') }}</th>
                            <th>{{ $t('commun.actions') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="p in prestations" :key="p.id">
                            <td><strong>{{ p.titre }}</strong></td>
                            <td>{{ p.categorie }}</td>
                            <td class="montant">{{ formaterPrix(p.tarif) }}</td>
                            <td><a class="lien-action" :href="lienPDF(p.id)">{{ $t('contrats.prestations.fiche_pdf') }}</a></td>
                        </tr>
                    </tbody>
                </table>
            </section>

            <!-- Nouvelle campagne -->
            <fenetre-modale :ouverte="campagneOuverte" :titre="$t('contrats.campagnes.nouvelle')" @fermer="campagneOuverte = false">
                <form @submit.prevent="lancerCampagne">
                    <p class="discret">{{ $t('contrats.campagnes.aide') }}</p>
                    <div class="form-group">
                        <label for="cible-type">{{ $t('contrats.campagnes.type_element') }}</label>
                        <select id="cible-type" v-model="formulaireCampagne.cible_type" @change="formulaireCampagne.cible_id = null">
                            <option value="projet">{{ $t('contrats.campagnes.cible.projet') }}</option>
                            <option value="annonce">{{ $t('contrats.campagnes.cible.annonce') }}</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="cible">{{ $t('contrats.campagnes.element') }}</label>
                        <select id="cible" v-model="formulaireCampagne.cible_id" required>
                            <option :value="null" disabled>{{ $t('contrats.campagnes.choisir') }}</option>
                            <option v-for="e in elementsDisponibles" :key="e.id" :value="e.id">{{ e.titre }}</option>
                        </select>
                        <p v-if="elementsDisponibles.length === 0" class="discret">{{ $t('contrats.campagnes.aucun_element.' + formulaireCampagne.cible_type) }}</p>
                    </div>
                    <div class="form-group">
                        <label for="budget">{{ $t('contrats.campagnes.budget_mensuel', { budget: formaterPrix(formulaireCampagne.montant) }) }}</label>
                        <input id="budget" type="range" v-model.number="formulaireCampagne.montant" min="100" max="500" step="50">
                        <div class="bornes-budget"><span>{{ formaterPrix(100) }}</span><span>{{ formaterPrix(500) }}</span></div>
                    </div>
                    <p v-if="erreurCampagne" class="erreur">{{ erreurCampagne }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit" :disabled="achatEnCours">{{ $t('contrats.campagnes.payer', { montant: formaterPrix(formulaireCampagne.montant) }) }}</button>
                        <button type="button" class="button-draft" @click="campagneOuverte = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

            <fenetre-modale :ouverte="campagneAArreter !== null" :titre="$t('contrats.campagnes.arreter')" @fermer="campagneAArreter = null">
                <p>{{ $t('contrats.campagnes.arreter_confirmation', { element: campagneAArreter?.cible_titre || '' }) }}</p>
                <div class="form-buttons">
                    <button type="button" class="button-submit bouton-refus" @click="arreterCampagne">{{ $t('contrats.campagnes.arreter') }}</button>
                    <button type="button" class="button-draft" @click="campagneAArreter = null">{{ $t('commun.annuler') }}</button>
                </div>
            </fenetre-modale>
        </main>
    </div>

    <script src="<?= ressource("js/contrats.js") ?>"></script>
</body>

</html>
