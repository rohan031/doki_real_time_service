package payload

const (
	userUpdateProfileType = payloadType("user_update_profile")

	// root nodes are nodes that are only dependent on only one node
	// post, discussion they only depend on only user node
	// but comment is dependent on 2 nodes, created by(user) and comment on(post or comment)
	userCreateRootNodeType = payloadType("user_create_root_node")
)

type userUpdateProfile struct {
	Type           payloadType `json:"type" validate:"required"`
	From           string      `json:"from" validate:"required"`
	Name           string      `json:"name" validate:"required"`
	ProfilePicture string      `json:"profilePicture"`
	Bio            string      `json:"bio"`
}

func (payload *userUpdateProfile) SendPayload(data *[]byte, h hub, senderResource string) {
	user := payload.From

	userConnectedClients := h.GetAllConnectedClients(user)
	for res, conn := range userConnectedClients {
		if res != senderResource {
			conn.WriteToChannel(data)
		}
	}
}

type userCreateRootNode struct {
	Type     payloadType `json:"type" validate:"required"`
	From     string      `json:"from" validate:"required"`
	Id       string      `json:"id" validate:"required"`
	NodeType string      `json:"nodeType" validate:"required"`
}

func (payload *userCreateRootNode) SendPayload(data *[]byte, h hub, senderResource string) {
	user := payload.From

	userConnectedClients := h.GetAllConnectedClients(user)
	for res, conn := range userConnectedClients {
		if res != senderResource {
			conn.WriteToChannel(data)
		}
	}
}