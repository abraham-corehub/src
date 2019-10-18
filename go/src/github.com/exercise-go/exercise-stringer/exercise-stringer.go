package main

import "fmt"

//IPAddr is the IP Address
type IPAddr [4]byte

func (x IPAddr) String() string {
	return fmt.Sprintf("%v.%v.%v.%v", x[0], x[1], x[2], x[3])
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 0},
		"googleDNS": {8, 8, 8, 8},
	}

	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}

}
