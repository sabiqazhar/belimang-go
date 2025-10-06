package helper

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func FloatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	// Use string formatting to avoid float precision issues
	err := n.Scan(fmt.Sprintf("%.2f", f))
	if err != nil {
		return pgtype.Numeric{}
	} // .2f matches your DECIMAL(10, 2)
	return n
}
