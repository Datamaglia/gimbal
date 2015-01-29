package printer

import (
    "fmt"

    "github.com/datamaglia/gimbal/runner"
)

func ResultsToConsole(resultSet *runner.ResultSet) {
    for _, result := range resultSet.Results {
        if (result.Spec.Name == "") {
            fmt.Printf("%s\n", result.Spec.Url())
        } else {
            fmt.Printf("%s\n", result.Spec.Name)
        }
    }
}
