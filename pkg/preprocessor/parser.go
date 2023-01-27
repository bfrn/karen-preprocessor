package preprocessor

import (
	"encoding/json"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

func parseTfjsonStateModule(tfjsonModule *tfjson.StateModule, nodeTable map[string]Node, state string, parent string) (map[string]Node, error) {

	address := RootAddress + "." + tfjsonModule.Address
	if tfjsonModule.Address == RootAddress {
		address = strings.Split(address, ".")[0]
	}

	module, err := NewModule(address, nil)
	if err != nil {
		return nil, err
	}
	nodeTable[module.Address] = Node(module)

	nodeTable, err = parseTfjsonStateResource(tfjsonModule.Resources, nodeTable, state, module.Address)
	if err != nil {
		return nil, err
	}

	for _, tfjsonChildModule := range tfjsonModule.ChildModules {
		childAddress := RootAddress + "." + tfjsonChildModule.Address
		module.Children = append(module.Children, childAddress)
		nodeTable, err = parseTfjsonStateModule(tfjsonChildModule, nodeTable, state, module.Address)
		if err != nil {
			return nil, err
		}
	}
	return nodeTable, nil
}

func parseTfjsonStateResource(tfjsonResources []*tfjson.StateResource, nodeTable map[string]Node, state string, parent string) (map[string]Node, error) {
	for _, tfjsonRessource := range tfjsonResources {

		var dependencies []string
		for _, dependency := range tfjsonRessource.DependsOn {
			dependencyAddress := RootAddress + "." + dependency
			dependencies = append(dependencies, dependencyAddress)
		}

		resource, err := NewResource(RootAddress+"."+tfjsonRessource.Address, dependencies)
		if err != nil {
			return nil, err
		}
		resource.addState(state, tfjsonRessource.AttributeValues)

		if tfjsonRessource.SensitiveValues != nil {
			var sensitiveValues map[string]interface{}
			jsonData, err := tfjsonRessource.SensitiveValues.MarshalJSON()
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(jsonData, &sensitiveValues)
			if err != nil {
				return nil, err
			}
			err = resource.removeSensitiveValues(state, sensitiveValues)
			if err != nil {
				return nil, err
			}
		}
		if strings.HasSuffix(resource.Address, "]") {
			splittedAddress := strings.Split(resource.Address, "[")
			referenceResourceAddress := strings.Join(splittedAddress[0:len(splittedAddress)-1], "[")

			if referenceResource, containsReferenceResource := nodeTable[referenceResourceAddress]; containsReferenceResource {
				referenceResource.AddChild(resource.Address)
			} else {
				referenceResource, err := NewReferenceResource(referenceResourceAddress, []string{resource.Address}, dependencies)
				if err != nil {
					return nil, err
				}
				nodeTable[referenceResource.Address] = referenceResource
				nodeTable[parent].AddChild(referenceResource.Address)
			}
		} else {
			nodeTable[parent].AddChild(resource.Address)
		}
		nodeTable[resource.Address] = resource
	}
	return nodeTable, nil
}
