package spinner

type TestOptions struct {
    MaxAttempts int
    TimeElapsedDelta float64
}

// Update options using system defaults. Replace any unset values with their
// corresponding system default values.
func (t *TestOptions) UpdateDefaults() {
    if t.MaxAttempts == 0 {
        t.MaxAttempts = MAX_ATTEMPTS
    }
    if t.TimeElapsedDelta == 0 {
        t.TimeElapsedDelta = TIME_ELAPSED_DELTA
    }
}

// Update options using a set of default options as a pattern. Replace any
// unset values with their corresponding values from the defaults.
func (t *TestOptions) Update(defaults *TestOptions) {
    if t.MaxAttempts == 0 {
        t.MaxAttempts = defaults.MaxAttempts
    }
    if t.TimeElapsedDelta == 0 {
        t.TimeElapsedDelta = defaults.TimeElapsedDelta
    }
}
