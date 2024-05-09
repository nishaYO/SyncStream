package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func saveBufferToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = buf.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}

// parse event from request body of the producer
func parseRequest(body io.Reader, event *Event) error {
	_, err := fmt.Fscanf(body, "Key=%s\nValue=%s", &event.Key, &event.Value)
	log.Print(event)
	return err
}

func saveEventToBuffer(event Event) error {
	// convert key value pair into binary
	keyByte := []byte(event.Key)
	valueByte := []byte(event.Value)
	// length of key and value bytes
	keyByteLen := uint16(len(keyByte))
	valueByteLen := uint16(len(valueByte))

	cap := 4 + len(keyByte) + len(valueByte)
	data := make([]byte, 0, cap)
	data = append(data, byte(keyByteLen), byte(keyByteLen>>8)) // convert keyByteLen to little endian from big endian using bit manipulation
	data = append(data, keyByte...)                                
	data = append(data, byte(valueByteLen), byte(valueByteLen>>8)) 
	data = append(data, valueByte...)                              

	// write key value and len of both to the buffer
	if _, err := buf.Write(data); err != nil {

		if err.Error() == "ErrTooLarge" {
			// save the buffer to the file
			if err := saveBufferToFile("log.bin"); err != nil {
				return err
			}

			buf.Reset() // Clear the buffer

			// rewrite the event
			if err := saveEventToBuffer(event); err != nil {
				return err
			}
			log.Println("Buffer size exceeded, discarding current data batch")
		} else {
			return err
		}

	}

	return nil
}
