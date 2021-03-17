# RobotController (Service) by ARSLab

This is a service for monitoring and controlling robotic systems. It is based mainly on I2C communication to send commands or read information directly from the robot unit control. The service also provides a series of HTTP Rest API with which to execute the same commands and a Websocket (pure or via Socketio) with which to interact in realtime with the robot.

## Requirements
The service is written in GoLang version 1.15, then to be compiled you must have almost this version installed in your machine.

## Tests
The service was tested into a Raspberry PI 3 with Raspian (debian linux based) with a I2C board which is used to communicate with the robot system.
In addition, the service was used to develop a useful tool to display and interact with the robot in a virtual environment.

## Installation
<ul>
<li>Install GoLang (1.15)</li>
<li>Install the statik lib (<code>go get -u github.com/rakyll/statik</code>)</li>
<li>Run the <code>make install-dep</code> command</li>
<li>Run the <code>make build-ui</code> command</li>
</ul>

## Build and Run
Execute the command <code>make build</code> to build the final binary and <code>make run</code> to execute it (or execute directly the binary file).

# Project Structure and Description
The service is written in GoLang, therefore there are no concepts such as classes or objects (such as in c). The real "main" file is inside of <code>cmd</code> directory. There are created the robot instance, the webserver and the main loop (an empty loop).

## Models
In this directory you can find all <code>structs</code> that rappresent every payload/message or "object" structure.

## Robot
In this directory are defined the <code>connection</code> struct with its functions and the <code>robot</code> struct. 
In the <code>robot</code> struct are defined all commands to be send to the connection throught the <code>connection</code> instance.

## Webserver
In this directory is defined the <code>webserver</code> struct and its functions. When a webserver is created (with a <code>robot</code> pointer instance, address and port) the http routes and websocket server are defined. 

## Utilities
Some struct and fuctions created as utilities.

# TODOs
<ul>
<li>Allowing the run of the service with a virtual robot instance.</li>
<li>Fix the <code>set_position</code> operation in the robot struct.</li>
<li>Create a WebUI to send and controll the local robot instance.</li>
<li>Other...</li>
</ul>