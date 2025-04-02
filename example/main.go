//Use below command to install the package

//go get github.com/IIFLSecurities/bridgeGo@v0.1.1

//Use below command to resolve the dependencies

// go mod tidy

// Go through the below implementation to understand how to use the package
// This is a sample implementation of the package
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
	"unsafe"

	connector "github.com/IIFLSecurities/bridgeGo/connector"
)

// bridge "bridge/connector"

type MWBOCombined struct {
	Ltp                int32     `json:"ltp"`
	LastTradedQuantity uint32    `json:"lastTradedQuantity"`
	TradedVolume       uint32    `json:"tradedVolume"`
	High               int32     `json:"high"`
	Low                int32     `json:"low"`
	Open               int32     `json:"open"`
	Close              int32     `json:"close"`
	AverageTradedPrice int32     `json:"averageTradedPrice"`
	Reserved           uint16    `json:"reserved"`
	BestBidQuantity    uint32    `json:"bestBidQuantity"`
	BestBidPrice       int32     `json:"bestBidPrice"`
	BestAskQuantity    uint32    `json:"bestAskQuantity"`
	BestAskPrice       int32     `json:"bestAskPrice"`
	TotalBidQuantity   uint32    `json:"totalBidQuantity"`
	TotalAskQuantity   uint32    `json:"totalAskQuantity"`
	PriceDivisor       int32     `json:"priceDivisor"`
	LastTradedTime     int32     `json:"lastTradedTime"`
	MarketDepth        [10]Depth `json:"marketDepth"` // Array of 10 Depth structures
}

type Depth struct {
	Quantity        uint32 `json:"quantity"`
	Price           int32  `json:"price"`
	Orders          int16  `json:"orders"`
	TransactionType int16  `json:"transactionType"`
}

