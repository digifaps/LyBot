#!/usr/bin/env python3

import context
import paho.mqtt.client as mqtt
from Adafruit_MotorHAT import Adafruit_MotorHAT, Adafruit_DCMotor
import time
import atexit

# defines
DEADBAND = 20
DEADBAND_L = -20 

# create a default object, no changes to I2C address or frequency
mh = Adafruit_MotorHAT(addr=0x60,i2c_bus=1)

# recommended for auto-disabling motors on shutdown!
def turnOffMotors():
    mh.getMotor(1).run(Adafruit_MotorHAT.RELEASE)
    mh.getMotor(2).run(Adafruit_MotorHAT.RELEASE)
    mh.getMotor(3).run(Adafruit_MotorHAT.RELEASE)
    mh.getMotor(4).run(Adafruit_MotorHAT.RELEASE)

atexit.register(turnOffMotors)

# MQTT stuff
def on_message_motor_left(mosq, obj, msg):
    if int(msg.payload) > DEADBAND:
        myMotorL.setSpeed(int(msg.payload))
        myMotorL.run(Adafruit_MotorHAT.FORWARD);
    else:
        myMotorL.run(Adafruit_MotorHAT.RELEASE);
        print("MOTOR LEFT: " + msg.topic + " " + str(msg.qos) + " " + str(msg.payload))

def on_message_motor_right(mosq, obj, msg):
    if int(msg.payload) > DEADBAND:
        myMotorR.setSpeed(int(msg.payload))
        myMotorR.run(Adafruit_MotorHAT.BACKWARD);
    else:
        myMotorR.run(Adafruit_MotorHAT.RELEASE);
        print("MOTOR RIGHT: " + msg.topic + " " + str(msg.qos) + " " + str(msg.payload))

def on_message(mosq, obj, msg):
    print(msg.topic + " " + str(msg.qos) + " " + str(msg.payload))

mqttc = mqtt.Client()

myMotorL = mh.getMotor(1)
myMotorR = mh.getMotor(2)

# Add message callbacks that will only trigger on a specific subscription match.
mqttc.message_callback_add("$SYS/motor/left/#", on_message_motor_left)
mqttc.message_callback_add("$SYS/motor/right/#", on_message_motor_right)
mqttc.on_message = on_message
mqttc.connect("localhost", 1883, 60)
mqttc.subscribe("$SYS/#", 0)

mqttc.loop_forever()

# set the speed to start, from 0 (off) to 255 (max speed)
# myMotorL.setSpeed(150)
# myMotorR.setSpeed(150)
# myMotorL.run(Adafruit_MotorHAT.FORWARD);
# myMotorR.run(Adafruit_MotorHAT.BACKWARD);
# time.sleep(1)
# # turn off motor
# myMotorL.run(Adafruit_MotorHAT.RELEASE);
# myMotorR.run(Adafruit_MotorHAT.RELEASE);

# client.loop_forever()

