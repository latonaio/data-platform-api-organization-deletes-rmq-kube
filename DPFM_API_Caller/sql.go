package dpfm_api_caller

import (
	dpfm_api_input_reader "data-platform-api-organization-deletes-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-organization-deletes-rmq-kube/DPFM_API_Output_Formatter"
	"fmt"
	"strings"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
)

func (c *DPFMAPICaller) Organization(
	input *dpfm_api_input_reader.SDC,
	log *logger.Logger,
) *dpfm_api_output_formatter.Organization {

	where := strings.Join([]string{
		fmt.Sprintf("WHERE organization.BusinessPartner = %d ", input.Organization.BusinessPartner),
		fmt.Sprintf("AND organization.Organization = \"%s\" ", input.Organization.Organization),
	}, "")

	rows, err := c.db.Query(
		`SELECT 
    	organization.BusinessPartner,
    	organization.Organization
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_organization_organization_data as organization 
		` + where + ` ;`)
	if err != nil {
		log.Error("%+v", err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToOrganization(rows)
	if err != nil {
		log.Error("%+v", err)
		return nil
	}

	return data
}
