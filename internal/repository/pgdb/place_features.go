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
	replicaspg "github.com/netbill/replicas/pgdb"
)

const PlaceFeaturesTable = "place_features"
const PlaceFeaturesColumns = "id, code, description"

type PlaceFeature struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}

func (f *PlaceFeature) scan(row sq.RowScanner) error {
	if err := row.Scan(
		&f.ID,
		&f.Code,
		&f.Description,
	); err != nil {
		return fmt.Errorf("scanning place feature: %w", err)
	}
	return nil
}

type PlaceFeaturesQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceFeaturesQ(db pgx.DBTX) PlaceFeaturesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return PlaceFeaturesQ{
		db:       db,
		selector: builder.Select("pf." + strings.ReplaceAll(PlaceFeaturesColumns, ", ", ", pf.")).From(PlaceFeaturesTable + " pf"),
		inserter: builder.Insert(PlaceFeaturesTable),
		updater:  builder.Update(PlaceFeaturesTable),
		deleter:  builder.Delete(PlaceFeaturesTable),
		counter:  builder.Select("COUNT(*) AS count").From(PlaceFeaturesTable + " pf"),
	}
}

type PlaceFeaturesQInsertInput struct {
	Code        string
	Description string
}

func (q PlaceFeaturesQ) Insert(ctx context.Context, data PlaceFeaturesQInsertInput) (PlaceFeature, error) {
	query, args, err := q.inserter.
		SetMap(map[string]interface{}{
			"code":        data.Code,
			"description": data.Description,
		}).
		Suffix("RETURNING " + PlaceFeaturesColumns).
		ToSql()
	if err != nil {
		return PlaceFeature{}, fmt.Errorf("building insert query for %s: %w", PlaceFeaturesTable, err)
	}

	var inserted PlaceFeature
	if err := inserted.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return PlaceFeature{}, nil
		default:
			return PlaceFeature{}, err
		}
	}

	return inserted, nil
}

func (q PlaceFeaturesQ) FilterByID(id uuid.UUID) PlaceFeaturesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q PlaceFeaturesQ) FilterByCode(code string) PlaceFeaturesQ {
	q.selector = q.selector.Where(sq.Eq{"code": code})
	q.counter = q.counter.Where(sq.Eq{"code": code})
	q.updater = q.updater.Where(sq.Eq{"code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	return q
}

func (q PlaceFeaturesQ) FilterByPlaceID(placeID uuid.UUID) PlaceFeaturesQ {
	q.selector = q.selector.
		Join(PlaceFeatureLinksTable + " pfl ON pfl.feature_id = pf.id").
		Where(sq.Eq{"pfl.place_id": placeID})

	q.counter = q.counter.
		Join(PlaceFeatureLinksTable + " pfl ON pfl.feature_id = pf.id").
		Where(sq.Eq{"pfl.place_id": placeID})

	return q
}

func (q PlaceFeaturesQ) FilterByAccountID(accountID uuid.UUID) PlaceFeaturesQ {
	q.selector = q.selector.
		Join(PlaceFeatureLinksTable + " pfl ON pfl.feature_id = pf.id").
		Join(PlacesTable + " p ON p.id = pfl.place_id").
		Join(replicaspg.OrganizationMembersTable + " om ON om.organization_id = p.organization_id").
		Where(sq.Eq{"om.account_id": accountID})

	q.counter = q.counter.
		Join(PlaceFeatureLinksTable + " pfl ON pfl.feature_id = pf.id").
		Join(PlacesTable + " p ON p.id = pfl.place_id").
		Join(replicaspg.OrganizationMembersTable + " om ON om.organization_id = p.organization_id").
		Where(sq.Eq{"om.account_id": accountID})

	return q
}

func (q PlaceFeaturesQ) Get(ctx context.Context) (PlaceFeature, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceFeature{}, fmt.Errorf("building select query for %s: %w", PlaceFeaturesTable, err)
	}

	var f PlaceFeature
	if err := f.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return PlaceFeature{}, err
	}

	return f, nil
}

func (q PlaceFeaturesQ) Select(ctx context.Context) ([]PlaceFeature, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlaceFeaturesTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PlaceFeaturesTable, err)
	}
	defer rows.Close()

	var items []PlaceFeature
	for rows.Next() {
		var f PlaceFeature
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

func (q PlaceFeaturesQ) UpdateOne(ctx context.Context) (PlaceFeature, error) {
	query, args, err := q.updater.
		Suffix("RETURNING " + PlaceFeaturesColumns).
		ToSql()
	if err != nil {
		return PlaceFeature{}, fmt.Errorf("building update query for %s: %w", PlaceFeaturesTable, err)
	}

	var updated PlaceFeature
	if err := updated.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return PlaceFeature{}, err
	}

	return updated, nil
}

func (q PlaceFeaturesQ) UpdateMany(ctx context.Context) (int64, error) {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", PlaceFeaturesTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", PlaceFeaturesTable, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", PlaceFeaturesTable, err)
	}

	return affected, nil
}

func (q PlaceFeaturesQ) UpdateCode(code string) PlaceFeaturesQ {
	q.updater = q.updater.Set("code", code)
	return q
}

func (q PlaceFeaturesQ) UpdateDescription(description string) PlaceFeaturesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q PlaceFeaturesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlaceFeaturesTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PlaceFeaturesTable, err)
	}

	return nil
}

func (q PlaceFeaturesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlaceFeaturesTable, err)
	}

	var count uint
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PlaceFeaturesTable, err)
	}

	return count, nil
}
