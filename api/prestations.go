package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Prestation struct {
	ID        int     `json:"id"`
	Titre     string  `json:"titre"`
	Categorie string  `json:"categorie"`
	Tarif     float64 `json:"tarif"`
	Statut    string  `json:"statut"`
}

func gererListePrestations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, titre, categorie, tarif, statut FROM prestations")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var prestations []Prestation
	for rows.Next() {
		var p Prestation
		if err := rows.Scan(&p.ID, &p.Titre, &p.Categorie, &p.Tarif, &p.Statut); err != nil {
			log.Println("Erreur de scan :", err)
			continue
		}
		prestations = append(prestations, p)
	}

	json.NewEncoder(w).Encode(prestations)
}

func gererCreationPrestation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p Prestation
	json.NewDecoder(r.Body).Decode(&p)

	err := db.QueryRow("INSERT INTO prestations (titre, categorie, tarif, statut) VALUES ($1, $2, $3, $4) RETURNING id", p.Titre, p.Categorie, p.Tarif, p.Statut).Scan(&p.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func gererModificationPrestation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	var p Prestation
	json.NewDecoder(r.Body).Decode(&p)

	_, err := db.Exec("UPDATE prestations SET titre=$1, categorie=$2, tarif=$3, statut=$4 WHERE id=$5", p.Titre, p.Categorie, p.Tarif, p.Statut, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"statut": "modifié"})
}

func gererSuppressionPrestation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := db.Exec("DELETE FROM prestations WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
