package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aunum/log"
	apiErrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/client-go/util/homedir"
)

type Resource struct {
	Action   dynamic.ResourceInterface
	Unstruct *unstructured.Unstructured
	Mapping  *meta.RESTMapping
	Ns       string
	Name     string
}

var TenantName string

func GenerateContext(host string) (*kubernetes.Clientset, dynamic.Interface, http.RoundTripper, spdy.Upgrader, string, error) {
	var config *rest.Config
	var err error
	if host == "" {
		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, nil, nil, nil, "", err
		}
	} else {
		config = &rest.Config{
			Host: host,
		}
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, "", err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, "", err
	}

	roundTripper, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return nil, nil, nil, nil, "", err
	}
	return clientset, dynamicClient, roundTripper, upgrader, config.Host, err
}

func GroupResources(client kubernetes.Interface) ([]*restmapper.APIGroupResources, error) {
	groupResources, err := restmapper.GetAPIGroupResources(client.Discovery())
	if err != nil {
		m := fmt.Errorf("failed to get API group resources: %s", err)
		log.Errorf(m.Error())
		return nil, m
	}

	return groupResources, nil
}

func Decode(doc string) (*unstructured.Unstructured, *schema.GroupVersionKind, error) {
	log.Debugf("Decoding YAML document into an unstructured")
	unstruct := &unstructured.Unstructured{}
	_, gvk, err := scheme.Codecs.UniversalDeserializer().Decode([]byte(doc), nil, unstruct)
	if err != nil {
		m := fmt.Errorf("error deserializing data: %s", err)
		log.Errorf(m.Error())
		return nil, nil, m
	}

	return unstruct, gvk, nil
}

func Search(unstruct *unstructured.Unstructured, fields ...string) (interface{}, bool, error) {
	log.Debugf("Searching metadata for fields %v", fields)
	field, found, err := unstructured.NestedFieldCopy(unstruct.Object, fields...)
	if err != nil {
		m := fmt.Errorf("malformed yaml: %s", err)
		log.Errorf(m.Error())
		return "", found, m
	}

	return field, found, nil
}

func RestMapping(resources []*restmapper.APIGroupResources, gvk *schema.GroupVersionKind) (*meta.RESTMapping, error) {
	mapping, err := restmapper.NewDiscoveryRESTMapper(resources).RESTMapping(schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}, gvk.Version)
	if err != nil {
		m := fmt.Errorf("rest mapping failed: %s", err)
		log.Errorf(m.Error())
		return nil, m
	}

	return mapping, nil
}

func Apply(resources ...*Resource) error {
	for _, resource := range resources {
		fmt.Printf("trying to create object %s/%s on namespace %s\n", resource.Mapping.Resource, resource.Name, resource.Ns)
		_, createErr := resource.Action.Create(context.Background(), resource.Unstruct, metav1.CreateOptions{})
		if createErr != nil {
			if apiErrs.IsAlreadyExists(createErr) {
				fmt.Printf("\u2713 skipping creation resource already exists %v for %s\n", resource.Mapping.Resource, resource.Name)
			} else {
				fmt.Printf("error creating resource in the cluster: %s\n", createErr)
				return createErr
			}
		}
		fmt.Printf("\u2713 created object %s/%s on namespace %s\n", resource.Mapping.Resource, resource.Name, resource.Ns)
	}
	return nil
}

func dynamicResourceMapper(resources []*restmapper.APIGroupResources, dynamicClient dynamic.Interface, doc string, namespace string) (*Resource, error) {
	var ns string
	unstruct, gvk, err := Decode(doc)
	if err != nil {
		return nil, err
	}
	kind, kindfound, err := Search(unstruct, "kind")
	if err != nil {
		return nil, err
	}

	// Not all resources will have or work with a namespace.
	if kindfound {
		if kind.(string) == "CustomResourceDefinition" || kind.(string) == "ClusterRole" || kind.(string) == "ClusterRoleBinding" {
			ns = ""
		} else {
			ns = namespace
			(*unstruct).SetNamespace(namespace)
		}
	}

	name, found, err := Search(unstruct, "metadata", "name")
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.New("failed to find name field")
	}

	mapping, err := RestMapping(resources, gvk)
	if err != nil {
		return nil, err
	}

	// If the namespace returned to us is empty then the action call will work in a non-namespaced context.
	// Think of resources that are cluster-wide and not namespace specific.
	// If we have a non-empty namespace here then we'll take an action within that namespace.
	return &Resource{
		Action:   dynamicClient.Resource(mapping.Resource).Namespace(ns),
		Mapping:  mapping,
		Ns:       ns,
		Name:     name.(string),
		Unstruct: unstruct,
	}, nil

}

