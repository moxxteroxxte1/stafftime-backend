package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"math/rand"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/moxxteroxxte1/stafftime-backend/src/models"
)

type APIServer struct {
	listenAddr string
	database   *gorm.DB
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{listenAddr: listenAddr}
}

func (s *APIServer) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"message": "hello world"})
	})

	router.HandleFunc("/auth", s.HandleLogin)
	router.HandleFunc("/uploads/{path:(?:[0-9]+.(?:jpg|jpeg|png|webp))}", func(w http.ResponseWriter, r *http.Request) {
		buff, err := os.ReadFile(fmt.Sprintf("./uploads/%s", mux.Vars(r)["path"]))
		if err != nil {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(buff)
	}).Methods(http.MethodGet)

	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(s.jwtAuthMiddleware)
	protectedRouter.Use(s.rbacMiddleware)

	protectedRouter.HandleFunc("/users", s.handleUsers)
	protectedRouter.HandleFunc("/users/{userID}", s.HandleUserByID)
	protectedRouter.HandleFunc("/users/{userID}/upload", s.HandleUserUpload)
	protectedRouter.HandleFunc("/users/{userID}/image", s.HandleUserUpload)

	protectedRouter.HandleFunc("/users/{userID}/shifts", s.HandleShiftByUserID)
	protectedRouter.HandleFunc("/users/{userID}/shifts/{shiftID}", s.HandleShiftByID)

	protectedRouter.HandleFunc("/users/{userID}/payments", s.HandlePaymentByUserID)
	protectedRouter.HandleFunc("/users/{userID}/payments/{paymentID}", s.HandlePaymentByID)

	protectedRouter.HandleFunc("/users/{userID}/contracts", s.HandleContractByUserID)
	protectedRouter.HandleFunc("/users/{userID}/contracts/{contractID}", s.HandleContractByID)

	protectedRouter.HandleFunc("/shifts", s.handleShifts)
	protectedRouter.HandleFunc("/shifts/{shiftID}", s.HandleShiftByID)

	protectedRouter.HandleFunc("/payments", s.handlePayments)
	protectedRouter.HandleFunc("/payments/{paymentID}", s.HandlePaymentByID)

	protectedRouter.HandleFunc("/contracts", s.handleContracts)
	protectedRouter.HandleFunc("/contracts/{contractID}", s.HandleContractByID)

	protectedRouter.HandleFunc("/statuses", s.handleStatus)
	protectedRouter.HandleFunc("/statuses/{statusID}", s.HandleStatusByID)

	protectedRouter.HandleFunc("/keys", s.handleKeys)
	protectedRouter.HandleFunc("/keys/{keyID}", s.HandleKeysByID)

	log.Printf("connecting to database")
	s.database = Connect()

	log.Printf("starting server on %s", s.listenAddr)
	err := http.ListenAndServeTLS(s.listenAddr, "server.crt", "server.key", router)
	if err != nil {
		log.Println(err)
	}
}

func Connect() *gorm.DB {
	dsn := fmt.Sprintf("host=stafftime-backend-database user=%s password=%s dbname=%s port=5432 sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	fmt.Println("connected!")

	fmt.Println("migrating models")
	db.AutoMigrate(&models.Contract{}, &models.Payment{}, &models.Shift{}, &models.Status{}, &models.User{}, &models.Key{})

	var user = new(models.User)
	if userErr := db.First(&user).Error; userErr != nil {
		bytes, bcryptErr := bcrypt.GenerateFromPassword([]byte("password"), 14)
		if bcryptErr != nil {
			log.Fatal(bcryptErr)
		}

		username := fmt.Sprintf("admin%d", rand.Intn(100))
		log.Println(fmt.Sprintf(`Username: %s`, username))
		user := models.User{
			Username:  username,
			Firstname: "admin",
			Lastname:  "admin",
			Email:     "email",
			Password:  string(bytes),
		}

		db.Create(&user)
	}

	return db
}
