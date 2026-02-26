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
	"github.com/paulmach/orb"
)

const placesTable = "places"
const placesColumns = `id, class_id, organization_id, status, verified, point, address, name,
	description, icon_key, banner_key, website, phone, version, created_at, updated_at`

func scanPlaceRow(row sq.RowScanner) (r repository.PlaceRow, err error) {
	var description pgtype.Text
	var iconKey pgtype.Text
	var bannerKey pgtype.Text
	var website pgtype.Text
	var phone pgtype.Text

	err = row.Scan(
		&r.ID,
		&r.ClassID,
		&r.OrganizationID,
		&r.Status,
		&r.Verified,
		&r.Point,
		&r.Address,
		&r.Name,
		&description,
		&iconKey,
		&bannerKey,
		&website,
		&phone,
		&r.Version,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.PlaceRow{}, nil
	case err != nil:
		return repository.PlaceRow{}, fmt.Errorf("scanning place row: %w", err)
	}

	if description.Valid {
		r.Description = &description.String
	}
	if iconKey.Valid {
		r.IconKey = &iconKey.String
	}
	if bannerKey.Valid {
		r.BannerKey = &bannerKey.String
	}
	if website.Valid {
		r.Website = &website.String
	}
	if phone.Valid {
		r.Phone = &phone.String
	}

	return r, nil
}

type placesQ struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlacesQ(db *pgdbx.DB) repository.PlacesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &placesQ{
		db:       db,
		selector: b.Select(placesColumns).From(placesTable),
		inserter: b.Insert(placesTable),
		updater:  b.Update(placesTable),
		deleter:  b.Delete(placesTable),
		counter:  b.Select("COUNT(*)").From(placesTable),
	}
}

func (q *placesQ) New() repository.PlacesQ {
	return NewPlacesQ(q.db)
}

func (q *placesQ) Insert(ctx context.Context, data repository.PlaceRow) (repository.PlaceRow, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"class_id":        data.ClassID,
		"organization_id": data.OrganizationID,
		"status":          data.Status,
		"verified":        data.Verified,
		"point":           data.Point,
		"address":         data.Address,
		"name":            data.Name,
		"description":     data.Description,
		"icon_key":        data.IconKey,
		"banner_key":      data.BannerKey,
		"website":         data.Website,
		"phone":           data.Phone,
	}).Suffix("RETURNING " + placesColumns).ToSql()
	if err != nil {
		return repository.PlaceRow{}, fmt.Errorf("building insert query for %s: %w", placesTable, err)
	}

	return scanPlaceRow(q.db.QueryRow(ctx, query, args...))
}

func (q *placesQ) FilterByID(id ...uuid.UUID) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q *placesQ) FilterByOrganizationID(orgID *uuid.UUID) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"organization_id": orgID})
	q.counter = q.counter.Where(sq.Eq{"organization_id": orgID})
	q.updater = q.updater.Where(sq.Eq{"organization_id": orgID})
	q.deleter = q.deleter.Where(sq.Eq{"organization_id": orgID})
	return q
}

func (q *placesQ) FilterByClassID(includeChildren, includeParents bool, classIDs ...uuid.UUID) repository.PlacesQ {
	if len(classIDs) == 0 {
		return q
	}

	if !includeChildren && !includeParents {
		q.selector = q.selector.Where(sq.Eq{"class_id": classIDs})
		q.counter = q.counter.Where(sq.Eq{"class_id": classIDs})
		q.updater = q.updater.Where(sq.Eq{"class_id": classIDs})
		q.deleter = q.deleter.Where(sq.Eq{"class_id": classIDs})
		return q
	}

	subSQL := `
		WITH RECURSIVE
		anc AS (
			SELECT id, parent_id FROM place_classes WHERE id = ANY(?)
			UNION ALL
			SELECT pc.id, pc.parent_id FROM place_classes pc
			JOIN anc ON anc.parent_id = pc.id
			WHERE ? = TRUE
		),
		des AS (
			SELECT id FROM place_classes WHERE id = ANY(?)
			UNION ALL
			SELECT pc.id FROM place_classes pc
			JOIN des ON pc.parent_id = des.id
			WHERE ? = TRUE
		)
		SELECT DISTINCT id FROM (SELECT id FROM anc UNION SELECT id FROM des) t
	`

	expr := sq.Expr(
		"class_id IN ("+subSQL+")",
		classIDs, includeParents,
		classIDs, includeChildren,
	)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)
	return q
}

