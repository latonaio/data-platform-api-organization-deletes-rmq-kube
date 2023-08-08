package dpfm_api_caller

import (
	dpfm_api_input_reader "data-platform-api-exchange-rate-deletes-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-exchange-rate-deletes-rmq-kube/DPFM_API_Output_Formatter"
	"fmt"
	"strings"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
)

func (c *DPFMAPICaller) ExchangeRate(
	input *dpfm_api_input_reader.SDC,
	log *logger.Logger,
) *dpfm_api_output_formatter.ExchangeRate {

	where := strings.Join([]string{
		fmt.Sprintf("WHERE exchangeRate.CurrencyTo = \"%s\ ", input.ExchangeRate.CurrencyTo),
		fmt.Sprintf("AND exchangeRate.CurrencyFrom = \"%s\" ", input.ExchangeRate.CurrencyFrom),
		fmt.Sprintf("AND exchangeRate.ValidityStartDate = \"%s\" ", input.ExchangeRate.ValidityStartDate),
		fmt.Sprintf("AND exchangeRate.ValidityEndDate = \"%s\" ", input.ExchangeRate.ValidityEndDate),
	}, "")

	rows, err := c.db.Query(
		`SELECT 
    	exchangeRate.ExchangeRate
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_exchangeRate_exchangeRate_data as exchangeRate 
		` + where + ` ;`)
	if err != nil {
		log.Error("%+v", err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToExchangeRate(rows)
	if err != nil {
		log.Error("%+v", err)
		return nil
	}

	return data
}
