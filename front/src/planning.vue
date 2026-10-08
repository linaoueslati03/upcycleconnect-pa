<template>
    <main>

        <div class="planning-header">
            <h1>Mon planning</h1>

            <div class="planning-filtres">

                <div>
                    <label for="semaine">Semaine</label>

                    <select
                        id="semaine"
                        v-model="semaineSelectionnee"
                    >
                        <option
                            v-for="semaine in semaines"
                            :key="semaine.numero"
                            :value="semaine.numero"
                        >
                            Semaine du {{ formatDate(semaine.date) }}
                        </option>
                    </select>
                </div>

                <div>
                    <label for="annee">Année</label>

                    <select
                        id="annee"
                        v-model="anneeSelectionnee"
                    >
                        <option
                            v-for="annee in annees"
                            :key="annee"
                            :value="annee"
                        >
                            {{ annee }}
                        </option>
                    </select>
                </div>

            </div>
        </div>

        <div class="semaine-titre">

            <button @click="semainePrecedente">
                ←
            </button>

            <h2>
                Semaine du {{ formatDate(debutSemaine) }}
            </h2>

            <button @click="semaineSuivante">
                →
            </button>

        </div>

        <div class="planning">

            <div
                class="jour"
                v-for="jour in jours"
                :key="jour"
            >
                <h3>{{ jour }}</h3>

                <div class="evenements">
                </div>
            </div>

        </div>

    </main>
</template>

<script setup>
import { computed, ref, watch } from "vue";

const jours = [
    "Lundi",
    "Mardi",
    "Mercredi",
    "Jeudi",
    "Vendredi",
    "Samedi",
    "Dimanche"
];

const aujourdHui = new Date();

const semaineActuelle = getNumeroSemaine(aujourdHui);

const semaineSelectionnee = ref(semaineActuelle);

const anneeSelectionnee = ref(
    aujourdHui.getFullYear()
);

const annees = [];

for (
    let annee = aujourdHui.getFullYear() - 2;
    annee <= aujourdHui.getFullYear() + 2;
    annee++
) {
    annees.push(annee);
}

const semaines = computed(() => {
    const resultat = [];

    for (let i = 1; i <= 53; i++) {
        resultat.push({
            numero: i,
            date: getDateDepuisSemaine(
                i,
                anneeSelectionnee.value
            )
        });
    }

    return resultat;
});

const debutSemaine = computed(() => {
    return getDateDepuisSemaine(
        semaineSelectionnee.value,
        anneeSelectionnee.value
    );
});

function getNumeroSemaine(date) {
    const copie = new Date(
        Date.UTC(
            date.getFullYear(),
            date.getMonth(),
            date.getDate()
        )
    );

    const jour = copie.getUTCDay() || 7;

    copie.setUTCDate(
        copie.getUTCDate() + 4 - jour
    );

    const debutAnnee = new Date(
        Date.UTC(
            copie.getUTCFullYear(),
            0,
            1
        )
    );

    return Math.ceil(
        (((copie - debutAnnee) / 86400000) + 1) / 7
    );
}

function getDebutSemaine(date) {
    const copie = new Date(date);

    const jour = copie.getDay() || 7;

    copie.setDate(
        copie.getDate() - jour + 1
    );

    copie.setHours(0, 0, 0, 0);

    return copie;
}

function getDateDepuisSemaine(semaine, annee) {
    const date = new Date(annee, 0, 4);

    const debut = getDebutSemaine(date);

    debut.setDate(
        debut.getDate() + (semaine - 1) * 7
    );

    return debut;
}

function formatDate(date) {
    const jour = String(
        date.getDate()
    ).padStart(2, "0");

    const mois = String(
        date.getMonth() + 1
    ).padStart(2, "0");

    return `${jour}/${mois}`;
}

function semainePrecedente() {
    if (semaineSelectionnee.value > 1) {
        semaineSelectionnee.value--;
    } else {
        anneeSelectionnee.value--;
        semaineSelectionnee.value = 52;
    }
}

function semaineSuivante() {
    if (semaineSelectionnee.value < 53) {
        semaineSelectionnee.value++;
    } else {
        anneeSelectionnee.value++;
        semaineSelectionnee.value = 1;
    }
}

watch(
    anneeSelectionnee,
    () => {
        if (semaineSelectionnee.value > semaines.value.length) {
            semaineSelectionnee.value = 1;
        }
    }
);
</script>