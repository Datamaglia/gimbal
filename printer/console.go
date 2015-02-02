package printer

import (
	"fmt"
)

func ResultsToConsole(resultSet *ResultSet) {
	if resultSet.Spec.Name == "" {
		fmt.Printf("%v\n", resultSet.Spec.Url())
	} else {
		fmt.Printf("%v\n", resultSet.Spec.Name)
	}
	for _, result := range resultSet.Results {
		fmt.Printf("  %v\n", result.Message)
		fmt.Printf("    Expected: %v\n", result.Expected)
		fmt.Printf("    Observed: %v\n", result.Observed)
	}
}
