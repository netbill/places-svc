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

const placeClassesColumns = "id, parent_id, code, name, description, icon, created_at, updated_at"
const placeClassesColumnsP = "pc.id, pc.parent_id, pc.code, pc.name, pc.description, pc.icon, pc.created_at, pc.updated_at"

func scanPlaceClass(row sq.RowScanner) (pc repository.PlaceClassRow, err error) {
	var parentID pgtype.UUID
	var icon pgtype.Text

	err = row.Scan(
		&pc.ID,
		&parentID,
		&pc.Code,
		&pc.Name,
		&pc.Description,
		&icon,
		&pc.CreatedAt,
		&pc.UpdatedAt,
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
	if icon.Valid {
		pc.Icon = &icon.String
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
	now := time.Now().UTC()

	if data.ID == uuid.Nil {
		return repository.PlaceClassRow{}, fmt.Errorf("missing id")
	}
	if data.Code == "" {
		return repository.PlaceClassRow{}, fmt.Errorf("missing code")
	}
	if data.Name == "" {
		return repository.PlaceClassRow{}, fmt.Errorf("missing name")
	}

	query, args, err := q.inserter.SetMap(map[string]any{
		"id":          data.ID,
		"parent_id":   data.ParentID,
		"code":        data.Code,
		"name":        data.Name,
		"description": data.Description,
		"icon":        data.Icon,
		"created_at":  now,
		"updated_at":  now,
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

func (q *placeClasses) FilterRoots() repository.PlaceClassesQ {
	q.selector = q.selector.Where(sq.Expr("pc.parent_id IS NULL"))
	q.counter = q.counter.Where(sq.Expr("pc.parent_id IS NULL"))
	q.updater = q.updater.Where(sq.Expr("pc.parent_id IS NULL"))
	q.deleter = q.deleter.Where(sq.Expr("pc.parent_id IS NULL"))
	return q
}

func (q *placeClasses) FilterByCode(code string) repository.PlaceClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.code": code})
	q.counter = q.counter.Where(sq.Eq{"pc.code": code})
	q.updater = q.updater.Where(sq.Eq{"pc.code": code})
	q.deleter = q.deleter.Where(sq.Eq{"pc.code": code})
	return q
}

func (q *placeClasses) FilterNameLike(name string) repository.PlaceClassesQ {
	q.selector = q.selector.Where(sq.ILike{"pc.name": "%" + name + "%"})
	q.counter = q.counter.Where(sq.ILike{"pc.name": "%" + name + "%"})
	return q
}

func (q *placeClasses) FilterDescriptionLike(description string) repository.PlaceClassesQ {
	q.selector = q.selector.Where(sq.ILike{"pc.description": "%" + description + "%"})
	q.counter = q.counter.Where(sq.ILike{"pc.description": "%" + description + "%"})
	return q
}

func (q *placeClasses) FilterByText(text string) repository.PlaceClassesQ {
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

func (q *placeClasses) FilterByClassID(classID uuid.UUID, includeChild, includeParent bool) repository.PlaceClassesQ {
	if !includeChild && !includeParent {
		q.selector = q.selector.Where(sq.Eq{"pc.id": classID})
		q.counter = q.counter.Where(sq.Eq{"pc.id": classID})
		return q
	}

	subSQL := `
		WITH RECURSIVE
		anc AS (
			SELECT id, parent_id
			FROM place_classes
			WHERE id = ?
			UNION ALL
			SELECT pc2.id, pc2.parent_id
			FROM place_classes pc2
			JOIN anc ON anc.parent_id = pc2.id
			WHERE (? = TRUE)
		),
		des AS (
			SELECT id
			FROM place_classes
			WHERE id = ?
			UNION ALL
			SELECT pc2.id
			FROM place_classes pc2
			JOIN des ON pc2.parent_id = des.id
			WHERE (? = TRUE)
		)
		SELECT DISTINCT id
		FROM (
			SELECT id FROM anc
			UNION
			SELECT id FROM des
		) t
	`

	expr := sq.Expr(
		"pc.id IN ("+subSQL+")",
		classID, includeParent,
		classID, includeChild,
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

	sqlq := "SELECT EXISTS (" + subSQL + ")"

	var ok bool
	if err = q.db.QueryRow(ctx, sqlq, subArgs...).Scan(&ok); err != nil {
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
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.Suffix("RETURNING " + placeClassesColumns).ToSql()
	if err != nil {
		return repository.PlaceClassRow{}, fmt.Errorf("building update query for %s: %w", placeClassesTable, err)
	}

	return scanPlaceClass(q.db.QueryRow(ctx, query, args...))
}

func (q *placeClasses) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", placeClassesTable, err)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", placeClassesTable, err)
	}

	return res.RowsAffected(), nil
}

func (q *placeClasses) UpdateParent(parentID *uuid.UUID) repository.PlaceClassesQ {
	q.updater = q.updater.Set("parent_id", parentID)
	return q
}

func (q *placeClasses) UpdateCode(code string) repository.PlaceClassesQ {
	q.updater = q.updater.Set("code", code)
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

func (q *placeClasses) UpdateIcon(icon *string) repository.PlaceClassesQ {
	q.updater = q.updater.Set("icon", icon)
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
