package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"p3/constants"
	"p3/model"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	PORT           *int
	FETCH_INTERVAL *int
	ENABLE_BINDING *bool
	ENDPOINT       = ""
)

func init() {
	PORT = flag.Int("port", 8080, "port to run http server on")
	FETCH_INTERVAL = flag.Int("fetch-interval", 10, "interval between data fetches (sec)")
	ENABLE_BINDING = flag.Bool("enable-binding", false, "set to true to enable sending req while monitoring")
	flag.Parse()
}

func main() {
	r := gin.Default()

	r.POST("", func(c *gin.Context) {
		var req model.TransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.Response{
				Status:  "Bad Request",
				Message: err.Error(),
			})
			return
		}

		req.Timestamp = int(time.Now().Unix())

		logInfo(fmt.Sprintf("Request Received\tSymbol: %s\tPrice: %d\tTimestamp: %d", req.Symbol, req.Price, req.Timestamp))

		go broadcastTransaction(&req)

		res := model.Response{
			Status: "Success",
		}

		c.JSON(http.StatusOK, res)
	})

	r.Run(fmt.Sprintf(":%d", *PORT))
}

func broadcastTransaction(req *model.TransactionRequest) {
	bReq, err := json.Marshal(req)
	if err != nil {
		logInfo(fmt.Sprintf("Error\tSymbol: %s\tPrice: %d\tTimestamp: %d\nbroadcastTransaction|Error marshaling JSON: %v", req.Symbol, req.Price, req.Timestamp, err))
		return
	}

	httpRes, err := http.Post("https://mock-node-wgqbnxruha-as.a.run.app/broadcast", "application/json", bytes.NewBuffer(bReq))
	if err != nil {
		logInfo(fmt.Sprintf("Error\tSymbol: %s\tPrice: %d\tTimestamp: %d\nbroadcastTransaction|Error making POST request: %v", req.Symbol, req.Price, req.Timestamp, err))
		return
	}
	defer httpRes.Body.Close()

	bRes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		logInfo(fmt.Sprintf("Error\tSymbol: %s\tPrice: %d\tTimestamp: %d\nbroadcastTransaction|Error reading resonse body: %v", req.Symbol, req.Price, req.Timestamp, err))
		return
	}

	var res model.TransactionHash
	if httpRes.StatusCode == http.StatusOK {
		if err := json.Unmarshal(bRes, &res); err != nil {
			logInfo(fmt.Sprintf("Error\tSymbol: %s\tPrice: %d\tTimestamp: %d\nbroadcastTransaction|Error unmarshaling response JSON: %v", req.Symbol, req.Price, req.Timestamp, err))
			return
		}
		logInfo(fmt.Sprintf("Broadcast Tx Success\tSymbol: %s\tPrice: %d\tTimestamp: %d\tTxHash: %s", req.Symbol, req.Price, req.Timestamp, res.TxHash))
	} else {
		logInfo(fmt.Sprintf("Broadcast Tx Fail\tSymbol: %s\tPrice: %d\tTimestamp: %d\nStatusCode: %d\tResponse: %s", req.Symbol, req.Price, req.Timestamp, httpRes.StatusCode, string(bRes)))
		return
	}

	monitorTransaction(req, res.TxHash)
}

func monitorTransaction(req *model.TransactionRequest, txHash string) {
	ticker := time.NewTicker(time.Duration(*FETCH_INTERVAL) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			httpRes, err := http.Get("https://mock-node-wgqbnxruha-as.a.run.app/check/" + txHash)
			if err != nil {
				logInfo(fmt.Sprintf("Error\tSymbol: %s\tPrice: %d\tTimestamp: %d\nbroadcastTransaction|Error making GET request: %v", req.Symbol, req.Price, req.Timestamp, err))
				continue
			}
			defer httpRes.Body.Close()

			bRes, err := io.ReadAll(httpRes.Body)
			if err != nil {
				logInfo(fmt.Sprintf("Error\tSymbol: %s\tPrice: %d\tTimestamp: %d\nbroadcastTransaction|Error reading resonse body: %v", req.Symbol, req.Price, req.Timestamp, err))
				continue
			}

			var res model.TransactionStatus
			if httpRes.StatusCode == http.StatusOK {
				if err := json.Unmarshal(bRes, &res); err != nil {
					logInfo(fmt.Sprintf("Error\tSymbol: %s\tPrice: %d\tTimestamp: %d\nbroadcastTransaction|Error unmarshaling response JSON: %v", req.Symbol, req.Price, req.Timestamp, err))
					continue
				}

				logInfo(fmt.Sprintf("Fetch Tx Status\t\tSymbol: %s\tPrice: %d\tTimestamp: %d\tTxHash: %s\nStatus: %s", req.Symbol, req.Price, req.Timestamp, txHash, res.TxStatus))

				if *ENABLE_BINDING {
					sendEndpointRequest(&model.EndpointRequest{TxHash: txHash, TxStatus: res.TxStatus})
				}

				if res.TxStatus != constants.PENDING {
					return
				}
			} else {
				logInfo(fmt.Sprintf("Fetch Tx Status Fail\tSymbol: %s\tPrice: %d\tTimestamp: %d\tTxHash: %s\nStatusCode: %d\tResponse: %s", req.Symbol, req.Price, req.Timestamp, txHash, httpRes.StatusCode, string(bRes)))
				continue
			}
		}
	}
}

func sendEndpointRequest(req *model.EndpointRequest) {
	bReq, err := json.Marshal(req)
	if err != nil {
		logInfo(fmt.Sprintf("Error\tTxHash: %s\tTxStatus: %s\nsendEndpointRequest|Error marshaling JSON: %v", req.TxHash, req.TxStatus, err))
		return
	}

	httpRes, err := http.Post(ENDPOINT, "application/json", bytes.NewBuffer(bReq))
	if err != nil {
		logInfo(fmt.Sprintf("Error\tTxHash: %s\tTxStatus: %s\nsendEndpointRequest|Error making POST request: %v", req.TxHash, req.TxStatus, err))
		return
	}
	defer httpRes.Body.Close()

	bRes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		logInfo(fmt.Sprintf("Error\tTxHash: %s\tTxStatus: %s\nsendEndpointRequest|Error reading resonse body: %v", req.TxHash, req.TxStatus, err))
		return
	}

	var res interface{}
	if httpRes.StatusCode == http.StatusOK {
		if err := json.Unmarshal(bRes, &res); err != nil {
			logInfo(fmt.Sprintf("Error\tTxHash: %s\tTxStatus: %s\nsendEndpointRequest|Error unmarshaling response JSON: %v", req.TxHash, req.TxStatus, err))
			return
		}
		logInfo(fmt.Sprintf("POST Endpoint Success\tTxHash: %s\tTxStatus: %s", req.TxHash, req.TxStatus))
	} else {
		logInfo(fmt.Sprintf("POST Endpoint Fail\tTxHash: %s\tTxStatus: %s\nStatusCode: %d\tResponse: %s", req.TxHash, req.TxStatus, httpRes.StatusCode, string(bRes)))
		return
	}
}

func logInfo(s string) {
	fmt.Println("------------------------------------------------------------")
	fmt.Println(s)
	fmt.Println("------------------------------------------------------------")
}
