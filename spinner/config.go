package spinner

import (
    "encoding/json"
    "io/ioutil"
)

type ConfigSettings struct {
    ConcurrentRequests int
}

type TestConfig struct {
    Specs []*TestSpec
    DefaultRequest *RequestSpec
    DefaultOptions *TestOptions
    Settings *ConfigSettings
}

// Load a configuration from a file in JSON format.
func LoadJsonConfig(filename string) *TestConfig {
    raw, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }

    config := new(TestConfig)
    err = json.Unmarshal(raw, config)
    if err != nil {
        // TODO Provide a more pleasant error message
        panic(err)
    }

    config.SetDefaults()

    return config
}

// Update default options and request values from the system defaults, then
// update the options for spec from those default values.
func (t *TestConfig) SetDefaults() {
    if t.DefaultOptions == nil {
        t.DefaultOptions = new(TestOptions)
    }
    defaults := t.DefaultOptions
    defaults.UpdateDefaults()

    if t.DefaultRequest == nil {
        t.DefaultRequest = new(RequestSpec)
    }

    for _, spec := range t.Specs {
        // Request
        spec.Request.Update(t.DefaultRequest)

        // Options
        if spec.Options == nil {
            spec.Options = new(TestOptions)
        }
        options := spec.Options
        options.Update(defaults)
    }
}
