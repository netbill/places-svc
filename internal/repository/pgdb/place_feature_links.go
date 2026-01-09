package pgdb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/netbill/pgx"
)

const PlaceFeatureLinksTable = "place_feature_links"

type PlaceFeatureLink struct {
	PlaceID   uuid.UUID `json:"place_id"`
	FeatureID uuid.UUID `json:"feature_id"`
}

type PlaceFeatureLinksQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceFeatureLinksQ(db pgx.DBTX) PlaceFeatureLinksQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return PlaceFeatureLinksQ{
		db:       db,
		selector: builder.Select("place_id, feature_id").From(PlaceFeatureLinksTable),
		inserter: builder.Insert(PlaceFeatureLinksTable),
		deleter:  builder.Delete(PlaceFeatureLinksTable),
		counter:  builder.Select("COUNT(*) AS count").From(PlaceFeatureLinksTable),
	}
}

type PlaceFeatureLinksQInsertInput struct {
	PlaceID   uuid.UUID
	FeatureID uuid.UUID
}

func (q PlaceFeatureLinksQ) Insert(ctx context.Context, data PlaceFeatureLinksQInsertInput) error {
	query, args, err := q.inserter.
		SetMap(map[string]interface{}{
			"place_id":   data.PlaceID,
			"feature_id": data.FeatureID,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", PlaceFeatureLinksTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing insert query for %s: %w", PlaceFeatureLinksTable, err)
	}

	return nil
}

func (q PlaceFeatureLinksQ) FilterByPlaceID(placeID uuid.UUID) PlaceFeatureLinksQ {
	q.selector = q.selector.Where(sq.Eq{"place_id": placeID})
	q.counter = q.counter.Where(sq.Eq{"place_id": placeID})
	q.deleter = q.deleter.Where(sq.Eq{"place_id": placeID})
	return q
}

func (q PlaceFeatureLinksQ) FilterByFeatureID(featureID uuid.UUID) PlaceFeatureLinksQ {
	q.selector = q.selector.Where(sq.Eq{"feature_id": featureID})
	q.counter = q.counter.Where(sq.Eq{"feature_id": featureID})
	q.deleter = q.deleter.Where(sq.Eq{"feature_id": featureID})
	return q
}

func (q PlaceFeatureLinksQ) Select(ctx context.Context) ([]PlaceFeatureLink, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlaceFeatureLinksTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PlaceFeatureLinksTable, err)
	}
	defer rows.Close()

	var links []PlaceFeatureLink
	for rows.Next() {
		var l PlaceFeatureLink
		if err := rows.Scan(&l.PlaceID, &l.FeatureID); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (q PlaceFeatureLinksQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlaceFeatureLinksTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PlaceFeatureLinksTable, err)
	}

	return nil
}

func (q PlaceFeatureLinksQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlaceFeatureLinksTable, err)
	}

	var count uint
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PlaceFeatureLinksTable, err)
	}

	return count, nil
}
