// Code generated by modelgen. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package dbmodels

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
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

// Parkingspot is an object representing the database table.
type Parkingspot struct {
	Parkingspotid      int64           `db:"parkingspotid,pk" `
	Userid             int64           `db:"userid" `
	Parkingspotuuid    uuid.UUID       `db:"parkingspotuuid" `
	Postalcode         string          `db:"postalcode" `
	Countrycode        string          `db:"countrycode" `
	City               string          `db:"city" `
	State              string          `db:"state" `
	Streetaddress      string          `db:"streetaddress" `
	Longitude          decimal.Decimal `db:"longitude" `
	Latitude           decimal.Decimal `db:"latitude" `
	Hasshelter         bool            `db:"hasshelter" `
	Hasplugin          bool            `db:"hasplugin" `
	Haschargingstation bool            `db:"haschargingstation" `
	Priceperhour       decimal.Decimal `db:"priceperhour" `

	R parkingspotR `db:"-" `
}

// ParkingspotSlice is an alias for a slice of pointers to Parkingspot.
// This should almost always be used instead of []*Parkingspot.
type ParkingspotSlice []*Parkingspot

// Parkingspots contains methods to work with the parkingspot table
var Parkingspots = psql.NewTablex[*Parkingspot, ParkingspotSlice, *ParkingspotSetter]("", "parkingspot")

// ParkingspotsQuery is a query on the parkingspot table
type ParkingspotsQuery = *psql.ViewQuery[*Parkingspot, ParkingspotSlice]

// ParkingspotsStmt is a prepared statment on parkingspot
type ParkingspotsStmt = bob.QueryStmt[*Parkingspot, ParkingspotSlice]

// parkingspotR is where relationships are stored.
type parkingspotR struct {
	UseridUser             *User         // parkingspot.parkingspot_userid_fkey
	ParkingspotidTimeunits TimeunitSlice // timeunit.timeunit_parkingspotid_fkey
}

// ParkingspotSetter is used for insert/upsert/update operations
// All values are optional, and do not have to be set
// Generated columns are not included
type ParkingspotSetter struct {
	Parkingspotid      omit.Val[int64]           `db:"parkingspotid,pk" `
	Userid             omit.Val[int64]           `db:"userid" `
	Parkingspotuuid    omit.Val[uuid.UUID]       `db:"parkingspotuuid" `
	Postalcode         omit.Val[string]          `db:"postalcode" `
	Countrycode        omit.Val[string]          `db:"countrycode" `
	City               omit.Val[string]          `db:"city" `
	State              omit.Val[string]          `db:"state" `
	Streetaddress      omit.Val[string]          `db:"streetaddress" `
	Longitude          omit.Val[decimal.Decimal] `db:"longitude" `
	Latitude           omit.Val[decimal.Decimal] `db:"latitude" `
	Hasshelter         omit.Val[bool]            `db:"hasshelter" `
	Hasplugin          omit.Val[bool]            `db:"hasplugin" `
	Haschargingstation omit.Val[bool]            `db:"haschargingstation" `
	Priceperhour       omit.Val[decimal.Decimal] `db:"priceperhour" `
}

func (s ParkingspotSetter) SetColumns() []string {
	vals := make([]string, 0, 14)
	if !s.Parkingspotid.IsUnset() {
		vals = append(vals, "parkingspotid")
	}

	if !s.Userid.IsUnset() {
		vals = append(vals, "userid")
	}

	if !s.Parkingspotuuid.IsUnset() {
		vals = append(vals, "parkingspotuuid")
	}

	if !s.Postalcode.IsUnset() {
		vals = append(vals, "postalcode")
	}

	if !s.Countrycode.IsUnset() {
		vals = append(vals, "countrycode")
	}

	if !s.City.IsUnset() {
		vals = append(vals, "city")
	}

	if !s.State.IsUnset() {
		vals = append(vals, "state")
	}

	if !s.Streetaddress.IsUnset() {
		vals = append(vals, "streetaddress")
	}

	if !s.Longitude.IsUnset() {
		vals = append(vals, "longitude")
	}

	if !s.Latitude.IsUnset() {
		vals = append(vals, "latitude")
	}

	if !s.Hasshelter.IsUnset() {
		vals = append(vals, "hasshelter")
	}

	if !s.Hasplugin.IsUnset() {
		vals = append(vals, "hasplugin")
	}

	if !s.Haschargingstation.IsUnset() {
		vals = append(vals, "haschargingstation")
	}

	if !s.Priceperhour.IsUnset() {
		vals = append(vals, "priceperhour")
	}

	return vals
}

