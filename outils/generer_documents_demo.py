"""Génère de vrais PDF archivés pour la démonstration, en passant par les routes de l'API
(les mêmes que celles qu'utilise la plateforme), pour que les pages « Documents » (admin)
et « Mes documents » (Mon compte) ne soient pas vides en attendant la génération automatique.

Génère :
  - les fiches PDF de quelques prestations publiées (formations, ateliers, événements) ;
  - des attestations d'inscription (dont celles du compte particulier de test) ;
  - des contrats et documents administratifs pour les comptes de test et quelques professionnels.

Utilisation (API démarrée) :
    python3 outils/generer_documents_demo.py --email admin@upcycleconnect.fr
Le mot de passe est demandé au clavier (ou lu dans la variable d'environnement UPC_MOT_DE_PASSE).
Le script refuse de tourner si des documents existent déjà, sauf avec --forcer.

Rattrapage (une seule fois, après la mise en place de la génération automatique) :
    python3 outils/generer_documents_demo.py --email admin@upcycleconnect.fr --rattrapage
génère les PDF manquants pour l'existant : fiches des prestations publiées ou terminées,
attestations de toutes les inscriptions, bons des dépôts qui ont des codes. Relancé, il ne crée aucun doublon.
"""

import argparse
import getpass
import json
import os
import sys
import urllib.error
import urllib.parse
import urllib.request

NB_FICHES_PAR_TYPE = 3
NB_ATTESTATIONS = 12


class Api:
    def __init__(self, base):
        self.base = base.rstrip("/")
        self.jeton = None

    def appeler(self, methode, chemin, corps=None):
        donnees = json.dumps(corps).encode() if corps is not None else None
        requete = urllib.request.Request(self.base + chemin, data=donnees, method=methode)
        requete.add_header("Content-Type", "application/json")
        if self.jeton:
            requete.add_header("Authorization", self.jeton)
        try:
            with urllib.request.urlopen(requete) as reponse:
                return json.loads(reponse.read() or b"null")
        except urllib.error.HTTPError as erreur:
            message = erreur.read().decode(errors="replace")
            raise RuntimeError(f"{methode} {chemin} → {erreur.code} {message}") from None

    def connecter(self, email, mot_de_passe):
        self.jeton = self.appeler("POST", "/login", {"email": email, "mot_de_passe": mot_de_passe})["token"]


def chercher_utilisateur(api, recherche, role_id=None):
    parametres = {"recherche": recherche}
    if role_id:
        parametres["role_id"] = role_id
    resultat = api.appeler("GET", "/utilisateurs?" + urllib.parse.urlencode(parametres))["utilisateurs"]
    exact = [u for u in resultat if u["email"] == recherche]
    return (exact or resultat or [None])[0]


def generer_fiches(api):
    nb = 0
    for type_offre, chemin in (("formation", "/formations"), ("atelier", "/ateliers"), ("evenement", "/evenements")):
        offres = api.appeler("GET", chemin)
        if isinstance(offres, dict):
            offres = next((v for v in offres.values() if isinstance(v, list)), [])
        publiees = [o for o in offres if o.get("statut", "publie") == "publie"][:NB_FICHES_PAR_TYPE]
        for offre in publiees:
            api.appeler("POST", f"/documents/offres/{type_offre}/{offre['id']}")
            nb += 1
    return nb


def generer_attestations(api):
    # Aucune route ne liste les inscriptions : on essaie des numéros répartis sur la plage
    # (les inscriptions ont des identifiants séquentiels) et on ignore ceux qui n'existent pas.
    nb = 0
    for inscription_id in range(1, 5000, 37):
        if nb >= NB_ATTESTATIONS:
            break
        try:
            api.appeler("POST", f"/inscriptions/{inscription_id}/attestation")
            nb += 1
        except RuntimeError as erreur:
            if "→ 404" not in str(erreur):
                raise
    return nb


