// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package resources

import (
	"strconv"

	"go.mondoo.com/cnquery/llx"
	"go.mondoo.com/cnquery/providers-sdk/v1/plugin"
	"go.mondoo.com/cnquery/providers/gitlab/connection"
)

func (g *mqlGitlabGroup) id() (string, error) {
	return "gitlab.group/" + strconv.FormatInt(g.Id.Data, 10), nil
}

// init initializes the gitlab group with the arguments
// see https://docs.gitlab.com/ee/api/groups.html#new-group
func initGitlabGroup(runtime *plugin.Runtime, args map[string]*llx.RawData) (map[string]*llx.RawData, plugin.Resource, error) {
	if len(args) > 2 {
		return args, nil, nil
	}

	conn := runtime.Connection.(*connection.GitLabConnection)
	grp, _, err := conn.Client().Groups.GetGroup(conn.GroupPath, nil)
	if err != nil {
		return nil, nil, err
	}

	args["id"] = llx.IntData(int64(grp.ID))
	args["name"] = llx.StringData(grp.Name)
	args["path"] = llx.StringData(grp.Path)
	args["description"] = llx.StringData(grp.Description)
	args["webURL"] = llx.StringData(string(grp.WebURL))
	args["visibility"] = llx.StringData(string(grp.Visibility))
	args["requireTwoFactorAuthentication"] = llx.BoolData(grp.RequireTwoFactorAuth)
	args["preventForkingOutsideGroup"] = llx.BoolData(grp.PreventForkingOutsideGroup)
	args["mentionsDisabled"] = llx.BoolData(grp.MentionsDisabled)
	args["emailsDisabled"] = llx.BoolData(grp.EmailsDisabled)

	return args, nil, nil
}

// GetProjects list all projects that belong to a group
// see https://docs.gitlab.com/ee/api/projects.html
func (g *mqlGitlabGroup) projects() ([]interface{}, error) {
	conn := g.MqlRuntime.Connection.(*connection.GitLabConnection)

	if g.Path.Error != nil {
		return nil, g.Path.Error
	}
	path := g.Path.Data

	grp, _, err := conn.Client().Groups.GetGroup(path, nil)
	if err != nil {
		return nil, err
	}

	var mqlProjects []interface{}
	for i := range grp.Projects {
		prj := grp.Projects[i]

		mqlProject, err := CreateResource(g.MqlRuntime, "gitlab.project", map[string]*llx.RawData{
			"allowMergeOnSkippedPipeline": llx.BoolData(prj.AllowMergeOnSkippedPipeline),
			"archived":                    llx.BoolData(prj.Archived),
			"autoDevopsEnabled":           llx.BoolData(prj.AutoDevopsEnabled),
			"containerRegistryEnabled":    llx.BoolData(prj.ContainerRegistryEnabled),
			"createdAt":                   llx.TimeDataPtr(prj.CreatedAt),
			"defaultBranch":               llx.StringData(prj.DefaultBranch),
			"description":                 llx.StringData(prj.Description),
			"emailsDisabled":              llx.BoolData(!prj.EmailsEnabled),
			"fullName":                    llx.StringData(prj.NameWithNamespace),
			"id":                          llx.IntData(int64(prj.ID)),
			"issuesEnabled":               llx.BoolData(prj.IssuesEnabled),
			"mergeRequestsEnabled":        llx.BoolData(prj.MergeRequestsEnabled),
			"mirror":                      llx.BoolData(prj.Mirror),
			"name":                        llx.StringData(prj.Name),
			"onlyAllowMergeIfAllDiscussionsAreResolved": llx.BoolData(prj.OnlyAllowMergeIfAllDiscussionsAreResolved),
			"onlyAllowMergeIfPipelineSucceeds":          llx.BoolData(prj.OnlyAllowMergeIfPipelineSucceeds),
			"packagesEnabled":                           llx.BoolData(prj.PackagesEnabled),
			"path":                                      llx.StringData(prj.Path),
			"requirementsEnabled":                       llx.BoolData(prj.RequirementsEnabled),
			"serviceDeskEnabled":                        llx.BoolData(prj.ServiceDeskEnabled),
			"snippetsEnabled":                           llx.BoolData(prj.SnippetsEnabled),
			"visibility":                                llx.StringData(string(prj.Visibility)),
			"webURL":                                    llx.StringData(prj.WebURL),
			"wikiEnabled":                               llx.BoolData(prj.WikiEnabled),
		})
		if err != nil {
			return nil, err
		}
		mqlProjects = append(mqlProjects, mqlProject)
	}

	return mqlProjects, nil
}

func (g *mqlGitlabProject) id() (string, error) {
	return "gitlab.project/" + strconv.FormatInt(g.Id.Data, 10), nil
}
