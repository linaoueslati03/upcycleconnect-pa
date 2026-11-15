<?php
require "inclus/espace.php";
$titrePage = "nav.notifications";
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
                <h1>{{ $t('admin.notifications.titre') }}</h1>
            </div>
            <p class="discret">{{ $t('admin.notifications.introduction') }}</p>

            <section class="carte section-categories">
                <h2>{{ $t('admin.notifications.nouvelle') }}</h2>
                <form class="formulaire-notification" @submit.prevent="envoyer">
                    <div class="form-group">
                        <label for="cible">{{ $t('admin.notifications.cible') }}</label>
                        <select id="cible" v-model="formulaire.cible" @change="changerCible">
                            <option v-for="c in cibles" :key="c" :value="c">{{ $t('admin.notifications.cibles.' + c) }}</option>
                        </select>
                    </div>

                    <div v-if="formulaire.cible === 'utilisateur'" class="form-group">
                        <label for="recherche-destinataire">{{ $t('admin.notifications.destinataire') }}</label>
                        <p v-if="destinataire" class="destinataire-choisi">
                            <strong>{{ destinataire.prenom }} {{ destinataire.nom }}</strong> · {{ destinataire.email }}
                            <button type="button" class="lien-action" @click="retirerDestinataire">{{ $t('admin.notifications.changer') }}</button>
                        </p>
                        <template v-else>
                            <input id="recherche-destinataire" type="search" v-model="recherche" @input="rechercherDestinataires"
                                :placeholder="$t('admin.notifications.rechercher')">
                            <ul v-if="resultatsRecherche.length" class="resultats-recherche">
                                <li v-for="u in resultatsRecherche" :key="u.id">
                                    <button type="button" class="lien-action" @click="choisirDestinataire(u)">
                                        {{ u.prenom }} {{ u.nom }} · {{ u.email }} ({{ $t('role.' + u.role_id) }})
                                    </button>
                                </li>
                            </ul>
                        </template>
                    </div>

                    <div v-if="formulaire.cible === 'prediction'" class="form-group">
                        <label for="classe">{{ $t('admin.notifications.classe') }}</label>
                        <select id="classe" v-model="formulaire.classe_predite" @change="actualiserApercu">
                            <option v-for="c in classes" :key="c" :value="c">{{ $t('admin.reporting.predictions.classe.' + c) }}</option>
                        </select>
                    </div>

                    <p class="apercu-destinataires" :class="{ vide: apercu === 0 }">
                        <template v-if="apercu === null">{{ $t('admin.notifications.apercu_attente') }}</template>
                        <template v-else>{{ $t('admin.notifications.apercu', { nombre: apercu }) }}</template>
                    </p>

                    <div class="champs-langues">
                        <fieldset>
                            <legend>{{ $t('admin.notifications.francais') }}</legend>
                            <div class="form-group">
                                <label for="titre-fr">{{ $t('commun.titre') }}</label>
                                <input id="titre-fr" type="text" v-model="formulaire.titre_fr" maxlength="100" required>
                            </div>
                            <div class="form-group">
                                <label for="message-fr">{{ $t('admin.notifications.message') }}</label>
                                <textarea id="message-fr" v-model="formulaire.message_fr" maxlength="500" rows="3" required></textarea>
                            </div>
                        </fieldset>
                        <fieldset>
                            <legend>{{ $t('admin.notifications.anglais') }}</legend>
                            <div class="form-group">
                                <label for="titre-en">{{ $t('commun.titre') }}</label>
                                <input id="titre-en" type="text" v-model="formulaire.titre_en" maxlength="100"
                                    :placeholder="$t('admin.notifications.repris_du_francais')">
                            </div>
                            <div class="form-group">
                                <label for="message-en">{{ $t('admin.notifications.message') }}</label>
                                <textarea id="message-en" v-model="formulaire.message_en" maxlength="500" rows="3"
                                    :placeholder="$t('admin.notifications.repris_du_francais')"></textarea>
                            </div>
                        </fieldset>
                    </div>

                    <p v-if="erreur" class="erreur">{{ erreur }}</p>
                    <p v-if="resultat" :class="resultat.statut === 'envoyee' ? 'succes' : 'erreur'">{{ messageResultat }}</p>

                    <div class="form-buttons">
                        <button type="submit" class="button-submit" :disabled="envoiEnCours || !apercu">
                            {{ envoiEnCours ? $t('admin.notifications.envoi_en_cours') : $t('admin.notifications.envoyer') }}
                        </button>
                    </div>
                </form>
            </section>

            <section class="carte section-categories">
                <h2>{{ $t('admin.notifications.historique') }}</h2>
                <p v-if="historique.notifications.length === 0">{{ $t('admin.notifications.aucune') }}</p>
                <table v-else class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('admin.notifications.colonne.date') }}</th>
                            <th>{{ $t('admin.notifications.colonne.notification') }}</th>
                            <th>{{ $t('admin.notifications.colonne.cible') }}</th>
                            <th class="montant">{{ $t('admin.notifications.colonne.destinataires') }}</th>
                            <th>{{ $t('commun.statut') }}</th>
                            <th>{{ $t('admin.notifications.colonne.auteur') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="n in historique.notifications" :key="n.id">
                            <td class="date">{{ formaterDate(n.date_envoi) }}</td>
                            <td><strong>{{ n.titre_fr }}</strong><br><span class="discret">{{ n.message_fr }}</span></td>
                            <td>{{ libelleCible(n) }}</td>
                            <td class="montant">{{ n.nb_destinataires }}</td>
                            <td>
                                <span class="badge-statut" :class="'notif-' + n.statut">{{ $t('admin.notifications.statut.' + n.statut) }}</span>
                                <br v-if="n.erreur"><span v-if="n.erreur" class="discret">{{ n.erreur }}</span>
                            </td>
                            <td>{{ n.auteur || '—' }}</td>
                        </tr>
                    </tbody>
                </table>

                <div class="pagination">
                    <button class="lien-action" :disabled="page <= 1" @click="chargerHistorique(page - 1)">← {{ $t('commun.precedent') }}</button>
                    <span>{{ $t('admin.notifications.page', { page, pages: nombrePages, total: historique.total }) }}</span>
                    <button class="lien-action" :disabled="page >= nombrePages" @click="chargerHistorique(page + 1)">{{ $t('commun.suivant') }} →</button>
                </div>
            </section>

        </main>
    </div>

    <script src="js/notifications.js"></script>
</body>

</html>
