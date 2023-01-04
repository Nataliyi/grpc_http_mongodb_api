package services

import (
	watcher_models "api/models/watcher"
	"log"

	"cloud.google.com/go/pubsub"
)

type PubSubWatcher struct {
	ProjectID string
	Topic     string
	// pub-sub client
	Client        *pubsub.Client
	InformChannel watcher_models.WatcherChannel
}

// Init new PubSubWatcher
func NewPubSubWatcher(project string,
	topic string, credsPath string) (*PubSubWatcher, error) {

	// ctx := context.Background()
	// Client, err := pubsub.NewClient(ctx, project, opts)
	// if err != nil {
	// 	return nil, fmt.Errorf("pubsub: NewClient error: %v", err)
	// }
	w := &PubSubWatcher{
		ProjectID: project,
		Topic:     topic,
		// Client:          Client,
		InformChannel: make(watcher_models.WatcherChannel),
	}
	return w, nil
}

// Close the InformChannel and the Client
func (w *PubSubWatcher) Close() {
	close(w.InformChannel)
	// w.Client.Close()
}

// Start listen messages.
// The func should run like a goroutine.
func (w *PubSubWatcher) Listen() {

	// goroutine for receiving
	for mess := range w.InformChannel {
		// send the message to the pub sub topic
		w.send(mess)
	}
}

// Return the active receiving channel
func (w *PubSubWatcher) GetChannel() watcher_models.WatcherChannel {
	return w.InformChannel
}

// Send the message to the pub-sub topic
func (w *PubSubWatcher) send(message string) error {
	// could implement a pub-sub method here
	log.Printf("PubSubWatcher: Received the PUB SUB Message: %s", message)
	return nil
}
