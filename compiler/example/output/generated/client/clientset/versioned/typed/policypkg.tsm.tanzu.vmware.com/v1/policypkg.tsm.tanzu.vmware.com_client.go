/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/apis/policypkg.tsm.tanzu.vmware.com/v1"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/client/clientset/versioned/scheme"

	rest "k8s.io/client-go/rest"
)

type PolicypkgTsmV1Interface interface {
	RESTClient() rest.Interface
	ACPConfigsGetter
	AccessControlPoliciesGetter
	AdditionalPolicyDatasGetter
	RandomPolicyDatasGetter
	VMpoliciesGetter
}

// PolicypkgTsmV1Client is used to interact with features provided by the policypkg.tsm.tanzu.vmware.com group.
type PolicypkgTsmV1Client struct {
	restClient rest.Interface
}

func (c *PolicypkgTsmV1Client) ACPConfigs() ACPConfigInterface {
	return newACPConfigs(c)
}

func (c *PolicypkgTsmV1Client) AccessControlPolicies() AccessControlPolicyInterface {
	return newAccessControlPolicies(c)
}

func (c *PolicypkgTsmV1Client) AdditionalPolicyDatas() AdditionalPolicyDataInterface {
	return newAdditionalPolicyDatas(c)
}

func (c *PolicypkgTsmV1Client) RandomPolicyDatas() RandomPolicyDataInterface {
	return newRandomPolicyDatas(c)
}

func (c *PolicypkgTsmV1Client) VMpolicies() VMpolicyInterface {
	return newVMpolicies(c)
}

// NewForConfig creates a new PolicypkgTsmV1Client for the given config.
func NewForConfig(c *rest.Config) (*PolicypkgTsmV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &PolicypkgTsmV1Client{client}, nil
}

// NewForConfigOrDie creates a new PolicypkgTsmV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *PolicypkgTsmV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new PolicypkgTsmV1Client for the given RESTClient.
func New(c rest.Interface) *PolicypkgTsmV1Client {
	return &PolicypkgTsmV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *PolicypkgTsmV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}