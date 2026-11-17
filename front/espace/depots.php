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
                <button v-if="roleConnecte() === 'professionnel'" class="btn-primary" @click="ouvrirScanner">Scanner un objet</button>
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

                            <ol class="suivi-etapes">
                                <li v-for="(etape, index) in etapes" :key="etape"
                                    :class="{ atteinte: index <= indexStatut(d.statut) }">
                                    {{ $t('depot.' + etape) }}
                                </li>
                            </ol>

                            <p v-if="d.code_ouverture">
                                {{ $t('depots.code_ouverture') }} : <span class="code">{{ d.code_ouverture }}</span><br>
                                {{ $t('depots.code_barre') }} : <span class="code">{{ d.code_barre }}</span>
                            </p>
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
                        <span class="evenement-titre">{{ d.description_objet }}</span>
                        <span class="evenement-date">{{ formaterDate(d.date_demande) }}</span>
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
                    <h2>Scanner un objet en conteneur</h2>
                    <p>Entrez le code de l'objet pour vérifier sa disponibilité et le récupérer.</p>
                    <button class="btn-primary" @click="ouvrirScanner">Scanner un objet</button>
                    <p v-if="succes && !scannerOuvert" class="succes">{{ succes }}</p>
                </section>

                <fenetre-modale :ouverte="scannerOuvert" titre="Scanner un objet en conteneur" @fermer="scannerOuvert = false">
                    <form @submit.prevent="verifierCode">
                        <div class="form-group">
                            <label for="code_objet">Code de l'objet (Scan ou manuel)</label>
                            <input id="code_objet" type="text" v-model="codeObjet" required placeholder="Ex: OBJ-12345">
                        </div>
                        <div class="form-buttons">
                            <button type="submit" class="button-submit">Vérifier le code</button>
                        </div>
                    </form>

                    <p v-if="erreur" class="erreur">{{ erreur }}</p>

                    <div v-if="objetTrouve" class="objet-trouve">
                        <h3>Objet trouvé</h3>
                        <p><strong>Type :</strong> {{ objetTrouve.type }}</p>
                        <p><strong>Description :</strong> {{ objetTrouve.description }}</p>
                        <p><strong>Conteneur / Box :</strong> {{ objetTrouve.conteneur }}</p>
                        <button @click="confirmerRecuperation" class="button-submit">Confirmer la récupération</button>
                    </div>
                    <p v-if="succes" class="succes">{{ succes }}</p>
                </fenetre-modale>
            </div>

        </main>
    </div>

    <script src="<?= ressource("js/depots.js") ?>"></script>
</body>

</html>
