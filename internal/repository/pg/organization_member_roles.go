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

const organizationMemberRolesTable = "organization_member_roles"

const organizationMemberRolesColumns = "member_id, role_id, source_created_at, replica_created_at"

const organizationMemberRolesColumnsP = "mr.member_id, mr.role_id, mr.source_created_at, mr.replica_created_at"

func scanOrganizationMemberRole(row sq.RowScanner) (mr repository.OrgMemberRoleRow, err error) {
	err = row.Scan(
		&mr.MemberID,
		&mr.RoleID,
		&mr.SourceCreatedAt,
		&mr.ReplicaCreatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrgMemberRoleRow{}, nil
	case err != nil:
		return repository.OrgMemberRoleRow{}, fmt.Errorf("scanning organization member role: %w", err)
	}

	return mr, nil
}

type organizationMemberRoles struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgMemberRolesQ(db *pgdbx.DB) repository.OrgMemberRolesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &organizationMemberRoles{
		db:       db,
		selector: b.Select(organizationMemberRolesColumnsP).From(organizationMemberRolesTable + " mr"),
		inserter: b.Insert(organizationMemberRolesTable),
		updater:  b.Update(organizationMemberRolesTable + " mr"),
		deleter:  b.Delete(organizationMemberRolesTable + " mr"),
		counter:  b.Select("COUNT(*)").From(organizationMemberRolesTable + " mr"),
	}
}

func (q *organizationMemberRoles) New() repository.OrgMemberRolesQ {
	return NewOrgMemberRolesQ(q.db)
}

func (q *organizationMemberRoles) Insert(
	ctx context.Context,
	data repository.OrgMemberRoleRow,
) (repository.OrgMemberRoleRow, error) {
	now := time.Now().UTC()

	query, args, err := q.inserter.SetMap(map[string]any{
		"member_id":          data.MemberID,
		"role_id":            data.RoleID,
		"source_created_at":  data.SourceCreatedAt.UTC(),
		"replica_created_at": now,
	}).Suffix("RETURNING " + organizationMemberRolesColumns).ToSql()
	if err != nil {
		return repository.OrgMemberRoleRow{}, fmt.Errorf(
			"building insert query for %s: %w",
			organizationMemberRolesTable,
			err,
		)
	}

	return scanOrganizationMemberRole(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationMemberRoles) Get(ctx context.Context) (repository.OrgMemberRoleRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrgMemberRoleRow{}, fmt.Errorf(
			"building select query for %s: %w",
			organizationMemberRolesTable,
			err,
		)
	}

	return scanOrganizationMemberRole(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationMemberRoles) Select(
	ctx context.Context,
) ([]repository.OrgMemberRoleRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf(
			"building select query for %s: %w",
			organizationMemberRolesTable,
			err,
		)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"executing select query for %s: %w",
			organizationMemberRolesTable,
			err,
		)
	}
	defer rows.Close()

	out := make([]repository.OrgMemberRoleRow, 0)
	for rows.Next() {
		mr, err := scanOrganizationMemberRole(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, mr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *organizationMemberRoles) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationMemberRolesTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationMemberRolesTable, err)
	}

	return nil
}

func (q *organizationMemberRoles) FilterByMemberID(
	memberID ...uuid.UUID,
) repository.OrgMemberRolesQ {
	q.selector = q.selector.Where(sq.Eq{"mr.member_id": memberID})
	q.counter = q.counter.Where(sq.Eq{"mr.member_id": memberID})
	q.updater = q.updater.Where(sq.Eq{"mr.member_id": memberID})
	q.deleter = q.deleter.Where(sq.Eq{"mr.member_id": memberID})
	return q
}

func (q *organizationMemberRoles) FilterByRoleID(roleID ...uuid.UUID) repository.OrgMemberRolesQ {
	q.selector = q.selector.Where(sq.Eq{"mr.role_id": roleID})
	q.counter = q.counter.Where(sq.Eq{"mr.role_id": roleID})
	q.updater = q.updater.Where(sq.Eq{"mr.role_id": roleID})
	q.deleter = q.deleter.Where(sq.Eq{"mr.role_id": roleID})
	return q
}

func (q *organizationMemberRoles) UpdateOne(ctx context.Context) (repository.OrgMemberRoleRow, error) {
	query, args, err := q.updater.Suffix("RETURNING " + organizationMemberRolesColumns).ToSql()
	if err != nil {
		return repository.OrgMemberRoleRow{}, fmt.Errorf(
			"building update query for %s: %w",
			organizationMemberRolesTable,
			err,
		)
	}

	return scanOrganizationMemberRole(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationMemberRoles) UpdateMany(ctx context.Context) (int64, error) {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf(
			"building update query for %s: %w",
			organizationMemberRolesTable,
			err,
		)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf(
			"executing update query for %s: %w",
			organizationMemberRolesTable,
			err,
		)
	}

	return res.RowsAffected(), nil
}

func (q *organizationMemberRoles) UpdateSourceCreatedAt(v time.Time) repository.OrgMemberRolesQ {
	q.updater = q.updater.Set("source_created_at", v.UTC())
	return q
}
