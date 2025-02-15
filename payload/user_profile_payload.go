package payload

const (
	userUpdateProfileType = payloadType("user_update_profile")

	// root nodes are nodes that are only dependent on only one node
	// post, discussion they only depend on only user node
	// but comment is dependent on 2 nodes, created by(user) and comment on(post or comment)
	userCreateRootNodeType      = payloadType("user_create_root_node")
	userNodeLikeActionType      = payloadType("user_node_like_action")
	userCreateSecondaryNodeType = payloadType("user_create_secondary_node")
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

type parentNode struct {
	NodeId   string `json:"nodeId" validate:"required"`
	NodeType string `json:"nodeType" validate:"required"`
}

type userNodeLikeAction struct {
	Type         payloadType  `json:"type" validate:"required"`
	From         string       `json:"from" validate:"required"`
	To           string       `json:"to" validate:"required"`
	IsLike       bool         `json:"isLike"`
	LikeCount    int          `json:"likeCount"`
	CommentCount int          `json:"commentCount"`
	NodeId       string       `json:"nodeId" validate:"required"`
	NodeType     string       `json:"nodeType" validate:"required"`
	Parents      []parentNode `json:"parents,string" validate:"required"`
}

func (payload *userNodeLikeAction) SendPayload(data *[]byte, h hub, senderResource string) {
	nodeOwner := payload.To
	actionBy := payload.From

	// this prevents sending messages twice when user interacts with self nodes
	if nodeOwner != actionBy {
		nodeOwnerConnectedClients := h.GetAllConnectedClients(nodeOwner)
		for _, conn := range nodeOwnerConnectedClients {
			conn.WriteToChannel(data)
		}
	}

	actionByConnectedClients := h.GetAllConnectedClients(actionBy)
	for res, conn := range actionByConnectedClients {
		if res != senderResource {
			conn.WriteToChannel(data)
		}
	}
}

type userCreateSecondaryNode struct {
	Type                 payloadType  `json:"type" validate:"required"`
	From                 string       `json:"from" validate:"required"`
	To                   string       `json:"to" validate:"required"`
	NodeId               string       `json:"nodeId" validate:"required"`
	NodeType             string       `json:"nodeType" validate:"required"`
	Mentions             []string     `json:"mentions"`
	ReplyOnNodeCreatedBy string       `json:"replyOnNodeCreatedBy"`
	Parents              []parentNode `json:"parents,string" validate:"required"`
}

func (payload *userCreateSecondaryNode) SendPayload(data *[]byte, h hub, senderResource string) {
	nodeCreator := payload.From
	parentNodeCreator := payload.To

	if nodeCreator != parentNodeCreator {
		parentNodeCreatorConnectedClients := h.GetAllConnectedClients(parentNodeCreator)
		for _, conn := range parentNodeCreatorConnectedClients {
			conn.WriteToChannel(data)
		}
	}

	nodeCreatorConnectedClients := h.GetAllConnectedClients(nodeCreator)
	for _, conn := range nodeCreatorConnectedClients {
		conn.WriteToChannel(data)
	}

	for _, userMentioned := range payload.Mentions {
		if userMentioned == nodeCreator || userMentioned == parentNodeCreator {
			continue
		}

		userMentionedAllConnectedClients := h.GetAllConnectedClients(userMentioned)
		for _, conn := range userMentionedAllConnectedClients {
			conn.WriteToChannel(data)
		}
	}

	if payload.ReplyOnNodeCreatedBy != nodeCreator && payload.ReplyOnNodeCreatedBy != parentNodeCreator {
		replyOnNodeCreatedByConnectedClients := h.GetAllConnectedClients(payload.ReplyOnNodeCreatedBy)
		for _, conn := range replyOnNodeCreatedByConnectedClients {
			conn.WriteToChannel(data)
		}
	}
}