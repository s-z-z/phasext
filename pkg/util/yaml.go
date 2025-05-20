package util

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// UpdateYAML updates the input YAML with values from the provided map
func UpdateYAML(input []byte, values map[string]any) ([]byte, error) {
	// Parse input YAML
	var node yaml.Node
	if err := yaml.Unmarshal(input, &node); err != nil {
		return nil, fmt.Errorf("failed to parse input YAML: %w", err)
	}

	// Ensure we're working with a document node
	if node.Kind != yaml.DocumentNode || len(node.Content) == 0 {
		return nil, fmt.Errorf("invalid YAML document")
	}

	// Process each value in the map
	for key, value := range values {
		// Split the key into parts for nested access
		parts := strings.Split(key, ".")
		if err := updateNode(node.Content[0], parts, value); err != nil {
			return nil, fmt.Errorf("failed to update key %s: %w", key, err)
		}
	}

	// Marshal back to YAML
	output, err := yaml.Marshal(&node)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return output, nil
}

// updateNode recursively updates the YAML node based on the key path
func updateNode(node *yaml.Node, parts []string, value any) error {
	if len(parts) == 0 {
		return fmt.Errorf("empty key path")
	}

	// If we're at the last part of the path, set the value
	if len(parts) == 1 {
		return setNodeValue(node, parts[0], value)
	}

	// For mapping nodes
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			if keyNode.Value == parts[0] {
				return updateNode(node.Content[i+1], parts[1:], value)
			}
		}
		// If key doesn't exist, create new mapping
		newNode := &yaml.Node{Kind: yaml.MappingNode}
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: parts[0]},
			newNode)
		return updateNode(newNode, parts[1:], value)
	}

	return fmt.Errorf("cannot traverse non-mapping node at %s", parts[0])
}

// setNodeValue sets the value for a specific key in the node
func setNodeValue(node *yaml.Node, key string, value any) error {
	// For mapping nodes
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			if keyNode.Value == key {
				return setSingleNodeValue(node.Content[i+1], value)
			}
		}
		// Key doesn't exist, create new entry
		newValueNode := &yaml.Node{}
		if err := setSingleNodeValue(newValueNode, value); err != nil {
			return err
		}
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: key},
			newValueNode)
		return nil
	}

	return fmt.Errorf("cannot set value in non-mapping node")
}

// setSingleNodeValue sets the value for a single node
func setSingleNodeValue(node *yaml.Node, value any) error {
	switch v := value.(type) {
	case string:
		node.Kind = yaml.ScalarNode
		node.Tag = "!!str"
		node.Value = v
	case bool:
		node.Kind = yaml.ScalarNode
		node.Tag = "!!bool"
		node.Value = fmt.Sprintf("%t", v)
	case []string:
		node.Kind = yaml.SequenceNode
		node.Tag = "!!seq"
		node.Content = nil
		for _, s := range v {
			node.Content = append(node.Content,
				&yaml.Node{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: s,
				})
		}
	case []any:
		node.Kind = yaml.SequenceNode
		node.Tag = "!!seq"
		node.Content = nil
		for _, item := range v {
			newNode := &yaml.Node{}
			if err := setSingleNodeValue(newNode, item); err != nil {
				return err
			}
			node.Content = append(node.Content, newNode)
		}
	default:
		return fmt.Errorf("unsupported value type: %T", value)
	}
	return nil
}
