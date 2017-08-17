// Copyright 2017 Verizon
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

package persistence

// KeyValueStore Interface defines how we interact with key value backends.
type KeyValueStore interface {
	Create(key, value string) error
	CreateWithLease(key, value string, ttl int64) (int64, error)
	Read(key string) (string, error)
	ReadAll(key string) (map[string]string, error)
	Update(key, value string) error
	RefreshLease(int64) error
	Delete(key string) error
}
