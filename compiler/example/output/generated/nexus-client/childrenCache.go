package nexus_client

import (
	"sync"
)

// Package level cache to track parent to child relationships.
// key: node crd type string of the parent.
// value: *nodeCache
var parentChildCache *sync.Map = &sync.Map{}

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
	if parentNodeTypeFound == false {
		parentNodeCache = &nodeCache{}
		parentChildCache.Store(parentNodeType, parentNodeCache)
	} else {
		parentNodeCache = (parentNodeTypeVal).(*nodeCache)
	}

	childNodeTypeVal, parentNodeFound := parentNodeCache.node.Load(parentName)
	var childTypeCache *childNodeTypeCache
	if parentNodeFound == false {
		childTypeCache = &childNodeTypeCache{}
		parentNodeCache.node.Store(parentName, childTypeCache)
	} else {
		childTypeCache = childNodeTypeVal.(*childNodeTypeCache)
	}

	childrenCacheVal, childNodeTypeFound := childTypeCache.childType.Load(childNodeType)
	var cache *childrenCache
	if childNodeTypeFound == false {
		cache = &childrenCache{}
		childTypeCache.childType.Store(childNodeType, cache)
	} else {
		cache = childrenCacheVal.(*childrenCache)
	}

	cache.children.Store(childName, struct{}{})
}

func RemoveChild(parentNodeType, parentName, childNodeType, childName string) {
	// Remove entry on the last child being removed.
	nodeVal, parentNodeTypeFound := parentChildCache.Load(parentNodeType)
	if parentNodeTypeFound == false {
		return
	}
	parentNodeCache := (nodeVal).(*nodeCache)

	childNodeTypeVal, parentNodeFound := parentNodeCache.node.Load(parentName)
	if parentNodeFound == false {
		return
	}
	childNodeTypeCache := childNodeTypeVal.(*childNodeTypeCache)

	childrenCacheVal, childNodeTypeFound := childNodeTypeCache.childType.Load(childNodeType)
	if childNodeTypeFound == false {
		return
	}
	childrenCache := childrenCacheVal.(*childrenCache)

	childrenCache.children.Delete(childName)
	/*
	   	if len(childrenCache.children) > 0 {
	   		return
	   	}

	   delete(childNodeTypeCache.childType, childNodeType)

	   	if len(childNodeTypeCache.childType) > 0 {
	   		return
	   	}

	   delete(parentNodeCache.node, parentName)

	   	if len(parentNodeCache.node) > 0 {
	   		return
	   	}

	   delete(parentChildCache, parentNodeType)
	*/
}

func DeleteParent(parentNodeType, parentName string) {
	// recursively walk and delete all children
}

func IsChildExists(parentNodeType, parentName, childNodeType, childName string) bool {
	nodeVal, parentNodeTypeFound := parentChildCache.Load(parentNodeType)
	if parentNodeTypeFound == false {
		return false
	}
	parentNodeCache := (nodeVal).(*nodeCache)

	childNodeTypeVal, parentNodeFound := parentNodeCache.node.Load(parentName)
	if parentNodeFound == false {
		return false
	}
	childNodeTypeCache := childNodeTypeVal.(*childNodeTypeCache)

	childrenCacheVal, childNodeTypeFound := childNodeTypeCache.childType.Load(childNodeType)
	if childNodeTypeFound == false {
		return false
	}
	childrenCache := childrenCacheVal.(*childrenCache)

	if _, ok := childrenCache.children.Load(childName); ok {
		return true
	}
	return false
}

func GetChildren(parentNodeType, parentName, childNodeType string) (children []string) {
	nodeVal, parentNodeTypeFound := parentChildCache.Load(parentNodeType)
	if parentNodeTypeFound == false {
		return
	}
	parentNodeCache := (nodeVal).(*nodeCache)

	childNodeTypeVal, parentNodeFound := parentNodeCache.node.Load(parentName)
	if parentNodeFound == false {
		return
	}
	childNodeTypeCache := childNodeTypeVal.(*childNodeTypeCache)

	childrenCacheVal, childNodeTypeFound := childNodeTypeCache.childType.Load(childNodeType)
	if childNodeTypeFound == false {
		return
	}
	childrenCache := childrenCacheVal.(*childrenCache)

	childrenCache.children.Range(func(k, v interface{}) bool {
		children = append(children, k.(string))
		return true
	})
	return
}
