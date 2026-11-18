<?php
require "inclus/espace.php";
$titrePage = "nav.depots";
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
                <h1>{{ $t('depots.titre') }}</h1>
                <button v-if="roleConnecte() === 'particulier'" class="btn-primary" @click="ouvrirDemande">+ {{ $t('depots.demander') }}</button>
                <button v-if="roleConnecte() === 'professionnel'" class="btn-primary" @click="ouvrirScanner">{{ $t('depots.pro.scanner') }}</button>
            </div>

            <div v-if="roleConnecte() === 'particulier'">
                <p v-if="confirmation && !demandeOuverte" class="succes">{{ $t('depots.demande_envoyee') }}</p>
                <div>

                    <section>
                        <h2>{{ $t('depots.suivi') }}</h2>
                        <p v-if="erreurListe" class="erreur">{{ erreurListe }}</p>
                        <p v-if="depots.length === 0" class="discret">{{ $t('depots.aucun') }}</p>

                        <div v-for="d in depots" :key="d.id" class="carte section">
                            <h3>{{ d.description_objet }}</h3>
                            <p class="discret">{{ $t('depots.conteneur_demande_le', { conteneur: nomConteneur(d.conteneur_id), date: formaterDate(d.date_demande) }) }}</p>
                            <p v-if="d.categorie_id || d.poids_kg" class="discret">
                                {{ nomMateriau(d.categorie_id) }}<template v-if="d.categorie_id && d.poids_kg"> · </template><template v-if="d.poids_kg">{{ formaterPoids(d.poids_kg) }}<template v-if="d.poids_estime"> ({{ $t('depots.poids_estime') }})</template></template>
                            </p>

                            <ol class="suivi-etapes">
                                <li v-for="(etape, index) in etapes" :key="etape"
                                    :class="{ atteinte: index <= indexStatut(d.statut) }">
                                    {{ $t('depot.' + etape) }}
                                </li>
                            </ol>

                            <template v-if="d.code_ouverture">
                                <p>{{ $t('depots.code_ouverture') }} : <span class="code">{{ d.code_ouverture }}</span></p>
                                <figure v-if="d.code_barre && d.statut !== 'recuperee'" class="code-barre">
                                    <svg v-code-barre="d.code_barre" role="img" :aria-label="$t('depots.code_barre') + ' ' + d.code_barre"></svg>
                                    <figcaption>{{ $t('depots.code_barre_aide') }}</figcaption>
                                </figure>
                                <p v-else>{{ $t('depots.code_barre') }} : <span class="code">{{ d.code_barre }}</span></p>
                            </template>
                            <p v-else class="discret">{{ $t('depots.codes_apres_validation') }}</p>
                        </div>
                    </section>

                    <fenetre-modale :ouverte="demandeOuverte" :titre="$t('depots.demander')" @fermer="demandeOuverte = false">
                        <p class="discret">{{ $t('depots.explication') }}</p>

                        <form @submit.prevent="demanderDepot">
                            <div class="form-group">
                                <label for="description">{{ $t('depots.description_objet') }}</label>
                                <textarea id="description" v-model="formulaire.description_objet" required
                                    :placeholder="$t('depots.placeholder_objet')"></textarea>
                            </div>

                            <div class="deux-champs">
                                <div class="form-group">
                                    <label for="materiau">{{ $t('depots.materiau') }}</label>
                                    <select id="materiau" v-model="formulaire.categorie_id" required>
                                        <option :value="null" disabled>{{ $t('depots.choisir_materiau') }}</option>
                                        <option v-for="c in categoriesActives" :key="c.id" :value="c.id">{{ nomMateriau(c.id) }}</option>
                                    </select>
                                </div>
                                <div class="form-group">
                                    <label for="poids">{{ $t('depots.poids') }}</label>
                                    <input id="poids" type="number" v-model="formulaire.poids_kg" min="0.1" max="1000" step="0.1" inputmode="decimal"
                                        :placeholder="$t('depots.poids_exemple')">
                                </div>
                            </div>
                            <p class="discret aide-champ">{{ $t('depots.poids_aide') }}</p>

                            <div class="form-group">
                                <label for="conteneur">{{ $t('depots.conteneur_souhaite') }}</label>
                                <select id="conteneur" v-model="formulaire.conteneur_id" required>
                                    <option :value="null" disabled>{{ $t('depots.choisir_conteneur') }}</option>
                                    <option v-for="c in conteneurs" :key="c.id" :value="c.id" :disabled="c.statut === 'maintenance'">
                                        {{ c.statut === 'maintenance' ? c.site + ' (' + $t('depots.conteneur_maintenance') + ')' : c.site }}
                                    </option>
                                </select>
                            </div>

                            <p v-if="erreur" class="erreur">{{ erreur }}</p>
                            <div class="form-buttons">
                                <button type="submit" class="button-submit">{{ $t('depots.envoyer') }}</button>
                                <button type="button" class="button-draft" @click="demandeOuverte = false">{{ $t('commun.annuler') }}</button>
                            </div>
                        </form>
                    </fenetre-modale>

                </div>
            </div>

            <div v-if="roleConnecte() === 'salarie'">
                <div class="filtres">
                    <button v-for="filtre in filtres" :key="filtre.valeur"
                        class="pill" :class="{ active: statutSelectionne === filtre.valeur }"
                        @click="statutSelectionne = filtre.valeur; chargerDepotsAValider()">
                        {{ $t(filtre.nom) }}
                    </button>
                </div>

                <p v-if="erreur" class="erreur">{{ erreur }}</p>
                <p v-if="!erreur && depots.length === 0">{{ $t('salaries.depots.aucun') }}</p>

                <div class="evenements">
                    <div v-for="d in depots" :key="d.id" class="evenement-row">
                        <span class="evenement-titre">{{ d.description_objet }}
                            <small v-if="d.categorie_id" class="discret">· {{ nomMateriau(d.categorie_id) }}</small>
                        </span>
                        <span class="evenement-date">{{ formaterDate(d.date_demande) }}</span>
                        <label v-if="d.statut === 'demande'" class="poids-validation">
                            <span>{{ $t('depots.poids_verifie') }}</span>
                            <input type="number" v-model="poidsCorriges[d.id]" min="0.1" max="1000" step="0.1" inputmode="decimal"
                                :aria-label="$t('depots.poids_verifie') + ' · ' + d.description_objet">
                        </label>
                        <span v-else-if="d.poids_kg" class="discret">{{ formaterPoids(d.poids_kg) }}</span>
                        <span class="badge-statut">{{ $t('depot.' + d.statut) }}</span>
                        <span v-if="d.code_ouverture">{{ $t('salaries.depots.code', { code: d.code_ouverture }) }}</span>
                        <button v-if="statutSuivant[d.statut]" class="btn-modifier" @click="avancer(d)">
                            {{ $t(libellesAction[d.statut]) }}
                        </button>
                    </div>
                </div>
            </div>

            <div v-if="roleConnecte() === 'professionnel'">
                <section class="carte">
                    <h2>{{ $t('depots.pro.titre') }}</h2>
                    <p>{{ $t('depots.pro.explication') }}</p>
                    <button class="btn-primary" @click="ouvrirScanner">{{ $t('depots.pro.scanner') }}</button>
                    <p v-if="succes && !scannerOuvert" class="succes">{{ succes }}</p>
                </section>

                <fenetre-modale :ouverte="scannerOuvert" :titre="$t('depots.pro.titre')" @fermer="scannerOuvert = false">
                    <form @submit.prevent="verifierCode">
                        <div class="form-group">
                            <label for="code_objet">{{ $t('depots.pro.code_barre') }}</label>
                            <input id="code_objet" type="text" v-model="codeObjet" required autocomplete="off"
                                :placeholder="$t('depots.pro.code_exemple')">
                        </div>
                        <div class="form-buttons">
                            <button type="submit" class="button-submit">{{ $t('depots.pro.verifier') }}</button>
                        </div>
                    </form>

                    <p v-if="erreur" class="erreur">{{ erreur }}</p>

                    <div v-if="objetTrouve" class="objet-trouve">
                        <h3>{{ $t('depots.pro.objet_trouve') }}</h3>
                        <dl class="fiche-profil">
                            <dt>{{ $t('commun.description') }}</dt><dd>{{ objetTrouve.description_objet }}</dd>
                            <dt>{{ $t('depots.materiau') }}</dt><dd>{{ objetTrouve.materiau || '—' }}</dd>
                            <dt>{{ $t('depots.poids') }}</dt><dd>{{ objetTrouve.poids_kg ? formaterPoids(objetTrouve.poids_kg) : '—' }}</dd>
                            <dt>{{ $t('depots.pro.conteneur') }}</dt><dd>{{ objetTrouve.site || '—' }}</dd>
                        </dl>
                        <button @click="confirmerRecuperation" class="button-submit" :disabled="recuperationEnCours">{{ $t('depots.pro.confirmer') }}</button>
                    </div>
                    <p v-if="succes" class="succes">{{ succes }}</p>
                </fenetre-modale>
            </div>

        </main>
    </div>

    <script src="https://cdnjs.cloudflare.com/ajax/libs/jsbarcode/3.11.6/JsBarcode.all.min.js"></script>
    <script src="<?= ressource("js/depots.js") ?>"></script>
</body>

</html>
