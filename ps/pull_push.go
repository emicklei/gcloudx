package ps

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/cel-go/cel"
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
	sub.ReceiveSettings = pubsub.ReceiveSettings{
		MaxExtensionPeriod: 10 * time.Second,
	}
	pushClient := new(http.Client)

	if args.SubscriptionFilter != "" {
		env, _ := cel.NewEnv(
			cel.Variable("attributes", cel.MapType(cel.StringType, cel.StringType)),
		)
		log.Println(env)
	}

	log.Println("waiting for messages from subscription:", args.Subscription, "...")
	if err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("received message: %s subscription: %s \n", msg.ID, args.Subscription)
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
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, args.PushURL, bytes.NewReader(dataOut))
		if err != nil {
			log.Printf("create request failed: %v", err)
			return
		}
		resp, err := pushClient.Do(req)
		if err != nil {
			log.Printf("send POST request failed: %v", err)
			return
		}
		if resp.StatusCode == http.StatusOK {
			log.Printf("pushed message: %s attributes: %v\n", msg.ID, msg.Attributes)
			msg.Ack()
		} else {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("failed to push message: %s error: %v body:%s\n", msg.ID, resp.Status, string(body))
			if args.AbortOnError {
				log.Fatal(err)
			}
			if args.AlwaysACK {
				log.Printf("despite the error, the message is acknowledged, id:%s", msg.ID)
				msg.Ack()
			}
		}
	}); err != nil {
		log.Println("Receive err:", err)
	}
	return err
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
