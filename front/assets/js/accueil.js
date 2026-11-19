// Page d'accueil publique : chiffres, dernières annonces et prochains rendez-vous (GET /api/accueil).
// Langues affichées en drapeaux (bibliothèque flag-icons) ; une langue sans drapeau connu garde son code.
// Un visiteur déjà connecté voit « Mon espace » à la place de « Se connecter » / « Créer un compte ».

const PAYS_DES_LANGUES = { fr: "fr", en: "gb", es: "es", de: "de", it: "it", pt: "pt", nl: "nl" };

demarrerApp({
    data() {
        const role = localStorage.getItem("token") ? roleConnecte() : null;
        return {
            role,
            espace: ESPACES[role] || "",
            langues: [],
            langue: langueCourante(),
            menuLangue: false,
            donnees: null,
            chargement: true,
            erreur: false,
            annee: new Date().getFullYear(),
            clesChiffres: ["membres", "objets_recuperes", "annonces_en_ligne", "rendez_vous_a_venir"],
            profils: ["particulier", "professionnel", "salarie"],
        };
    },

    computed: {
        // « Voir toutes les annonces » : la page Annonces si le rôle y a accès, sinon son espace ou la connexion
        lienAnnonces() {
            if (["particulier", "professionnel"].includes(this.role)) {
                return "/espace/annonces";
            }
            return this.espace || PAGE_CONNEXION;
        },
        lienCatalogue() {
            if (this.role === "particulier") {
                return "/espace/catalogue";
            }
            return this.espace || PAGE_CONNEXION;
        },
    },

    async mounted() {
        // La liste des langues se ferme quand on clique ailleurs
        document.addEventListener("click", (evenement) => {
            if (this.menuLangue && !this.$refs.choixLangue.contains(evenement.target)) {
                this.menuLangue = false;
            }
        });

        try {
            const [langues, accueil] = await Promise.all([appelerApi("/langues"), appelerApi("/accueil")]);
            if (langues.ok) {
                // Langue par défaut (français) en premier, comme sur la maquette
                this.langues = (await langues.json()).sort((a, b) =>
                    (b.code === LANGUE_PAR_DEFAUT) - (a.code === LANGUE_PAR_DEFAUT) || a.code.localeCompare(b.code));
            }
            if (!accueil.ok) {
                throw new Error();
            }
            this.donnees = await accueil.json();
        } catch (erreur) {
            this.erreur = true;
        }
        this.chargement = false;
    },

    methods: {
        // Changement de langue sans recharger la page : on charge les textes de la langue choisie,
        // puis on bascule dessus, avec un court fondu pour éviter le clignotement.
        async choisirLangueVitrine(code) {
            this.menuLangue = false;
            if (code === this.langue) {
                return;
            }
            document.body.classList.add("changement-langue");
            if (localStorage.getItem("token")) {
                appelerApi("/moi/langue", { method: "PUT", body: JSON.stringify({ code }) }).catch(() => {});
            }
            localStorage.setItem("langue", code);
            await Promise.all([chargerMessages(code), new Promise((fin) => setTimeout(fin, 150))]);
            i18n.global.locale.value = code;
            this.langue = code;
            traduireElementsStatiques();
            requestAnimationFrame(() => document.body.classList.remove("changement-langue"));
        },

        // Code de la langue → code du pays du drapeau (l'anglais n'a pas de pays « en »)
        pays(code) {
            return PAYS_DES_LANGUES[code] || "";
        },

        nombre(valeur) {
            return Number(valeur).toLocaleString(langueCourante());
        },

        prix(annonce) {
            if (annonce.type === "don" || !annonce.prix) {
                return t("commun.gratuit");
            }
            return Number(annonce.prix).toLocaleString(langueCourante(), { style: "currency", currency: "EUR" });
        },

        categorie(annonce) {
            if (!annonce.categorie_code) {
                return t("commun.non_precisee");
            }
            return t("categorie." + annonce.categorie_code, annonce.categorie_libelle || annonce.categorie_code);
        },

        date(valeur) {
            return new Date(valeur).toLocaleDateString(langueCourante(), { day: "numeric", month: "long", year: "numeric" });
        },

        jour(valeur) {
            return new Date(valeur).toLocaleDateString(langueCourante(), { day: "2-digit" });
        },

        mois(valeur) {
            return new Date(valeur).toLocaleDateString(langueCourante(), { month: "short" });
        },

        heure(valeur) {
            return new Date(valeur).toLocaleString(langueCourante(), {
                weekday: "long", hour: "2-digit", minute: "2-digit",
            });
        },
    },
});
