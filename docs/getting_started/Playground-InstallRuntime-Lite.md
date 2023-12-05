# Install Sock Shop data model on Nexus Runtime
[[Prev]](Playground-InstallCLI-Lite.md) [[Exit]](../../README.md)  [[Next]](Playground-SockShop-Lite.md)

![NexusRuntime0](../images/Playground-9-Nexus-Runtime-0.png)

![NexusRuntime1](../images/Playground-9-Nexus-Runtime-1.png)

![NexusRuntime3](../images/Playground-9-Nexus-Runtime-3.png)

![NexusRuntime2](../images/Playground-10-Nexus-Runtime-2.png)


Nexus data model is installed on a Nexus Runtime software stack.

### 1. Build Runtime Artifacts
#### Specify a tag to use for locally built runtime artifacts
```
# The tag can be anything. Here we use a tag called "letsplay"
echo letsplay > TAG
```
#### Start Build
```
sudo make runtime.clean; make runtime.build
```

### 2. Install Runtime

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


[[Prev]](Playground-InstallCLI-Lite.md) [[Exit]](../../README.md)  [[Next]](Playground-SockShop-Lite.md)
