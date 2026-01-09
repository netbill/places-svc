package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/netbill/pgx"
)

const PlacesTable = "places"

const PlacesColumns = `
	id,
	parent_id,
	organization_id,
	class_id,
	status,
	verified,
	point,
	address,
	name,
	description,
	icon,
	banner,
	website,
	phone,
	created_at,
	updated_at
`

type Place struct {
	ID             uuid.UUID     `json:"id"`
	ParentID       uuid.NullUUID `json:"parent_id"`
	OrganizationID uuid.UUID     `json:"organization_id"`
	ClassID        uuid.UUID     `json:"class_id"`

	Status   string `json:"status"`
	Verified bool   `json:"verified"`

	Point   any    `json:"point"` // geography, как и в других местах
	Address string `json:"address"`

	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Icon        sql.NullString `json:"icon"`
	Banner      sql.NullString `json:"banner"`
	Website     sql.NullString `json:"website"`
	Phone       sql.NullString `json:"phone"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Place) scan(row sq.RowScanner) error {
	if err := row.Scan(
		&p.ID,
		&p.ParentID,
		&p.OrganizationID,
		&p.ClassID,
		&p.Status,
		&p.Verified,
		&p.Point,
		&p.Address,
		&p.Name,
		&p.Description,
		&p.Icon,
		&p.Banner,
		&p.Website,
		&p.Phone,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return fmt.Errorf("scanning place: %w", err)
	}
	return nil
}

type PlacesQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlacesQ(db pgx.DBTX) PlacesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return PlacesQ{
		db:       db,
		selector: builder.Select(PlacesColumns).From(PlacesTable),
		inserter: builder.Insert(PlacesTable),
		updater:  builder.Update(PlacesTable),
		deleter:  builder.Delete(PlacesTable),
		counter:  builder.Select("COUNT(*) AS count").From(PlacesTable),
	}
}

type PlacesQInsertInput struct {
	ParentID       *uuid.UUID
	OrganizationID uuid.UUID
	ClassID        uuid.UUID

	Status   string
	Verified bool

	Point   any
	Address string

	Name        string
	Description *string
	Icon        *string
	Banner      *string
	Website     *string
	Phone       *string
}

func (q PlacesQ) Insert(ctx context.Context, data PlacesQInsertInput) (Place, error) {
	set := map[string]interface{}{
		"organization_id": data.OrganizationID,
		"class_id":        data.ClassID,
		"status":          data.Status,
		"verified":        data.Verified,
		"point":           data.Point,
		"address":         data.Address,
		"name":            data.Name,
		"description":     data.Description,
		"icon":            data.Icon,
		"banner":          data.Banner,
		"website":         data.Website,
		"phone":           data.Phone,
	}

	if data.ParentID != nil {
		set["parent_id"] = *data.ParentID
	}

	query, args, err := q.inserter.
		SetMap(set).
		Suffix("RETURNING " + PlacesColumns).
		ToSql()
	if err != nil {
		return Place{}, fmt.Errorf("building insert query for %s: %w", PlacesTable, err)
	}

	var inserted Place
	if err := inserted.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Place{}, nil
		default:
			return Place{}, err
		}
	}

	return inserted, nil
}

func (q PlacesQ) FilterByID(id uuid.UUID) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q PlacesQ) FilterByOrganizationID(orgID uuid.UUID) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"organization_id": orgID})
	q.counter = q.counter.Where(sq.Eq{"organization_id": orgID})
	q.updater = q.updater.Where(sq.Eq{"organization_id": orgID})
	q.deleter = q.deleter.Where(sq.Eq{"organization_id": orgID})
	return q
}

func (q PlacesQ) FilterByClassID(classID uuid.UUID) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"class_id": classID})
	q.counter = q.counter.Where(sq.Eq{"class_id": classID})
	return q
}

func (q PlacesQ) FilterByStatus(status string) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q PlacesQ) FilterVerified(verified bool) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"verified": verified})
	q.counter = q.counter.Where(sq.Eq{"verified": verified})
	return q
}

func (q PlacesQ) FilterByParentID(parentID uuid.UUID) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"parent_id": parentID})
	q.counter = q.counter.Where(sq.Eq{"parent_id": parentID})
	return q
}

func (q PlacesQ) Get(ctx context.Context) (Place, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Place{}, fmt.Errorf("building select query for %s: %w", PlacesTable, err)
	}

	var p Place
	if err := p.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Place{}, err
	}

	return p, nil
}

func (q PlacesQ) Select(ctx context.Context) ([]Place, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlacesTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PlacesTable, err)
	}
	defer rows.Close()

	var places []Place
	for rows.Next() {
		var p Place
		if err := p.scan(rows); err != nil {
			return nil, err
		}
		places = append(places, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return places, nil
}

func (q PlacesQ) UpdateOne(ctx context.Context) (Place, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.
		Suffix("RETURNING " + PlacesColumns).
		ToSql()
	if err != nil {
		return Place{}, fmt.Errorf("building update query for %s: %w", PlacesTable, err)
	}

	var updated Place
	if err := updated.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Place{}, err
	}

	return updated, nil
}

func (q PlacesQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", PlacesTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", PlacesTable, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", PlacesTable, err)
	}

	return affected, nil
}

func (q PlacesQ) UpdateStatus(status string) PlacesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q PlacesQ) UpdateVerified(verified bool) PlacesQ {
	q.updater = q.updater.Set("verified", verified)
	return q
}

func (q PlacesQ) UpdateName(name string) PlacesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q PlacesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlacesTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PlacesTable, err)
	}

	return nil
}

func (q PlacesQ) Page(limit, offset uint) PlacesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q PlacesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlacesTable, err)
	}

	var count uint
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PlacesTable, err)
	}

	return count, nil
}
