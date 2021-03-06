package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/nsf/termbox-go"
	"github.com/simulatedsimian/joystick"
)

// MQTT

type controller struct {
	axis   [8]int
	button [8]uint32
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	//	fmt.Printf("TOPIC: %s\n", msg.Topic())
	//	fmt.Printf("MSG: %s\n", msg.Payload())
}

func mqtt_init() *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("jetson-nano")
	opts.SetDefaultPublishHandler(f)
	return opts
}

func mqtt_connectClient(opts *MQTT.ClientOptions) MQTT.Client {
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return c
}

//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
func mqtt_publish(c MQTT.Client, topic string, text string) {

	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	token := c.Publish(topic, 0, false, text)
	token.Wait()

	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

}

func mqtt_disconnect(c MQTT.Client) {
	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("/"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}

// joystick
func printAt(x, y int, s string) {
	for _, r := range s {
		termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
}

func readJoystick(js joystick.Joystick, c MQTT.Client, cntl *controller) {
	jinfo, err := js.Read()

	if err != nil {
		printAt(1, 5, "Error: "+err.Error())
		return
	}

	printAt(1, 5, "Buttons:")
	for button := 0; button < js.ButtonCount(); button++ {
		if jinfo.Buttons&(1<<uint32(button)) != 0 {
			//cntl.axis[1] = jinfo.Buttons & (1 << uint32(button))
			printAt(10+button, 5, "X")
		} else {
			printAt(10+button, 5, ".")
		}
	}

	mqtt_publish(c, "axis5", fmt.Sprintf("Value: %7d", jinfo.AxisData[5]))

	for axis := 0; axis < js.AxisCount(); axis++ {
		cntl.axis[axis] = jinfo.AxisData[axis]
		printAt(1, axis+7, fmt.Sprintf("Axis %2d Value: %7d", axis, jinfo.AxisData[axis]))
	}

	return
}

func init_cntl(js joystick.Joystick) *controller {
	var cntl *controller
	for axis := 0; axis < js.AxisCount(); axis++ {
		cntl.axis[axis] = 0
	}
	for button := 0; button < js.ButtonCount(); button++ {
		cntl.button[button] = 0
	}
	return cntl
}

func main() {

	opts := mqtt_init()
	c := mqtt_connectClient(opts)
	// var cntl *controller

	jsid := 0
	if len(os.Args) > 1 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		jsid = i
	}

	js, jserr := joystick.Open(jsid)

	if jserr != nil {
		fmt.Println(jserr)
		return
	}
	cntl := init_cntl(js)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 40)

	for doQuit := false; !doQuit; {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Ch == 'q' {
					doQuit = true
				}
			}
			if ev.Type == termbox.EventResize {
				termbox.Flush()
			}

		case <-ticker.C:
			printAt(1, 0, "-- Press 'q' to Exit --")
			printAt(1, 1, fmt.Sprintf("Joystick Name: %s", js.Name()))
			printAt(1, 2, fmt.Sprintf("   Axis Count: %d", js.AxisCount()))
			printAt(1, 3, fmt.Sprintf(" Button Count: %d", js.ButtonCount()))
			readJoystick(js, c, cntl)
			termbox.Flush()
		}
	}
}
