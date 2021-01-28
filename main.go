package main

import (
	"bytes"
	"context"
	"fmt"
	"kitalulus/handler"
	"kitalulus/helper"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// LogResponseWriter sturct
type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

// LogMiddleware sturct
type LogMiddleware struct {
	logger *log.Logger
}

// NewLogResponseWriter - return logresponsewriter with response writer inside
func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

// CustomLogMiddleware - return logmiddleware
func CustomLogMiddleware(logger *log.Logger) *LogMiddleware {
	return &LogMiddleware{logger: logger}
}

// Init - init middleware
func (m *LogMiddleware) Init() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			logRespWriter := NewLogResponseWriter(w)
			next.ServeHTTP(logRespWriter, r)

			m.logger.Printf(
				"duration=%s url=%s method=%s",
				time.Since(startTime).String(),
				r.URL.String(),
				r.Method)
		})
	}
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	fmt.Println("Hello World")
	logger := log.New(os.Stdout, "", log.LstdFlags)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://" + goDotEnvVariable("MONGO_DB_USERNAME") + ":" + goDotEnvVariable("MONGO_DB_PASSWORD") + "@" + goDotEnvVariable("MONGO_DB_HOST") + "/?authSource=admin&readPreference=primary&ssl=false")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err.Error())
	}

	router := mux.NewRouter()
	var h = handler.Handler{R: helper.Response{}, C: client}
	router.HandleFunc("/question", h.QuestionList()).Methods("GET")
	router.HandleFunc("/question/{uuid}", h.QuestionShow()).Methods("GET")
	router.HandleFunc("/question", h.QuestionCreate()).Methods("POST")
	router.HandleFunc("/question/remove", h.QuestionRemove()).Methods("POST")
	router.HandleFunc("/question/update/{id}", h.QuestionUpdate()).Methods("POST")

	// logging
	logMiddleware := CustomLogMiddleware(logger)
	router.Use(logMiddleware.Init())

	logger.Fatalln(http.ListenAndServe(":12345", router))
}
