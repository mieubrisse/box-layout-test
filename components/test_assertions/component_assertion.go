package test_assertions

import (
	"github.com/mieubrisse/box-layout-test/components"
	"testing"
)

type ComponentAssertion interface {
	Check(t *testing.T, component components.Component)
}

func FlattenAssertionGroups(assertionGroups ...[]ComponentAssertion) []ComponentAssertion {
	numAssertions := 0
	for _, group := range assertionGroups {
		numAssertions += len(group)
	}

	result := make([]ComponentAssertion, 0, numAssertions)
	for _, group := range assertionGroups {
		result = append(result, group...)
	}
	return result
}

func CheckAll(t *testing.T, assertionGroup []ComponentAssertion, component components.Component) {
	for _, assertion := range assertionGroup {
		assertion.Check(t, component)
	}
}
