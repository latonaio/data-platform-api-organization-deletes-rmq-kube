package dpfm_api_output_formatter

import (
	"database/sql"
	"fmt"
)

func ConvertToExchangeRate(rows *sql.Rows) (*exchangeRate, error) {
	defer rows.Close()
	exchangeRate := ExchangeRate{}
	i := 0

	for rows.Next() {
		i++
		err := rows.Scan(
			&exchangeRate.CurrencyTo,
			&exchangeRate.CurrencyFrom,
			&exchangeRate.ValidityStartDate,
			&exchangeRate.ValidityEndDate,
			&exchangeRate.IsMarkedForDeletion,
		)
		if err != nil {
			fmt.Printf("err = %+v \n", err)
			return &exchangeRate, err
		}

	}
	if i == 0 {
		fmt.Printf("DBに対象のレコードが存在しません。")
		return &exchangeRate, nil
	}

	return &exchangeRate, nil
}