func runInstall(file, host, namespace string) error {
	clientset, dynamicClient, _, _, _, err := GenerateContext(host)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	stringData := string(data)
	resources := make([]*Resource, 0)
	groupResources, err := GroupResources(clientset)
	if err != nil {
		return err
	}
	for _, doc := range strings.Split(stringData, "---") {
		//fmt.Printf("Installing Object: %s", doc)
		if doc != "" {
			resource, err := dynamicResourceMapper(groupResources, dynamicClient, doc, namespace)
			if err != nil {
				return err
			}
			resources = append(resources, resource)
		}

	}
	err = Apply(resources...)
	if err != nil {
		return err
	}
	return nil
}
func InstallObject(file, host, namespace string) error {
	filestat, err := os.Stat(file)
	if err != nil {
		return err
	}
	switch filestat.IsDir(); filestat.IsDir() {
	case true:
		files, err := ioutil.ReadDir(file)
		if err != nil {
			return err
		}
		for _, newfile := range files {
			err = InstallObject(filepath.Join(file, newfile.Name()), host, namespace)
			if err != nil {
				return err
			}
		}
	case false:
		err = runInstall(file, host, namespace)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckAPIServer(hostName string, pollTimeout time.Duration) error {
	waitTimer := time.NewTimer(pollTimeout * time.Minute)
	for {
		select {
		case <-waitTimer.C:
			err := fmt.Errorf("timeout in bringing up apisever")
			return err
		default:
			resp, err := http.Get(fmt.Sprintf("%s/api/v1/namespaces", hostName))
			if err != nil {
				break
			}
			if resp.StatusCode != 200 {
				break
			}
			resp, err = http.Get(fmt.Sprintf("%s/apis/apiextensions.k8s.io/v1/customresourcedefinitions", hostName))
			if err != nil {
				break
			}
			if resp.StatusCode != 200 {
				break
			}
			return nil
		}
	}
}
func GetLoadBalancer(clientset kubernetes.Interface, pollTimeout time.Duration, TenantName string, LoadBalancerName string) (string, error) {
	var lbIp string = ""
	waitTimer := time.NewTimer(pollTimeout * time.Minute)
	for {
		select {
		case <-waitTimer.C:
			err := fmt.Errorf("timeout in checking if apiserver has loadbalancer ip.")
			return "", err
		default:
			svc, err := clientset.CoreV1().Services(TenantName).Get(context.Background(), LoadBalancerName, metav1.GetOptions{})
			if err != nil {
				break
			}
			print(".")
			if len(svc.Status.LoadBalancer.Ingress) != 0 {
				for _, section := range svc.Status.LoadBalancer.Ingress {
					if section.Hostname != "" {
						lbIp = section.Hostname
					} else {
						if section.IP != "" {
							lbIp = section.IP
						}
					}
				}
				if lbIp == "" {
					break
				}
				_, err := http.Get(fmt.Sprintf("http://%s/api/v1/namespaces", lbIp))
				if err != nil {
					break
				}
				fmt.Printf("\nloadbalancer ip to connect to local apiserver: ...........%s\n", lbIp)
				return lbIp, nil
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func GetFreePort() (string, error) {
	var localPort string
	for i := 49152; i <= 65535; i++ {
		timeout := time.Second
		conn, _ := net.DialTimeout("tcp", net.JoinHostPort("localhost", fmt.Sprintf("%d", i)), timeout)
		if conn != nil {
			defer conn.Close()
			break
		}
		localPort = fmt.Sprintf("%d", i)
		return localPort, nil
	}
	return "", errors.New("could not find free port for port-forwarding...")
}

func CreatePortForwarder(k8shostname, port, podName, TenantName string, stopChan, readyChan chan struct{}, out, errOut io.Writer, roundTripper http.RoundTripper, upgrader spdy.Upgrader) (*portforward.PortForwarder, error) {
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", TenantName, podName)
	hostIP := strings.TrimLeft(k8shostname, "htps:/")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)
	ports := []string{fmt.Sprintf("%s:8001", port)}
	forwarder, err := portforward.New(dialer, ports, stopChan, readyChan, out, errOut)
	if err != nil {
		err = fmt.Errorf("failure in setting up port-forwarding due to: %s", err)
		return nil, err
	}
	return forwarder, nil

}

func GetPodRunning(clientset kubernetes.Interface, TenantName, podS string) (string, error) {
	var podName string
	var pollTimeout time.Duration = 10
	waitTimer := time.NewTimer(pollTimeout * time.Minute)
	for {
		select {
		case <-waitTimer.C:
			err := fmt.Errorf("timeout in checking if pod is up %s.", podS)
			return "", err
		default:
			pods, err := clientset.CoreV1().Pods(TenantName).List(context.Background(), metav1.ListOptions{LabelSelector: fmt.Sprintf("app=%s", podS), FieldSelector: "status.phase=Running"})
			if err != nil {
				break
			}
			if len(pods.Items) == 0 {
				break
			}
			for _, pod := range pods.Items {
				podName = pod.Name
			}
			return podName, nil
		}
	}
}
