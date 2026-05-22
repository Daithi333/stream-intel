package model

type TaxiTrip struct {
	EventId        string  `json:"event_id"`
	EventTs        string  `json:"event_ts"`
	VendorId       int     `json:"vendor_id"`
	PickupTs       string  `json:"pickup_ts"`
	DropoffTs      string  `json:"dropoff_ts"`
	PuLocationId   int     `json:"pu_location_id"`
	DoLocationId   int     `json:"do_location_id"`
	PassengerCount int     `json:"passenger_count"`
	TripDistance   float64 `json:"trip_distance"`
	FareAmount     float64 `json:"fare_amount"`
	TipAmount      float64 `json:"tip_amount"`
	TotalAmount    float64 `json:"total_amount"`
	PaymentType    int     `json:"payment_type"`
	RatecodeId     int     `json:"ratecode_id"`
}
