package connector

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	pack "github.com/eclipse/paho.mqtt.golang/packets"
)

const (
	mw           = "prod/marketfeed/mw/v1/"
	index        = "prod/marketfeed/index/v1/"
	oi           = "prod/marketfeed/oi/v1/"
	marketStatus = "prod/marketfeed/marketStatus/v1/"
	lpp          = "prod/marketfeed/lpp/v1/"
	high52Week   = "prod/marketfeed/high52week/v1/"
	low52Week    = "prod/marketfeed/low52week/v1/"
	upperCircuit = "prod/marketfeed/uppercircuit/v1/"
	lowerCircuit = "prod/marketfeed/lowercircuit/v1/"
	order        = "prod/updates/order/v1/"
	trade        = "prod/updates/trade/v1/"
)

var validateTokenUrl = "https://idaas.iiflsecurities.com/v1/access/check/token"

type Connect struct {
	client              mqtt.Client
	OnDisconnect        onDisconnectHnadler
	MWHandler           messageHandler
	IndexHandler        messageHandler
	OpenInterstHandler  messageHandler
	MarketStatusHandler messageHandler
	LppHandler          messageHandler
	High52WeekHandler   messageHandler
	Low52WeekHandler    messageHandler
	UpperCircuitHandler messageHandler
	LowerCircuitHandler messageHandler
	OrderUpdatesHandler messageHandler
	TradeUpdatesHandler messageHandler
}

var instance *Connect

var once sync.Once

type messageHandler func([]byte, string)
type onDisconnectHnadler func(error)

var subackReturnCodes = map[uint8]string{
	0x00: "Success",
	0x01: "Success",
	0x02: "Success",
	0x80: "Failure",
}

// GetInstance provides the single instance of the connect struct.
// Public properties of the `connect` instance that the user must set:
//
//   - `OnDisconnect`: A callback function that is triggered when the connection is lost.
//     The user must set this to handle disconnection events.
//     Example:
//     instance.OnDisconnect = func(err error) {
//     fmt.Println("Disconnected:", err)
//     }
//
//   - `MWHandler`: A callback function to handle MarketWatch data messages.
//     The user must set this to process incoming MarketWatch data.
//     Example:
//     instance.MWHandler = func(payload []byte, topic string) {
//     fmt.Printf("MarketWatch Data: %s, Topic: %s\n", payload, topic)
//     }
//
//   - `IndexHandler`: A callback function to handle Index data messages.
//     The user must set this to process incoming Index data.
//     Example:
//     instance.IndexHandler = func(payload []byte, topic string) {
//     fmt.Printf("Index Data: %s, Topic: %s\n", payload, topic)
//     }
//
//   - `OpenInterstHandler`: A callback function to handle Open Interest data messages.
//     The user must set this to process incoming Open Interest data.
//     Example:
//     instance.OpenInterstHandler = func(payload []byte, topic string) {
//     fmt.Printf("Open Interest Data: %s, Topic: %s\n", payload, topic)
//     }
//
//   - `MarketStatusHandler`: A callback function to handle Market Status data messages.
//     The user must set this to process incoming Market Status data.
//     Example:
//     instance.MarketStatusHandler = func(payload []byte, topic string) {
//     fmt.Printf("Market Status Data: %s, Topic: %s\n", payload, topic)
//     }
//
//   - `LppHandler`: A callback function to handle LPP data messages.
//     The user must set this to process incoming LPP data.
//     Example:
//     instance.LppHandler = func(payload []byte, topic string) {
//     fmt.Printf("LPP Data: %s, Topic: %s\n", payload, topic)
//     }
//
//   - `High52WeekHandler`: A callback function to handle High 52 Week data messages.
//     The user must set this to process incoming High 52 Week data.
//     Example:
//     instance.High52WeekHandler = func(payload []byte, topic string) {
//     fmt.Printf("High 52 Week Data: %s, Topic: %s\n", payload, topic)
//     }
//
//   - `Low52WeekHandler`: A callback function to handle Low 52 Week data messages.
//     The user must set this to process incoming Low 52 Week data.
//     Example:
//     instance.Low52WeekHandler = func(payload []byte, topic string) {
//     fmt.Printf("Low 52 Week Data: %s, Topic: %s\n", payload, topic)
//     }
//
//   - `UpperCircuitHandler`: A callback function to handle Upper Circuit data messages.
//     The user must set this to process incoming Upper Circuit data.
//     Example:
//     instance.UpperCircuitHandler = func(payload []byte, topic string) {
//     fmt.Printf("Upper Circuit Data: %s, Topic: %s\n", payload, topic)
//     }
//
//   - `LowerCircuitHandler`: A callback function to handle Lower Circuit data messages.
//     The user must set this to process incoming Lower Circuit data.
//     Example:
//     instance.LowerCircuitHandler = func(payload []byte, topic string) {
//     fmt.Printf("Lower Circuit Data: %s, Topic: %s\n", payload, topic)
//     }
//
//   - `OrderUpdatesHandler`: A callback function to handle Order Updates messages.
//     The user must set this to process incoming Order Updates data.
//     Example:
//     instance.OrderUpdatesHandler = func(payload []byte, topic string) {
//     fmt.Printf("Order Updates Data: %s, Topic: %s\n", payload, topic)
//     }
//
//   - `TradeUpdatesHandler`: A callback function to handle Trade Updates messages.
//     The user must set this to process incoming Trade Updates data.
//     Example:
//     instance.TradeUpdatesHandler = func(payload []byte, topic string) {
//     fmt.Printf("Trade Updates Data: %s, Topic: %s\n", payload, topic)
//     }
func GetInstance() *Connect {
	once.Do(func() {
		instance = &Connect{}
	})
	return instance
}

