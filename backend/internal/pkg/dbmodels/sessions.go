// Code generated by modelgen. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package dbmodels

import (
	"context"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/bob/expr"
)

// Session is an object representing the database table.
type Session struct {
	Token  string    `db:"token,pk" `
	Data   []byte    `db:"data" `
	Expiry time.Time `db:"expiry" `
}

// SessionSlice is an alias for a slice of pointers to Session.
// This should almost always be used instead of []*Session.
type SessionSlice []*Session

// Sessions contains methods to work with the sessions table
var Sessions = psql.NewTablex[*Session, SessionSlice, *SessionSetter]("", "sessions")

// SessionsQuery is a query on the sessions table
type SessionsQuery = *psql.ViewQuery[*Session, SessionSlice]

// SessionsStmt is a prepared statment on sessions
type SessionsStmt = bob.QueryStmt[*Session, SessionSlice]

// SessionSetter is used for insert/upsert/update operations
// All values are optional, and do not have to be set
// Generated columns are not included
type SessionSetter struct {
	Token  omit.Val[string]    `db:"token,pk" `
	Data   omit.Val[[]byte]    `db:"data" `
	Expiry omit.Val[time.Time] `db:"expiry" `
}

func (s SessionSetter) SetColumns() []string {
	vals := make([]string, 0, 3)
	if !s.Token.IsUnset() {
		vals = append(vals, "token")
	}

	if !s.Data.IsUnset() {
		vals = append(vals, "data")
	}

	if !s.Expiry.IsUnset() {
		vals = append(vals, "expiry")
	}

	return vals
}

func (s SessionSetter) Overwrite(t *Session) {
	if !s.Token.IsUnset() {
		t.Token, _ = s.Token.Get()
	}
	if !s.Data.IsUnset() {
		t.Data, _ = s.Data.Get()
	}
	if !s.Expiry.IsUnset() {
		t.Expiry, _ = s.Expiry.Get()
	}
}

func (s SessionSetter) InsertMod() bob.Mod[*dialect.InsertQuery] {
	vals := make([]bob.Expression, 3)
	if s.Token.IsUnset() {
		vals[0] = psql.Raw("DEFAULT")
	} else {
		vals[0] = psql.Arg(s.Token)
	}

	if s.Data.IsUnset() {
		vals[1] = psql.Raw("DEFAULT")
	} else {
		vals[1] = psql.Arg(s.Data)
	}

	if s.Expiry.IsUnset() {
		vals[2] = psql.Raw("DEFAULT")
	} else {
		vals[2] = psql.Arg(s.Expiry)
	}

	return im.Values(vals...)
}

func (s SessionSetter) Apply(q *dialect.UpdateQuery) {
	um.Set(s.Expressions()...).Apply(q)
}

func (s SessionSetter) Expressions(prefix ...string) []bob.Expression {
	exprs := make([]bob.Expression, 0, 3)

	if !s.Token.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "token")...),
			psql.Arg(s.Token),
		}})
	}

	if !s.Data.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "data")...),
			psql.Arg(s.Data),
		}})
	}

	if !s.Expiry.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "expiry")...),
			psql.Arg(s.Expiry),
		}})
	}

	return exprs
}

type sessionColumnNames struct {
	Token  string
	Data   string
	Expiry string
}

var SessionColumns = buildSessionColumns("sessions")

type sessionColumns struct {
	tableAlias string
	Token      psql.Expression
	Data       psql.Expression
	Expiry     psql.Expression
}

func (c sessionColumns) Alias() string {
	return c.tableAlias
}

func (sessionColumns) AliasedAs(alias string) sessionColumns {
	return buildSessionColumns(alias)
}

func buildSessionColumns(alias string) sessionColumns {
	return sessionColumns{
		tableAlias: alias,
		Token:      psql.Quote(alias, "token"),
		Data:       psql.Quote(alias, "data"),
		Expiry:     psql.Quote(alias, "expiry"),
	}
}

type sessionWhere[Q psql.Filterable] struct {
	Token  psql.WhereMod[Q, string]
	Data   psql.WhereMod[Q, []byte]
	Expiry psql.WhereMod[Q, time.Time]
}

func (sessionWhere[Q]) AliasedAs(alias string) sessionWhere[Q] {
	return buildSessionWhere[Q](buildSessionColumns(alias))
}

func buildSessionWhere[Q psql.Filterable](cols sessionColumns) sessionWhere[Q] {
	return sessionWhere[Q]{
		Token:  psql.Where[Q, string](cols.Token),
		Data:   psql.Where[Q, []byte](cols.Data),
		Expiry: psql.Where[Q, time.Time](cols.Expiry),
	}
}

// FindSession retrieves a single record by primary key
// If cols is empty Find will return all columns.
func FindSession(ctx context.Context, exec bob.Executor, TokenPK string, cols ...string) (*Session, error) {
	if len(cols) == 0 {
		return Sessions.Query(
			ctx, exec,
			SelectWhere.Sessions.Token.EQ(TokenPK),
		).One()
	}

	return Sessions.Query(
		ctx, exec,
		SelectWhere.Sessions.Token.EQ(TokenPK),
		sm.Columns(Sessions.Columns().Only(cols...)),
	).One()
}

// SessionExists checks the presence of a single record by primary key
func SessionExists(ctx context.Context, exec bob.Executor, TokenPK string) (bool, error) {
	return Sessions.Query(
		ctx, exec,
		SelectWhere.Sessions.Token.EQ(TokenPK),
	).Exists()
}

// PrimaryKeyVals returns the primary key values of the Session
func (o *Session) PrimaryKeyVals() bob.Expression {
	return psql.Arg(o.Token)
}

// Update uses an executor to update the Session
func (o *Session) Update(ctx context.Context, exec bob.Executor, s *SessionSetter) error {
	return Sessions.Update(ctx, exec, s, o)
}

// Delete deletes a single Session record with an executor
func (o *Session) Delete(ctx context.Context, exec bob.Executor) error {
	return Sessions.Delete(ctx, exec, o)
}

// Reload refreshes the Session using the executor
func (o *Session) Reload(ctx context.Context, exec bob.Executor) error {
	o2, err := Sessions.Query(
		ctx, exec,
		SelectWhere.Sessions.Token.EQ(o.Token),
	).One()
	if err != nil {
		return err
	}

	*o = *o2

	return nil
}

func (o SessionSlice) UpdateAll(ctx context.Context, exec bob.Executor, vals SessionSetter) error {
	return Sessions.Update(ctx, exec, &vals, o...)
}

func (o SessionSlice) DeleteAll(ctx context.Context, exec bob.Executor) error {
	return Sessions.Delete(ctx, exec, o...)
}

func (o SessionSlice) ReloadAll(ctx context.Context, exec bob.Executor) error {
	var mods []bob.Mod[*dialect.SelectQuery]

	TokenPK := make([]string, len(o))

	for i, o := range o {
		TokenPK[i] = o.Token
	}

	mods = append(mods,
		SelectWhere.Sessions.Token.In(TokenPK...),
	)

	o2, err := Sessions.Query(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, old := range o {
		for _, new := range o2 {
			if new.Token != old.Token {
				continue
			}

			*old = *new
			break
		}
	}

	return nil
}
