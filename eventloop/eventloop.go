package main

import (
	"bpm/log"
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

type Runnable[T interface{}] func(event Event[T])

// Event structure of the data sent through a Channel
type Event[T interface{}] struct {
	Data         T
	Id           string
	Channel      *Channel[T]
	Subscription *Subscription[T]
}

// EventLoop manages all Event that needs to be send to Channel through Subscription
type EventLoop[T interface{}] struct {
	Channels map[string]*Channel[T]
}

// CreateChannel will create a new Channel in this event loop which can be used to send Event to it
func (el *EventLoop[T]) CreateChannel(name string) *Channel[T] {
	ch := &Channel[T]{Name: name, Subs: []*Subscription[T]{}, Ch: make(chan Event[T], 1024)}
	el.Channels[name] = ch
	go ch.Listen()
	return ch
}

// Register subscribe a Subscription to events published to the Channel specified
func (el *EventLoop[T]) Register(ch, id string, action Runnable[T]) (*Subscription[T], error) {
	if channel, ok := el.Channels[ch]; ok {
		log.Debugf("registering new listener '%s' on channel '%s'\n", id, ch)
		return channel.Subscribe(id, action), nil
	}
	return nil, fmt.Errorf("unknown channel '%s'\n", ch)

}

// Push a new Event on the specified Channel. All Subscription will be notified : a new event will be dispatched
// to all Subscription that are listening on the specified Channel
func (el *EventLoop[T]) Push(ch string, data T) (string, error) {
	if channel, ok := el.Channels[ch]; ok {
		return channel.Push(data), nil
	}
	return "", fmt.Errorf("unknown channel '%s'\n", ch)
}

// Remove a Subscription from the specified Channel of this event loop. Returns an error if the Channel does not exist,
// or if there is no Subscription with that id for that Channel
func (el *EventLoop[T]) Remove(ch, id string) error {
	if channel, ok := el.Channels[ch]; ok {
		log.Debugf("removing listener '%s' from channel '%s'\n", id, channel.Name)
		return channel.Remove(id)
	} else {
		return fmt.Errorf("unknown channel '%s'\n", ch)
	}
}

// Close this event loop. Stops all Channel from receiving Event and dispatching those to Subscription. Calling Push
// after Close will result in an error
func (el *EventLoop[T]) Close() {
	for _, channel := range el.Channels {
		close(channel.Ch)
		for _, subs := range channel.Subs {
			close(subs.Ch)
		}
	}
}

type Subscription[T interface{}] struct {
	Id     string
	Action Runnable[T]
	Ch     chan Event[T]
}

func (s *Subscription[T]) Listen() {
	for e := range s.Ch {
		log.Debugf("received event '%s' in subscription '%s'\n", e.Id, s.Id)
		e.Subscription = s
		s.Action(e)
		log.Debugf("event ")
	}
	log.Debugf("stop listening on subscription '%s'", s.Id)
}

type Channel[T interface{}] struct {
	Name string
	Subs []*Subscription[T]
	Ch   chan Event[T]
}

// Subscribe creates a new Subscription for this Channel.
// When data is Push to this Channel, the Runnable will execute asynchronously
func (c *Channel[T]) Subscribe(id string, action Runnable[T]) *Subscription[T] {
	s := Subscription[T]{
		Id:     id,
		Action: action,
		Ch:     make(chan Event[T], 1024),
	}
	c.Subs = append(c.Subs, &s)
	go s.Listen()
	return &s
}

// Remove removes the Subscription with the provided id from this channel
func (c *Channel[T]) Remove(id string) error {
	for i, ch := range c.Subs {
		if ch.Id == id {
			c.Subs[i] = c.Subs[len(c.Subs)-1]
			c.Subs = c.Subs[:len(c.Subs)-1]
			return nil
		}
	}
	return fmt.Errorf("subscription '%s' not found for channel '%s'", id, c.Name)
}

// Push data to the Channel
func (c *Channel[T]) Push(data T) string {
	id := uuid.New()
	log.Debugf("pushing event '%s' to channel '%s'\n", id, c.Name)
	c.Ch <- Event[T]{
		Data:    data,
		Id:      id.String(),
		Channel: c,
	}
	log.Debugf("'%s' pushed\n", id)
	return id.String()
}

// Listen start listening for data on this Channel
func (c *Channel[T]) Listen() {
	for e := range c.Ch {
		log.Debugf("received '%s' on channel %s\n", e.Id, c.Name)
		for _, subs := range c.Subs {
			log.Debugf("sending to subscription '%s'\n", subs.Id)
			subs.Ch <- e
			log.Debugf("sent\n")
		}
	}
	log.Debugf("Stop listening on channel '%s'", c.Name)
}

func New[T interface{}]() EventLoop[T] {
	channels := make(map[string]*Channel[T])
	return EventLoop[T]{Channels: channels}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("====== Start listening =====")
	eventLoop := New[string]()
loop:
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Split(text, "\n")[0]
		splits := strings.Split(text, " ")
		switch splits[0] {
		case "":
			continue
		case "stop":
			break loop
		case "close":
			eventLoop.Close()

		case "create":
			{
				if len(splits) < 2 {
					log.Error("command 'create' requires 1 argument: create <channel>")
					break
				}
				name := splits[1]
				log.Debugf("creating new channel '%s'\n", name)
				eventLoop.CreateChannel(name)
				continue
			}

		case "register":
			{
				if len(splits) < 3 {
					log.Error("command 'register' requires 2 arguments: register <channel> <listener id>")
					break
				}
				ch := splits[1]
				id := splits[2]
				if _, err := eventLoop.Register(ch, id, func(e Event[string]) {
					fmt.Printf("[%s:%s:%s] %s\n", e.Channel.Name, e.Subscription.Id, e.Id, e.Data)
				}); err != nil {
					log.ErrorE(err)
				}
				continue
			}

		case "push":
			{
				if len(splits) < 3 {
					log.Error("command 'push' requires 2 arguments: push <channel> <data>")
					break
				}
				ch := splits[1]
				data := strings.Join(splits[2:], " ")
				if _, err := eventLoop.Push(ch, data); err != nil {
					log.ErrorE(err)
				}
				continue
			}

		case "remove":
			{
				if len(splits) < 3 {
					log.Error("command 'remove' requires 2 arguments: remove <channel> <listener id>")
					break
				}
				ch := splits[1]
				id := splits[2]
				if err := eventLoop.Remove(ch, id); err != nil {
					log.ErrorE(err)
				}
			}

		case "debug":
			{
				for _, channel := range eventLoop.Channels {
					log.Infof("%s: %#v\n", channel.Name, *channel)
				}
			}

		case "log":
			{
				if len(splits) < 2 {
					log.Errorf("command 'log' requires 1 argument: log <log level>")
					break
				}
				log.SetGlobalLogLevelFromString(splits[1])
			}
		default:
			log.Errorf("unknown command '%s'\n", splits[0])
		}
	}
	fmt.Printf("===== Stop listening =====\n")
}
