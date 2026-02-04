package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/netbill/pgdbx"
	"github.com/netbill/places-svc/internal/repository"
)

const organizationsTable = "organizations"

const organizationsColumns = "id, status, verified, name, icon, banner, source_created_at, source_updated_at, replica_created_at, replica_updated_at"
const organizationsColumnsP = "o.id, o.status, o.verified, o.name, o.icon, o.banner, o.source_created_at, o.source_updated_at, o.replica_created_at, o.replica_updated_at"

func scanOrganization(row sq.RowScanner) (o repository.OrganizationRow, err error) {
	var icon pgtype.Text
	var banner pgtype.Text

	err = row.Scan(
		&o.ID,
		&o.Status,
		&o.Verified,
		&o.Name,
		&icon,
		&banner,
		&o.SourceCreatedAt,
		&o.SourceUpdatedAt,
		&o.ReplicaCreatedAt,
		&o.ReplicaUpdatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrganizationRow{}, nil
	case err != nil:
		return repository.OrganizationRow{}, fmt.Errorf("scanning organization: %w", err)
	}

	if icon.Valid {
		o.Icon = &icon.String
	}
	if banner.Valid {
		o.Banner = &banner.String
	}

	return o, nil
}

type organizations struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrganizationsQ(db *pgdbx.DB) repository.OrganizationsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &organizations{
		db:       db,
		selector: b.Select(organizationsColumnsP).From(organizationsTable + " o"),
		inserter: b.Insert(organizationsTable),
		updater:  b.Update(organizationsTable + " o"),
		deleter:  b.Delete(organizationsTable + " o"),
		counter:  b.Select("COUNT(*)").From(organizationsTable + " o"),
	}
}

func (q *organizations) New() repository.OrganizationsQ {
	return NewOrganizationsQ(q.db)
}

func (q *organizations) Insert(ctx context.Context, data repository.OrganizationRow) (repository.OrganizationRow, error) {
	now := time.Now().UTC()

	query, args, err := q.inserter.SetMap(map[string]any{
		"id":                 data.ID,
		"status":             data.Status,
		"verified":           data.Verified,
		"name":               data.Name,
		"icon":               data.Icon,
		"banner":             data.Banner,
		"source_created_at":  data.SourceCreatedAt.UTC(),
		"source_updated_at":  data.SourceUpdatedAt.UTC(),
		"replica_created_at": now,
		"replica_updated_at": now,
	}).Suffix("RETURNING " + organizationsColumns).ToSql()
	if err != nil {
		return repository.OrganizationRow{}, fmt.Errorf("building insert query for %s: %w", organizationsTable, err)
	}

	return scanOrganization(q.db.QueryRow(ctx, query, args...))
}

func (q *organizations) Get(ctx context.Context) (repository.OrganizationRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrganizationRow{}, fmt.Errorf("building select query for %s: %w", organizationsTable, err)
	}

	return scanOrganization(q.db.QueryRow(ctx, query, args...))
}

func (q *organizations) Select(ctx context.Context) ([]repository.OrganizationRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", organizationsTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", organizationsTable, err)
	}
	defer rows.Close()

	out := make([]repository.OrganizationRow, 0)
	for rows.Next() {
		o, err := scanOrganization(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *organizations) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationsTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationsTable, err)
	}

	return nil
}

func (q *organizations) FilterByID(id ...uuid.UUID) repository.OrganizationsQ {
	q.selector = q.selector.Where(sq.Eq{"o.id": id})
	q.counter = q.counter.Where(sq.Eq{"o.id": id})
	q.updater = q.updater.Where(sq.Eq{"o.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"o.id": id})
	return q
}

func (q *organizations) FilterByStatus(status string) repository.OrganizationsQ {
	q.selector = q.selector.Where(sq.Eq{"o.status": status})
	q.counter = q.counter.Where(sq.Eq{"o.status": status})
	q.updater = q.updater.Where(sq.Eq{"o.status": status})
	q.deleter = q.deleter.Where(sq.Eq{"o.status": status})
	return q
}

func (q *organizations) FilterByVerified(verified bool) repository.OrganizationsQ {
	q.selector = q.selector.Where(sq.Eq{"o.verified": verified})
	q.counter = q.counter.Where(sq.Eq{"o.verified": verified})
	q.updater = q.updater.Where(sq.Eq{"o.verified": verified})
	q.deleter = q.deleter.Where(sq.Eq{"o.verified": verified})
	return q
}

func (q *organizations) FilterByName(name string) repository.OrganizationsQ {
	q.selector = q.selector.Where(sq.Eq{"o.name": name})
	q.counter = q.counter.Where(sq.Eq{"o.name": name})
	q.updater = q.updater.Where(sq.Eq{"o.name": name})
	q.deleter = q.deleter.Where(sq.Eq{"o.name": name})
	return q
}

func (q *organizations) UpdateOne(ctx context.Context) (repository.OrganizationRow, error) {
	q.updater = q.updater.Set("replica_updated_at", time.Now().UTC())

	query, args, err := q.updater.Suffix("RETURNING " + organizationsColumns).ToSql()
	if err != nil {
		return repository.OrganizationRow{}, fmt.Errorf("building update query for %s: %w", organizationsTable, err)
	}

	return scanOrganization(q.db.QueryRow(ctx, query, args...))
}

func (q *organizations) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("replica_updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", organizationsTable, err)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", organizationsTable, err)
	}

	return res.RowsAffected(), nil
}

func (q *organizations) UpdateStatus(status string) repository.OrganizationsQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q *organizations) UpdateVerified(verified bool) repository.OrganizationsQ {
	q.updater = q.updater.Set("verified", verified)
	return q
}

func (q *organizations) UpdateName(name string) repository.OrganizationsQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q *organizations) UpdateIcon(icon *string) repository.OrganizationsQ {
	q.updater = q.updater.Set("icon", icon)
	return q
}

func (q *organizations) UpdateBanner(banner *string) repository.OrganizationsQ {
	q.updater = q.updater.Set("banner", banner)
	return q
}

func (q *organizations) UpdateSourceUpdatedAt(updatedAt time.Time) repository.OrganizationsQ {
	q.updater = q.updater.Set("source_updated_at", updatedAt.UTC())
	return q
}