// ConnectHost is a method that connects to the broker
// Parameters:
// - connectReq: A JSON string containing the connection request details
// Returns:
// - A JSON string containing the connection response
// Sample request body:
//
//	{
//	    "host": "broker.example.com",
//	    "port": 8883,
//	    "password": "your_password"
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0
//	}
func (c *Connect) ConnectHost(connectReq string) (string, error) {
	var request connectRequest
	var response connectResponse

	if connectReq == "" {
		response.Message = "Connection Request should not be empty"
		response.Status = 101
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}

	err := json.Unmarshal([]byte(connectReq), &request)
	if err != nil {
		return "", err
	}

	if request.Host == "" {
		response.Message = "Parameter 'host' should not be empty"
		response.Status = 101
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil

	}
	if request.Port == 0 {
		response.Message = "Parameter 'port' should not be empty"
		response.Status = 101
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}
	if request.Password == "" {
		response.Message = "Parameter 'password' should not be empty"
		response.Status = 101
		jsonData, err := json.Marshal(response)
		if err == nil {
			return "", err
		}
		return string(jsonData), nil
	}

	userName := getUserName(request.Password)
	status, res, err := validateToken(userName, request.Password)

	if err != nil {
		return "", err
	} else {
		if status != 0 {
			response.Message = res
			response.Status = 1
			jsonData, err := json.Marshal(response)
			if err != nil {
				return "", err
			}
			return string(jsonData), nil
		}
	}
	currentTime := time.Now()
	formattedTime := currentTime.Format("020106150405")
	opts := mqtt.NewClientOptions().AddBroker("ssl://" + request.Host + ":" + strconv.Itoa(request.Port)).SetClientID(userName + "_go_" + formattedTime)
	opts.KeepAlive = 20
	opts.SetUsername(userName)
	opts.SetPassword("OPENID~~" + request.Password + "~")
	tlsConfig := newTlsConfig()
	opts.SetTLSConfig(tlsConfig)
	opts.SetProtocolVersion(4)
	opts.AutoReconnect = false
	opts.OnConnectionLost = onDisconnect
	if c.client == nil {
		c.client = mqtt.NewClient(opts)
	}

	if c.client.IsConnected() {
		response.Message = "Client is already connected"
		response.Status = 105
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}

	token := c.client.Connect()
	if token.Wait() && token.Error() != nil {
		connectToken := token.(*mqtt.ConnectToken)
		response.Message = pack.ConnackReturnCodes[uint8(connectToken.ReturnCode())]
		response.Status = int16(connectToken.ReturnCode())
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), token.Error()
	}

	connectToken := token.(*mqtt.ConnectToken)
	response.Message = pack.ConnackReturnCodes[uint8(connectToken.ReturnCode())]
	response.Status = int16(connectToken.ReturnCode())
	jsonData, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// SubscribeFeed subscribes to the MarketWatch data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeFeed(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, mw)
}

