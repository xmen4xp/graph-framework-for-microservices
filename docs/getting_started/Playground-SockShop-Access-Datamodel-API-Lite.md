# Access SockShop data model through API


[[Prev]](Playground-SockShop-Install-Datamodel-Lite.md) [[Exit]](../../README.md) [[Next]](Playground-SockShop-Wrap-Lite.md)

***Installation Complete !***

## Access your API's

[Nexus Runtime based on K0s](#Nexus-Runtime-based-on-K0s)

[Nexus Runtime based on KIND](#Nexus-Runtime-based-on-KIND)

## Nexus Runtime based on K0s

### Kubectl API Access

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

### REST API is available [here](http://localhost:8082/sockshop.com/docs#/)

![RESTAPI](../images/Playground-11-Nexus-API-1.png)

## Nexus Runtime based on KIND

### Kubectl API Access

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

### REST API is available: `http://localhost:<PORT>/sockshop.com/docs#/`

![RESTAPI](../images/Playground-11-Nexus-API-1.png)

[[Prev]](Playground-SockShop-Install-Datamodel-Lite.md) [[Exit]](../../README.md) [[Next]](Playground-SockShop-Wrap-Lite.md)
