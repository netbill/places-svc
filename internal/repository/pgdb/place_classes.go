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

const PlaceClassesTable = "place_classes"

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

type PlaceClassesQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceClassesQ(db pgx.DBTX) PlaceClassesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return PlaceClassesQ{
		db:       db,
		selector: builder.Select(PlacesClassesColumnsP).From(PlaceClassesTable + " pc"),
		inserter: builder.Insert(PlaceClassesTable),
		updater:  builder.Update(PlaceClassesTable),
		deleter:  builder.Delete(PlaceClassesTable),
		counter:  builder.Select("COUNT(*) AS count").From(PlaceClassesTable + " pc"),
	}
}

type PlaceClassesInsertInput struct {
	ParentID    *uuid.UUID // nil => root
	Code        string
	Name        string
	Description string
	Icon        *string
}

func (q PlaceClassesQ) Insert(ctx context.Context, data PlaceClassesInsertInput) (PlaceClass, error) {
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
		return PlaceClass{}, fmt.Errorf("building insert query for %s: %w", PlaceClassesTable, err)
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

func (q PlaceClassesQ) FilterByID(id uuid.UUID) PlaceClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.id": id})
	q.counter = q.counter.Where(sq.Eq{"pc.id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q PlaceClassesQ) FilterByParentID(parentID uuid.UUID) PlaceClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.parent_id": parentID})
	q.counter = q.counter.Where(sq.Eq{"pc.parent_id": parentID})
	q.updater = q.updater.Where(sq.Eq{"parent_id": parentID})
	q.deleter = q.deleter.Where(sq.Eq{"parent_id": parentID})
	return q
}

func (q PlaceClassesQ) FilterRoots() PlaceClassesQ {
	q.selector = q.selector.Where(sq.Expr("pc.parent_id IS NULL"))
	q.counter = q.counter.Where(sq.Expr("pc.parent_id IS NULL"))
	return q
}

func (q PlaceClassesQ) FilterByCode(code string) PlaceClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.code": code})
	q.counter = q.counter.Where(sq.Eq{"pc.code": code})
	q.updater = q.updater.Where(sq.Eq{"code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	return q
}

func (q PlaceClassesQ) FilterNameLike(name string) PlaceClassesQ {
	q.selector = q.selector.Where(sq.Like{"pc.name": "%" + name + "%"})
	q.counter = q.counter.Where(sq.Like{"pc.name": "%" + name + "%"})
	return q
}

func (q PlaceClassesQ) FilterByParentIDTree(parentID uuid.UUID, maxGeneration int) PlaceClassesQ {
	// maxGeneration: 0 = unlimited
	var depthCond string
	if maxGeneration > 0 {
		// depth starts from 1 for direct children
		depthCond = fmt.Sprintf("WHERE t.depth <= %d", maxGeneration)
	} else {
		depthCond = ""
	}

	// descendants of parentID (excluding the parent itself)
	// t.depth = 1 => direct child
	exprSQL := fmt.Sprintf(`
		pc.id IN (
			WITH RECURSIVE t AS (
				SELECT c.id, c.parent_id, 1 AS depth
				FROM %s c
				WHERE c.parent_id = $1
				UNION ALL
				SELECT c2.id, c2.parent_id, t.depth + 1
				FROM %s c2
				JOIN t ON c2.parent_id = t.id
			)
			SELECT t.id FROM t
			%s
		)
	`, PlaceClassesTable, PlaceClassesTable, depthCond)

	cond := sq.Expr(exprSQL, parentID)

	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}

func (q PlaceClassesQ) FilterByChildIDTree(childID uuid.UUID, maxGeneration int) PlaceClassesQ {
	// maxGeneration: 0 = unlimited
	var depthCond string
	if maxGeneration > 0 {
		// depth starts from 1 for direct parent
		depthCond = fmt.Sprintf("WHERE t.depth <= %d", maxGeneration)
	} else {
		depthCond = ""
	}

	// ancestors of childID (excluding the child itself)
	// t.depth = 1 => direct parent
	exprSQL := fmt.Sprintf(`
		pc.id IN (
			WITH RECURSIVE t AS (
				SELECT p.id, p.parent_id, 1 AS depth
				FROM %s c
				JOIN %s p ON p.id = c.parent_id
				WHERE c.id = $1 AND c.parent_id IS NOT NULL
				UNION ALL
				SELECT p2.id, p2.parent_id, t.depth + 1
				FROM %s p2
				JOIN t ON p2.id = t.parent_id
			)
			SELECT t.id FROM t
			%s
		)
	`, PlaceClassesTable, PlaceClassesTable, PlaceClassesTable, depthCond)

	cond := sq.Expr(exprSQL, childID)

	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}

func (q PlaceClassesQ) OrderName(asc bool) PlaceClassesQ {
	if asc {
		q.selector = q.selector.OrderBy("pc.name ASC", "pc.id ASC")
	} else {
		q.selector = q.selector.OrderBy("pc.name DESC", "pc.id DESC")
	}
	return q
}

func (q PlaceClassesQ) Get(ctx context.Context) (PlaceClass, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceClass{}, fmt.Errorf("building select query for %s: %w", PlaceClassesTable, err)
	}

	var c PlaceClass
	if err := c.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return PlaceClass{}, err
	}
	return c, nil
}

func (q PlaceClassesQ) Select(ctx context.Context) ([]PlaceClass, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlaceClassesTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PlaceClassesTable, err)
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

func (q PlaceClassesQ) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", PlaceClassesTable, err)
	}

	sqlq := "SELECT EXISTS (" + subSQL + ")"

	var ok bool
	if err = q.db.QueryRowContext(ctx, sqlq, subArgs...).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning exists for %s: %w", PlaceClassesTable, err)
	}

	return ok, nil
}

func (q PlaceClassesQ) Page(limit, offset uint) PlaceClassesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q PlaceClassesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlaceClassesTable, err)
	}

	var count uint
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PlaceClassesTable, err)
	}

	return count, nil
}

func (q PlaceClassesQ) UpdateOne(ctx context.Context) (PlaceClass, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.
		Suffix("RETURNING " + PlacesClassesColumns).
		ToSql()
	if err != nil {
		return PlaceClass{}, fmt.Errorf("building update query for %s: %w", PlaceClassesTable, err)
	}

	var updated PlaceClass
	if err := updated.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return PlaceClass{}, err
	}

	return updated, nil
}

func (q PlaceClassesQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", PlaceClassesTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", PlaceClassesTable, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", PlaceClassesTable, err)
	}

	return affected, nil
}

func (q PlaceClassesQ) UpdateParent(parentID *uuid.UUID) PlaceClassesQ {
	if parentID == nil {
		q.updater = q.updater.Set("parent_id", nil)
	} else {
		q.updater = q.updater.Set("parent_id", *parentID)
	}
	return q
}

func (q PlaceClassesQ) UpdateCode(code string) PlaceClassesQ {
	q.updater = q.updater.Set("code", code)
	return q
}

func (q PlaceClassesQ) UpdateName(name string) PlaceClassesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q PlaceClassesQ) UpdateDescription(description string) PlaceClassesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q PlaceClassesQ) UpdateIcon(icon *string) PlaceClassesQ {
	if icon == nil {
		q.updater = q.updater.Set("icon", nil)
	} else {
		q.updater = q.updater.Set("icon", *icon)
	}
	return q
}

func (q PlaceClassesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlaceClassesTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PlaceClassesTable, err)
	}

	return nil
}
