package workers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type JobResult struct {
	ID      string                 `json:"id"`
	Status  string                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

const (
	StatusQueued = "queued"
	StatusDone   = "done"
	StatusError  = "error"
)

type Worker struct {
	queue   chan map[string]interface{}
	results map[string]*JobResult
}

var workerInstance *Worker

func InitializePaymentAcledaTaskWorker() {
	workerInstance = &Worker{
		queue:   make(chan map[string]interface{}, 100),
		results: make(map[string]*JobResult),
	}

	// Start the worker goroutine
	go workerInstance.processQueue()
}

func (w *Worker) processQueue() {
	for payload := range w.queue {
		ctx := context.Background()
		_ = ctx

		jobID, _ := payload["job_id"].(string)
		if jobID == "" {
			jobID = generateJobID()
		}

		// Placeholder: implement Acleda payment processing here
		w.results[jobID] = &JobResult{
			ID:      jobID,
			Status:  StatusDone,
			Message: "Acleda payment processed successfully",
			Data:    payload,
		}
		log.Printf("Job %s completed", jobID)
	}
}

func PaymentHandler(c *fiber.Ctx) error {
	// Placeholder for Acleda payment processing
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Acleda payment handler - not implemented yet",
	})
}

func EnqueueHandler(c *fiber.Ctx) error {
	if workerInstance == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "worker not initialized"})
	}

	// Parse payload from request
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}

	jobID := generateJobID()
	payload["job_id"] = jobID

	workerInstance.results[jobID] = &JobResult{ID: jobID, Status: StatusQueued, Message: "Job queued"}
	select {
	case workerInstance.queue <- payload:
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"job_id": jobID, "status": StatusQueued})
	default:
		workerInstance.results[jobID] = &JobResult{ID: jobID, Status: StatusError, Error: "queue full"}
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "queue full"})
	}
}

func StatusHandler(c *fiber.Ctx) error {
	if workerInstance == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "worker not initialized"})
	}
	jobID := c.Query("id")
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job ID is required"})
	}
	result, ok := workerInstance.results[jobID]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Job not found"})
	}
	return c.JSON(result)
}

func generateJobID() string {
	return "job-" + time.Now().Format("20060102-150405-999999999")
}
