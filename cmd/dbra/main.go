package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"

	"github.com/dima-kgd/risca-tools/internal/emul"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	fmt.Println("RiscA Emulator v.1.0.0")
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

	app := tview.NewApplication()
	regs := tview.NewTable()
	regs.SetBorder(true)
	regs.SetTitle("Registers")
	regs.SetCell(0, 0, tview.NewTableCell("PC"))
	regs.SetCell(0, 1, tview.NewTableCell(fmt.Sprintf("0x%08X", cpu.PC)))
	for i, r := range cpu.Registers {
		regs.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("R%d", i)))
		regs.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprintf("0x%08X", r)))
	}
	disasm := tview.NewTextView()
	disasm.SetBorder(true).SetTitle("Disassembly")
	disasm.SetDynamicColors(true)
	stackview := tview.NewTextView()
	stackview.SetBorder(true).SetTitle("Stack")
	stackview.SetDynamicColors(true)
	memview := tview.NewTextView()
	memview.SetBorder(true).SetTitle("Memory")
	memview.SetDynamicColors(true)
	programview := tview.NewTextView()
	programview.SetBorder(true).SetTitle("Program")
	programview.SetDynamicColors(true)
	inputfield := tview.NewInputField()
	inputfield.SetLabel("Input: ")

	textText := ""
	for i := range 10 {
		textText += fmt.Sprintf("0x%08X: %v\n", cpu.Peek(int32(i)).Address, cpu.Peek(int32(i)))
	}
	disasm.SetText(textText)

	right := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(regs, 0, 1, false).AddItem(stackview, 0, 1, false)
	left := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(disasm, 0, 3, true).AddItem(memview, 0, 2, false)

	debugPage := tview.NewFlex().AddItem(left, 0, 3, true).AddItem(right, 40, 1, false)
	programPage := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(programview, 0, 1, false).AddItem(inputfield, 3, 0, true)

	pages := tview.NewPages().AddPage("debug", debugPage, true, true).AddPage("program", programPage, true, false)

	updateRegs := func() {
		regs.Clear()
		regs.SetCell(0, 0, tview.NewTableCell("PC"))
		regs.SetCell(0, 1, tview.NewTableCell(fmt.Sprintf("0x%08X", cpu.PC)))
		for i, r := range cpu.Registers {
			regs.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("R%d", i)))
			regs.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprintf("0x%08X", r)))
		}
	}

	updateDisasm := func() {
		disasm.Clear()
		pc := int(cpu.PC)
		//_, _, _, th := disasm.GetInnerRect()
		start := pc - 10 //5 instructins
		if start < 0 {
			start = 0
		}
		for i := 0; i < 20; i++ {
			idx := start + i
			instr := cpu.Peek(int32(idx))
			if int(instr.Address) == pc {
				fmt.Fprintf(disasm, "[black:yellow]> 0x%08X: %v[-:-]\n", instr.Address, instr)
			} else {
				fmt.Fprintf(disasm, "  0x%08X: %v\n", instr.Address, instr)
			}
		}
	}

	refresh := func() {
		updateRegs()
		updateDisasm()
	}

	refresh()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF10:
			cpu.Step()
			refresh()
			return nil
		}
		switch event.Rune() {
		case 's':
			cpu.Step()
			regs.SetCell(0, 0, tview.NewTableCell("PC"))
			regs.SetCell(0, 1, tview.NewTableCell(fmt.Sprintf("0x%08X", cpu.PC)))
			for i, r := range cpu.Registers {
				regs.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("R%d", i)))
				regs.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprintf("0x%08X", r)))
			}
			textText := ""
			_, _, _, th := disasm.GetInnerRect()
			for i := range th {
				if i == 0 {
					textText += fmt.Sprintf("[red]0x%08X: %v\n", cpu.Peek(int32(i)).Address, cpu.Peek(int32(i)))
				} else {
					textText += fmt.Sprintf("[while]0x%08X: %v\n", cpu.Peek(int32(i)).Address, cpu.Peek(int32(i)))
				}
			}
			disasm.SetText(textText)
			return nil
		case 'd':
			pages.SwitchToPage("debug")
		case 'p':
			pages.SwitchToPage("program")
		case 'q':
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
