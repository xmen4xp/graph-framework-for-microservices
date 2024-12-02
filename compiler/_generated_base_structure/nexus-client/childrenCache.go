//nolint:stylecheck,revive // Inherited from opensource.s
package nexus_client

import (
	"sync"
)

// Package level cache to track parent to child relationships.
// key: node crd type string of the parent.
// value: *nodeCache
var parentChildCache = &sync.Map{}

type nodeCache struct {
	// Key: Node hashed name string. Answers the question if a node exists or it has children.
	// Value: *childNodeTypeCache
	node sync.Map
}

type childNodeTypeCache struct {
	// Key: node crd type string of the child.
	// Value: *childrenCache
	childType sync.Map
}

type childrenCache struct {
	// Key: Node hashed name string. Answers the question if a child with the name exists.
	// Value: struct{}. Indicates presence.
	children sync.Map
}

func AddChild(parentNodeType, parentName, childNodeType, childName string) {
	parentNodeTypeVal, parentNodeTypeFound := parentChildCache.Load(parentNodeType)
	var parentNodeCache *nodeCache
	var typeAsserted bool

	if !parentNodeTypeFound {
		parentNodeCache = &nodeCache{}
		parentChildCache.Store(parentNodeType, parentNodeCache)
	} else {
		parentNodeCache, typeAsserted = (parentNodeTypeVal).(*nodeCache)
		if !typeAsserted {
			logger.Fatalf("invalid parentNodeTypeVal found in nodeCache for parentNodeType %s", parentNodeType)
		}
	}

	childNodeTypeVal, parentNodeFound := parentNodeCache.node.Load(parentName)
	var childTypeCache *childNodeTypeCache
	if !parentNodeFound {
		childTypeCache = &childNodeTypeCache{}
		parentNodeCache.node.Store(parentName, childTypeCache)
	} else {
		childTypeCache, typeAsserted = childNodeTypeVal.(*childNodeTypeCache)
		if !typeAsserted {
			logger.Fatalf("invalid childNodeTypeVal found in childNodeTypeCache for parentName %s", parentName)
		}
	}

	childrenCacheVal, childNodeTypeFound := childTypeCache.childType.Load(childNodeType)
	var cache *childrenCache
	if !childNodeTypeFound {
		cache = &childrenCache{}
		childTypeCache.childType.Store(childNodeType, cache)
	} else {
		cache, typeAsserted = childrenCacheVal.(*childrenCache)
		if !typeAsserted {
			logger.Fatalf("invalid childrenCacheVal found in childrenCacheVal forchildNodeType %s", childNodeType)
		}
	}

	cache.children.Store(childName, struct{}{})
}

func RemoveChild(parentNodeType, parentName, childNodeType, childName string) {
	var typeAsserted bool

	// Remove entry on the last child being removed.
	nodeVal, parentNodeTypeFound := parentChildCache.Load(parentNodeType)
	if !parentNodeTypeFound {
		return
	}
	parentNodeCache, typeAsserted := (nodeVal).(*nodeCache)
	if !typeAsserted {
		logger.Fatalf("invalid nodeVal found in nodeCache for parentNodeType %s", parentNodeType)
	}

	childNodeTypeVal, parentNodeFound := parentNodeCache.node.Load(parentName)
	if !parentNodeFound {
		return
	}
	childNodeTypeCache, typeAsserted := childNodeTypeVal.(*childNodeTypeCache)
	if !typeAsserted {
		logger.Fatalf("invalid childNodeTypeVal found in childNodeTypeVal for parentName %s", parentName)
	}

	childrenCacheVal, childNodeTypeFound := childNodeTypeCache.childType.Load(childNodeType)
	if !childNodeTypeFound {
		return
	}
	childrenCache, typeAsserted := childrenCacheVal.(*childrenCache)
	if !typeAsserted {
		logger.Fatalf("invalid childrenCacheVal found in childrenCache for childNodeType %s", childNodeType)
	}

	childrenCache.children.Delete(childName)
}

func IsChildExists(parentNodeType, parentName, childNodeType, childName string) bool {
	var typeAsserted bool

	nodeVal, parentNodeTypeFound := parentChildCache.Load(parentNodeType)
	if !parentNodeTypeFound {
		return false
	}
	parentNodeCache, typeAsserted := (nodeVal).(*nodeCache)
	if !typeAsserted {
		logger.Fatalf("invalid nodeVal found in nodeCache for parentNodeType %s", parentNodeType)
	}

	childNodeTypeVal, parentNodeFound := parentNodeCache.node.Load(parentName)
	if !parentNodeFound {
		return false
	}
	childNodeTypeCache, typeAsserted := childNodeTypeVal.(*childNodeTypeCache)
	if !typeAsserted {
		logger.Fatalf("invalid childNodeTypeVal found in childNodeTypeCache for parentName %s", parentName)
	}

	childrenCacheVal, childNodeTypeFound := childNodeTypeCache.childType.Load(childNodeType)
	if !childNodeTypeFound {
		return false
	}
	childrenCache, typeAsserted := childrenCacheVal.(*childrenCache)
	if !typeAsserted {
		logger.Fatalf("invalid childrenCacheVal found in childrenCache for childNodeType %s", childNodeType)
	}

	if _, ok := childrenCache.children.Load(childName); ok {
		return true
	}
	return false
}

func GetChildren(parentNodeType, parentName, childNodeType string) (children []string) {
	var typeAsserted bool

	nodeVal, parentNodeTypeFound := parentChildCache.Load(parentNodeType)
	if !parentNodeTypeFound {
		return children
	}
	parentNodeCache, typeAsserted := (nodeVal).(*nodeCache)
	if !typeAsserted {
		logger.Fatalf("invalid nodeVal found in nodeCache for parentNodeType %s", parentNodeType)
	}

	childNodeTypeVal, parentNodeFound := parentNodeCache.node.Load(parentName)
	if !parentNodeFound {
		return children
	}
	childNodeTypeCache, typeAsserted := childNodeTypeVal.(*childNodeTypeCache)
	if !typeAsserted {
		logger.Fatalf("invalid childNodeTypeVal found in childNodeTypeCache for parentName %s", parentName)
	}

	childrenCacheVal, childNodeTypeFound := childNodeTypeCache.childType.Load(childNodeType)
	if !childNodeTypeFound {
		return children
	}
	childrenCache, typeAsserted := childrenCacheVal.(*childrenCache)
	if !typeAsserted {
		logger.Fatalf("invalid childrenCacheVal found in childrenCache for childNodeType %s", childNodeType)
	}

	childrenCache.children.Range(func(k, v interface{}) bool {
		children = append(children, k.(string))
		return true
	})
	return children
}
