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

const PlacePossibilitiesTable = "place_possibilities"
const PlacePossibilitiesColumns = "code, description"

type PlacePossibility struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (f *PlacePossibility) scan(row sq.RowScanner) error {
	if err := row.Scan(
		&f.Code,
		&f.Description,
	); err != nil {
		return fmt.Errorf("scanning place possibility: %w", err)
	}
	return nil
}

type PlacePossibilitiesQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlacePossibilitiesQ(db pgx.DBTX) PlacePossibilitiesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return PlacePossibilitiesQ{
		db:       db,
		selector: builder.Select("pf." + strings.ReplaceAll(PlacePossibilitiesColumns, ", ", ", pf.")).From(PlacePossibilitiesTable + " pf"),
		inserter: builder.Insert(PlacePossibilitiesTable),
		updater:  builder.Update(PlacePossibilitiesTable),
		deleter:  builder.Delete(PlacePossibilitiesTable),
		counter:  builder.Select("COUNT(*) AS count").From(PlacePossibilitiesTable + " pf"),
	}
}

type PlacePossibilitiesQInsertInput struct {
	Code        string
	Description string
}

func (q PlacePossibilitiesQ) Insert(ctx context.Context, data PlacePossibilitiesQInsertInput) (PlacePossibility, error) {
	query, args, err := q.inserter.
		SetMap(map[string]interface{}{
			"code":        data.Code,
			"description": data.Description,
		}).
		Suffix("RETURNING " + PlacePossibilitiesColumns).
		ToSql()
	if err != nil {
		return PlacePossibility{}, fmt.Errorf("building insert query for %s: %w", PlacePossibilitiesTable, err)
	}

	var inserted PlacePossibility
	if err := inserted.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return PlacePossibility{}, nil
		default:
			return PlacePossibility{}, err
		}
	}

	return inserted, nil
}

func (q PlacePossibilitiesQ) FilterByCode(code string) PlacePossibilitiesQ {
	q.selector = q.selector.Where(sq.Eq{"code": code})
	q.counter = q.counter.Where(sq.Eq{"code": code})
	q.updater = q.updater.Where(sq.Eq{"code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	return q
}

func (q PlacePossibilitiesQ) FilterByPlaceID(placeID uuid.UUID) PlacePossibilitiesQ {
	q.selector = q.selector.
		Join(PlacePossibilityLinksTable + " pfl ON pfl.possibility_id = pf.id").
		Where(sq.Eq{"pfl.place_id": placeID})

	q.counter = q.counter.
		Join(PlacePossibilityLinksTable + " pfl ON pfl.possibility_id = pf.id").
		Where(sq.Eq{"pfl.place_id": placeID})

	return q
}

func (q PlacePossibilitiesQ) Get(ctx context.Context) (PlacePossibility, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlacePossibility{}, fmt.Errorf("building select query for %s: %w", PlacePossibilitiesTable, err)
	}

	var f PlacePossibility
	if err := f.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return PlacePossibility{}, err
	}

	return f, nil
}

func (q PlacePossibilitiesQ) Select(ctx context.Context) ([]PlacePossibility, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlacePossibilitiesTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PlacePossibilitiesTable, err)
	}
	defer rows.Close()

	var items []PlacePossibility
	for rows.Next() {
		var f PlacePossibility
		if err := f.scan(rows); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (q PlacePossibilitiesQ) UpdateOne(ctx context.Context) (PlacePossibility, error) {
	query, args, err := q.updater.
		Suffix("RETURNING " + PlacePossibilitiesColumns).
		ToSql()
	if err != nil {
		return PlacePossibility{}, fmt.Errorf("building update query for %s: %w", PlacePossibilitiesTable, err)
	}

	var updated PlacePossibility
	if err := updated.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return PlacePossibility{}, err
	}

	return updated, nil
}

func (q PlacePossibilitiesQ) UpdateMany(ctx context.Context) (int64, error) {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", PlacePossibilitiesTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", PlacePossibilitiesTable, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", PlacePossibilitiesTable, err)
	}

	return affected, nil
}

func (q PlacePossibilitiesQ) UpdateCode(code string) PlacePossibilitiesQ {
	q.updater = q.updater.Set("code", code)
	return q
}

func (q PlacePossibilitiesQ) UpdateDescription(description string) PlacePossibilitiesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q PlacePossibilitiesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlacePossibilitiesTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PlacePossibilitiesTable, err)
	}

	return nil
}

func (q PlacePossibilitiesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlacePossibilitiesTable, err)
	}

	var count uint
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PlacePossibilitiesTable, err)
	}

	return count, nil
}

func (q PlacePossibilitiesQ) Exists(ctx context.Context) (bool, error) {
	query, args, err := q.selector.
		Columns("1").
		Limit(1).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", PlacePossibilitiesTable, err)
	}

	var one int
	err = q.db.QueryRowContext(ctx, query, args...).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