func (s ParkingspotSetter) Overwrite(t *Parkingspot) {
	if !s.Parkingspotid.IsUnset() {
		t.Parkingspotid, _ = s.Parkingspotid.Get()
	}
	if !s.Userid.IsUnset() {
		t.Userid, _ = s.Userid.Get()
	}
	if !s.Parkingspotuuid.IsUnset() {
		t.Parkingspotuuid, _ = s.Parkingspotuuid.Get()
	}
	if !s.Postalcode.IsUnset() {
		t.Postalcode, _ = s.Postalcode.Get()
	}
	if !s.Countrycode.IsUnset() {
		t.Countrycode, _ = s.Countrycode.Get()
	}
	if !s.City.IsUnset() {
		t.City, _ = s.City.Get()
	}
	if !s.State.IsUnset() {
		t.State, _ = s.State.Get()
	}
	if !s.Streetaddress.IsUnset() {
		t.Streetaddress, _ = s.Streetaddress.Get()
	}
	if !s.Longitude.IsUnset() {
		t.Longitude, _ = s.Longitude.Get()
	}
	if !s.Latitude.IsUnset() {
		t.Latitude, _ = s.Latitude.Get()
	}
	if !s.Hasshelter.IsUnset() {
		t.Hasshelter, _ = s.Hasshelter.Get()
	}
	if !s.Hasplugin.IsUnset() {
		t.Hasplugin, _ = s.Hasplugin.Get()
	}
	if !s.Haschargingstation.IsUnset() {
		t.Haschargingstation, _ = s.Haschargingstation.Get()
	}
	if !s.Priceperhour.IsUnset() {
		t.Priceperhour, _ = s.Priceperhour.Get()
	}
}

func (s ParkingspotSetter) InsertMod() bob.Mod[*dialect.InsertQuery] {
	vals := make([]bob.Expression, 14)
	if s.Parkingspotid.IsUnset() {
		vals[0] = psql.Raw("DEFAULT")
	} else {
		vals[0] = psql.Arg(s.Parkingspotid)
	}

	if s.Userid.IsUnset() {
		vals[1] = psql.Raw("DEFAULT")
	} else {
		vals[1] = psql.Arg(s.Userid)
	}

	if s.Parkingspotuuid.IsUnset() {
		vals[2] = psql.Raw("DEFAULT")
	} else {
		vals[2] = psql.Arg(s.Parkingspotuuid)
	}

	if s.Postalcode.IsUnset() {
		vals[3] = psql.Raw("DEFAULT")
	} else {
		vals[3] = psql.Arg(s.Postalcode)
	}

	if s.Countrycode.IsUnset() {
		vals[4] = psql.Raw("DEFAULT")
	} else {
		vals[4] = psql.Arg(s.Countrycode)
	}

	if s.City.IsUnset() {
		vals[5] = psql.Raw("DEFAULT")
	} else {
		vals[5] = psql.Arg(s.City)
	}

	if s.State.IsUnset() {
		vals[6] = psql.Raw("DEFAULT")
	} else {
		vals[6] = psql.Arg(s.State)
	}

	if s.Streetaddress.IsUnset() {
		vals[7] = psql.Raw("DEFAULT")
	} else {
		vals[7] = psql.Arg(s.Streetaddress)
	}

	if s.Longitude.IsUnset() {
		vals[8] = psql.Raw("DEFAULT")
	} else {
		vals[8] = psql.Arg(s.Longitude)
	}

	if s.Latitude.IsUnset() {
		vals[9] = psql.Raw("DEFAULT")
	} else {
		vals[9] = psql.Arg(s.Latitude)
	}

	if s.Hasshelter.IsUnset() {
		vals[10] = psql.Raw("DEFAULT")
	} else {
		vals[10] = psql.Arg(s.Hasshelter)
	}

	if s.Hasplugin.IsUnset() {
		vals[11] = psql.Raw("DEFAULT")
	} else {
		vals[11] = psql.Arg(s.Hasplugin)
	}

	if s.Haschargingstation.IsUnset() {
		vals[12] = psql.Raw("DEFAULT")
	} else {
		vals[12] = psql.Arg(s.Haschargingstation)
	}

	if s.Priceperhour.IsUnset() {
		vals[13] = psql.Raw("DEFAULT")
	} else {
		vals[13] = psql.Arg(s.Priceperhour)
	}

	return im.Values(vals...)
}

