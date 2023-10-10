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
		if loadBalancerIndex >= len(dstAddresses) {
			loadBalancerIndex = 0
		}
		dst := dstAddresses[loadBalancerIndex]
		select {
		case c := <-conn:
			log.Println("Proxying to", dst)
			destination, err := net.Dial("tcp", dst)
			if err != nil {
				loadBalancerIndex += 1
				log.Println("Unable to dial host", err)
				c.Close()
			} else {
				defer destination.Close()
				go func() {
					_, err = io.Copy(c, destination)
					if err != nil {
						log.Println("Unable to copy client", err)
						c.Close()
					}
				}()

				_, err = io.Copy(destination, c)
				if err != nil {
					log.Println("Unable to copy destination", err)
					c.Close()
				}
				loadBalancerIndex += 1
				c.Close()
			}

		default:
		}
	}

}
