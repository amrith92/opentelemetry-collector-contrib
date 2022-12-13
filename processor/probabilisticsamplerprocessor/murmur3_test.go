// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package probabilisticsamplerprocessor

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/idutils"
)

// TestHashHasNoCollisions ensures that the hash function supports different key lengths even if in
// practice it is only expected to receive keys with length 16 (trace id length in OC proto).
func TestHashHasNoCollisions(t *testing.T) {
	// Statistically a random selection of such small number of keys should not result in
	// collisions, but, of course it is possible that they happen, a different random source
	// should avoid that.
	r := rand.New(rand.NewSource(1))
	fullKey := idutils.UInt64ToTraceID(r.Uint64(), r.Uint64())
	seen := make(map[uint32]bool)
	for i := 1; i <= len(fullKey); i++ {
		key := fullKey[:i]
		hash := hash(key, 1)
		require.False(t, seen[hash], "Unexpected duplicated hash")
		seen[hash] = true
	}
}