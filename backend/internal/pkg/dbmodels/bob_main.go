// Code generated by BobGen psql v0.28.1. DO NOT EDIT.
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
	Auths       string
	Resettokens string
	Sessions    string
	Users       string
}{
	Auths:       "auth",
	Resettokens: "resettoken",
	Sessions:    "sessions",
	Users:       "users",
}

var ColumnNames = struct {
	Auths       authColumnNames
	Resettokens resettokenColumnNames
	Sessions    sessionColumnNames
	Users       userColumnNames
}{
	Auths: authColumnNames{
		Authid:       "authid",
		Authuuid:     "authuuid",
		Email:        "email",
		Passwordhash: "passwordhash",
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
	Auths       authWhere[Q]
	Resettokens resettokenWhere[Q]
	Sessions    sessionWhere[Q]
	Users       userWhere[Q]
} {
	return struct {
		Auths       authWhere[Q]
		Resettokens resettokenWhere[Q]
		Sessions    sessionWhere[Q]
		Users       userWhere[Q]
	}{
		Auths:       buildAuthWhere[Q](AuthColumns),
		Resettokens: buildResettokenWhere[Q](ResettokenColumns),
		Sessions:    buildSessionWhere[Q](SessionColumns),
		Users:       buildUserWhere[Q](UserColumns),
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
	Auths       joinSet[authJoins[Q]]
	Resettokens joinSet[resettokenJoins[Q]]
	Users       joinSet[userJoins[Q]]
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
		Auths:       buildJoinSet[authJoins[Q]](AuthColumns, buildAuthJoins),
		Resettokens: buildJoinSet[resettokenJoins[Q]](ResettokenColumns, buildResettokenJoins),
		Users:       buildJoinSet[userJoins[Q]](UserColumns, buildUserJoins),
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
