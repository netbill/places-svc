package pg

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/pgdbx"

	"github.com/netbill/places-svc/internal/repository"
)

const organizationRolePermissionLinksTable = "organization_role_permission_links"

const organizationRolePermissionLinksColumns = "role_id, permission_code"
const organizationRolePermissionLinksColumnsP = "l.role_id, l.permission_code"

func scanOrganizationRolePermissionLink(row sq.RowScanner) (l repository.OrganizationRolePermissionLinkRow, err error) {
	err = row.Scan(
		&l.RoleID,
		&l.PermissionCode,
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
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgRolePermissionLinksQ(db *pgdbx.DB) repository.OrgRolePermissionLinksQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &organizationRolePermissionLinks{
		db:       db,
		selector: b.Select(organizationRolePermissionLinksColumnsP).From(organizationRolePermissionLinksTable + " l"),
		inserter: b.Insert(organizationRolePermissionLinksTable),
		deleter:  b.Delete(organizationRolePermissionLinksTable + " l"),
		counter:  b.Select("COUNT(*) AS count").From(organizationRolePermissionLinksTable + " l"),
	}
}

func (q *organizationRolePermissionLinks) New() repository.OrgRolePermissionLinksQ {
	return NewOrgRolePermissionLinksQ(q.db)
}

func (q *organizationRolePermissionLinks) Upsert(
	ctx context.Context,
	roleID uuid.UUID,
	codes ...string,
) ([]repository.OrganizationRolePermissionLinkRow, error) {
	uniq := make([]string, 0, len(codes))
	seen := make(map[string]struct{}, len(codes))

	for _, c := range codes {
		if c == "" {
			continue
		}
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		uniq = append(uniq, c)
	}

	const sqlq = `
		WITH del AS (
			DELETE FROM organization_role_permission_links
			WHERE role_id = $1
		)
		INSERT INTO organization_role_permission_links (role_id, permission_code)
		SELECT $1, x.code
		FROM UNNEST($2::text[]) AS x(code)
		RETURNING role_id, permission_code
	`

	rows, err := q.db.Query(ctx, sqlq, roleID, uniq)
	if err != nil {
		return nil, fmt.Errorf("upsert role permission links: %w", err)
	}
	defer rows.Close()

	out := make([]repository.OrganizationRolePermissionLinkRow, 0, len(uniq))
	for rows.Next() {
		var r repository.OrganizationRolePermissionLinkRow
		if err := rows.Scan(&r.RoleID, &r.PermissionCode); err != nil {
			return nil, fmt.Errorf("scanning role permission link: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
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

func (q *organizationRolePermissionLinks) Select(ctx context.Context) ([]repository.OrganizationRolePermissionLinkRow, error) {
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

func (q *organizationRolePermissionLinks) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf(
			"building exists query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	sqlq := "SELECT EXISTS (" + subSQL + ")"

	var ok bool
	if err = q.db.QueryRow(ctx, sqlq, subArgs...).Scan(&ok); err != nil {
		return false, fmt.Errorf(
			"scanning exists for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	return ok, nil
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

func (q *organizationRolePermissionLinks) FilterByRoleID(roleID ...uuid.UUID) repository.OrgRolePermissionLinksQ {
	q.selector = q.selector.Where(sq.Eq{"l.role_id": roleID})
	q.counter = q.counter.Where(sq.Eq{"l.role_id": roleID})
	q.deleter = q.deleter.Where(sq.Eq{"l.role_id": roleID})
	return q
}

func (q *organizationRolePermissionLinks) FilterByPermissionCode(code ...string) repository.OrgRolePermissionLinksQ {
	q.selector = q.selector.Where(sq.Eq{"l.permission_code": code})
	q.counter = q.counter.Where(sq.Eq{"l.permission_code": code})
	q.deleter = q.deleter.Where(sq.Eq{"l.permission_code": code})
	return q
}

func (q *organizationRolePermissionLinks) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf(
			"building count query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	var n uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf(
			"scanning count for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	return n, nil
}

func (q *organizationRolePermissionLinks) Page(limit, offset uint) repository.OrgRolePermissionLinksQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q *organizationRolePermissionLinks) Exist(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", organizationRolePermissionLinksTable, err)
	}

	sqlq := "SELECT EXISTS (" + subSQL + ")"

	var ok bool
	if err = q.db.QueryRow(ctx, sqlq, subArgs...).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning exists for %s: %w", PlacesTable, err)
	}

	return ok, nil
}
