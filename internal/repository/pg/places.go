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

const PlacesTable = "places"

const PlacesColumns = `id, class_id, organization_id,status, verified, point, address, name,
	description, icon, banner, website, phone, created_at, updated_at`

func scanPlaceRow(row sq.RowScanner) (r repository.PlaceRow, err error) {
	var organizationID pgtype.UUID
	var description pgtype.Text
	var icon pgtype.Text
	var banner pgtype.Text
	var website pgtype.Text
	var phone pgtype.Text

	err = row.Scan(
		&r.ID,
		&r.ClassID,
		&organizationID,
		&r.Status,
		&r.Verified,
		&r.Point,
		&r.Address,
		&r.Name,
		&description,
		&icon,
		&banner,
		&website,
		&phone,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.PlaceRow{}, nil
	case err != nil:
		return repository.PlaceRow{}, fmt.Errorf("scanning place row: %w", err)
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
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return placesQ{
		db:       db,
		selector: builder.Select(PlacesColumns).From(PlacesTable),
		inserter: builder.Insert(PlacesTable),
		updater:  builder.Update(PlacesTable),
		deleter:  builder.Delete(PlacesTable),
		counter:  builder.Select("COUNT(*) AS count").From(PlacesTable),
	}
}

func (q placesQ) New() repository.PlacesQ {
	return NewPlacesQ(q.db)
}

func (q placesQ) Insert(ctx context.Context, data repository.PlaceRow) (repository.PlaceRow, error) {
	set := map[string]interface{}{
		"class_id":    data.ClassID,
		"status":      data.Status,
		"verified":    data.Verified,
		"point":       data.Point,
		"address":     data.Address,
		"name":        data.Name,
		"description": data.Description,
		"icon":        data.IconKey,
		"banner":      data.BannerKey,
		"website":     data.Website,
		"phone":       data.Phone,
	}

	if data.OrganizationID != nil {
		set["organization_id"] = *data.OrganizationID
	}

	query, args, err := q.inserter.
		SetMap(set).
		Suffix("RETURNING " + PlacesColumns).
		ToSql()
	if err != nil {
		return repository.PlaceRow{}, fmt.Errorf("building insert query for %s: %w", PlacesTable, err)
	}

	return scanPlaceRow(q.db.QueryRow(ctx, query, args...))
}

func (q placesQ) FilterByID(id uuid.UUID) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q placesQ) FilterByOrganizationID(orgID *uuid.UUID) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"organization_id": orgID})
	q.counter = q.counter.Where(sq.Eq{"organization_id": orgID})
	q.updater = q.updater.Where(sq.Eq{"organization_id": orgID})
	q.deleter = q.deleter.Where(sq.Eq{"organization_id": orgID})
	return q
}

func (q placesQ) FilterByClassID(includeChild, includeParent bool, classIDs ...uuid.UUID) repository.PlacesQ {
	if len(classIDs) == 0 {
		return q
	}

	if !includeChild && !includeParent {
		q.selector = q.selector.Where(sq.Eq{"class_id": classIDs})
		q.counter = q.counter.Where(sq.Eq{"class_id": classIDs})
		q.deleter = q.deleter.Where(sq.Eq{"class_id": classIDs})
		q.updater = q.updater.Where(sq.Eq{"class_id": classIDs})
		return q
	}

	subSQL := `
		WITH RECURSIVE
		anc AS (
			SELECT id, parent_id
			FROM place_classes
			WHERE id = ANY(?)
			UNION ALL
			SELECT pc.id, pc.parent_id
			FROM place_classes pc
			JOIN anc ON anc.parent_id = pc.id
			WHERE (? = TRUE)
		),
		des AS (
			SELECT id
			FROM place_classes
			WHERE id = ANY(?)
			UNION ALL
			SELECT pc.id
			FROM place_classes pc
			JOIN des ON pc.parent_id = des.id
			WHERE (? = TRUE)
		)
		SELECT DISTINCT id
		FROM (
			SELECT id FROM anc
			UNION
			SELECT id FROM des
		) t
	`

	// Важно: чтобы ANY(?) работал, аргумент должен быть массивом для драйвера.
	// Для pgdbx/stdlib обычно ок передать []uuid.UUID, но если будут проблемы — нужно pq.Array(...) (для lib/pq).
	expr := sq.Expr(
		"class_id IN ("+subSQL+")",
		classIDs, includeParent,
		classIDs, includeChild,
	)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.deleter = q.deleter.Where(expr)
	q.updater = q.updater.Where(expr)
	return q
}

