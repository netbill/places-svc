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

const placeClassesTable = "place_classes"
const placeClassesColumns = "id, parent_id, name, description, icon_key, version, created_at, updated_at, deprecated_at"
const placeClassesColumnsP = "pc.id, pc.parent_id, pc.name, pc.description, pc.icon_key, pc.version, pc.created_at, pc.updated_at, pc.deprecated_at"

func scanPlaceClass(row sq.RowScanner) (pc repository.PlaceClassRow, err error) {
	var parentID pgtype.UUID
	var iconKey pgtype.Text
	var deprecatedAt pgtype.Timestamptz

	err = row.Scan(
		&pc.ID,
		&parentID,
		&pc.Name,
		&pc.Description,
		&iconKey,
		&pc.Version,
		&pc.CreatedAt,
		&pc.UpdatedAt,
		&deprecatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.PlaceClassRow{}, nil
	case err != nil:
		return repository.PlaceClassRow{}, fmt.Errorf("scanning place class: %w", err)
	}

	if parentID.Valid {
		v := uuid.UUID(parentID.Bytes)
		pc.ParentID = &v
	}
	if iconKey.Valid {
		pc.IconKey = &iconKey.String
	}
	if deprecatedAt.Valid {
		t := deprecatedAt.Time
		pc.DeprecatedAt = &t
	}

	return pc, nil
}

type placeClasses struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceClassesQ(db *pgdbx.DB) repository.PlaceClassesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &placeClasses{
		db:       db,
		selector: b.Select(placeClassesColumnsP).From(placeClassesTable + " pc"),
		inserter: b.Insert(placeClassesTable),
		updater:  b.Update(placeClassesTable + " pc"),
		deleter:  b.Delete(placeClassesTable + " pc"),
		counter:  b.Select("COUNT(*)").From(placeClassesTable + " pc"),
	}
}

func (q *placeClasses) New() repository.PlaceClassesQ {
	return NewPlaceClassesQ(q.db)
}

func (q *placeClasses) Insert(ctx context.Context, data repository.PlaceClassRow) (repository.PlaceClassRow, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"parent_id":   data.ParentID,
		"name":        data.Name,
		"description": data.Description,
		"icon_key":    data.IconKey,
	}).Suffix("RETURNING " + placeClassesColumns).ToSql()
	if err != nil {
		return repository.PlaceClassRow{}, fmt.Errorf("building insert query for %s: %w", placeClassesTable, err)
	}

	return scanPlaceClass(q.db.QueryRow(ctx, query, args...))
}

func (q *placeClasses) FilterByID(id ...uuid.UUID) repository.PlaceClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.id": id})
	q.counter = q.counter.Where(sq.Eq{"pc.id": id})
	q.updater = q.updater.Where(sq.Eq{"pc.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"pc.id": id})
	return q
}

func (q *placeClasses) FilterByParentID(parentID ...uuid.UUID) repository.PlaceClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.parent_id": parentID})
	q.counter = q.counter.Where(sq.Eq{"pc.parent_id": parentID})
	q.updater = q.updater.Where(sq.Eq{"pc.parent_id": parentID})
	q.deleter = q.deleter.Where(sq.Eq{"pc.parent_id": parentID})
	return q
}

func (q *placeClasses) FilterBestMatch(text string) repository.PlaceClassesQ {
	if text == "" {
		return q
	}
	pattern := "%" + text + "%"
	expr := sq.Or{
		sq.Expr("pc.name ILIKE ?", pattern),
		sq.Expr("pc.description ILIKE ?", pattern),
	}
	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	return q
}

