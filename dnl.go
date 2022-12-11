package dnl

import (
	"errors"
	"fmt"
)

type Provider interface {
	Send(id string, msg string) error
	SetOnMessageFunc(onMessageFunc func(string, string) error)
}

type DNL interface {
	AddChannel(id string)
	RemoveChannel(id string)
	GetChannel(id string) chan string
	Send(id string, msg string) error
}

type dnl struct {
	channels map[string]chan string
	provider Provider
}

func NewWithProvider(provider Provider) DNL {
	d := &dnl{
		channels: make(map[string]chan string),
		provider: provider,
	}

	provider.SetOnMessageFunc(d.onMsgFunc)

	return d
}

func (d *dnl) onMsgFunc(id, msg string) error {
	// Check if channel exists
	if _, ok := d.channels[id]; ok {
		d.channels[id] <- msg

		return nil
	}

	return nil
}

func (d *dnl) SetProvider(provider Provider) {
	d.provider = provider
}

// AddChannel Add a new channel to the DWL
func (d *dnl) AddChannel(id string) {
	d.channels[id] = make(chan string)
}

// RemoveChannel Remove a channel from the DWL
func (d *dnl) RemoveChannel(id string) {
	delete(d.channels, id)
}

// GetChannel Get a channel from the DWL
func (d *dnl) GetChannel(id string) chan string {
	return d.channels[id]
}

// Send a message to a channel
func (d *dnl) Send(id string, msg string) error {
	// Check if channel exists
	if _, ok := d.channels[id]; ok {
		d.channels[id] <- msg

		return nil
	}

	// If channel doesn't exist, then we need to send a pub/sub message to the provider
	if d.provider != nil {
		err := d.provider.Send(id, msg)
		if err != nil {
			return fmt.Errorf("error sending message to provider: %s", err)
		}

		return nil
	}

	// Return an error if channel doesn't exist
	return errors.New("the provided id is not associated with a channel. Please add the channel first or provide a provider")
}

func (d *dnl) Close() {
	for _, channel := range d.channels {
		close(channel)
	}
}
