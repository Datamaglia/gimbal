package main

import (
    "github.com/datamaglia/gimbal/spinner"
)

func main() {
    config := spinner.LoadJsonConfig("test.json")
    spinner.ExecuteTestConfig(config)
}
