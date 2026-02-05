package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/pgdbx"
	"github.com/netbill/places-svc/internal/repository"
)

const organizationRolesTable = "organization_roles"

const organizationRolesColumns = "id, organization_id, rank, source_created_at, source_updated_at, replica_created_at, replica_updated_at"

const organizationRolesColumnsP = "r.id, r.organization_id, r.rank, r.source_created_at, r.source_updated_at, r.replica_created_at, r.replica_updated_at"

func scanOrganizationRole(row sq.RowScanner) (r repository.OrgRoleRow, err error) {
	err = row.Scan(
		&r.ID,
		&r.OrganizationID,
		&r.Rank,
		&r.SourceCreatedAt,
		&r.SourceUpdatedAt,
		&r.ReplicaCreatedAt,
		&r.ReplicaUpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrgRoleRow{}, nil
	case err != nil:
		return repository.OrgRoleRow{}, fmt.Errorf("scanning organization role: %w", err)
	}

	return r, nil
}

type organizationRoles struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgRolesQ(db *pgdbx.DB) repository.OrgRolesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &organizationRoles{
		db:       db,
		selector: b.Select(organizationRolesColumnsP).From(organizationRolesTable + " r"),
		inserter: b.Insert(organizationRolesTable),
		updater:  b.Update(organizationRolesTable + " r"),
		deleter:  b.Delete(organizationRolesTable + " r"),
		counter:  b.Select("COUNT(*)").From(organizationRolesTable + " r"),
	}
}

func (q *organizationRoles) New() repository.OrgRolesQ {
	return NewOrgRolesQ(q.db)
}

func (q *organizationRoles) Insert(
	ctx context.Context,
	data repository.OrgRoleRow,
) (repository.OrgRoleRow, error) {
	now := time.Now().UTC()

	query, args, err := q.inserter.SetMap(map[string]any{
		"id":                 data.ID,
		"organization_id":    data.OrganizationID,
		"rank":               data.Rank,
		"source_created_at":  data.SourceCreatedAt.UTC(),
		"source_updated_at":  data.SourceUpdatedAt.UTC(),
		"replica_created_at": now,
		"replica_updated_at": now,
	}).Suffix("RETURNING " + organizationRolesColumns).ToSql()
	if err != nil {
		return repository.OrgRoleRow{}, fmt.Errorf(
			"building insert query for %s: %w",
			organizationRolesTable,
			err,
		)
	}

	return scanOrganizationRole(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationRoles) Get(ctx context.Context) (repository.OrgRoleRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrgRoleRow{}, fmt.Errorf(
			"building select query for %s: %w",
			organizationRolesTable,
			err,
		)
	}

	return scanOrganizationRole(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationRoles) Select(ctx context.Context) ([]repository.OrgRoleRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf(
			"building select query for %s: %w",
			organizationRolesTable,
			err,
		)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"executing select query for %s: %w",
			organizationRolesTable,
			err,
		)
	}
	defer rows.Close()

	out := make([]repository.OrgRoleRow, 0)
	for rows.Next() {
		r, err := scanOrganizationRole(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *organizationRoles) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationRolesTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationRolesTable, err)
	}

	return nil
}

func (q *organizationRoles) FilterByID(id ...uuid.UUID) repository.OrgRolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.id": id})
	q.counter = q.counter.Where(sq.Eq{"r.id": id})
	q.updater = q.updater.Where(sq.Eq{"r.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"r.id": id})
	return q
}

func (q *organizationRoles) FilterByOrganizationID(
	organizationID ...uuid.UUID,
) repository.OrgRolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.organization_id": organizationID})
	q.counter = q.counter.Where(sq.Eq{"r.organization_id": organizationID})
	q.updater = q.updater.Where(sq.Eq{"r.organization_id": organizationID})
	q.deleter = q.deleter.Where(sq.Eq{"r.organization_id": organizationID})
	return q
}

func (q *organizationRoles) FilterByRank(rank uint) repository.OrgRolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.rank": rank})
	q.counter = q.counter.Where(sq.Eq{"r.rank": rank})
	q.updater = q.updater.Where(sq.Eq{"r.rank": rank})
	q.deleter = q.deleter.Where(sq.Eq{"r.rank": rank})
	return q
}

