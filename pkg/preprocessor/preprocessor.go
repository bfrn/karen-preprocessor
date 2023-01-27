package preprocessor

import (
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

const (
	RootAddress = "_root"
)

//ParseStateFile takes a json formatted state file and generates a node table from it.
func ParseStateFile(stateFile []byte) (map[string]Node, error) {
	tfjsonState := new(tfjson.State)
	err := tfjsonState.UnmarshalJSON(stateFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse the given state file: %s", err.Error())
	}
	nodeTable := make(map[string]Node)

	rootModule := tfjsonState.Values.RootModule
	rootModule.Address = RootAddress

	nodeTable, err = parseTfjsonStateModule(rootModule, nodeTable, State_current, RootAddress)
	if err != nil {
		return nil, err
	}
	return nodeTable, nil
}

// ParsePlanFile takes a json formatted plan file and generates a node table from it.
func ParsePlanFile(planFile []byte, tfConfigUrl string, tfConfigMainPath string) (map[string]Node, error) {
	plan := new(tfjson.Plan)
	err := plan.UnmarshalJSON(planFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse the given plan file: %s", err.Error())
	}
	nodeTable := make(map[string]Node)

	isCurrentStatePresent := plan.PriorState != nil

	var rootModule *tfjson.StateModule
	var state string
	if isCurrentStatePresent {
		rootModule = plan.PriorState.Values.RootModule
		state = State_current
	} else {
		rootModule = plan.PlannedValues.RootModule
		state = State_planned
	}

	rootModule.Address = RootAddress
	nodeTable, err = parseTfjsonStateModule(rootModule, nodeTable, state, RootAddress)
	if err != nil {
		return nil, err
	}

	if isCurrentStatePresent && (plan.ResourceChanges != nil) {
		nodeTable, err = addResourceChangesInformation(nodeTable, plan.ResourceChanges, State_planned)
		if err != nil {
			return nil, err
		}
	}
	if plan.Config != nil {
		nodeTable, err = addConfigInformation(nodeTable, plan.Config, tfConfigUrl, tfConfigMainPath)
		if err != nil {
			return nil, err
		}
	}

	return nodeTable, nil
}

// EnrichStateFile takes an existing node table and enriches the nodes with the information that the provided json formatted plan file contains.
// It is assumed that the provided node Table only contains the planned values from a tfjson state file.
func EnrichStateFile(nodeTable map[string]Node, planFile []byte, tfConfigUrl string, tfConfigMainPath string) (map[string]Node, error) {
	plan := new(tfjson.Plan)
	err := plan.UnmarshalJSON(planFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse the given plan file: %s", err.Error())
	}

	isPlannedStatePresent := plan.PlannedValues != nil
	if isPlannedStatePresent && (plan.ResourceChanges != nil) {
		nodeTable, err = addResourceChangesInformation(nodeTable, plan.ResourceChanges, State_planned)
		if err != nil {
			return nil, err
		}
	}
	if plan.Config != nil {
		nodeTable, err = addConfigInformation(nodeTable, plan.Config, tfConfigUrl, tfConfigMainPath)
		if err != nil {
			return nil, err
		}
	}

	return nodeTable, nil
}
