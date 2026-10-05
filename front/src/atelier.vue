<template>

    <main>

        <div class="page-header">
            <h1>Mes ateliers</h1>

            <button
                class="btn-primary"
                @click="creerAtelier"
            >
                + Créer un atelier
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

        <div class="ateliers">

            <div
                v-for="atelier in ateliersFiltres"
                :key="atelier.id"
                class="atelier-row"
            >

                <span class="atelier-titre">
                    {{ atelier.titre }}
                </span>

                <span class="atelier-date">
                    {{ formaterDate(atelier.date_debut) }}
                </span>

                <span
                    class="badge-statut"
                    :class="'statut-' + atelier.statut"
                >
                    {{ formaterStatut(atelier.statut) }}
                </span>

                <button
                    class="btn-modifier"
                    @click="modifierAtelier(atelier.id)"
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

const ateliers = ref([]);
const statutSelectionne = ref("tous");

const filtres = [
    { nom: "Tous", valeur: "tous" },
    { nom: "Brouillon", valeur: "brouillon" },
    { nom: "En attente", valeur: "en_attente" },
    { nom: "Publié", valeur: "publie" },
    { nom: "Annulé", valeur: "annule" },
    { nom: "Terminé", valeur: "termine" }
];

const ateliersFiltres = computed(() => {

    if (statutSelectionne.value === "tous") {
        return ateliers.value;
    }

    return ateliers.value.filter(
        atelier =>
            atelier.statut === statutSelectionne.value
    );

});

async function chargerAteliers() {

    try {

        const reponse = await fetch(
            "http://localhost:8081/api/ateliers"
        );

        if (!reponse.ok) {
            throw new Error(
                "Erreur lors du chargement des ateliers"
            );
        }

        ateliers.value = await reponse.json();

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

function creerAtelier() {
    emit("creer");
}

function modifierAtelier(id) {
    console.log("Modifier l'atelier :", id);
}

onMounted(() => {
    chargerAteliers();
});

</script>