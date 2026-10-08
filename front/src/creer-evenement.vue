<template>

    <main id="main-content">

        <h1>Créer un événement</h1>


        <form @submit.prevent="enregistrer">


            <div class="form-group">

                <label for="titre">
                    Titre
                </label>

                <input
                    type="text"
                    id="titre"
                    v-model="formulaire.titre"
                    required
                >

            </div>


            <div class="form-group">

                <label for="description">
                    Description
                </label>

                <textarea
                    id="description"
                    v-model="formulaire.description"
                ></textarea>

            </div>


            <div class="form-group">

                <label for="date-debut">
                    Date de début
                </label>

                <input
                    type="datetime-local"
                    id="date-debut"
                    v-model="formulaire.date_debut"
                    required
                >

            </div>


            <div class="form-group">

                <label for="date-fin">
                    Date de fin
                </label>

                <input
                    type="datetime-local"
                    id="date-fin"
                    v-model="formulaire.date_fin"
                >

            </div>


            <div class="form-group">

                <label for="lieu">
                    Lieu
                </label>

                <input
                    type="text"
                    id="lieu"
                    v-model="formulaire.lieu"
                >

            </div>


            <div class="form-group">

                <label for="site">
                    Site
                </label>

                <select
                    id="site"
                    v-model="formulaire.site"
                >

                    <option value="">
                        Sélectionner un site
                    </option>

                    <option value="paris-10">
                        Paris 10e
                    </option>

                    <option value="paris-11">
                        Paris 11e
                    </option>

                    <option value="paris-13">
                        Paris 13e
                    </option>

                    <option value="montreuil">
                        Montreuil
                    </option>

                    <option value="suisse">
                        Suisse
                    </option>

                </select>

            </div>


            <div class="form-buttons">

                <button
                    type="button"
                    class="button-draft"
                    @click="enregistrerBrouillon"
                >
                    Enregistrer en brouillon
                </button>


                <button
                    type="button"
                    class="button-submit"
                    @click="soumettreValidation"
                >
                    Soumettre à validation
                </button>

            </div>


        </form>


        <div
            v-if="confirmation"
            id="confirmation"
        >
            Événement enregistré !
        </div>

    </main>

</template>


<script setup>

import { ref } from "vue";


const emit = defineEmits(["retour"]);


const formulaire = ref({

    titre: "",

    description: "",

    date_debut: "",

    date_fin: "",

    lieu: "",

    site: ""

});


const confirmation = ref(false);


async function envoyerEvenement(statut) {

    try {

        const reponse = await fetch(
            "http://localhost:8081/api/evenements",
            {
                method: "POST",

                headers: {
                    "Content-Type": "application/json"
                },

                body: JSON.stringify({

                    titre: formulaire.value.titre,

                    description: formulaire.value.description,

                    date_debut: formulaire.value.date_debut,

                    date_fin: formulaire.value.date_fin || null,

                    lieu: formulaire.value.lieu,

                    site: formulaire.value.site,

                    statut: statut,

                    createur_id: 1 //ça sera celui du user connecté mais pour l'instant on met 1 pck sinon matrche pas

                })

            }
        );


        if (!reponse.ok) {

            throw new Error(
                "Erreur lors de la création de l'événement"
            );

        }


        confirmation.value = true;

    } catch (erreur) {

        console.error(erreur);

    }

}


function enregistrerBrouillon() {

    envoyerEvenement("brouillon");

}


function soumettreValidation() {

    envoyerEvenement("en_attente");

}

</script>