func (q placesQ) FilterByStatus(status ...string) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q placesQ) FilterByVerified(verified bool) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"verified": verified})
	q.counter = q.counter.Where(sq.Eq{"verified": verified})
	q.updater = q.updater.Where(sq.Eq{"verified": verified})
	q.deleter = q.deleter.Where(sq.Eq{"verified": verified})
	return q
}

func (q placesQ) FilterByParentID(parentID uuid.UUID) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"parent_id": parentID})
	q.counter = q.counter.Where(sq.Eq{"parent_id": parentID})
	q.updater = q.updater.Where(sq.Eq{"parent_id": parentID})
	q.deleter = q.deleter.Where(sq.Eq{"parent_id": parentID})
	return q
}

func (q placesQ) FilterByRadius(center orb.Point, radiusMeters uint) repository.PlacesQ {
	// orb.Point = [lng, lat]
	lng := center[0]
	lat := center[1]

	// point — geography
	// ST_MakePoint(lng, lat) -> geometry, затем ST_SetSRID(..., 4326)
	expr := sq.Expr(
		`ST_DWithin(point, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)`,
		lng, lat, int64(radiusMeters),
	)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)

	return q
}

func (q placesQ) FilterByText(text string) repository.PlacesQ {
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

func (q placesQ) FilterLikeAddress(address string) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"address": "%s" + address + "%"})
	q.counter = q.counter.Where(sq.Eq{"address": "%s" + address + "%"})
	return q
}

func (q placesQ) FilterLikeName(name string) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"name": "%s" + name + "%"})
	q.counter = q.counter.Where(sq.Eq{"name": "%s" + name + "%"})
	return q
}

func (q placesQ) FilterLikeDescription(description string) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"description": "%" + description + "%"})
	q.counter = q.counter.Where(sq.Eq{"description": "%" + description + "%"})
	return q
}

func (q placesQ) Get(ctx context.Context) (repository.PlaceRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.PlaceRow{}, fmt.Errorf("building select query for %s: %w", PlacesTable, err)
	}

	return scanPlaceRow(q.db.QueryRow(ctx, query, args...))
}

func (q placesQ) Select(ctx context.Context) ([]repository.PlaceRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlacesTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PlacesTable, err)
	}
	defer rows.Close()

	var places []repository.PlaceRow
	for rows.Next() {
		p, err := scanPlaceRow(rows)
		if err != nil {
			return nil, err
		}
		places = append(places, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return places, nil
}

func (q placesQ) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", PlacesTable, err)
	}

	sqlq := "SELECT EXISTS (" + subSQL + ")"

	var ok bool
	if err = q.db.QueryRow(ctx, sqlq, subArgs...).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning exists for %s: %w", PlacesTable, err)
	}

	return ok, nil
}

func (q placesQ) UpdateOne(ctx context.Context) (repository.PlaceRow, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.
		Suffix("RETURNING " + PlacesColumns).
		ToSql()
	if err != nil {
		return repository.PlaceRow{}, fmt.Errorf("building update query for %s: %w", PlacesTable, err)
	}

	return scanPlaceRow(q.db.QueryRow(ctx, query, args...))
}

func (q placesQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", PlacesTable, err)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", PlacesTable, err)
	}

	return res.RowsAffected(), nil
}

func (q placesQ) UpdateClassID(classID uuid.UUID) repository.PlacesQ {
	q.updater = q.updater.Set("class_id", classID)
	return q
}

func (q placesQ) UpdateStatus(status string) repository.PlacesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q placesQ) UpdateVerified(verified bool) repository.PlacesQ {
	q.updater = q.updater.Set("verified", verified)
	return q
}

func (q placesQ) UpdateName(name string) repository.PlacesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q placesQ) UpdateAddress(address string) repository.PlacesQ {
	q.updater = q.updater.Set("address", address)
	return q
}

func (q placesQ) UpdateDescription(description *string) repository.PlacesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q placesQ) UpdateIcon(icon *string) repository.PlacesQ {
	q.updater = q.updater.Set("icon", icon)
	return q
}

func (q placesQ) UpdateBanner(banner *string) repository.PlacesQ {
	q.updater = q.updater.Set("banner", banner)
	return q
}

func (q placesQ) UpdateWebsite(website *string) repository.PlacesQ {
	q.updater = q.updater.Set("website", website)
	return q
}

func (q placesQ) UpdatePhone(phone *string) repository.PlacesQ {
	q.updater = q.updater.Set("phone", phone)
	return q
}

func (q placesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlacesTable, err)
	}

	if _, err := q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PlacesTable, err)
	}

	return nil
}

func (q placesQ) Page(limit, offset uint) repository.PlacesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q placesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlacesTable, err)
	}

	var count uint
	if err := q.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PlacesTable, err)
	}

	return count, nil
}
