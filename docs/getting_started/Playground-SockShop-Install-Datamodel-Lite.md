# Install Sock Shop data model on Nexus Runtime

[[Prev]](Playground-SockShop-Compile-Datamodel-Lite.md) [[Exit]](../../README.md) [[Next]](Playground-SockShop-Access-Datamodel-API-Lite.md)


## Export KUBECONFIG to Nexus Runtime

The KUBECONFIG to export depends runtime being used in this playgroud.

***Option 1***: If you running a K0s based Nexus runtime:

Run this make target to get the shell export command to execute:
```
make -C $NEXUS_REPO_DIR runtime.k0s.kubeconfig.export
```

***Option 2***: If you running a KIND based Nexus runtime:

Run this make target to get the shell export command to execute:
```
CLUSTER_NAME=<name> make -C $NEXUS_REPO_DIR runtime.k0s.kubeconfig.export
```

***NOTE: Remember to execute the printed "export" command on your shell.***

## Install data model
```
make dm.install
```

[[Prev]](Playground-SockShop-Compile-Datamodel-Lite.md) [[Exit]](../../README.md) [[Next]](Playground-SockShop-Access-Datamodel-API-Lite.md)
