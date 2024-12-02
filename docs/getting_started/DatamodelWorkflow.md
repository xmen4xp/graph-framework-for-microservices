# Getting Started

With Nexus you can create your own extensible, distributed platform that:

* Implements a datamodel as K8s CRDs.

* Generates accessors that facilitate easy consumption of the datamodel.

* Bootstrap a cloud native application that consumes that datamodel and is ready-to-go in no time.


## Datamodel Workflow

Nexus Datamodel provides the framework, language, toolkit and runtime to implement state/data required by application, that is:

* hierarchical
* distributed
* consistent
* customizable

Nexus datamodel framework supports formalization of objects/nodes and their relationship in a spec that is expressed using Golang syntax and is fully compliant with the Golang compiler.

**NOTE: If you would prefer to follow a pre-cooked Datamodel, as you try this workflow, execute commands from the code sections below** 

This guided workflow will walk you through setting up a datamodel that is local to your application or Importable datamodel

* [Nexus Install](DatamodelWorkflow.md#nexus-install)

* [Nexus Pre-req verify](DatamodelWorkflow.md#nexus-pre-req-verify)

* [Setup Datamodel Workspace](DatamodelWorkflow.md#setup-datamodel-workspace)

* [Datamodel Build](DatamodelWorkflow.md#datamodel-build)
  
* [Datamodel Install](DatamodelWorkflow.md#datamodel-install)
  
* [Datamodel Playground](DatamodelWorkflow.md#datamodel-playground)


## Nexus Install

1. Install Nexus CLI
    ```
    GOPRIVATE="gitlab.eng.vmware.com" go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/cmd/plugin/nexus@master
    ```

<details><summary>FAQ</summary>
1. The above command shows unable to connect to gitlab.eng.vmware.com ?

Verify that you have permissions to the repo
            
    git ls-remote git@gitlab.eng.vmware.com:nsx-allspark_users/nexus-sdk/cli.git

Update gitconfig to use ssh instead of https

    git config --global url.git@gitlab.eng.vmware.com:.insteadOf https://gitlab.eng.vmware.com

</details>

2. Verify nexus sdk pre-requisites are satisfied

   ```
   nexus prereq verify
   ```

## Setup Datamodel Workspace

1. Create and `cd` to your workspace directory to create, compile and install datamodel
    ```
    mkdir -p $HOME/test-datamodel/orgchart && cd $HOME/test-datamodel/orgchart       
    ```
   
     
1. Initialize datamodel workspace
    ```
    nexus datamodel init --name orgchart --group orgchart.org
    ```

1. Start writing datamodel specification for your application.
   
   **To understand the workflow we can use the below example datamodel. To write your own datamodel please refer** [here](../../compiler/DSL.md)
   
**Example Orgchart DSL**

   The Orgchart Application has 3 levels. 1. Leader, 2. Manager and 3. Engineer

```shell
echo 'package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

var LeaderRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/leader/{root.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/leaders",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:LeaderRestAPISpec
type Leader struct {

	// Tags "Root" as a node in datamodel graph
	nexus.Node

	Name          string
	Designation   string
	DirectReports Manager `nexus:"children"`
}

type Manager struct {

	// Tags "Root" as a node in datamodel graph
	nexus.Node

	Name          string
	Designation   string
	DirectReports Engineer `nexus:"children"`
}

type Engineer struct {

	// Tags "Root" as a node in datamodel graph
	nexus.Node

	Name        string
	Designation string
}
' > $HOME/test-datamodel/orgchart/root.go
```

## Datamodel Build

   ```
   nexus datamodel build --name orgchart
   ```

This generates libraries, types, runtime and metadata required to implement the datamodel at runtime.


## Datamodel Install

### Pre-requisites
To install datamodel we need to install nexus runtime as a pre-requisite
  
#### [Install Runtime](RuntimeWorkflow.md)


   ```
   DOCKER_REPO=orgchart VERSION=latest make docker_build
   ```

   ```
   kind load docker-image orgchart:latest --name <kind cluster name>
   ```

   ```
   nexus datamodel install image orgchart:latest --namespace <name>
   ```

## Datamodel Playground

#### [App Datamodel Usage Workflow](AppDatamodelUsageWorkflow.md)
