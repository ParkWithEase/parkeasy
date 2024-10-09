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

// Car is an object representing the database table.
type Car struct {
	Carid        int64     `db:"carid,pk" `
	Userid       int64     `db:"userid" `
	Caruuid      uuid.UUID `db:"caruuid" `
	Licenseplate string    `db:"licenseplate" `
	Make         string    `db:"make" `
	Model        string    `db:"model" `
	Color        string    `db:"color" `

	R carR `db:"-" `
}

// CarSlice is an alias for a slice of pointers to Car.
// This should almost always be used instead of []*Car.
type CarSlice []*Car

// Cars contains methods to work with the car table
var Cars = psql.NewTablex[*Car, CarSlice, *CarSetter]("", "car")

// CarsQuery is a query on the car table
type CarsQuery = *psql.ViewQuery[*Car, CarSlice]

// CarsStmt is a prepared statment on car
type CarsStmt = bob.QueryStmt[*Car, CarSlice]

// carR is where relationships are stored.
type carR struct {
	UseridUser *User // car.car_userid_fkey
}

// CarSetter is used for insert/upsert/update operations
// All values are optional, and do not have to be set
// Generated columns are not included
type CarSetter struct {
	Carid        omit.Val[int64]     `db:"carid,pk" `
	Userid       omit.Val[int64]     `db:"userid" `
	Caruuid      omit.Val[uuid.UUID] `db:"caruuid" `
	Licenseplate omit.Val[string]    `db:"licenseplate" `
	Make         omit.Val[string]    `db:"make" `
	Model        omit.Val[string]    `db:"model" `
	Color        omit.Val[string]    `db:"color" `
}

func (s CarSetter) SetColumns() []string {
	vals := make([]string, 0, 7)
	if !s.Carid.IsUnset() {
		vals = append(vals, "carid")
	}

	if !s.Userid.IsUnset() {
		vals = append(vals, "userid")
	}

	if !s.Caruuid.IsUnset() {
		vals = append(vals, "caruuid")
	}

	if !s.Licenseplate.IsUnset() {
		vals = append(vals, "licenseplate")
	}

	if !s.Make.IsUnset() {
		vals = append(vals, "make")
	}

	if !s.Model.IsUnset() {
		vals = append(vals, "model")
	}

	if !s.Color.IsUnset() {
		vals = append(vals, "color")
	}

	return vals
}

func (s CarSetter) Overwrite(t *Car) {
	if !s.Carid.IsUnset() {
		t.Carid, _ = s.Carid.Get()
	}
	if !s.Userid.IsUnset() {
		t.Userid, _ = s.Userid.Get()
	}
	if !s.Caruuid.IsUnset() {
		t.Caruuid, _ = s.Caruuid.Get()
	}
	if !s.Licenseplate.IsUnset() {
		t.Licenseplate, _ = s.Licenseplate.Get()
	}
	if !s.Make.IsUnset() {
		t.Make, _ = s.Make.Get()
	}
	if !s.Model.IsUnset() {
		t.Model, _ = s.Model.Get()
	}
	if !s.Color.IsUnset() {
		t.Color, _ = s.Color.Get()
	}
}

func (s CarSetter) InsertMod() bob.Mod[*dialect.InsertQuery] {
	vals := make([]bob.Expression, 7)
	if s.Carid.IsUnset() {
		vals[0] = psql.Raw("DEFAULT")
	} else {
		vals[0] = psql.Arg(s.Carid)
	}

	if s.Userid.IsUnset() {
		vals[1] = psql.Raw("DEFAULT")
	} else {
		vals[1] = psql.Arg(s.Userid)
	}

	if s.Caruuid.IsUnset() {
		vals[2] = psql.Raw("DEFAULT")
	} else {
		vals[2] = psql.Arg(s.Caruuid)
	}

	if s.Licenseplate.IsUnset() {
		vals[3] = psql.Raw("DEFAULT")
	} else {
		vals[3] = psql.Arg(s.Licenseplate)
	}

	if s.Make.IsUnset() {
		vals[4] = psql.Raw("DEFAULT")
	} else {
		vals[4] = psql.Arg(s.Make)
	}

	if s.Model.IsUnset() {
		vals[5] = psql.Raw("DEFAULT")
	} else {
		vals[5] = psql.Arg(s.Model)
	}

	if s.Color.IsUnset() {
		vals[6] = psql.Raw("DEFAULT")
	} else {
		vals[6] = psql.Arg(s.Color)
	}

	return im.Values(vals...)
}

func (s CarSetter) Apply(q *dialect.UpdateQuery) {
	um.Set(s.Expressions()...).Apply(q)
}

func (s CarSetter) Expressions(prefix ...string) []bob.Expression {
	exprs := make([]bob.Expression, 0, 7)

	if !s.Carid.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "carid")...),
			psql.Arg(s.Carid),
		}})
	}

	if !s.Userid.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "userid")...),
			psql.Arg(s.Userid),
		}})
	}

	if !s.Caruuid.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "caruuid")...),
			psql.Arg(s.Caruuid),
		}})
	}

	if !s.Licenseplate.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "licenseplate")...),
			psql.Arg(s.Licenseplate),
		}})
	}

	if !s.Make.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "make")...),
			psql.Arg(s.Make),
		}})
	}

	if !s.Model.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "model")...),
			psql.Arg(s.Model),
		}})
	}

	if !s.Color.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "color")...),
			psql.Arg(s.Color),
		}})
	}

	return exprs
}

