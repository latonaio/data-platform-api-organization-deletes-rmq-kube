package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-organization-deletes-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-organization-deletes-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-organization-deletes-rmq-kube/config"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	database "github.com/latonaio/golang-mysql-network-connector"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
	"golang.org/x/xerrors"
)

type DPFMAPICaller struct {
	ctx  context.Context
	conf *config.Conf
	rmq  *rabbitmq.RabbitmqClient
	db   *database.Mysql
}

func NewDPFMAPICaller(
	conf *config.Conf, rmq *rabbitmq.RabbitmqClient, db *database.Mysql,
) *DPFMAPICaller {
	return &DPFMAPICaller{
		ctx:  context.Background(),
		conf: conf,
		rmq:  rmq,
		db:   db,
	}
}

func (c *DPFMAPICaller) AsyncDeletes(
	accepter []string,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	log *logger.Logger,
) (interface{}, []error) {
	var response interface{}
	if input.APIType == "deletes" {
		response = c.deleteSqlProcess(input, output, accepter, log)
	}

	return response, nil
}

func (c *DPFMAPICaller) deleteSqlProcess(
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	accepter []string,
	log *logger.Logger,
) *dpfm_api_output_formatter.Message {
	var organizationData *dpfm_api_output_formatter.Organization
	for _, a := range accepter {
		switch a {
		case "Organization":
			h := c.organizationDelete(input, output, log)
			organizationData = h
			if h == nil {
				continue
			}
		}
	}

	return &dpfm_api_output_formatter.Message{
		Organization: organizationData,
	}
}

func (c *DPFMAPICaller) organizationDelete(
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	log *logger.Logger,
) *dpfm_api_output_formatter.Organization {
	sessionID := input.RuntimeSessionID
	organization := c.OrganizationRead(input, log)
	organization.SupplyChainRelationshipID = input.Organization.SupplyChainRelationshipID
	organization.Buyer = input.Organization.Buyer
	organization.Seller = input.Organization.Seller
	organization.ConditionRecord = input.Organization.ConditionRecord
	organization.ConditionSequentialNumber = input.Organization.ConditionSequentialNumber
	organization.Product = input.Organization.Product
	organization.ConditionValidityStartDate = input.Organization.ConditionValidityStartDate
	organization.ConditionValidityEndDate = input.Organization.ConditionValidityEndDate
	organization.IsMarkedForDeletion = input.Organization.IsMarkedForDeletion
	res, err := c.rmq.SessionKeepRequest(nil, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": organization, "function": "OrganizationOrganization", "runtime_session_id": sessionID})
	if err != nil {
		err = xerrors.Errorf("rmq error: %w", err)
		log.Error("%+v", err)
		return nil
	}
	res.Success()
	if !checkResult(res) {
		output.SQLUpdateResult = getBoolPtr(false)
		output.SQLUpdateError = "Organization Data cannot delete"
		return nil
	}

	return organization
}

func checkResult(msg rabbitmq.RabbitmqMessage) bool {
	data := msg.Data()
	d, ok := data["result"]
	if !ok {
		return false
	}
	result, ok := d.(string)
	if !ok {
		return false
	}
	return result == "success"
}

func getBoolPtr(b bool) *bool {
	return &b
}
