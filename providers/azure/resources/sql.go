// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package resources

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"go.mondoo.com/cnquery/llx"
	"go.mondoo.com/cnquery/providers-sdk/v1/plugin"
	"go.mondoo.com/cnquery/providers-sdk/v1/util/convert"
	"go.mondoo.com/cnquery/providers/azure/connection"
	"go.mondoo.com/cnquery/types"

	sql "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
)

func (a *mqlAzureSubscriptionSql) id() (string, error) {
	return "azure.subscription.sql/" + a.SubscriptionId.Data, nil
}

func initAzureSubscriptionSql(runtime *plugin.Runtime, args map[string]*llx.RawData) (map[string]*llx.RawData, plugin.Resource, error) {
	if len(args) > 0 {
		return args, nil, nil
	}

	conn := runtime.Connection.(*connection.AzureConnection)
	args["subscriptionId"] = llx.StringData(conn.SubId())

	return args, nil, nil
}

func (a *mqlAzureSubscriptionSqlDatabase) id() (string, error) {
	return a.Id.Data, nil
}

func (a *mqlAzureSubscriptionSqlDatabaseusage) id() (string, error) {
	return a.Id.Data, nil
}

func (a *mqlAzureSubscriptionSqlServer) id() (string, error) {
	return a.Id.Data, nil
}

func (a *mqlAzureSubscriptionSqlServerAdministrator) id() (string, error) {
	return a.Id.Data, nil
}

func (a *mqlAzureSubscriptionSqlServerVulnerabilityassessmentsettings) id() (string, error) {
	return a.Id.Data, nil
}

func (a *mqlAzureSubscriptionSql) servers() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	subId := a.SubscriptionId.Data
	dbClient, err := sql.NewServersClient(subId, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	pager := dbClient.NewListPager(&sql.ServersClientListOptions{})
	res := []interface{}{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, dbServer := range page.Value {
			properties, err := convert.JsonToDict(dbServer.Properties)
			if err != nil {
				return nil, err
			}

			mqlAzureDbServer, err := CreateResource(a.MqlRuntime, "azure.subscription.sql.server",
				map[string]*llx.RawData{
					"id":         llx.StringData(convert.ToString(dbServer.ID)),
					"name":       llx.StringData(convert.ToString(dbServer.Name)),
					"location":   llx.StringData(convert.ToString(dbServer.Location)),
					"tags":       llx.MapData(convert.PtrMapStrToInterface(dbServer.Tags), types.String),
					"type":       llx.StringData(convert.ToString(dbServer.Type)),
					"properties": llx.DictData(properties),
				})
			if err != nil {
				return nil, err
			}
			res = append(res, mqlAzureDbServer)
		}
	}

	return res, nil
}

func (a *mqlAzureSubscriptionSqlServer) databases() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}
	dbDatabaseClient, err := sql.NewDatabasesClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	pager := dbDatabaseClient.NewListByServerPager(resourceID.ResourceGroup, server, &sql.DatabasesClientListByServerOptions{})
	res := []interface{}{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, entry := range page.Value {
			args := map[string]*llx.RawData{
				"id":               llx.StringData(convert.ToString(entry.ID)),
				"name":             llx.StringData(convert.ToString(entry.Name)),
				"type":             llx.StringData(convert.ToString(entry.Type)),
				"collation":        llx.StringData(convert.ToString(entry.Properties.Collation)),
				"creationDate":     llx.TimeData(*entry.Properties.CreationDate),
				"databaseId":       llx.StringData(convert.ToString(entry.Properties.DatabaseID)),
				"createMode":       llx.StringData(convert.ToString((*string)(entry.Properties.CreateMode))),
				"sourceDatabaseId": llx.StringData(convert.ToString(entry.Properties.SourceDatabaseID)),
				"recoveryServicesRecoveryPointResourceId": llx.StringData(convert.ToString(entry.Properties.RecoveryServicesRecoveryPointID)),
				"edition":                       llx.StringData(convert.ToString(entry.SKU.Tier)),
				"maxSizeBytes":                  llx.IntData(convert.ToInt64(entry.Properties.MaxSizeBytes)),
				"requestedServiceObjectiveName": llx.StringData(convert.ToString(entry.Properties.RequestedServiceObjectiveName)),
				"serviceLevelObjective":         llx.StringData(convert.ToString(entry.Properties.CurrentServiceObjectiveName)),
				"status":                        llx.StringData(convert.ToString((*string)(entry.Properties.Status))),
				"elasticPoolName":               llx.StringData(convert.ToString(entry.Properties.ElasticPoolID)),
				"defaultSecondaryLocation":      llx.StringData(convert.ToString(entry.Properties.DefaultSecondaryLocation)),
				"failoverGroupId":               llx.StringData(convert.ToString(entry.Properties.FailoverGroupID)),
				"readScale":                     llx.StringData(convert.ToString((*string)(entry.Properties.ReadScale))),
				"sampleName":                    llx.StringData(convert.ToString((*string)(entry.Properties.SampleName))),
				"zoneRedundant":                 llx.BoolData(convert.ToBool(entry.Properties.ZoneRedundant)),
				"earliestRestoreDate":           llx.TimeDataPtr(entry.Properties.EarliestRestoreDate),
				"sourceDatabaseDeletionDate":    llx.TimeDataPtr(entry.Properties.SourceDatabaseDeletionDate),
				"restorePointInTime":            llx.TimeDataPtr(entry.Properties.RestorePointInTime),
			}

			mqlAzureDatabase, err := CreateResource(a.MqlRuntime, "azure.subscription.sql.database", args)
			if err != nil {
				return nil, err
			}
			res = append(res, mqlAzureDatabase)
		}
	}
	return res, nil
}

