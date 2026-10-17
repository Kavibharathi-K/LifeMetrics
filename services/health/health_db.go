package health

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateOAuthState(
	ctx context.Context,
	state string,
	userID int64,
	expiresAt time.Time,
) error {

	_, err := r.DB.Exec(
		ctx,
		"CALL health_oauth_state_add($1, $2, $3)",
		state,
		userID,
		expiresAt,
	)

	return err
}

func (r *Repository) GetOAuthState(
	ctx context.Context,
	state string,
) (int64, time.Time, error) {

	var userID int64
	var expiresAt time.Time

	err := r.DB.QueryRow(
		ctx,
		"SELECT * FROM health_oauth_state_get($1)",
		state,
	).Scan(
		&userID,
		&expiresAt,
	)

	return userID, expiresAt, err
}

func (r *Repository) DeleteOAuthState(
	ctx context.Context,
	state string,
) error {

	_, err := r.DB.Exec(
		ctx,
		"CALL health_oauth_state_delete($1)",
		state,
	)

	return err
}

func (r *Repository) UpsertHealthConnection(
	ctx context.Context,
	userID int64,
	provider string,
	providerUserID string,
	accessToken string,
	refreshToken string,
	expiresAt time.Time,
) error {

	_, err := r.DB.Exec(
		ctx,
		"CALL health_connection_upsert($1, $2, $3, $4, $5, $6)",
		userID,
		provider,
		providerUserID,
		accessToken,
		refreshToken,
		expiresAt,
	)

	return err
}

func (r *Repository) GetHealthConnection(
	ctx context.Context,
	userID int64,
	provider string,
) (*HealthConnection, error) {

	connection := &HealthConnection{}

	err := r.DB.QueryRow(
		ctx,
		"SELECT * FROM health_connection_get($1, $2)",
		userID,
		provider,
	).Scan(
		&connection.ProviderUserID,
		&connection.AccessToken,
		&connection.RefreshToken,
		&connection.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (r *Repository) UpsertDailyActivity(
	ctx context.Context,
	userID int64,
	activityDate time.Time,
	steps *int64,
	distanceMeters *float64,
	activeEnergyKcal *float64,
	totalCaloriesKcal *float64,
	activeMinutes *float64,
) error {
	_, err := r.DB.Exec(
		ctx,
		`CALL health_daily_activity_upsert($1, $2, $3, $4, $5, $6, $7)`,
		userID,
		activityDate,
		steps,
		distanceMeters,
		activeEnergyKcal,
		totalCaloriesKcal,
		activeMinutes,
	)

	return err
}

func (r *Repository) UpsertDailyVitals(

	ctx context.Context,
	userID int64,
	vitalDate time.Time,
	restingHeartRateBPM *float64,
	hrvMS *float64,
	zone1Minutes *float64,
	zone2Minutes *float64,
	zone3Minutes *float64,
	zone4Minutes *float64,
	zone5Minutes *float64,

) error {

	_, err := r.DB.Exec(
		ctx,
		`CALL health_daily_vitals_upsert($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		userID,
		vitalDate,
		restingHeartRateBPM,
		hrvMS,
		zone1Minutes,
		zone2Minutes,
		zone3Minutes,
		zone4Minutes,
		zone5Minutes,
	)
	return err

}

func (r *Repository) UpsertSleepSession(
	ctx context.Context,
	userID int64,
	googleDataPointID string,
	startTime time.Time,
	endTime time.Time,
	durationMinutes float64,
	sleepType string,
	isNap bool,
	isMainSleep bool,
	processed bool,
	manuallyEdited bool,
	stagesStatus string,
	externalID string,
	minutesInSleepPeriod float64,
	minutesAfterWakeUp float64,
	minutesToFallAsleep float64,
	minutesAsleep float64,
	minutesAwake float64,
) (int64, error) {
	var sleepSessionID int64

	err := r.DB.QueryRow(
		ctx,
		`SELECT health_sleep_session_upsert(
            $1, $2, $3, $4, $5, $6, $7, $8, $9,
            $10, $11, $12, $13, $14, $15, $16, $17
        )`,
		userID,
		googleDataPointID,
		startTime,
		endTime,
		durationMinutes,
		sleepType,
		isNap,
		isMainSleep,
		processed,
		manuallyEdited,
		stagesStatus,
		externalID,
		minutesInSleepPeriod,
		minutesAfterWakeUp,
		minutesToFallAsleep,
		minutesAsleep,
		minutesAwake,
	).Scan(&sleepSessionID)

	if err != nil {
		return 0, fmt.Errorf("upsert sleep session: %w", err)
	}

	return sleepSessionID, nil
}

func (r *Repository) UpsertSleepStage(
	ctx context.Context,
	sleepSessionID int64,
	stage string,
	startTime time.Time,
	endTime time.Time,
	durationMinutes float64,
) (int64, error) {
	var sleepStageID int64

	err := r.DB.QueryRow(
		ctx,
		`SELECT health_sleep_stage_upsert($1, $2, $3, $4, $5)`,
		sleepSessionID,
		stage,
		startTime,
		endTime,
		durationMinutes,
	).Scan(&sleepStageID)

	if err != nil {
		return 0, fmt.Errorf("upsert sleep stage: %w", err)
	}

	return sleepStageID, nil
}
