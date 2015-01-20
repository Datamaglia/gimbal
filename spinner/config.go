package spinner

import (
    "encoding/json"
    "io/ioutil"
)

type ConfigSettings struct {
    ConcurrentRequests int
    SuppressSuccess bool
    SuppressWarning bool
    SuppressFailure bool
    SuppressOutput bool
}

func (c *ConfigSettings) UpdateDefaults() {
    if c.ConcurrentRequests < 1 {
        c.ConcurrentRequests = CONCURRENT_REQUESTS
    }
}

type TestConfig struct {
    Specs []*TestSpec
    DefaultRequest *RequestSpec
    DefaultResponse *ResponseSpec
    DefaultOptions *TestOptions
    Settings *ConfigSettings

    TotalTests int
    TotalFailures int
    TotalWarnings int
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
    if t.Settings == nil {
        t.Settings = new(ConfigSettings)
    }
    t.Settings.UpdateDefaults()

    if t.DefaultOptions == nil {
        t.DefaultOptions = new(TestOptions)
    }
    t.DefaultOptions.UpdateDefaults()

    if t.DefaultRequest == nil {
        t.DefaultRequest = new(RequestSpec)
    }
    if t.DefaultResponse == nil {
        t.DefaultResponse = new(ResponseSpec)
    }

    for _, spec := range t.Specs {
        // Request
        spec.Request.Update(t.DefaultRequest)
        spec.Response.Update(t.DefaultResponse)

        // Options
        if spec.Options == nil {
            spec.Options = new(TestOptions)
        }
        options := spec.Options
        options.Update(t.DefaultOptions)
    }
}