func (s ParkingspotSetter) Apply(q *dialect.UpdateQuery) {
	um.Set(s.Expressions()...).Apply(q)
}

func (s ParkingspotSetter) Expressions(prefix ...string) []bob.Expression {
	exprs := make([]bob.Expression, 0, 14)

	if !s.Parkingspotid.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "parkingspotid")...),
			psql.Arg(s.Parkingspotid),
		}})
	}

	if !s.Userid.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "userid")...),
			psql.Arg(s.Userid),
		}})
	}

	if !s.Parkingspotuuid.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "parkingspotuuid")...),
			psql.Arg(s.Parkingspotuuid),
		}})
	}

	if !s.Postalcode.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "postalcode")...),
			psql.Arg(s.Postalcode),
		}})
	}

	if !s.Countrycode.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "countrycode")...),
			psql.Arg(s.Countrycode),
		}})
	}

	if !s.City.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "city")...),
			psql.Arg(s.City),
		}})
	}

	if !s.State.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "state")...),
			psql.Arg(s.State),
		}})
	}

	if !s.Streetaddress.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "streetaddress")...),
			psql.Arg(s.Streetaddress),
		}})
	}

	if !s.Longitude.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "longitude")...),
			psql.Arg(s.Longitude),
		}})
	}

	if !s.Latitude.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "latitude")...),
			psql.Arg(s.Latitude),
		}})
	}

	if !s.Hasshelter.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "hasshelter")...),
			psql.Arg(s.Hasshelter),
		}})
	}

	if !s.Hasplugin.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "hasplugin")...),
			psql.Arg(s.Hasplugin),
		}})
	}

	if !s.Haschargingstation.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "haschargingstation")...),
			psql.Arg(s.Haschargingstation),
		}})
	}

	if !s.Priceperhour.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "priceperhour")...),
			psql.Arg(s.Priceperhour),
		}})
	}

	return exprs
}

type parkingspotColumnNames struct {
	Parkingspotid      string
	Userid             string
	Parkingspotuuid    string
	Postalcode         string
	Countrycode        string
	City               string
	State              string
	Streetaddress      string
	Longitude          string
	Latitude           string
	Hasshelter         string
	Hasplugin          string
	Haschargingstation string
	Priceperhour       string
}

var ParkingspotColumns = buildParkingspotColumns("parkingspot")

type parkingspotColumns struct {
	tableAlias         string
	Parkingspotid      psql.Expression
	Userid             psql.Expression
	Parkingspotuuid    psql.Expression
	Postalcode         psql.Expression
	Countrycode        psql.Expression
	City               psql.Expression
	State              psql.Expression
	Streetaddress      psql.Expression
	Longitude          psql.Expression
	Latitude           psql.Expression
	Hasshelter         psql.Expression
	Hasplugin          psql.Expression
	Haschargingstation psql.Expression
	Priceperhour       psql.Expression
}

func (c parkingspotColumns) Alias() string {
	return c.tableAlias
}

func (parkingspotColumns) AliasedAs(alias string) parkingspotColumns {
	return buildParkingspotColumns(alias)
}

func buildParkingspotColumns(alias string) parkingspotColumns {
	return parkingspotColumns{
		tableAlias:         alias,
		Parkingspotid:      psql.Quote(alias, "parkingspotid"),
		Userid:             psql.Quote(alias, "userid"),
		Parkingspotuuid:    psql.Quote(alias, "parkingspotuuid"),
		Postalcode:         psql.Quote(alias, "postalcode"),
		Countrycode:        psql.Quote(alias, "countrycode"),
		City:               psql.Quote(alias, "city"),
		State:              psql.Quote(alias, "state"),
		Streetaddress:      psql.Quote(alias, "streetaddress"),
		Longitude:          psql.Quote(alias, "longitude"),
		Latitude:           psql.Quote(alias, "latitude"),
		Hasshelter:         psql.Quote(alias, "hasshelter"),
		Hasplugin:          psql.Quote(alias, "hasplugin"),
		Haschargingstation: psql.Quote(alias, "haschargingstation"),
		Priceperhour:       psql.Quote(alias, "priceperhour"),
	}
}

