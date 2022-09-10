package ps

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
)

func Publish(args PubSubArguments) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, args.Project)
	if err != nil {
		log.Printf("pubsub.NewClient: %v", err)
		return err
	}
	defer client.Close()
	fmt.Println("reading from", args.File)
	data, err := ioutil.ReadFile(args.File)
	if err != nil {
		log.Printf("reading file: %v", err)
	}
	t := client.Topic(args.Topic)
	fmt.Println("publishing to", args.Topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"origin": "gcloudx",
		},
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		log.Printf("Get: %v", err)
	}
	fmt.Printf("Published message with custom attributes; msg ID: %v\n", id)
	return nil
}

func Pull(args PubSubArguments) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, args.Project)
	if err != nil {
		log.Printf("pubsub.NewClient: %v", err)
		return err
	}
	defer client.Close()
	sub := client.Subscription(args.Subscription)
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	var data []byte
	fmt.Println("receiving from", args.Subscription)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		data = msg.Data
		msg.Ack()
		cancel()
	})
	if err != nil {
		log.Printf("Receive: %v", err)
	}
	if len(data) == 0 {
		return err
	}
	if len(args.File) > 0 {
		fmt.Println("writing to", args.File)
		err = os.WriteFile(args.File, data, os.ModePerm)
		if err != nil {
			log.Printf("Write: %v", err)
		}
	}
	return nil
}
