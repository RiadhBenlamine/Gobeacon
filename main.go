package main

import (
	"fmt"
)

func main() {
	device := Device{}
	info := device.run()
	fmt.Println(info)

}
