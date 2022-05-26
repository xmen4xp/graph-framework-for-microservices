/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validation

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/pkg/util/proto"
)

func ValidateModel(obj interface{}, schema proto.Schema, name string) []error {
	rootValidation, err := itemFactory(proto.NewPath(name), obj)
	if err != nil {
		return []error{err}
	}
	schema.Accept(rootValidation)
	return rootValidation.Errors()
}