type parkingspotWhere[Q psql.Filterable] struct {
	Parkingspotid      psql.WhereMod[Q, int64]
	Userid             psql.WhereMod[Q, int64]
	Parkingspotuuid    psql.WhereMod[Q, uuid.UUID]
	Postalcode         psql.WhereMod[Q, string]
	Countrycode        psql.WhereMod[Q, string]
	City               psql.WhereMod[Q, string]
	State              psql.WhereMod[Q, string]
	Streetaddress      psql.WhereMod[Q, string]
	Longitude          psql.WhereMod[Q, decimal.Decimal]
	Latitude           psql.WhereMod[Q, decimal.Decimal]
	Hasshelter         psql.WhereMod[Q, bool]
	Hasplugin          psql.WhereMod[Q, bool]
	Haschargingstation psql.WhereMod[Q, bool]
	Priceperhour       psql.WhereMod[Q, decimal.Decimal]
}

func (parkingspotWhere[Q]) AliasedAs(alias string) parkingspotWhere[Q] {
	return buildParkingspotWhere[Q](buildParkingspotColumns(alias))
}

func buildParkingspotWhere[Q psql.Filterable](cols parkingspotColumns) parkingspotWhere[Q] {
	return parkingspotWhere[Q]{
		Parkingspotid:      psql.Where[Q, int64](cols.Parkingspotid),
		Userid:             psql.Where[Q, int64](cols.Userid),
		Parkingspotuuid:    psql.Where[Q, uuid.UUID](cols.Parkingspotuuid),
		Postalcode:         psql.Where[Q, string](cols.Postalcode),
		Countrycode:        psql.Where[Q, string](cols.Countrycode),
		City:               psql.Where[Q, string](cols.City),
		State:              psql.Where[Q, string](cols.State),
		Streetaddress:      psql.Where[Q, string](cols.Streetaddress),
		Longitude:          psql.Where[Q, decimal.Decimal](cols.Longitude),
		Latitude:           psql.Where[Q, decimal.Decimal](cols.Latitude),
		Hasshelter:         psql.Where[Q, bool](cols.Hasshelter),
		Hasplugin:          psql.Where[Q, bool](cols.Hasplugin),
		Haschargingstation: psql.Where[Q, bool](cols.Haschargingstation),
		Priceperhour:       psql.Where[Q, decimal.Decimal](cols.Priceperhour),
	}
}

type parkingspotJoins[Q dialect.Joinable] struct {
	typ                    string
	UseridUser             func(context.Context) modAs[Q, userColumns]
	ParkingspotidTimeunits func(context.Context) modAs[Q, timeunitColumns]
}

func (j parkingspotJoins[Q]) aliasedAs(alias string) parkingspotJoins[Q] {
	return buildParkingspotJoins[Q](buildParkingspotColumns(alias), j.typ)
}

func buildParkingspotJoins[Q dialect.Joinable](cols parkingspotColumns, typ string) parkingspotJoins[Q] {
	return parkingspotJoins[Q]{
		typ:                    typ,
		UseridUser:             parkingspotsJoinUseridUser[Q](cols, typ),
		ParkingspotidTimeunits: parkingspotsJoinParkingspotidTimeunits[Q](cols, typ),
	}
}

// FindParkingspot retrieves a single record by primary key
// If cols is empty Find will return all columns.
func FindParkingspot(ctx context.Context, exec bob.Executor, ParkingspotidPK int64, cols ...string) (*Parkingspot, error) {
	if len(cols) == 0 {
		return Parkingspots.Query(
			ctx, exec,
			SelectWhere.Parkingspots.Parkingspotid.EQ(ParkingspotidPK),
		).One()
	}

	return Parkingspots.Query(
		ctx, exec,
		SelectWhere.Parkingspots.Parkingspotid.EQ(ParkingspotidPK),
		sm.Columns(Parkingspots.Columns().Only(cols...)),
	).One()
}

