package dbtype

import (
	"database/sql/driver"
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

	v := pgtype.Range[time.Time]{
		Lower:     r.Start,
		Upper:     r.End,
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}

	t, ok := m.TypeForValue(&v)
	if !ok {
		return nil, fmt.Errorf("cannot find registered type for %T", v)
	}

	var buf []byte
	buf, err := m.Encode(t.OID, pgtype.TextFormatCode, &v, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *Tstzrange) Scan(src any) error {
	m := pgtype.NewMap()

	var v pgtype.Range[time.Time]

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

	r.Start = v.Lower
	r.End = v.Upper
	return nil
}
