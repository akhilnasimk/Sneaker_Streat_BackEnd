package constent

// updating the status of a order
var AllowedOrderStatus = map[string]bool{ //this are the available status
	"processing": true,
	"shipped":    true,
	"delivered":  true,
}

var AllowedTransitions = map[string]string{ //this are allowed transactions
	"pending":    "processing",
	"processing": "shipped",
	"shipped":    "delivered",
}
