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

package client

import (
	"context"
	"fmt"
	n4d "github.com/edio/n4dgrpc/client/mock"
	"github.com/edio/n4dgrpc/convertions"
	"github.com/golang/mock/gomock"
	mesh "github.com/linkerd/linkerd/mesh/core/src/main/protobuf"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"testing"
)

// should return [path],nil if leaf is bound
func Test_bind_leaf(t *testing.T) {
	// setup
	reset()

	// given
	expectedPath := path("/$/inet/8.8.8.8/53")

	// and
	ctrl := gomock.NewController(t)
	interpreterMock := n4d.NewMockInterpreterClient(ctrl)
	interpreterMock.EXPECT().GetBoundTree(
		gomock.Any(),
		gomock.Any(),
	).Return(leaf(expectedPath), nil)
	interpreter = interpreterMock

	// when
	paths, err := bind(context.Background(), path("/dns/google"), path("/default"))

	// then
	assert.NoError(t, err)
	assert.Len(t, paths, 1, "Paths should be array of exactly one element")
	assert.Equal(t, expectedPath, paths[0], "Paths should be array of exactly one element")
}

// should return nil,err if unexpected error has happened
func Test_bind_error(t *testing.T) {
	// setup
	reset()

	// given
	ctrl := gomock.NewController(t)
	interpreterMock := n4d.NewMockInterpreterClient(ctrl)
	interpreterMock.EXPECT().GetBoundTree(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, fmt.Errorf("Unexpected error"))
	interpreter = interpreterMock

	// when
	paths, err := bind(context.Background(), path("/dns/google"), path("/default"))

	// then
	assert.Error(t, err)
	assert.Nil(t, paths)
}

// should return [],ErrNegBinding if Neg
func Test_bind_neg(t *testing.T) {
	// setup
	reset()

	// given
	ctrl := gomock.NewController(t)
	interpreterMock := n4d.NewMockInterpreterClient(ctrl)
	interpreterMock.EXPECT().GetBoundTree(
		gomock.Any(),
		gomock.Any(),
	).Return(neg(), nil)
	interpreter = interpreterMock

	// when
	paths, err := bind(context.Background(), path("/dns/google"), path("/nonexistentroot"))

	// then
	assert.Error(t, err)
	assert.IsType(t, (*ErrNegBinding)(nil), err)
	assert.NotNil(t, paths)
	assert.Len(t, paths, 0)
}

// should return nil,err if no stream is returned
func Test_resolve_noStreamError(t *testing.T) {
	// setup
	reset()

	// given
	ctrl := gomock.NewController(t)
	resolverMock := n4d.NewMockResolverClient(ctrl)
	resolverMock.EXPECT().StreamReplicas(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, fmt.Errorf("Some error"))
	resolver = resolverMock

	// when
	paths, err := resolve(context.Background(), path("/$/inet/google.com/53"))

	// then
	assert.Error(t, err)
	assert.Nil(t, paths)
}

// should return endpoint,nil for bound resolution
func Test_resolve_bound(t *testing.T) {
	// setup
	reset()

	// given
	expectedEndpoint := endpoint([]byte{8, 8, 8, 8}, 53)

	// and
	ctrl := gomock.NewController(t)
	resolverStreamMock := n4d.NewMockResolver_StreamReplicasClient(ctrl)
	resolverStreamMock.EXPECT().Recv().Return(replicas(expectedEndpoint), nil)
	resolverStreamMock.EXPECT().CloseSend()

	// and
	resolverMock := n4d.NewMockResolverClient(ctrl)
	resolverMock.EXPECT().StreamReplicas(
		gomock.Any(),
		gomock.Any(),
	).Return(resolverStreamMock, nil)
	resolver = resolverMock

	// when
	endpoints, err := resolve(context.Background(), path("/$/inet/google.dns/53"))

	// then
	assert.NoError(t, err)
	assert.NotNil(t, endpoints)
	assert.Len(t, endpoints, 1)
	assert.Equal(t, expectedEndpoint, endpoints[0])
}

// should return endpoint,nil for bound resolution if io.EOF happens but replica is returned
func Test_resolve_bound_and_EOF(t *testing.T) {
	// setup
	reset()

	// given
	expectedEndpoint := endpoint([]byte{8, 8, 8, 8}, 53)

	// and
	ctrl := gomock.NewController(t)
	resolverStreamMock := n4d.NewMockResolver_StreamReplicasClient(ctrl)
	resolverStreamMock.EXPECT().Recv().Return(replicas(expectedEndpoint), io.EOF)
	resolverStreamMock.EXPECT().CloseSend()

	// and
	resolverMock := n4d.NewMockResolverClient(ctrl)
	resolverMock.EXPECT().StreamReplicas(
		gomock.Any(),
		gomock.Any(),
	).Return(resolverStreamMock, nil)
	resolver = resolverMock

	// when
	endpoints, err := resolve(context.Background(), path("/$/inet/google.dns/53"))

	// then
	assert.NoError(t, err)
	assert.NotNil(t, endpoints)
	assert.Len(t, endpoints, 1)
	assert.Equal(t, expectedEndpoint, endpoints[0])
}

// should return nil,EOF if io.EOF happens and no replica is returned
func Test_resolve_nil_and_EOF(t *testing.T) {
	// setup
	reset()

	// given
	ctrl := gomock.NewController(t)
	resolverStreamMock := n4d.NewMockResolver_StreamReplicasClient(ctrl)
	resolverStreamMock.EXPECT().Recv().Return(nil, io.EOF)
	resolverStreamMock.EXPECT().CloseSend()

	// and
	resolverMock := n4d.NewMockResolverClient(ctrl)
	resolverMock.EXPECT().StreamReplicas(
		gomock.Any(),
		gomock.Any(),
	).Return(resolverStreamMock, nil)
	resolver = resolverMock

	// when
	endpoints, err := resolve(context.Background(), path("/$/inet/google.dns/53"))

	// then
	assert.Error(t, err)
	assert.Equal(t, err, io.EOF)
	assert.Nil(t, endpoints)
}

// Create BoundTreeRsp with Leaf node
func leaf(path *mesh.Path) *mesh.BoundTreeRsp {
	return &mesh.BoundTreeRsp{
		Tree: &mesh.BoundNameTree{
			Node: &mesh.BoundNameTree_Leaf_{
				Leaf: &mesh.BoundNameTree_Leaf{
					Id: path,
				},
			},
		},
	}
}

// Create BoundTreeRsp with Neg node
func neg() *mesh.BoundTreeRsp {
	return &mesh.BoundTreeRsp{
		Tree: &mesh.BoundNameTree{
			Node: &mesh.BoundNameTree_Neg_{
				Neg: &mesh.BoundNameTree_Neg{},
			},
		},
	}
}

func replicas(endpoints ...*mesh.Endpoint) *mesh.Replicas {
	return &mesh.Replicas{
		Result: &mesh.Replicas_Bound_{
			Bound: &mesh.Replicas_Bound{
				Endpoints: endpoints,
			},
		},
	}
}

// Create Endpoint with specified addresses
func endpoint(ip4 []byte, port int) *mesh.Endpoint {
	return &mesh.Endpoint{
		Port:    80,
		Address: ip4,
		InetAf:  mesh.Endpoint_INET4,
		Meta:    &mesh.Endpoint_Meta{NodeName: "whatever"},
	}
}

// Create Path from string without returning an error
func path(str string) *mesh.Path {
	path, err := convertions.StrToPath(str)
	if err != nil {
		log.Fatalln(err)
	}
	return path
}

func reset() {
	Close()
}