func (q *organizationRoles) UpdateOne(ctx context.Context) (repository.OrgRoleRow, error) {
	q.updater = q.updater.Set("replica_updated_at", time.Now().UTC())

	query, args, err := q.updater.
		Suffix("RETURNING " + organizationRolesColumns).
		ToSql()
	if err != nil {
		return repository.OrgRoleRow{}, fmt.Errorf(
			"building update query for %s: %w",
			organizationRolesTable,
			err,
		)
	}

	return scanOrganizationRole(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationRoles) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("replica_updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf(
			"building update query for %s: %w",
			organizationRolesTable,
			err,
		)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf(
			"executing update query for %s: %w",
			organizationRolesTable,
			err,
		)
	}

	return res.RowsAffected(), nil
}

func (q *organizationRoles) UpdateRank(rank uint) repository.OrgRolesQ {
	q.updater = q.updater.Set("rank", rank)
	return q
}

func (q *organizationRoles) UpdateSourceUpdatedAt(
	updatedAt time.Time,
) repository.OrgRolesQ {
	q.updater = q.updater.Set("source_updated_at", updatedAt.UTC())
	return q
}

func (q *organizationRoles) OrderByRank(asc bool) repository.OrgRolesQ {
	if asc {
		q.selector = q.selector.OrderBy("r.rank ASC", "r.id ASC")
	} else {
		q.selector = q.selector.OrderBy("r.rank DESC", "r.id DESC")
	}
	return q
}

func (q *organizationRoles) UpdateRolesRanks(
	ctx context.Context,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
	updatedAt time.Time,
) error {
	roles, err := NewOrgRolesQ(q.db).
		FilterByOrganizationID(organizationID).
		OrderByRank(true).
		Select(ctx)
	if err != nil {
		return fmt.Errorf("select roles by organization: %w", err)
	}
	if len(roles) == 0 {
		return fmt.Errorf("no roles in organization %s", organizationID)
	}

	n := uint(len(roles))

	idToRole := make(map[uuid.UUID]repository.OrgRoleRow, n)
	for i := range roles {
		idToRole[roles[i].ID] = roles[i]
	}

	usedRank := make(map[uint]uuid.UUID, len(order))
	for roleID, newRank := range order {
		if newRank >= n {
			return fmt.Errorf("rank %d out of range [0..%d]", newRank, n-1)
		}
		if _, ok := idToRole[roleID]; !ok {
			return fmt.Errorf("role %s not in organization %s", roleID, organizationID)
		}
		if prev, ok := usedRank[newRank]; ok && prev != roleID {
			return fmt.Errorf("duplicate rank %d for roles %s and %s", newRank, prev, roleID)
		}
		usedRank[newRank] = roleID
	}

	target := make([]uuid.UUID, n)
	filled := make([]bool, n)

	for rnk, id := range usedRank {
		target[rnk] = id
		filled[rnk] = true
	}

	rest := make([]uuid.UUID, 0, n-uint(len(order)))
	for i := range roles {
		id := roles[i].ID
		if _, ok := order[id]; ok {
			continue
		}
		rest = append(rest, id)
	}

	j := 0
	for i := 0; uint(i) < n; i++ {
		if filled[i] {
			continue
		}
		target[i] = rest[j]
		j++
	}

	changed := make([]uuid.UUID, 0, n)
	newRanks := make([]int32, 0, n)

	for newRank, id := range target {
		if roles[newRank].ID != id {
			changed = append(changed, id)
			newRanks = append(newRanks, int32(newRank))
		}
	}

	if len(changed) > 0 {
		const sqlUpdate = `
			UPDATE organization_roles r
			SET
				rank = v.rank,
				replica_updated_at = $3
			FROM (
				SELECT UNNEST($1::uuid[]) AS id, UNNEST($2::int4[]) AS rank
			) v
			WHERE r.id = v.id
			  AND r.organization_id = $4
		`

		if _, err := q.db.Exec(
			ctx,
			sqlUpdate,
			changed,
			newRanks,
			updatedAt.UTC(),
			organizationID,
		); err != nil {
			return fmt.Errorf("updating roles ranks: %w", err)
		}
	}

	return nil
}
