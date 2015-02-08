package printer

import (
	"github.com/datamaglia/gimbal/wrapper"
	"github.com/wsxiaoys/terminal/color"
)

const successSymbol string = "\u2713"
const failureSymbol string = "\u2718"

func ResultsToConsole(w *wrapper.Wrapper) {
	prefix := ""
	switch w.Status() {
	case wrapper.SUCCESS:
		prefix = "@g" + successSymbol + " "
	case wrapper.WARNING:
		prefix = "@y" + successSymbol + " "
	case wrapper.FAILURE:
		prefix = "@r" + failureSymbol + " "
	}

	if w.Spec.Name == "" {
		color.Printf(prefix+"%v\n", w.Spec.Url())
	} else {
		color.Printf(prefix+"%v\n", w.Spec.Name)
	}

	for _, result := range w.Results {
		symbolPrefix := ""
		colorPrefix := ""
		switch result.Status {
		case wrapper.SUCCESS:
			symbolPrefix = "@g  " + successSymbol + " "
			colorPrefix = "@g"
		case wrapper.WARNING:
			symbolPrefix = "@y  " + successSymbol + " "
			colorPrefix = "@y"
		case wrapper.FAILURE:
			symbolPrefix = "@r  " + failureSymbol + " "
			colorPrefix = "@r"
		}
		color.Printf(symbolPrefix+"%v\n", result.CheckName)
		color.Printf(colorPrefix+"    Observed: %v\n", result.Observed)

		if result.Status != wrapper.SUCCESS {
			color.Printf(colorPrefix+"    Expected: %v\n", result.Expected)
		}
	}
}
