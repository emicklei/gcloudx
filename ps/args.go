package ps

type PubSubArguments struct {
	Project      string
	File         string
	Topic        string
	Subscription string
	UseEmulator  bool
	PushURL      string
}