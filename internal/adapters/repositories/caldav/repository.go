package caldav

import (
	"context"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maypok86/otter/v2"
	"github.com/pkg/errors"

	"github.com/hexarchy/itmo-calendar/internal/entities"
	"github.com/hexarchy/itmo-calendar/pkg/transactions"
)

const (
	cacheTTL      = 5 * time.Minute
	cacheMaxSize  = 10_000
	cacheInitSize = 256
)

type Repository struct {
	db    *pgxpool.Pool
	cache *otter.Cache[int64, entities.CalDav]
}

func New(db *pgxpool.Pool) *Repository {
	cache := otter.Must(&otter.Options[int64, entities.CalDav]{
		MaximumSize:      cacheMaxSize,
		InitialCapacity:  cacheInitSize,
		ExpiryCalculator: otter.ExpiryWriting[int64, entities.CalDav](cacheTTL),
	})

	return &Repository{db: db, cache: cache}
}

// Create inserts or updates a user's iCal data and invalidates the cache.
func (r *Repository) Create(ctx context.Context, caldav entities.CalDav) error {
	const query = `
INSERT INTO caldav (isu, ical)
VALUES ($1, $2)
ON CONFLICT (isu) DO UPDATE SET ical = EXCLUDED.ical
`

	db := transactions.FromContext(ctx, r.db)
	_, err := db.Exec(ctx, query, caldav.ISU, []byte(caldav.ICal.Serialize()))
	if err != nil {
		return errors.Wrap(err, "caldav repository: create")
	}

	r.cache.Invalidate(caldav.ISU)
	return nil
}

// Get retrieves a user's iCal data by ISU with in-memory caching.
func (r *Repository) Get(ctx context.Context, isu int64) (entities.CalDav, error) {
	if v, ok := r.cache.GetIfPresent(isu); ok {
		return v, nil
	}

	const query = `SELECT isu, ical FROM caldav WHERE isu = $1`
	var caldav entities.CalDav
	var ical []byte
	db := transactions.FromContext(ctx, r.db)
	err := db.QueryRow(ctx, query, isu).Scan(&caldav.ISU, &ical)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.CalDav{}, errors.Wrapf(entities.ErrNotFound, "caldav repository: get isu %d", isu)
		}
		return entities.CalDav{}, errors.Wrap(err, "caldav repository: get")
	}

	caldav.ICal, err = ics.ParseCalendar(strings.NewReader(string(ical)))
	if err != nil {
		return entities.CalDav{}, errors.Wrap(err, "caldav repository: parse calendar")
	}

	r.cache.Set(isu, caldav)
	return caldav, nil
}