def generer_contrats(api, comptes):
    modeles = []
    pro = comptes.get("pro")
    if pro:
        modeles.append(("contrat", pro, "Contrat d'abonnement Premium", [
            "Entre UpcycleConnect, plateforme d'économie circulaire, et " + pro["prenom"] + " " + pro["nom"] + ".",
            "Objet : abonnement Premium mensuel donnant accès aux statistiques avancées et à la mise en avant des projets.",
            "Montant : 20,00 € par mois, prélevés à chaque date anniversaire.",
            "Durée : sans engagement, résiliable à tout moment depuis l'espace professionnel.",
            "Commission sur les ventes : selon le taux en vigueur affiché dans les conditions générales.",
        ]))
        modeles.append(("contrat", pro, "Contrat de récupération d'objets en conteneur", [
            "Le professionnel s'engage à récupérer les objets réservés dans un délai de 7 jours.",
            "La récupération se fait avec le code-barre fourni par UpcycleConnect.",
            "Tout objet non récupéré dans ce délai est remis à disposition des autres professionnels.",
        ]))
    particulier = comptes.get("particulier")
    if particulier:
        modeles.append(("document_administratif", particulier, "Confirmation d'adhésion à UpcycleConnect", [
            "Nous confirmons l'inscription de " + particulier["prenom"] + " " + particulier["nom"] + " sur UpcycleConnect.",
            "Le compte donne accès aux annonces, aux dépôts en conteneur, aux formations et à l'espace conseils.",
            "Les données personnelles sont traitées conformément au RGPD et peuvent être supprimées sur demande.",
        ]))
    salarie = comptes.get("salarie")
    if salarie:
        modeles.append(("contrat", salarie, "Convention d'animation d'ateliers", [
            "Le salarié anime les ateliers et formations qui lui sont attribués dans le planning.",
            "Chaque session est déclarée dans l'espace salarié et validée par un administrateur.",
            "Les supports pédagogiques restent la propriété d'UpcycleConnect.",
        ]))
    for pro_genere in comptes.get("autres_pros", []):
        modeles.append(("contrat", pro_genere, "Contrat d'abonnement Premium", [
            "Entre UpcycleConnect et " + pro_genere["prenom"] + " " + pro_genere["nom"] + ".",
            "Abonnement Premium mensuel de 20,00 €, sans engagement.",
        ]))

    for type_document, beneficiaire, titre, lignes in modeles:
        api.appeler("POST", "/admin/documents", {
            "type_document": type_document,
            "utilisateur_id": beneficiaire["id"],
            "titre": titre,
            "contenu": "\n".join(lignes),
        })
    return len(modeles)


def main():
    parser = argparse.ArgumentParser(description=__doc__.split("\n")[0])
    parser.add_argument("--email", required=True, help="e-mail d'un compte administrateur")
    parser.add_argument("--api", default="http://localhost:8081/api", help="adresse de l'API")
    parser.add_argument("--forcer", action="store_true", help="génère même si des documents existent déjà")
    parser.add_argument("--rattrapage", action="store_true",
                        help="génère seulement les PDF manquants de l'existant (fiches, attestations, bons de dépôt)")
    args = parser.parse_args()

    api = Api(args.api)
    mot_de_passe = os.environ.get("UPC_MOT_DE_PASSE") or getpass.getpass("Mot de passe administrateur : ")
    api.connecter(args.email, mot_de_passe)

    if args.rattrapage:
        bilan = api.appeler("POST", "/admin/documents/rattrapage")
        print(f"Fiches : {bilan['fiches']} | attestations : {bilan['attestations']} | "
              f"bons de dépôt : {bilan['bons_depot']} | erreurs : {bilan['erreurs']}")
        print("Total archivé :", api.appeler("GET", "/admin/documents")["total"])
        api.appeler("POST", "/logout")
        return

    existants = api.appeler("GET", "/admin/documents")["total"]
    if existants and not args.forcer:
        sys.exit(f"{existants} document(s) déjà archivé(s) : relancez avec --forcer pour en ajouter.")

    comptes = {
        "pro": chercher_utilisateur(api, "pro@upcycleconnect.fr"),
        "particulier": chercher_utilisateur(api, "particulier@upcycleconnect.fr"),
        "salarie": chercher_utilisateur(api, "salarie@upcycleconnect.fr"),
    }
    autres = api.appeler("GET", "/utilisateurs?role_id=2")["utilisateurs"]
    comptes["autres_pros"] = [u for u in autres if not u["email"].endswith("@upcycleconnect.fr")][:3]

    print("Fiches de prestation :", generer_fiches(api))
    print("Attestations d'inscription :", generer_attestations(api))
    print("Contrats et documents administratifs :", generer_contrats(api, comptes))
    print("Total archivé :", api.appeler("GET", "/admin/documents")["total"])
    api.appeler("POST", "/logout")


if __name__ == "__main__":
    main()
