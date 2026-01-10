package pgdb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/netbill/pgx"
)

const PlacePossibilityLinksTable = "place_possibility_links"

type PlacePossibilityLink struct {
	PlaceID       uuid.UUID `json:"place_id"`
	PossibilityID uuid.UUID `json:"possibility_id"`
}

type PlacePossibilityLinksQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlacePossibilityLinksQ(db pgx.DBTX) PlacePossibilityLinksQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return PlacePossibilityLinksQ{
		db:       db,
		selector: builder.Select("place_id, possibility_id").From(PlacePossibilityLinksTable),
		inserter: builder.Insert(PlacePossibilityLinksTable),
		deleter:  builder.Delete(PlacePossibilityLinksTable),
		counter:  builder.Select("COUNT(*) AS count").From(PlacePossibilityLinksTable),
	}
}

type PlacePossibilityLinksQInsertInput struct {
	PlaceID       uuid.UUID
	PossibilityID uuid.UUID
}

func (q PlacePossibilityLinksQ) Insert(ctx context.Context, data PlacePossibilityLinksQInsertInput) error {
	query, args, err := q.inserter.
		SetMap(map[string]interface{}{
			"place_id":       data.PlaceID,
			"possibility_id": data.PossibilityID,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", PlacePossibilityLinksTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing insert query for %s: %w", PlacePossibilityLinksTable, err)
	}

	return nil
}

func (q PlacePossibilityLinksQ) FilterByPlaceID(placeID uuid.UUID) PlacePossibilityLinksQ {
	q.selector = q.selector.Where(sq.Eq{"place_id": placeID})
	q.counter = q.counter.Where(sq.Eq{"place_id": placeID})
	q.deleter = q.deleter.Where(sq.Eq{"place_id": placeID})
	return q
}

func (q PlacePossibilityLinksQ) FilterByPossibilityID(possibilityID uuid.UUID) PlacePossibilityLinksQ {
	q.selector = q.selector.Where(sq.Eq{"possibility_id": possibilityID})
	q.counter = q.counter.Where(sq.Eq{"possibility_id": possibilityID})
	q.deleter = q.deleter.Where(sq.Eq{"possibility_id": possibilityID})
	return q
}

func (q PlacePossibilityLinksQ) Select(ctx context.Context) ([]PlacePossibilityLink, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlacePossibilityLinksTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PlacePossibilityLinksTable, err)
	}
	defer rows.Close()

	var links []PlacePossibilityLink
	for rows.Next() {
		var l PlacePossibilityLink
		if err := rows.Scan(&l.PlaceID, &l.PossibilityID); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (q PlacePossibilityLinksQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlacePossibilityLinksTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PlacePossibilityLinksTable, err)
	}

	return nil
}

func (q PlacePossibilityLinksQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlacePossibilityLinksTable, err)
	}

	var count uint
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PlacePossibilityLinksTable, err)
	}

	return count, nil
}
