<?php
require "inclus/espace.php";
$titrePage = "nav.documents";
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
                <h1>{{ $t('admin.documents.titre') }}</h1>
                <button v-if="onglet === 'documents'" class="btn-primary" @click="ouvrirNouveau">+ {{ $t('admin.documents.nouveau') }}</button>
            </div>

            <div class="filtres">
                <button v-for="o in onglets" :key="o" class="pill" :class="{ active: onglet === o }" @click="changerOnglet(o)">
                    {{ $t('admin.documents.onglet.' + o) }}
                </button>
            </div>

            <p v-if="succes" class="succes">{{ succes }}</p>
            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <!-- Documents PDF archivés -->
            <template v-if="onglet === 'documents'">
                <p class="discret">{{ $t('admin.documents.explication') }}</p>
                <div class="filtres-utilisateurs">
                    <input type="search" v-model="filtresDocuments.recherche" :placeholder="$t('admin.documents.rechercher')"
                        :aria-label="$t('admin.documents.rechercher')">
                </div>
                <div class="filtres">
                    <button class="pill" :class="{ active: filtresDocuments.type === '' }" @click="choisirType('')">
                        {{ $t('commun.tous') }} <span class="nombre-pill">{{ totalTousTypes }}</span>
                    </button>
                    <button v-for="c in documents.par_type" :key="c.libelle" class="pill"
                        :class="{ active: filtresDocuments.type === c.libelle }" @click="choisirType(c.libelle)">
                        {{ $t('admin.documents.type.' + c.libelle) }} <span class="nombre-pill">{{ c.nombre }}</span>
                    </button>
                </div>

                <section class="carte">
                    <p class="discret">{{ $t('admin.documents.resultats', { total: documents.total }) }}</p>
                    <p v-if="documents.total === 0 && !chargement" class="vide-documents">{{ $t('admin.documents.aucun') }}</p>
                    <table v-else class="tableau">
                        <thead>
                            <tr>
                                <th>{{ $t('admin.documents.colonne.date') }}</th>
                                <th>{{ $t('commun.type') }}</th>
                                <th>{{ $t('commun.titre') }}</th>
                                <th>{{ $t('admin.documents.colonne.beneficiaire') }}</th>
                                <th>{{ $t('admin.documents.colonne.auteur') }}</th>
                                <th class="montant">{{ $t('admin.documents.colonne.taille') }}</th>
                                <th>{{ $t('commun.actions') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="d in documents.documents" :key="d.id">
                                <td class="date">{{ formaterDate(d.date_creation) }}</td>
                                <td><span class="badge-statut badge-document" :class="'document-' + d.type_document">{{ $t('admin.documents.type.' + d.type_document) }}</span></td>
                                <td>{{ d.titre }}</td>
                                <td>{{ d.beneficiaire ? d.beneficiaire.nom : '—' }}</td>
                                <td>{{ d.auteur ? d.auteur.nom : '—' }}</td>
                                <td class="montant">{{ formaterTaille(d.taille_octets) }}</td>
                                <td class="actions-document">
                                    <button class="lien-action" @click="ouvrir(d)">{{ $t('admin.documents.ouvrir') }}</button>
                                    <button class="lien-action" @click="telecharger(d)">{{ $t('admin.documents.telecharger') }}</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                    <div v-if="pagesDocuments > 1" class="pagination">
                        <button class="lien-action" :disabled="filtresDocuments.page <= 1" @click="changerPageDocuments(-1)">← {{ $t('commun.precedent') }}</button>
                        <span>{{ $t('admin.utilisateurs.page', { page: filtresDocuments.page, pages: pagesDocuments }) }}</span>
                        <button class="lien-action" :disabled="filtresDocuments.page >= pagesDocuments" @click="changerPageDocuments(1)">{{ $t('commun.suivant') }} →</button>
                    </div>
                </section>
            </template>

            <!-- Codes de dépôt -->
            <template v-if="onglet === 'codes'">
                <p class="discret">{{ $t('admin.documents.codes.explication') }}</p>
                <div class="filtres-utilisateurs">
                    <input type="search" v-model="filtresCodes.recherche" :placeholder="$t('admin.documents.codes.rechercher')"
                        :aria-label="$t('admin.documents.codes.rechercher')">
                    <select v-model="filtresCodes.statut" :aria-label="$t('commun.statut')">
                        <option value="">{{ $t('admin.utilisateurs.tous_statuts') }}</option>
                        <option v-for="s in statutsCodes" :key="s" :value="s">{{ $t('admin.documents.codes.statut.' + s) }}</option>
                    </select>
                </div>

                <section class="carte">
                    <p class="discret">{{ $t('admin.documents.codes.resultats', { total: codes.total }) }}</p>
                    <table class="tableau">
                        <thead>
                            <tr>
                                <th>{{ $t('admin.documents.codes.objet') }}</th>
                                <th>{{ $t('admin.documents.codes.deposant') }}</th>
                                <th>{{ $t('commun.site') }}</th>
                                <th>{{ $t('depots.code_ouverture') }}</th>
                                <th>{{ $t('depots.code_barre') }}</th>
                                <th>{{ $t('commun.statut') }}</th>
                                <th>{{ $t('admin.documents.codes.recuperateur') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="c in codes.codes" :key="c.depot_id">
                                <td>{{ c.description_objet }}</td>
                                <td>{{ c.deposant.nom }}</td>
                                <td>{{ c.site || '—' }}</td>
                                <td><span class="code">{{ c.code_ouverture || '—' }}</span></td>
                                <td><span class="code">{{ c.code_barre || '—' }}</span></td>
                                <td>
                                    {{ $t('admin.documents.codes.statut.' + c.statut) }}
                                    <span v-if="c.date_recuperation" class="discret date-code">{{ formaterDate(c.date_recuperation) }}</span>
                                    <span v-else-if="c.date_validation" class="discret date-code">{{ formaterDate(c.date_validation) }}</span>
                                </td>
                                <td>{{ c.recuperateur ? c.recuperateur.nom : '—' }}</td>
                            </tr>
                        </tbody>
                    </table>
                    <div v-if="pagesCodes > 1" class="pagination">
                        <button class="lien-action" :disabled="filtresCodes.page <= 1" @click="changerPageCodes(-1)">← {{ $t('commun.precedent') }}</button>
                        <span>{{ $t('admin.utilisateurs.page', { page: filtresCodes.page, pages: pagesCodes }) }}</span>
                        <button class="lien-action" :disabled="filtresCodes.page >= pagesCodes" @click="changerPageCodes(1)">{{ $t('commun.suivant') }} →</button>
                    </div>
                </section>
            </template>

            <!-- Nouveau contrat ou document administratif -->
            <fenetre-modale :ouverte="formulaireOuvert" :titre="$t('admin.documents.nouveau')" large @fermer="formulaireOuvert = false">
                <form @submit.prevent="creer">
                    <div class="form-group">
                        <label for="type-document">{{ $t('commun.type') }}</label>
                        <select id="type-document" v-model="formulaire.type_document" required>
                            <option value="contrat">{{ $t('admin.documents.type.contrat') }}</option>
                            <option value="document_administratif">{{ $t('admin.documents.type.document_administratif') }}</option>
                        </select>
                    </div>

                    <div class="form-group choix-beneficiaire">
                        <label for="beneficiaire">{{ $t('admin.documents.colonne.beneficiaire') }}</label>
                        <div v-if="beneficiaire" class="beneficiaire-choisi">
                            <span><strong>{{ beneficiaire.prenom }} {{ beneficiaire.nom }}</strong> · {{ beneficiaire.email }}</span>
                            <button type="button" class="lien-action" @click="changerBeneficiaire">{{ $t('commun.modifier') }}</button>
                        </div>
                        <template v-else>
                            <input id="beneficiaire" type="search" v-model="rechercheBeneficiaire" autocomplete="off"
                                :placeholder="$t('admin.documents.chercher_beneficiaire')">
                            <ul v-if="propositions.length" class="propositions-beneficiaire">
                                <li v-for="u in propositions" :key="u.id">
                                    <button type="button" @click="choisirBeneficiaire(u)">
                                        <strong>{{ u.prenom }} {{ u.nom }}</strong>
                                        <span class="discret">{{ u.email }} · {{ $t('role.' + u.role_id) }}</span>
                                    </button>
                                </li>
                            </ul>
                            <p v-else-if="rechercheBeneficiaire.trim().length >= 2 && !rechercheEnCours" class="discret">{{ $t('admin.documents.aucun_beneficiaire') }}</p>
                        </template>
                    </div>

                    <div class="form-group">
                        <label for="titre-document">{{ $t('commun.titre') }}</label>
                        <input id="titre-document" type="text" v-model="formulaire.titre" maxlength="255" required>
                    </div>
                    <div class="form-group">
                        <label for="contenu-document">{{ $t('admin.documents.contenu') }}</label>
                        <textarea id="contenu-document" v-model="formulaire.contenu" rows="8" required
                            :placeholder="$t('admin.documents.contenu_aide')"></textarea>
                    </div>
                    <p v-if="erreurFormulaire" class="erreur">{{ erreurFormulaire }}</p>
                    <div class="form-buttons">
                        <button type="submit" class="button-submit" :disabled="enregistrement">{{ $t('admin.documents.generer') }}</button>
                        <button type="button" class="button-draft" @click="formulaireOuvert = false">{{ $t('commun.annuler') }}</button>
                    </div>
                </form>
            </fenetre-modale>

        </main>
    </div>

    <script src="<?= ressource("js/documents.js") ?>"></script>
</body>

</html>
