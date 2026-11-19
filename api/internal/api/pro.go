package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Espace professionnel : entreprise, abonnement Premium, campagnes publicitaires (mise en avant),
// statistiques d'impact, alertes de collecte, récupération d'objets par code-barre.
// Les achats (abonnement, campagne) passent par les commandes (commandes.go), comme les autres paiements.

var formatSiret = regexp.MustCompile(`^[0-9]{14}$`)

// Un SIRET valide a 14 chiffres et respecte la formule de Luhn.
func siretValide(siret string) bool {
	if !formatSiret.MatchString(siret) {
		return false
	}
	somme := 0
	for i := 0; i < 14; i++ {
		chiffre := int(siret[13-i] - '0')
		if i%2 == 1 {
			chiffre *= 2
			if chiffre > 9 {
				chiffre -= 9
			}
		}
		somme += chiffre
	}
	return somme%10 == 0
}

// ---------- Entreprise ----------

type EntrepriseEntree struct {
	Siret   string `json:"siret"`
	Adresse string `json:"adresse"`
}

func (s *Serveur) gererModificationEntreprise(w http.ResponseWriter, r *http.Request) {
	utilisateurID, ok := s.exigerRole(w, r, roleProfessionnel)
	if !ok {
		return
	}
	var entree EntrepriseEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	entree.Siret = strings.ReplaceAll(strings.TrimSpace(entree.Siret), " ", "")
	entree.Adresse = strings.TrimSpace(entree.Adresse)
	if !siretValide(entree.Siret) {
		envoyerErreur(w, http.StatusBadRequest, "SIRET invalide (14 chiffres attendus)")
		return
	}
	if entree.Adresse == "" || len(entree.Adresse) > 255 {
		envoyerErreur(w, http.StatusBadRequest, "adresse de l'entreprise obligatoire (255 caractères maximum)")
		return
	}
	if _, err := s.db.Exec("UPDATE utilisateurs SET siret = $1, adresse_entreprise = $2 WHERE id = $3",
		entree.Siret, entree.Adresse, utilisateurID); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusOK, map[string]string{"siret": entree.Siret, "adresse_entreprise": entree.Adresse})
}

// ---------- Abonnement Premium ----------

type AbonnementPro struct {
	Abonnement        string     `json:"abonnement"`
	PrixMensuel       float64    `json:"prix_mensuel"`
	ID                *int       `json:"id"`
	MontantMensuel    *float64   `json:"montant_mensuel"`
	DateDebut         *time.Time `json:"date_debut"`
	ProchaineEcheance *time.Time `json:"prochaine_echeance"`
}

