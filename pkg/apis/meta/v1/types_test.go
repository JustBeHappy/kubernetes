/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ugorji/go/codec"
)

func TestVerbsMarshalJSON(t *testing.T) {
	cases := []struct {
		input  APIResource
		result string
	}{
		{APIResource{}, `{"name":"","namespaced":false,"kind":"","verbs":null}`},
		{APIResource{Verbs: Verbs([]string{})}, `{"name":"","namespaced":false,"kind":"","verbs":[]}`},
		{APIResource{Verbs: Verbs([]string{"delete"})}, `{"name":"","namespaced":false,"kind":"","verbs":["delete"]}`},
	}

	for i, c := range cases {
		result, err := json.Marshal(&c.input)
		if err != nil {
			t.Errorf("[%d] Failed to marshal input: '%v': %v", i, c.input, err)
		}
		if string(result) != c.result {
			t.Errorf("[%d] Failed to marshal input: '%v': expected %+v, got %q", i, c.input, c.result, string(result))
		}
	}
}

func TestVerbsUnmarshalJSON(t *testing.T) {
	cases := []struct {
		input  string
		result APIResource
	}{
		{`{}`, APIResource{}},
		{`{"verbs":null}`, APIResource{}},
		{`{"verbs":[]}`, APIResource{Verbs: Verbs([]string{})}},
		{`{"verbs":["delete"]}`, APIResource{Verbs: Verbs([]string{"delete"})}},
	}

	for i, c := range cases {
		var result APIResource
		if err := codec.NewDecoderBytes([]byte(c.input), new(codec.JsonHandle)).Decode(&result); err != nil {
			t.Errorf("[%d] Failed to unmarshal input '%v': %v", i, c.input, err)
		}
		if !reflect.DeepEqual(result, c.result) {
			t.Errorf("[%d] Failed to unmarshal input '%v': expected %+v, got %+v", i, c.input, c.result, result)
		}
	}
}

func TestVerbsUgorjiUnmarshalJSON(t *testing.T) {
	cases := []struct {
		input  string
		result APIResource
	}{
		{`{}`, APIResource{}},
		{`{"verbs":null}`, APIResource{}},
		{`{"verbs":[]}`, APIResource{Verbs: Verbs([]string{})}},
		{`{"verbs":["delete"]}`, APIResource{Verbs: Verbs([]string{"delete"})}},
	}

	for i, c := range cases {
		var result APIResource
		if err := json.Unmarshal([]byte(c.input), &result); err != nil {
			t.Errorf("[%d] Failed to unmarshal input '%v': %v", i, c.input, err)
		}
		if !reflect.DeepEqual(result, c.result) {
			t.Errorf("[%d] Failed to unmarshal input '%v': expected %+v, got %+v", i, c.input, c.result, result)
		}
	}
}

func TestVerbsProto(t *testing.T) {
	cases := []APIResource{
		{},
		{Verbs: Verbs([]string{})},
		{Verbs: Verbs([]string{"delete"})},
	}

	for _, input := range cases {
		data, err := input.Marshal()
		if err != nil {
			t.Fatalf("Failed to marshal input: '%v': %v", input, err)
		}
		resource := APIResource{}
		if err := resource.Unmarshal(data); err != nil {
			t.Fatalf("Failed to unmarshal output: '%v': %v", input, err)
		}
		if !reflect.DeepEqual(input, resource) {
			t.Errorf("Marshal->Unmarshal is not idempotent: '%v' vs '%v'", input, resource)
		}
	}
}