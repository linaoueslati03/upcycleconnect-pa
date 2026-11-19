<?php
require "inclus/espace.php";
$titrePage = "paiement.titre";
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
                <h1>{{ $t('paiement.titre') }}</h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <section v-if="commande" class="carte carte-paiement" :class="'paiement-' + commande.statut">

                <!-- En attente : récapitulatif et paiement -->
                <!-- Retour de Stripe : on attend la confirmation (webhook) -->
                <div v-if="commande.statut === 'en_attente' && commande.fournisseur !== 'simulation'" class="resultat-paiement">
                    <p>{{ $t('paiement.confirmation_en_cours') }}</p>
                </div>

                <template v-else-if="commande.statut === 'en_attente'">
                    <p class="bandeau-demo" role="note">
                        <strong>{{ $t('paiement.demo_titre') }}</strong> {{ $t('paiement.demo_texte') }}
                    </p>
                    <h2>{{ $t('paiement.recapitulatif') }}</h2>
                    <dl class="recapitulatif-paiement">
                        <dt>{{ $t('paiement.article') }}</dt><dd>{{ $t('paiement.type.' + commande.type) }} : {{ commande.libelle }}</dd>
                        <dt>{{ $t('paiement.montant') }}</dt><dd class="montant-paiement">{{ formaterPrix(commande.montant) }}</dd>
                    </dl>
                    <p class="discret">{{ $t('paiement.reserve.' + commande.type, { temps: tempsRestant }) }}</p>
                    <div class="form-buttons">
                        <button class="button-submit bouton-payer" :disabled="enCours" @click="payer">
                            {{ enCours ? $t('paiement.en_cours') : $t('paiement.payer', { montant: formaterPrix(commande.montant) }) }}
                        </button>
                        <button class="button-draft" :disabled="enCours" @click="annuler">{{ $t('commun.annuler') }}</button>
                    </div>
                </template>

                <!-- Payée -->
                <template v-else-if="commande.statut === 'payee'">
                    <div class="resultat-paiement">
                        <span class="pastille-resultat succes-paiement" aria-hidden="true">✓</span>
                        <h2>{{ $t('paiement.reussi_titre') }}</h2>
                        <p>{{ $t('paiement.reussi.' + commande.type, { libelle: commande.libelle }) }}</p>
                        <p class="discret">{{ commande.type === 'abonnement' || commande.type === 'publicite' ? $t('paiement.facture_contrat_disponibles') : $t('paiement.facture_disponible') }}</p>
                        <div class="form-buttons">
                            <button v-if="commande.facture_id" class="btn-primary" @click="ouvrirFacture">{{ $t('paiement.voir_facture') }}</button>
                            <a class="button-draft" href="mon-compte#mes-documents">{{ $t('compte.documents.titre') }}</a>
                            <a class="lien-action" :href="pageRetour">{{ $t('paiement.continuer.' + commande.type) }}</a>
                        </div>
                    </div>
                </template>

                <!-- Annulée ou expirée -->
                <template v-else>
                    <div class="resultat-paiement">
                        <span class="pastille-resultat" aria-hidden="true">✕</span>
                        <h2>{{ $t('paiement.' + commande.statut + '_titre') }}</h2>
                        <p>{{ $t('paiement.aucun_debit') }}</p>
                        <a class="btn-primary" :href="pageRetour">{{ $t('paiement.retour') }}</a>
                    </div>
                </template>

            </section>

        </main>
    </div>

    <script src="<?= ressource("js/paiement.js") ?>"></script>
</body>

</html>