// ParkingspotExists checks the presence of a single record by primary key
func ParkingspotExists(ctx context.Context, exec bob.Executor, ParkingspotidPK int64) (bool, error) {
	return Parkingspots.Query(
		ctx, exec,
		SelectWhere.Parkingspots.Parkingspotid.EQ(ParkingspotidPK),
	).Exists()
}

// PrimaryKeyVals returns the primary key values of the Parkingspot
func (o *Parkingspot) PrimaryKeyVals() bob.Expression {
	return psql.Arg(o.Parkingspotid)
}

// Update uses an executor to update the Parkingspot
func (o *Parkingspot) Update(ctx context.Context, exec bob.Executor, s *ParkingspotSetter) error {
	return Parkingspots.Update(ctx, exec, s, o)
}

// Delete deletes a single Parkingspot record with an executor
func (o *Parkingspot) Delete(ctx context.Context, exec bob.Executor) error {
	return Parkingspots.Delete(ctx, exec, o)
}

// Reload refreshes the Parkingspot using the executor
func (o *Parkingspot) Reload(ctx context.Context, exec bob.Executor) error {
	o2, err := Parkingspots.Query(
		ctx, exec,
		SelectWhere.Parkingspots.Parkingspotid.EQ(o.Parkingspotid),
	).One()
	if err != nil {
		return err
	}
	o2.R = o.R
	*o = *o2

	return nil
}

func (o ParkingspotSlice) UpdateAll(ctx context.Context, exec bob.Executor, vals ParkingspotSetter) error {
	return Parkingspots.Update(ctx, exec, &vals, o...)
}

func (o ParkingspotSlice) DeleteAll(ctx context.Context, exec bob.Executor) error {
	return Parkingspots.Delete(ctx, exec, o...)
}

func (o ParkingspotSlice) ReloadAll(ctx context.Context, exec bob.Executor) error {
	var mods []bob.Mod[*dialect.SelectQuery]

	ParkingspotidPK := make([]int64, len(o))

	for i, o := range o {
		ParkingspotidPK[i] = o.Parkingspotid
	}

	mods = append(mods,
		SelectWhere.Parkingspots.Parkingspotid.In(ParkingspotidPK...),
	)

	o2, err := Parkingspots.Query(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, old := range o {
		for _, new := range o2 {
			if new.Parkingspotid != old.Parkingspotid {
				continue
			}
			new.R = old.R
			*old = *new
			break
		}
	}

	return nil
}

func parkingspotsJoinUseridUser[Q dialect.Joinable](from parkingspotColumns, typ string) func(context.Context) modAs[Q, userColumns] {
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

func parkingspotsJoinParkingspotidTimeunits[Q dialect.Joinable](from parkingspotColumns, typ string) func(context.Context) modAs[Q, timeunitColumns] {
	return func(ctx context.Context) modAs[Q, timeunitColumns] {
		return modAs[Q, timeunitColumns]{
			c: TimeunitColumns,
			f: func(to timeunitColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, Timeunits.Name(ctx).As(to.Alias())).On(
						to.Parkingspotid.EQ(from.Parkingspotid),
					))
				}

				return mods
			},
		}
	}
}

// UseridUser starts a query for related objects on users
func (o *Parkingspot) UseridUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) UsersQuery {
	return Users.Query(ctx, exec, append(mods,
		sm.Where(UserColumns.Userid.EQ(psql.Arg(o.Userid))),
	)...)
}

func (os ParkingspotSlice) UseridUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) UsersQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = psql.ArgGroup(o.Userid)
	}

	return Users.Query(ctx, exec, append(mods,
		sm.Where(psql.Group(UserColumns.Userid).In(PKArgs...)),
	)...)
}

// ParkingspotidTimeunits starts a query for related objects on timeunit
func (o *Parkingspot) ParkingspotidTimeunits(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) TimeunitsQuery {
	return Timeunits.Query(ctx, exec, append(mods,
		sm.Where(TimeunitColumns.Parkingspotid.EQ(psql.Arg(o.Parkingspotid))),
	)...)
}

