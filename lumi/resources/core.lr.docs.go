// Code generated by lumi. DO NOT EDIT.
package resources

import (
	"go.mondoo.io/mondoo/lumi/lr/docs"
)

var ResourceDocs = docs.LrDocs{
	Resources: map[string]*docs.LrDocsEntry{
		"k8s.daemonset": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"kubernetes"},
			},
		},
		"aws.config.rule": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.dms": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.dynamodb.table": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.iam.virtualmfadevice": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"vsphere": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Display Information about the vSphere environment",
					Query: "vsphere.about",
				},
				{
					Title: "Display ESXi host moid and properties",
					Query: "vsphere.datacenters { hosts { moid properties } }",
				},
				{
					Title: "Display NTP server for all ESXi hosts",
					Query: "vsphere.datacenters { hosts { ntp.server } }",
				},
				{
					Title: "Ensure a specific NTP Server is set",
					Query: "vsphere.datacenters { hosts { ntp.server.any(_ == \"10.31.21.2\") } }",
				},
				{
					Title: "Ensure specific VmkNics properties for all management VmkNics",
					Query: "vsphere.datacenters {\n  hosts {\n    vmknics.where(tags == \"Management\") {\n      properties['Enabled'] == true\n      properties['MTU'] == 1500\n      properties['VDSName'] != /(?i)storage/\n    }\n  }\n}\n",
				},
			},
		},
		"arista.eos.stp": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			},
		},
		"msgraph.beta.security.securityscore": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"vsphere.license": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"aws.codebuild": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.ec2.instance.device": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.ec2.vpnconnection": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.rds.dbinstance": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.keyvault": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"esxi.ntpconfig": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"gitlab.group": {
			Maturity: "experimental",
		},
		"arista.eos.stp.mst": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			},
		},
		"azurerm.monitor": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"aws.accessAnalyzer": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "The `aws.accessAnalyzer` resource returns a list of `aws.accessAnalyzer.analyzer` objects representing all of the AWS IAM Access Analyzers configured across the AWS account.\n",
			}, Refs: []docs.LrDocsRefs{
				{
					Title: "Using AWS IAM Access Analyzer",
					Url:   "https://docs.aws.amazon.com/IAM/latest/UserGuide/what-is-access-analyzer.html",
				},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Return a list of `aws.accessAnalyzer.analyzer` objects representing all of the AWS IAM Access Analyzers configured across the AWS account",
					Query: "aws.accessAnalyzer.analyzers",
				},
			},
		},
		"aws.config": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.elb": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.keyvault.certificate": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.network.watcher": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"gitlab.project": {
			Maturity: "experimental",
		},
		"k8s.pod": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"kubernetes"},
			},
		},
		"kernel": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "List all kernel modules",
					Query: "kernel.modules { name loaded size }",
				},
				{
					Title: "List all loaded kernel modules",
					Query: "kernel.modules.where( loaded == true ) { name }",
				},
			},
		},
		"yum.repo": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Check if a yum repo is enabled",
					Query: "yum.repo('salt-latest') {\n  enabled\n}\n",
				},
			},
		},
		"aws.config.recorder": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.emr": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.postgresql.database": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"aws.cloudwatch.loggroup": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.codebuild.project": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"ms365.teams": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"arista.eos.runningConfig": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			},
		},
		"azurerm.sql.server": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"aws.securityhub": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.keyvault.vault": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"k8s.node": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"kubernetes"},
			},
		},
		"arista.eos.role": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			},
		},
		"aws.dynamodb.limit": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.mariadb.database": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"arista.eos.ntpSetting": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			},
		},
		"aws.s3.bucket": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "Buckets and objects are AWS resources, and Amazon S3 provides APIs for you to manage them\n",
			}, Refs: []docs.LrDocsRefs{
				{
					Title: "Amazon S3 Product Page",
					Url:   "https://aws.amazon.com/s3/",
				},
				{
					Title: "AWS Documentation: Buckets overview",
					Url:   "https://docs.aws.amazon.com/AmazonS3/latest/userguide/UsingBucket.html",
				},
			},
		},
		"aws.secretsmanager.secret": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"esxi": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Check that all kernel modules are signed",
					Query: "esxi.host {\n  kernelModules {\n    signedStatus == \"Signed\"\n  }\n}\n",
				},
			},
		},
		"aws.ec2.image": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.ec2.volume": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.redshift.cluster": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.s3.bucket.corsrule": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.network.securitygroup": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"dns.mxRecord": {
			Maturity: "experimental",
		},
		"k8s.apiresource": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"kubernetes"},
			},
		},
		"package": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Check if a package is installed",
					Query: "package('git').installed",
				},
			},
		},
		"terraform.module": {
			Maturity: "experimental",
			Platform: &docs.LrDocsPlatform{
				Name: []string{"terraform"},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Display all loaded Terraform modules",
					Query: "terraform.modules { key version source}",
				},
			},
		},
		"aws.vpc.routetable": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"windows.firewall": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Check a specific Windows Firewall rule",
					Query: "windows.firewall.rules.where ( displayName == \"File and Printer Sharing (Echo Request - ICMPv4-In)\") {\n  enabled == 1\n}\n",
				},
			},
		},
		"arista.eos.interface": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			},
		},
		"arista.eos.snmpSetting": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			},
		},
		"aws.efs.filesystem": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.sql.configuration": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azuread.serviceprincipal": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.storage.account": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"groups": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Ensure the user is not part of group",
					Query: "groups.where(name == 'wheel').list { members.all( name != 'username') }",
				},
			},
		},
		"msgraph.beta.domaindnsrecord": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"aws.guardduty": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.lambda": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.s3": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "Amazon Simple Storage Service (Amazon S3) is an object storage service\n",
			}, Refs: []docs.LrDocsRefs{
				{
					Title: "Amazon S3 Product Page",
					Url:   "https://aws.amazon.com/s3/",
				},
				{
					Title: "AWS Documentation: What is Amazon S3?",
					Url:   "https://docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html",
				},
			},
		},
		"aws.ec2.vgwtelemetry": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"users": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Display all users and their uid",
					Query: "users.list { uid name }",
				},
				{
					Title: "Ensure user exists",
					Query: "users.one( name == 'root')",
				},
				{
					Title: "Ensure user does not exist",
					Query: "users.none(name == \"vagrant\")",
				},
				{
					Title: "Search for a specific SID and check for its values",
					Query: "users.where( sid == /S-1-5-21-\\d+-\\d+-\\d+-501/ ).list {\n  name != \"Guest\"\n}\n",
				},
			},
		},
		"aws.cloudtrail": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.cloudwatch.metric": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.ec2.networkacl": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.sagemaker.notebookinstance": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.mysql": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.storage.container": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"file": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Test if a directory exists",
					Query: "file('/etc') {\n  exists\n  permissions.isDirectory\n}\n",
				},
			},
		},
		"msgraph.beta.rolemanagement.roledefinition": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"terraform.file": {
			Maturity: "experimental",
			Platform: &docs.LrDocsPlatform{
				Name: []string{"terraform"},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Display all files and their blocks",
					Query: "terraform.files { path blocks { nameLabel } }",
				},
			},
		},
		"aws.cloudwatch.loggroup.metricsfilter": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.iam.group": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.rds.dbcluster": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.sns": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.kms.key": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.lambda.function": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.secretsmanager": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.acm": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "The `aws.acm` resource can be used to assess the configuration of the AWS Certificates Manager service within the account. This resource returns a list of `aws.acm.certificate` objects for all ACM certificates found within the account.",
			}, Refs: []docs.LrDocsRefs{
				{
					Title: "What Is AWS Certificate Manager?",
					Url:   "https://docs.aws.amazon.com/acm/latest/userguide/acm-overview.html",
				},
				{
					Title: "Security in AWS Certificate Manager",
					Url:   "https://docs.aws.amazon.com/acm/latest/userguide/security.html",
				},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Return a list of `aws.acm.certificates` within the AWS account",
					Query: "aws.acm.certificates",
				},
				{
					Title: "Return a list of `aws.acm.certificates` within the AWS account along with values for specified fields",
					Query: "aws.acm.certificates {\narn\nnotBefore\nnotAfter \ncreatedAt\ndomainName\nstatus\nsubject\ncertificate() \n}\n",
				},
				{
					Title: "Checks whether ACM Certificates in your account are marked for expiration within 90 days",
					Query: "aws.acm.certificates.\n  where( status != /PENDING_VALIDATION/ ).\n  all (notAfter - notBefore <= 90 * time.day)\n",
				},
			},
		},
		"k8s.job": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"kubernetes"},
			},
		},
		"msgraph.beta.security": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"aws.ec2.internetgateway": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.es": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"dns": {
			Maturity: "experimental",
		},
		"aws.guardduty.detector": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.vpc": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"dns.record": {
			Maturity: "experimental",
		},
		"k8s.namespace": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"kubernetes"},
			},
		},
		"msgraph.beta.user": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"powershell": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Run custom powershell command",
					Query: "powershell('Get-WmiObject -Class Win32_volume -Filter \"DriveType=3\"| Select Label') {\n  stdout == /PAGEFILE/\n  stderr == ''\n}\n",
				},
				{
					Title: "Check the timezone",
					Query: "powershell('tzutil /g') {\n  stdout.trim == 'GMT Standard Time'\n  stderr == ''\n}\n",
				},
			},
		},
		"aws.vpc.flowlog": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"equinix.metal.project": {
			Maturity: "experimental",
			Platform: &docs.LrDocsPlatform{
				Name: []string{"equinix"},
			},
		},
		"vsphere.cluster": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"aws.cloudwatch.metricsalarm": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.ec2.networkacl.entry": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.sns.subscription": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azuread.group": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.resource": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.sql.firewallrule": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"ipmi.chassis": {
			Maturity: "experimental",
		},
		"os": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Show all environment variables",
					Query: "os.env",
				},
				{
					Title: "Retrieve a single environment variable",
					Query: "os.env['windir']",
				},
			},
		},
		"aws.account": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "The `aws.account` resource provides configuration for AWS accounts including the account number, and configured aliases.\n",
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Return the account id (number) and any configured account aliases",
					Query: "aws.account { id aliases }",
				},
			},
		},
		"aws.elasticache": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.iam.user": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"sshd.config": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Check the ssh banner setting",
					Query: "sshd.config.params['Banner'] == '/etc/ssh/sshd-banner'",
				},
			},
		},
		"aws.apigateway.restapi": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "The `aws.apigateway.restapi` resource provides an object representing an individual REST API configured within the AWS account. For usage see the `aws.apigateway` resource documentation.\n",
			},
		},
		"aws.autoscaling.group": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.cloudwatch": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.dynamodb.globaltable": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.ec2": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.redshift": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.monitor.logprofile": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"aws.iam.role": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"msgraph.beta.domain": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"aws.autoscaling": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.kms": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"k8s": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"kubernetes"},
			},
		},
		"github.organization": {
			Maturity: "experimental",
		},
		"k8s.deployment": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"kubernetes"},
			},
		},
		"msgraph.beta.rolemanagement": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"aws.apigateway": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "The `aws.apigateway` resource can be used to assess the configuration of the AWS API Gateway service.",
			}, Refs: []docs.LrDocsRefs{
				{
					Title: "What is Amazon API Gateway?",
					Url:   "https://docs.aws.amazon.com/apigateway/latest/developerguide/welcome.html",
				},
				{
					Title: "Security in Amazon API Gateway",
					Url:   "https://docs.aws.amazon.com/apigateway/latest/developerguide/security.html",
				},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Return a list of `aws.apigateway.restapi` objects for all REST APIs configured with the account across all enabled regions",
					Query: "aws.apigateway.restApis",
				},
				{
					Title: "Return a list of `aws.apigateway.restapi` objects for all REST APIs configured with the account across all enabled regions and the value for specified fields",
					Query: "aws.apigateway.restApis {\n  createdDate\n  description\n  stages\n  region\n  arn\n  id\n  name\n}\n",
				},
				{
					Title: "Checks that all methods in Amazon API Gateway have caching enabled and encrypted",
					Query: "aws.apigateway.restApis.all(stages.all(\n  methodSettings['CachingEnabled'] == true && \n    methodSettings['CacheDataEncrypted'] == true\n))\n",
				},
				{
					Title: "Checks that all methods in Amazon API Gateway have logging enabled",
					Query: "aws.apigateway.restApis.all(stages.all(\nmethodSettings['LoggingLevel'] == \"ERROR\" || methodSettings['LoggingLevel'] == \"INFO\"\n))\n",
				},
			},
		},
		"aws.apigateway.stage": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "The `aws.apigateway.stage` resource provides an object representing an individual stage configured on a REST API. For usage see the `aws.apigateway` resource documentation.\n",
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Return a list of AWS API Gateway REST APIs configured across all enabled regions in the AWS account and the values for the arn and stages",
					Query: "aws.apigateway.restApis { \n  arn \n  stages \n}\n",
				},
			},
		},
		"aws.iam.policy": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"msgraph.beta.devicemanagement": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"vsphere.host": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Verify the Software AcceptanceLevel for ESXi host",
					Query: "// targeting a single ESXi host\nesxi.host {\n  acceptanceLevel == 'VMwareCertified' || acceptanceLevel == 'VMwareAccepted' || acceptanceLevel == 'PartnerSupported'\n}\n\n// targeting the vSphere API\nvsphere.datacenters {\n  hosts {\n    acceptanceLevel == 'VMwareCertified' || acceptanceLevel == 'VMwareAccepted' || acceptanceLevel == 'PartnerSupported'\n  }\n}\n",
				},
				{
					Title: "Verify that each vib is \"VMwareCertified\" or \"VMwareAccepted\"",
					Query: "esxi.host {\n  packages {\n    acceptanceLevel == 'VMwareCertified' || acceptanceLevel == 'VMwareAccepted' || acceptanceLevel == 'PartnerSupported'\n  }\n}\n",
				},
			},
		},
		"arista.eos.runningConfig.section": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			},
		},
		"aws.cloudtrail.trail": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.ec2.securitygroup.ippermission": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.sql.databaseusage": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"platform": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Platform Name and Release",
					Query: "platform { name release }",
				},
			},
		},
		"azurerm.web.appsiteconfig": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"msgraph.beta": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"aws.iam": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.securityhub.hub": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.mysql.server": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"terraform.fileposition": {
			Maturity: "experimental",
			Platform: &docs.LrDocsPlatform{
				Name: []string{"terraform"},
			},
		},
		"vsphere.vm": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"arista.eos.spt.mstInterface": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Verify the Spanning-Tree Version",
					Query: "arista.eos.stp.mstInstances {\n protocol == \"mstp\"\n}\n",
				},
			},
		},
		"aws.s3control": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"auditpol": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "List all audit policies",
					Query: "auditpol.list { inclusionsetting exclusionsetting  subcategory }",
				},
				{
					Title: "Check a specific auditpol configuration",
					Query: "auditpol.where(subcategory == 'Sensitive Privilege Use').list {\n  inclusionsetting == 'Success and Failure'\n}\n",
				},
			},
		},
		"aws.acm.certificate": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "The `aws.acm.certificate` resource provides an object representing an individual AWS ACM certificate. For usage see the `aws.acm` documentation.\n",
			},
		},
		"msgraph.beta.organization": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"msgraph.beta.policies": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"terraform.block": {
			Maturity: "experimental",
			Platform: &docs.LrDocsPlatform{
				Name: []string{"terraform"},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Display all Terraform blocks and their arguments",
					Query: "terraform.blocks { nameLabel arguments }",
				},
			},
		},
		"aws.ec2.snapshot": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.emr.cluster": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.sagemaker.endpoint": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azuread.domain": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.network.securityrule": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"equinix.metal.sshkey": {
			Maturity: "experimental",
			Platform: &docs.LrDocsPlatform{
				Name: []string{"equinix"},
			},
		},
		"azuread.user": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.mariadb.server": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.web": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.mariadb": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"esxi.timezone": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"k8s.cronjob": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"kubernetes"},
			},
		},
		"azurerm.postgresql.server": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"esxi.vib": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"k8s.container": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"kubernetes"},
			},
		},
		"aws.sagemaker": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azuread": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"msgraph.beta.rolemanagement.roleassignment": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"vsphere.vswitch.standard": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"aws.iam.usercredentialreportentry": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"windows": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Check the OS Edition",
					Query: "windows.computerInfo['WindowsInstallationType'] == 'Server Core'",
				},
			},
		},
		"aws.dynamodb": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.web.appsiteauthsettings": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"ms365.exchangeonline": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"msgraph.beta.application": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"registrykey.property": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Verify a registry key property",
					Query: "registrykey.property(path: 'HKEY_LOCAL_MACHINE\\Software\\Policies\\Microsoft\\Windows\\EventLog\\System', name: 'MaxSize') {\n  value >= 32768\n}\n",
				},
			},
		},
		"arista.eos": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Arista EOS Version",
					Query: "arista.eos.version",
				},
				{
					Title: "Verify that Arista EOS Version is 18.x",
					Query: "arista.eos.version['version'] == /18\\./",
				},
				{
					Title: "Display EOS interfaces",
					Query: "arista.eos.interfaces { name mtu bandwidth status }",
				},
				{
					Title: "Display all connected EOS interfaces",
					Query: "arista.eos.interfaces.where ( status['linkStatus'] == \"connected\") {  name mtu bandwidth stauts}",
				},
				{
					Title: "EOS Hostname",
					Query: "arista.eos.hostname",
				},
			},
		},
		"arista.eos.ipInterface": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			},
		},
		"azurerm.web.appsite": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"ipmi": {
			Maturity: "experimental",
		},
		"vsphere.datacenter": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"aws.elb.loadbalancer": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"equinix.metal.device": {
			Maturity: "experimental",
			Platform: &docs.LrDocsPlatform{
				Name: []string{"equinix"},
			},
		},
		"ms365.sharepointonline": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"azurerm.postgresql": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.sql.server.administrator": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"aws.accessanalyzer.analyzer": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.efs": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.rds": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.s3.bucket.policy": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "Bucket policies grant permission to your Amazon S3 resources\n",
			}, Refs: []docs.LrDocsRefs{
				{
					Title: "AWS Documentation: Using bucket policies",
					Url:   "https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucket-policies.html",
				},
			},
		},
		"azurerm.mysql.database": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"mount": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "List all mount points",
					Query: "mount.list { path device fstype options }",
				},
				{
					Title: "Ensure the mountpoint exists",
					Query: "mount.one( path == \"/\" )",
				},
				{
					Title: "Check mountpoint configuration",
					Query: "mount.where( path == \"/\" ).list {\n  device == '/dev/mapper/vg00-lv_root'\n  fstype == 'xfs'\n  options['rw'] != null\n  options['relatime'] != null\n  options['seclabel'] != null\n  options['attr2'] != null\n  options['inode64'] != null\n  options['noquota'] != null\n}\n",
				},
			},
		},
		"msgraph.beta.devicemanagement.deviceconfiguration": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"aws.ec2.networkacl.entry.portrange": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.network": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"windows.feature": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Check that a Windows features is installed",
					Query: "windows.feature('SNMP-Service').installed",
				},
				{
					Title: "Check that a specific feature is not installed",
					Query: "windows.feature('Windows-Defender').installed == false",
				},
			},
		},
		"aws.ec2.instance": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.sns.topic": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.compute.disk": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"equinix.metal.user": {
			Maturity: "experimental",
			Platform: &docs.LrDocsPlatform{
				Name: []string{"equinix"},
			},
		},
		"vsphere.vmnic": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"vsphere.vswitch.dvs": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"arista.eos.user": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"arista-eos"},
			},
		},
		"aws.s3.bucket.grant": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.compute.vm": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"msgraph.beta.devicemanagement.devicecompliancepolicy": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"microsoft365"},
			},
		},
		"aws.rds.snapshot": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.compute": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.sql": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"aws": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			}, Docs: &docs.LrDocsDocumentation{
				Description: "The AWS Resource is used to query AWS accounts and the services and resources within them. The `aws` resource returns a list of enabled regions within the AWS account, as well as a list of `aws.vpc` objects representing all VPCs configured across all enabled regions.\n",
			}, Refs: []docs.LrDocsRefs{
				{
					Title: "AWS Documentation: Managing AWS Regions",
					Url:   "https://docs.aws.amazon.com/general/latest/gr/rande-manage.html",
				},
				{
					Title: "AWS Documentation: Security in Amazon Virtual Private Cloud",
					Url:   "https://docs.aws.amazon.com/vpc/latest/userguide/security.html",
				},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "List all enabled regions within the AWS account",
					Query: "aws.regions",
				},
				{
					Title: "List of `aws.vpc` objects for all VPCs across all enabled regions",
					Query: "aws.vpcs",
				},
				{
					Title: "List of `aws.vpc` objects for all VPCs across all enabled regions and the values for specified fields",
					Query: "aws.vpcs {\n  arn \n  id \n  state \n  isDefault \n  region \n  flowLogs\n  routeTables\n}\n",
				},
				{
					Title: "Ensure VPC flow logging is enabled in all VPCs",
					Query: "aws.vpcs\n  .all(\n    flowLogs.any(status == \"ACTIVE\")\n  )\n",
				},
			},
		},
		"aws.es.domain": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.storage": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"esxi.service": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"azuread.application": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.network.interface": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"secpol": {
			Snippets: []docs.LrDocsSnippet{
				{
					Title: "Check that a specific SID is included in the privilege rights",
					Query: "secpol.privilegerights['SeRemoteShutdownPrivilege'].contains( _ == 'S-1-5-32-544')",
				},
			},
		},
		"terraform": {
			Maturity: "experimental",
			Platform: &docs.LrDocsPlatform{
				Name: []string{"terraform"},
			}, Snippets: []docs.LrDocsSnippet{
				{
					Title: "Display all Terraform blocks and their arguments",
					Query: "terraform.blocks { nameLabel arguments }",
				},
				{
					Title: "Display all data blocks",
					Query: "terraform.datasources { nameLabel arguments }",
				},
				{
					Title: "Display all resource blocks",
					Query: "terraform.resources { nameLabel arguments }",
				},
			},
		},
		"vsphere.vmknic": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"azurerm.keyvault.secret": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"azurerm.sql.database": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"aws.iam.policyversion": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.monitor.diagnosticsetting": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"esxi.kernelmodule": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"vmware-esxi", "vmware-vsphere"},
			},
		},
		"aws.ec2.securitygroup": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"aws.sagemaker.notebookinstance.details": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"aws"},
			},
		},
		"azurerm.keyvault.key": {
			Platform: &docs.LrDocsPlatform{
				Name: []string{"azure"},
			},
		},
		"equinix.metal.organization": {
			Maturity: "experimental",
			Platform: &docs.LrDocsPlatform{
				Name: []string{"equinix"},
			},
		},
	},
}
