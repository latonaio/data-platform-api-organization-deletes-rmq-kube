package dpfm_api_output_formatter

import (
	"database/sql"
	"fmt"
)

func ConvertToOrganization(rows *sql.Rows) (*Organization, error) {
	defer rows.Close()
	organization := Organization{}
	i := 0

	for rows.Next() {
		i++
		err := rows.Scan(
			&organization.BusinessPartner,
			&organization.Organization,
			&organization.IsMarkedForDeletion,
		)
		if err != nil {
			fmt.Printf("err = %+v \n", err)
			return &organization, err
		}

	}
	if i == 0 {
		fmt.Printf("DBに対象のレコードが存在しません。")
		return &organization, nil
	}

	return &organization, nil
}
