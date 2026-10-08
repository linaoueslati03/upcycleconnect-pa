<template>

    <aside>

        <button
            id="toggle-button"
            aria-label="Accueil"
            @click="page = 'planning'"
        >
            <img
                src="../assets/LOGO-pa.png"
                alt="Logo"
            >
        </button>

        <ul>

            <li>
                <button
                    class="nav-button"
                    @click="page = 'planning'"
                >
                    <img
                        src="../assets/people-svgrepo-com (1).svg"
                        alt="Formations"
                    >
                    <span>Formations</span>
                </button>
            </li>

            <li>
                <button
                    class="nav-button"
                    @click="page = 'evenement'"
                >
                    <img
                        src="../assets/party-horn-svgrepo-com.svg"
                        alt="Événement"
                    >
                    <span>Événement</span>
                </button>
            </li>

            <li>
                <button
                    class="nav-button"
                    @click="page = 'atelier'"
                >
                    <img
                        src="../assets/pencil-square-svgrepo-com.svg"
                        alt="Ateliers"
                    >
                    <span>Ateliers</span>
                </button>
            </li>

            <li>
                <button
                    class="nav-button"
                    @click="page = 'articles'"
                >
                    <img
                        src="../assets/news-svgrepo-com.svg"
                        alt="Articles"
                    >
                    <span>Articles</span>
                </button>
            </li>

        </ul>

        <button id="logout-button">
            <img
                src="../assets/logout-svgrepo-com.svg"
                alt="Logout"
            >
        </button>

    </aside>


    <Planning
        v-if="page === 'planning'"
    />


    <Evenements
        v-if="page === 'evenement'"
        @creer="page = 'creer-evenement'"
    />


    <CreerEvenement
        v-if="page === 'creer-evenement'"
        @retour="page = 'evenement'"
    />


    <Atelier
        v-if="page === 'atelier'"
        @creer="page = 'creer-atelier'"
    />


    <CreerAtelier
        v-if="page === 'creer-atelier'"
        @retour="page = 'atelier'"
    />


    <Articles
        v-if="page === 'articles'"
    />

</template>


<script setup>

import { ref } from "vue";

import Planning from "./planning.vue";
import Evenements from "./evenements.vue";
import CreerEvenement from "./creer-evenement.vue";
import Atelier from "./atelier.vue";
import CreerAtelier from "./creer-atelier.vue";
import Articles from "./articles.vue";

const page = ref("planning");

</script>