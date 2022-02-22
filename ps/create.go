package ps

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

func CreateTopic(args PubSubArguments) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, args.Project)
	if err != nil {
		log.Printf("pubsub.NewClient: %v", err)
		return err
	}
	defer client.Close()
	t, err := client.CreateTopic(ctx, args.Topic)
	defer t.Stop()
	if err != nil {
		log.Printf("pubsub.CreateTopic: %v", err)
		return err
	}
	log.Println("created topic:", t.ID())
	return nil
}

func CreateSubscription(args PubSubArguments) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, args.Project)
	if err != nil {
		log.Printf("pubsub.NewClient: %v", err)
		return err
	}
	defer client.Close()
	t := client.Topic(args.Topic)
	defer t.Stop()
	s, err := client.CreateSubscription(ctx, args.Subscription, pubsub.SubscriptionConfig{
		Topic: t,
	})
	if err != nil {
		log.Printf("pubsub.CreateSubscription: %v", err)
		return err
	}
	log.Println("created subscription:", s.ID())
	return nil
}
