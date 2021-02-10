package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/unix-streamdeck/api"
)

type button struct {
	widget.BaseWidget
	editor *editor

	keyID int
	key   api.Key
}

func newButton(key api.Key, id int, e *editor) *button {
	b := &button{key: key, keyID: id, editor: e}
	b.ExtendBaseWidget(b)
	return b
}

func (b *button) CreateRenderer() fyne.WidgetRenderer {
	icon := canvas.NewImageFromFile(b.key.Icon)
	text := &canvas.Image{}

	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 2
	border.SetMinSize(fyne.NewSize(float32(b.editor.iconSize), float32(b.editor.iconSize)))

	bg := canvas.NewRectangle(color.Black)
	render := &buttonRenderer{border: border, text: text, icon: icon, bg: bg,
		objects: []fyne.CanvasObject{bg, icon, text, border}, b: b}
	render.Refresh()
	return render
}

func (b *button) Tapped(ev *fyne.PointEvent) {
	b.editor.editButton(b)
}

func (b *button) updateKey() {
	b.editor.config.Pages[b.editor.currentPage][b.keyID] = b.key
	if b.editor.config.Pages[b.editor.currentPage][b.keyID].IconHandler == "Default" {
		b.editor.config.Pages[b.editor.currentPage][b.keyID].IconHandler = ""
	}
	if b.editor.config.Pages[b.editor.currentPage][b.keyID].KeyHandler == "Default" {
		b.editor.config.Pages[b.editor.currentPage][b.keyID].KeyHandler = ""
	}
}

const (
	buttonInset = 5
)

type buttonRenderer struct {
	border, bg *canvas.Rectangle
	icon, text *canvas.Image

	objects []fyne.CanvasObject

	b *button
}

func (r *buttonRenderer) Layout(s fyne.Size) {
	size := s.Subtract(fyne.NewSize(buttonInset*2, buttonInset*2))
	offset := fyne.NewPos(buttonInset, buttonInset)

	for _, obj := range r.objects {
		obj.Move(offset)
		obj.Resize(size)
	}
}

func (r *buttonRenderer) MinSize() fyne.Size {
	iconSize := fyne.NewSize(r.b.editor.iconSize, r.b.editor.iconSize)
	return iconSize.Add(fyne.NewSize(buttonInset*2, buttonInset*2))
}

func (r *buttonRenderer) Refresh() {
	if r.b.editor.currentButton == r.b {
		r.border.StrokeColor = theme.FocusColor()
	} else {
		r.border.StrokeColor = &color.Gray{128}
	}

	r.text.Image = r.textToImage()
	r.text.Refresh()
	if r.b.key.Icon != r.icon.File {
		r.icon.File = r.b.key.Icon
		r.icon.Refresh()
	}

	r.border.Refresh()
}

func (r *buttonRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *buttonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *buttonRenderer) Destroy() {
	// nothing
}

func (r *buttonRenderer) textToImage() image.Image {
	textImg := image.NewNRGBA(image.Rect(0, 0, r.b.editor.iconSize, r.b.editor.iconSize))
	img, err := api.DrawText(textImg, r.b.key.Text, r.b.key.TextSize, r.b.key.TextAlignment)
	if err != nil {
		fyne.LogError("Failed to draw text to imge", err)
	}
	return img
}
