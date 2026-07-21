package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"familiz/internal/utils"
)

// Authenticate vérifie le token JWT et l'ajoute au contexte
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Récupérer le header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token manquant", http.StatusUnauthorized)
			return
		}

		// 2. Vérifier le format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Format de token invalide. Utilisez 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// 3. Valider le token
		token, err := ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Token invalide ou expiré", http.StatusUnauthorized)
			return
		}

		// 4. Extraire les claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Token invalide", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Token invalide: user_id manquant", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Token invalide: email manquant", http.StatusUnauthorized)
			return
		}

		// 5. Ajouter les infos dans le contexte de la requête
		ctx := context.WithValue(r.Context(), utils.UserIDKey, int(userIDFloat))
		ctx = context.WithValue(ctx, utils.UserEmailKey, email)
		ctx = context.WithValue(ctx, utils.UserRoleKey, claims["role"])

		// 6. Passer à la suite
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
