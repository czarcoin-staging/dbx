// Copyright (C) 2019 Storj Labs, Inc.
// Copyright (C) 2017 Space Monkey, Inc.
// See LICENSE for copying information.

package sql

import (
	"fmt"
	"strconv"
	"strings"

	"storj.io/dbx/consts"
	"storj.io/dbx/ir"
)

type cockroach struct {
}

func Cockroach() Dialect {
	return &cockroach{}
}

func (p *cockroach) Name() string {
	return "cockroach"
}

func (p *cockroach) Features() Features {
	return Features{
		Returning:           true,
		PositionalArguments: true,
		NoLimitToken:        "ALL",
	}
}

func (p *cockroach) RowId() string {
	return ""
}

func (p *cockroach) ColumnType(field *ir.Field) string {
	switch field.Type {
	case consts.SerialField:
		return "serial"
	case consts.Serial64Field:
		return "bigserial"
	case consts.IntField:
		return "integer"
	case consts.Int64Field:
		return "bigint"
	case consts.UintField:
		return "integer"
	case consts.Uint64Field:
		return "bigint"
	case consts.FloatField:
		return "real"
	case consts.Float64Field:
		return "double precision"
	case consts.TextField:
		if field.Length > 0 {
			return fmt.Sprintf("varchar(%d)", field.Length)
		}
		return "text"
	case consts.BoolField:
		return "boolean"
	case consts.TimestampField:
		return "timestamp with time zone"
	case consts.TimestampUTCField:
		return "timestamp"
	case consts.BlobField:
		return "bytea"
	case consts.DateField:
		return "date"
	default:
		panic(fmt.Sprintf("unhandled field type %s", field.Type))
	}
}

func (p *cockroach) Rebind(sql string) string {
	out := make([]byte, 0, len(sql)+10)

	j := 1
	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		if ch != '?' {
			out = append(out, ch)
			continue
		}

		out = append(out, '$')
		out = append(out, strconv.Itoa(j)...)
		j++
	}

	return string(out)
}

func (p *cockroach) ArgumentPrefix() string { return "$" }

var cockroachEscaper = strings.NewReplacer(
	`'`, `\'`,
	`\`, `\\`,
)

func (p *cockroach) EscapeString(s string) string {
	return cockroachEscaper.Replace(s)
}

func (p *cockroach) BoolLit(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
