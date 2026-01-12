package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/netbill/pgx"
)

const OrganizationMembersTable = "organization_members"

const OrganizationMemberColumns = "id, account_id, organization_id, created_at, updated_at"
const OrganizationMemberColumnsM = "m.id, m.account_id, m.organization_id, m.created_at, m.updated_at"

func (m *OrgMember) scan(row sq.RowScanner) error {
	err := row.Scan(
		&m.ID,
		&m.AccountID,
		&m.OrganizationID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("scanning member: %w", err)
	}
	return nil
}

type OrgMembersQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgMembersQ(db pgx.DBTX) OrgMembersQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return OrgMembersQ{
		db:       db,
		selector: builder.Select(OrganizationMemberColumnsM).From(OrganizationMembersTable + " m"),
		inserter: builder.Insert(OrganizationMembersTable),
		// важно: updater/deleter без алиаса для стабильности
		updater: builder.Update(OrganizationMembersTable),
		deleter: builder.Delete(OrganizationMembersTable),
		counter: builder.Select("COUNT(*) AS count").From(OrganizationMembersTable + " m"),
	}
}

type OrgMember struct {
	ID             uuid.UUID `json:"id"`
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q OrgMembersQ) Upsert(
	ctx context.Context,
	data OrgMember,
) error {

	const sqlUpsert = `
		INSERT INTO organization_members (account_id, organization_id)
		VALUES ($1, $2)
		ON CONFLICT (account_id, organization_id)
		DO UPDATE SET
			updated_at = now()
	`

	if _, err := q.db.ExecContext(
		ctx,
		sqlUpsert,
		data.AccountID,
		data.OrganizationID,
	); err != nil {
		return fmt.Errorf("upserting organization member: %w", err)
	}

	return nil
}

type InsertMemberParams struct {
	AccountID      uuid.UUID
	OrganizationID uuid.UUID
}

func (q OrgMembersQ) Insert(ctx context.Context, data InsertMemberParams) (OrgMember, error) {
	query, args, err := q.inserter.SetMap(map[string]interface{}{
		"account_id":      data.AccountID,
		"organization_id": data.OrganizationID,
	}).Suffix("RETURNING " + OrganizationMemberColumns).ToSql()
	if err != nil {
		return OrgMember{}, fmt.Errorf("building insert query for %s: %w", OrganizationMembersTable, err)
	}

	var inserted OrgMember
	if err := inserted.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return OrgMember{}, err
	}

	return inserted, nil
}

func (q OrgMembersQ) Exists(ctx context.Context) (bool, error) {
	existsQ := q.selector.
		Columns("1").
		RemoveLimit().
		RemoveOffset().
		Prefix("SELECT EXISTS (").
		Suffix(") AS exists").
		Limit(1)

	query, args, err := existsQ.ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", OrganizationMembersTable, err)
	}

	var ok bool
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning exists for %s: %w", OrganizationMembersTable, err)
	}

	return ok, nil
}

func (q OrgMembersQ) Get(ctx context.Context) (OrgMember, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return OrgMember{}, fmt.Errorf("building select query for %s: %w", OrganizationMembersTable, err)
	}

	var m OrgMember
	err = m.scan(q.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return OrgMember{}, nil
		default:
			return OrgMember{}, err
		}
	}

	return m, nil
}

func (q OrgMembersQ) Select(ctx context.Context) ([]OrgMember, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", OrganizationMembersTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", OrganizationMembersTable, err)
	}
	defer rows.Close()

	var out []OrgMember
	for rows.Next() {
		var m OrgMember
		if err := m.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q OrgMembersQ) FilterByID(id uuid.UUID) OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.id": id})
	q.counter = q.counter.Where(sq.Eq{"m.id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q OrgMembersQ) FilterByAccountID(accountID uuid.UUID) OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.account_id": accountID})
	q.counter = q.counter.Where(sq.Eq{"m.account_id": accountID})
	q.updater = q.updater.Where(sq.Eq{"account_id": accountID})
	q.deleter = q.deleter.Where(sq.Eq{"account_id": accountID})
	return q
}

func (q OrgMembersQ) FilterByOrganizationID(organizationID uuid.UUID) OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.organization_id": organizationID})
	q.counter = q.counter.Where(sq.Eq{"m.organization_id": organizationID})
	q.updater = q.updater.Where(sq.Eq{"organization_id": organizationID})
	q.deleter = q.deleter.Where(sq.Eq{"organization_id": organizationID})
	return q
}

func (q OrgMembersQ) FilterByPermissionCode(code string) OrgMembersQ {
	expr := sq.Expr(`
		EXISTS (
			SELECT 1
			FROM organization_member_roles mr
			JOIN organization_role_permission_links rp ON rp.role_id = mr.role_id
			JOIN organization_role_permissions perm ON perm.id = rp.permission_id
			WHERE mr.member_id = m.id
			  AND perm.code = ?
		)
	`, code)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)

	return q
}

func (q OrgMembersQ) UpdateOne(ctx context.Context) (OrgMember, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.Suffix("RETURNING " + OrganizationMemberColumns).ToSql()
	if err != nil {
		return OrgMember{}, fmt.Errorf("building update query for %s: %w", OrganizationMembersTable, err)
	}

	var updated OrgMember
	if err := updated.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return OrgMember{}, err
	}

	return updated, nil
}

func (q OrgMembersQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", OrganizationMembersTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", OrganizationMembersTable, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", OrganizationMembersTable, err)
	}

	return affected, nil
}

func (q OrgMembersQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", OrganizationMembersTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", OrganizationMembersTable, err)
	}

	return nil
}

func (q OrgMembersQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", OrganizationMembersTable, err)
	}

	var count uint
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", OrganizationMembersTable, err)
	}

	return count, nil
}

func (q OrgMembersQ) Page(limit uint, offset uint) OrgMembersQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}
