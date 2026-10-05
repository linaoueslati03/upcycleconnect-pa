<template>

    <main>

        <div class="page-header">

            <h1>Articles</h1>

            <button class="btn-primary" @click="afficherFormulaire">
                Créer un article
            </button>

        </div>

        <CreerArticle
            v-if="formulaireVisible"
            @fermer="formulaireVisible = false"
        />

    </main>

</template>

<script setup>

import { ref } from "vue";
import CreerArticle from "./creer-article.vue";

const formulaireVisible = ref(false);

function afficherFormulaire() {

    formulaireVisible.value = true;

}

</script>