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
const organizationsColumns = "id, status, name, icon_key, banner_key, version, source_created_at, source_updated_at, replica_created_at, replica_updated_at"
const organizationsColumnsO = "o.id, o.status, o.name, o.icon_key, o.banner_key, o.version, o.source_created_at, o.source_updated_at, o.replica_created_at, o.replica_updated_at"

func scanOrganization(row sq.RowScanner) (o repository.OrganizationRow, err error) {
	var iconKey pgtype.Text
	var bannerKey pgtype.Text

	err = row.Scan(
		&o.ID,
		&o.Status,
		&o.Name,
		&iconKey,
		&bannerKey,
		&o.Version,
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

	if iconKey.Valid {
		o.IconKey = &iconKey.String
	}
	if bannerKey.Valid {
		o.BannerKey = &bannerKey.String
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
		selector: b.Select(organizationsColumnsO).From(organizationsTable + " o"),
		inserter: b.Insert(organizationsTable),
		updater:  b.Update(organizationsTable + " o"),
		deleter:  b.Delete(organizationsTable + " o"),
		counter:  b.Select("COUNT(*)").From(organizationsTable + " o"),
	}
}

func (q *organizations) New() repository.OrganizationsQ {
	return NewOrganizationsQ(q.db)
}

func (q *organizations) Insert(ctx context.Context, data repository.OrganizationRow) error {
	query, args, err := q.inserter.SetMap(map[string]any{
		"id":                data.ID,
		"status":            data.Status,
		"name":              data.Name,
		"icon_key":          data.IconKey,
		"banner_key":        data.BannerKey,
		"version":           data.Version,
		"source_created_at": data.SourceCreatedAt.UTC(),
		"source_updated_at": data.SourceUpdatedAt.UTC(),
	}).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", organizationsTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing insert query for %s: %w", organizationsTable, err)
	}

	return nil
}

func (q *organizations) Get(ctx context.Context) (repository.OrganizationRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrganizationRow{}, fmt.Errorf("building select query for %s: %w", organizationsTable, err)
	}

	return scanOrganization(q.db.QueryRow(ctx, query, args...))
}

func (q *organizations) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, err
	}
	sql := "SELECT EXISTS (" + subSQL + ")"

	var exists bool
	if err = q.db.QueryRow(ctx, sql, subArgs...).Scan(&exists); err != nil {
		return false, fmt.Errorf("sql=%s args=%v: %w", sql, subArgs, err)
	}
	return exists, nil
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

func (q *organizations) UpdateOne(ctx context.Context) error {
	q.updater = q.updater.Set("replica_updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", organizationsTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing update query for %s: %w", organizationsTable, err)
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

func (q *organizations) FilterByName(name string) repository.OrganizationsQ {
	q.selector = q.selector.Where(sq.Eq{"o.name": name})
	q.counter = q.counter.Where(sq.Eq{"o.name": name})
	q.updater = q.updater.Where(sq.Eq{"o.name": name})
	q.deleter = q.deleter.Where(sq.Eq{"o.name": name})
	return q
}

func (q *organizations) UpdateStatus(status string) repository.OrganizationsQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q *organizations) UpdateName(name string) repository.OrganizationsQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q *organizations) UpdateIcon(icon *string) repository.OrganizationsQ {
	q.updater = q.updater.Set("icon_key", icon)
	return q
}

func (q *organizations) UpdateBanner(banner *string) repository.OrganizationsQ {
	q.updater = q.updater.Set("banner_key", banner)
	return q
}

func (q *organizations) UpdateVersion(version int32) repository.OrganizationsQ {
	q.updater = q.updater.Set("version", version)
	return q
}

func (q *organizations) UpdateSourceUpdatedAt(updatedAt time.Time) repository.OrganizationsQ {
	q.updater = q.updater.Set("source_updated_at", updatedAt.UTC())
	return q
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
