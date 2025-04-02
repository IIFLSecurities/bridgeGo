package connector

type connectRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

type connectResponse struct {
	Message string `json:"message"`
	Status  int16  `json:"status"`
}

type subscribeRequest struct {
	SubscriptionList []string `json:"subscriptionList"`
}

type subscribeResponse struct {
	Status             int16                `json:"status"`
	Message            string               `json:"message"`
	SubscriptionResult []subscriptionResult `json:"subscriptionResult"`
}

type subscriptionResult struct {
	ResultCode int16  `json:"resultCode"`
	Result     string `json:"result"`
	Topic      string `json:"topic"`
}

type unsubscribeRequest struct {
	UnSubscriptionList []string `json:"UnsubscriptionList"`
}

type unsubscribeResponse struct {
	Status  int16  `json:"status"`
	Message string `json:"message"`
}

type disconnectResponse struct {
	Status  int16  `json:"status"`
	Message string `json:"message"`
}
