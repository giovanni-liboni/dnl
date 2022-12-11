package dnl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

type ProviderRedis struct {
	client        *redis.Client
	onMessageFunc func(string, string) error
}

const (
	redisChannel = "dnl"
)

type ChannelPayload struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// NewProviderRedis Pub/sub client for Redis
func NewProviderRedis(addr string) *ProviderRedis {
	pr := &ProviderRedis{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}

	go pr.startSubscriptions()

	return pr
}

func (p *ProviderRedis) SetOnMessageFunc(onMessageFunc func(string, string) error) {
	p.onMessageFunc = onMessageFunc
}

func (p *ProviderRedis) Send(id string, msg string) error {
	if p.client == nil {
		return errors.New("no redis client")
	}

	payload := ChannelPayload{
		ID:   id,
		Data: msg,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = p.client.Publish(redisChannel, payloadBytes).Err()
	if err != nil {
		return fmt.Errorf("error sending message to redis: %s", err)
	}

	return nil
}

func (p *ProviderRedis) startSubscriptions() {
	pubsub := p.client.Subscribe(redisChannel)
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {

		}
	}(pubsub)

	ch := pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			// Unmarshal the message
			var payload ChannelPayload
			err := json.Unmarshal([]byte(msg.Payload), &payload)
			if err != nil {
				fmt.Println("error unmarshalling message:", err)
				continue
			}

			err = p.onMessageFunc(payload.ID, payload.Data)
			if err != nil {
				fmt.Println("error handling message:", err)
			}
		}
	}
}
