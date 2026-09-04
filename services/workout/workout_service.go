package workout

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/jackc/pgx/v5"
)

type WorkoutService interface {
	SendTodaysWorkoutEmail() error
}

type workoutService struct {
	repo         *Repository
	emailService EmailService
	now          func() time.Time
}

func NewWorkoutService(repo *Repository, emailService EmailService) WorkoutService {
	return &workoutService{
		repo:         repo,
		emailService: emailService,
		now:          time.Now,
	}
}

func (s *workoutService) SendTodaysWorkoutEmail() error {
	now := s.now()
	dayOfWeek := workoutDayOfWeek(now)

	userIDs, err := s.repo.GetWorkoutEmailUsers(context.Background())
	if err != nil {
		return err
	}

	for _, userID := range userIDs {

		settings, err := s.repo.GetWorkoutEmailSettings(
			context.Background(),
			userID,
		)

		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}

		if err != nil {
			return err
		}

		if settings.Email == "" || settings.EmailTime == "" {
			continue
		}

		if settings.EmailTime != now.Format("15:04") {
			continue
		}

		schedule, err := s.repo.GetWorkoutScheduleByDay(
			context.Background(),
			userID,
			dayOfWeek,
		)

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		emailData := WorkoutScheduleEmail{
			DayName:  workoutDayName(dayOfWeek),
			RestDay:  errors.Is(err, pgx.ErrNoRows),
			Schedule: schedule,
		}

		subject := "Today's workout"

		if emailData.RestDay {
			subject = "Today's workout: Rest day"
		} else {
			subject = fmt.Sprintf(
				"Today's workout: %s",
				schedule.WorkoutName,
			)
		}

		htmlBody, err := renderWorkoutScheduleEmail(emailData)
		if err != nil {
			return err
		}

		err = s.emailService.SendWorkoutEmail(
			settings.Email,
			subject,
			htmlBody,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func workoutDayOfWeek(t time.Time) int {
	if t.Weekday() == time.Sunday {
		return 7
	}

	return int(t.Weekday())
}

func workoutDayName(dayOfWeek int) string {
	names := map[int]string{
		1: "Monday",
		2: "Tuesday",
		3: "Wednesday",
		4: "Thursday",
		5: "Friday",
		6: "Saturday",
		7: "Sunday",
	}

	return names[dayOfWeek]
}

type WorkoutScheduleEmail struct {
	DayName  string
	RestDay  bool
	Schedule WorkoutScheduleResponse
}

func renderWorkoutScheduleEmail(data WorkoutScheduleEmail) (string, error) {
	tmpl, err := template.ParseFiles(
		"assets/templates/emails/workout_schedule.html",
	)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}
