package workout

import (
	"log"

	"github.com/robfig/cron/v3"
)

func StartWorkoutScheduler(service WorkoutService) {
	c := cron.New()

	_, err := c.AddFunc("* * * * *", func() {
		log.Println("Running workout scheduler...")

		if err := service.SendTodaysWorkoutEmail(); err != nil {
			log.Printf("Workout email failed: %v", err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	c.Start()
}