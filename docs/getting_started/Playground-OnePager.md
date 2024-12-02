# Playground TL;DR

Welcome to the Nexus playground TL;DR version.

This tutorial will walk through a minimalistic workflow to try out Nexus playground.

The goal is to get you goin with Nexus in the shortest possible time.

## 1. Clone Nexus Repo
```
git clone git@github.com:intel-sandbox/applications.development.framework.nexus.git
cd applications.development.framework.nexus/
export NEXUS_REPO_DIR=${PWD}
```

## 2. Build & Install Nexus CLI

### 2.1 Build Nexus CLI

### For Linux
```
make cli.build.linux
```
For MacOS
```
make cli.build.darwin
```

### 2.2 Install Nexus CLI
### For Linux
```
sudo make cli.install.linux
```
For MacOS
```
sudo make cli.install.darwin
```

## 3. Build and Install Nexus Runtime

### Build Nexus Runtime

#### Specify a tag to use for locally built runtime artifacts
```
# The tag can be anything. Here we use a tag called "letsplay"
echo letsplay > TAG
```

#### Build Nexus Runtime
```
sudo make runtime.clean; make runtime.build
```

### Install Nexus Runtime

There are 2 options to install a Nexus Runtime:

1. [Runtime based on K0s](#Runtime-on-K0s-K8s-cluster) - this runs docker-compose based workflow to install a minimal K8s stack for use by the runtime
2. [Runtime based on KIND](#Runtime-on-KIND-K8s-cluster) - this is docker based as well but makes use of full fledged K8s stack for use by the runtime

Please pick ONE of the options and stick to it for rest of the playground workflow.

#### Runtime on K0s K8s cluster
```
make runtime.install.k0s
```
#### Runtime on KIND K8s cluster
```
CLUSTER_NAME=<name> CLUSTER_PORT=<starting-port> make runtime.install.kind
```
where

CLUSTER_NAME --> Custom name for the Nexus runtime

CLUSTER_PORT --> Starting port of the range of ports(assume 100 ports) to be used by Nexus runtime.

Some examples

Example 1: Create a runtime called "foo" that can use ports from 8000+ for its runtime
```
CLUSTER_NAME=foo CLUSTER_PORT=8000 make runtime.install.kind
```
Example 2: Create a runtime called "foo" that can use ports from 9000+ for its runtime
```
CLUSTER_NAME=foo CLUSTER_PORT=9000 make runtime.install.kind
```
Example 3: Create a runtime called "bar" that can use ports from 10000+ for its runtime
```
CLUSTER_NAME=bar CLUSTER_PORT=10000 make runtime.install.kind
```

## 4 Setup Workspace for our playground application: Sock Shop
```
export NEXUS_REPO_DIR=${PWD};
mkdir sock-shop-saas-app;cd sock-shop-saas-app;
go mod init demo;
go mod tidy;
go mod edit -replace sockshop=./datamodel;
mkdir datamodel; cd datamodel
```

## 5 Declare data model specification for Sock Shop

### 5.1 Initialize data model 

So we want to call our data model, sockshop.

```
# NOTE: NEXUS_REPO_DIR is the local directory where Nexus repo is clone to.
nexus datamodel init --name sockshop --group sockshop.com --local-dir $NEXUS_REPO_DIR
```

With the above init, the entire code setup needed to define a fully functional data model is now ready.

### 5.2 Write the model specification fo Sock Shop

Copy & paste the below code snippet to file root.to in the datamodel directory

#### File: root.go

```
package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type SockShop struct {
	nexus.SingletonNode

	OrgName  string
	Location string
	Website  string

	Inventory      Socks    `nexus:"children"`
	PO             Orders   `nexus:"children"`
	ShippingLedger Shipping `nexus:"children"`
}

var SocksRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/sock/{root.Socks}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/socks",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:SocksRestAPISpec
type Socks struct {
	nexus.Node

	Brand string
	Color string
	Size  int
}

var OrderRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/order/{root.Orders}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
	},
}

// nexus-rest-api-gen:OrderRestAPISpec
type Orders struct {
	nexus.Node

	SockName string
	Address  string

	Cart     Socks    `nexus:"link"`
	Shipping Shipping `nexus:"link"`
}

type Shipping struct {
	nexus.Node

	TrackingId int
}
```
### 5.3 Build data model for Sock Shop

Nexus compiler can be invoked to build the datamodel with the following command:

```
COMPILER_TAG=letsplay make datamodel_build
```

This will generate all the artifacts needed at install and runtime.

The generated artifacts are available in the $PWD/build directory.

## 6 Install Sock Shop data model on Nexus Runtime

### 6.1 Export KUBECONFIG to Nexus Runtime

The KUBECONFIG to export depends runtime being used in this playgroud.

***Option 1***: If you running a K0s based Nexus runtime:

Run this make target to get the shell export command to execute:
```
make -C $NEXUS_REPO_DIR runtime.k0s.kubeconfig.export
```

***Option 2***: If you running a KIND based Nexus runtime:

Run this make target to get the shell export command to execute:
```
CLUSTER_NAME=<name> make -C $NEXUS_REPO_DIR runtime.kind.kubeconfig.export
```

***NOTE: Remember to execute the printed "export" command on your shell.***

### 6.2 Install data model
```
make dm.install
```

## 7 Access your API's

[Nexus Runtime based on K0s](#Nexus-Runtime-based-on-K0s)

[Nexus Runtime based on KIND](#Nexus-Runtime-based-on-KIND)

### Nexus Runtime based on K0s

#### Kubectl API Access

Lets instantiate sock shop by creating the SockShop node, via kubectl.

```
kubectl -s localhost:8082 apply -f - <<EOF
apiVersion: root.sockshop.com/v1
kind: SockShop
metadata:
  name: default
spec:
  orgName: Unicorn
  location: Seattle
  website: Unicorn.inc
EOF
```

#### REST API is available [here](http://localhost:8082/sockshop.com/docs#/)

![RESTAPI](../images/Playground-11-Nexus-API-1.png)

### Nexus Runtime based on KIND

#### Kubectl API Access

Lets instantiate sock shop by creating the SockShop node, via kubectl.

***NOTE: Replace \<PORT> with CLUSTER_PORT used when creating runtime.***

```
kubectl -s localhost:<PORT> apply -f - <<EOF
apiVersion: root.sockshop.com/v1
kind: SockShop
metadata:
  name: default
spec:
  orgName: Unicorn
  location: Seattle
  website: Unicorn.inc
EOF
```

#### REST API is available: `http://localhost:<PORT+1>/sockshop.com/docs#/`

