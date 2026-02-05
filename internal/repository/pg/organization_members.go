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

const organizationMembersTable = "organization_members"

const organizationMembersColumns = "id, account_id, organization_id, head, source_created_at, source_updated_at, replica_created_at, replica_updated_at"

const organizationMembersColumnsP = "m.id, m.account_id, m.organization_id, m.head, m.source_created_at, m.source_updated_at, m.replica_created_at, m.replica_updated_at"

func scanOrganizationMember(row sq.RowScanner) (m repository.OrgMemberRow, err error) {
	err = row.Scan(
		&m.ID,
		&m.AccountID,
		&m.OrganizationID,
		&m.Head,
		&m.SourceCreatedAt,
		&m.SourceUpdatedAt,
		&m.ReplicaCreatedAt,
		&m.ReplicaUpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrgMemberRow{}, nil
	case err != nil:
		return repository.OrgMemberRow{}, fmt.Errorf("scanning organization member: %w", err)
	}

	return m, nil
}

type organizationMembers struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgMembersQ(db *pgdbx.DB) repository.OrgMembersQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &organizationMembers{
		db:       db,
		selector: b.Select(organizationMembersColumnsP).From(organizationMembersTable + " m"),
		inserter: b.Insert(organizationMembersTable),
		updater:  b.Update(organizationMembersTable + " m"),
		deleter:  b.Delete(organizationMembersTable + " m"),
		counter:  b.Select("COUNT(*)").From(organizationMembersTable + " m"),
	}
}

func (q *organizationMembers) New() repository.OrgMembersQ {
	return NewOrgMembersQ(q.db)
}

func (q *organizationMembers) Insert(
	ctx context.Context,
	data repository.OrgMemberRow,
) (repository.OrgMemberRow, error) {
	now := time.Now().UTC()

	query, args, err := q.inserter.SetMap(map[string]any{
		"id":                 data.ID,
		"account_id":         data.AccountID,
		"organization_id":    data.OrganizationID,
		"head":               data.Head,
		"source_created_at":  data.SourceCreatedAt.UTC(),
		"source_updated_at":  data.SourceUpdatedAt.UTC(),
		"replica_created_at": now,
		"replica_updated_at": now,
	}).Suffix("RETURNING " + organizationMembersColumns).ToSql()
	if err != nil {
		return repository.OrgMemberRow{}, fmt.Errorf(
			"building insert query for %s: %w",
			organizationMembersTable,
			err,
		)
	}

	return scanOrganizationMember(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationMembers) Get(ctx context.Context) (repository.OrgMemberRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrgMemberRow{}, fmt.Errorf(
			"building select query for %s: %w",
			organizationMembersTable,
			err,
		)
	}

	return scanOrganizationMember(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationMembers) Select(ctx context.Context) ([]repository.OrgMemberRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf(
			"building select query for %s: %w",
			organizationMembersTable,
			err,
		)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"executing select query for %s: %w",
			organizationMembersTable,
			err,
		)
	}
	defer rows.Close()

	out := make([]repository.OrgMemberRow, 0)
	for rows.Next() {
		m, err := scanOrganizationMember(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *organizationMembers) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationMembersTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationMembersTable, err)
	}

	return nil
}

func (q *organizationMembers) FilterByID(id ...uuid.UUID) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.id": id})
	q.counter = q.counter.Where(sq.Eq{"m.id": id})
	q.updater = q.updater.Where(sq.Eq{"m.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"m.id": id})
	return q
}

func (q *organizationMembers) FilterByAccountID(accountID ...uuid.UUID) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.account_id": accountID})
	q.counter = q.counter.Where(sq.Eq{"m.account_id": accountID})
	q.updater = q.updater.Where(sq.Eq{"m.account_id": accountID})
	q.deleter = q.deleter.Where(sq.Eq{"m.account_id": accountID})
	return q
}

func (q *organizationMembers) FilterByOrganizationID(
	organizationID ...uuid.UUID,
) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.organization_id": organizationID})
	q.counter = q.counter.Where(sq.Eq{"m.organization_id": organizationID})
	q.updater = q.updater.Where(sq.Eq{"m.organization_id": organizationID})
	q.deleter = q.deleter.Where(sq.Eq{"m.organization_id": organizationID})
	return q
}

func (q *organizationMembers) FilterByHead(head bool) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.head": head})
	q.counter = q.counter.Where(sq.Eq{"m.head": head})
	q.updater = q.updater.Where(sq.Eq{"m.head": head})
	q.deleter = q.deleter.Where(sq.Eq{"m.head": head})
	return q
}