func (a *mqlAzureSubscriptionSqlServer) firewallRules() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	dbFirewallClient, err := sql.NewFirewallRulesClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}

	pager := dbFirewallClient.NewListByServerPager(resourceID.ResourceGroup, server, &sql.FirewallRulesClientListByServerOptions{})
	res := []interface{}{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, entry := range page.Value {
			mqlFireWallRule, err := CreateResource(a.MqlRuntime, "azure.subscription.sql.firewallrule",
				map[string]*llx.RawData{
					"id":             llx.StringData(convert.ToString(entry.ID)),
					"name":           llx.StringData(convert.ToString(entry.Name)),
					"type":           llx.StringData(convert.ToString(entry.Type)),
					"startIpAddress": llx.StringData(convert.ToString(entry.Properties.StartIPAddress)),
					"endIpAddress":   llx.StringData(convert.ToString(entry.Properties.EndIPAddress)),
				})
			if err != nil {
				return nil, err
			}
			res = append(res, mqlFireWallRule)
		}
	}
	return res, nil
}

func (a *mqlAzureSubscriptionSqlServer) virtualNetworkRules() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	client, err := sql.NewVirtualNetworkRulesClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}

	pager := client.NewListByServerPager(resourceID.ResourceGroup, server, &sql.VirtualNetworkRulesClientListByServerOptions{})
	res := []interface{}{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, entry := range page.Value {
			properties, err := convert.JsonToDict(entry)
			if err != nil {
				return nil, err
			}

			mqlAzure, err := CreateResource(a.MqlRuntime, "azure.subscription.sql.virtualNetworkRule",
				map[string]*llx.RawData{
					"id":                     llx.StringData(convert.ToString(entry.ID)),
					"name":                   llx.StringData(convert.ToString(entry.Name)),
					"type":                   llx.StringData(convert.ToString(entry.Type)),
					"properties":             llx.DictData(properties),
					"virtualNetworkSubnetId": llx.StringData(convert.ToString(entry.Properties.VirtualNetworkSubnetID)),
				})
			if err != nil {
				return nil, err
			}
			res = append(res, mqlAzure)
		}
	}
	return res, nil
}

func (a *mqlAzureSubscriptionSqlServer) azureAdAdministrators() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}
	administratorClient, err := sql.NewServerAzureADAdministratorsClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	pager := administratorClient.NewListByServerPager(resourceID.ResourceGroup, server, &sql.ServerAzureADAdministratorsClientListByServerOptions{})
	res := []interface{}{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, entry := range page.Value {
			mqlAzureSqlAdministrator, err := CreateResource(a.MqlRuntime, "azure.subscription.sql.server.administrator",
				map[string]*llx.RawData{
					"id":                llx.StringData(convert.ToString(entry.ID)),
					"name":              llx.StringData(convert.ToString(entry.Name)),
					"type":              llx.StringData(convert.ToString(entry.Type)),
					"administratorType": llx.StringData(convert.ToString((*string)(entry.Properties.AdministratorType))),
					"login":             llx.StringData(convert.ToString(entry.Properties.Login)),
					"sid":               llx.StringData(convert.ToString(entry.Properties.Sid)),
					"tenantId":          llx.StringData(convert.ToString(entry.Properties.TenantID)),
				})
			if err != nil {
				return nil, err
			}
			res = append(res, mqlAzureSqlAdministrator)
		}
	}

	return res, nil
}

