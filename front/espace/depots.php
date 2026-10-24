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
            </div>

            <div v-if="roleConnecte() === 'particulier'">
                <div class="deux-colonnes">

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

                    <div class="carte">
                        <h2>{{ $t('depots.demander') }}</h2>
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
                            <p v-if="confirmation" class="succes">{{ $t('depots.demande_envoyee') }}</p>

                            <button type="submit" class="button-submit">{{ $t('depots.envoyer') }}</button>
                        </form>
                    </div>

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
                <div class="mise-en-page-admin">
                    <section class="carte">
                        <h2>Scanner un objet en conteneur</h2>
                        <p>Entrez le code de l'objet pour vérifier sa disponibilité et le récupérer.</p>
                    
                        <form @submit.prevent="verifierCode">
                            <div class="form-group">
                                <label for="code_objet">Code de l'objet (Scan ou manuel)</label>
                                <input id="code_objet" type="text" v-model="codeObjet" required placeholder="Ex: OBJ-12345">
                            </div>
                            <button type="submit" class="btn-primary">Vérifier le code</button>
                        </form>

                        <p v-if="erreur" class="erreur" style="margin-top: 15px;">{{ erreur }}</p>
                    </section>

                    <section v-if="objetTrouve" class="carte" style="margin-top: 20px; background-color: #f9f9f9; border-left: 4px solid #4CAF50;">
                        <h2>Objet trouvé</h2>
                        <p><strong>Type :</strong> {{ objetTrouve.type }}</p>
                        <p><strong>Description :</strong> {{ objetTrouve.description }}</p>
                        <p><strong>Conteneur / Box :</strong> {{ objetTrouve.conteneur }}</p>
                    
                        <button @click="confirmerRecuperation" class="btn-primary" style="margin-top: 15px; background-color: #4CAF50;">Confirmer la récupération</button>
                        <p v-if="succes" style="color: green; margin-top: 10px; font-weight: bold;">{{ succes }}</p>
                    </section>
                </div>
            </div>

        </main>
    </div>

    <script src="js/depots.js"></script>
</body>

</html>
