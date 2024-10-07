// Code generated by BobGen psql v0.28.1. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package dbmodels

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/bob/expr"
	"github.com/stephenafamo/bob/mods"
	"github.com/stephenafamo/bob/orm"
)

// Auth is an object representing the database table.
type Auth struct {
	Authid       int32     `db:"authid,pk" `
	Authuuid     uuid.UUID `db:"authuuid" `
	Email        string    `db:"email" `
	Passwordhash string    `db:"passwordhash" `

	R authR `db:"-" `
}

// AuthSlice is an alias for a slice of pointers to Auth.
// This should almost always be used instead of []*Auth.
type AuthSlice []*Auth

// Auths contains methods to work with the auth table
var Auths = psql.NewTablex[*Auth, AuthSlice, *AuthSetter]("", "auth")

// AuthsQuery is a query on the auth table
type AuthsQuery = *psql.ViewQuery[*Auth, AuthSlice]

// AuthsStmt is a prepared statment on auth
type AuthsStmt = bob.QueryStmt[*Auth, AuthSlice]

// authR is where relationships are stored.
type authR struct {
	AuthuuidUser *User // users.users_authuuid_fkey
}

// AuthSetter is used for insert/upsert/update operations
// All values are optional, and do not have to be set
// Generated columns are not included
type AuthSetter struct {
	Authid       omit.Val[int32]     `db:"authid,pk" `
	Authuuid     omit.Val[uuid.UUID] `db:"authuuid" `
	Email        omit.Val[string]    `db:"email" `
	Passwordhash omit.Val[string]    `db:"passwordhash" `
}

func (s AuthSetter) SetColumns() []string {
	vals := make([]string, 0, 4)
	if !s.Authid.IsUnset() {
		vals = append(vals, "authid")
	}

	if !s.Authuuid.IsUnset() {
		vals = append(vals, "authuuid")
	}

	if !s.Email.IsUnset() {
		vals = append(vals, "email")
	}

	if !s.Passwordhash.IsUnset() {
		vals = append(vals, "passwordhash")
	}

	return vals
}

func (s AuthSetter) Overwrite(t *Auth) {
	if !s.Authid.IsUnset() {
		t.Authid, _ = s.Authid.Get()
	}
	if !s.Authuuid.IsUnset() {
		t.Authuuid, _ = s.Authuuid.Get()
	}
	if !s.Email.IsUnset() {
		t.Email, _ = s.Email.Get()
	}
	if !s.Passwordhash.IsUnset() {
		t.Passwordhash, _ = s.Passwordhash.Get()
	}
}

func (s AuthSetter) InsertMod() bob.Mod[*dialect.InsertQuery] {
	vals := make([]bob.Expression, 4)
	if s.Authid.IsUnset() {
		vals[0] = psql.Raw("DEFAULT")
	} else {
		vals[0] = psql.Arg(s.Authid)
	}

	if s.Authuuid.IsUnset() {
		vals[1] = psql.Raw("DEFAULT")
	} else {
		vals[1] = psql.Arg(s.Authuuid)
	}

	if s.Email.IsUnset() {
		vals[2] = psql.Raw("DEFAULT")
	} else {
		vals[2] = psql.Arg(s.Email)
	}

	if s.Passwordhash.IsUnset() {
		vals[3] = psql.Raw("DEFAULT")
	} else {
		vals[3] = psql.Arg(s.Passwordhash)
	}

	return im.Values(vals...)
}

func (s AuthSetter) Apply(q *dialect.UpdateQuery) {
	um.Set(s.Expressions()...).Apply(q)
}

func (s AuthSetter) Expressions(prefix ...string) []bob.Expression {
	exprs := make([]bob.Expression, 0, 4)

	if !s.Authid.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "authid")...),
			psql.Arg(s.Authid),
		}})
	}

	if !s.Authuuid.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "authuuid")...),
			psql.Arg(s.Authuuid),
		}})
	}

	if !s.Email.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "email")...),
			psql.Arg(s.Email),
		}})
	}

	if !s.Passwordhash.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "passwordhash")...),
			psql.Arg(s.Passwordhash),
		}})
	}

	return exprs
}

type authColumnNames struct {
	Authid       string
	Authuuid     string
	Email        string
	Passwordhash string
}

var AuthColumns = buildAuthColumns("auth")

type authColumns struct {
	tableAlias   string
	Authid       psql.Expression
	Authuuid     psql.Expression
	Email        psql.Expression
	Passwordhash psql.Expression
}

func (c authColumns) Alias() string {
	return c.tableAlias
}

func (authColumns) AliasedAs(alias string) authColumns {
	return buildAuthColumns(alias)
}

