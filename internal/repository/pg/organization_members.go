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
)

const organizationMembersTable = "organization_members"
const organizationMembersColumns = "id, account_id, organization_id, head, label, position, version, source_created_at, source_updated_at, replica_created_at, replica_updated_at"
const organizationMembersColumnsM = "m.id, m.account_id, m.organization_id, m.head, m.label, m.position, m.version, m.source_created_at, m.source_updated_at, m.replica_created_at, m.replica_updated_at"

func scanOrganizationMember(row sq.RowScanner) (m repository.OrgMemberRow, err error) {
	var label pgtype.Text
	var position pgtype.Text

	err = row.Scan(
		&m.ID,
		&m.AccountID,
		&m.OrganizationID,
		&m.Head,
		&label,
		&position,
		&m.Version,
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

	if label.Valid {
		m.Label = &label.String
	}
	if position.Valid {
		m.Position = &position.String
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
		selector: b.Select(organizationMembersColumnsM).From(organizationMembersTable + " m"),
		inserter: b.Insert(organizationMembersTable),
		updater:  b.Update(organizationMembersTable + " m"),
		deleter:  b.Delete(organizationMembersTable + " m"),
		counter:  b.Select("COUNT(*)").From(organizationMembersTable + " m"),
	}
}

func (q *organizationMembers) New() repository.OrgMembersQ {
	return NewOrgMembersQ(q.db)
}

func (q *organizationMembers) Insert(ctx context.Context, data repository.OrgMemberRow) error {
	query, args, err := q.inserter.SetMap(map[string]any{
		"id":                data.ID,
		"account_id":        data.AccountID,
		"organization_id":   data.OrganizationID,
		"head":              data.Head,
		"label":             data.Label,
		"position":          data.Position,
		"version":           data.Version,
		"source_created_at": data.SourceCreatedAt.UTC(),
		"source_updated_at": data.SourceUpdatedAt.UTC(),
	}).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", organizationMembersTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing insert query for %s: %w", organizationMembersTable, err)
	}

	return nil
}

func (q *organizationMembers) Get(ctx context.Context) (repository.OrgMemberRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrgMemberRow{}, fmt.Errorf("building select query for %s: %w", organizationMembersTable, err)
	}

	return scanOrganizationMember(q.db.QueryRow(ctx, query, args...))
}

func (q *organizationMembers) Select(ctx context.Context) ([]repository.OrgMemberRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", organizationMembersTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", organizationMembersTable, err)
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

func (q *organizationMembers) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, err
	}
	sql := "SELECT EXISTS (" + subSQL + ")"

	var exists bool
	if err = q.db.QueryRow(ctx, sql, subArgs...).Scan(&exists); err != nil {
		return false, fmt.Errorf("sql=%s args=%v: %w", sql, subArgs, err)
	}
	return exists, nil
}

func (q *organizationMembers) UpdateOne(ctx context.Context) error {
	q.updater = q.updater.Set("replica_updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", organizationMembersTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing update query for %s: %w", organizationMembersTable, err)
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

func (q *organizationMembers) FilterByOrganizationID(organizationID ...uuid.UUID) repository.OrgMembersQ {
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

func (q *organizationMembers) UpdateVersion(v int32) repository.OrgMembersQ {
	q.updater = q.updater.Set("version", v)
	return q
}

func (q *organizationMembers) UpdateLabel(label *string) repository.OrgMembersQ {
	q.updater = q.updater.Set("label", label)
	return q
}

func (q *organizationMembers) UpdatePosition(position *string) repository.OrgMembersQ {
	q.updater = q.updater.Set("position", position)
	return q
}

func (q *organizationMembers) UpdateSourceUpdatedAt(updatedAt time.Time) repository.OrgMembersQ {
	q.updater = q.updater.Set("source_updated_at", updatedAt.UTC())
	return q
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