func (a *mqlAzureSubscriptionSqlServer) connectionPolicy() (interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	connectionClient, err := sql.NewServerConnectionPoliciesClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	policy, err := connectionClient.Get(ctx, resourceID.ResourceGroup, server, sql.ConnectionPolicyNameDefault, &sql.ServerConnectionPoliciesClientGetOptions{})
	if err != nil {
		return nil, err
	}

	return convert.JsonToDict(policy)
}

func (a *mqlAzureSubscriptionSqlServer) securityAlertPolicy() (interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	secAlertClient, err := sql.NewServerSecurityAlertPoliciesClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	policy, err := secAlertClient.Get(ctx, resourceID.ResourceGroup, server, sql.SecurityAlertPolicyNameDefault, &sql.ServerSecurityAlertPoliciesClientGetOptions{})
	if err != nil {
		return nil, err
	}

	return convert.JsonToDict(policy.ServerSecurityAlertPolicy.Properties)
}

func (a *mqlAzureSubscriptionSqlServer) auditingPolicy() (interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}
	auditClient, err := sql.NewServerBlobAuditingPoliciesClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	policy, err := auditClient.Get(ctx, resourceID.ResourceGroup, server, &sql.ServerBlobAuditingPoliciesClientGetOptions{})
	if err != nil {
		return nil, err
	}

	return convert.JsonToDict(policy.ServerBlobAuditingPolicy.Properties)
}

func (a *mqlAzureSubscriptionSqlServer) threatDetectionPolicy() (interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	serverClient, err := sql.NewServerAdvancedThreatProtectionSettingsClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	threatPolicy, err := serverClient.Get(ctx, resourceID.ResourceGroup, server, sql.AdvancedThreatProtectionNameDefault, &sql.ServerAdvancedThreatProtectionSettingsClientGetOptions{})
	if err != nil {
		return nil, err
	}

	return convert.JsonToDict(threatPolicy.Properties)
}

func (a *mqlAzureSubscriptionSqlServer) encryptionProtector() (interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	client, err := sql.NewEncryptionProtectorsClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	policy, err := client.Get(ctx, resourceID.ResourceGroup, server, sql.EncryptionProtectorNameCurrent, &sql.EncryptionProtectorsClientGetOptions{})
	if err != nil {
		return nil, err
	}
	return convert.JsonToDict(policy.EncryptionProtector.Properties)
}

func (a *mqlAzureSubscriptionSqlServer) vulnerabilityAssessmentSettings() (*mqlAzureSubscriptionSqlServerVulnerabilityassessmentsettings, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	serverClient, err := sql.NewServerVulnerabilityAssessmentsClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	vaSettings, err := serverClient.Get(ctx, resourceID.ResourceGroup, server, sql.VulnerabilityAssessmentNameDefault, &sql.ServerVulnerabilityAssessmentsClientGetOptions{})
	if err != nil {
		return nil, err
	}
	res, err := CreateResource(a.MqlRuntime, "azure.subscription.sql.server.vulnerabilityassessmentsettings",
		map[string]*llx.RawData{
			"id":                      llx.StringData(convert.ToString(vaSettings.ID)),
			"name":                    llx.StringData(convert.ToString(vaSettings.Name)),
			"type":                    llx.StringData(convert.ToString(vaSettings.Type)),
			"storageContainerPath":    llx.StringData(convert.ToString(vaSettings.Properties.StorageContainerPath)),
			"storageAccountAccessKey": llx.StringData(convert.ToString(vaSettings.Properties.StorageAccountAccessKey)),
			"storageContainerSasKey":  llx.StringData(convert.ToString(vaSettings.Properties.StorageContainerSasKey)),
			"recurringScanEnabled":    llx.BoolData(convert.ToBool(vaSettings.Properties.RecurringScans.IsEnabled)),
			"recurringScanEmails":     llx.ArrayData(llx.TArr2Raw(convert.ToListFromPtrs(vaSettings.Properties.RecurringScans.Emails)), types.String),
			"mailSubscriptionAdmins":  llx.BoolData(convert.ToBool(vaSettings.Properties.RecurringScans.EmailSubscriptionAdmins)),
		})
	if err != nil {
		return nil, err
	}
	return res.(*mqlAzureSubscriptionSqlServerVulnerabilityassessmentsettings), err
}

