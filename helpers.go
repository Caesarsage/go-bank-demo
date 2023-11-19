package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// Helpers

type ApiError struct {
	Error string `json:"error"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id %s", idStr)
	}

	return id, nil
}

func createJWTToken(account *Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": account.ID,
		"email":					account.Email,
		"exp":            time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

// Middleware
func withJWTAuth(handleFunc http.HandlerFunc, s Storage) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)

		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid token"})
			return
		}

		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Permission denied"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		fmt.Println(claims)

		userID, err := getID(r)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Permission denied"})
			return
		}

		account, err := s.GetAccountById(userID)

		if int64(account.ID) != int64(claims["id"].(float64)) {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Permission denied"})
			return
		}

		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid token"})
			return
		}
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

// func generateAccountNumber() int {
// 	b := make([]byte, 13) // adjust size for desired length
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	number := fmt.Sprintf("%x", b)
// 	return number
// }

