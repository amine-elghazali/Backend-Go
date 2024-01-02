package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {

		idConv, err := getId(r)
		if err != nil {
			return err
		}
		account, err := s.store.GetAccountByID(idConv)
		if err != nil {
			return err
		}

		return WriteJson(w, http.StatusOK, account)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("following method not allowed : %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	createAccountReq := new(CreateAccountRequest) // Allocates memory for the struct and returns a pointer.
	// createAccountReq := CreateAccountRequest{}	//Creates an actual instance of the struct (not a pointer). so if u want to use in the Decode methode , u need to pass a reference & to it

	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	token, err := createJWT(account)
	fmt.Println("token : ", token)

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	idConv, err := getId(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(idConv); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]int{"id deleted ": idConv})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	transferRequest := new(TransferRequet)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, transferRequest)
}

func getId(r *http.Request) (int, error) {
	id := mux.Vars(r)["id"]
	// We can then do something like : DB.getId(id) ...
	// fmt.Print(id)
	idConv, err := strconv.Atoi(id)

	if err != nil {
		// log.Fatal("Conversion error !")
		return idConv, fmt.Errorf("conversion error ! make sure to enter a vali ID")
	}

	return idConv, nil
}

func WriteJson(w http.ResponseWriter, status int, v interface{}) error {
	// Start Writing JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func permissionDenied(w http.ResponseWriter) {
	WriteJson(w, http.StatusForbidden, ApiError{Error: "permission denied"})

}

// Creating a Decoraator for JWT middleware Check
func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT Auth middleware check")

		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)

		if err != nil || !token.Valid {
			permissionDenied(w)
			return
		}

		userId, err := getId(r)

		if err != nil {
			permissionDenied(w)
			return
		}

		account, err := s.GetAccountByID(userId)

		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}

}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func createJWT(account *Account) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	})
	// It's better to use your own Claim Structure here , bcw this JWT Standard Claim is just a map,
	// using ur own struct Much better so you can use pre defined types ...

	secret := os.Getenv("JWT_SECRET")

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(secret))
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// fmt.Print("Hello")
			// Handle the error :
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHttpHandleFunc(s.handleGetAccountById), s.store))
	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer))

	log.Println("SERVER RUNNING ON PORT : ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed : %s", r.Method)
}