func main() {

	connector := connector.GetInstance()

	connector.MWHandler = func(b []byte, s string) {
		// fmt.Println("MWHandler", s, string(b))
		// Create a buffer from the byte array
		buffer := bytes.NewBuffer(b)

		// Create an instance of MWBOCombined
		var data MWBOCombined
		fmt.Printf("Size of Person struct: %d\n", unsafe.Sizeof(data))

		// Get the alignment of the struct
		fmt.Printf("Alignment of Person struct: %d\n", unsafe.Alignof(data))
		// Decode the byte array into the struct
		err := binary.Read(buffer, binary.LittleEndian, &data)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			return
		}
		fmt.Printf("%+v\n", data)
	}

	connector.IndexHandler = func(b []byte, s string) {
		fmt.Println("IndexHandler", s, string(b))
	}

	connector.Low52WeekHandler = func(b []byte, s string) {
		fmt.Println("Low52WeekHandler", s, string(b))
	}

	connector.High52WeekHandler = func(b []byte, s string) {
		fmt.Println("High52WeekHandler", s, string(b))
	}

	connector.OpenInterstHandler = func(b []byte, s string) {
		fmt.Println("OpenInterstHandler", s, string(b))
	}

	connector.LppHandler = func(b []byte, s string) {
		fmt.Println("LppHandler", s, string(b))
	}

	connector.UpperCircuitHandler = func(b []byte, s string) {
		fmt.Println("UpperCircuitHandler", s, string(b))
	}

	connector.LowerCircuitHandler = func(b []byte, s string) {
		fmt.Println("LowerCircuitHandler", s, string(b))
	}

	connector.MarketStatusHandler = func(b []byte, s string) {
		fmt.Println("MarketStatusHandler", s, string(b))
	}

	connector.OrderUpdatesHandler = func(b []byte, s string) {
		fmt.Println("OrderHandler", s, string(b))
	}

	connector.TradeUpdatesHandler = func(b []byte, s string) {
		fmt.Println("TradeHandler", s, string(b))
	}

	data := map[string]interface{}{
		"host":     "bridge.iiflcapital.com",
		"port":     9906,
		"password": "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIxVks4TEhlRnRvSmp6YWk1RmJlSGNPbDI3ekpGanBScTE2Vmt4eGJBZ0ZjIn0.eyJleHAiOjE3NTkxMjQ2OTgsImlhdCI6MTc0MzU3MjY5OCwianRpIjoiN2YzOGI4ZGUtZDk2YS00ZjU4LTlhNGMtNWNiMjk2OTk3NjU3IiwiaXNzIjoiaHR0cHM6Ly9rZXljbG9hay5paWZsc2VjdXJpdGllcy5jb20vcmVhbG1zL0lJRkwiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNDk4NTU2NTAtMDcyNS00ZGUzLWFlMmQtN2ExMGYxYzIwODI4IiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiSUlGTCIsInNpZCI6IjkzMjY3ZGY2LThhYzQtNGE0NC1hNzU4LTUzODY4YWYyNzllMyIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovLzEwLjEyNS42OC4xNDQ6ODA4MC8iXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbImRlZmF1bHQtcm9sZXMtaWlmbCIsIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJJSUZMIjp7InJvbGVzIjpbIkdVRVNUX1VTRVIiLCJBQ1RJVkVfVVNFUiJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJ1Y2MiOiI5MzA4MDA0OCIsInNvbGFjZV9ncm91cCI6IlNVQlNDUklCRVJfQ0xJRU5UIiwibmFtZSI6IlNBVEhFRVNIIEtVTUFSIEdPUElTRVRUWSBOQSIsInByZWZlcnJlZF91c2VybmFtZSI6IjkzMDgwMDQ4IiwiZ2l2ZW5fbmFtZSI6IlNBVEhFRVNIIEtVTUFSIEdPUElTRVRUWSIsImZhbWlseV9uYW1lIjoiTkEiLCJlbWFpbCI6InNhdGlzaGt1bWFyZ29waXNldHR5MTExQGdtYWlsLmNvbSJ9.4753i1DPabO9sa0zSj2RuDLSXFbdQBU-otX5F767ffOAMkrs5qnqcF30QkcVNRcm_FuQhWSVTjoNOk0rXS80ZXpSJfgBUIUsp1OOAefXnaGAyVPuapR7-mnqk6SbGj6UpXW8TS8um_N4b5dJ8GCel_JQV0ebc5B4y0Nxa8xWQ4Ee8rkjXtYy15LVmWaR8LmtKXp8zaw5O2yNLn7z6vOMjsD4tLLnJ-6UI-Wb18bRlPXIQTHVukAErlSaNH9Eu0FQ5fuugY8kuIgduliomHkegWG67rVbasuHYitThhv2ZmQrE1EWIHZU5cUGCPLHER7VhNacvpw4RI3Jv598d9Itnw"}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	//Connect to Host
	res, err := connector.ConnectHost(string(jsonData))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(res)
	} else {
		fmt.Println(res)
	}

	//	Subscribe Feed
	subReq := map[string]interface{}{
		"SubscriptionList": []string{"nseeq/2885", "nsee/2886", "bseeq/2887"},
	}
	jsonData, err = json.Marshal(subReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	subRes, err := connector.SubscribeFeed(string(jsonData))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(subRes)
	}

	time.Sleep(1 * time.Second)
	//Unsubscribe Feed
	UnsubReq := map[string]interface{}{
		"UnsubscriptionList": []string{"nseeq/2885", "nsee/2886", "bseeq/2887"},
	}
	jsonData, err = json.Marshal(UnsubReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	unsubRes, err := connector.UnsubscribeFeed(string(jsonData))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(unsubRes)
	}

	//Disconnet
	disRes, err := connector.DisconnectHost()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(disRes)
	}
	// bridge.Subscribe("prod/marketfeed/mw/v1/nseeq/2885")
	// time.Sleep(20 * time.Second)
	// bridge.Unsubscribe("prod/marketfeed/mw/v1/nseeq/2885")
	// time.Sleep(20 * time.Second)
	// bridge.Disconnect()
	time.Sleep(1000 * time.Second)
}
