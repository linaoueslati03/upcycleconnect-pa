<template>

    <main id="main-content">

        <h1>Créer un atelier</h1>

        <form>

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
            Atelier enregistré !
        </div>

    </main>

</template>

<script setup>

import { ref } from "vue";

const formulaire = ref({
    titre: "",
    description: "",
    date_debut: "",
    date_fin: "",
    lieu: ""
});

const confirmation = ref(false);

async function envoyerAtelier(statut) {

    try {

        const reponse = await fetch(
            "http://localhost:8081/api/ateliers",
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
                    statut: statut,
                    createur_id: 1

                })

            }
        );

        if (!reponse.ok) {

            throw new Error(
                "Erreur lors de la création de l'atelier"
            );

        }

        confirmation.value = true;

    } catch (erreur) {

        console.error(erreur);

    }

}

function enregistrerBrouillon() {
    envoyerAtelier("brouillon");
}

function soumettreValidation() {
    envoyerAtelier("en_attente");
}

</script>