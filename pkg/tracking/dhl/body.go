package dhl

type body struct {
	Shipments []shipment `json:"seguimiento"`
}

type shipment struct {
	Date string `json:"Fecha"`
	Time string `json:"Hora"`
	Text string `json:"Descripcion"`
	City string `json:"Ciudad"`
}
