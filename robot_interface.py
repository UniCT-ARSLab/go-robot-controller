import sys
import struct
import smbus
import time
import numpy as np
import RPi.GPIO as GPIO


class RobotInterface:

    def __init__(self):
        GPIO.setmode(GPIO.BCM)
        GPIO.setup(4,GPIO.OUT)
        GPIO.output(4,GPIO.HIGH)

        self.bus = smbus.SMBus(1)
        
    def __send_command(self, *args):
        i2c_data = []
        for c in struct.pack("<BhhhB", *args):
            #i2c_data.append( ord(c) )
            i2c_data.append(c)

        #print( len(i2c_data), i2c_data )

        try:
            self.bus.write_i2c_block_data(0x34, 0x60, i2c_data)
            time.sleep(0.1)
        except:
            print("Exception in write_i2c_block_data... retrying")
            self.board_reset()

    def set_speed(self, speed):
        self.__send_command(0x8C, speed, 0, 0, 0)

        
    def forward_to_distance(self, dist):
        self.__send_command(0x85, dist, 0, 0, 0)

        
    def rotate_relative(self, angle):
        self.__send_command(0x88, angle, 0, 0, 0)

        
    def board_reset(self):
        #print("Resetting I2C-to-CAN")
        GPIO.output(4,GPIO.LOW)
        time.sleep(0.1)
        GPIO.output(4,GPIO.HIGH)
        time.sleep(2)
        #print("Reset done")

    def get_position(self):
        try:
            pos = self.bus.read_i2c_block_data(0x34, 1, 6) # Read a block of 6 bytes from address 0x34, offset 1
            pos_data = b""
            for a in pos:
                pos_data = pos_data + bytes([a])

            [x,y,t] = struct.unpack("<hhh", pos_data)
            #print(x,",",y,",",t)
            t = t / 100.0
            coordinate = [x,y,t]
            return coordinate
            
        except:
            print("Exception in read_i2c_block_data... retrying")
            self.board_reset()
            return None

        
if __name__ == "__main__":
    robot_if = RobotInterface()
    robot_if.board_reset()
    robot_if.set_speed(200)
    
    starting_pos = np.asarray( robot_if.get_position() )
    print("Starting position:",starting_pos)
    # robot_if.forward_to_distance(500) # mm
    # robot_if.rotate_relative(-10) # clockwise

    
    for _ in range(2):
        robot_if.forward_to_distance(250) # mm
        #time.sleep(4)
        pos = np.asarray( robot_if.get_position() ) - starting_pos
        print(pos)
"""
    for _ in range(3):
        robot_if.rotate_relative(-60)
        time.sleep(4)
        pos = np.asarray( robot_if.get_position() ) - starting_pos
        ##print(pos)
    
    for _ in range(2):
        #robot_if.forward_to_distance(250) # mm
        time.sleep(4)
        pos = np.asarray( robot_if.get_position() ) - starting_pos
        #print(pos)
"""