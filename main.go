package main

import (
	"log"
	"strconv"

	"payment-airpay/infrastructure/configuration"
	"payment-airpay/infrastructure/controllers"
	"payment-airpay/infrastructure/database"
	"payment-airpay/infrastructure/database/clients"
	"payment-airpay/infrastructure/database/repositories"
	"payment-airpay/infrastructure/dependencies"
	"payment-airpay/infrastructure/publishers"
	"payment-airpay/infrastructure/queue"
	"payment-airpay/infrastructure/workers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	defer func() {
		queue.CloseRabbitMQ()
	}()

	log.Println("Acleda Worker is starting...")

	// Initialize configurations
	log.Println("Initializing configuration...")
	configuration.InitializeAppConfig()
	log.Println("Configuration initialized")

	log.Println("Initializing YugabyteDB...")
	database.InitializeYugabyteDB()
	log.Println("YugabyteDB initialized")

	// Initialize RabbitMQ
	log.Println("Initializing RabbitMQ...")
	queue.InitializeRabbitMQ()
	log.Println("RabbitMQ initialized")

	// Initialize Redis for publishing
	log.Println("Initializing publishers...")
	publishers.InitializeRedis()
	log.Println("Publishers initialized")

	// Initialize worker
	log.Println("Initializing worker...")
	workers.InitializePaymentAcledaTaskWorker()
	log.Println("Worker initialized")

	// Initialize fiber app with HTML template engine
	engine := html.New("./infrastructure/views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Initialize Yugabyte client wrapper
	yugabyteClient := clients.NewYugabyteClient(database.YugabyteDBClient)

	// Initialize Acleda controller
	acledaController := controllers.NewAcledaController(
		dependencies.ProvideAcledaGateway(),
		dependencies.ProvidePaymentAcledaService(),
		repositories.NewPaymentAcledaRepositoryYugabyteDB(yugabyteClient),
	)

	// Register routes
	app.Post("/payment/acleda", workers.PaymentHandler)
	app.Post("/payment/acleda/async", workers.EnqueueHandler)
	app.Get("/jobs/status", workers.StatusHandler)

	// Setup Acleda controller routes
	app.Post("/api/v1/acleda/payment-links", acledaController.CreatePaymentLink)
	app.Get("/payment-page/acleda/:id", acledaController.PaymentPage)
	app.Get("/api/v1/acleda/payments/:id/status", acledaController.GetPaymentStatus)

	// Start server
	port := strconv.Itoa(configuration.AppConfig.ApplicationPort)
	if port == "0" || port == "" {
		port = "8080" // default port
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
