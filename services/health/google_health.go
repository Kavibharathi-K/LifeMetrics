package health

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

type DailyRestingHeartRateResponse struct {
	DataPoints    []DailyRestingHeartRateDataPoint `json:"dataPoints"`
	NextPageToken string                           `json:"nextPageToken"`
}

type DailyRestingHeartRateDataPoint struct {
	DailyRestingHeartRate DailyRestingHeartRateValue `json:"dailyRestingHeartRate"`
}

type DailyRestingHeartRateValue struct {
	Date struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	} `json:"date"`

	BeatsPerMinute string `json:"beatsPerMinute"`
}

type DailyHRVResponse struct {
	DataPoints    []DailyHRVDataPoint `json:"dataPoints"`
	NextPageToken string              `json:"nextPageToken"`
}

type DailyHRVDataPoint struct {
	DailyHeartRateVariability DailyHRVValue `json:"dailyHeartRateVariability"`
}

type DailyHRVValue struct {
	Date struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	} `json:"date"`

	AverageHRVMilliseconds *float64 `json:"averageHeartRateVariabilityMilliseconds"`

	DeepSleepRMSSDMilliseconds *float64 `json:"deepSleepRootMeanSquareOfSuccessiveDifferencesMilliseconds"`
}

type SleepResponse struct {
	DataPoints    []SleepDataPoint `json:"dataPoints"`
	NextPageToken string           `json:"nextPageToken"`
}

type SleepDataPoint struct {
	Name  string           `json:"name"`
	Sleep SleepSessionData `json:"sleep"`
}

