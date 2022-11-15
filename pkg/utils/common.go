package utils

import (
	"api-gw/pkg/authn"
	"api-gw/pkg/client"
	"api-gw/pkg/config"
	"api-gw/pkg/envoy"
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const DISPLAY_NAME_LABEL = "nexus/display_name"

func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsServerConfigValid(conf *config.Config) bool {
	if conf != nil {
		if conf.Server.Address != "" && conf.Server.CertPath != "" && conf.Server.KeyPath != "" {
			return true
		}
	}
	return false
}

func DumpReq(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Warn(err)
	}
	fmt.Println(string(requestDump))
}

func GetEnvoyInitParams() (*envoy.JwtAuthnConfig, map[string]*envoy.UpstreamConfig, map[string]*envoy.HeaderMatchedUpstream, error) {
	var jwt *envoy.JwtAuthnConfig
	jwts, err := client.NexusClient.Authentication().ListOIDCs(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Errorln(err)
		return nil, nil, nil, fmt.Errorf("failed to fetch OIDCs: %s", err)
	} else {
		if jwts != nil && len(jwts) > 0 {
			if len(jwts) > 1 {
				return nil, nil, nil, fmt.Errorf("more than 1 oidc objects found")
			}
			var issuer string
			issuer, err = authn.GetIssuer(jwts[0])
			if err != nil {
				log.Errorln(err)
				return nil, nil, nil, fmt.Errorf("failed to get issuer: %s", err)
			}

			var jwksUri string
			jwksUri, err = authn.GetJwksUri(jwts[0])
			if err != nil {
				log.Errorln(err)
				return nil, nil, nil, fmt.Errorf("failed to get jwks_uri: %s", err)
			}

			var callbackEndpoint string
			callbackEndpoint, err = authn.GetCallbackEndpoint(jwts[0])
			if err != nil {
				log.Errorln(err)
				return nil, nil, nil, fmt.Errorf("failed to get callback endpoint: %s", err)
			}

			jwt = &envoy.JwtAuthnConfig{
				IdpName:          jwts[0].Name,
				Issuer:           issuer,
				JwksUri:          jwksUri,
				CallbackEndpoint: callbackEndpoint,
				JwtClaimUsername: jwts[0].Spec.JwtClaimUsername,
			}
		}
	}

	var upstreams = make(map[string]*envoy.UpstreamConfig)
	var headerMatchedUpstreams = make(map[string]*envoy.HeaderMatchedUpstream)
	allUpstreams, err := client.NexusClient.Admin().ListProxyRules(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Errorln(err)
		return nil, nil, nil, fmt.Errorf("failed to get proxyrules: %s", err)
	} else {
		for _, u := range allUpstreams {
			switch u.Spec.MatchCondition.Type {
			case "jwt":
				upstreams[u.Name] = &envoy.UpstreamConfig{
					Name:          u.Name,
					JwtClaimKey:   u.Spec.MatchCondition.Key,
					JwtClaimValue: u.Spec.MatchCondition.Value,
					Host:          u.Spec.Upstream.Host,
					Port:          u.Spec.Upstream.Port,
				}
			case "header":
				headerMatchedUpstreams[u.Name] = &envoy.HeaderMatchedUpstream{
					Name:        u.Name,
					HeaderName:  u.Spec.MatchCondition.Key,
					HeaderValue: u.Spec.MatchCondition.Value,
					Host:        u.Spec.Upstream.Host,
					Port:        u.Spec.Upstream.Port,
				}
			default:
				log.Errorln("invalid proxyrule match condition found")
				return nil, nil, nil, fmt.Errorf("invalid proxyrule match condition found")
			}
		}
	}
	return jwt, upstreams, headerMatchedUpstreams, nil
}

func GetDatamodelName(crdType string) string {
	p, _ := publicsuffix.EffectiveTLDPlusOne(crdType)
	return p
}

func GetCrdType(kind, groupName string) string {
	return GetGroupResourceName(kind) + "." + groupName // eg roots.root.helloworld.com
}

func GetGroupResourceName(kind string) string {
	return strings.ToLower(ToPlural(kind)) // eg roots
}

// GetParentHierarchy constructs the parent in the format <roots.orgchart.vmware.org:default>
func GetParentHierarchy(parents []string, labels map[string]string) (hierarchy []string) {
	for _, parent := range parents {
		for key, val := range labels {
			if parent == key {
				hierarchy = append(hierarchy, key+":"+val)
			}
		}
	}
	return
}

/*
	ConstructGVR constructs group, version, resource for a CRD Type.

Eg: For a given CRD type: roots.vmware.org and ApiVersion: vmware.org/v1,

	      group => vmware.org
		  resource => roots
		  version => v1
*/
func ConstructGVR(crdType string) schema.GroupVersionResource {
	parts := strings.Split(crdType, ".")
	return schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}
}
