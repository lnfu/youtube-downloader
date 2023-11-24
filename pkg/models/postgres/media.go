package postgres

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/lnfu/youtube-downloader/pkg/models"
)

type MediaModel struct {
	DB *sql.DB
}

func (m *MediaModel) Insert(originId, mediaType, accessKey string) error {
	stmt, err := m.DB.Prepare(`
		INSERT INTO media(vid, type, media_status, access_key)
		VALUES($1, $2, 'running', $3)
		RETURNING id;
	`)

	_, err = stmt.Exec(originId, mediaType, accessKey)

	return err
}

func (m *MediaModel) DownloadComplete(originId, mediaType string) error {
	stmt, err := m.DB.Prepare(`
		UPDATE media
		SET media_status = 'done'
		WHERE vid = $1 AND type = $2;
	`)
	_, err = stmt.Exec(originId, mediaType)
	return err
}

func (m *MediaModel) DownloadFailure(originId, mediaType string) error {
	stmt, err := m.DB.Prepare(`
		UPDATE media
		SET media_status = 'failure'
		WHERE vid = $1 AND type = $2;
	`)
	_, err = stmt.Exec(originId, mediaType)
	return err
}

func (m *MediaModel) Get(originId, mediaType string, currentTime time.Time) (*models.Media, error) {

	stmt := `
		SELECT *
		FROM media
		WHERE vid = $1 AND type = $2;
	`

	row := m.DB.QueryRow(stmt, originId, mediaType)

	mm := &models.Media{}

	err := row.Scan(
		&mm.Id,
		&mm.OriginId,
		&mm.Type,
		&mm.MediaStatus,
		&mm.AccessKey,
		&mm.RecentlyAccessTime,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	stmt_2, err := m.DB.Prepare(`
		UPDATE media
		SET recently_access_time = $3
		WHERE vid = $1 AND type = $2;
	`)
	_, err = stmt_2.Exec(originId, mediaType, currentTime)

	return mm, nil

}
