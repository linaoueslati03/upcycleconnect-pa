package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func envoyerJSON(w http.ResponseWriter, code int, donnees any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(donnees)
}

func envoyerErreur(w http.ResponseWriter, code int, message string) {
	envoyerJSON(w, code, map[string]string{"erreur": message})
}

func idDepuisChemin(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}
