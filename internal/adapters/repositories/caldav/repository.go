package caldav

import (
	"context"
	"strings"

	"github.com/hexarchy/itmo-calendar/internal/entities"

	ics "github.com/arran4/golang-ical"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserts or updates a user's iCal data.
func (r *Repository) Create(ctx context.Context, caldav entities.CalDav) error {
	const query = `
INSERT INTO caldav (isu, ical)
VALUES ($1, $2)
ON CONFLICT (isu) DO UPDATE SET ical = EXCLUDED.ical
`

	_, err := r.db.Exec(ctx, query, caldav.ISU, []byte(caldav.ICal.Serialize()))
	if err != nil {
		return errors.Wrap(err, "caldav repository: create")
	}
	return nil
}

// Get retrieves a user's iCal data by ISU.
func (r *Repository) Get(ctx context.Context, isu int64) (entities.CalDav, error) {
	const query = `SELECT isu, ical FROM caldav WHERE isu = $1`
	var caldav entities.CalDav
	var ical []byte
	err := r.db.QueryRow(ctx, query, isu).Scan(&caldav.ISU, &ical)
	if err != nil {
		return entities.CalDav{}, errors.Wrap(err, "caldav repository: get")
	}

	caldav.ICal, err = ics.ParseCalendar(strings.NewReader(string(ical)))
	if err != nil {
		return entities.CalDav{}, errors.Wrap(err, "caldav repository: parse calendar")
	}

	return caldav, nil
}
