package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/mattn/go-runewidth"
	"os"
	"strings"
)

func suddenDeath(msg string) {
	lines := strings.Split(msg, "\n")
	widths := []int{}

	maxWidth := 0
	for _, line := range lines {
		width := runewidth.StringWidth(line)
		widths = append(widths, width)
		if maxWidth < width {
			maxWidth = width
		}
	}

	ct.ChangeColor(ct.Red, true, ct.None, false)
	fmt.Println("＿" + strings.Repeat("人", maxWidth/2+2) + "＿")
	for i, line := range lines {
		ct.ChangeColor(ct.Red, true, ct.None, false)
		fmt.Print("＞　")
		ct.ChangeColor(ct.Yellow, true, ct.None, false)
		fmt.Print(line + strings.Repeat(" ", maxWidth-widths[i]))
		ct.ChangeColor(ct.Red, true, ct.None, false)
		fmt.Println("　＜")
	}
	ct.ChangeColor(ct.Red, true, ct.None, false)
	fmt.Println("￣" + strings.Repeat("Ｙ", maxWidth/2+2) + "￣")
	ct.ResetColor()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: example [message]\n")
		os.Exit(1)
	}
	suddenDeath(os.Args[1])
}
