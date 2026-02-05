package background

import (
	"freepass-2026/internal/service"
	"log"
	"time"
)

type CleanupJob struct {
	orderService service.IOrderService
	interval     time.Duration
	stopChan     chan bool
}

func NewCleanupJob(orderService service.IOrderService) *CleanupJob {
	return &CleanupJob{
		orderService: orderService,
		interval:     3 * 24 * time.Hour, // 3 days
		stopChan:     make(chan bool),
	}
}

func (c *CleanupJob) Start() {
	go func() {
		ticker := time.NewTicker(c.interval) // set ticker
		defer ticker.Stop()                  // stop ticker when done

		log.Println("Clean up job started - running every 3 days")

		for {
			select {
			case <-ticker.C:
				log.Println("Running clean up job...")
				err := c.orderService.CleanupCanceledOrders()
				if err != nil {
					log.Printf("Error during clean up: %v", err)
				} else {
					log.Println("Clean up job completed successfully")
				}
			case <-c.stopChan:
				log.Println("Clean up job stopped")
				return
			}
		}
	}()

}

func (c *CleanupJob) Stop() {
	c.stopChan <- true
}
