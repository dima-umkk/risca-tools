package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dima-kgd/risca-tools/internal/emul"
)

func main() {
	fmt.Println("RiscA DisAssembler v.1.0.0")
	ifileName := flag.String("i", "input.bin", "Input file")
	flag.Usage = func() {
		fmt.Printf("Usage: TODO")
	}
	flag.Parse()

	ifile, err := os.Open(*ifileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer ifile.Close()

	cpu := emul.CPU{Bus: emul.BusDevice{Devices: []emul.Device{emul.NewMemoryDevice(0x0000_0000, 0x0000_FFFF), emul.NewUartDevice()}}}

	var curAddr uint32 = 0
	for {
		var bt uint8
		err = binary.Read(ifile, binary.LittleEndian, &bt)
		if err != nil {
			break
		}
		cpu.Bus.Write8(curAddr, bt)
		curAddr += 1
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">")
	text, _ := reader.ReadString('\n')
	for text != "q" {
		switch strings.TrimSpace(text) {
		case "s":
			i := cpu.Peek()
			fmt.Printf("0x%08X: %v\n", cpu.PC, i)
			cpu.Step()
		case "r":
			for i, r := range cpu.Registers {
				fmt.Printf("R%d 0x%08X %d\n", i, r, int32(r))
			}
		default:
			fmt.Print("Unknown command\n")
		}
		fmt.Print(">")
		text, _ = reader.ReadString('\n')
	}

}
