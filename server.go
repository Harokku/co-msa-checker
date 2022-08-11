package main

import (
	"co-msa-checker/database"
	"co-msa-checker/router"
	"co-msa-checker/utils"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func main() {
	// -------------------------
	// Variable definition
	// -------------------------

	var (
		err  error
		port string //server port from env
		app  *fiber.App
	)

	log.Printf("Starting environment init...")
	initStartTime := time.Now() //Startup timer start

	// -------------------------
	// .env loading
	// -------------------------

	// Load YAML with godotenv pkg
	err = godotenv.Load("env.yaml")

	if err != nil {
		log.Printf("WARNING:\tError loading .env file: %v\tUsing only already set env variables\n", err)
	}

	// Read server port from env
	port, err = utils.ReadEnv("PORT")
	if err != nil {
		log.Fatalf("FATAL:\terror setting server port: %v", err)
	}
	log.Printf("Server port set to: %v", port)

	// -------------------------
	// Database connection and init
	// -------------------------

	database.Connect()
	defer func(DbConnection *sql.DB) {
		err := DbConnection.Close()
		if err != nil {
			log.Printf("WARNING:\tCannot close db pool connection, check for memory leak\terr:%v", err)
		}
	}(database.DbConnection)

	// -------------------------
	// Fiber definition and server start
	// -------------------------

	app = fiberApp()

	// -------------------------
	// Init completed, starting listener
	// -------------------------

	initDuration := time.Since(initStartTime) //calculate total startup time
	log.Printf("Enviromnent initialized in %s", initDuration)

	err = app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("FATAL:\terror starting server: %s", err)
	}
}

func fiberApp() *fiber.App {
	var (
		app *fiber.App
	)
	app = fiber.New()
	app.Use(logger.New())  //logger init
	app.Use(cors.New())    //CORS init
	app.Use(recover.New()) //recover init

	// -------------------------
	// Static routes
	// -------------------------

	app.Static("/", "./static")

	// -------------------------
	// Debug routes
	// -------------------------

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	// -------------------------
	// Router init (config in router pkg)
	// -------------------------
	router.Setup(app)

	return app
}
