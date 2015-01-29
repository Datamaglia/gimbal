# Gimbal

Gimbal is a system for testing HTTP APIs and web sites. It allows the developer
to create a set of specifications of the general form "send X, expect Y back".

See <http://gimbal.datamaglia.com/> for more information.

## Specification

### name

String. A descriptive name for the spec. Not used internally but will be logged
in some cases for debugging and on failures. Parent names will be prepended, so
if a spec named A has a child spec named B when the child spec is run its name
will appear as "A :: B". Default is "".

### concurrentRequests

Integer. The maximum number of concurrent requests that will be made. If this is set to
1, requests will be run in serial. Default is 1.

### outputLevel

String. The type of messages that should be included in the output. One of
"success", "warning", and "failure". Levels include subsequent levels. Default
is "success".

### maxAttempts

Integer. The maximum number of times to attempt a request before it is declared
a failure. Default is 10.

### timeElapsedDelta

Float. If the time a request takes is greater than `maxTimeElapsed`, but by less
than this amount, the result will be a warning, not a failure. This allows for
some "wiggle room". Default is 0.0.

### host

String. The host to contact, such as `httpbin.org`.

### port

Integer. The port on which to contact the host. Default is 80 when `SSL` is
false and 443 when `SSL` is true.

### ssl

Boolean. Whether to use SSL for the request. Default is false.

### method

String. The HTTP verb to use for the request. Default is "GET".

### statusCode

Integer. The expected status code of the response. Default is 200.

### maxTimeElapsed

Float. The maximum time the request can take and still be considered a success,
in seconds. Default is 10.0.

### requestHeaders

Object. Mapping of header names to lists of header values that will be included
in the request. Default is `{}`.

### responseHeaders

Object. Mapping of header names to lists of header values that are expected to
be included in the response. Unless `exactResponseHeaders` is true, the response
can contain additional headers, so this represents a minimum. Default is `{}`.

### exactResponseHeaders

Boolean. If true, the `responseHeaders` value will be treated as an exact
enumeration of headers and values, if the response includes additional headers
the spec will fail. Default is false.

### uri

String. The URI to request from the host. Default is "/".

### specs

List. A recursive list of specs to run. Each of these will inherit properties
from the parent as necessary. If this is set, then the spec that contains the
list will not run on its own. In other words, if you think about the
configuration as a tree, only leaf nodes actually cause a test to be run.
Default is [].

### requestData

String. The raw data to include with the request. Default is "".

### responseData

String. The raw data that should be included with the response. If set to the
empty string, no check is performed. Default is "".

### vars (not yet implemented)

Object. A map of variable assignments that can be used in the given spec and all
its children. The keys must be strings. For now, the values must be strings or
numbers. Eventually lists will be supported and a `forEach` option will be added
to run a given spec once for each value in the list. Default is `{}`.
