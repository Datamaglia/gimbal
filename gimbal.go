package main

import (
    "flag"
    "os"

    "github.com/datamaglia/gimbal/spinner"
)

var filename = flag.String("f", "", "Read the config from a file")
var quiet = flag.Bool("q", false, "Suppress all output")

func main() {
    flag.Parse()

    config := spinner.LoadJsonConfig(*filename)
    if *quiet {
        config.Settings.SuppressOutput = true
    }

    spinner.ExecuteTestConfig(config)

    if config.TotalWarnings != 0 || config.TotalFailures != 0 {
        os.Exit(10)
    } else {
        os.Exit(0)
    }
}