type carColumnNames struct {
	Carid        string
	Userid       string
	Caruuid      string
	Licenseplate string
	Make         string
	Model        string
	Color        string
}

var CarColumns = buildCarColumns("car")

type carColumns struct {
	tableAlias   string
	Carid        psql.Expression
	Userid       psql.Expression
	Caruuid      psql.Expression
	Licenseplate psql.Expression
	Make         psql.Expression
	Model        psql.Expression
	Color        psql.Expression
}

func (c carColumns) Alias() string {
	return c.tableAlias
}

func (carColumns) AliasedAs(alias string) carColumns {
	return buildCarColumns(alias)
}

func buildCarColumns(alias string) carColumns {
	return carColumns{
		tableAlias:   alias,
		Carid:        psql.Quote(alias, "carid"),
		Userid:       psql.Quote(alias, "userid"),
		Caruuid:      psql.Quote(alias, "caruuid"),
		Licenseplate: psql.Quote(alias, "licenseplate"),
		Make:         psql.Quote(alias, "make"),
		Model:        psql.Quote(alias, "model"),
		Color:        psql.Quote(alias, "color"),
	}
}

type carWhere[Q psql.Filterable] struct {
	Carid        psql.WhereMod[Q, int64]
	Userid       psql.WhereMod[Q, int64]
	Caruuid      psql.WhereMod[Q, uuid.UUID]
	Licenseplate psql.WhereMod[Q, string]
	Make         psql.WhereMod[Q, string]
	Model        psql.WhereMod[Q, string]
	Color        psql.WhereMod[Q, string]
}

func (carWhere[Q]) AliasedAs(alias string) carWhere[Q] {
	return buildCarWhere[Q](buildCarColumns(alias))
}

func buildCarWhere[Q psql.Filterable](cols carColumns) carWhere[Q] {
	return carWhere[Q]{
		Carid:        psql.Where[Q, int64](cols.Carid),
		Userid:       psql.Where[Q, int64](cols.Userid),
		Caruuid:      psql.Where[Q, uuid.UUID](cols.Caruuid),
		Licenseplate: psql.Where[Q, string](cols.Licenseplate),
		Make:         psql.Where[Q, string](cols.Make),
		Model:        psql.Where[Q, string](cols.Model),
		Color:        psql.Where[Q, string](cols.Color),
	}
}

type carJoins[Q dialect.Joinable] struct {
	typ        string
	UseridUser func(context.Context) modAs[Q, userColumns]
}

func (j carJoins[Q]) aliasedAs(alias string) carJoins[Q] {
	return buildCarJoins[Q](buildCarColumns(alias), j.typ)
}

func buildCarJoins[Q dialect.Joinable](cols carColumns, typ string) carJoins[Q] {
	return carJoins[Q]{
		typ:        typ,
		UseridUser: carsJoinUseridUser[Q](cols, typ),
	}
}

// FindCar retrieves a single record by primary key
// If cols is empty Find will return all columns.
func FindCar(ctx context.Context, exec bob.Executor, CaridPK int64, cols ...string) (*Car, error) {
	if len(cols) == 0 {
		return Cars.Query(
			ctx, exec,
			SelectWhere.Cars.Carid.EQ(CaridPK),
		).One()
	}

	return Cars.Query(
		ctx, exec,
		SelectWhere.Cars.Carid.EQ(CaridPK),
		sm.Columns(Cars.Columns().Only(cols...)),
	).One()
}

// CarExists checks the presence of a single record by primary key
func CarExists(ctx context.Context, exec bob.Executor, CaridPK int64) (bool, error) {
	return Cars.Query(
		ctx, exec,
		SelectWhere.Cars.Carid.EQ(CaridPK),
	).Exists()
}

// PrimaryKeyVals returns the primary key values of the Car
func (o *Car) PrimaryKeyVals() bob.Expression {
	return psql.Arg(o.Carid)
}

// Update uses an executor to update the Car
func (o *Car) Update(ctx context.Context, exec bob.Executor, s *CarSetter) error {
	return Cars.Update(ctx, exec, s, o)
}

// Delete deletes a single Car record with an executor
func (o *Car) Delete(ctx context.Context, exec bob.Executor) error {
	return Cars.Delete(ctx, exec, o)
}

// Reload refreshes the Car using the executor
func (o *Car) Reload(ctx context.Context, exec bob.Executor) error {
	o2, err := Cars.Query(
		ctx, exec,
		SelectWhere.Cars.Carid.EQ(o.Carid),
	).One()
	if err != nil {
		return err
	}
	o2.R = o.R
	*o = *o2

	return nil
}

func (o CarSlice) UpdateAll(ctx context.Context, exec bob.Executor, vals CarSetter) error {
	return Cars.Update(ctx, exec, &vals, o...)
}