func (os ParkingspotSlice) ParkingspotidTimeunits(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) TimeunitsQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = psql.ArgGroup(o.Parkingspotid)
	}

	return Timeunits.Query(ctx, exec, append(mods,
		sm.Where(psql.Group(TimeunitColumns.Parkingspotid).In(PKArgs...)),
	)...)
}

func (o *Parkingspot) Preload(name string, retrieved any) error {
	if o == nil {
		return nil
	}

	switch name {
	case "UseridUser":
		rel, ok := retrieved.(*User)
		if !ok {
			return fmt.Errorf("parkingspot cannot load %T as %q", retrieved, name)
		}

		o.R.UseridUser = rel

		if rel != nil {
			rel.R.UseridParkingspots = ParkingspotSlice{o}
		}
		return nil
	case "ParkingspotidTimeunits":
		rels, ok := retrieved.(TimeunitSlice)
		if !ok {
			return fmt.Errorf("parkingspot cannot load %T as %q", retrieved, name)
		}

		o.R.ParkingspotidTimeunits = rels

		for _, rel := range rels {
			if rel != nil {
				rel.R.ParkingspotidParkingspot = o
			}
		}
		return nil
	default:
		return fmt.Errorf("parkingspot has no relationship %q", name)
	}
}

func PreloadParkingspotUseridUser(opts ...psql.PreloadOption) psql.Preloader {
	return psql.Preload[*User, UserSlice](orm.Relationship{
		Name: "UseridUser",
		Sides: []orm.RelSide{
			{
				From: "parkingspot",
				To:   TableNames.Users,
				ToExpr: func(ctx context.Context) bob.Expression {
					return Users.Name(ctx)
				},
				FromColumns: []string{
					ColumnNames.Parkingspots.Userid,
				},
				ToColumns: []string{
					ColumnNames.Users.Userid,
				},
			},
		},
	}, Users.Columns().Names(), opts...)
}

func ThenLoadParkingspotUseridUser(queryMods ...bob.Mod[*dialect.SelectQuery]) psql.Loader {
	return psql.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadParkingspotUseridUser(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load ParkingspotUseridUser", retrieved)
		}

		err := loader.LoadParkingspotUseridUser(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadParkingspotUseridUser loads the parkingspot's UseridUser into the .R struct
func (o *Parkingspot) LoadParkingspotUseridUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.UseridUser = nil

	related, err := o.UseridUser(ctx, exec, mods...).One()
	if err != nil {
		return err
	}

	related.R.UseridParkingspots = ParkingspotSlice{o}

	o.R.UseridUser = related
	return nil
}

// LoadParkingspotUseridUser loads the parkingspot's UseridUser into the .R struct
func (os ParkingspotSlice) LoadParkingspotUseridUser(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
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

			rel.R.UseridParkingspots = append(rel.R.UseridParkingspots, o)

			o.R.UseridUser = rel
			break
		}
	}

	return nil
}

