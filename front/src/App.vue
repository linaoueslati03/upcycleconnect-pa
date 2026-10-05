<template>

    <aside>

        <button id="toggle-button" aria-label="Accueil" @click="page = 'planning'">

            <img src="../assets/LOGO-pa.png" alt="Logo">

        </button>

        <ul>

            <li>

                <button class="nav-button" @click="page = 'planning'">

                    <img src="../assets/people-svgrepo-com (1).svg" alt="Formations">

                    <span>Formations</span>

                </button>

            </li>

            <li>

                <button class="nav-button" @click="page = 'evenement'">

                    <img src="../assets/party-horn-svgrepo-com.svg" alt="Événement">

                    <span>Événement</span>

                </button>

            </li>

            <li>

                <button class="nav-button">

                    <img src="../assets/pencil-square-svgrepo-com.svg" alt="Ateliers">

                    <span>Ateliers</span>

                </button>

            </li>

            <li>

                <button class="nav-button" @click="page = 'articles'">

                    <img src="../assets/news-svgrepo-com.svg" alt="Articles">

                    <span>Articles</span>

                </button>

            </li>

        </ul>

        <button id="logout-button">

            <img src="../assets/logout-svgrepo-com.svg" alt="Logout">

        </button>

    </aside>

    <Planning v-if="page === 'planning'" />

    <Articles v-if="page === 'articles'" />

</template>

<script setup>

import { ref } from "vue";

import Planning from "./planning.vue";
import Articles from "./articles.vue";

const page = ref("planning");

</script>