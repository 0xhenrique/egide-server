package repository

import (
	"database/sql"
	"time"

	"egide-server/internal/models"
)

type HealthCheckRepository struct {
	db *sql.DB
}

func NewHealthCheckRepository(db *sql.DB) *HealthCheckRepository {
	return &HealthCheckRepository{
		db: db,
	}
}

func (r *HealthCheckRepository) Create(check *models.HealthCheck) (int64, error) {
	query := `
		INSERT INTO health_checks (timestamp, response_time_ms, status_code, success, error, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(
		query,
		check.Timestamp,
		check.ResponseTimeMs,
		check.StatusCode,
		check.Success,
		check.Error,
		now,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetChecksInRange returns health checks within a time range
func (r *HealthCheckRepository) GetChecksInRange(start, end time.Time) ([]*models.HealthCheck, error) {
	query := `
		SELECT id, timestamp, response_time_ms, status_code, success, error, created_at
		FROM health_checks
		WHERE timestamp BETWEEN ? AND ?
		ORDER BY timestamp ASC
	`

	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []*models.HealthCheck
	for rows.Next() {
		var check models.HealthCheck
		err := rows.Scan(
			&check.ID,
			&check.Timestamp,
			&check.ResponseTimeMs,
			&check.StatusCode,
			&check.Success,
			&check.Error,
			&check.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		checks = append(checks, &check)
	}

	return checks, rows.Err()
}

// DeleteOldChecks removes health checks older than the specified duration
func (r *HealthCheckRepository) DeleteOldChecks(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	query := `DELETE FROM health_checks WHERE timestamp < ?`
	_, err := r.db.Exec(query, cutoff)
	return err
}
