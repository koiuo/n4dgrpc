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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStrToPath_simple(t *testing.T) {
	path, err := StrToPath("/default")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(path.Elems))
	assert.Equal(t, "default", string(path.Elems[0]))
}

func TestStrToPath_multipleElements(t *testing.T) {
	path, err := StrToPath("/service/consul")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(path.Elems))
	assert.Equal(t, "service", string(path.Elems[0]))
	assert.Equal(t, "consul", string(path.Elems[1]))
}

func TestStrToPath_noslash_err(t *testing.T) {
	path, err := StrToPath("default")

	assert.Nil(t, path)
	assert.Error(t, err)
}

func TestStrToPath_root_err(t *testing.T) {
	path, err := StrToPath("/")

	assert.Nil(t, path)
	assert.Error(t, err)
}

func TestToIPv4TCP(t *testing.T) {
	ip, err := toIPv4TCP([]byte{8, 8, 8, 8}, 53)

	assert.NoError(t, err)
	assert.Equal(t, "8.8.8.8:53", ip)
}

func TestToIPv4TCP_tooShort(t *testing.T) {
	ip, err := toIPv4TCP([]byte{8, 8, 8}, 53)

	assert.Error(t, err)
	assert.Empty(t, ip)
}

func TestToIPv4TCP_tooLong(t *testing.T) {
	ip, err := toIPv4TCP([]byte{8, 8, 8, 8, 8}, 53)

	assert.Error(t, err)
	assert.Empty(t, ip)
}
