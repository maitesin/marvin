package app

import "fmt"

const errMsgDeliveryNotFound = "delivery %q not found"

type DeliveryNotFound struct {
	identifier string
}

func NewDeliveryNotFound(identifier string) DeliveryNotFound {
	return DeliveryNotFound{identifier: identifier}
}

func (dnf DeliveryNotFound) Error() string {
	return fmt.Sprintf(errMsgDeliveryNotFound, dnf.identifier)
}
