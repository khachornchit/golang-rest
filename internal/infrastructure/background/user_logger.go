package background

import (
	"context"
	"golang-rest/internal/core/ports"
	"log"
	"sync"
	"time"
)

func StartUserLogger(ctx context.Context, wg *sync.WaitGroup, userRepository ports.UserRepositoryInterface) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("User logger shutting down.")
				return
			case <-ticker.C:
				users, err := userRepository.GetAllUsers()
				if err != nil {
					log.Printf("Error fetching users: %v\n", err)
					continue
				}
				log.Printf("Total users: %d\n", len(users))
			}
		}
	}()
}