func (q *placesQ) FilterByStatus(status ...string) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q *placesQ) FilterByVerified(verified bool) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"verified": verified})
	q.counter = q.counter.Where(sq.Eq{"verified": verified})
	q.updater = q.updater.Where(sq.Eq{"verified": verified})
	q.deleter = q.deleter.Where(sq.Eq{"verified": verified})
	return q
}

func (q *placesQ) FilterByRadius(center orb.Point, radiusMeters uint) repository.PlacesQ {
	expr := sq.Expr(
		`ST_DWithin(point, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)`,
		center[0], center[1], int64(radiusMeters),
	)
	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	return q
}

func (q *placesQ) FilterBestMatch(text string) repository.PlacesQ {
	if text == "" {
		return q
	}
	pattern := "%" + text + "%"
	expr := sq.Or{
		sq.Expr("name ILIKE ?", pattern),
		sq.Expr("address ILIKE ?", pattern),
		sq.Expr("description ILIKE ?", pattern),
	}
	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	return q
}

func (q *placesQ) FilterByOrgStatus(status string) repository.PlacesQ {
	expr := sq.Expr(
		"organization_id IN (SELECT id FROM organizations WHERE status = ?)",
		status,
	)
	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	return q
}

func (q *placesQ) Get(ctx context.Context) (repository.PlaceRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.PlaceRow{}, fmt.Errorf("building select query for %s: %w", placesTable, err)
	}

	return scanPlaceRow(q.db.QueryRow(ctx, query, args...))
}

func (q *placesQ) Select(ctx context.Context) ([]repository.PlaceRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", placesTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", placesTable, err)
	}
	defer rows.Close()

	out := make([]repository.PlaceRow, 0)
	for rows.Next() {
		p, err := scanPlaceRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *placesQ) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", placesTable, err)
	}

	var ok bool
	if err = q.db.QueryRow(ctx, "SELECT EXISTS ("+subSQL+")", subArgs...).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning exists for %s: %w", placesTable, err)
	}

	return ok, nil
}

func (q *placesQ) UpdateOne(ctx context.Context) (repository.PlaceRow, error) {
	q.updater = q.updater.
		Set("updated_at", time.Now().UTC()).
		Set("version", sq.Expr("version + 1"))

	query, args, err := q.updater.Suffix("RETURNING " + placesColumns).ToSql()
	if err != nil {
		return repository.PlaceRow{}, fmt.Errorf("building update query for %s: %w", placesTable, err)
	}

	return scanPlaceRow(q.db.QueryRow(ctx, query, args...))
}

func (q *placesQ) UpdateClassID(classID uuid.UUID) repository.PlacesQ {
	q.updater = q.updater.Set("class_id", classID)
	return q
}

func (q *placesQ) UpdateStatus(status string) repository.PlacesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q *placesQ) UpdateVerified(verified bool) repository.PlacesQ {
	q.updater = q.updater.Set("verified", verified)
	return q
}

func (q *placesQ) UpdateName(name string) repository.PlacesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q *placesQ) UpdateAddress(address string) repository.PlacesQ {
	q.updater = q.updater.Set("address", address)
	return q
}

func (q *placesQ) UpdateDescription(description *string) repository.PlacesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q *placesQ) UpdateIconKey(icon *string) repository.PlacesQ {
	q.updater = q.updater.Set("icon_key", icon)
	return q
}

func (q *placesQ) UpdateBannerKey(banner *string) repository.PlacesQ {
	q.updater = q.updater.Set("banner_key", banner)
	return q
}

func (q *placesQ) UpdateWebsite(website *string) repository.PlacesQ {
	q.updater = q.updater.Set("website", website)
	return q
}

func (q *placesQ) UpdatePhone(phone *string) repository.PlacesQ {
	q.updater = q.updater.Set("phone", phone)
	return q
}

func (q *placesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", placesTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", placesTable, err)
	}

	return nil
}

func (q *placesQ) Page(limit, offset uint) repository.PlacesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q *placesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", placesTable, err)
	}

	var count uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", placesTable, err)
	}

	return count, nil
}