// SubscribeIndex subscribes to the Index data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeIndex(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, index)
}

// SubscribeOpenInterest subscribes to the Index data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeOpenInterest(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, oi)
}

// SubscribeMarketStatus subscribes to the MarketStatus data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeMarketStatus(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, marketStatus)
}

// SubscribeLpp subscribes to the Lpp data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeLpp(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, lpp)
}

// SubscribeHigh52Week subscribes to the High52Week data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeHigh52Week(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, high52Week)
}

// SubscribeLow52Week subscribes to the Low52Week data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeLow52Week(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, low52Week)
}

// SubscribeUpperCircuit subscribes to the UpperCircuit data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeUpperCircuit(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, upperCircuit)
}

// SubscribeLowerCircuit subscribes to the LowerCircuit data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeLowerCircuit(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, lowerCircuit)
}

// SubscribeOrderUpdates subscribes to the OrderUpdates data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeOrderUpdates(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, order)
}

// SubscribeTradeUpdates subscribes to the TradeUpdates data
// Parameters:
// - subscribeReq: A JSON string containing the subscription request details
// Returns:
// - A JSON string containing the subscription response
//
// Sample request body:
//
//	{
//	    "SubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Success",
//	    "Status": 0,
//	    "SubscriptionResult": [
//	        {
//	            "ResultCode": 0,
//	            "Result": "Subscribed successfully",
//	            "Topic": "nseeq/2885"
//	        }
//	    ]
//	}
func (c *Connect) SubscribeTradeUpdates(subscribeReq string) (string, error) {
	return c.subscribe(subscribeReq, trade)
}

func (c *Connect) subscribe(subscribeReq string, feedType string) (string, error) {

	var response subscribeResponse
	if c.client == nil || c.client.IsConnected() == false {
		response.Message = "Client not connected"
		response.Status = 106
		response.SubscriptionResult = nil
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}

	if subscribeReq == "" {
		response.Message = "Subscription Request should not be empty"
		response.Status = 101
		response.SubscriptionResult = nil
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}
	var request subscribeRequest

	err := json.Unmarshal([]byte(subscribeReq), &request)

	if request.SubscriptionList == nil || len(request.SubscriptionList) > 1024 {
		response.Message = "TopicList cannot be nil and no. of topics should be less than 1024"
		response.Status = 102
		response.SubscriptionResult = nil
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}

	filter := map[string]byte{}
	for _, value := range request.SubscriptionList {

		match, err := regexp.MatchString(`^[A-Za-z0-9\/]+$`, value)
		if err != nil {
			return "", err
		}
		if match {
			filter[feedType+value] = 0
		} else {
			response.SubscriptionResult = append(response.SubscriptionResult, subscriptionResult{ResultCode: 104, Result: "Invalid Topic", Topic: value})
		}
	}
	if len(filter) != 0 {
		token := c.client.SubscribeMultiple(filter, messagehandler)
		if token.Wait() && token.Error() != nil {
			response.Status = -1
			response.Message = token.Error().Error()
			jsonData, err := json.Marshal(response)
			if err != nil {
				return "", nil
			}
			return string(jsonData), token.Error()
		}
		subscribeToken := token.(*mqtt.SubscribeToken)
		re := regexp.MustCompile("(v1/)")
		for key, value := range subscribeToken.Result() {
			response.SubscriptionResult = append(response.SubscriptionResult, subscriptionResult{ResultCode: int16(value), Result: subackReturnCodes[uint8(value)], Topic: re.Split(key, -1)[1]})
		}

		response.Message = "Success"
		response.Status = 0
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}

	response.Message = "Subscription Failed"
	response.Status = 103
	jsonData, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// UnsubscribeFeed unsubscribes from the MarketWatch data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeFeed(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, mw)
}

