package dbtype

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Tstzrange struct {
	Start time.Time
	End   time.Time
}

func (r Tstzrange) Value() (driver.Value, error) {
	m := pgtype.NewMap()

	v := pgtype.Range[pgtype.Timestamptz]{
		Lower: pgtype.Timestamptz{
			Time:  r.Start,
			Valid: true,
		},
		Upper: pgtype.Timestamptz{
			Time:  r.End,
			Valid: true,
		},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}

	t, ok := m.TypeForValue(v)
	if !ok {
		return nil, fmt.Errorf("cannot find registered type for %T", v)
	}

	var buf []byte
	buf, err := m.Encode(t.OID, pgtype.TextFormatCode, v, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *Tstzrange) Scan(src any) error {
	m := pgtype.NewMap()

	var v pgtype.Range[pgtype.Timestamptz]

	t, ok := m.TypeForValue(&v)
	if !ok {
		return fmt.Errorf("cannot find registered type for %T", v)
	}

	var buf []byte
	if src != nil {
		switch src := src.(type) {
		case string:
			buf = []byte(src)
		case []byte:
			buf = src
		default:
			buf = []byte(fmt.Sprint(src))
		}
	}
	err := m.Scan(t.OID, pgtype.TextFormatCode, buf, &v)
	if err != nil {
		return err
	}

	if !v.Valid {
		return errors.New("got NULL but expected tstzrange")
	}

	if v.Lower.InfinityModifier != pgtype.Finite || v.Upper.InfinityModifier != pgtype.Finite {
		return errors.New("at least one of the timestamptz boundary is not finite")
	}

	r.Start = v.Lower.Time
	r.End = v.Upper.Time
	return nil
}
