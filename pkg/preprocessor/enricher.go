package preprocessor

import (
	"fmt"
	"path"

	tfjson "github.com/hashicorp/terraform-json"
)

func addResourceChangesInformation(nodeTable map[string]Node, tfjsonResourceChanges []*tfjson.ResourceChange, stateToAdd string) (map[string]Node, error) {
	var err error
	for _, tfjsonResourceChange := range tfjsonResourceChanges {
		address := RootAddress + "." + tfjsonResourceChange.Address
		node := nodeTable[address]
		resource, ok := node.(*Resource)
		if !ok {
			err := fmt.Errorf("could not cast Node to Resource")
			return nil, err
		}
		resource = addActionsToNode(tfjsonResourceChange, resource)

		requestedStateIsNotPresent := (tfjsonResourceChange.Change.Before != nil && stateToAdd == State_current) ||
			(tfjsonResourceChange.Change.After != nil && stateToAdd == State_planned)
		if requestedStateIsNotPresent {
			resource, err = addStateFromResourceChangeToNode(tfjsonResourceChange, resource, stateToAdd)
			if err != nil {
				return nil, err
			}
		}

		nodeTable[address] = Node(resource)
	}
	return nodeTable, nil
}

func addActionsToNode(tfjsonResourceChange *tfjson.ResourceChange, resource *Resource) *Resource {
	if tfjsonResourceChange.Change.Actions.Create() {
		resource.addAction("Create")
	}
	if tfjsonResourceChange.Change.Actions.CreateBeforeDestroy() {
		resource.addAction("CreateBeforeDestroy")
	}
	if tfjsonResourceChange.Change.Actions.Delete() {
		resource.addAction("Delete")
	}
	if tfjsonResourceChange.Change.Actions.DestroyBeforeCreate() {
		resource.addAction("DestroyBeforeCreate")
	}
	if tfjsonResourceChange.Change.Actions.NoOp() {
		resource.addAction("NoOp")
	}
	if tfjsonResourceChange.Change.Actions.Read() {
		resource.addAction("Read")
	}
	if tfjsonResourceChange.Change.Actions.Replace() {
		resource.addAction("Replace")
	}
	if tfjsonResourceChange.Change.Actions.Update() {
		resource.addAction("Update")
	}
	return resource
}

func addStateFromResourceChangeToNode(tfjsonResourceChange *tfjson.ResourceChange, resource *Resource, stateToAdd string) (*Resource, error) {
	switch stateToAdd {
	case State_current:
		priorValues := tfjsonResourceChange.Change.Before
		priorValuesMap, ok := priorValues.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("could not cast values \\'Before\\' to map[string]interface{}")
		}
		resource.addState(stateToAdd, priorValuesMap)
		sensitiveValues, ok := tfjsonResourceChange.Change.BeforeSensitive.(map[string]interface{})
		if ok {
			err := resource.removeSensitiveValues(stateToAdd, sensitiveValues)
			if err != nil {
				return nil, err
			}
		}
	case State_planned:
		planedValues := tfjsonResourceChange.Change.After
		planedValuesMap, ok := planedValues.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("could not cast values \\'After\\' to map[string]interface{}")
		}
		resource.addState(stateToAdd, planedValuesMap)
		sensitiveValues, ok := tfjsonResourceChange.Change.AfterSensitive.(map[string]interface{})
		if ok {
			err := resource.removeSensitiveValues(stateToAdd, sensitiveValues)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("provided state \\'%s\\' is not valid", stateToAdd)
	}
	return resource, nil
}