func ThenLoadParkingspotParkingspotidTimeunits(queryMods ...bob.Mod[*dialect.SelectQuery]) psql.Loader {
	return psql.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadParkingspotParkingspotidTimeunits(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load ParkingspotParkingspotidTimeunits", retrieved)
		}

		err := loader.LoadParkingspotParkingspotidTimeunits(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadParkingspotParkingspotidTimeunits loads the parkingspot's ParkingspotidTimeunits into the .R struct
func (o *Parkingspot) LoadParkingspotParkingspotidTimeunits(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.ParkingspotidTimeunits = nil

	related, err := o.ParkingspotidTimeunits(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, rel := range related {
		rel.R.ParkingspotidParkingspot = o
	}

	o.R.ParkingspotidTimeunits = related
	return nil
}

// LoadParkingspotParkingspotidTimeunits loads the parkingspot's ParkingspotidTimeunits into the .R struct
func (os ParkingspotSlice) LoadParkingspotParkingspotidTimeunits(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	timeunits, err := os.ParkingspotidTimeunits(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, o := range os {
		o.R.ParkingspotidTimeunits = nil
	}

	for _, o := range os {
		for _, rel := range timeunits {
			if o.Parkingspotid != rel.Parkingspotid {
				continue
			}

			rel.R.ParkingspotidParkingspot = o

			o.R.ParkingspotidTimeunits = append(o.R.ParkingspotidTimeunits, rel)
		}
	}

	return nil
}

func attachParkingspotUseridUser0(ctx context.Context, exec bob.Executor, count int, parkingspot0 *Parkingspot, user1 *User) (*Parkingspot, error) {
	setter := &ParkingspotSetter{
		Userid: omit.From(user1.Userid),
	}

	err := Parkingspots.Update(ctx, exec, setter, parkingspot0)
	if err != nil {
		return nil, fmt.Errorf("attachParkingspotUseridUser0: %w", err)
	}

	return parkingspot0, nil
}

func (parkingspot0 *Parkingspot) InsertUseridUser(ctx context.Context, exec bob.Executor, related *UserSetter) error {
	user1, err := Users.Insert(ctx, exec, related)
	if err != nil {
		return fmt.Errorf("inserting related objects: %w", err)
	}

	_, err = attachParkingspotUseridUser0(ctx, exec, 1, parkingspot0, user1)
	if err != nil {
		return err
	}

	parkingspot0.R.UseridUser = user1

	user1.R.UseridParkingspots = append(user1.R.UseridParkingspots, parkingspot0)

	return nil
}

func (parkingspot0 *Parkingspot) AttachUseridUser(ctx context.Context, exec bob.Executor, user1 *User) error {
	var err error

	_, err = attachParkingspotUseridUser0(ctx, exec, 1, parkingspot0, user1)
	if err != nil {
		return err
	}

	parkingspot0.R.UseridUser = user1

	user1.R.UseridParkingspots = append(user1.R.UseridParkingspots, parkingspot0)

	return nil
}

func insertParkingspotParkingspotidTimeunits0(ctx context.Context, exec bob.Executor, timeunits1 []*TimeunitSetter, parkingspot0 *Parkingspot) (TimeunitSlice, error) {
	for i := range timeunits1 {
		timeunits1[i].Parkingspotid = omit.From(parkingspot0.Parkingspotid)
	}

	ret, err := Timeunits.InsertMany(ctx, exec, timeunits1...)
	if err != nil {
		return ret, fmt.Errorf("insertParkingspotParkingspotidTimeunits0: %w", err)
	}

	return ret, nil
}

func attachParkingspotParkingspotidTimeunits0(ctx context.Context, exec bob.Executor, count int, timeunits1 TimeunitSlice, parkingspot0 *Parkingspot) (TimeunitSlice, error) {
	setter := &TimeunitSetter{
		Parkingspotid: omit.From(parkingspot0.Parkingspotid),
	}

	err := Timeunits.Update(ctx, exec, setter, timeunits1...)
	if err != nil {
		return nil, fmt.Errorf("attachParkingspotParkingspotidTimeunits0: %w", err)
	}

	return timeunits1, nil
}

func (parkingspot0 *Parkingspot) InsertParkingspotidTimeunits(ctx context.Context, exec bob.Executor, related ...*TimeunitSetter) error {
	if len(related) == 0 {
		return nil
	}

	timeunits1, err := insertParkingspotParkingspotidTimeunits0(ctx, exec, related, parkingspot0)
	if err != nil {
		return err
	}

	parkingspot0.R.ParkingspotidTimeunits = append(parkingspot0.R.ParkingspotidTimeunits, timeunits1...)

	for _, rel := range timeunits1 {
		rel.R.ParkingspotidParkingspot = parkingspot0
	}
	return nil
}

func (parkingspot0 *Parkingspot) AttachParkingspotidTimeunits(ctx context.Context, exec bob.Executor, related ...*Timeunit) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	timeunits1 := TimeunitSlice(related)

	_, err = attachParkingspotParkingspotidTimeunits0(ctx, exec, len(related), timeunits1, parkingspot0)
	if err != nil {
		return err
	}

	parkingspot0.R.ParkingspotidTimeunits = append(parkingspot0.R.ParkingspotidTimeunits, timeunits1...)

	for _, rel := range related {
		rel.R.ParkingspotidParkingspot = parkingspot0
	}

	return nil
}
