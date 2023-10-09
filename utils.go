package main

import (
	"io"
	"log"
	"net"
)

func handleErrors(prompt string, err error) {
	if err != nil {
		log.Fatalln(prompt, err)
	}
}

func proxy(conn chan net.Conn, dstAddresses []string) {
	loadBalancerIndex := 0

	for {
		dst := dstAddresses[loadBalancerIndex]
		select {
		case c := <-conn:
			log.Println("Proxying to", dst)
			destination, err := net.Dial("tcp", dst)
			handleErrors("Error dialing backend", err)
			defer destination.Close()
			go func() {
				_, err = io.Copy(c, destination)
				handleErrors("Error copying client", err)
			}()

			_, err = io.Copy(destination, c)
			handleErrors("Error copying client", err)
			loadBalancerIndex += 1
			if loadBalancerIndex >= len(dstAddresses) {
				loadBalancerIndex = 0
			}
			c.Close()
		default:
		}
	}

}
