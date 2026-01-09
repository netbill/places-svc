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

const PlacesClassesTable = "place_classes"

const PlacesClassesColumns = "id, parent_id, code, name, description, icon, created_at, updated_at"
const PlacesClassesColumnsP = "pc.id, pc.parent_id, pc.code, pc.name, pc.description, pc.icon, pc.created_at, pc.updated_at"

type PlaceClass struct {
	ID          uuid.UUID      `json:"id"`
	ParentID    uuid.NullUUID  `json:"parent_id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Icon        sql.NullString `json:"icon"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (c *PlaceClass) scan(row sq.RowScanner) error {
	if err := row.Scan(
		&c.ID,
		&c.ParentID,
		&c.Code,
		&c.Name,
		&c.Description,
		&c.Icon,
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		return fmt.Errorf("scanning place class: %w", err)
	}
	return nil
}

type PlacesClassesQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlacesClassesQ(db pgx.DBTX) PlacesClassesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return PlacesClassesQ{
		db:       db,
		selector: builder.Select(PlacesClassesColumnsP).From(PlacesClassesTable + " pc"),
		inserter: builder.Insert(PlacesClassesTable),
		updater:  builder.Update(PlacesClassesTable),
		deleter:  builder.Delete(PlacesClassesTable),
		counter:  builder.Select("COUNT(*) AS count").From(PlacesClassesTable + " pc"),
	}
}

type PlacesClassesQInsertInput struct {
	ParentID    *uuid.UUID // nil => root
	Code        string
	Name        string
	Description string
	Icon        *string
}

func (q PlacesClassesQ) Insert(ctx context.Context, data PlacesClassesQInsertInput) (PlaceClass, error) {
	set := map[string]interface{}{
		"code":        data.Code,
		"name":        data.Name,
		"description": data.Description,
		"icon":        data.Icon,
	}
	if data.ParentID != nil {
		set["parent_id"] = *data.ParentID
	} else {
		set["parent_id"] = nil
	}

	query, args, err := q.inserter.SetMap(set).
		Suffix("RETURNING " + PlacesClassesColumns).
		ToSql()
	if err != nil {
		return PlaceClass{}, fmt.Errorf("building insert query for %s: %w", PlacesClassesTable, err)
	}

	var inserted PlaceClass
	if err := inserted.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return PlaceClass{}, nil
		default:
			return PlaceClass{}, err
		}
	}

	return inserted, nil
}

func (q PlacesClassesQ) FilterByID(id uuid.UUID) PlacesClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.id": id})
	q.counter = q.counter.Where(sq.Eq{"pc.id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q PlacesClassesQ) FilterByParentID(parentID uuid.UUID) PlacesClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.parent_id": parentID})
	q.counter = q.counter.Where(sq.Eq{"pc.parent_id": parentID})
	q.updater = q.updater.Where(sq.Eq{"parent_id": parentID})
	q.deleter = q.deleter.Where(sq.Eq{"parent_id": parentID})
	return q
}

func (q PlacesClassesQ) FilterRoots() PlacesClassesQ {
	q.selector = q.selector.Where(sq.Expr("pc.parent_id IS NULL"))
	q.counter = q.counter.Where(sq.Expr("pc.parent_id IS NULL"))
	return q
}

func (q PlacesClassesQ) FilterByCode(code string) PlacesClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.code": code})
	q.counter = q.counter.Where(sq.Eq{"pc.code": code})
	q.updater = q.updater.Where(sq.Eq{"code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	return q
}

func (q PlacesClassesQ) FilterNameLike(name string) PlacesClassesQ {
	q.selector = q.selector.Where(sq.Like{"pc.name": "%" + name + "%"})
	q.counter = q.counter.Where(sq.Like{"pc.name": "%" + name + "%"})
	return q
}

func (q PlacesClassesQ) OrderName(asc bool) PlacesClassesQ {
	if asc {
		q.selector = q.selector.OrderBy("pc.name ASC", "pc.id ASC")
	} else {
		q.selector = q.selector.OrderBy("pc.name DESC", "pc.id DESC")
	}
	return q
}

func (q PlacesClassesQ) Get(ctx context.Context) (PlaceClass, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceClass{}, fmt.Errorf("building select query for %s: %w", PlacesClassesTable, err)
	}

	var c PlaceClass
	if err := c.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return PlaceClass{}, err
	}
	return c, nil
}

func (q PlacesClassesQ) Select(ctx context.Context) ([]PlaceClass, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlacesClassesTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PlacesClassesTable, err)
	}
	defer rows.Close()

	var items []PlaceClass
	for rows.Next() {
		var c PlaceClass
		if err := c.scan(rows); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (q PlacesClassesQ) Page(limit, offset uint) PlacesClassesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q PlacesClassesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlacesClassesTable, err)
	}

	var count uint
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PlacesClassesTable, err)
	}

	return count, nil
}

func (q PlacesClassesQ) UpdateOne(ctx context.Context) (PlaceClass, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.
		Suffix("RETURNING " + PlacesClassesColumns).
		ToSql()
	if err != nil {
		return PlaceClass{}, fmt.Errorf("building update query for %s: %w", PlacesClassesTable, err)
	}

	var updated PlaceClass
	if err := updated.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return PlaceClass{}, err
	}

	return updated, nil
}

func (q PlacesClassesQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", PlacesClassesTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", PlacesClassesTable, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", PlacesClassesTable, err)
	}

	return affected, nil
}

func (q PlacesClassesQ) UpdateParent(parentID *uuid.UUID) PlacesClassesQ {
	if parentID == nil {
		q.updater = q.updater.Set("parent_id", nil)
	} else {
		q.updater = q.updater.Set("parent_id", *parentID)
	}
	return q
}

func (q PlacesClassesQ) UpdateCode(code string) PlacesClassesQ {
	q.updater = q.updater.Set("code", code)
	return q
}

func (q PlacesClassesQ) UpdateName(name string) PlacesClassesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q PlacesClassesQ) UpdateDescription(description string) PlacesClassesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q PlacesClassesQ) UpdateIcon(icon *string) PlacesClassesQ {
	if icon == nil {
		q.updater = q.updater.Set("icon", nil)
	} else {
		q.updater = q.updater.Set("icon", *icon)
	}
	return q
}

func (q PlacesClassesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlacesClassesTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PlacesClassesTable, err)
	}

	return nil
}
