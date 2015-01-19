package main

import (
    "flag"

    "github.com/datamaglia/gimbal/spinner"
)

var filename = flag.String("f", "", "Read the config from a file")

func main() {
    flag.Parse()

    config := spinner.LoadJsonConfig(*filename)
    spinner.ExecuteTestConfig(config)
}
