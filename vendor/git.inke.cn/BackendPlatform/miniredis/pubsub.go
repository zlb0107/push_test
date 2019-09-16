package miniredis

import (
	"regexp"
	"sort"
	"sync"

	"git.inke.cn/BackendPlatform/miniredis/server"
)

// PubsubMessage is what gets broadcasted over pubsub channels.
type PubsubMessage struct {
	Channel string
	Message string
}

// Subscriber has the (p)subscriptions.
type Subscriber struct {
	publish  chan PubsubMessage
	channels map[string]struct{}
	patterns map[string]*regexp.Regexp
	mu       sync.Mutex
}

// Make a new subscriber. The channel is not buffered, so you will need to keep
// reading using Messages(). Use Close() when done, or unsubscribe.
func newSubscriber() *Subscriber {
	return &Subscriber{
		publish:  make(chan PubsubMessage),
		channels: map[string]struct{}{},
		patterns: map[string]*regexp.Regexp{},
	}
}

// Close the listening channel
func (s *Subscriber) Close() {
	close(s.publish)
}

// Count the total number of channels and patterns
func (s *Subscriber) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.count()
}

func (s *Subscriber) count() int {
	return len(s.channels) + len(s.patterns)
}

// Subscribe to a channel. Returns the total number of (p)subscriptions after
// subscribing.
func (s *Subscriber) Subscribe(c string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.channels[c] = struct{}{}
	return s.count()
}

// Unsubscribe a channel. Returns the total number of (p)subscriptions after
// unsubscribing.
func (s *Subscriber) Unsubscribe(c string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.channels, c)
	return s.count()
}

// Subscribe to a pattern. Returns the total number of (p)subscriptions after
// subscribing.
func (s *Subscriber) Psubscribe(pat string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.patterns[pat] = compileChannelPattern(pat)
	return s.count()
}

// Unsubscribe a pattern. Returns the total number of (p)subscriptions after
// unsubscribing.
func (s *Subscriber) Punsubscribe(pat string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.patterns, pat)
	return s.count()
}

// List all subscribed channels, in alphabetical order
func (s *Subscriber) Channels() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var cs []string
	for c := range s.channels {
		cs = append(cs, c)
	}
	sort.Strings(cs)
	return cs
}

// List all subscribed patterns, in alphabetical order
func (s *Subscriber) Patterns() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var ps []string
	for p := range s.patterns {
		ps = append(ps, p)
	}
	sort.Strings(ps)
	return ps
}

// Publish a message. Will return return how often we sent the message (can be
// a match for a subscription and for a psubscription.
func (s *Subscriber) Publish(c, msg string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	found := 0

subs:
	for sub := range s.channels {
		if sub == c {
			s.publish <- PubsubMessage{c, msg}
			found++
			break subs
		}
	}

pats:
	for _, pat := range s.patterns {
		if pat.MatchString(c) {
			s.publish <- PubsubMessage{c, msg}
			found++
			break pats
		}
	}

	return found
}

// The channel to read messages for this subscriber
func (s *Subscriber) Messages() <-chan PubsubMessage {
	return s.publish
}

// List all pubsub channels. If `pat` isn't empty channels names must match the
// pattern. Channels are returned alphabetically.
func activeChannels(subs []*Subscriber, pat string) []string {
	channels := map[string]struct{}{}
	for _, s := range subs {
		for c := range s.channels {
			channels[c] = struct{}{}
		}
	}

	var cpat *regexp.Regexp
	if pat != "" {
		cpat = compileChannelPattern(pat)
	}

	var cs []string
	for k := range channels {
		if cpat != nil && !cpat.MatchString(k) {
			continue
		}
		cs = append(cs, k)
	}
	sort.Strings(cs)
	return cs
}

// Count all subscribed (not psubscribed) clients for the given channel
// pattern. Channels are returned alphabetically.
func countSubs(subs []*Subscriber, channel string) int {
	n := 0
	for _, p := range subs {
		for c := range p.channels {
			if c == channel {
				n++
				break
			}
		}
	}
	return n
}

// Count the total of all client psubscriptions.
func countPsubs(subs []*Subscriber) int {
	n := 0
	for _, p := range subs {
		n += len(p.patterns)
	}
	return n
}

func monitorPublish(conn *server.Peer, msgs <-chan PubsubMessage) {
	for msg := range msgs {
		conn.Block(func(c *server.Writer) {
			c.WriteLen(3)
			c.WriteBulk("message")
			c.WriteBulk(msg.Channel)
			c.WriteBulk(msg.Message)
			c.Flush()
		})
	}
}