func addConfigInformation(nodeTable map[string]Node, tfjsonConfig *tfjson.Config, url string, basepath string) (map[string]Node, error) {
	var err error
	if tfjsonConfig.ProviderConfigs != nil {
		nodeTable, err := addProviderInfo(nodeTable, tfjsonConfig.ProviderConfigs, RootAddress)
		if err != nil {
			return nodeTable, err
		}
	}
	rootModule := tfjsonConfig.RootModule
	if rootModule != nil {
		location := url + basepath
		nodeTable[RootAddress].SetLocation(location)

		tfjsonConfigResources := rootModule.Resources
		if tfjsonConfigResources != nil {
			nodeTable, err = addConfigResourceInfoToNodes(nodeTable, tfjsonConfigResources, RootAddress, location)
			if err != nil {
				return nodeTable, err
			}
		}
		tfjsonConfigChildModuleCalls := rootModule.ModuleCalls
		for childModuleName, tfjsonConfigChildModuleCall := range tfjsonConfigChildModuleCalls {
			childModuleAddress := RootAddress + ".module." + childModuleName
			addConfigModuleInfoToNodes(nodeTable, tfjsonConfigChildModuleCall, childModuleAddress, url, basepath)

		}
	}
	return nodeTable, nil
}

func addConfigModuleInfoToNodes(nodeTable map[string]Node, tfjsonModuleCall *tfjson.ModuleCall, address string, url string, currentPath string) (map[string]Node, error) {
	var err error

	location := url + currentPath
	if _, ok := nodeTable[address]; !ok {
		module, err := NewModule(address, []string{})
		if err != nil {
			return nil, err
		}
		nodeTable[address] = module
	}
	nodeTable[address].SetLocation(location)
	childPath := path.Join(currentPath, tfjsonModuleCall.Source)

	nodeTable[address].SetLocation(location)

	tfjsonConfigModule := tfjsonModuleCall.Module

	tfjsonConfigResources := tfjsonConfigModule.Resources
	childLocation := url + childPath
	nodeTable, err = addConfigResourceInfoToNodes(nodeTable, tfjsonConfigResources, address, childLocation)
	if err != nil {
		return nodeTable, err
	}

	_ = tfjsonConfigModule.Resources
	tfjsonConfigChildModuleCalls := tfjsonConfigModule.ModuleCalls
	for childModuleName, tfjsonConfigChildModuleCall := range tfjsonConfigChildModuleCalls {
		childModuleAddress := address + ".module." + childModuleName
		nodeTable, err = addConfigModuleInfoToNodes(nodeTable, tfjsonConfigChildModuleCall, childModuleAddress, url, childPath)
		if err != nil {
			return nodeTable, err
		}
	}
	return nodeTable, nil
}

func addConfigResourceInfoToNodes(nodeTable map[string]Node, tfjsonConfigResources []*tfjson.ConfigResource, parentAddress string, location string) (map[string]Node, error) {
	for _, tfjsonConfigResource := range tfjsonConfigResources {
		address := parentAddress + "." + tfjsonConfigResource.Address
		if _, ok := nodeTable[address]; !ok {
			module, err := NewResource(address, []string{})
			if err != nil {
				return nil, err
			}
			nodeTable[address] = module
		}
		nodeTable[address].SetLocation(location)
	}
	return nodeTable, nil
}

func addProviderInfo(nodeTable map[string]Node, tfjsonProviderConfigs map[string]*tfjson.ProviderConfig, parentAdress string) (map[string]Node, error) {
	for _, tfjsonProviderConfig := range tfjsonProviderConfigs {
		address := parentAdress + ".provider." + tfjsonProviderConfig.Name
		provider, err := NewProvider(address, []string{})
		if err != nil {
			return nodeTable, err
		}
		provider.AddAttribute("alias", tfjsonProviderConfig.Alias)
		provider.AddAttribute("name", tfjsonProviderConfig.Name)
		provider.AddAttribute("versionConstraint", tfjsonProviderConfig.VersionConstraint)
		provider.AddAttribute("moduleAddress", tfjsonProviderConfig.ModuleAddress)
		provider.AddAttribute("attributes", tfjsonProviderConfig.Expressions)
		nodeTable[address] = Node(provider)
	}
	return nodeTable, nil
}
