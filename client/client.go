// Copyright © 2017 Dmytro Kostiuchenko edio@archlinux.us
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Stateful client
// Call Connect() before calling anything else
// You may want also to set custom to exported variables
package client

import (
	"context"
	mesh "github.com/linkerd/linkerd/mesh/core/src/main/protobuf"
	"google.golang.org/grpc"
	"time"
	"fmt"
	"log"
)

// exported
var (
	// Timeout for initial dial
	DialTimeout time.Duration = 1 * time.Second
	// Timeout for completing operation (shared if multiple calls to n4d are involved)
	OpTimeout time.Duration = 1 * time.Second
)

// private
var (
	connection *grpc.ClientConn
	ctx = context.Background()
)

// Connect to n4d
func Connect(namerdAddress string) (err error) {
	lctx, cancel := context.WithTimeout(ctx, DialTimeout)
	defer cancel()
	connection, err = grpc.DialContext(lctx, namerdAddress, grpc.WithInsecure())
	return
}

// Bind a name in a namespace specified by root
func Bind(root *mesh.Path, name *mesh.Path) ([]*mesh.Path, error) {
	lctx, cancel := context.WithTimeout(ctx, OpTimeout)
	defer cancel()
	return bind(lctx, root, name)
}

func bind(ctx context.Context, root *mesh.Path, name *mesh.Path) ([]*mesh.Path, error) {
	bindReq := &mesh.BindReq{
		Root: root,
		Name: name,
	}

	interpreter := mesh.NewInterpreterClient(connection)
	resp, err := interpreter.GetBoundTree(ctx, bindReq)
	if err != nil {
		return nil, err
	}

	switch resp.Tree.Node.(type) {
	case *mesh.BoundNameTree_Leaf_:
		return []*mesh.Path{resp.Tree.GetLeaf().Id}, nil
	case *mesh.BoundNameTree_Empty_:
	case *mesh.BoundNameTree_Neg_:
		return []*mesh.Path{}, nil
	default:
		return nil, fmt.Errorf("Not supported yet: %v", resp)
	}
	return nil, fmt.Errorf("Something unexpected has happened")
}

// Resolve a name in a namespace specified by root
func Resolve(root *mesh.Path, name *mesh.Path) ([]*mesh.Endpoint, error) {
	lctx, cancel := context.WithTimeout(ctx, OpTimeout)
	defer cancel()

	boundPaths, err := bind(lctx, root, name)
	if (err != nil || len(boundPaths) == 0) {
		return nil, fmt.Errorf("Not bound: %v", err)
	}

	var endpoints []*mesh.Endpoint
	for _, path := range boundPaths {
		// TODO return typed errors to distinguish downstream
		e, err := resolve(lctx, path)
		if err != nil {
			log.Printf("Error resolving [%v]: %v", path, err)
		}
		endpoints = append(endpoints, e...)
	}

	return endpoints, nil
}

func resolve(ctx context.Context, boundPath *mesh.Path) ([]*mesh.Endpoint, error) {
	replicasReq := &mesh.ReplicasReq{
		Id: boundPath,
	}

	resolver := mesh.NewResolverClient(connection)
	resp, err := resolver.GetReplicas(ctx, replicasReq)
	if err != nil {
		return nil, err
	}

	// TODO handle pending
	endpoints := []*mesh.Endpoint{}
	if resp.GetBound() != nil{
		endpoints = resp.GetBound().Endpoints
	}
	return endpoints, nil
}

