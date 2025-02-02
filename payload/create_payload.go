package payload

var payloadMap = make(map[payloadType]func() Payload)

// CreatePayload is factory method to create different payloads based on type
func CreatePayload(data *[]byte, from string) (Payload, error) {
	// base payload to get type from data
	var base basePayload
	if !unmarshalAndValidate(data, &base) {
		return nil, &InvalidPayload{
			reason: "Invalid payload received.",
		}
	}

	// basic routing check to see if sender is user only
	if base.From != from {
		return nil, &InvalidPayload{
			reason: "Client username and payload from mismatch.",
		}
	}

	// get factory method to generate the  actual payload based on type
	factory, exists := payloadMap[base.Type]
	if !exists {
		return nil, &InvalidPayload{
			reason: "Invalid payload received",
		}
	}

	// validate the payload and reject if not proper
	payload := factory()
	if unmarshalAndValidate(data, &payload) {
		return nil, &InvalidPayload{
			reason: "Invalid payload received",
		}
	}

	return payload, nil
}