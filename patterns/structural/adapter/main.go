package main

import "fmt"

type client struct{}

func (c *client) insertLightningConnectorIntoComp(com computer) {
	fmt.Println("Client inserts lightning connector into computer.")
	com.insertIntoLightningPort()
}

type computer interface {
	insertIntoLightningPort()
}

type mac struct{}

func (m *mac) insertIntoLightningPort() {
	fmt.Println("Lightning connector is plugged into mac computer.")
}

type windows struct{}

func (w *windows) insertIntoUSBPort() {
	fmt.Println("USB connector is plugged into windows machine.")
}

type windowsAdapter struct {
	windowsMachine *windows
}

func (w *windowsAdapter) insertIntoLightningPort() {
	fmt.Println("Adapter converts lightning signal to USB.")
	w.windowsMachine.insertIntoUSBPort()
}

func main() {
	client := &client{}
	mac := &mac{}

	client.insertLightningConnectorIntoComp(mac)

	windowsMachine := &windows{}
	windowsMachineAdapter := &windowsAdapter{
		windowsMachine: windowsMachine,
	}

	client.insertLightningConnectorIntoComp(windowsMachineAdapter)
}
