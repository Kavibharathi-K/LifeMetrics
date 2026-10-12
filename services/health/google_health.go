package health

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const googleHealthAPIBaseURL = "https://health.googleapis.com/v4"

// DailyActivityResponse represents the daily rollup response from Google Health.
type DailyActivityResponse struct {
	RollupDataPoints []struct {
		CivilStartTime struct {
			Date struct {
				Year  int `json:"year"`
				Month int `json:"month"`
				Day   int `json:"day"`
			} `json:"date"`
		} `json:"civilStartTime"`

		Steps struct {
			CountSum string `json:"countSum"`
		} `json:"steps"`

		Distance struct {
			MillimetersSum string `json:"millimetersSum"`
		} `json:"distance"`

		ActiveEnergyBurned struct {
			KcalSum *float64 `json:"kcalSum"`
		} `json:"activeEnergyBurned"`

		TotalCalories struct {
			KcalSum *float64 `json:"kcalSum"`
		} `json:"totalCalories"`

		ActiveMinutes struct {
			ActiveMinutesRollupByActivityLevel []struct {
				ActivityLevel    string `json:"activityLevel"`
				ActiveMinutesSum string `json:"activeMinutesSum"`
			} `json:"activeMinutesRollupByActivityLevel"`
		} `json:"activeMinutes"`
	} `json:"rollupDataPoints"`
}

// FetchDailyRollup retrieves a daily rollup for one Google Health data type.
func (h *Handler) FetchDailyRollup(
	ctx context.Context,
	accessToken string,
	dataType string,
	today time.Time,
) ([]byte, error) {
	tomorrow := today.AddDate(0, 0, 1)

	requestBody := fmt.Sprintf(`{
		"range": {
			"start": {
				"date": {
					"year": %d,
					"month": %d,
					"day": %d
				},
				"time": {
					"hours": 0,
					"minutes": 0,
					"seconds": 0,
					"nanos": 0
				}
			},
			"end": {
				"date": {
					"year": %d,
					"month": %d,
					"day": %d
				},
				"time": {
					"hours": 0,
					"minutes": 0,
					"seconds": 0,
					"nanos": 0
				}
			}
		},
		"windowSizeDays": 1
	}`,
		today.Year(),
		int(today.Month()),
		today.Day(),
		tomorrow.Year(),
		int(tomorrow.Month()),
		tomorrow.Day(),
	)

	apiURL := fmt.Sprintf(
		"%s/users/me/dataTypes/%s/dataPoints:dailyRollUp",
		googleHealthAPIBaseURL,
		dataType,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiURL,
		bytes.NewBufferString(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create Google Health request for %s: %w",
			dataType,
			err,
		)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call Google Health API for %s: %w",
			dataType,
			err,
		)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read Google Health response for %s: %w",
			dataType,
			err,
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"Google Health API for %s returned status %d: %s",
			dataType,
			resp.StatusCode,
			string(body),
		)
	}

	return body, nil
}

// SyncDailyActivityData synchronizes today's activity metrics from Google Health.
func (h *Handler) SyncDailyActivityData(
	ctx context.Context,
	userID int64,
) error {
	// 1. Obtain a valid Google Health access token.
	accessToken, err := h.GetValidAccessToken(ctx, userID)
	if err != nil {
		return fmt.Errorf(
			"failed to get Google Health access token: %w",
			err,
		)
	}

	// 2. Calculate today's date in the Asia/Kolkata timezone.
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return fmt.Errorf("failed to load timezone: %w", err)
	}

	today := time.Now().In(loc)

	// 3. Initialize collected metrics.
	// nil means the metric was not available in the response.
	var (
		steps             *int64
		distanceMeters    *float64
		activeEnergyKcal  *float64
		totalCaloriesKcal *float64
		activeMinutes     *float64
		activityDate      time.Time
	)

	// 4. Fetch all five daily activity metrics.
	dataTypes := []string{
		"steps",
		"distance",
		"active-energy-burned",
		"total-calories",
		"active-minutes",
	}

	for _, dataType := range dataTypes {
		body, err := h.FetchDailyRollup(
			ctx,
			accessToken,
			dataType,
			today,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to fetch %s: %w",
				dataType,
				err,
			)
		}

		var result DailyActivityResponse

		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf(
				"failed to parse %s response: %w",
				dataType,
				err,
			)
		}

		// A missing rollup means no data was returned for this metric.
		if len(result.RollupDataPoints) == 0 {
			continue
		}

		point := result.RollupDataPoints[0]

		// Use the civil date returned by Google Health.
		if activityDate.IsZero() {
			date := point.CivilStartTime.Date

			activityDate = time.Date(
				date.Year,
				time.Month(date.Month),
				date.Day,
				0, 0, 0, 0,
				time.UTC,
			)
		}

		switch dataType {
		case "steps":
			if point.Steps.CountSum == "" {
				continue
			}

			value, err := strconv.ParseInt(
				point.Steps.CountSum,
				10,
				64,
			)
			if err != nil {
				return fmt.Errorf(
					"failed to parse steps count: %w",
					err,
				)
			}

			steps = &value

		case "distance":
			if point.Distance.MillimetersSum == "" {
				continue
			}

			millimeters, err := strconv.ParseFloat(
				point.Distance.MillimetersSum,
				64,
			)
			if err != nil {
				return fmt.Errorf(
					"failed to parse distance: %w",
					err,
				)
			}

			// Convert millimeters to meters.
			meters := millimeters / 1000
			distanceMeters = &meters

		case "active-energy-burned":
			if point.ActiveEnergyBurned.KcalSum == nil {
				continue
			}

			value := *point.ActiveEnergyBurned.KcalSum
			activeEnergyKcal = &value

		case "total-calories":
			if point.TotalCalories.KcalSum == nil {
				continue
			}

			value := *point.TotalCalories.KcalSum
			totalCaloriesKcal = &value

		case "active-minutes":
			var totalMinutes int64
			hasMinutes := false

			// Sum the minutes across all returned activity levels.
			for _, activityLevel := range point.ActiveMinutes.ActiveMinutesRollupByActivityLevel {

				if activityLevel.ActiveMinutesSum == "" {
					continue
				}

				minutes, err := strconv.ParseInt(
					activityLevel.ActiveMinutesSum,
					10,
					64,
				)
				if err != nil {
					return fmt.Errorf(
						"failed to parse active minutes for %s: %w",
						activityLevel.ActivityLevel,
						err,
					)
				}

				totalMinutes += minutes
				hasMinutes = true
			}

			if hasMinutes {
				value := float64(totalMinutes)
				activeMinutes = &value
			}
		}
	}

	// 5. Ensure at least one metric was returned.
	if activityDate.IsZero() {
		return fmt.Errorf(
			"Google Health returned no daily activity data for %s",
			today.Format("2006-01-02"),
		)
	}

	// 6. Save all collected metrics with one database call.
	err = h.Repo.UpsertDailyActivity(
		ctx,
		userID,
		activityDate,
		steps,
		distanceMeters,
		activeEnergyKcal,
		totalCaloriesKcal,
		activeMinutes,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to save daily activity to database: %w",
			err,
		)
	}

	return nil
}