type SleepSessionData struct {
	Interval SleepInterval `json:"interval"`
	Type     string        `json:"type"`
	Stages   []SleepStage  `json:"stages"`

	Metadata SleepMetadata `json:"metadata"`
	Summary  SleepSummary  `json:"summary"`

	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type SleepInterval struct {
	StartTime      string         `json:"startTime"`
	StartUTCOffset string         `json:"startUtcOffset"`
	EndTime        string         `json:"endTime"`
	EndUTCOffset   string         `json:"endUtcOffset"`
	CivilStartTime *CivilDateTime `json:"civilStartTime"`
	CivilEndTime   *CivilDateTime `json:"civilEndTime"`
}

type CivilDateTime struct {
	Date CivilDate  `json:"date"`
	Time *CivilTime `json:"time"`
}

type CivilDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type CivilTime struct {
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
	Seconds int `json:"seconds"`
	Nanos   int `json:"nanos"`
}

type SleepStage struct {
	StartTime      string `json:"startTime"`
	StartUTCOffset string `json:"startUtcOffset"`
	EndTime        string `json:"endTime"`
	EndUTCOffset   string `json:"endUtcOffset"`
	Type           string `json:"type"`
	CreateTime     string `json:"createTime"`
	UpdateTime     string `json:"updateTime"`
}

type SleepMetadata struct {
	Processed      bool   `json:"processed"`
	Nap            bool   `json:"nap"`
	MainSleep      bool   `json:"mainSleep"`
	StagesStatus   string `json:"stagesStatus"`
	ManuallyEdited bool   `json:"manuallyEdited"`
	ExternalID     string `json:"externalId"`
}

type SleepSummary struct {
	MinutesInSleepPeriod string              `json:"minutesInSleepPeriod"`
	MinutesAfterWakeUp   string              `json:"minutesAfterWakeUp"`
	MinutesToFallAsleep  string              `json:"minutesToFallAsleep"`
	MinutesAsleep        string              `json:"minutesAsleep"`
	MinutesAwake         string              `json:"minutesAwake"`
	StagesSummary        []SleepStageSummary `json:"stagesSummary"`
}

type SleepStageSummary struct {
	Type    string `json:"type"`
	Minutes string `json:"minutes"`
	Count   string `json:"count"`
}

// FetchDailyRollup retrieves a daily rollup for one Google Health data type.
func (h *Handler) FetchHealthData(ctx context.Context, accessToken string, dataType string, action string, today time.Time) ([]byte, error) {
	var (
		method      string
		apiURL      string
		requestBody []byte
		err         error
	)

	switch action {
	case "dailyRollUp":
		tomorrow := today.AddDate(0, 0, 1)

		body := fmt.Sprintf(`{
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

		requestBody = []byte(body)
		method = http.MethodPost
		apiURL = fmt.Sprintf(
			"%s/users/me/dataTypes/%s/dataPoints:dailyRollUp",
			googleHealthAPIBaseURL,
			dataType,
		)

	case "list":
		method = http.MethodGet
		apiURL = fmt.Sprintf(
			"%s/users/me/dataTypes/%s/dataPoints",
			googleHealthAPIBaseURL,
			dataType,
		)

	case "reconcile":
		method = http.MethodGet
		apiURL = fmt.Sprintf(
			"%s/users/me/dataTypes/%s/dataPoints:reconcile",
			googleHealthAPIBaseURL,
			dataType,
		)
	case "sleepYesterday":
		yesterday := today.AddDate(0, 0, -1)

		startOfYesterday := time.Date(
			yesterday.Year(),
			yesterday.Month(),
			yesterday.Day(),
			0, 0, 0, 0,
			yesterday.Location(),
		)

		startUTC := startOfYesterday.UTC()
		endUTC := startOfYesterday.AddDate(0, 0, 1).UTC()

		filter := fmt.Sprintf(
			`sleep.interval.end_time >= "%s" AND sleep.interval.end_time < "%s"`,
			startUTC.Format(time.RFC3339),
			endUTC.Format(time.RFC3339),
		)

		log.Printf("Google Health sleep filter: %s", filter)

		method = http.MethodGet
		apiURL = fmt.Sprintf(
			"%s/users/me/dataTypes/%s/dataPoints?filter=%s",
			googleHealthAPIBaseURL,
			dataType,
			url.QueryEscape(filter),
		)

	default:
		return nil, fmt.Errorf(
			"unsupported Google Health API action: %s",
			action,
		)
	}

	var requestReader io.Reader
	if requestBody != nil {
		requestReader = bytes.NewReader(requestBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, apiURL, requestReader)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create Google Health request for %s: %w",
			dataType,
			err,
		)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call Google Health API for %s using %s: %w",
			dataType,
			action,
			err,
		)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read Google Health response for %s: %w",
			dataType,
			err,
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"Google Health API for %s using %s returned status %d: %s",
			dataType,
			action,
			resp.StatusCode,
			string(responseBody),
		)
	}

	return responseBody, nil
}

// SyncDailyActivityData synchronizes today's activity metrics from Google Health.
func (h *Handler) SyncDailyActivityData(ctx context.Context, userID int64) error {
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
		body, err := h.FetchHealthData(
			ctx,
			accessToken,
			dataType,
			"dailyRollUp",
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

func (h *Handler) SyncDailyVitalsData(ctx context.Context, userID int64) error {
	accessToken, err := h.GetValidAccessToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("get valid Google Health access token: %w", err)
	}

	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return fmt.Errorf("load timezone: %w", err)
	}

	today := time.Now().In(location)

	var restingHeartRate *float64
	var hrvMS *float64

	// 1. Fetch daily resting heart rate.
	response, err := h.FetchHealthData(
		ctx,
		accessToken,
		"daily-resting-heart-rate",
		"list",
		today,
	)
	if err != nil {
		return fmt.Errorf("fetch daily resting heart rate: %w", err)
	}

	var restingHRResult DailyRestingHeartRateResponse
	if err := json.Unmarshal(response, &restingHRResult); err != nil {
		return fmt.Errorf("parse daily resting heart rate response: %w", err)
	}

	for _, point := range restingHRResult.DataPoints {
		value := point.DailyRestingHeartRate
		date := value.Date

		if date.Year == today.Year() &&
			date.Month == int(today.Month()) &&
			date.Day == today.Day() {

			bpm, err := strconv.ParseFloat(value.BeatsPerMinute, 64)
			if err != nil {
				return fmt.Errorf("parse resting heart rate: %w", err)
			}

			restingHeartRate = &bpm

			fmt.Printf("Today's resting heart rate: %.0f BPM\n", bpm)
			break
		}
	}

	if restingHeartRate == nil {
		fmt.Println("No resting heart rate data available for today")
	}

	// 2. Fetch daily HRV.
	hrvResponse, err := h.FetchHealthData(
		ctx,
		accessToken,
		"daily-heart-rate-variability",
		"list",
		today,
	)
	if err != nil {
		return fmt.Errorf("fetch daily HRV: %w", err)
	}

	var hrvResult DailyHRVResponse
	if err := json.Unmarshal(hrvResponse, &hrvResult); err != nil {
		return fmt.Errorf("parse daily HRV response: %w", err)
	}

	for _, point := range hrvResult.DataPoints {
		value := point.DailyHeartRateVariability
		date := value.Date

		if date.Year == today.Year() &&
			date.Month == int(today.Month()) &&
			date.Day == today.Day() {

			// Google Health returns average HRV as a JSON number.
			if value.AverageHRVMilliseconds == nil {
				continue
			}

			hrv := *value.AverageHRVMilliseconds
			hrvMS = &hrv

			fmt.Printf("Today's average HRV: %.3f ms\n", hrv)
			break
		}
	}

	if hrvMS == nil {
		fmt.Println("No HRV data available for today")
	}

	// 3. Save the available metrics.
	// Zones are not being fetched yet. Passing nil preserves their
	// existing values through the procedure's COALESCE logic.
	if restingHeartRate != nil || hrvMS != nil {
		err = h.Repo.UpsertDailyVitals(
			ctx,
			userID,
			today,
			restingHeartRate,
			hrvMS,
			nil, // Zone 1
			nil, // Zone 2
			nil, // Zone 3
			nil, // Zone 4
			nil, // Zone 5
		)
		if err != nil {
			return fmt.Errorf("save daily vitals: %w", err)
		}

		fmt.Println("Daily vitals saved successfully")
	} else {
		fmt.Println("No daily vitals data available to save")
	}

	return nil
}

func (h *Handler) SyncSleepData(ctx context.Context, userID int64) error {
	accessToken, err := h.GetValidAccessToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("get Google Health access token: %w", err)
	}

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return fmt.Errorf("load timezone: %w", err)
	}

	today := time.Now().In(loc)

	response, err := h.FetchHealthData(
		ctx,
		accessToken,
		"sleep",
		"sleepYesterday",
		today,
	)
	if err != nil {
		log.Printf("Failed to fetch Google Health sleep data: %v", err)
		return fmt.Errorf("fetch Google Health sleep data: %w", err)
	}

	log.Printf("Raw Google Health sleep response: %s", string(response))

	var sleepResponse SleepResponse
	if err := json.Unmarshal(response, &sleepResponse); err != nil {
		return fmt.Errorf("parse Google Health sleep response: %w", err)
	}

	savedCount := 0

	for _, point := range sleepResponse.DataPoints {
		sleep := point.Sleep

		startTime, err := time.Parse(time.RFC3339, sleep.Interval.StartTime)
		if err != nil {
			return fmt.Errorf("parse sleep start time for %s: %w", point.Name, err)
		}

		endTime, err := time.Parse(time.RFC3339, sleep.Interval.EndTime)
		if err != nil {
			return fmt.Errorf("parse sleep end time for %s: %w", point.Name, err)
		}

		durationMinutes := endTime.Sub(startTime).Minutes()

		minutesInSleepPeriod, err := strconv.ParseFloat(
			sleep.Summary.MinutesInSleepPeriod, 64,
		)
		if err != nil {
			return fmt.Errorf("parse minutes in sleep period for %s: %w", point.Name, err)
		}

		minutesAfterWakeUp, err := strconv.ParseFloat(
			sleep.Summary.MinutesAfterWakeUp, 64,
		)
		if err != nil {
			return fmt.Errorf("parse minutes after wake up for %s: %w", point.Name, err)
		}

		minutesToFallAsleep, err := strconv.ParseFloat(
			sleep.Summary.MinutesToFallAsleep, 64,
		)
		if err != nil {
			return fmt.Errorf("parse minutes to fall asleep for %s: %w", point.Name, err)
		}

		minutesAsleep, err := strconv.ParseFloat(
			sleep.Summary.MinutesAsleep, 64,
		)
		if err != nil {
			return fmt.Errorf("parse minutes asleep for %s: %w", point.Name, err)
		}

		minutesAwake, err := strconv.ParseFloat(
			sleep.Summary.MinutesAwake, 64,
		)
		if err != nil {
			return fmt.Errorf("parse minutes awake for %s: %w", point.Name, err)
		}

		var stagesStatus string
		if sleep.Metadata.StagesStatus != "" {
			stagesStatus = sleep.Metadata.StagesStatus
		}

		sleepSessionID, err := h.Repo.UpsertSleepSession(
			ctx,
			userID,
			point.Name,
			startTime,
			endTime,
			durationMinutes,
			sleep.Type,
			sleep.Metadata.Nap,
			sleep.Metadata.MainSleep,
			sleep.Metadata.Processed,
			sleep.Metadata.ManuallyEdited,
			stagesStatus,
			sleep.Metadata.ExternalID,
			minutesInSleepPeriod,
			minutesAfterWakeUp,
			minutesToFallAsleep,
			minutesAsleep,
			minutesAwake,
		)
		if err != nil {
			return fmt.Errorf("save sleep session %s: %w", point.Name, err)
		}

		// Save each individual sleep stage.
		for _, stage := range sleep.Stages {
			stageStartTime, err := time.Parse(time.RFC3339, stage.StartTime)
			if err != nil {
				return fmt.Errorf(
					"parse sleep stage start time for session %s: %w",
					point.Name,
					err,
				)
			}

			stageEndTime, err := time.Parse(time.RFC3339, stage.EndTime)
			if err != nil {
				return fmt.Errorf(
					"parse sleep stage end time for session %s: %w",
					point.Name,
					err,
				)
			}

			stageDurationMinutes := stageEndTime.Sub(stageStartTime).Minutes()
			_, err = h.Repo.UpsertSleepStage(
				ctx,
				sleepSessionID,
				stage.Type,
				stageStartTime,
				stageEndTime,
				stageDurationMinutes,
			)
			if err != nil {
				return fmt.Errorf(
					"save sleep stage for session %s: %w",
					point.Name,
					err,
				)
			}
		}
		savedCount++
	}
	fmt.Printf("Total sessions returned by API: %d\n", len(sleepResponse.DataPoints))
	fmt.Printf("Sleep sessions saved: %d\n", savedCount)
	fmt.Printf("Next page token present: %t\n", sleepResponse.NextPageToken != "")

	return nil
}