// UnsubscribeIndex unsubscribes from the Index data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeIndex(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, index)
}

// UnsubscribeOpenInterest unsubscribes from the OpenInterest data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeOpenInterest(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, oi)
}

// UnsubscribeMarketStatus unsubscribes from the MarketStatus data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeMarketStatus(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, marketStatus)
}

// UnsubscribeLpp unsubscribes from the Lpp data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeLpp(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, lpp)
}

// UnsubscribeHigh52Week unsubscribes from the High52Week data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeHigh52Week(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, high52Week)
}

// UnsubscribeLow52Week unsubscribes from the Low52Week data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeLow52Week(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, low52Week)
}

// UnsubscribeUpperCircuit unsubscribes from the UpperCircuit data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeUpperCircuit(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, upperCircuit)
}

// UnsubscribeLowerCircuit unsubscribes from the LowerCircuit data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeLowerCircuit(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, lowerCircuit)
}

// UnsubscribeOrderUpdates unsubscribes from the OrderUpdates data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeOrderUpdates(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, order)
}

// UnsubscribeTradeUpdates unsubscribes from the TradeUpdates data
// Parameters:
// - unsubscribeReq: A JSON string containing the unsubscription request details
// Returns:
// - A JSON string containing the unsubscription response
//
// Sample request body:
//
//	{
//	    "UnsubscriptionList": ["nseeq/2885", "nsee/2886", "bseeq/2887"]
//	}
//
// Sample response:
//
//	{
//	    "Message": "Unsubscribed Successfully",
//	    "Status": 0
//	}
func (c *Connect) UnsubscribeTradeUpdates(unsubscribeReq string) (string, error) {
	return c.unsubscribe(unsubscribeReq, trade)
}

func (c *Connect) unsubscribe(unsubscribeReq string, feedType string) (string, error) {

	var response unsubscribeResponse
	if c.client == nil || c.client.IsConnected() == false {
		response.Message = "Client not connected"
		response.Status = 106
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}

	var request unsubscribeRequest

	err := json.Unmarshal([]byte(unsubscribeReq), &request)

	if request.UnSubscriptionList == nil || len(request.UnSubscriptionList) > 1024 {
		response.Message = "TopicList cannot be nil and no. of topics should be less than 1024"
		response.Status = 102
		// response.UnSubscriptionResult = nil
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}

	var filter []string
	for _, value := range request.UnSubscriptionList {
		match, err := regexp.MatchString(`^[A-Za-z0-9\/]+$`, value)
		if err != nil {
			return "", err
		}
		if match {
			filter = append(filter, feedType+value)
		}
	}

	if len(filter) != 0 {
		token := c.client.Unsubscribe(filter...)
		if token.Wait() && token.Error() != nil {
			response.Message = token.Error().Error()
			response.Status = -1
			jsonData, err := json.Marshal(response)
			if err != nil {
				return "", err
			}
			return string(jsonData), token.Error()

		} else {
			response.Message = "Unsubscribed Successfully"
			response.Status = 0
			jsonData, err := json.Marshal(response)
			if err != nil {
				return "", err
			}
			return string(jsonData), nil
		}
	}

	response.Message = "Unsubscription Failed"
	response.Status = 103
	jsonData, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// DisconnectHost disconnects from the broker
// Returns:
// - A JSON string containing the disconnection response
//
// Sample response:
//
//	{
//	    "Message": "Disconnected Successfully",
//	    "Status": 0
//	}
func (c *Connect) DisconnectHost() (string, error) {

	var response disconnectResponse
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
	} else {
		response.Message = "Client not connected"
		response.Status = 106
		jsonData, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}

	if c.client.IsConnected() {
		response.Message = "Disconnection Failed"
		response.Status = 1
	} else {
		response.Message = "Disconnected Successfully"
		response.Status = 0
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	c.client = nil
	return string(jsonData), nil
}

// IsConnected checks if the client is currently connected to the broker.
// Returns:
// - true: If the client is connected.
// - false: If the client is not connected or the client instance is nil.
func (c *Connect) IsConnected() bool {
	if c.client != nil {
		return c.client.IsConnected()
	}
	return false
}

func getUserName(tokenString string) string {

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "bridgeGo"
	}

	payload := parts[1]
	// Calculate padding
	padding := (4 - len(payload)%4) % 4
	payload += strings.Repeat("=", padding)

	// Decode the payload (second part)
	decodedpayload, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return "bridgeGo"
	}

	var payloadMap map[string]interface{}
	if err := json.Unmarshal(decodedpayload, &payloadMap); err != nil {
		return "bridgeGo"
	}
	// Retrieve the claim value by its key (e.g., "user_id")
	claimValue, ok := payloadMap["preferred_username"].(string)
	if !ok {
		return "bridgeGo"
	}
	return claimValue
}

