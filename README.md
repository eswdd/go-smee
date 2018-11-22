# go-smee

Provides a Go client for smee.io

Based on previous work by [@cryptix](https://github.com/cryptix/goSSEClient)

## Usage

Import:
```
import "github.com/eswdd/go-smee"
```

Optionally create a new smee.io channel
```
source, err = CreateSmeeChannel()
```

Setup a new client:
```
target := make(chan SSEvent)
client := NewSmeeClient(source, target)

sub, err := client.Start()
if err != nil {
    panic(err)
}

// do what you want with the events
for ev := range target {
    fmt.Printf("Received event: id=%v, name=%v, payload=%v\n", ev.Id, ev.Name, string(ev.Data))
}
```

Experience suggests that if subscribing to Github hooks, you'll want to ignore any with `Name` set.

