To start: `go run main.go -port 8080 -fetch-interval 5`

Available flags:
- port (int): port to run http server on
- fetch-interval (int): interval between data fetches (sec)
- enable-binding: set to enable sending req while monitoring

The app is REST API server, being a microservice to be integrated.

API spec is POST 'localhost:{{port}}' with payload
```json
{
    "symbol": "string",
    "price": uint64
}
```

After receiving request, it will immediatly return 200 ok. Goroutine is used to branch out a thread, broadcasting the transaction in the background and keep fetching transaction status every {{fetch-interval}} seconds. The app will keep logging activity info including each transaction monitoring status. By integrating it further, I have implemented function 'sendEndpointRequest' so it will POST to endpoint with payload `{ "tx_hash": "string", "tx_status": "string" }` every time the hash get fetched, to use it: update var ENDPOINT in the code and set the flag `-enable-binding` when run the app.

While monitoring the status, if the status is not 'PENDING' then it will stop monitoring that transaction. I'm not sure about the use cases or business needs so I'm not sure how each status should be handled, but for 'FAILED' it should get retry again, or save all the statuses to db and notify user further if it failed.