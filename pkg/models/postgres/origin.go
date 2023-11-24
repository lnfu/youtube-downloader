package postgres

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"github.com/lnfu/youtube-downloader/pkg/models"
)

type OriginModel struct {
	DB *sql.DB
}

func (m *OriginModel) Insert(id string) error {
	stmt, err := m.DB.Prepare(`
		INSERT INTO origin(id, title, duration, info_status)
		VALUES($1, NULL, NULL, 'running');
	`)
	_, err = stmt.Exec(id)

	return err
}

func (m *OriginModel) UpdateInfo(id, title string, duration int) error {
	stmt, err := m.DB.Prepare(`
		UPDATE origin
		SET title = $2, duration = $3, info_status = 'done'
		WHERE id = $1;
	`)
	_, err = stmt.Exec(id, title, duration)
	return err
}

func (m *OriginModel) SetStatusFailure(id string) error {
	stmt, err := m.DB.Prepare(`
		UPDATE origin
		SET status = 'failure'
		WHERE id = $1;
	`)
	_, err = stmt.Exec(id)
	return err
}

func (m *OriginModel) Get(id string) (*models.Origin, error) {

	stmt := `
		SELECT *
		FROM origin
		WHERE id = $1;
	`

	row := m.DB.QueryRow(stmt, id)

	mo := &models.Origin{}

	err := row.Scan(
		&mo.Id,
		&mo.Title,
		&mo.Duration,
		&mo.InfoStatus,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return mo, nil

}
