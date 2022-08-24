package controllers

var datamodelCrdExample = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: datamodels.nexus.org
  creationTimestamp: null
spec:
  conversion:
    strategy: None
  group: nexus.org
  names:
    kind: Datamodel
    listKind: DatamodelList
    plural: datamodels
    shortNames:
    - datamodel
    singular: datamodel
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            properties:
              name:
                type: string
              title:
                default: Nexus API GW APIs
                type: string
              url:
                type: string
            type: object
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: null
  storedVersions:
  - v1
`
