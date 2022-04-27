package ui

import (
	"fyne.io/fyne/v2"
	"github.com/m00p1ng/learn-go/pixl/apptype"
	"github.com/m00p1ng/learn-go/pixl/pxcanvas"
	"github.com/m00p1ng/learn-go/pixl/swatch"
)

type AppInit struct {
	PixlCanvas *pxcanvas.PxCanvas
	PixlWindow fyne.Window
	State      *apptype.State
	Swatches   []*swatch.Swatch
}