func (a *mqlAzureSubscriptionSqlDatabase) transparentDataEncryption() (interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	database, err := resourceID.Component("databases")
	if err != nil {
		return nil, err
	}

	client, err := sql.NewTransparentDataEncryptionsClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}

	policy, err := client.Get(ctx, resourceID.ResourceGroup, server, database, sql.TransparentDataEncryptionNameCurrent, &sql.TransparentDataEncryptionsClientGetOptions{})
	if err != nil {
		return nil, err
	}

	return convert.JsonToDict(policy.LogicalDatabaseTransparentDataEncryption.Properties)
}

func (a *mqlAzureSubscriptionSqlDatabase) advisor() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	database, err := resourceID.Component("databases")
	if err != nil {
		return nil, err
	}
	client, err := sql.NewDatabaseAdvisorsClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}

	// it's an OData API, supports $expand. We can get the recommendedActions for all advisors here.
	expandRecommendedActions := "recommendedActions"
	advisors, err := client.ListByDatabase(ctx, resourceID.ResourceGroup, server, database, &sql.DatabaseAdvisorsClientListByDatabaseOptions{Expand: &expandRecommendedActions})
	if err != nil {
		return nil, err
	}

	res := []interface{}{}
	for _, entry := range advisors.AdvisorArray {
		dict, err := convert.JsonToDict(entry)
		if err != nil {
			return nil, err
		}

		res = append(res, dict)
	}

	return res, nil
}

func (a *mqlAzureSubscriptionSqlDatabase) threatDetectionPolicy() (interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	database, err := resourceID.Component("databases")
	if err != nil {
		return nil, err
	}
	client, err := sql.NewDatabaseSecurityAlertPoliciesClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}

	policy, err := client.Get(ctx, resourceID.ResourceGroup, server, database, sql.SecurityAlertPolicyNameDefault, &sql.DatabaseSecurityAlertPoliciesClientGetOptions{})
	if err != nil {
		return nil, err
	}

	return convert.JsonToDict(policy.DatabaseSecurityAlertPolicy.Properties)
}

func (a *mqlAzureSubscriptionSqlDatabase) connectionPolicy() (interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	connectionClient, err := sql.NewServerConnectionPoliciesClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	policy, err := connectionClient.Get(ctx, resourceID.ResourceGroup, server, sql.ConnectionPolicyNameDefault, &sql.ServerConnectionPoliciesClientGetOptions{})
	if err != nil {
		return nil, err
	}

	return convert.JsonToDict(policy.ServerConnectionPolicy.Properties)
}

func (a *mqlAzureSubscriptionSqlDatabase) auditingPolicy() (interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	database, err := resourceID.Component("databases")
	if err != nil {
		return nil, err
	}

	auditClient, err := sql.NewDatabaseBlobAuditingPoliciesClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}

	policy, err := auditClient.Get(ctx, resourceID.ResourceGroup, server, database, &sql.DatabaseBlobAuditingPoliciesClientGetOptions{})
	if err != nil {
		return nil, err
	}

	return convert.JsonToDict(policy.DatabaseBlobAuditingPolicy.Properties)
}

func (a *mqlAzureSubscriptionSqlDatabase) usage() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AzureConnection)
	ctx := context.Background()
	token := conn.Token()
	id := a.Id.Data
	resourceID, err := ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	server, err := resourceID.Component("servers")
	if err != nil {
		return nil, err
	}

	database, err := resourceID.Component("databases")
	if err != nil {
		return nil, err
	}

	client, err := sql.NewDatabaseUsagesClient(resourceID.SubscriptionID, token, &arm.ClientOptions{})
	if err != nil {
		return nil, err
	}
	pager := client.NewListByDatabasePager(resourceID.ResourceGroup, server, database, &sql.DatabaseUsagesClientListByDatabaseOptions{})
	res := []interface{}{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, entry := range page.Value {
			mqlAzureSqlUsage, err := CreateResource(a.MqlRuntime, "azure.subscription.sql.databaseusage",
				map[string]*llx.RawData{
					"id":           llx.StringData(convert.ToString(entry.ID)),
					"name":         llx.StringData(convert.ToString(entry.Name)),
					"resourceName": llx.StringData(convert.ToString(entry.Name)),
					"displayName":  llx.StringData(convert.ToString(entry.Properties.DisplayName)),
					"currentValue": llx.FloatData(convert.ToFloat64(entry.Properties.CurrentValue)),
					"limit":        llx.FloatData(convert.ToFloat64(entry.Properties.Limit)),
					"unit":         llx.StringData(convert.ToString(entry.Properties.Unit)),
				})
			if err != nil {
				return nil, err
			}
			res = append(res, mqlAzureSqlUsage)
		}
	}

	return res, nil
}