func (q *placeClasses) FilterByClassID(classID uuid.UUID, includeChildren, includeParents bool) repository.PlaceClassesQ {
	if !includeChildren && !includeParents {
		q.selector = q.selector.Where(sq.Eq{"pc.id": classID})
		q.counter = q.counter.Where(sq.Eq{"pc.id": classID})
		return q
	}

	subSQL := `
		WITH RECURSIVE
		anc AS (
			SELECT id, parent_id FROM place_classes WHERE id = ?
			UNION ALL
			SELECT c.id, c.parent_id FROM place_classes c
			JOIN anc ON anc.parent_id = c.id
			WHERE ? = TRUE
		),
		des AS (
			SELECT id FROM place_classes WHERE id = ?
			UNION ALL
			SELECT c.id FROM place_classes c
			JOIN des ON c.parent_id = des.id
			WHERE ? = TRUE
		)
		SELECT DISTINCT id FROM (SELECT id FROM anc UNION SELECT id FROM des) t
	`

	expr := sq.Expr(
		"pc.id IN ("+subSQL+")",
		classID, includeParents,
		classID, includeChildren,
	)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	return q
}

func (q *placeClasses) OrderName(asc bool) repository.PlaceClassesQ {
	if asc {
		q.selector = q.selector.OrderBy("pc.name ASC", "pc.id ASC")
	} else {
		q.selector = q.selector.OrderBy("pc.name DESC", "pc.id DESC")
	}
	return q
}

func (q *placeClasses) OrderRoot(asc bool) repository.PlaceClassesQ {
	if asc {
		q.selector = q.selector.OrderBy("pc.parent_id ASC NULLS FIRST", "pc.name ASC")
	} else {
		q.selector = q.selector.OrderBy("pc.parent_id DESC NULLS LAST", "pc.name DESC")
	}
	return q
}

func (q *placeClasses) Get(ctx context.Context) (repository.PlaceClassRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.PlaceClassRow{}, fmt.Errorf("building select query for %s: %w", placeClassesTable, err)
	}

	return scanPlaceClass(q.db.QueryRow(ctx, query, args...))
}

func (q *placeClasses) Select(ctx context.Context) ([]repository.PlaceClassRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", placeClassesTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", placeClassesTable, err)
	}
	defer rows.Close()

	out := make([]repository.PlaceClassRow, 0)
	for rows.Next() {
		pc, err := scanPlaceClass(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, pc)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *placeClasses) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", placeClassesTable, err)
	}

	var ok bool
	if err = q.db.QueryRow(ctx, "SELECT EXISTS ("+subSQL+")", subArgs...).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning exists for %s: %w", placeClassesTable, err)
	}

	return ok, nil
}

func (q *placeClasses) Page(limit, offset uint) repository.PlaceClassesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q *placeClasses) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", placeClassesTable, err)
	}

	var count uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", placeClassesTable, err)
	}

	return count, nil
}

func (q *placeClasses) UpdateOne(ctx context.Context) (repository.PlaceClassRow, error) {
	q.updater = q.updater.
		Set("updated_at", time.Now().UTC()).
		Set("version", sq.Expr("version + 1"))

	query, args, err := q.updater.Suffix("RETURNING " + placeClassesColumns).ToSql()
	if err != nil {
		return repository.PlaceClassRow{}, fmt.Errorf("building update query for %s: %w", placeClassesTable, err)
	}

	return scanPlaceClass(q.db.QueryRow(ctx, query, args...))
}

func (q *placeClasses) UpdateParent(parentID *uuid.UUID) repository.PlaceClassesQ {
	q.updater = q.updater.Set("parent_id", parentID)
	return q
}

func (q *placeClasses) UpdateName(name string) repository.PlaceClassesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q *placeClasses) UpdateDescription(description string) repository.PlaceClassesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q *placeClasses) UpdateIconKey(key *string) repository.PlaceClassesQ {
	q.updater = q.updater.Set("icon_key", key)
	return q
}

func (q *placeClasses) UpdateDeprecatedAt(t *time.Time) repository.PlaceClassesQ {
	q.updater = q.updater.Set("deprecated_at", t)
	return q
}

func (q *placeClasses) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", placeClassesTable, err)
	}

	if _, err := q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", placeClassesTable, err)
	}

	return nil
}
