package main

import (
	"net"
	"fmt"
	"bufio"
	"strconv"
	"math"
	"encoding/binary"
	"time"
)

// Define BYOND constant type return values
const TYPE_NULL = 0
const TYPE_FLOAT = 0x2a
const TYPE_STRING = 0x06

/*
* Usage: data should be a string like "?ping". Question mark is required.
* Return values: Length of returned data (int). Raw data (byte splice)
*/
func byond_topic(data string, address string, port int) (error, int, []byte) {
	// BYOND packet forging nonsense. b_msg will be our final message to send.
	pref := "\x00\x83\x00"
	leng := byte((len(data)+6))
	pad := "\x00\x00\x00\x00\x00"
	suff := "\x00"
	b_msg := pref+string(leng)+pad+data+suff

	// Make the TCP connection
	conn, err := net.DialTimeout("tcp", address+":"+strconv.Itoa(port),3*time.Second)
	if err != nil {
		// Connection timeout, or other error.
		return err,-1,make([]byte,1,1)
	}

	// Ship request to server.
	fmt.Fprintf(conn, b_msg)
	// BYOND response length can be up to 0xFFFF. We want to ensure we have room.
	dat := make([]byte, 65541, 65541)
	read,err := bufio.NewReader(conn).Read(dat)
	conn.Close()
	if err != nil {
		return err,-1,make([]byte,1,1)
	}
	// Below for posterity is code to use the length value BYOND sends back. Normally I always read the entire response at once, so I will instead use the number of bytes returned by Read() above in a slice.
	//read = uint16(dat[2])<<8 | uint16(dat[3])
	rtype := int(dat[4])
	return nil,rtype,dat[5:read]
}

func float32toint(flt []byte) int {
	return int(math.Float32frombits(binary.LittleEndian.Uint32(flt)))
}

func main() {
	err,rtype,pdata := byond_topic("?ping", "play.cm-ss13.com", 1400)
	if err != nil {
		panic(err)
	}
	switch(rtype) {
	case TYPE_NULL:
		fmt.Println("NULL type returned from server")
	case TYPE_FLOAT:
		players := float32toint(pdata)
		fmt.Println(strconv.Itoa(players)+" players online")
		fmt.Println("Float returned from server")
	case TYPE_STRING:
		fmt.Println(string(pdata))
		fmt.Println("String returned from server")
	}
}