func (s *Serveur) gererMonAbonnement(w http.ResponseWriter, r *http.Request) {
	utilisateurID, ok := s.exigerRole(w, r, roleProfessionnel)
	if !ok {
		return
	}
	var a AbonnementPro
	err := s.db.QueryRow(`SELECT u.abonnement, p.prix_premium_mensuel FROM utilisateurs u, parametres_finances p
		WHERE u.id = $1 AND p.id = 1`, utilisateurID).Scan(&a.Abonnement, &a.PrixMensuel)
	if err == nil {
		err = s.db.QueryRow(`SELECT id, montant_mensuel, date_debut FROM abonnements
			WHERE utilisateur_id = $1 AND statut = 'actif'`, utilisateurID).Scan(&a.ID, &a.MontantMensuel, &a.DateDebut)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if a.DateDebut != nil {
		echeance := *a.DateDebut
		for !echeance.After(time.Now()) {
			echeance = echeance.AddDate(0, 1, 0)
		}
		a.ProchaineEcheance = &echeance
	}
	envoyerJSON(w, http.StatusOK, a)
}

func verifierAchatAbonnement(tx *sql.Tx, utilisateurID int) (string, float64, int, string) {
	var abonnement string
	var prix float64
	var actif bool
	err := tx.QueryRow(`
		SELECT u.abonnement, p.prix_premium_mensuel,
		       EXISTS (SELECT 1 FROM abonnements WHERE utilisateur_id = u.id AND statut = 'actif')
		FROM utilisateurs u, parametres_finances p WHERE u.id = $1 AND p.id = 1 FOR UPDATE OF u`,
		utilisateurID).Scan(&abonnement, &prix, &actif)
	switch {
	case err != nil:
		return "", 0, http.StatusInternalServerError, "erreur serveur"
	case abonnement == abonnementPremium || actif:
		return "", 0, http.StatusConflict, "vous êtes déjà abonné Premium"
	}
	return "Premium (1 mois)", prix, 0, ""
}

// Résiliation : prend effet immédiatement (sans Stripe, il n'y a pas de prélèvement mensuel à arrêter).
func (s *Serveur) gererResiliationAbonnement(w http.ResponseWriter, r *http.Request) {
	utilisateurID, ok := s.exigerRole(w, r, roleProfessionnel)
	if !ok {
		return
	}
	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()
	resultat, err := tx.Exec(`UPDATE abonnements SET statut = 'resilie', date_fin = CURRENT_DATE
		WHERE utilisateur_id = $1 AND statut = 'actif'`, utilisateurID)
	if err == nil {
		if n, _ := resultat.RowsAffected(); n == 0 {
			envoyerErreur(w, http.StatusConflict, "aucun abonnement Premium actif")
			return
		}
		_, err = tx.Exec("UPDATE utilisateurs SET abonnement = 'freemium' WHERE id = $1", utilisateurID)
	}
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	notifierAvecLien(utilisateurID, "Abonnement résilié", "Ton abonnement Premium a été résilié. Tu repasses en Freemium.", "mon-compte")
	w.WriteHeader(http.StatusNoContent)
}

// ---------- Campagnes publicitaires (mise en avant) ----------

const (
	budgetCampagneMin = 100
	budgetCampagneMax = 500
)

type CampagnePro struct {
	ID             int        `json:"id"`
	CibleType      string     `json:"cible_type"`
	CibleID        int        `json:"cible_id"`
	CibleTitre     *string    `json:"cible_titre"`
	MontantMensuel float64    `json:"montant_mensuel"`
	Statut         string     `json:"statut"`
	DateDebut      time.Time  `json:"date_debut"`
	DateFin        *time.Time `json:"date_fin"`
}

// Titre de l'élément à mettre en avant, s'il appartient bien au professionnel.
func titreCible(ex interface {
	QueryRow(string, ...any) *sql.Row
}, utilisateurID int, cibleType string, cibleID int) (string, int, string) {
	var titre string
	var err error
	switch cibleType {
	case "projet":
		err = ex.QueryRow("SELECT titre FROM projets_upcycling WHERE id = $1 AND utilisateur_id = $2 AND partage_public",
			cibleID, utilisateurID).Scan(&titre)
	case "annonce":
		err = ex.QueryRow(`SELECT titre FROM annonces WHERE id = $1 AND utilisateur_id = $2 AND statut = 'en_ligne'`,
			cibleID, utilisateurID).Scan(&titre)
	default:
		return "", http.StatusBadRequest, "type d'élément invalide (projet ou annonce attendu)"
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", http.StatusNotFound, "élément introuvable parmi vos projets publics ou vos annonces en ligne"
	}
	if err != nil {
		return "", http.StatusInternalServerError, "erreur serveur"
	}
	return titre, 0, ""
}

// Crée la campagne « en attente » le temps du paiement ; renvoie son identifiant.
func verifierAchatPublicite(tx *sql.Tx, utilisateurID int, entree CommandeEntree) (int, string, float64, int, string) {
	if entree.Montant < budgetCampagneMin || entree.Montant > budgetCampagneMax {
		return 0, "", 0, http.StatusBadRequest, fmt.Sprintf("budget mensuel entre %d et %d €", budgetCampagneMin, budgetCampagneMax)
	}
	titre, code, message := titreCible(tx, utilisateurID, entree.CibleType, entree.CibleID)
	if code != 0 {
		return 0, "", 0, code, message
	}
	var dejaEnCours bool
	if err := tx.QueryRow(`SELECT EXISTS (SELECT 1 FROM campagnes_publicitaires
		WHERE cible_type = $1 AND cible_id = $2 AND statut IN ('active', 'en_attente'))`,
		entree.CibleType, entree.CibleID).Scan(&dejaEnCours); err != nil {
		return 0, "", 0, http.StatusInternalServerError, "erreur serveur"
	}
	if dejaEnCours {
		return 0, "", 0, http.StatusConflict, "cet élément est déjà mis en avant"
	}
	var campagneID int
	if err := tx.QueryRow(`INSERT INTO campagnes_publicitaires (utilisateur_id, cible_type, cible_id, montant_mensuel, statut, date_debut)
		VALUES ($1, $2, $3, $4, 'en_attente', CURRENT_DATE) RETURNING id`,
		utilisateurID, entree.CibleType, entree.CibleID, entree.Montant).Scan(&campagneID); err != nil {
		return 0, "", 0, http.StatusInternalServerError, "erreur serveur"
	}
	return campagneID, titre, entree.Montant, 0, ""
}

func (s *Serveur) gererMesCampagnes(w http.ResponseWriter, r *http.Request) {
	utilisateurID, ok := s.exigerRole(w, r, roleProfessionnel)
	if !ok {
		return
	}
	if err := libererCommandesExpirees(s.db); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	lignes, err := s.db.Query(`
		SELECT c.id, c.cible_type, c.cible_id,
		       CASE c.cible_type WHEN 'projet' THEN (SELECT titre FROM projets_upcycling WHERE id = c.cible_id)
		                         ELSE (SELECT titre FROM annonces WHERE id = c.cible_id) END,
		       c.montant_mensuel, c.statut, c.date_debut, c.date_fin
		FROM campagnes_publicitaires c
		WHERE c.utilisateur_id = $1 AND c.statut <> 'en_attente'
		ORDER BY c.statut = 'active' DESC, c.date_debut DESC, c.id DESC`, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()
	campagnes := make([]CampagnePro, 0)
	for lignes.Next() {
		var c CampagnePro
		if err := lignes.Scan(&c.ID, &c.CibleType, &c.CibleID, &c.CibleTitre, &c.MontantMensuel, &c.Statut,
			&c.DateDebut, &c.DateFin); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		campagnes = append(campagnes, c)
	}
	envoyerJSON(w, http.StatusOK, campagnes)
}

func (s *Serveur) gererArretCampagne(w http.ResponseWriter, r *http.Request) {
	utilisateurID, ok := s.exigerRole(w, r, roleProfessionnel)
	if !ok {
		return
	}
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}
	resultat, err := s.db.Exec(`UPDATE campagnes_publicitaires SET statut = 'terminee', date_fin = CURRENT_DATE
		WHERE id = $1 AND utilisateur_id = $2 AND statut = 'active'`, id, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if n, _ := resultat.RowsAffected(); n == 0 {
		envoyerErreur(w, http.StatusNotFound, "campagne active introuvable")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Éléments mis en avant par une campagne active (affichés « À la une » aux particuliers et aux pros).
type ProjetALaUne struct {
	ID            int     `json:"id"`
	Titre         string  `json:"titre"`
	Description   *string `json:"description"`
	Professionnel string  `json:"professionnel"`
}

func (s *Serveur) gererALaUne(w http.ResponseWriter, r *http.Request) {
	if utilisateurID, _ := s.roleConnecte(r); utilisateurID == 0 {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}
	resultat := struct {
		Annonces []int          `json:"annonces"`
		Projets  []ProjetALaUne `json:"projets"`
	}{Annonces: make([]int, 0), Projets: make([]ProjetALaUne, 0)}

	lignes, err := s.db.Query(`SELECT cible_id FROM campagnes_publicitaires WHERE statut = 'active' AND cible_type = 'annonce'`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	for lignes.Next() {
		var id int
		if err := lignes.Scan(&id); err == nil {
			resultat.Annonces = append(resultat.Annonces, id)
		}
	}
	lignes.Close()

	lignes, err = s.db.Query(`
		SELECT p.id, p.titre, p.description, u.prenom || ' ' || u.nom
		FROM campagnes_publicitaires c
		JOIN projets_upcycling p ON p.id = c.cible_id
		JOIN utilisateurs u ON u.id = p.utilisateur_id
		WHERE c.statut = 'active' AND c.cible_type = 'projet' AND p.partage_public
		ORDER BY c.montant_mensuel DESC, c.date_debut DESC
		LIMIT 6`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()
	for lignes.Next() {
		var p ProjetALaUne
		if err := lignes.Scan(&p.ID, &p.Titre, &p.Description, &p.Professionnel); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		resultat.Projets = append(resultat.Projets, p)
	}
	envoyerJSON(w, http.StatusOK, resultat)
}

// ---------- Statistiques (Premium) ----------

type ImpactMateriau struct {
	Code    string   `json:"code"`
	Libelle string   `json:"libelle"`
	Nombre  int      `json:"nombre"`
	PoidsKg float64  `json:"poids_kg"`
	Co2Kg   float64  `json:"co2_kg"`
	Facteur *float64 `json:"facteur"`
	Source  *string  `json:"source"`
}

type ImpactMois struct {
	Mois    string  `json:"mois"`
	Nombre  int     `json:"nombre"`
	PoidsKg float64 `json:"poids_kg"`
	Co2Kg   float64 `json:"co2_kg"`
}

type DisponibleSite struct {
	Site     string `json:"site"`
	Materiau string `json:"materiau"`
	Nombre   int    `json:"nombre"`
}

// Objet déposé dans un conteneur et pas encore récupéré (le code-barre n'est pas montré :
// le professionnel le scanne sur place, sur le bon de dépôt).
type ObjetDisponible struct {
	ID               int       `json:"id"`
	DescriptionObjet string    `json:"description_objet"`
	Site             string    `json:"site"`
	Materiau         *string   `json:"materiau"`
	PoidsKg          *float64  `json:"poids_kg"`
	DateDepot        time.Time `json:"date_depot"`
}

type StatistiquesPro struct {
	Recuperations struct {
		Nombre       int     `json:"nombre"`
		PoidsKg      float64 `json:"poids_kg"`
		Co2Kg        float64 `json:"co2_kg"`
		PoidsEstimes int     `json:"poids_estimes"`
		SansPoids    int     `json:"sans_poids"`
	} `json:"recuperations"`
	ParMateriau []ImpactMateriau `json:"par_materiau"`
	ParMois     []ImpactMois     `json:"par_mois"`
	Achats      struct {
		Nombre   int     `json:"nombre"`
		Depenses float64 `json:"depenses"`
	} `json:"achats"`
	Disponibles       []DisponibleSite  `json:"disponibles"`
	ObjetsDisponibles []ObjetDisponible `json:"objets_disponibles"`
	CampagnesActives  int               `json:"campagnes_actives"`
}

func (s *Serveur) gererStatistiquesPro(w http.ResponseWriter, r *http.Request) {
	utilisateurID, ok := s.exigerPremium(w, r)
	if !ok {
		return
	}
	var st StatistiquesPro
	st.ParMateriau = make([]ImpactMateriau, 0)
	st.ParMois = make([]ImpactMois, 0)
	st.Disponibles = make([]DisponibleSite, 0)
	st.ObjetsDisponibles = make([]ObjetDisponible, 0)

	err := s.db.QueryRow(`
		SELECT count(*), COALESCE(sum(d.poids_kg), 0), COALESCE(sum(d.poids_kg * c.co2_evite_par_kg), 0),
		       count(*) FILTER (WHERE d.poids_estime), count(*) FILTER (WHERE d.poids_kg IS NULL)
		FROM depots d LEFT JOIN categories_materiaux c ON c.id = d.categorie_id
		WHERE d.professionnel_recuperateur_id = $1 AND d.statut = 'recuperee'`, utilisateurID,
	).Scan(&st.Recuperations.Nombre, &st.Recuperations.PoidsKg, &st.Recuperations.Co2Kg,
		&st.Recuperations.PoidsEstimes, &st.Recuperations.SansPoids)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	lignes, err := s.db.Query(`
		SELECT c.code, c.libelle, count(d.id), COALESCE(sum(d.poids_kg), 0),
		       COALESCE(sum(d.poids_kg * c.co2_evite_par_kg), 0), c.co2_evite_par_kg, c.source_facteur
		FROM categories_materiaux c
		LEFT JOIN depots d ON d.categorie_id = c.id AND d.professionnel_recuperateur_id = $1 AND d.statut = 'recuperee'
		GROUP BY c.id ORDER BY c.libelle`, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	for lignes.Next() {
		var m ImpactMateriau
		if err := lignes.Scan(&m.Code, &m.Libelle, &m.Nombre, &m.PoidsKg, &m.Co2Kg, &m.Facteur, &m.Source); err != nil {
			lignes.Close()
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		st.ParMateriau = append(st.ParMateriau, m)
	}
	lignes.Close()

	lignes, err = s.db.Query(`
		SELECT to_char(m.debut, 'YYYY-MM'), count(d.id), COALESCE(sum(d.poids_kg), 0),
		       COALESCE(sum(d.poids_kg * c.co2_evite_par_kg), 0)
		FROM generate_series(date_trunc('month', now()) - interval '11 months', date_trunc('month', now()), interval '1 month') AS m(debut)
		LEFT JOIN depots d ON d.professionnel_recuperateur_id = $1 AND d.statut = 'recuperee'
		     AND d.date_recuperation >= m.debut AND d.date_recuperation < m.debut + interval '1 month'
		LEFT JOIN categories_materiaux c ON c.id = d.categorie_id
		GROUP BY m.debut ORDER BY m.debut`, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	for lignes.Next() {
		var m ImpactMois
		if err := lignes.Scan(&m.Mois, &m.Nombre, &m.PoidsKg, &m.Co2Kg); err != nil {
			lignes.Close()
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		st.ParMois = append(st.ParMois, m)
	}
	lignes.Close()

	err = s.db.QueryRow(`SELECT count(*), COALESCE(sum(montant), 0) FROM paiements
		WHERE utilisateur_id = $1 AND source = 'vente' AND statut = 'reussi'`, utilisateurID,
	).Scan(&st.Achats.Nombre, &st.Achats.Depenses)
	if err == nil {
		err = s.db.QueryRow(`SELECT count(*) FROM campagnes_publicitaires WHERE utilisateur_id = $1 AND statut = 'active'`,
			utilisateurID).Scan(&st.CampagnesActives)
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	// Objets déposés, en attente d'un professionnel : l'état des stocks en temps réel
	lignes, err = s.db.Query(`
		SELECT normaliser_site(co.site), COALESCE(c.libelle, '—'), count(*)
		FROM depots d JOIN conteneurs co ON co.id = d.conteneur_id
		LEFT JOIN categories_materiaux c ON c.id = d.categorie_id
		WHERE d.statut = 'deposee'
		GROUP BY 1, 2 ORDER BY 1, 2`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	for lignes.Next() {
		var d DisponibleSite
		if err := lignes.Scan(&d.Site, &d.Materiau, &d.Nombre); err != nil {
			lignes.Close()
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		st.Disponibles = append(st.Disponibles, d)
	}
	lignes.Close()

	lignes, err = s.db.Query(`
		SELECT d.id, d.description_objet, normaliser_site(co.site), c.libelle, d.poids_kg,
		       COALESCE(d.date_validation, d.date_demande)
		FROM depots d JOIN conteneurs co ON co.id = d.conteneur_id
		LEFT JOIN categories_materiaux c ON c.id = d.categorie_id
		WHERE d.statut = 'deposee'
		ORDER BY 3, d.id DESC
		LIMIT 100`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()
	for lignes.Next() {
		var o ObjetDisponible
		if err := lignes.Scan(&o.ID, &o.DescriptionObjet, &o.Site, &o.Materiau, &o.PoidsKg, &o.DateDepot); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		st.ObjetsDisponibles = append(st.ObjetsDisponibles, o)
	}
	envoyerJSON(w, http.StatusOK, st)
}

// ---------- Alertes de collecte (Premium) ----------

type AlerteCollecte struct {
	ID           int       `json:"id"`
	CategorieID  *int      `json:"categorie_id"`
	Materiau     *string   `json:"materiau"`
	Site         *string   `json:"site"`
	MotCle       *string   `json:"mot_cle"`
	DateCreation time.Time `json:"date_creation"`
}

type AlerteEntree struct {
	CategorieID *int   `json:"categorie_id"`
	Site        string `json:"site"`
	MotCle      string `json:"mot_cle"`
}

const alertesMaximum = 10

func (s *Serveur) gererMesAlertes(w http.ResponseWriter, r *http.Request) {
	utilisateurID, ok := s.exigerPremium(w, r)
	if !ok {
		return
	}
	lignes, err := s.db.Query(`
		SELECT a.id, a.categorie_id, c.libelle, a.site, a.mot_cle, a.date_creation
		FROM alertes_collecte a LEFT JOIN categories_materiaux c ON c.id = a.categorie_id
		WHERE a.utilisateur_id = $1 ORDER BY a.date_creation DESC`, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()
	alertes := make([]AlerteCollecte, 0)
	for lignes.Next() {
		var a AlerteCollecte
		if err := lignes.Scan(&a.ID, &a.CategorieID, &a.Materiau, &a.Site, &a.MotCle, &a.DateCreation); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		alertes = append(alertes, a)
	}
	envoyerJSON(w, http.StatusOK, alertes)
}

func videVersNull(texte string) any {
	if texte = strings.TrimSpace(texte); texte == "" {
		return nil
	}
	return texte
}

func (s *Serveur) gererCreationAlerte(w http.ResponseWriter, r *http.Request) {
	utilisateurID, ok := s.exigerPremium(w, r)
	if !ok {
		return
	}
	var entree AlerteEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if entree.CategorieID == nil && strings.TrimSpace(entree.Site) == "" && strings.TrimSpace(entree.MotCle) == "" {
		envoyerErreur(w, http.StatusBadRequest, "choisissez au moins un matériau, un site ou un mot-clé")
		return
	}
	if len(entree.MotCle) > 100 {
		envoyerErreur(w, http.StatusBadRequest, "mot-clé trop long (100 caractères maximum)")
		return
	}
	if entree.CategorieID != nil {
		if _, message, err := s.verifierCategorie(tableCategoriesMateriaux, *entree.CategorieID, 0); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		} else if message != "" {
			envoyerErreur(w, http.StatusBadRequest, message)
			return
		}
	}
	var nombre int
	if err := s.db.QueryRow("SELECT count(*) FROM alertes_collecte WHERE utilisateur_id = $1", utilisateurID).Scan(&nombre); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if nombre >= alertesMaximum {
		envoyerErreur(w, http.StatusConflict, fmt.Sprintf("%d alertes maximum", alertesMaximum))
		return
	}
	var id int
	if err := s.db.QueryRow(`INSERT INTO alertes_collecte (utilisateur_id, categorie_id, site, mot_cle)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		utilisateurID, entree.CategorieID, videVersNull(entree.Site), videVersNull(entree.MotCle)).Scan(&id); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (s *Serveur) gererSuppressionAlerte(w http.ResponseWriter, r *http.Request) {
	utilisateurID, ok := s.exigerPremium(w, r)
	if !ok {
		return
	}
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}
	resultat, err := s.db.Exec("DELETE FROM alertes_collecte WHERE id = $1 AND utilisateur_id = $2", id, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if n, _ := resultat.RowsAffected(); n == 0 {
		envoyerErreur(w, http.StatusNotFound, "alerte introuvable")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Appelée quand un objet est déposé dans un conteneur : prévient les professionnels Premium concernés.
func (s *Serveur) declencherAlertesCollecte(d Depot) {
	lignes, err := s.db.Query(`
		SELECT DISTINCT a.utilisateur_id, COALESCE(co.site, ''), COALESCE(c.libelle, '')
		FROM depots dp
		LEFT JOIN conteneurs co ON co.id = dp.conteneur_id
		LEFT JOIN categories_materiaux c ON c.id = dp.categorie_id
		JOIN alertes_collecte a ON (a.categorie_id IS NULL OR a.categorie_id = dp.categorie_id)
		     AND (a.site IS NULL OR normaliser_site(a.site) = normaliser_site(co.site))
		     AND (a.mot_cle IS NULL OR dp.description_objet ILIKE '%' || a.mot_cle || '%')
		JOIN utilisateurs u ON u.id = a.utilisateur_id AND u.abonnement = 'premium' AND u.statut = 'actif'
		WHERE dp.id = $1`, d.ID)
	if err != nil {
		log.Printf("alertes de collecte du dépôt %d : %v", d.ID, err)
		return
	}
	defer lignes.Close()
	for lignes.Next() {
		var proID int
		var site, materiau string
		if err := lignes.Scan(&proID, &site, &materiau); err != nil {
			log.Printf("alertes de collecte du dépôt %d : %v", d.ID, err)
			return
		}
		message := "« " + d.DescriptionObjet + " »"
		if materiau != "" {
			message += " (" + strings.ToLower(materiau) + ")"
		}
		if site != "" {
			message += " vient d'être déposé à " + site + "."
		}
		notifierAvecLien(proID, "Objet disponible", message, "depots")
	}
}

// ---------- Récupération par code-barre ----------

type ObjetARecuperer struct {
	ID               int      `json:"id"`
	DescriptionObjet string   `json:"description_objet"`
	Site             *string  `json:"site"`
	Materiau         *string  `json:"materiau"`
	PoidsKg          *float64 `json:"poids_kg"`
}

// Seuls les objets « déposés » (présents dans un conteneur) peuvent être trouvés par leur code-barre.
func (s *Serveur) gererRechercheCodeBarre(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleProfessionnel); !ok {
		return
	}
	code := strings.ToLower(strings.TrimSpace(r.PathValue("code")))
	var o ObjetARecuperer
	err := s.db.QueryRow(`
		SELECT d.id, d.description_objet, co.site, c.libelle, d.poids_kg
		FROM depots d
		LEFT JOIN conteneurs co ON co.id = d.conteneur_id
		LEFT JOIN categories_materiaux c ON c.id = d.categorie_id
		WHERE lower(d.code_barre) = $1 AND d.statut = 'deposee'`, code,
	).Scan(&o.ID, &o.DescriptionObjet, &o.Site, &o.Materiau, &o.PoidsKg)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "aucun objet déposé ne correspond à ce code-barre")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusOK, o)
}

// ---------- Paiements de la personne connectée (page Contrats & Facturation) ----------

type MonPaiement struct {
	ID           int       `json:"id"`
	Source       string    `json:"source"`
	Libelle      *string   `json:"libelle"`
	Montant      float64   `json:"montant"`
	Statut       string    `json:"statut"`
	DatePaiement time.Time `json:"date_paiement"`
	FactureID    *int      `json:"facture_id"`
}

func (s *Serveur) gererMesPaiements(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}
	lignes, err := s.db.Query(`
		SELECT p.id, p.source, c.libelle, p.montant, p.statut, p.date_paiement,
		       (SELECT d.id FROM documents d WHERE d.type_document = 'facture' AND d.commande_id = c.id LIMIT 1)
		FROM paiements p LEFT JOIN commandes c ON c.paiement_id = p.id
		WHERE p.utilisateur_id = $1
		ORDER BY p.date_paiement DESC, p.id DESC
		LIMIT 50`, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()
	paiements := make([]MonPaiement, 0)
	for lignes.Next() {
		var p MonPaiement
		if err := lignes.Scan(&p.ID, &p.Source, &p.Libelle, &p.Montant, &p.Statut, &p.DatePaiement, &p.FactureID); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		paiements = append(paiements, p)
	}
	envoyerJSON(w, http.StatusOK, paiements)
}

// ---------- Contrat PDF (abonnement, campagne) ----------

func (s *Serveur) genererContratCommande(commandeID int) (int, error) {
	var existant int
	err := s.db.QueryRow("SELECT id FROM documents WHERE type_document = 'contrat' AND commande_id = $1 LIMIT 1", commandeID).Scan(&existant)
	if err == nil {
		return existant, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	var c Commande
	if err := scannerCommande(s.db.QueryRow(`SELECT `+colonnesCommande+` FROM commandes WHERE id = $1`, commandeID), &c); err != nil {
		return 0, err
	}
	var prenom, nom, email string
	var siret, adresse sql.NullString
	if err := s.db.QueryRow("SELECT prenom, nom, email, siret, adresse_entreprise FROM utilisateurs WHERE id = $1", c.UtilisateurID).
		Scan(&prenom, &nom, &email, &siret, &adresse); err != nil {
		return 0, err
	}
	date := time.Now()
	if c.DatePaiement != nil {
		date = *c.DatePaiement
	}
	parties := "Entre UpcycleConnect et " + prenom + " " + nom + " (" + email + ")"
	if siret.Valid {
		parties += ", SIRET " + siret.String
	}
	if adresse.Valid {
		parties += ", " + adresse.String
	}
	lignes := []string{parties + ".", "Date de prise d'effet : " + date.Format("02/01/2006") + "."}

	var titre string
	if c.Type == "abonnement" {
		titre = "Contrat d'abonnement Premium"
		lignes = append(lignes,
			"Objet : abonnement Premium donnant accès aux statistiques d'impact, à l'état des stocks en temps réel et aux alertes de collecte.",
			fmt.Sprintf("Prix : %.2f € par mois.", c.Montant),
			"Durée : mensuelle, sans engagement. Résiliable à tout moment depuis Mon compte, avec effet immédiat.",
		)
	} else {
		var cibleType string
		var debut time.Time
		if err := s.db.QueryRow("SELECT cible_type, date_debut FROM campagnes_publicitaires WHERE id = $1", c.ReferenceID).
			Scan(&cibleType, &debut); err != nil {
			return 0, err
		}
		element := map[string]string{"projet": "le projet", "annonce": "l'annonce"}[cibleType]
		titre = "Contrat de campagne publicitaire"
		lignes = append(lignes,
			"Objet : mise en avant de "+element+" « "+c.Libelle+" » dans la rubrique « À la une » de la plateforme.",
			fmt.Sprintf("Budget : %.2f € par mois (barème de 100 à 500 € par campagne).", c.Montant),
			"Début de la campagne : "+debut.Format("02/01/2006")+". Arrêt possible à tout moment depuis Contrats & Facturation.",
		)
	}
	contenu, err := genererPDF(titre, lignes)
	if err != nil {
		return 0, err
	}
	titreDocument := titre + " : " + c.Libelle
	if c.Type == "abonnement" {
		titreDocument = titre + " du " + date.Format("02/01/2006")
	}
	var id int
	err = s.db.QueryRow(`
		INSERT INTO documents (type_document, titre, utilisateur_id, commande_id, cree_par, contenu_pdf)
		VALUES ('contrat', $1, $2, $3, $2, $4) RETURNING id`,
		titreDocument, c.UtilisateurID, c.ID, contenu).Scan(&id)
	return id, err
}

// ---------- Facteurs d'impact des matériaux (admin) ----------

type ImpactEntree struct {
	Facteur *float64 `json:"co2_evite_par_kg"`
	Source  string   `json:"source_facteur"`
}

func (s *Serveur) gererModificationImpactMateriau(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}
	var entree ImpactEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if entree.Facteur == nil || *entree.Facteur < 0 || *entree.Facteur > 100 {
		envoyerErreur(w, http.StatusBadRequest, "facteur entre 0 et 100 kg CO2e par kg")
		return
	}
	if strings.TrimSpace(entree.Source) == "" {
		envoyerErreur(w, http.StatusBadRequest, "la source du facteur est obligatoire")
		return
	}
	resultat, err := s.db.Exec(`UPDATE categories_materiaux SET co2_evite_par_kg = $1, source_facteur = $2 WHERE id = $3`,
		*entree.Facteur, strings.TrimSpace(entree.Source), id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if n, _ := resultat.RowsAffected(); n == 0 {
		envoyerErreur(w, http.StatusNotFound, "matériau introuvable")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
