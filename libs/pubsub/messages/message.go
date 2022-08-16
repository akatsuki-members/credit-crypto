package messages

// Header contains event metadata.
type Header struct {
	ID          string // id correlation id.
	Domain      string // domain the business domain of the message.
	EventType   string // EventType within the domain what type of event it is.
	Version     string // Version it is the event type version.
	Application string // AppName name of the sender application
	MessageID   string // MessageID id used for message acknowledge.
}

// Event contains data related to the event.
type Event struct {
	Header Header
	Data   []byte
}
