<?php
require "inclus/espace.php";
$titrePage = "nav.validations";
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
                <h1>{{ $t('admin.validations.titre') }}</h1>
            </div>

            <p class="discret">{{ $t('admin.validations.explication') }}</p>
            <p v-if="erreur" class="erreur">{{ erreur }}</p>
            <p v-if="succes" class="succes">{{ succes }}</p>

            <div class="filtres">
                <button v-for="filtre in filtres" :key="filtre.valeur"
                    class="pill" :class="{ active: filtreType === filtre.valeur }"
                    @click="filtreType = filtre.valeur">
                    {{ $t(filtre.nom) }} ({{ compter(filtre.valeur) }})
                </button>
            </div>

            <p v-if="offresFiltrees.length === 0">{{ $t('admin.validations.aucune') }}</p>

            <section v-for="o in offresFiltrees" :key="o.type + o.id" class="carte fiche-validation">
                <span class="badge-statut" :class="'type-' + o.type">{{ $t('admin.validations.type.' + o.type) }}</span>
                <h2>{{ o.titre }}</h2>
                <p class="discret">{{ $t('admin.validations.propose_par', { nom: o.createur, date: formaterDate(o.created_at) }) }}</p>

                <ul class="details-offre">
                    <li><strong>{{ $t('commun.date_debut') }} :</strong> {{ formaterDateHeure(o.date_debut) }}</li>
                    <li v-if="o.date_fin"><strong>{{ $t('commun.date_fin') }} :</strong> {{ formaterDateHeure(o.date_fin) }}</li>
                    <li><strong>{{ $t('commun.lieu') }} :</strong> {{ o.lieu || $t('commun.non_precisee') }}<template v-if="o.site"> ({{ o.site }})</template></li>
                    <li v-if="o.type === 'formation'"><strong>{{ $t('admin.validations.places') }} :</strong> {{ o.nb_places ?? '—' }}</li>
                    <li v-if="o.type === 'formation'"><strong>{{ $t('commun.tarif') }} :</strong> {{ o.tarif ? o.tarif + ' €' : $t('commun.gratuit') }}</li>
                </ul>

                <p class="description-offre">{{ o.description || '—' }}</p>

                <form v-if="refus.cle === o.type + o.id" @submit.prevent="decider(o, 'refuser')">
                    <div class="form-group">
                        <label :for="'motif-' + o.type + o.id">{{ $t('admin.validations.motif_refus') }}</label>
                        <input :id="'motif-' + o.type + o.id" type="text" v-model="refus.motif" required
                            :placeholder="$t('admin.validations.motif_exemple')">
                    </div>
                    <div class="boutons-formulaire">
                        <button type="submit" class="button-submit bouton-refus">{{ $t('admin.validations.confirmer_refus') }}</button>
                        <button type="button" class="lien-action" @click="annulerRefus">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
                <div v-else class="boutons-formulaire">
                    <button class="button-submit" @click="decider(o, 'valider')">{{ $t('commun.valider') }}</button>
                    <button class="lien-action lien-supprimer" @click="commencerRefus(o)">{{ $t('commun.refuser') }}</button>
                </div>
            </section>

        </main>
    </div>

    <script src="js/validations.js"></script>
</body>

</html>
