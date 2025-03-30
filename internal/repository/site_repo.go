package repository

import (
	"database/sql"
	"errors"
	"time"

	"egide-server/internal/models"
)

type SiteRepository struct {
	db *sql.DB
}

func NewSiteRepository(db *sql.DB) *SiteRepository {
	return &SiteRepository{
		db: db,
	}
}

func (r *SiteRepository) Create(site *models.Site) (int64, error) {
	query := `
		INSERT INTO sites (user_id, domain, protection_mode, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(
		query,
		site.UserID,
		site.Domain,
		site.ProtectionMode,
		site.Active,
		now,
		now,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *SiteRepository) FindByID(id int64) (*models.Site, error) {
	query := `
		SELECT id, user_id, domain, protection_mode, active, created_at, updated_at
		FROM sites
		WHERE id = ?
	`

	var site models.Site
	var protectionMode string

	err := r.db.QueryRow(query, id).Scan(
		&site.ID,
		&site.UserID,
		&site.Domain,
		&protectionMode,
		&site.Active,
		&site.CreatedAt,
		&site.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("site not found")
		}
		return nil, err
	}

	site.ProtectionMode = models.ProtectionMode(protectionMode)
	return &site, nil
}

func (r *SiteRepository) FindByUserID(userID int64) ([]*models.Site, error) {
	query := `
		SELECT id, user_id, domain, protection_mode, active, created_at, updated_at
		FROM sites
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []*models.Site
	for rows.Next() {
		var site models.Site
		var protectionMode string

		err := rows.Scan(
			&site.ID,
			&site.UserID,
			&site.Domain,
			&protectionMode,
			&site.Active,
			&site.CreatedAt,
			&site.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		site.ProtectionMode = models.ProtectionMode(protectionMode)
		sites = append(sites, &site)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sites, nil
}

func (r *SiteRepository) FindByDomain(userID int64, domain string) (*models.Site, error) {
	query := `
		SELECT id, user_id, domain, protection_mode, active, created_at, updated_at
		FROM sites
		WHERE user_id = ? AND domain = ?
	`

	var site models.Site
	var protectionMode string

	err := r.db.QueryRow(query, userID, domain).Scan(
		&site.ID,
		&site.UserID,
		&site.Domain,
		&protectionMode,
		&site.Active,
		&site.CreatedAt,
		&site.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("site not found")
		}
		return nil, err
	}

	site.ProtectionMode = models.ProtectionMode(protectionMode)
	return &site, nil
}

func (r *SiteRepository) Update(site *models.Site) error {
	query := `
		UPDATE sites
		SET domain = ?, protection_mode = ?, active = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		site.Domain,
		site.ProtectionMode,
		site.Active,
		time.Now(),
		site.ID,
	)

	return err
}

func (r *SiteRepository) Delete(id int64) error {
	query := `DELETE FROM sites WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
