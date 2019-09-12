package main

import (
	"fmt"
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"net"
	"time"
)

const WorkStarted = "work-started"

func main() {

	adapter := firmata.NewAdaptor("/dev/ttyACM0")
	sensor  := gpio.NewPIRMotionDriver(adapter, "5")
	led 	:= gpio.NewLedDriver(adapter, "13")

	work := func() {

		stash := logrus.New()

		conn, err := net.Dial("tcp", "media.local:5044")

		if err != nil {
			stash.Info("Error occurred")
			stash.Error(err)
		}

		hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type":"event"}))

		stash.Hooks.Add(hook)

		var start time.Time

		stash.WithFields(logrus.Fields{
			"value": WorkStarted,
		}).Info("Motion bot has started to work")

		_ = sensor.On(gpio.MotionDetected, func(data interface{}) {
			fmt.Println(gpio.MotionDetected)

			_ = led.On()

			start = time.Now()

			stash.WithFields(logrus.Fields{
				"value": gpio.MotionDetected,
			}).Info("Motion has been detected")
		})

		_ = sensor.On(gpio.MotionStopped, func(data interface{}) {

			fmt.Println(gpio.MotionStopped)

			_ = led.Off()

			end := time.Now()

			stash.WithFields(logrus.Fields{
				"value": gpio.MotionStopped,
				"elapsed_time": end.Sub(start).String(),
			}).Info("Motion has stopped")
		})
	}

	robot := gobot.NewRobot("motion-bot",
		[]gobot.Connection{adapter},
		[]gobot.Device{sensor, led},
		work,
	)

	_ = robot.Start()
}