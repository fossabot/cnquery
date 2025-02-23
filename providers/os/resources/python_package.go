// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package resources

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"go.mondoo.com/cnquery/llx"
	"go.mondoo.com/cnquery/providers-sdk/v1/plugin"
	"go.mondoo.com/cnquery/providers/os/connection/shared"
)

func (k *mqlPythonPackage) id() (string, error) {
	file := k.GetFile()
	if file.Error != nil {
		return "", file.Error
	}

	mqlFile := file.Data
	metadataPath := mqlFile.Path.Data
	return metadataPath, nil
}

func initPythonPackage(runtime *plugin.Runtime, args map[string]*llx.RawData) (map[string]*llx.RawData, plugin.Resource, error) {
	if len(args) > 1 {
		return args, nil, nil
	}
	if x, ok := args["path"]; ok {
		path, ok := x.Value.(string)
		if !ok {
			return nil, nil, errors.New("Wrong type for 'path' in python.package initialization, it must be a string")
		}

		file, err := CreateResource(runtime, "file", map[string]*llx.RawData{
			"path": llx.StringData(path),
		})
		if err != nil {
			return nil, nil, err
		}
		args["file"] = llx.ResourceData(file, "file")

		delete(args, "path")
	}
	return args, nil, nil
}

func (k *mqlPythonPackage) name() (string, error) {
	k.populateData()
	return k.Name.Data, nil
}

func (k *mqlPythonPackage) version() (string, error) {
	k.populateData()
	return k.Version.Data, nil
}

func (k *mqlPythonPackage) license() (string, error) {
	k.populateData()
	return k.License.Data, nil
}

func (k *mqlPythonPackage) author() (string, error) {
	k.populateData()
	return k.Author.Data, nil
}

func (k *mqlPythonPackage) summary() (string, error) {
	k.populateData()
	return k.Summary.Data, nil
}

func (k *mqlPythonPackage) dependencies() ([]interface{}, error) {
	k.populateData()
	return k.Dependencies.Data, nil
}

func (k *mqlPythonPackage) populateData() error {
	file := k.GetFile()
	if file.Error != nil {
		return file.Error
	}
	mqlFile := file.Data
	conn := k.MqlRuntime.Connection.(shared.Connection)
	afs := &afero.Afero{Fs: conn.FileSystem()}
	metadataPath := mqlFile.Path.Data
	ppd, err := parseMIME(afs, metadataPath)
	if err != nil {
		return fmt.Errorf("error parsing python package data: %s", err)
	}

	deps := make([]interface{}, len(ppd.dependencies))
	for i, dep := range ppd.dependencies {
		deps[i] = dep
	}

	k.Name = plugin.TValue[string]{Data: ppd.name, State: plugin.StateIsSet}
	k.Version = plugin.TValue[string]{Data: ppd.version, State: plugin.StateIsSet}
	k.Author = plugin.TValue[string]{Data: ppd.author, State: plugin.StateIsSet}
	k.Summary = plugin.TValue[string]{Data: ppd.summary, State: plugin.StateIsSet}
	k.License = plugin.TValue[string]{Data: ppd.license, State: plugin.StateIsSet}
	k.Dependencies = plugin.TValue[[]interface{}]{Data: deps, State: plugin.StateIsSet}

	return nil
}
