package transptypes

import(
	"https://github.com/vantihovich/go_tasks/tree/master/patterns/creational/factory/composer"
)

type vehicle struct {
	transport
}

func newVehicle() commonTransport {
	return &vehicle{
		transport: transport{
			name:       "TranspKings",
			maxWeight:  "20",
			transpType: "truck",
		},
	}
}

type ship struct {
	transport
}

func newShip() commonTransport {
	return &ship{
		transport: transport{
			name:       "Very Big Ship",
			maxWeight:  "400000",
			transpType: "ship",
		},
	}
}
