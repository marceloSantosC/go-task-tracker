package model

import (
	"fmt"
	"testing"
)

func TestTaskStatus(t *testing.T) {

	var testTable = []struct {
		status   TaskStatus
		expected string
	}{
		{0, "To do"},
		{1, "In progress"},
		{2, "Done"},
	}

	for _, testData := range testTable {

		testName := fmt.Sprintf("For Input (%d), Expect: %s", testData.status, testData.expected)

		t.Run(testName, func(t *testing.T) {
			answer := testData.status.String()
			if answer != testData.expected {
				t.Errorf("with input (%d) got %s, but expected %s", testData.status, answer, testData.expected)
			}
		})
	}
}
