<?php
require "inclus/espace.php";
$titrePage = "nav.mon_compte";
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
                <h1>{{ $t('nav.mon_compte') }}</h1>
            </div>

            <p v-if="motDePasseObligatoire" class="erreur">{{ $t('compte.changement_obligatoire') }}</p>

            <div class="deux-colonnes">

                <div class="carte">
                    <h2>{{ $t('compte.profil') }}</h2>
                    <dl class="fiche-profil">
                        <dt>{{ $t('commun.nom') }}</dt><dd>{{ profil.nom || '—' }}</dd>
                        <dt>{{ $t('commun.prenom') }}</dt><dd>{{ profil.prenom || '—' }}</dd>
                        <dt>{{ $t('commun.email') }}</dt><dd>{{ profil.email || '—' }}</dd>
                        <dt>{{ $t('compte.langue_affichage') }}</dt><dd>{{ libelleLangue || '—' }}</dd>
                    </dl>
                    <p v-if="erreurProfil && !profilOuvert" class="erreur">{{ erreurProfil }}</p>
                    <p v-if="messageProfil" class="succes">{{ messageProfil }}</p>
                    <button class="btn-primary" @click="ouvrirProfil">{{ $t('commun.modifier') }}</button>
                </div>

                <div>
                    <div class="carte section">
                        <h2>{{ $t('commun.mot_de_passe') }}</h2>
                        <p class="discret">{{ $t('commun.regle_mot_de_passe') }}</p>
                        <p v-if="messageMotDePasse" class="succes">{{ messageMotDePasse }}</p>
                        <button class="btn-primary" @click="ouvrirMotDePasse">{{ $t('compte.changer_mot_de_passe') }}</button>
                    </div>

                    <div class="carte section">
                        <h2>{{ $t('compte.notifications.titre') }}</h2>
                        <p>{{ $t('compte.notifications.explication') }}</p>
                        <button id="activer-notifications" type="button" class="btn-primary">
                            {{ $t('compte.notifications.activer') }}
                        </button>
                    </div>

                    <div class="carte section">
                        <h2>{{ $t('compte.tutoriel') }}</h2>
                        <p>{{ $t('compte.tutoriel_explication') }}</p>
                        <a class="btn-primary" :href="lienTutoriel()">{{ $t('compte.revoir_tutoriel') }}</a>
                    </div>

                    <div v-if="roleConnecte() === 'professionnel'" class="carte section">
                        <h2>{{ $t('compte.entreprise.titre') }}</h2>
                        <dl class="fiche-profil">
                            <dt>{{ $t('compte.entreprise.siret') }}</dt><dd>{{ entreprise.siret || '—' }}</dd>
                            <dt>{{ $t('compte.entreprise.adresse') }}</dt><dd>{{ entreprise.adresse || '—' }}</dd>
                        </dl>
                        <p v-if="messageEntreprise" class="succes">{{ messageEntreprise }}</p>
                        <button class="btn-primary" @click="ouvrirEntreprise">{{ $t('commun.modifier') }}</button>
                    </div>

                    <div v-if="roleConnecte() === 'professionnel' && abonnement" class="carte section carte-abonnement" :class="{ premium: abonnement.abonnement === 'premium' }">
                        <h2>{{ $t('compte.abonnement.titre') }}</h2>
                        <template v-if="abonnement.abonnement === 'premium'">
                            <p><span class="badge-premium">Premium</span></p>
                            <dl v-if="abonnement.date_debut" class="fiche-profil">
                                <dt>{{ $t('compte.abonnement.depuis') }}</dt><dd>{{ formaterDateDocument(abonnement.date_debut) }}</dd>
                                <dt>{{ $t('compte.abonnement.montant') }}</dt><dd>{{ formaterPrix(abonnement.montant_mensuel) }} / {{ $t('compte.abonnement.mois') }}</dd>
                                <dt>{{ $t('compte.abonnement.echeance') }}</dt><dd>{{ formaterDateDocument(abonnement.prochaine_echeance) }}</dd>
                            </dl>
                            <button class="button-draft" @click="resiliationOuverte = true">{{ $t('compte.abonnement.resilier') }}</button>
                        </template>
                        <template v-else>
                            <p>{{ $t('compte.abonnement.freemium') }}</p>
                            <ul class="avantages-premium">
                                <li v-for="a in ['statistiques', 'stocks', 'alertes']" :key="a">{{ $t('compte.abonnement.avantage.' + a) }}</li>
                            </ul>
                            <button class="btn-primary" :disabled="achatEnCours" @click="passerPremium">
                                {{ $t('compte.abonnement.passer', { prix: formaterPrix(abonnement.prix_mensuel) }) }}
                            </button>
                        </template>
                        <p v-if="messageAbonnement" class="succes">{{ messageAbonnement }}</p>
                        <p v-if="erreurAbonnement" class="erreur">{{ erreurAbonnement }}</p>
                    </div>
                </div>

            </div>

            <section class="carte mes-documents" id="mes-documents">
                <h2>{{ $t('compte.documents.titre') }}</h2>
                <p class="discret">{{ $t('compte.documents.explication') }}</p>
                <p v-if="erreurDocuments" class="erreur">{{ erreurDocuments }}</p>
                <p v-else-if="documents.length === 0" class="discret">{{ $t('compte.documents.aucun') }}</p>
                <ul v-else class="liste-documents">
                    <li v-for="d in documents" :key="d.id">
                        <img src="../assets/img/documents.svg" alt="" aria-hidden="true">
                        <div class="infos-document">
                            <strong>{{ d.titre }}</strong>
                            <span class="discret">{{ $t('admin.documents.type.' + d.type_document) }} · {{ formaterDateDocument(d.date_creation) }}</span>
                        </div>
                        <div class="actions-document">
                            <button class="lien-action" @click="ouvrirDocument(d)">{{ $t('admin.documents.ouvrir') }}</button>
                            <button class="lien-action" @click="telechargerDocument(d)">{{ $t('admin.documents.telecharger') }}</button>
                        </div>
                    </li>
                </ul>
            </section>

            <fenetre-modale :ouverte="profilOuvert" :titre="$t('compte.profil')" @fermer="profilOuvert = false">
                <form @submit.prevent="enregistrerProfil">
                    <div class="form-group">
                        <label for="nom">{{ $t('commun.nom') }}</label>
                        <input id="nom" type="text" v-model="brouillonProfil.nom" required>
                    </div>
                    <div class="form-group">
                        <label for="prenom">{{ $t('commun.prenom') }}</label>
                        <input id="prenom" type="text" v-model="brouillonProfil.prenom" required>
                    </div>
                    <div class="form-group">
                        <label for="email">{{ $t('commun.email') }}</label>
                        <input id="email" type="email" v-model="brouillonProfil.email" required>
                    </div>
                    <div class="form-group">
                        <label for="langue">{{ $t('compte.langue_affichage') }}</label>
                        <select id="langue" v-model="brouillonProfil.langue_preferee_id">
                            <option v-for="l in langues" :key="l.id" :value="l.id">{{ l.libelle }}</option>
                        </select>
                    </div>
                    <p v-if="erreurProfil" class="erreur">{{ erreurProfil }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                        <button type="button" class="button-draft" @click="profilOuvert = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

            <fenetre-modale :ouverte="entrepriseOuverte" :titre="$t('compte.entreprise.titre')" @fermer="entrepriseOuverte = false">
                <form @submit.prevent="enregistrerEntreprise">
                    <div class="form-group">
                        <label for="siret">{{ $t('compte.entreprise.siret') }}</label>
                        <input id="siret" type="text" v-model="brouillonEntreprise.siret" required inputmode="numeric"
                            pattern="[0-9 ]{14,17}" maxlength="17" :placeholder="$t('compte.entreprise.siret_exemple')">
                    </div>
                    <div class="form-group">
                        <label for="adresse">{{ $t('compte.entreprise.adresse') }}</label>
                        <input id="adresse" type="text" v-model="brouillonEntreprise.adresse" required maxlength="255" autocomplete="street-address">
                    </div>
                    <p v-if="erreurEntreprise" class="erreur">{{ erreurEntreprise }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('commun.enregistrer') }}</button>
                        <button type="button" class="button-draft" @click="entrepriseOuverte = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

            <fenetre-modale :ouverte="resiliationOuverte" :titre="$t('compte.abonnement.resilier')" @fermer="resiliationOuverte = false">
                <p>{{ $t('compte.abonnement.resilier_confirmation') }}</p>
                <div class="form-buttons">
                    <button type="button" class="button-submit bouton-refus" @click="resilier">{{ $t('compte.abonnement.resilier_valider') }}</button>
                    <button type="button" class="button-draft" @click="resiliationOuverte = false">{{ $t('commun.annuler') }}</button>
                </div>
            </fenetre-modale>

            <fenetre-modale :ouverte="motDePasseOuvert" :titre="$t('compte.changer_mot_de_passe')" @fermer="motDePasseOuvert = false">
                <form @submit.prevent="changerMotDePasse">
                    <p v-if="motDePasseObligatoire" class="erreur">{{ $t('compte.changement_obligatoire') }}</p>
                    <div class="form-group">
                        <label for="ancien">{{ $t('compte.mot_de_passe_actuel') }}</label>
                        <input id="ancien" type="password" v-model="motsDePasse.ancien_mot_de_passe"
                            autocomplete="current-password" required>
                    </div>
                    <div class="form-group">
                        <label for="nouveau">{{ $t('compte.nouveau_mot_de_passe') }}</label>
                        <input id="nouveau" type="password" v-model="motsDePasse.nouveau_mot_de_passe"
                            autocomplete="new-password" minlength="12" maxlength="72" required
                            aria-describedby="regle-mot-de-passe">
                        <small id="regle-mot-de-passe">{{ $t('commun.regle_mot_de_passe') }}</small>
                    </div>
                    <p v-if="erreurMotDePasse" class="erreur">{{ erreurMotDePasse }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit">{{ $t('compte.changer_mot_de_passe') }}</button>
                        <button type="button" class="button-draft" @click="motDePasseOuvert = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

        </main>
    </div>

    <script src="<?= ressource("js/mon-compte.js") ?>"></script>
</body>

</html>