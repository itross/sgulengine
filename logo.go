package sgulengine

import "github.com/common-nighthawk/go-figure"

// PrintLogo prints out the SgulENGINE logo.
func PrintLogo() {
	myFigure := figure.NewFigure("SgulENGINE", "doom", true).Scroll(10000, 200, "left")
	myFigure.Print()
}
