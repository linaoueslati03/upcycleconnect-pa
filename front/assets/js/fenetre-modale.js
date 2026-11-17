// Fenêtre modale commune à tous les formulaires.
// Utilisation : <fenetre-modale :ouverte="formulaireOuvert" :titre="..." @fermer="formulaireOuvert = false"> <form>…</form> </fenetre-modale>
// Elle se ferme avec la croix, la touche Échap ou un clic sur le fond, bloque le défilement de la page
// derrière elle, place le curseur dans le premier champ et le rend à l'élément d'origine à la fermeture.

let fenetresOuvertes = 0;

const FenetreModale = {
    props: {
        ouverte: { type: Boolean, default: false },
        titre: { type: String, default: "" },
        large: { type: Boolean, default: false },
    },

    emits: ["fermer"],

    data() {
        return { elementPrecedent: null };
    },

    template: `
        <Teleport to="body">
            <Transition name="modale" @after-enter="placerCurseur">
                <div v-if="ouverte" class="modale-fond" @mousedown.self="fermer">
                    <div class="modale-boite" :class="{ 'modale-large': large }" role="dialog" aria-modal="true"
                        :aria-label="titre" ref="boite" @keydown="gererClavier">
                        <header class="modale-entete">
                            <h2>{{ titre }}</h2>
                            <button type="button" class="modale-fermer" :aria-label="$t('commun.fermer')" :title="$t('commun.fermer')" @click="fermer">
                                <span aria-hidden="true"></span>
                            </button>
                        </header>
                        <div class="modale-corps">
                            <slot></slot>
                        </div>
                    </div>
                </div>
            </Transition>
        </Teleport>
    `,

    watch: {
        ouverte: {
            immediate: true,
            handler(ouverte, avant) {
                if (ouverte && !avant) {
                    this.elementPrecedent = document.activeElement;
                    fenetresOuvertes++;
                } else if (!ouverte && avant) {
                    fenetresOuvertes = Math.max(0, fenetresOuvertes - 1);
                    if (this.elementPrecedent && document.contains(this.elementPrecedent)) {
                        this.elementPrecedent.focus({ preventScroll: true });
                    }
                }
                document.body.classList.toggle("modale-ouverte", fenetresOuvertes > 0);
            },
        },
    },

    unmounted() {
        if (this.ouverte) {
            fenetresOuvertes = Math.max(0, fenetresOuvertes - 1);
            document.body.classList.toggle("modale-ouverte", fenetresOuvertes > 0);
        }
    },

    methods: {
        fermer() {
            this.$emit("fermer");
        },

        placerCurseur() {
            const champ = this.$refs.boite?.querySelector(
                "input:not([type=hidden]):not([disabled]), select:not([disabled]), textarea:not([disabled])");
            (champ || this.$refs.boite)?.focus({ preventScroll: true });
        },

        gererClavier(evenement) {
            if (evenement.key === "Escape") {
                evenement.stopPropagation();
                this.fermer();
                return;
            }
            if (evenement.key !== "Tab") {
                return;
            }
            // Garde le focus clavier dans la fenêtre
            const focusables = [...this.$refs.boite.querySelectorAll(
                "a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled])")]
                .filter((e) => e.offsetParent !== null);
            if (focusables.length === 0) {
                return;
            }
            const premier = focusables[0];
            const dernier = focusables[focusables.length - 1];
            if (evenement.shiftKey && document.activeElement === premier) {
                dernier.focus();
                evenement.preventDefault();
            } else if (!evenement.shiftKey && document.activeElement === dernier) {
                premier.focus();
                evenement.preventDefault();
            }
        },
    },
};
