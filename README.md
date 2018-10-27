# pathfinder
Go package to match path to path template and extract path params

> Static path parts take precedence over parametrized parts. I.e. static parts are matched first and if matched parametrized parts are not looked up.

Parametrized part is marked up with `:` and will be the key in matched params with leading `:`.

Example:
```go
n := &pathfinder.Node{}

n.Add("hello/world", "world") // static route
n.Add(":p0/:p1", "param") // parametrized route

payload, params, _ := n.Lookup("hello/world") // will match static route without params
payload, params, _ := n.Lookup("param0/param1") // will match parametrized route with :p0=param0 and :p1=param1
payload, params, _ := n.Lookup("hello/param1") // will not match any route, as the first part will match a static part 'hello'
```
