<?php
require "inclus/espace.php";
$titrePage = "Dépôt en conteneur";
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
                <h1>Dépôt en conteneur</h1>
            </div>

            <div class="deux-colonnes">

                <section>
                    <h2>Suivi de mes dépôts</h2>
                    <p v-if="erreurListe" class="erreur">{{ erreurListe }}</p>
                    <p v-if="depots.length === 0" class="discret">Aucune demande de dépôt pour le moment.</p>

                    <div v-for="d in depots" :key="d.id" class="carte section">
                        <h3>{{ d.description_objet }}</h3>
                        <p class="discret">Conteneur : {{ nomConteneur(d.conteneur_id) }} · demandé le {{ formaterDate(d.date_demande) }}</p>

                        <ol class="suivi-etapes">
                            <li v-for="(etape, index) in etapes" :key="etape.statut"
                                :class="{ atteinte: index <= indexStatut(d.statut) }">
                                {{ etape.libelle }}
                            </li>
                        </ol>

                        <p v-if="d.code_ouverture">
                            Code d'ouverture du conteneur : <span class="code">{{ d.code_ouverture }}</span><br>
                            Code-barre de l'objet : <span class="code">{{ d.code_barre }}</span>
                        </p>
                        <p v-else class="discret">Les codes vous seront communiqués dès que la demande sera validée.</p>
                    </div>
                </section>

                <div class="carte">
                    <h2>Demander un dépôt</h2>
                    <p class="discret">Votre demande est validée par notre équipe, puis vous recevez un code pour ouvrir le conteneur.</p>

                    <form @submit.prevent="demanderDepot">
                        <div class="form-group">
                            <label for="description">Description de l'objet</label>
                            <textarea id="description" v-model="formulaire.description_objet" required
                                placeholder="Type d'objet, état, dimensions…"></textarea>
                        </div>

                        <div class="form-group">
                            <label for="conteneur">Conteneur souhaité</label>
                            <select id="conteneur" v-model="formulaire.conteneur_id" required>
                                <option :value="null" disabled>Choisir un conteneur</option>
                                <option v-for="c in conteneurs" :key="c.id" :value="c.id">{{ c.site }}</option>
                            </select>
                        </div>

                        <p v-if="erreur" class="erreur">{{ erreur }}</p>
                        <p v-if="confirmation" class="succes">Demande envoyée !</p>

                        <button type="submit" class="button-submit">Envoyer la demande</button>
                    </form>
                </div>

            </div>

        </main>
    </div>

    <script src="js/depots.js"></script>
</body>

</html>
