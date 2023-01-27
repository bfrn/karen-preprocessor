package preprocessor

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	Node_type_module             = "Module"
	Node_type_resource           = "Resource"
	Node_type_reference_resource = "ReferenceResource"
	Node_type_provider           = "Provider"
)

type Node interface {
	GetNodeType() string
	GetAddress() string
	SetLocation(location string)
	AddChild(address string)
	AddAttribute(key string, attribute interface{})
}

type node struct {
	Address    string                 `json:"address"`
	NodeType   string                 `json:"nodeType"`
	Location   string                 `json:"location"`
	Children   []string               `json:"children,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

func newNodeData(address string, nodeType string, children []string) *node {
	nodeData := new(node)

	nodeData.Address = address
	nodeData.NodeType = nodeType
	nodeData.Children = children
	nodeData.Attributes = make(map[string]interface{})
	return nodeData
}

func (nodeData *node) GetNodeType() string {
	return nodeData.NodeType
}

func (nodeData *node) GetAddress() string {
	return nodeData.Address
}

func (nodeData *node) SetLocation(location string) {
	nodeData.Location = location
}

func (nodeData *node) AddChild(address string) {
	nodeData.Children = append(nodeData.Children, address)
}

func (nodeData *node) AddAttribute(key string, attribute interface{}) {
	nodeData.Attributes[key] = attribute
}

func (nodeData *node) MarshalBinary() ([]byte, error) {
	return json.Marshal(nodeData)
}

// Module represents a terraform module
type Module struct {
	*node
}

func NewModule(
	address string,
	children []string,
) (*Module, error) {
	module := new(Module)
	module.node = newNodeData(address, Node_type_module, children)

	return module, nil
}

const (
	Resource_action_no_op   = "noop"
	Resource_action_created = "created"
)

const (
	State_current = "Current_State"
	State_planned = "Planned_State"
)

func (module *Module) MarshalBinary() ([]byte, error) {
	return json.Marshal(module)
}

//Resource represents a terraform resource.
type Resource struct {
	*node
	// Dependencies contain the addresses of the ressources on which this ressource depends
	Dependencies []string `json:"dependencies,omitempty"`
	// The  actions which are performed on the ressource when the plan file is executed.
	Actions []string `json:"actions,omitempty"`
	// States contain the attributes of a resource that are associated with a specific state
	States map[string]map[string]interface{} `json:"states,omitempty"`
}

func NewResource(
	address string,
	dependencies []string,
) (*Resource, error) {
	resource := new(Resource)
	resource.node = newNodeData(address, Node_type_resource, nil)

	resource.Dependencies = dependencies
	resource.States = make(map[string]map[string]interface{})
	return resource, nil
}

func (resource *Resource) addState(state string, attributes map[string]interface{}) {
	resource.States[state] = attributes
}

func (resource *Resource) removeSensitiveValues(stateName string, sensitiveValues map[string]interface{}) error {
	state := resource.States[stateName]

	var remove func(interface{}, []string) error
	remove = func(state interface{}, key []string) error {
		currentKey := key[0]
		if len(key) == 1 {
			switch casted := state.(type) {
			case map[string]interface{}:
				delete(casted, currentKey)
				state = casted
			case []interface{}:
				idx, err := strconv.Atoi(currentKey)
				if err != nil {
					return err
				}
				state = append(casted[:idx], casted[idx+1:]...)
			default:
				return fmt.Errorf("cannot remove value from type %t", casted)
			}
		} else {
			if strings.HasPrefix(currentKey, "[") && strings.HasSuffix(currentKey, "]") {
				nested, ok := state.([]interface{})
				if !ok {
					return fmt.Errorf("cannot cast variable of type %T to type %T", state, nested)
				}
				idx, err := strconv.Atoi(currentKey)
				if err != nil {
					return err
				}
				next := nested[idx]
				return remove(next, key[1:])
			} else {
				nested, ok := state.(map[string]interface{})
				if !ok {
					return fmt.Errorf("cannot cast variable of type %T to type %T", state, nested)
				}
				next := nested[currentKey]
				return remove(next, key[1:])
			}
		}

		return nil
	}

	var search func(string, interface{}) error
	search = func(key string, value interface{}) error {
		switch casted := value.(type) {
		case map[string]interface{}:
			for nestedKey, nested := range casted {
				nextKey := key + "." + nestedKey
				err := search(nextKey, nested)
				if err != nil {
					return err
				}
			}
			return nil
		case []interface{}:
			for idx, nested := range casted {
				nextKey := key + ".[" + strconv.Itoa(idx) + "]"
				err := search(nextKey, nested)
				if err != nil {
					return err
				}
			}
			return nil
		case bool:
			if casted == true {
				err := remove(state, strings.Split(key, "."))
				if err != nil {
					return err
				}
			}
			return nil
		default:
			return nil
		}
	}

	for key, value := range sensitiveValues {
		err := search(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resource *Resource) addAction(action string) {
	resource.Actions = append(resource.Actions, action)
}

func (resource *Resource) MarshalBinary() ([]byte, error) {
	return json.Marshal(resource)
}

type ReferenceResource struct {
	*node
	Dependencies []string `json:"dependencies,omitempty"`
}

func NewReferenceResource(
	address string,
	children []string,
	dependencies []string,
) (*ReferenceResource, error) {
	referenceResource := new(ReferenceResource)
	referenceResource.node = newNodeData(address, Node_type_reference_resource, children)
	referenceResource.Dependencies = dependencies

	return referenceResource, nil
}

// Provider represents a terraform provider
type Provider struct {
	*node
}

func NewProvider(
	address string,
	children []string,
) (*Provider, error) {
	provider := new(Provider)
	provider.node = newNodeData(address, Node_type_provider, children)
	return provider, nil
}

func (provider *Provider) MarshalBinary() ([]byte, error) {
	return json.Marshal(provider)
}