func buildAuthColumns(alias string) authColumns {
	return authColumns{
		tableAlias:   alias,
		Authid:       psql.Quote(alias, "authid"),
		Authuuid:     psql.Quote(alias, "authuuid"),
		Email:        psql.Quote(alias, "email"),
		Passwordhash: psql.Quote(alias, "passwordhash"),
	}
}

type authWhere[Q psql.Filterable] struct {
	Authid       psql.WhereMod[Q, int32]
	Authuuid     psql.WhereMod[Q, uuid.UUID]
	Email        psql.WhereMod[Q, string]
	Passwordhash psql.WhereMod[Q, string]
}

func (authWhere[Q]) AliasedAs(alias string) authWhere[Q] {
	return buildAuthWhere[Q](buildAuthColumns(alias))
}

func buildAuthWhere[Q psql.Filterable](cols authColumns) authWhere[Q] {
	return authWhere[Q]{
		Authid:       psql.Where[Q, int32](cols.Authid),
		Authuuid:     psql.Where[Q, uuid.UUID](cols.Authuuid),
		Email:        psql.Where[Q, string](cols.Email),
		Passwordhash: psql.Where[Q, string](cols.Passwordhash),
	}
}

type authJoins[Q dialect.Joinable] struct {
	typ          string
	AuthuuidUser func(context.Context) modAs[Q, userColumns]
}

func (j authJoins[Q]) aliasedAs(alias string) authJoins[Q] {
	return buildAuthJoins[Q](buildAuthColumns(alias), j.typ)
}

func buildAuthJoins[Q dialect.Joinable](cols authColumns, typ string) authJoins[Q] {
	return authJoins[Q]{
		typ:          typ,
		AuthuuidUser: authsJoinAuthuuidUser[Q](cols, typ),
	}
}

// FindAuth retrieves a single record by primary key
// If cols is empty Find will return all columns.
func FindAuth(ctx context.Context, exec bob.Executor, AuthidPK int32, cols ...string) (*Auth, error) {
	if len(cols) == 0 {
		return Auths.Query(
			ctx, exec,
			SelectWhere.Auths.Authid.EQ(AuthidPK),
		).One()
	}

	return Auths.Query(
		ctx, exec,
		SelectWhere.Auths.Authid.EQ(AuthidPK),
		sm.Columns(Auths.Columns().Only(cols...)),
	).One()
}

// AuthExists checks the presence of a single record by primary key
func AuthExists(ctx context.Context, exec bob.Executor, AuthidPK int32) (bool, error) {
	return Auths.Query(
		ctx, exec,
		SelectWhere.Auths.Authid.EQ(AuthidPK),
	).Exists()
}

// PrimaryKeyVals returns the primary key values of the Auth
func (o *Auth) PrimaryKeyVals() bob.Expression {
	return psql.Arg(o.Authid)
}

// Update uses an executor to update the Auth
func (o *Auth) Update(ctx context.Context, exec bob.Executor, s *AuthSetter) error {
	return Auths.Update(ctx, exec, s, o)
}

// Delete deletes a single Auth record with an executor
func (o *Auth) Delete(ctx context.Context, exec bob.Executor) error {
	return Auths.Delete(ctx, exec, o)
}

// Reload refreshes the Auth using the executor
func (o *Auth) Reload(ctx context.Context, exec bob.Executor) error {
	o2, err := Auths.Query(
		ctx, exec,
		SelectWhere.Auths.Authid.EQ(o.Authid),
	).One()
	if err != nil {
		return err
	}
	o2.R = o.R
	*o = *o2

	return nil
}

func (o AuthSlice) UpdateAll(ctx context.Context, exec bob.Executor, vals AuthSetter) error {
	return Auths.Update(ctx, exec, &vals, o...)
}

func (o AuthSlice) DeleteAll(ctx context.Context, exec bob.Executor) error {
	return Auths.Delete(ctx, exec, o...)
}

func (o AuthSlice) ReloadAll(ctx context.Context, exec bob.Executor) error {
	var mods []bob.Mod[*dialect.SelectQuery]

	AuthidPK := make([]int32, len(o))

	for i, o := range o {
		AuthidPK[i] = o.Authid
	}

	mods = append(mods,
		SelectWhere.Auths.Authid.In(AuthidPK...),
	)

	o2, err := Auths.Query(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, old := range o {
		for _, new := range o2 {
			if new.Authid != old.Authid {
				continue
			}
			new.R = old.R
			*old = *new
			break
		}
	}

	return nil
}

func authsJoinAuthuuidUser[Q dialect.Joinable](from authColumns, typ string) func(context.Context) modAs[Q, userColumns] {
	return func(ctx context.Context) modAs[Q, userColumns] {
		return modAs[Q, userColumns]{
			c: UserColumns,
			f: func(to userColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, Users.Name(ctx).As(to.Alias())).On(
						to.Authuuid.EQ(from.Authuuid),
					))
				}

				return mods
			},
		}
	}
}

