import argparse
import bisect
import math
import random
import sys
import unicodedata
from datetime import date, datetime, time, timedelta
from decimal import Decimal

import psycopg

from config import parametres_connexion
from profils import PART_SITE_HABITUEL, PROFILS, SITES, TYPES_DEMANDE, tirer_demande, tirer_poids

DOMAINE = "example.com"
HASH_MOT_DE_PASSE = "$2a$10$HDK2AXzOOGuvVHeeI759R.Xz2oREQp2BAQRL381FMBAvIrXZZGqGO"
JOURS_HISTORIQUE = 540
NB_SALARIES = 6

POINTS_INSCRIPTION = 5
POINTS_ANNONCE_CEDEE = 15
POINTS_DEPOT_RECUPERE = 10

PRENOMS = [
    "Emma", "Louise", "Jade", "Alice", "Chloé", "Lina", "Léa", "Manon", "Inès", "Camille",
    "Sarah", "Zoé", "Anna", "Julia", "Nora", "Yasmine", "Clara", "Maëlys", "Rose", "Lou",
    "Gabriel", "Léo", "Raphaël", "Louis", "Arthur", "Jules", "Adam", "Hugo", "Nathan", "Paul",
    "Tom", "Noah", "Sacha", "Mohamed", "Ethan", "Théo", "Victor", "Malo", "Karim", "Antoine",
]
NOMS = [
    "Martin", "Bernard", "Dubois", "Thomas", "Robert", "Richard", "Petit", "Durand", "Leroy", "Moreau",
    "Simon", "Laurent", "Lefebvre", "Michel", "Garcia", "David", "Bertrand", "Roux", "Vincent", "Fournier",
    "Morel", "Girard", "André", "Mercier", "Dupont", "Lambert", "Bonnet", "François", "Martinez", "Legrand",
    "Garnier", "Faure", "Rousseau", "Blanc", "Guerin", "Muller", "Henry", "Benali", "Nguyen", "Haddad",
]
TITRES = {
    "formation": [
        "Upcycling textile : les bases", "Restaurer un meuble ancien", "Fabriquer avec des palettes",
        "Soudure et récupération de métal", "Couture zéro déchet", "Luminaires en matériaux récupérés",
        "Réparer son petit électroménager", "Tapisserie de sièges", "Mosaïque en verre recyclé",
        "Initiation à l'ébénisterie durable", "Créer sa marque upcyclée", "Vendre ses créations upcyclées",
    ],
    "atelier": [
        "Réparer ses vêtements", "Bijoux en capsules", "Tote bag en tissu récupéré", "Cadres en bois flotté",
        "Pots de fleurs en pneus", "Lampe en bocal", "Patchwork de chutes", "Étagère en caisses de vin",
        "Customiser un jean", "Repair Café vélo",
    ],
    "evenement": [
        "Repair Café", "Troc de printemps", "Marché des créateurs upcycling", "Conférence économie circulaire",
        "Journée portes ouvertes", "Défi zéro déchet", "Expo des projets de la communauté", "Nettoyage de quartier",
    ],
}
LIEUX_EVENEMENT = ["Salle polyvalente", "Hall d'accueil", "Cour intérieure", "Salle de conférence"]
DESCRIPTIONS = {
    "formation": [
        "Session pratique en petit groupe, matériel fourni.",
        "Formation pour débutants, aucun prérequis.",
        "Session de perfectionnement animée avec un artisan partenaire.",
    ],
    "atelier": [
        "Atelier convivial de deux heures, ouvert à tous.",
        "Venez avec vos objets, on les transforme ensemble.",
        "Atelier découverte, outils et conseils sur place.",
    ],
    "evenement": [
        "Événement gratuit ouvert à toute la communauté.",
        "Rencontre entre particuliers, artisans et salariés d'UpcycleConnect.",
        "Animations, démonstrations et échanges autour de l'upcycling.",
    ],
    "annonce": [
        "Bon état général, à récupérer sur place.",
        "Quelques traces d'usage, idéal pour un projet d'upcycling.",
        "Lot complet, disponible rapidement.",
    ],
}
OBJETS = {
    "bois": ["Lot de planches en chêne", "Chaises bistrot à restaurer", "Palettes", "Portes anciennes", "Caisses de vin"],
    "metal": ["Radiateur en fonte", "Tôles de zinc", "Pieds de table en acier", "Jantes de vélo", "Chutes de cuivre"],
    "textile": ["Rouleaux de tissu ancien", "Chutes de cuir", "Jeans usés", "Rideaux en lin", "Sacs de laine"],
    "verre": ["Bocaux en verre", "Portes vitrées", "Bouteilles colorées", "Miroirs anciens", "Vitraux"],
    "plastique": ["Bidons propres", "Caisses en plastique", "Bouchons triés", "Bâches", "Jouets à réparer"],
}


