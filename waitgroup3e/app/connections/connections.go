package connections

import (
	"main/app/inventory"
)

type Connectors interface {
	Open(*inventory.Host)  (error)
	Close()
}