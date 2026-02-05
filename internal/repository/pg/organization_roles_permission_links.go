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

const organizationRolePermissionLinksTable = "organization_role_permission_links"

const organizationRolePermissionLinksColumns = "role_id, permission_code, source_created_at, replica_created_at"

const organizationRolePermissionLinksColumnsP = "l.role_id, l.permission_code, l.source_created_at, l.replica_created_at"

func scanOrganizationRolePermissionLink(row sq.RowScanner) (l repository.OrganizationRolePermissionLinkRow, err error) {
	err = row.Scan(
		&l.RoleID,
		&l.PermissionCode,
		&l.SourceCreatedAt,
		&l.ReplicaCreatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrganizationRolePermissionLinkRow{}, nil
	case err != nil:
		return repository.OrganizationRolePermissionLinkRow{}, fmt.Errorf("scanning organization role permission link: %w", err)
	}

	return l, nil
}

type organizationRolePermissionLinks struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgRolePermissionLinksQ(db *pgdbx.DB) repository.OrgRolePermissionLinksQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &organizationRolePermissionLinks{
		db:       db,
		selector: b.Select(organizationRolePermissionLinksColumnsP).From(organizationRolePermissionLinksTable + " l"),
		inserter: b.Insert(organizationRolePermissionLinksTable),
		updater:  b.Update(organizationRolePermissionLinksTable + " l"),
		deleter:  b.Delete(organizationRolePermissionLinksTable + " l"),
		counter:  b.Select("COUNT(*)").From(organizationRolePermissionLinksTable + " l"),
	}
}

func (q *organizationRolePermissionLinks) New() repository.OrgRolePermissionLinksQ {
	return NewOrgRolePermissionLinksQ(q.db)
}

func (q *organizationRolePermissionLinks) Insert(
	ctx context.Context,
	data repository.OrganizationRolePermissionLinkRow,
) (repository.OrganizationRolePermissionLinkRow, error) {
	now := time.Now().UTC()

	query, args, err := q.inserter.SetMap(map[string]any{
		"role_id":            data.RoleID,
		"permission_code":    data.PermissionCode,
		"source_created_at":  data.SourceCreatedAt.UTC(),
		"replica_created_at": now,
	}).Suffix("RETURNING " + organizationRolePermissionLinksColumns).ToSql()
	if err != nil {
		return repository.OrganizationRolePermissionLinkRow{}, fmt.Errorf(
			"building insert query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	return scanOrganizationRolePermissionLink(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationRolePermissionLinks) Get(ctx context.Context) (repository.OrganizationRolePermissionLinkRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrganizationRolePermissionLinkRow{}, fmt.Errorf(
			"building select query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	return scanOrganizationRolePermissionLink(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationRolePermissionLinks) Select(
	ctx context.Context,
) ([]repository.OrganizationRolePermissionLinkRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf(
			"building select query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"executing select query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}
	defer rows.Close()

	out := make([]repository.OrganizationRolePermissionLinkRow, 0)
	for rows.Next() {
		l, err := scanOrganizationRolePermissionLink(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *organizationRolePermissionLinks) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationRolePermissionLinksTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationRolePermissionLinksTable, err)
	}

	return nil
}

func (q *organizationRolePermissionLinks) FilterByRoleID(
	roleID ...uuid.UUID,
) repository.OrgRolePermissionLinksQ {
	q.selector = q.selector.Where(sq.Eq{"l.role_id": roleID})
	q.counter = q.counter.Where(sq.Eq{"l.role_id": roleID})
	q.updater = q.updater.Where(sq.Eq{"l.role_id": roleID})
	q.deleter = q.deleter.Where(sq.Eq{"l.role_id": roleID})
	return q
}

func (q *organizationRolePermissionLinks) FilterByPermissionCode(
	code ...string,
) repository.OrgRolePermissionLinksQ {
	q.selector = q.selector.Where(sq.Eq{"l.permission_code": code})
	q.counter = q.counter.Where(sq.Eq{"l.permission_code": code})
	q.updater = q.updater.Where(sq.Eq{"l.permission_code": code})
	q.deleter = q.deleter.Where(sq.Eq{"l.permission_code": code})
	return q
}

func (q *organizationRolePermissionLinks) UpdateOne(
	ctx context.Context,
) (repository.OrganizationRolePermissionLinkRow, error) {
	query, args, err := q.updater.Suffix("RETURNING " + organizationRolePermissionLinksColumns).ToSql()
	if err != nil {
		return repository.OrganizationRolePermissionLinkRow{}, fmt.Errorf(
			"building update query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	return scanOrganizationRolePermissionLink(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationRolePermissionLinks) UpdateMany(ctx context.Context) (int64, error) {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf(
			"building update query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf(
			"executing update query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	return res.RowsAffected(), nil
}

func (q *organizationRolePermissionLinks) UpdateSourceCreatedAt(
	v time.Time,
) repository.OrgRolePermissionLinksQ {
	q.updater = q.updater.Set("source_created_at", v.UTC())
	return q
}
