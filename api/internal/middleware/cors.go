// Package middleware contient les traitements appliqués à toutes les requêtes.
package middleware

import "net/http"

// CORS autorise le front (servi sur un autre port) à appeler l'API depuis le navigateur.
func CORS(suivant http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		suivant.ServeHTTP(w, r)
	})
}
