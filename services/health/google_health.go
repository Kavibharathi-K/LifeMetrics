package health

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const googleHealthAPIBaseURL = "https://health.googleapis.com/v4"

func (h *Handler) GetSteps(
	ctx context.Context,
	userID int64,
) ([]byte, error) {

	accessToken, err := h.GetValidAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Test date: yesterday
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return nil, err
	}

	yesterday := time.Now().In(loc).AddDate(0, 0, -1)
	tomorrow := yesterday.AddDate(0, 0, 1)

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
		yesterday.Year(),
		int(yesterday.Month()),
		yesterday.Day(),
		tomorrow.Year(),
		int(tomorrow.Month()),
		tomorrow.Day(),
	)

	url := googleHealthAPIBaseURL +
		"/users/me/dataTypes/steps/dataPoints:dailyRollUp"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBufferString(requestBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+accessToken,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"Accept",
		"application/json",
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"google health API returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	return body, nil
}
