package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/netbill/pgx"
)

const PlaceTimetablesTable = "place_timetables"

const PlaceTimetablesColumns = "id, place_id, start_min, end_min"

type PlaceTimetableRow struct {
	ID       uuid.UUID `json:"id"`
	PlaceID  uuid.UUID `json:"place_id"`
	StartMin int       `json:"start_min"`
	EndMin   int       `json:"end_min"`
}

func (t *PlaceTimetableRow) scan(row sq.RowScanner) error {
	if err := row.Scan(&t.ID, &t.PlaceID, &t.StartMin, &t.EndMin); err != nil {
		return fmt.Errorf("scanning %s: %w", PlaceTimetablesTable, err)
	}
	return nil
}

type PlaceTimetablesQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceTimetablesQ(db pgx.DBTX) PlaceTimetablesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return PlaceTimetablesQ{
		db:       db,
		selector: b.Select(PlaceTimetablesColumns).From(PlaceTimetablesTable),
		inserter: b.Insert(PlaceTimetablesTable),
		updater:  b.Update(PlaceTimetablesTable),
		deleter:  b.Delete(PlaceTimetablesTable),
		counter:  b.Select("COUNT(*) AS count").From(PlaceTimetablesTable),
	}
}

func (q PlaceTimetablesQ) New() PlaceTimetablesQ { return NewPlaceTimetablesQ(q.db) }

type PlaceTimetablesQInsertInput struct {
	ID       uuid.UUID
	PlaceID  uuid.UUID
	StartMin int
	EndMin   int
}

func (q PlaceTimetablesQ) Insert(ctx context.Context, in ...PlaceTimetablesQInsertInput) error {
	if len(in) == 0 {
		return nil
	}

	ins := q.inserter.Columns("id", "place_id", "start_min", "end_min")
	for _, r := range in {
		ins = ins.Values(r.ID, r.PlaceID, r.StartMin, r.EndMin)
	}

	query, args, err := ins.ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", PlaceTimetablesTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing insert query for %s: %w", PlaceTimetablesTable, err)
	}

	return nil
}

func (q PlaceTimetablesQ) Upsert(ctx context.Context, in ...PlaceTimetablesQInsertInput) error {
	if len(in) == 0 {
		return nil
	}

	const cols = "(id, place_id, start_min, end_min)"

	var (
		args []any
		ph   []string
		i    = 1
	)

	for _, r := range in {
		ph = append(ph, fmt.Sprintf("($%d,$%d,$%d,$%d)", i, i+1, i+2, i+3))
		args = append(args, r.ID, r.PlaceID, r.StartMin, r.EndMin)
		i += 4
	}

	query := fmt.Sprintf(`
		INSERT INTO %s %s VALUES %s
		ON CONFLICT (id) DO UPDATE
		SET place_id = EXCLUDED.place_id,
		    start_min = EXCLUDED.start_min,
		    end_min   = EXCLUDED.end_min
	`, PlaceTimetablesTable, cols, strings.Join(ph, ","))

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing upsert query for %s: %w", PlaceTimetablesTable, err)
	}

	return nil
}

func (q PlaceTimetablesQ) Get(ctx context.Context) (PlaceTimetableRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceTimetableRow{}, fmt.Errorf("building select query for %s: %w", PlaceTimetablesTable, err)
	}

	var out PlaceTimetableRow
	if err := out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return PlaceTimetableRow{}, err
	}

	return out, nil
}

func (q PlaceTimetablesQ) Select(ctx context.Context) ([]PlaceTimetableRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlaceTimetablesTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PlaceTimetablesTable, err)
	}
	defer rows.Close()

	var out []PlaceTimetableRow
	for rows.Next() {
		var t PlaceTimetableRow
		if err := t.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q PlaceTimetablesQ) UpdateOne(ctx context.Context) (PlaceTimetableRow, error) {
	query, args, err := q.updater.Suffix("RETURNING " + PlaceTimetablesColumns).ToSql()
	if err != nil {
		return PlaceTimetableRow{}, fmt.Errorf("building update query for %s: %w", PlaceTimetablesTable, err)
	}

	var updated PlaceTimetableRow
	if err := updated.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return PlaceTimetableRow{}, err
	}

	return updated, nil
}

func (q PlaceTimetablesQ) UpdateMany(ctx context.Context) (int64, error) {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", PlaceTimetablesTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", PlaceTimetablesTable, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", PlaceTimetablesTable, err)
	}

	return affected, nil
}

func (q PlaceTimetablesQ) UpdatePlaceID(placeID uuid.UUID) PlaceTimetablesQ {
	q.updater = q.updater.Set("place_id", placeID)
	return q
}

func (q PlaceTimetablesQ) UpdateStartMin(v int) PlaceTimetablesQ {
	q.updater = q.updater.Set("start_min", v)
	return q
}

func (q PlaceTimetablesQ) UpdateEndMin(v int) PlaceTimetablesQ {
	q.updater = q.updater.Set("end_min", v)
	return q
}

func (q PlaceTimetablesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlaceTimetablesTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PlaceTimetablesTable, err)
	}

	return nil
}

func (q PlaceTimetablesQ) FilterByID(id uuid.UUID) PlaceTimetablesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	return q
}

func (q PlaceTimetablesQ) FilterByPlaceID(placeID uuid.UUID) PlaceTimetablesQ {
	q.selector = q.selector.Where(sq.Eq{"place_id": placeID})
	q.updater = q.updater.Where(sq.Eq{"place_id": placeID})
	q.deleter = q.deleter.Where(sq.Eq{"place_id": placeID})
	q.counter = q.counter.Where(sq.Eq{"place_id": placeID})
	return q
}

func (q PlaceTimetablesQ) FilterBetween(start, end int) PlaceTimetablesQ {
	const week = 7 * 24 * 60

	norm := func(x int) int {
		x %= week
		if x < 0 {
			x += week
		}
		return x
	}

	s, e := norm(start), norm(end)
	if s == e {
		q.selector = q.selector.Where("1=0")
		q.updater = q.updater.Where("1=0")
		q.deleter = q.deleter.Where("1=0")
		q.counter = q.counter.Where("1=0")
		return q
	}

	var cond any
	if s < e {
		cond = sq.And{
			sq.Lt{"start_min": e},
			sq.Gt{"end_min": s},
		}
	} else {
		cond = sq.Or{
			sq.Gt{"end_min": s},
			sq.Lt{"start_min": e},
		}
	}

	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}

func (q PlaceTimetablesQ) Page(limit, offset uint64) PlaceTimetablesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q PlaceTimetablesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlaceTimetablesTable, err)
	}

	var cnt uint64
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&cnt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("scanning count for %s: %w", PlaceTimetablesTable, err)
	}

	return cnt, nil
}
