# IIFL Capital Go client module.

The officail Go module to connect and get market data. BridgeGo module exposes four functionalities (connection,subscription,unsusbcription,disconnection) and event handlers.
Register the event handlers to receive and process data.

## Installation

```
go get github.com/IIFLSecurities/bridgeGo@v0.1.1
```

## This is a sample implementation of the package

```go
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	connector "github.com/IIFLSecurities/bridgeGo/connector"
)

// Declare all feed data type structures

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

type OpenInterestData struct {
	OpenInterest int32 `json:"openInterest"`
	DayHighOi    int32 `json:"dayHighOi"`
	DayLowOi     int32 `json:"dayLowOi"`
	PreviousOi   int32 `json:"previousOi"`
}

type LppData struct {
	LppHigh      uint32 `json:"lppHigh"`
	LppLow       uint32 `json:"lppLow"`
	PriceDivisor int32  `json:"priceDivisor"`
}

type UpperCircuitData struct {
	InstrumentId uint32 `json:"instrumentId"`
	UpperCircuit uint32 `json:"upperCircuit"`
	PriceDivisor int32  `json:"priceDivisor"`
}

type LowerCircuitData struct {
	InstrumentId uint32 `json:"instrumentId"`
	LowerCircuit uint32 `json:"lowerCircuit"`
	PriceDivisor int32  `json:"priceDivisor"`
}

type MarketStatusData struct {
	MarketStatusCode uint16 `json:"MarketStatusCode"`
}

type High52WeekData struct {
	InstrumentId uint32 `json:"instrumentId"`
	High52Week   uint32 `json:"52WeekHigh"`
	PriceDivisor int32  `json:"priceDivisor"`
}

type Low52WeekData struct {
	InstrumentId uint32 `json:"instrumentId"`
	Low52Week    uint32 `json:"52WeekLow"`
	PriceDivisor int32  `json:"priceDivisor"`
}

func main() {

	//Generate insntance of connector
	connector := connector.GetInstance()

 //Add OnDisconnet Hanlder so that module will call this function when disconnection happens
	connector.OnDisconnect = func(err error) {
		fmt.Println(err.Error())
	}


//Add MarketWatch handler to recieve MarketWatch data and process
	connector.MWHandler = func(b []byte, s string) {
		buffer := bytes.NewBuffer(b)

		// Create an instance of MWBOCombined
		var data MWBOCombined
		// Decode the byte array into the struct
		err := binary.Read(buffer, binary.LittleEndian, &data)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			return
		}
		fmt.Printf("MarketFeed Data of "+s+" : %+v\n", data)
	}


//Add Low52WeekHandler to recieve Low52Week data and process
	connector.Low52WeekHandler = func(b []byte, s string) {
		buffer := bytes.NewBuffer(b)

		// Create an instance of MWBOCombined
		var data Low52WeekData
		// Decode the byte array into the struct
		err := binary.Read(buffer, binary.LittleEndian, &data)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			return
		}
		if data.InstrumentId != 444540 {
		}

		fmt.Printf("Low52WeekData Data of "+s+" : %+v\n", data)
	}

//Add High52WeekHandler to recieve High52Week data and process
	connector.High52WeekHandler = func(b []byte, s string) {
		buffer := bytes.NewBuffer(b)

		// Create an instance of MWBOCombined
		var data High52WeekData
		// Decode the byte array into the struct
		err := binary.Read(buffer, binary.LittleEndian, &data)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			return
		}
		fmt.Printf("High52WeekData Data of "+s+" : %+v\n", data)
	}

//Add OpenInterstHandler to recieve OpenInterst data and process
	connector.OpenInterstHandler = func(b []byte, s string) {
		buffer := bytes.NewBuffer(b)
		// Create an instance of MWBOCombined
		var data OpenInterestData
		// Decode the byte array into the struct
		err := binary.Read(buffer, binary.LittleEndian, &data)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			return
		}

		fmt.Printf("OpenInterest Data of "+s+" : %+v\n", data)

	}

//Add LppHandler to recieve Lpp data and process
	connector.LppHandler = func(b []byte, s string) {
		buffer := bytes.NewBuffer(b)

		// Create an instance of MWBOCombined
		var data LppData
		// Decode the byte array into the struct
		err := binary.Read(buffer, binary.LittleEndian, &data)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			return
		}
		fmt.Printf("Lpp Data of "+s+" : %+v\n", data)
	}

//Add UpperCircuitHandler to recieve UpperCircuitdata and process
	connector.UpperCircuitHandler = func(b []byte, s string) {
		buffer := bytes.NewBuffer(b)
		// Create an instance of MWBOCombined
		var data UpperCircuitData
		// Decode the byte array into the struct
		err := binary.Read(buffer, binary.LittleEndian, &data)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			return
		}
		if data.InstrumentId != 444540 {
			return
		}
		fmt.Printf("UpperCircuit Data of "+s+" : %+v\n", data)
	}

//Add LowerCircuitHandler to receive LowerCircuit data and process
	connector.LowerCircuitHandler = func(b []byte, s string) {
		buffer := bytes.NewBuffer(b)

		// Create an instance of MWBOCombined
		var data LowerCircuitData
		// Decode the byte array into the struct
		err := binary.Read(buffer, binary.LittleEndian, &data)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			return
		}
		if data.InstrumentId != 444540 {
			return
		}
		fmt.Printf("LowerCircuit Data of "+s+" : %+v\n", data)
	}

//Add MarketStatusHandler to receive MarketStatus data and process
	connector.MarketStatusHandler = func(b []byte, s string) {
		buffer := bytes.NewBuffer(b)

		// Create an instance of MWBOCombined
		var data MarketStatusData
		// Decode the byte array into the struct
		err := binary.Read(buffer, binary.LittleEndian, &data)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			return
		}
		fmt.Printf("MarketStatus Data of "+s+" : %+v\n", data)
	}

//Add OrderUpdatesHandler to receive OrderUpdates  data and process
	connector.OrderUpdatesHandler = func(b []byte, s string) {
		fmt.Printf("Order Update  of "+s+" : %+v\n", string(b))
	}

//Add TradeUpdatesHandler to receive TradeUpdates data and process
	connector.TradeUpdatesHandler = func(b []byte, s string) {
		fmt.Printf("Trade Update  of "+s+" : %+v\n", string(b))
	}

	data := map[string]interface{}{
		"host":     "bridge.iiflcapital.com",
		"port":     9906,
		"password": "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIxVks4TEhlRnRvSmp6YWk1RmJlSGNPbDI3ekpGanBScTE2Vmt4eGJBZ0ZjIn0.eyJleHAiOjE3NTkxMjY4MDIsImlhdCI6MTc0MzU3NDgwMiwianRpIjoiNzE2YzUzMzUtZDY3MS00YzYwLWJmMTUtMzZkNWZkYzIzZTA4IiwiaXNzIjoiaHR0cHM6Ly9rZXljbG9hay5paWZsc2VjdXJpdGllcy5jb20vcmVhbG1zL0lJRkwiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNDk4NTU2NTAtMDcyNS00ZGUzLWFlMmQtN2ExMGYxYzIwODI4IiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiSUlGTCIsInNpZCI6IjA5YjYwZjM4LTE4N2EtNGFjMi1hMWJlLTRmMTgwNjkyNWZkMiIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovLzEwLjEyNS42OC4xNDQ6ODA4MC8iXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbImRlZmF1bHQtcm9sZXMtaWlmbCIsIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJJSUZMIjp7InJvbGVzIjpbIkdVRVNUX1VTRVIiLCJBQ1RJVkVfVVNFUiJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJ1Y2MiOiI5MzA4MDA0OCIsInNvbGFjZV9ncm91cCI6IlNVQlNDUklCRVJfQ0xJRU5UIiwibmFtZSI6IlNBVEhFRVNIIEtVTUFSIEdPUElTRVRUWSBOQSIsInByZWZlcnJlZF91c2VybmFtZSI6IjkzMDgwMDQ4IiwiZ2l2ZW5fbmFtZSI6IlNBVEhFRVNIIEtVTUFSIEdPUElTRVRUWSIsImZhbWlseV9uYW1lIjoiTkEiLCJlbWFpbCI6InNhdGlzaGt1bWFyZ29waXNldHR5MTExQGdtYWlsLmNvbSJ9.oUofcKL1c2T-EI8njv9XW96pdduK7HmUtZRO4A2Vf4l4CbxsX561VJRsC93UGD8u7zXOJHuciD-YLo_IYEylqW6rzwx4oBsJyonB1_xhankRwkxVp7vqNbbrF6RKbuK32MvG3sTAJIk2mh2m9XddS1GFfxuC285DbuBi58Iiw4gpmQSV50dffJ1RgnuIyrtRY5WNzsqQwhgMM7RhGZ3lWL4ZuhEKXUQixs4VPLts2q6FiHfx_JH9HM98uOZQODRSCo7mOIpjpos6GgbbPMxLVjGw-iAmEKPqusPTvoNddwlEvsOWxCAHai6Wq62QfPQnnypXB95_FdEIycBz2wU7Sw",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	//Call ConnectHost function to connect
	res, err := connector.ConnectHost(string(jsonData))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(res)
	} else {
		fmt.Println(res)
	}

	// res1, err1 := connector.ConnectHost(string(jsonData))
	// if err1 != nil {
	// 	fmt.Println(err1.Error())
	// 	fmt.Println(res1)
	// } else {
	// 	fmt.Println(res1)
	// }

	//Call SubscribeFeed function to subscribe and receive data
	subReq := map[string]interface{}{
		"SubscriptionList": []string{"nsef/54452", "nseeq/2885", "bseeq*/1234"}}

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

	// //	Subscribe LPP
	// subReq := map[string]interface{}{
	// 	"SubscriptionList": []string{"nsefo/54452"}}
	// jsonData, err = json.Marshal(subReq)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// subRes, err := connector.SubscribeLpp(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(subRes)
	// }

	// //	Subscribe OpenInterst
	// subReq2 := map[string]interface{}{
	// 	"SubscriptionList": []string{"nsefo/54452"}}

	// jsonData, err = json.Marshal(subReq2)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// subRes2, err := connector.SubscribeOpenInterest(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(subRes2)
	// }

	// //subscribe 52High
	// subReq3 := map[string]interface{}{
	// 	"SubscriptionList": []string{"mcxcomm"}}
	// jsonData, err = json.Marshal(subReq3)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// subRes3, err := connector.SubscribeHigh52Week(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(subRes3)
	// }

	// //Subscribe 52Low
	// subReq4 := map[string]interface{}{
	// 	"SubscriptionList": []string{"mcxcomm"}}
	// jsonData, err = json.Marshal(subReq4)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// subRes4, err := connector.SubscribeLow52Week(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(subRes4)
	// }

	// //Subscribe UpperCircuit
	// subReq5 := map[string]interface{}{
	// 	"SubscriptionList": []string{"mcxcomm"}}
	// jsonData, err = json.Marshal(subReq5)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// subRes5, err := connector.SubscribeUpperCircuit(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(subRes5)
	// }

	// //Subscribe LowerCircuit
	// subReq6 := map[string]interface{}{
	// 	"SubscriptionList": []string{"mcxcomm"}}
	// jsonData, err = json.Marshal(subReq6)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// subRes6, err := connector.SubscribeLowerCircuit(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(subRes6)
	// }

	// //Subcribe MarketStatus
	// subReq7 := map[string]interface{}{
	// 	"SubscriptionList": []string{"mcxcomm"}}
	// jsonData, err = json.Marshal(subReq7)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// subRes7, err := connector.SubscribeLowerCircuit(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(subRes7)
	// }

	// //Subcribre Order
	// subReq8 := map[string]interface{}{
	// 	"SubscriptionList": []string{"nishtha9"}}
	// jsonData, err = json.Marshal(subReq8)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// subRes8, err := connector.SubscribeOrderUpdates(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(subRes8)
	// }

	// //Subcribre Trade
	// subReq9 := map[string]interface{}{
	// 	"SubscriptionList": []string{"nishtha9"}}
	// jsonData, err = json.Marshal(subReq9)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// subRes9, err := connector.SubscribeTradeUpdates(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(subRes9)
	// }

	time.Sleep(15 * time.Second)

	//Call UnsubscribeFeed function to Unsubscribe and stops receiving data
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

	//Unsubscribe OpenInterest
	// UnSubReq1 := map[string]interface{}{
	// 	"UnsubscriptionList": []string{"nseeq/2885", "nsee/2886", "bseeq/2887"},
	// }
	// jsonData, err = json.Marshal(UnSubReq1)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// unsubRes1, err := connector.UnsubscribeOpenInterest(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(unsubRes1)
	// }

	// //Unsubscribe Lpp
	// UnSubReq2 := map[string]interface{}{
	// 	"UnsubscriptionList": []string{"nseeq"},
	// }
	// jsonData, err = json.Marshal(UnSubReq2)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// unsubRes2, err := connector.UnsubscribeLpp(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(unsubRes2)
	// }

	// //Unsubscribe High52Week
	// UnSubReq3 := map[string]interface{}{
	// 	"UnsubscriptionList": []string{"mcxcomm"},
	// }
	// jsonData, err = json.Marshal(UnSubReq3)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// unsubRes3, err := connector.UnsubscribeHigh52Week(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(unsubRes3)
	// }

	// //Unsubscribe Low52Week
	// UnSubReq4 := map[string]interface{}{
	// 	"UnsubscriptionList": []string{"mcxcomm"},
	// }
	// jsonData, err = json.Marshal(UnSubReq4)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// unsubRes4, err := connector.UnsubscribeLow52Week(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(unsubRes4)
	// }

	// //Unsubscribe UpperCircuit
	// UnSubReq5 := map[string]interface{}{
	// 	"UnsubscriptionList": []string{"mcxcomm"},
	// }
	// jsonData, err = json.Marshal(UnSubReq5)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// unsubRes5, err := connector.UnsubscribeUpperCircuit(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(unsubRes5)
	// }

	// //Unsubscribe LowerCircuit
	// UnSubReq6 := map[string]interface{}{
	// 	"UnsubscriptionList": []string{"mcxcomm"},
	// }
	// jsonData, err = json.Marshal(UnSubReq6)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// unsubRes6, err := connector.UnsubscribeLowerCircuit(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(unsubRes6)
	// }

	// //Unsubscribe MarketStatus
	// UnSubReq7 := map[string]interface{}{
	// 	"UnsubscriptionList": []string{"nseeq"},
	// }
	// jsonData, err = json.Marshal(UnSubReq7)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// unsubRes7, err := connector.UnsubscribeMarketStatus(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(unsubRes7)
	// }

	// //Unsubscribe Order
	// UnSubReq8 := map[string]interface{}{
	// 	"UnsubscriptionList": []string{"93080048"},
	// }
	// jsonData, err = json.Marshal(UnSubReq8)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// unsubRes8, err := connector.UnsubscribeOrderUpdates(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(unsubRes8)
	// }

	// //Unsubscribe Trade
	// UnSubReq9 := map[string]interface{}{
	// 	"UnsubscriptionList": []string{"93080048"},
	// }
	// jsonData, err = json.Marshal(UnSubReq9)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// unsubRes9, err := connector.UnsubscribeTradeUpdates(string(jsonData))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(unsubRes9)
	// }

	time.Sleep(150 * time.Second)

	//Call DisconnectHost to disconnect
	disRes, err := connector.DisconnectHost()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(disRes)
	}

	time.Sleep(1000 * time.Second)
}

```
