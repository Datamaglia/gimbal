package spec

import (
    "encoding/json"
    "io/ioutil"
)

func LoadJsonFile(filename string) *Spec {
    raw, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }

    return loadJsonString(raw)
}

func loadJsonString(content []byte) *Spec {
    spec := new(Spec)
    err := json.Unmarshal(content, spec)
    if err != nil {
        panic(err)
    }

    spec.setDefaults()

    return spec
}
