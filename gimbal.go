package main

import (
    "flag"
//    "os"

    "github.com/datamaglia/gimbal/spec"
    "github.com/datamaglia/gimbal/runner"
)

var filename = flag.String("f", "", "Read the config from a file")
var quiet = flag.Bool("q", false, "Suppress all output")

func main() {
    flag.Parse()

    config := spec.LoadJsonFile(*filename)
    runner.RunSpec(config)
}
