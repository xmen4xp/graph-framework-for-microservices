# Compile the data model

[[Prev]](Playground-SockShop-Complete-Datamodel-Lite.md) [[Exit]](../../README.md) [[Next]](Playground-SockShop-Install-Datamodel-Lite.md)

![SockShop](../images/Playground-8-Compile-Datamodel.png)

Let's build our data model and generate all the artifacts needed for the runtime.

## Compile data model

Nexus compiler can be invoked to build the datamodel with the following command:

```
DATAMODEL_DOCKER_REGISTRY=<container-registry-for-datamodel> TAG=<datamodel-tag> make docker.build
```

Example

```
DATAMODEL_DOCKER_REGISTRY=foo.com TAG=test make docker.build
```

This will generate a docker image packaged with all required artifacts.

[[Prev]](Playground-SockShop-Complete-Datamodel-Lite.md) [[Exit]](../../README.md) [[Next]](Playground-SockShop-Install-Datamodel-Lite.md)