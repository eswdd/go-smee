package smee

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

// CreateChannel creates a new channel on smee.io.
func CreateChannel() (*string, error) {
	loc, err := CreateChannelFromURL("https://smee.io/new")

	return loc, err
}

// CreateChannelFromURL creates a new channel on a specific server.
func CreateChannelFromURL(newURL string) (*string, error) {
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := httpClient.Head(newURL)
	if err != nil {
		return nil, err
	}

	loc := resp.Header.Get("Location")
	return &loc, nil
}

// ClientSubscription holds a channel required for termination of the client subscription.
type ClientSubscription struct {
	terminator chan<- interface{}
}

// Stop terminates the client subscription.
func (c *ClientSubscription) Stop() {
	c.terminator <- nil
}

// Client holds information about the source, target and logger.
type Client struct {
	source *string
	target chan<- SSEvent
	logger *log.Logger
}

// NewClient creates a new client.
func NewClient(source *string, target chan<- SSEvent, logger *log.Logger) *Client {
	c := new(Client)
	c.source = source
	c.target = target
	c.logger = logger
	return c
}

// Start starts the client.
func (c *Client) Start() (*ClientSubscription, error) {
	eventStream, err := openSSEUrl(*c.source)
	if err != nil {
		return nil, err
	}

	quit := make(chan interface{})
	go c.run(eventStream, quit)

	return &ClientSubscription{terminator: quit}, nil
}

func (c *Client) run(sseEventStream <-chan SSEvent, quit <-chan interface{}) {
	for {
		select {
		case event := <-sseEventStream:
			c.target <- event
		case <-quit:
			return
		}
	}
}

// SSEvent holds information about the server-sent event.
type SSEvent struct {
	ID   string
	Name string
	Data []byte
}

func openSSEUrl(url string) (<-chan SSEvent, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "text/event-stream")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("resp.StatusCode == %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "text/event-stream" {
		return nil, fmt.Errorf("invalid Content-Type == %s", resp.Header.Get("Content-Type"))
	}

	events := make(chan SSEvent)

	var buf bytes.Buffer

	go func() {
		ev := SSEvent{}
		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Bytes()

			switch {

			// start of event
			case bytes.HasPrefix(line, []byte("id:")):
				ev.ID = string(line[4:])

				// event name
			case bytes.HasPrefix(line, []byte("event:")):
				ev.Name = string(line[7:])

				// event data
			case bytes.HasPrefix(line, []byte("data:")):
				buf.Write(line[6:])

				// end of event
			case len(line) == 0:
				ev.Data = buf.Bytes()
				buf.Reset()
				events <- ev
				ev = SSEvent{}

			default:
				fmt.Fprintf(os.Stderr, "Error during EventReadLoop - Default triggerd! len:%d\n%s", len(line), line)
				close(events)

			}
		}

		if err = scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error during resp.Body read:%s\n", err)
			close(events)

		}
	}()

	return events, nil
}