def sans_accents(texte):
    texte = unicodedata.normalize("NFKD", texte)
    return "".join(c for c in texte if not unicodedata.combining(c)).replace(" ", "").replace("'", "").lower()


def ajouter_mois(jour, nb):
    mois = jour.month - 1 + nb
    return date(jour.year + mois // 12, mois % 12 + 1, min(jour.day, 28))


def argent(valeur):
    return Decimal(str(valeur)).quantize(Decimal("0.01"))


class Generateur:
    def __init__(self, rng, reference, nb_acteurs, categories, conteneurs, prix_premium, codes_existants):
        self.rng = rng
        self.reference = reference
        self.nb_acteurs = nb_acteurs
        self.categories = categories
        self.conteneurs = conteneurs
        self.prix_premium = prix_premium
        self.codes = set(codes_existants)
        self.utilisateurs = []
        self.salaries = []
        self.offres = {t: [] for t in ("formation", "atelier", "evenement")}
        self.dates_offres = {}
        self.inscriptions = []
        self.annonces = []
        self.abonnements = []
        self.depots = []
        self.paiements = []
        self.scores = []
        self.compteur_paiements = 0

    def debut_taux_actuel(self):
        debut_mois = self.reference.replace(day=1, hour=0, minute=0, second=0, microsecond=0)
        return datetime.combine(ajouter_mois(debut_mois.date(), -5), time())

    def taux_commission(self, quand):
        return Decimal("7.00") if quand < self.debut_taux_actuel() else Decimal("8.00")

    def nouvel_identifiant_paiement(self, prefixe):
        self.compteur_paiements += 1
        return f"pi_gen_{prefixe}_{self.compteur_paiements}"

    def nouveau_code(self, octets):
        while True:
            code = f"{self.rng.getrandbits(octets * 8):0{octets * 2}x}"
            if code not in self.codes:
                self.codes.add(code)
                return code

    def ajouter_score(self, utilisateur, delta, motif, quand):
        self.scores.append({"utilisateur": utilisateur, "delta": delta, "motif": motif, "date": quand})
        utilisateur["score"] += delta

    def creer_salaries(self):
        for i in range(NB_SALARIES):
            prenom, nom = self.rng.choice(PRENOMS), self.rng.choice(NOMS)
            salarie = {
                "role": "salarie",
                "prenom": prenom,
                "nom": nom,
                "email": f"{sans_accents(prenom)}.{sans_accents(nom)}.salarie{i + 1}@{DOMAINE}",
                "date_creation": self.reference - timedelta(days=JOURS_HISTORIQUE + 60),
                "statut": "actif",
                "score": 0,
                "abonnement": "freemium",
                "responsable": i < 2,
            }
            self.utilisateurs.append(salarie)
            self.salaries.append(salarie)

    def creer_offres(self):
        responsables = [s for s in self.salaries if s["responsable"]]
        intervalles = {"formation": 3.5, "atelier": 3.5, "evenement": 6.0}
        for type_offre, intervalle in intervalles.items():
            jour = self.reference - timedelta(days=JOURS_HISTORIQUE + 20)
            fin = self.reference + timedelta(days=45)
            while jour < fin:
                debut = datetime.combine(jour.date(), time(self.rng.choice([10, 14, 18])))
                offre = {
                    "type": type_offre,
                    "titre": self.rng.choice(TITRES[type_offre]),
                    "site": tirer_poids(self.rng, SITES),
                    "date_debut": debut,
                    "date_fin": debut + timedelta(hours=self.rng.choice([2, 3])),
                    "statut": "termine" if debut < self.reference else "publie",
                    "createur": self.rng.choice(self.salaries),
                    "valideur": self.rng.choice(responsables),
                    "inscrits": set(),
                    "nb_places": None,
                    "tarif": Decimal("0"),
                }
                if type_offre == "formation":
                    offre["nb_places"] = self.rng.choice([10, 12, 15, 20])
                    if self.rng.random() > 0.15:
                        offre["tarif"] = argent(self.rng.randrange(20, 101, 5))
                self.offres[type_offre].append(offre)
                jour += timedelta(days=intervalle * self.rng.uniform(0.6, 1.4))
            self.dates_offres[type_offre] = [o["date_debut"] for o in self.offres[type_offre]]

    def creer_acteurs(self):
        noms_profils = list(PROFILS)
        parts = [PROFILS[p]["part"] for p in noms_profils]
        for i in range(self.nb_acteurs):
            profil = self.rng.choices(noms_profils, weights=parts)[0]
            prenom, nom = self.rng.choice(PRENOMS), self.rng.choice(NOMS)
            jours = max(1.0, JOURS_HISTORIQUE * (1 - math.sqrt(self.rng.random())))
            self.utilisateurs.append({
                "role": PROFILS[profil]["role"],
                "profil": profil,
                "prenom": prenom,
                "nom": nom,
                "email": f"{sans_accents(prenom)}.{sans_accents(nom)}.{i + 1}@{DOMAINE}",
                "date_creation": self.reference - timedelta(days=jours),
                "statut": "suspendu" if self.rng.random() < 0.03 else "actif",
                "site": tirer_poids(self.rng, SITES),
                "categories": PROFILS[profil]["categories"],
                "score": 0,
                "abonnement": "freemium",
                "abonnement_en_cours": None,
            })
        self.pros = self.acteurs("professionnel")
        self.poids_pros = [3 if p["profil"] == "artisan_acheteur" else 1 for p in self.pros]
        self.vendeurs = sorted(self.acteurs("particulier"), key=lambda u: u["date_creation"])
        self.dates_vendeurs = [v["date_creation"] for v in self.vendeurs]
        cumul, total = [], 0.0
        for vendeur in self.vendeurs:
            total += PROFILS[vendeur["profil"]]["annonces_par_mois"]
            cumul.append(total)
        self.cumul_vendeurs = cumul

    def acteurs(self, role=None):
        return [u for u in self.utilisateurs if "profil" in u and (role is None or u["role"] == role)]

    def site_demande(self, acteur):
        if self.rng.random() < PART_SITE_HABITUEL:
            return acteur["site"]
        return tirer_poids(self.rng, SITES)

    def categorie(self, acteur):
        code = tirer_poids(self.rng, acteur["categories"])
        return code if code in self.categories else self.rng.choice(list(self.categories))

    def abonnement_actif(self, acteur, quand):
        abo = acteur["abonnement_en_cours"]
        if abo is None:
            return False
        return abo["date_fin"] is None or quand.date() < abo["date_fin"]

    def inscrire(self, acteur, type_offre, quand):
        dates = self.dates_offres[type_offre]
        debut = bisect.bisect_left(dates, quand + timedelta(days=1))
        fin = bisect.bisect_right(dates, quand + timedelta(days=45))
        candidates = [
            o for o in self.offres[type_offre][debut:fin]
            if id(acteur) not in o["inscrits"] and (o["nb_places"] is None or len(o["inscrits"]) < o["nb_places"])
        ]
        if not candidates:
            return False
        site = self.site_demande(acteur)
        du_site = [o for o in candidates if o["site"] == site]
        offre = self.rng.choice(du_site or candidates)
        offre["inscrits"].add(id(acteur))
        payant = offre["tarif"] > 0
        self.inscriptions.append({
            "utilisateur": acteur,
            "offre": offre,
            "date": quand,
            "statut_paiement": "paye" if payant else "gratuit",
            "montant": offre["tarif"],
        })
        if payant:
            self.paiements.append({
                "utilisateur": acteur, "source": type_offre, "reference": offre, "montant": offre["tarif"],
                "taux": None, "commission": None, "statut": "reussi",
                "stripe_id": self.nouvel_identifiant_paiement("ins"), "date": quand,
            })
        self.ajouter_score(acteur, POINTS_INSCRIPTION, "Inscription à une offre", quand)
        return True

    def deposer(self, acteur, quand):
        site = self.site_demande(acteur)
        conteneur = self.rng.choice(self.conteneurs.get(site) or [c for liste in self.conteneurs.values() for c in liste])
        depot = {
            "utilisateur": acteur,
            "conteneur": conteneur,
            "description": self.rng.choice(OBJETS[self.categorie(acteur)]),
            "date_demande": quand,
            "statut": "demande",
            "code_ouverture": None,
            "code_barre": None,
            "recuperateur": None,
            "date_validation": None,
            "date_recuperation": None,
        }
        if self.reference - quand > timedelta(days=20):
            validation = quand + timedelta(days=self.rng.uniform(1, 3))
            recuperation = min(validation + timedelta(days=self.rng.uniform(2, 15)), self.reference - timedelta(hours=1))
            depot.update({
                "statut": "recuperee",
                "code_ouverture": self.nouveau_code(3),
                "code_barre": self.nouveau_code(8),
                "recuperateur": self.rng.choices(self.pros, weights=self.poids_pros)[0],
                "date_validation": validation,
                "date_recuperation": recuperation,
            })
            self.ajouter_score(acteur, POINTS_DEPOT_RECUPERE, "Dépôt récupéré", recuperation)
        self.depots.append(depot)
        return True

    def acheter_annonce(self, acheteur, quand):
        nb = bisect.bisect_left(self.dates_vendeurs, quand - timedelta(days=2))
        if nb == 0 or (nb == 1 and self.vendeurs[0] is acheteur):
            return False
        vendeur = acheteur
        while vendeur is acheteur:
            vendeur = self.rng.choices(self.vendeurs[:nb], cum_weights=self.cumul_vendeurs[:nb])[0]
        code = self.categorie(acheteur)
        mini, maxi = (30, 400) if acheteur["role"] == "professionnel" else (10, 150)
        prix = argent(self.rng.randrange(mini, maxi + 1, 5))
        creation = max(vendeur["date_creation"] + timedelta(hours=1), quand - timedelta(days=self.rng.uniform(1, 20)))
        annonce = {
            "utilisateur": vendeur,
            "titre": self.rng.choice(OBJETS[code]),
            "type": "vente",
            "prix": prix,
            "categorie": code,
            "localisation": self.site_demande(acheteur),
            "statut": "cedee",
            "date_creation": creation,
        }
        self.annonces.append(annonce)
        taux = self.taux_commission(quand)
        self.paiements.append({
            "utilisateur": acheteur, "source": "vente", "reference": annonce, "montant": prix,
            "taux": taux, "commission": argent(prix * taux / 100), "statut": "reussi",
            "stripe_id": self.nouvel_identifiant_paiement("vente"), "date": quand,
        })
        self.ajouter_score(vendeur, POINTS_ANNONCE_CEDEE, "Annonce cédée", quand)
        return True

    def abonner(self, acteur, quand):
        debut = quand.date()
        fin = None
        if self.rng.random() < 0.35:
            fin_prevue = ajouter_mois(debut, self.rng.randint(2, 8))
            if fin_prevue < self.reference.date():
                fin = fin_prevue
        abonnement = {
            "utilisateur": acteur,
            "montant": self.prix_premium,
            "statut": "resilie" if fin else "actif",
            "date_debut": debut,
            "date_fin": fin,
        }
        self.abonnements.append(abonnement)
        acteur["abonnement_en_cours"] = abonnement
        if fin is None:
            acteur["abonnement"] = "premium"
        limite = min(datetime.combine(fin, time(23, 59)), self.reference) if fin else self.reference
        mois = 0
        while datetime.combine(ajouter_mois(debut, mois), time(9)) < limite:
            self.paiements.append({
                "utilisateur": acteur, "source": "abonnement", "reference": abonnement, "montant": self.prix_premium,
                "taux": None, "commission": None, "statut": "reussi",
                "stripe_id": self.nouvel_identifiant_paiement("abo"),
                "date": datetime.combine(ajouter_mois(debut, mois), time(9)),
            })
            mois += 1
        return True

    def realiser(self, acteur, type_demande, quand):
        if type_demande in ("formation", "atelier", "evenement"):
            return self.inscrire(acteur, type_demande, quand)
        if type_demande == "depot":
            return self.deposer(acteur, quand)
        if type_demande == "achat_annonce":
            return self.acheter_annonce(acteur, quand)
        return self.abonner(acteur, quand)

    def simuler_demandes(self, acteur):
        profil = acteur["profil"]
        par_jour = PROFILS[profil]["demandes_par_mois"] / 30
        quand = acteur["date_creation"] + timedelta(days=self.rng.expovariate(par_jour))
        derniere = None
        while quand < self.reference:
            permis = [t for t in TYPES_DEMANDE if not (t == "premium" and (
                acteur["role"] != "professionnel" or self.abonnement_actif(acteur, quand)))]
            type_demande = tirer_demande(self.rng, profil, derniere, permis)
            if self.realiser(acteur, type_demande, quand):
                derniere = type_demande
            quand += timedelta(days=self.rng.expovariate(par_jour))

    def publier_annonces(self, acteur):
        profil = PROFILS[acteur["profil"]]
        par_jour = profil["annonces_par_mois"] / 30
        quand = acteur["date_creation"] + timedelta(days=self.rng.expovariate(par_jour))
        while quand < self.reference:
            vente = self.rng.random() < profil["part_vente"]
            code = self.categorie(acteur)
            age = self.reference - quand
            if age < timedelta(days=3):
                statut = "en_attente"
            else:
                tirage = self.rng.random()
                if tirage < 0.05:
                    statut = "refusee"
                elif not vente and tirage < 0.55:
                    statut = "cedee"
                elif tirage < 0.65:
                    statut = "reservee"
                else:
                    statut = "en_ligne"
            self.annonces.append({
                "utilisateur": acteur,
                "titre": self.rng.choice(OBJETS[code]),
                "type": "vente" if vente else "don",
                "prix": argent(self.rng.randrange(10, 200, 5)) if vente else None,
                "categorie": code,
                "localisation": self.site_demande(acteur),
                "statut": statut,
                "date_creation": quand,
            })
            if statut == "cedee":
                cession = min(quand + timedelta(days=self.rng.uniform(3, 30)), self.reference - timedelta(hours=1))
                self.ajouter_score(acteur, POINTS_ANNONCE_CEDEE, "Annonce cédée", cession)
            quand += timedelta(days=self.rng.expovariate(par_jour))

    def generer(self):
        self.creer_salaries()
        self.creer_offres()
        self.creer_acteurs()
        for acteur in self.acteurs():
            self.simuler_demandes(acteur)
        for acteur in self.acteurs():
            self.publier_annonces(acteur)


def reserver_identifiants(cur, table, objets):
    if not objets:
        return
    cur.execute(
        "SELECT nextval(pg_get_serial_sequence(%s, 'id')) FROM generate_series(1, %s)",
        (table, len(objets)),
    )
    for objet, (identifiant,) in zip(objets, cur.fetchall()):
        objet["id"] = identifiant


def copier(cur, table, colonnes, lignes):
    with cur.copy(f"COPY {table} ({', '.join(colonnes)}) FROM STDIN") as copie:
        for ligne in lignes:
            copie.write_row(ligne)


SOUS_REQUETE_GENERES = f"SELECT id FROM utilisateurs WHERE email LIKE '%%@{DOMAINE}'"


def supprimer_donnees_generees(cur):
    g = SOUS_REQUETE_GENERES
    requetes = [
        f"DELETE FROM paiements WHERE utilisateur_id IN ({g})",
        f"DELETE FROM score_historique WHERE utilisateur_id IN ({g})",
        f"""DELETE FROM inscriptions WHERE utilisateur_id IN ({g})
            OR (type_offre = 'formation' AND offre_id IN (SELECT id FROM formations WHERE createur_id IN ({g})))
            OR (type_offre = 'atelier' AND offre_id IN (SELECT id FROM ateliers WHERE createur_id IN ({g})))
            OR (type_offre = 'evenement' AND offre_id IN (SELECT id FROM evenements WHERE createur_id IN ({g})))""",
        f"DELETE FROM depots WHERE utilisateur_id IN ({g})",
        f"UPDATE depots SET professionnel_recuperateur_id = NULL WHERE professionnel_recuperateur_id IN ({g})",
        f"DELETE FROM abonnements WHERE utilisateur_id IN ({g})",
        f"DELETE FROM annonces WHERE utilisateur_id IN ({g})",
        f"DELETE FROM formations WHERE createur_id IN ({g})",
        f"DELETE FROM ateliers WHERE createur_id IN ({g})",
        f"DELETE FROM evenements WHERE createur_id IN ({g})",
        f"UPDATE formations SET valide_par_id = NULL WHERE valide_par_id IN ({g})",
        f"UPDATE ateliers SET responsable_id = NULL WHERE responsable_id IN ({g})",
        f"UPDATE evenements SET valide_par_id = NULL WHERE valide_par_id IN ({g})",
        f"DELETE FROM utilisateurs WHERE id IN ({g})",
    ]
    for requete in requetes:
        cur.execute(requete)


def charger_referentiels(cur):
    cur.execute("SELECT code, id FROM roles")
    roles = dict(cur.fetchall())
    cur.execute("SELECT id FROM langues WHERE code = 'fr'")
    langue = cur.fetchone()
    cur.execute("SELECT count(*) FROM categories_materiaux")
    if cur.fetchone()[0] == 0:
        cur.execute("""INSERT INTO categories_materiaux (code, libelle) VALUES
            ('bois', 'Bois'), ('metal', 'Métal'), ('textile', 'Textile'), ('verre', 'Verre'), ('plastique', 'Plastique')""")
    cur.execute("SELECT code, id FROM categories_materiaux ORDER BY id")
    categories = {code: identifiant for code, identifiant in cur.fetchall() if code in OBJETS}
    if not categories:
        sys.exit("Aucune catégorie de matériaux connue (bois, metal, textile, verre, plastique) en base.")
    cur.execute("SELECT id, normaliser_site(site) FROM conteneurs ORDER BY id")
    conteneurs = {}
    for identifiant, site in cur.fetchall():
        conteneurs.setdefault(site, []).append(identifiant)
    for site in SITES:
        if site not in conteneurs:
            cur.execute("INSERT INTO conteneurs (site) VALUES (%s), (%s) RETURNING id", (site, site))
            conteneurs[site] = [r[0] for r in cur.fetchall()]
    conteneurs = {site: ids for site, ids in conteneurs.items() if site in SITES}
    cur.execute("SELECT prix_premium_mensuel FROM parametres_finances WHERE id = 1")
    ligne = cur.fetchone()
    prix_premium = argent(ligne[0]) if ligne else argent(20)
    cur.execute("SELECT code_ouverture FROM depots WHERE code_ouverture IS NOT NULL "
                "UNION SELECT code_barre FROM depots WHERE code_barre IS NOT NULL")
    codes = [r[0] for r in cur.fetchall()]
    return roles, (langue[0] if langue else None), categories, conteneurs, prix_premium, codes


def ecrire(cur, gen, roles, langue, categories):
    reserver_identifiants(cur, "utilisateurs", gen.utilisateurs)
    for type_offre in ("formation", "atelier", "evenement"):
        reserver_identifiants(cur, f"{type_offre}s", gen.offres[type_offre])
    reserver_identifiants(cur, "annonces", gen.annonces)
    reserver_identifiants(cur, "abonnements", gen.abonnements)

    copier(cur, "utilisateurs",
           ["id", "role_id", "nom", "prenom", "email", "mot_de_passe_hash", "statut", "langue_preferee_id",
            "upcycling_score", "tutoriel_vu", "abonnement", "date_creation"],
           ([u["id"], roles[u["role"]], u["nom"], u["prenom"], u["email"], HASH_MOT_DE_PASSE, u["statut"], langue,
             u["score"], True, u["abonnement"], u["date_creation"]] for u in gen.utilisateurs))
    copier(cur, "salaries", ["utilisateur_id", "est_responsable", "created_at"],
           ([s["id"], s["responsable"], s["date_creation"]] for s in gen.salaries))

    def creation_offre(o):
        return o["date_debut"] - timedelta(days=30)

    copier(cur, "formations",
           ["id", "titre", "description", "date_debut", "date_fin", "lieu", "nb_places", "tarif", "statut",
            "createur_id", "valide_par_id", "created_at", "updated_at"],
           ([o["id"], o["titre"], gen.rng.choice(DESCRIPTIONS["formation"]), o["date_debut"], o["date_fin"], o["site"],
             o["nb_places"], o["tarif"], o["statut"], o["createur"]["id"], o["valideur"]["id"],
             creation_offre(o), creation_offre(o)] for o in gen.offres["formation"]))
    copier(cur, "ateliers",
           ["id", "titre", "description", "date_debut", "date_fin", "lieu", "statut", "createur_id",
            "responsable_id", "created_at", "updated_at"],
           ([o["id"], o["titre"], gen.rng.choice(DESCRIPTIONS["atelier"]), o["date_debut"], o["date_fin"], o["site"],
             o["statut"], o["createur"]["id"], o["valideur"]["id"], creation_offre(o), creation_offre(o)]
            for o in gen.offres["atelier"]))
    copier(cur, "evenements",
           ["id", "titre", "description", "date_debut", "date_fin", "lieu", "site", "statut", "createur_id",
            "valide_par_id", "created_at", "updated_at"],
           ([o["id"], o["titre"], gen.rng.choice(DESCRIPTIONS["evenement"]), o["date_debut"], o["date_fin"],
             gen.rng.choice(LIEUX_EVENEMENT), o["site"], o["statut"], o["createur"]["id"], o["valideur"]["id"],
             creation_offre(o), creation_offre(o)] for o in gen.offres["evenement"]))
    copier(cur, "annonces",
           ["id", "utilisateur_id", "titre", "description", "type", "prix", "categorie_id", "localisation", "statut",
            "date_creation"],
           ([a["id"], a["utilisateur"]["id"], a["titre"], gen.rng.choice(DESCRIPTIONS["annonce"]), a["type"],
             a["prix"], categories[a["categorie"]], a["localisation"], a["statut"], a["date_creation"]]
            for a in gen.annonces))
    copier(cur, "abonnements",
           ["id", "utilisateur_id", "montant_mensuel", "statut", "stripe_abonnement_id", "date_debut", "date_fin"],
           ([a["id"], a["utilisateur"]["id"], a["montant"], a["statut"], f"sub_gen_{a['id']}", a["date_debut"],
             a["date_fin"]] for a in gen.abonnements))
    copier(cur, "inscriptions",
           ["utilisateur_id", "type_offre", "offre_id", "date_inscription", "statut_paiement", "montant_paye"],
           ([i["utilisateur"]["id"], i["offre"]["type"], i["offre"]["id"], i["date"], i["statut_paiement"],
             i["montant"]] for i in gen.inscriptions))
    copier(cur, "depots",
           ["utilisateur_id", "conteneur_id", "description_objet", "code_ouverture", "code_barre", "statut",
            "professionnel_recuperateur_id", "date_demande", "date_validation", "date_recuperation"],
           ([d["utilisateur"]["id"], d["conteneur"], d["description"], d["code_ouverture"], d["code_barre"],
             d["statut"], d["recuperateur"]["id"] if d["recuperateur"] else None, d["date_demande"],
             d["date_validation"], d["date_recuperation"]] for d in gen.depots))
    copier(cur, "paiements",
           ["utilisateur_id", "source", "reference_id", "montant", "taux_commission", "montant_commission", "statut",
            "stripe_id", "date_paiement"],
           ([p["utilisateur"]["id"], p["source"], p["reference"]["id"], p["montant"], p["taux"], p["commission"],
             p["statut"], p["stripe_id"], p["date"]] for p in gen.paiements))
    copier(cur, "score_historique", ["utilisateur_id", "delta", "motif", "date"],
           ([s["utilisateur"]["id"], s["delta"], s["motif"], s["date"]] for s in gen.scores))


def afficher_resume(cur):
    cur.execute(f"""SELECT r.code, count(*) FROM utilisateurs u JOIN roles r ON r.id = u.role_id
                    WHERE u.id IN ({SOUS_REQUETE_GENERES}) GROUP BY r.code ORDER BY r.code""")
    print("Utilisateurs générés :", ", ".join(f"{code} {nb}" for code, nb in cur.fetchall()))
    cur.execute(f"""SELECT type_demande, count(*) FROM vue_demandes_acteurs
                    WHERE utilisateur_id IN ({SOUS_REQUETE_GENERES}) GROUP BY 1 ORDER BY 1""")
    print("Demandes générées :", ", ".join(f"{t} {nb}" for t, nb in cur.fetchall()))
    cur.execute("SELECT count(*) FROM vue_variables_acteurs")
    print("Acteurs dans vue_variables_acteurs :", cur.fetchone()[0])


def lire_arguments():
    parser = argparse.ArgumentParser(description="Génère des acteurs et leur historique pour le bloc Sciences.")
    parser.add_argument("--acteurs", type=int, default=600, help="nombre d'acteurs à générer (défaut : 600)")
    parser.add_argument("--graine", type=int, default=42, help="graine du générateur aléatoire (défaut : 42)")
    parser.add_argument("--date-reference", type=date.fromisoformat, default=date.today(),
                        help="date de fin de l'historique, AAAA-MM-JJ (défaut : aujourd'hui)")
    parser.add_argument("--reinitialiser", action="store_true",
                        help="supprime d'abord les données déjà générées par ce script")
    return parser.parse_args()


def main():
    args = lire_arguments()
    reference = datetime.combine(args.date_reference, time())
    with psycopg.connect(**parametres_connexion()) as connexion, connexion.cursor() as cur:
        cur.execute("SELECT to_regclass('vue_variables_acteurs') IS NOT NULL AND to_regclass('paiements') IS NOT NULL")
        if not cur.fetchone()[0]:
            sys.exit("La base n'a pas toutes les migrations (il faut au moins 022_reporting.sql).\n"
                     "Appliquez-les depuis la racine du projet : docker compose exec postgres sh /db/migrer.sh")
        cur.execute(f"SELECT count(*) FROM ({SOUS_REQUETE_GENERES}) g")
        if cur.fetchone()[0] > 0:
            if not args.reinitialiser:
                sys.exit("Des données générées existent déjà : relancez avec --reinitialiser pour les remplacer.")
            supprimer_donnees_generees(cur)
        roles, langue, categories, conteneurs, prix_premium, codes = charger_referentiels(cur)
        gen = Generateur(random.Random(args.graine), reference, args.acteurs, categories, conteneurs,
                         prix_premium, codes)
        gen.generer()
        ecrire(cur, gen, roles, langue, categories)
        afficher_resume(cur)
    print("Terminé.")


if __name__ == "__main__":
    main()
