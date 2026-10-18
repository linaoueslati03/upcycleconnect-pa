package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

type reponseSurveillee struct {
	http.ResponseWriter
	remplacee bool
}

func (rep *reponseSurveillee) WriteHeader(code int) {
	dejaJSON := rep.Header().Get("Content-Type") == "application/json"
	if !dejaJSON && (code == http.StatusNotFound || code == http.StatusMethodNotAllowed) {
		rep.remplacee = true
		message := "route introuvable"
		if code == http.StatusMethodNotAllowed {
			message = "méthode non autorisée pour cette route"
		}
		ecrireErreur(rep.ResponseWriter, code, message)
		return
	}
	rep.ResponseWriter.WriteHeader(code)
}

func (rep *reponseSurveillee) Write(contenu []byte) (int, error) {
	if rep.remplacee {
		return len(contenu), nil
	}
	return rep.ResponseWriter.Write(contenu)
}

func ecrireErreur(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"erreur": message})
}

func Erreurs(suivant http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if erreur := recover(); erreur != nil {
				log.Println("erreur interne :", erreur)
				ecrireErreur(w, http.StatusInternalServerError, "erreur serveur")
			}
		}()
		suivant.ServeHTTP(&reponseSurveillee{ResponseWriter: w}, r)
	})
}
