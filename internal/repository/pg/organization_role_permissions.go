package pg

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/pgdbx"
	"github.com/netbill/places-svc/internal/repository"
)

const organizationRolePermissionsTable = "organization_role_permissions"

const organizationRolePermissionsColumns = "code, description"
const organizationRolePermissionsColumnsP = "p.code, p.description"

func scanOrganizationRolePermission(row sq.RowScanner) (p repository.OrgRolePermissionRow, err error) {
	err = row.Scan(
		&p.Code,
		&p.Description,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrgRolePermissionRow{}, nil
	case err != nil:
		return repository.OrgRolePermissionRow{}, fmt.Errorf("scanning organization role permission: %w", err)
	}

	return p, nil
}

type organizationRolePermissions struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrganizationRolePermissionsQ(db *pgdbx.DB) repository.OrgRolePermissionsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &organizationRolePermissions{
		db:       db,
		selector: b.Select(organizationRolePermissionsColumnsP).From(organizationRolePermissionsTable + " p"),
		inserter: b.Insert(organizationRolePermissionsTable),
		updater:  b.Update(organizationRolePermissionsTable + " p"),
		deleter:  b.Delete(organizationRolePermissionsTable + " p"),
		counter:  b.Select("COUNT(*)").From(organizationRolePermissionsTable + " p"),
	}
}

func (q *organizationRolePermissions) New() repository.OrgRolePermissionsQ {
	return NewOrganizationRolePermissionsQ(q.db)
}

func (q *organizationRolePermissions) Insert(
	ctx context.Context,
	data repository.OrgRolePermissionRow,
) (repository.OrgRolePermissionRow, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"code":        data.Code,
		"description": data.Description,
	}).Suffix("RETURNING " + organizationRolePermissionsColumns).ToSql()
	if err != nil {
		return repository.OrgRolePermissionRow{}, fmt.Errorf(
			"building insert query for %s: %w",
			organizationRolePermissionsTable,
			err,
		)
	}

	return scanOrganizationRolePermission(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationRolePermissions) Get(ctx context.Context) (repository.OrgRolePermissionRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrgRolePermissionRow{}, fmt.Errorf(
			"building select query for %s: %w",
			organizationRolePermissionsTable,
			err,
		)
	}

	return scanOrganizationRolePermission(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationRolePermissions) Select(ctx context.Context) ([]repository.OrgRolePermissionRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf(
			"building select query for %s: %w",
			organizationRolePermissionsTable,
			err,
		)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"executing select query for %s: %w",
			organizationRolePermissionsTable,
			err,
		)
	}
	defer rows.Close()

	out := make([]repository.OrgRolePermissionRow, 0)
	for rows.Next() {
		p, err := scanOrganizationRolePermission(rows)
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

func (q *organizationRolePermissions) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationRolePermissionsTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationRolePermissionsTable, err)
	}

	return nil
}

func (q *organizationRolePermissions) FilterByCode(code ...string) repository.OrgRolePermissionsQ {
	q.selector = q.selector.Where(sq.Eq{"p.code": code})
	q.counter = q.counter.Where(sq.Eq{"p.code": code})
	q.updater = q.updater.Where(sq.Eq{"p.code": code})
	q.deleter = q.deleter.Where(sq.Eq{"p.code": code})
	return q
}

func (q *organizationRolePermissions) UpdateOne(ctx context.Context) (repository.OrgRolePermissionRow, error) {
	query, args, err := q.updater.Suffix("RETURNING " + organizationRolePermissionsColumns).ToSql()
	if err != nil {
		return repository.OrgRolePermissionRow{}, fmt.Errorf(
			"building update query for %s: %w",
			organizationRolePermissionsTable,
			err,
		)
	}

	return scanOrganizationRolePermission(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationRolePermissions) UpdateMany(ctx context.Context) (int64, error) {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf(
			"building update query for %s: %w",
			organizationRolePermissionsTable,
			err,
		)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf(
			"executing update query for %s: %w",
			organizationRolePermissionsTable,
			err,
		)
	}

	return res.RowsAffected(), nil
}

func (q *organizationRolePermissions) UpdateDescription(description string) repository.OrgRolePermissionsQ {
	q.updater = q.updater.Set("description", description)
	return q
}