func newTlsConfig() *tls.Config {
	// Create a TLS configuration with InsecureSkipVerify set to true
	return &tls.Config{
		InsecureSkipVerify: true, // Skip server certificate verification
	}
}

func validateToken(userName string, token string) (int, string, error) {

	// Define the data to send in the POST request
	data := map[string]interface{}{
		"userId": userName,
		"token":  token,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return -1, "", err
	}

	response, err := http.Post(validateTokenUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return -1, "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, "", err
	}

	var responseMap map[string]interface{}

	er := json.Unmarshal(body, &responseMap)
	if er != nil {
		return -1, "", er
	}

	result := responseMap["result"].(map[string]interface{})
	if result["status"].(string) != "Success" {
		return -1, result["message"].(string), nil
	}

	return 0, "Success", nil
}

var messagehandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	re := regexp.MustCompile("(v1/)")
	parts := re.Split(topic, -1)
	payload := msg.Payload()

	switch parts[0] + "v1/" {

	case mw:
		if instance.MWHandler != nil {
			instance.MWHandler(payload, parts[1])
		}
	case index:
		if instance.IndexHandler != nil {
			instance.IndexHandler(payload, parts[1])
		}
	case oi:
		if instance.OpenInterstHandler != nil {
			instance.OpenInterstHandler(payload, parts[1])
		}
	case marketStatus:
		if instance.MarketStatusHandler != nil {
			instance.MarketStatusHandler(payload, parts[1])
		}
	case lpp:
		if instance.LppHandler != nil {
			instance.LppHandler(payload, parts[1])
		}
	case high52Week:
		if instance.High52WeekHandler != nil {
			instance.High52WeekHandler(payload, parts[1])
		}
	case low52Week:
		if instance.Low52WeekHandler != nil {
			instance.Low52WeekHandler(payload, parts[1])
		}
	case upperCircuit:
		if instance.UpperCircuitHandler != nil {
			instance.UpperCircuitHandler(payload, parts[1])
		}
	case lowerCircuit:
		if instance.LowerCircuitHandler != nil {
			instance.LowerCircuitHandler(payload, parts[1])
		}
	case order:
		if instance.OrderUpdatesHandler != nil {
			instance.OrderUpdatesHandler(payload, parts[1])
		}
	case trade:
		if instance.TradeUpdatesHandler != nil {
			instance.TradeUpdatesHandler(payload, parts[1])
		}
	default:
	}
}

var onDisconnect mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	if instance.OnDisconnect != nil {
		instance.client = nil
		instance.OnDisconnect(err)
	}
}
