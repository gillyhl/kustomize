// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package filtersutil

import (
	"encoding/json"

	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// ApplyToJSON applies the filter to the json objects.
func ApplyToJSON(filter kio.Filter, objs ...marshalerUnmarshaler) error {
	var nodes []*yaml.RNode
	for i := range objs {
		node, err := getRNode(objs[i])
		if err != nil {
			return err
		}
		nodes = append(nodes, node)
		l, err := filter.Filter([]*yaml.RNode{node})
		if err != nil {
			return err
		}
		err = setRNode(objs[i], l[0])
		if err != nil {
			return err
		}
	}

	_, err := filter.Filter(nodes)
	if err != nil {
		return err
	}

	for i := range nodes {
		err = setRNode(objs[i], nodes[i])
		if err != nil {
			return err
		}
	}

	return nil
}

type marshalerUnmarshaler interface {
	json.Unmarshaler
	json.Marshaler
}

// getRNode converts k into an RNode
func getRNode(k json.Marshaler) (*yaml.RNode, error) {
	j, err := k.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return yaml.Parse(string(j))
}

// setRNode marshals node into k
func setRNode(k json.Unmarshaler, node *yaml.RNode) error {
	s, err := node.String()
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(s), &m); err != nil {
		return err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return k.UnmarshalJSON(b)
}
