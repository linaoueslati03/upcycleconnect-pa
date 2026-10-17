<?php
require "inclus/espace.php";
$titrePage = "nav.moderation";
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
                <h1>{{ $t('admin.moderation.titre') }}</h1>
            </div>

            <p v-if="erreur" class="erreur">{{ erreur }}</p>

            <h2>{{ $t('admin.moderation.signalements') }}</h2>
            <div class="filtres">
                <button v-for="filtre in filtresSignalements" :key="filtre.valeur"
                    class="pill" :class="{ active: statutSignalements === filtre.valeur }"
                    @click="statutSignalements = filtre.valeur; chargerSignalements()">
                    {{ $t(filtre.nom) }}
                </button>
            </div>

            <p v-if="signalements.length === 0">{{ $t('admin.moderation.aucun_signalement') }}</p>
            <section v-else class="carte">
                <table class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('admin.moderation.date') }}</th>
                            <th>{{ $t('admin.moderation.sujet') }}</th>
                            <th>{{ $t('admin.moderation.message') }}</th>
                            <th>{{ $t('admin.moderation.auteur') }}</th>
                            <th>{{ $t('admin.moderation.signale_par') }}</th>
                            <th>{{ $t('admin.moderation.motif') }}</th>
                            <th>{{ $t('commun.statut') }}</th>
                            <th>{{ $t('commun.actions') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="sg in signalements" :key="sg.id">
                            <td>{{ formaterDate(sg.created_at) }}</td>
                            <td>{{ sg.titre_sujet }}</td>
                            <td :class="{ 'statut-suspendu': sg.statut_message === 'masque' }">{{ sg.contenu }}</td>
                            <td>{{ sg.auteur_message }}</td>
                            <td>{{ sg.signale_par }}</td>
                            <td>{{ sg.motif || '—' }}</td>
                            <td>{{ $t('admin.moderation.statut_signalement.' + sg.statut) }}</td>
                            <td>
                                <button v-if="sg.statut_message === 'visible'" class="lien-action lien-supprimer"
                                    @click="changerStatutMessage(sg.message_id, 'masque')">
                                    {{ $t('admin.moderation.masquer') }}
                                </button>
                                <button v-else class="lien-action" @click="changerStatutMessage(sg.message_id, 'visible')">
                                    {{ $t('admin.moderation.retablir') }}
                                </button>
                                <button v-if="sg.statut === 'a_traiter'" class="lien-action" @click="rejeterSignalement(sg.id)">
                                    {{ $t('admin.moderation.rejeter_signalement') }}
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </section>

            <h2>{{ $t('admin.moderation.sujets') }}</h2>
            <div class="filtres">
                <button v-for="filtre in filtresSujets" :key="filtre.valeur"
                    class="pill" :class="{ active: statutSujets === filtre.valeur }"
                    @click="statutSujets = filtre.valeur; chargerSujets()">
                    {{ $t(filtre.nom) }}
                </button>
            </div>

            <p v-if="sujets.length === 0">{{ $t('admin.moderation.aucun_sujet') }}</p>
            <section v-else class="carte">
                <table class="tableau">
                    <thead>
                        <tr>
                            <th>{{ $t('admin.moderation.date') }}</th>
                            <th>{{ $t('commun.titre') }}</th>
                            <th>{{ $t('admin.moderation.auteur') }}</th>
                            <th>{{ $t('commun.statut') }}</th>
                            <th>{{ $t('commun.actions') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="su in sujets" :key="su.id">
                            <td>{{ formaterDate(su.created_at) }}</td>
                            <td>{{ su.titre }}</td>
                            <td>{{ su.auteur }}</td>
                            <td>{{ $t('admin.moderation.statut_sujet.' + su.statut) }}</td>
                            <td>
                                <button v-if="su.statut !== 'valide'" class="lien-action" @click="changerStatutSujet(su.id, 'valide')">
                                    {{ $t('commun.valider') }}
                                </button>
                                <button v-if="su.statut !== 'rejete'" class="lien-action lien-supprimer" @click="changerStatutSujet(su.id, 'rejete')">
                                    {{ $t('commun.refuser') }}
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </section>

        </main>
    </div>

    <script src="js/moderation.js"></script>
</body>

</html>
