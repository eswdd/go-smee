# go-smee

Provides a Go client for [smee.io](https://smee.io).

Mimics the interface provided by <https://github.com/probot/smee>.

Based on previous work by [@cryptix](https://github.com/cryptix/goSSEClient).

## Usage

Import:

```go
import "github.com/eswdd/go-smee/smee"
```

Optionally create a new smee.io channel:

```go
source, err := smee.CreateSmeeChannel()
```

Or you can use your own [SMEE
server](https://github.com/probot/smee.io#deploying-your-own-smeeio):

```go
source, err := smee.CreateSmeeChannelFromURL("http://localhost:3000/new")
```

Setup a new client:

```go
target := make(chan smee.SSEvent)
client := smee.NewSmeeClient(source, target)

sub, err := client.Start()
if err != nil {
    panic(err)
}

// do what you want with the events
for ev := range target {
    fmt.Printf("Received event: id=%v, name=%v, payload=%v\n", ev.Id, ev.Name, string(ev.Data))
}
```

Experience suggests that if subscribing to Github hooks, you'll want to ignore
any with `Name` set.