func (o CarSlice) DeleteAll(ctx context.Context, exec bob.Executor) error {
	return Cars.Delete(ctx, exec, o...)
}

func (o CarSlice) ReloadAll(ctx context.Context, exec bob.Executor) error {
	var mods []bob.Mod[*dialect.SelectQuery]

	CaridPK := make([]int64, len(o))

	for i, o := range o {
		CaridPK[i] = o.Carid
	}

	mods = append(mods,
		SelectWhere.Cars.Carid.In(CaridPK...),
	)

	o2, err := Cars.Query(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, old := range o {
		for _, new := range o2 {
			if new.Carid != old.Carid {
				continue
			}
			new.R = old.R
			*old = *new
			break
		}
	}

	return nil
}

func carsJoinUseridUser[Q dialect.Joinable](from carColumns, typ string) func(context.Context) modAs[Q, userColumns] {
	return func(ctx context.Context) modAs[Q, userColumns] {
		return modAs[Q, userColumns]{
			c: UserColumns,
			f: func(to userColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, Users.Name(ctx).As(to.Alias())).On(
						to.Userid.EQ(from.Userid),
					))
				}

				return mods
			},
		}
	}
}

// UseridUser starts a query for related objects on users
func (o *Car) UseridUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) UsersQuery {
	return Users.Query(ctx, exec, append(mods,
		sm.Where(UserColumns.Userid.EQ(psql.Arg(o.Userid))),
	)...)
}

func (os CarSlice) UseridUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) UsersQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = psql.ArgGroup(o.Userid)
	}

	return Users.Query(ctx, exec, append(mods,
		sm.Where(psql.Group(UserColumns.Userid).In(PKArgs...)),
	)...)
}

func (o *Car) Preload(name string, retrieved any) error {
	if o == nil {
		return nil
	}

	switch name {
	case "UseridUser":
		rel, ok := retrieved.(*User)
		if !ok {
			return fmt.Errorf("car cannot load %T as %q", retrieved, name)
		}

		o.R.UseridUser = rel

		if rel != nil {
			rel.R.UseridCars = CarSlice{o}
		}
		return nil
	default:
		return fmt.Errorf("car has no relationship %q", name)
	}
}

func PreloadCarUseridUser(opts ...psql.PreloadOption) psql.Preloader {
	return psql.Preload[*User, UserSlice](orm.Relationship{
		Name: "UseridUser",
		Sides: []orm.RelSide{
			{
				From: "car",
				To:   TableNames.Users,
				ToExpr: func(ctx context.Context) bob.Expression {
					return Users.Name(ctx)
				},
				FromColumns: []string{
					ColumnNames.Cars.Userid,
				},
				ToColumns: []string{
					ColumnNames.Users.Userid,
				},
			},
		},
	}, Users.Columns().Names(), opts...)
}

func ThenLoadCarUseridUser(queryMods ...bob.Mod[*dialect.SelectQuery]) psql.Loader {
	return psql.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadCarUseridUser(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load CarUseridUser", retrieved)
		}

		err := loader.LoadCarUseridUser(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadCarUseridUser loads the car's UseridUser into the .R struct
func (o *Car) LoadCarUseridUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.UseridUser = nil

	related, err := o.UseridUser(ctx, exec, mods...).One()
	if err != nil {
		return err
	}

	related.R.UseridCars = CarSlice{o}

	o.R.UseridUser = related
	return nil
}

// LoadCarUseridUser loads the car's UseridUser into the .R struct
func (os CarSlice) LoadCarUseridUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	users, err := os.UseridUser(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, o := range os {
		for _, rel := range users {
			if o.Userid != rel.Userid {
				continue
			}

			rel.R.UseridCars = append(rel.R.UseridCars, o)

			o.R.UseridUser = rel
			break
		}
	}

	return nil
}

func attachCarUseridUser0(ctx context.Context, exec bob.Executor, count int, car0 *Car, user1 *User) (*Car, error) {
	setter := &CarSetter{
		Userid: omit.From(user1.Userid),
	}

	err := Cars.Update(ctx, exec, setter, car0)
	if err != nil {
		return nil, fmt.Errorf("attachCarUseridUser0: %w", err)
	}

	return car0, nil
}

func (car0 *Car) InsertUseridUser(ctx context.Context, exec bob.Executor, related *UserSetter) error {
	user1, err := Users.Insert(ctx, exec, related)
	if err != nil {
		return fmt.Errorf("inserting related objects: %w", err)
	}

	_, err = attachCarUseridUser0(ctx, exec, 1, car0, user1)
	if err != nil {
		return err
	}

	car0.R.UseridUser = user1

	user1.R.UseridCars = append(user1.R.UseridCars, car0)

	return nil
}

func (car0 *Car) AttachUseridUser(ctx context.Context, exec bob.Executor, user1 *User) error {
	var err error

	_, err = attachCarUseridUser0(ctx, exec, 1, car0, user1)
	if err != nil {
		return err
	}

	car0.R.UseridUser = user1

	user1.R.UseridCars = append(user1.R.UseridCars, car0)

	return nil
}