<template>

    <div class="form-article">

        <div class="form-group">

            <label for="titre">Titre</label>

            <input
                id="titre"
                type="text"
                placeholder="Titre de l'article"
                v-model="titre"
            >

        </div>

        <div class="form-group">

            <label for="categorie">Catégorie</label>

            <select
                id="categorie"
                v-model="categorie"
            >

                <option value="">Sélectionner une catégorie</option>
                <option value="actualites">Actualités</option>
                <option value="conseils">Conseils</option>
                <option value="evenements">Événements</option>
                <option value="upcycling">Upcycling</option>

            </select>

        </div>

        <div class="form-group">

            <label for="texte">Article</label>

            <textarea
                id="texte"
                placeholder="Écrivez votre article..."
                v-model="texte"
            ></textarea>

        </div>

        <div>

            <button
                type="button"
                class="button-submit"
                @click="publier"
            >
                Publier
            </button>

        </div>

    </div>

</template>

<script setup>

import { ref } from "vue";

const titre = ref("");
const categorie = ref("");
const texte = ref("");

function publier() {

    console.log({
        titre: titre.value,
        categorie: categorie.value,
        texte: texte.value
    });

}

</script>