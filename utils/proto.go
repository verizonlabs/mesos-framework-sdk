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

package utils

func ProtoString(s string) *string {
	return &s
}

func ProtoFloat64(f float64) *float64 {
	return &f
}

func ProtoInt64(i int64) *int64 {
	return &i
}

func ProtoBool(i bool) *bool {
	return &i
}

func ProtoInt32(i int32) *int32 {
	return &i
}

func ProtoUint32(i uint32) *uint32 {
	return &i
}