// AuthuuidUser starts a query for related objects on users
func (o *Auth) AuthuuidUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) UsersQuery {
	return Users.Query(ctx, exec, append(mods,
		sm.Where(UserColumns.Authuuid.EQ(psql.Arg(o.Authuuid))),
	)...)
}

func (os AuthSlice) AuthuuidUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) UsersQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = psql.ArgGroup(o.Authuuid)
	}

	return Users.Query(ctx, exec, append(mods,
		sm.Where(psql.Group(UserColumns.Authuuid).In(PKArgs...)),
	)...)
}

func (o *Auth) Preload(name string, retrieved any) error {
	if o == nil {
		return nil
	}

	switch name {
	case "AuthuuidUser":
		rel, ok := retrieved.(*User)
		if !ok {
			return fmt.Errorf("auth cannot load %T as %q", retrieved, name)
		}

		o.R.AuthuuidUser = rel

		if rel != nil {
			rel.R.AuthuuidAuth = o
		}
		return nil
	default:
		return fmt.Errorf("auth has no relationship %q", name)
	}
}

func PreloadAuthAuthuuidUser(opts ...psql.PreloadOption) psql.Preloader {
	return psql.Preload[*User, UserSlice](orm.Relationship{
		Name: "AuthuuidUser",
		Sides: []orm.RelSide{
			{
				From: "auth",
				To:   TableNames.Users,
				ToExpr: func(ctx context.Context) bob.Expression {
					return Users.Name(ctx)
				},
				FromColumns: []string{
					ColumnNames.Auths.Authuuid,
				},
				ToColumns: []string{
					ColumnNames.Users.Authuuid,
				},
			},
		},
	}, Users.Columns().Names(), opts...)
}

func ThenLoadAuthAuthuuidUser(queryMods ...bob.Mod[*dialect.SelectQuery]) psql.Loader {
	return psql.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadAuthAuthuuidUser(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load AuthAuthuuidUser", retrieved)
		}

		err := loader.LoadAuthAuthuuidUser(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadAuthAuthuuidUser loads the auth's AuthuuidUser into the .R struct
func (o *Auth) LoadAuthAuthuuidUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.AuthuuidUser = nil

	related, err := o.AuthuuidUser(ctx, exec, mods...).One()
	if err != nil {
		return err
	}

	related.R.AuthuuidAuth = o

	o.R.AuthuuidUser = related
	return nil
}

// LoadAuthAuthuuidUser loads the auth's AuthuuidUser into the .R struct
func (os AuthSlice) LoadAuthAuthuuidUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	users, err := os.AuthuuidUser(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, o := range os {
		for _, rel := range users {
			if o.Authuuid != rel.Authuuid {
				continue
			}

			rel.R.AuthuuidAuth = o

			o.R.AuthuuidUser = rel
			break
		}
	}

	return nil
}

func insertAuthAuthuuidUser0(ctx context.Context, exec bob.Executor, user1 *UserSetter, auth0 *Auth) (*User, error) {
	user1.Authuuid = omit.From(auth0.Authuuid)

	ret, err := Users.Insert(ctx, exec, user1)
	if err != nil {
		return ret, fmt.Errorf("insertAuthAuthuuidUser0: %w", err)
	}

	return ret, nil
}

func attachAuthAuthuuidUser0(ctx context.Context, exec bob.Executor, count int, user1 *User, auth0 *Auth) (*User, error) {
	setter := &UserSetter{
		Authuuid: omit.From(auth0.Authuuid),
	}

	err := Users.Update(ctx, exec, setter, user1)
	if err != nil {
		return nil, fmt.Errorf("attachAuthAuthuuidUser0: %w", err)
	}

	return user1, nil
}

func (auth0 *Auth) InsertAuthuuidUser(ctx context.Context, exec bob.Executor, related *UserSetter) error {
	user1, err := insertAuthAuthuuidUser0(ctx, exec, related, auth0)
	if err != nil {
		return err
	}

	auth0.R.AuthuuidUser = user1

	user1.R.AuthuuidAuth = auth0

	return nil
}

func (auth0 *Auth) AttachAuthuuidUser(ctx context.Context, exec bob.Executor, user1 *User) error {
	var err error

	_, err = attachAuthAuthuuidUser0(ctx, exec, 1, user1, auth0)
	if err != nil {
		return err
	}

	auth0.R.AuthuuidUser = user1

	user1.R.AuthuuidAuth = auth0

	return nil
}
