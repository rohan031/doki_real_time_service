package payload

// CreatePayload is factory method to create different payloads based on type
func CreatePayload(data *[]byte, from string) (Payload, error) {
	var base basePayload
	if !unmarshalAndValidate(data, &base) {
		return nil, &InvalidPayload{
			reason: "Invalid payload received.",
		}
	}

	if base.From != from {
		return nil, &InvalidPayload{
			reason: "Client username and payload from mismatch.",
		}
	}

	switch base.Type {
	case chatMessageType:
		var message chatMessage
		if unmarshalAndValidate(data, &message) {
			return &message, nil
		}

	case typingStatusType:
		var status typingStatus
		if unmarshalAndValidate(data, &status) {
			return &status, nil
		}
	case editMessageType:
		var message editMessage
		if unmarshalAndValidate(data, &message) {
			return &message, nil
		}
	case deleteMessageType:
		var message deleteMessage
		if unmarshalAndValidate(data, &message) {
			return &message, nil
		}
	default:
		return nil, &InvalidPayload{
			reason: "unknown payload type",
		}
	}

	return nil, &InvalidPayload{
		reason: "Invalid payload received",
	}
}