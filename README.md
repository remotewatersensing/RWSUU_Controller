# RWSUU_Controller

This is a Golang script that listens to the serial connection with the remote controller and posts the sensor data real time to ThingsBoard.
The data that is received by the remote controller follows this structure:

![RWSUU DAT structure](https://github.com/remotewatersensing/RWSUU-Diagrams/blob/main/diagrams/Datastructure.png?raw=true "RWSUU DAT structure")

It should be noted that an extra start bit (uint8 255) is added before this structure to note the beginning of the payload. This makes the total parsed payload 13 bytes 

## How to use
Make sure you have the latest version of Go installed. You need to install the go-serial dependency by running:

> go get "github.com/albenik/go-serial"

Run the script in the terminal with:

> go run main.go

Connect your remote controller and
