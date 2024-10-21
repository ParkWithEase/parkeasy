// Code generated by modelgen. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package dbmodels

import (
	"hash/maphash"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/clause"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

var TableNames = struct {
	Auths        string
	Bookings     string
	Cars         string
	Parkingspots string
	Resettokens  string
	Sessions     string
	Timeunits    string
	Users        string
}{
	Auths:        "auth",
	Bookings:     "booking",
	Cars:         "car",
	Parkingspots: "parkingspot",
	Resettokens:  "resettoken",
	Sessions:     "sessions",
	Timeunits:    "timeunit",
	Users:        "users",
}

var ColumnNames = struct {
	Auths        authColumnNames
	Bookings     bookingColumnNames
	Cars         carColumnNames
	Parkingspots parkingspotColumnNames
	Resettokens  resettokenColumnNames
	Sessions     sessionColumnNames
	Timeunits    timeunitColumnNames
	Users        userColumnNames
}{
	Auths: authColumnNames{
		Authid:       "authid",
		Authuuid:     "authuuid",
		Email:        "email",
		Passwordhash: "passwordhash",
	},
	Bookings: bookingColumnNames{
		Bookingid:   "bookingid",
		Buyeruserid: "buyeruserid",
		Bookinguuid: "bookinguuid",
		Paidamount:  "paidamount",
	},
	Cars: carColumnNames{
		Carid:        "carid",
		Userid:       "userid",
		Caruuid:      "caruuid",
		Licenseplate: "licenseplate",
		Make:         "make",
		Model:        "model",
		Color:        "color",
	},
	Parkingspots: parkingspotColumnNames{
		Parkingspotid:      "parkingspotid",
		Userid:             "userid",
		Parkingspotuuid:    "parkingspotuuid",
		Postalcode:         "postalcode",
		Countrycode:        "countrycode",
		City:               "city",
		State:              "state",
		Streetaddress:      "streetaddress",
		Longitude:          "longitude",
		Latitude:           "latitude",
		Hasshelter:         "hasshelter",
		Hasplugin:          "hasplugin",
		Haschargingstation: "haschargingstation",
		Priceperhour:       "priceperhour",
	},
	Resettokens: resettokenColumnNames{
		Token:    "token",
		Authuuid: "authuuid",
		Expiry:   "expiry",
	},
	Sessions: sessionColumnNames{
		Token:  "token",
		Data:   "data",
		Expiry: "expiry",
	},
	Timeunits: timeunitColumnNames{
		Starttime:       "starttime",
		Endtime:         "endtime",
		Parkingspotuuid: "parkingspotuuid",
		Bookingid:       "bookingid",
		Status:          "status",
	},
	Users: userColumnNames{
		Userid:     "userid",
		Useruuid:   "useruuid",
		Authuuid:   "authuuid",
		Fullname:   "fullname",
		Email:      "email",
		Isverified: "isverified",
		Addedat:    "addedat",
	},
}

var (
	SelectWhere = Where[*dialect.SelectQuery]()
	InsertWhere = Where[*dialect.InsertQuery]()
	UpdateWhere = Where[*dialect.UpdateQuery]()
	DeleteWhere = Where[*dialect.DeleteQuery]()
)

func Where[Q psql.Filterable]() struct {
	Auths        authWhere[Q]
	Bookings     bookingWhere[Q]
	Cars         carWhere[Q]
	Parkingspots parkingspotWhere[Q]
	Resettokens  resettokenWhere[Q]
	Sessions     sessionWhere[Q]
	Timeunits    timeunitWhere[Q]
	Users        userWhere[Q]
} {
	return struct {
		Auths        authWhere[Q]
		Bookings     bookingWhere[Q]
		Cars         carWhere[Q]
		Parkingspots parkingspotWhere[Q]
		Resettokens  resettokenWhere[Q]
		Sessions     sessionWhere[Q]
		Timeunits    timeunitWhere[Q]
		Users        userWhere[Q]
	}{
		Auths:        buildAuthWhere[Q](AuthColumns),
		Bookings:     buildBookingWhere[Q](BookingColumns),
		Cars:         buildCarWhere[Q](CarColumns),
		Parkingspots: buildParkingspotWhere[Q](ParkingspotColumns),
		Resettokens:  buildResettokenWhere[Q](ResettokenColumns),
		Sessions:     buildSessionWhere[Q](SessionColumns),
		Timeunits:    buildTimeunitWhere[Q](TimeunitColumns),
		Users:        buildUserWhere[Q](UserColumns),
	}
}

var (
	SelectJoins = getJoins[*dialect.SelectQuery]()
	UpdateJoins = getJoins[*dialect.UpdateQuery]()
	DeleteJoins = getJoins[*dialect.DeleteQuery]()
)

type joinSet[Q interface{ aliasedAs(string) Q }] struct {
	InnerJoin Q
	LeftJoin  Q
	RightJoin Q
}

func (j joinSet[Q]) AliasedAs(alias string) joinSet[Q] {
	return joinSet[Q]{
		InnerJoin: j.InnerJoin.aliasedAs(alias),
		LeftJoin:  j.LeftJoin.aliasedAs(alias),
		RightJoin: j.RightJoin.aliasedAs(alias),
	}
}

type joins[Q dialect.Joinable] struct {
	Auths        joinSet[authJoins[Q]]
	Bookings     joinSet[bookingJoins[Q]]
	Cars         joinSet[carJoins[Q]]
	Parkingspots joinSet[parkingspotJoins[Q]]
	Resettokens  joinSet[resettokenJoins[Q]]
	Timeunits    joinSet[timeunitJoins[Q]]
	Users        joinSet[userJoins[Q]]
}

func buildJoinSet[Q interface{ aliasedAs(string) Q }, C any, F func(C, string) Q](c C, f F) joinSet[Q] {
	return joinSet[Q]{
		InnerJoin: f(c, clause.InnerJoin),
		LeftJoin:  f(c, clause.LeftJoin),
		RightJoin: f(c, clause.RightJoin),
	}
}

func getJoins[Q dialect.Joinable]() joins[Q] {
	return joins[Q]{
		Auths:        buildJoinSet[authJoins[Q]](AuthColumns, buildAuthJoins),
		Bookings:     buildJoinSet[bookingJoins[Q]](BookingColumns, buildBookingJoins),
		Cars:         buildJoinSet[carJoins[Q]](CarColumns, buildCarJoins),
		Parkingspots: buildJoinSet[parkingspotJoins[Q]](ParkingspotColumns, buildParkingspotJoins),
		Resettokens:  buildJoinSet[resettokenJoins[Q]](ResettokenColumns, buildResettokenJoins),
		Timeunits:    buildJoinSet[timeunitJoins[Q]](TimeunitColumns, buildTimeunitJoins),
		Users:        buildJoinSet[userJoins[Q]](UserColumns, buildUserJoins),
	}
}

type modAs[Q any, C interface{ AliasedAs(string) C }] struct {
	c C
	f func(C) bob.Mod[Q]
}

func (m modAs[Q, C]) Apply(q Q) {
	m.f(m.c).Apply(q)
}

func (m modAs[Q, C]) AliasedAs(alias string) bob.Mod[Q] {
	m.c = m.c.AliasedAs(alias)
	return m
}

func randInt() int64 {
	out := int64(new(maphash.Hash).Sum64())

	if out < 0 {
		return -out % 10000
	}

	return out % 10000
}
