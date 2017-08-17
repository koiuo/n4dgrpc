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

package convertions

import (
	"errors"
	mesh "github.com/linkerd/linkerd/mesh/core/src/main/protobuf"
	"strings"
	"fmt"
	//"encoding/hex"
	"encoding/hex"
)

// Convert string to mesh.Path
func StrToPath(path string) (meshPath *mesh.Path, err error) {
	// perform only basic validations that impact parsing, server will validate rest
	if len(path) <= 1 {
		err = errors.New("Path must start with /")
		return
	}
	if path[0] != '/' {
		err = errors.New("Path must start with /")
		return
	}

	var elements [][]byte
	for _, el := range strings.Split(path[1:], "/") {
		elements = append(elements, []byte(el))
	}

	meshPath = &mesh.Path{
		Elems: elements,
	}
	return
}

func PathToStr(path *mesh.Path) (str string) {
	var strs []string = []string{""} // init with empty string to get "/" in beginning on Join

	for _, el := range path.Elems {
		strs = append(strs, string(el))
	}
	str = strings.Join(strs, "/")
	return
}

func EndpointToStr(endpoint *mesh.Endpoint) (string, error) {
	if endpoint.InetAf == mesh.Endpoint_INET6 {
		return toIPv6TCP(endpoint.Address, int(endpoint.Port))
	} else {
		return toIPv4TCP(endpoint.Address, int(endpoint.Port))
	}
}

func toIPv4TCP(ip []byte, port int) (string, error) {
	if len(ip) != 4 {
		return "", fmt.Errorf("Expected 4 bytes, got %d in %v", len(ip), ip)
	}
	return fmt.Sprintf("%d.%d.%d.%d:%d", ip[0], ip[1], ip[2], ip[3], port), nil
}

func toIPv6TCP(ip []byte, port int) (string, error) {
	if len(ip) != 16 {
		return "", fmt.Errorf("Expected 16 bytes, got %d in %v", len(ip), ip)
	}

	segments := make([]string, 8)

	for i := 0; i < 8; i++ {
		segments[i] = hex.EncodeToString(ip[i*2:i*2+2])
	}

	return fmt.Sprintf("[%s]:%d", strings.Join(segments, ":"), port), nil
}