package ps

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
)

// Pull messages from a Subscription and Push them to an endpoint.
func PullPush(args PubSubArguments) error {
	ctx := context.Background()
	pullClient, err := pubsub.NewClient(ctx, args.Project)
	if err != nil {
		log.Printf("pubsub.NewClient: %v", err)
		return err
	}
	defer pullClient.Close()
	sub := pullClient.Subscription(args.Subscription)
	pushClient := new(http.Client)
	for {
		log.Println("receiving from", args.Subscription)
		err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			log.Printf("received message: %s\n", msg.ID)
			msgOut := PubSubMessage{}
			msgOut.Message.Data = msg.Data
			msgOut.Message.ID = msg.ID
			msgOut.Subscription = args.Subscription
			msgOut.MessageId = msg.ID
			msgOut.Message.Attributes = msg.Attributes
			dataOut, err := json.Marshal(msgOut)
			if err != nil {
				log.Printf("payload marshal failed: %v", err)
				return
			}
			req, err := http.NewRequest(http.MethodPost, args.PushURL, bytes.NewReader(dataOut))
			if err != nil {
				log.Printf("create request failed: %v", err)
				return
			}
			resp, err := pushClient.Do(req)
			if err != nil {
				log.Printf("send request failed: %v", err)
				return
			}
			if resp.StatusCode == http.StatusOK {
				log.Printf("pushed message: %s\n", msg.ID)
				msg.Ack()
			} else {
				log.Printf("failed to push message (Nack-ed it): %s error: %v\n", msg.ID, resp.Status)
				msg.Nack()
			}
		})
		if err != nil {
			return err
		}
	}
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Message struct {
		Data       []byte            `json:"data,omitempty"`
		Attributes map[string]string `json:"attributes"`
		ID         string            `json:"id"`
	} `json:"message"`
	MessageId    string `json:"messageId"`
	Subscription string `json:"subscription"`
}
