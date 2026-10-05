<template>

    <main>

        <div class="page-header">

            <h1>Mes événements</h1>

            <button
                class="btn-primary"
                @click="creerEvenement"
            >
                + Créer un événement
            </button>

        </div>


        <div class="filtres">

            <button
                v-for="filtre in filtres"
                :key="filtre.valeur"
                class="pill"
                :class="{ active: statutSelectionne === filtre.valeur }"
                @click="statutSelectionne = filtre.valeur"
            >
                {{ filtre.nom }}
            </button>

        </div>


        <div class="evenements">

            <div
                v-for="evenement in evenementsFiltres"
                :key="evenement.id"
                class="evenement-row"
            >

                <span class="evenement-titre">
                    {{ evenement.titre }}
                </span>

                <span class="evenement-date">
                    {{ formaterDate(evenement.date_debut) }}
                </span>

                <span
                    class="badge-statut"
                    :class="'statut-' + evenement.statut"
                >
                    {{ formaterStatut(evenement.statut) }}
                </span>

                <button
                    class="btn-modifier"
                    @click="modifierEvenement(evenement.id)"
                >
                    Modifier
                </button>

            </div>

        </div>

    </main>

</template>


<script setup>

import { computed, onMounted, ref } from "vue";

const emit = defineEmits(["creer"]);


const evenements = ref([]);

const statutSelectionne = ref("tous");


const filtres = [
    {
        nom: "Tous",
        valeur: "tous"
    },
    {
        nom: "Brouillon",
        valeur: "brouillon"
    },
    {
        nom: "En attente",
        valeur: "en_attente"
    },
    {
        nom: "Publié",
        valeur: "publie"
    },
    {
        nom: "Annulé",
        valeur: "annule"
    },
    {
        nom: "Terminé",
        valeur: "termine"
    }
];


const evenementsFiltres = computed(() => {

    if (statutSelectionne.value === "tous") {

        return evenements.value;

    }

    return evenements.value.filter(
        evenement =>
            evenement.statut === statutSelectionne.value
    );

});


async function chargerEvenements() {

    try {

        const reponse = await fetch(
            "http://localhost:8081/api/evenements"
        );

        if (!reponse.ok) {

            throw new Error(
                "Erreur lors du chargement des événements"
            );

        }

        evenements.value = await reponse.json();

    } catch (erreur) {

        console.error(erreur);

    }

}


function formaterDate(date) {

    return new Date(date).toLocaleDateString(
        "fr-FR",
        {
            day: "numeric",
            month: "long",
            year: "numeric"
        }
    );

}


function formaterStatut(statut) {

    const statuts = {

        brouillon: "Brouillon",

        en_attente: "En attente",

        publie: "Publié",

        annule: "Annulé",

        termine: "Terminé"

    };

    return statuts[statut] || statut;

}


function creerEvenement() {

    emit("creer");

}


function modifierEvenement(id) {

    console.log("Modifier l'événement :", id);

}


onMounted(() => {

    chargerEvenements();

});

</script>