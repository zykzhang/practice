// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package tree 提供了以树形结构保存路由项的相关操作。
package tree

import (
	"fmt"
	"net/http"

	"github.com/issue9/mux/internal/handlers"
	"github.com/issue9/mux/params"
)

// Tree 以树节点的形式保存的路由。
//
// 多段路由项，会提取其中的相同的内容组成树状结构的节点。
// 比如以下路由项：
//  /posts/{id}/author
//  /posts/{id}/author/emails
//  /posts/{id}/author/profile
//  /posts/1/author
// 会被转换成以下结构
//  /posts
//     |
//     +---- 1/author
//     |
//     +---- {id}/author
//               |
//               +---- /profile
//               |
//               +---- /emails
type Tree struct {
	node
	disableOptions bool
}

// New 声明一个 Tree 实例
func New(disableOptions bool) *Tree {
	return &Tree{
		node:           node{},
		disableOptions: disableOptions,
	}
}

// Add 添加路由项。
//
// methods 可以为空，表示添加除 OPTIONS 之外所有支持的请求方法。
func (tree *Tree) Add(pattern string, h http.Handler, methods ...string) error {
	n, err := tree.getNode(pattern)
	if err != nil {
		return err
	}

	if n.handlers == nil {
		n.handlers = handlers.New(tree.disableOptions)
	}

	return n.handlers.Add(h, methods...)
}

// Clean 清除路由项
func (tree *Tree) Clean(prefix string) {
	tree.clean(prefix)
}

// Remove 移除路由项
//
// methods 可以为空，表示删除所有内容。
func (tree *Tree) Remove(pattern string, methods ...string) error {
	child := tree.find(pattern)
	if child == nil {
		return fmt.Errorf("不存在的节点 %v", pattern)
	}

	if child.handlers == nil {
		if len(child.children) == 0 {
			child.parent.children = removeNodes(child.parent.children, child.pattern)
			child.parent.buildIndexes()
		}
		return nil
	}

	if child.handlers.Remove(methods...) && len(child.children) == 0 {
		child.parent.children = removeNodes(child.parent.children, child.pattern)
		child.parent.buildIndexes()
	}
	return nil
}

// 获取指定的节点，若节点不存在，则在该位置生成一个新节点。
func (tree *Tree) getNode(pattern string) (*node, error) {
	ss, err := split(pattern)
	if err != nil {
		return nil, err
	}

	return tree.node.getNode(ss)
}

// SetAllow 设置指定节点的 allow 报头。
// 若节点不存在，则会自动生成该节点。
func (tree *Tree) SetAllow(pattern, allow string) error {
	n, err := tree.getNode(pattern)
	if err != nil {
		return err
	}

	if n.handlers == nil {
		n.handlers = handlers.New(tree.disableOptions)
	}

	n.handlers.SetAllow(allow)
	return nil
}

// URL 根据参数生成地址。
//
// 若节点不存在，则会自动生成。
func (tree *Tree) URL(pattern string, params map[string]string) (string, error) {
	node, err := tree.getNode(pattern)
	if err != nil {
		return "", err
	}
	return node.url(params)
}

// Handler 找到与当前内容匹配的 handlers.Handlers 实例。
func (tree *Tree) Handler(path string) (*handlers.Handlers, params.Params) {
	ps := make(params.Params, 5)
	node := tree.match(path, ps)

	if node == nil || node.handlers == nil || node.handlers.Len() == 0 {
		return nil, nil
	}

	return node.handlers, ps
}
