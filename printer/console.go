package printer

import (
	"github.com/wsxiaoys/terminal/color"
)

func ResultsToConsole(resultSet *ResultSet) {
	prefix := ""
	switch resultSet.Status() {
	case SUCCESS:
		prefix = "@g\u2713 "
	case WARNING:
		prefix = "@y\u2713 "
	case FAILURE:
		prefix = "@r\u2718 "
	}

	if resultSet.Spec.Name == "" {
		color.Printf(prefix+"%v\n", resultSet.Spec.Url())
	} else {
		color.Printf(prefix+"%v\n", resultSet.Spec.Name)
	}

	for _, result := range resultSet.Results {
		symbolPrefix := ""
		colorPrefix := ""
		switch result.Status {
		case SUCCESS:
			symbolPrefix = "@g  \u2713 "
			colorPrefix = "@g"
		case WARNING:
			symbolPrefix = "@y  \u2713 "
			colorPrefix = "@y"
		case FAILURE:
			symbolPrefix = "@r  \u2718 "
			colorPrefix = "@r"
		}
		color.Printf(symbolPrefix+"%v\n", result.Message)
		color.Printf(colorPrefix+"    Expected: %v\n", result.Expected)
		color.Printf(colorPrefix+"    Observed: %v\n", result.Observed)
	}
}
