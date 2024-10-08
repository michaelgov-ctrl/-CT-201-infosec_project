package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type customTheme struct{}

var _ fyne.Theme = (*customTheme)(nil)

func (ct customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 204, G: 0, B: 0, A: 160}
	case theme.ColorNameInputBackground:
		return color.White
	case theme.ColorNameForeground:
		return color.Black
	case theme.ColorNameInputBorder:
		return color.Black
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (ct customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (ct customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (ct customTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 50
	case theme.SizeNameInputBorder:
		return 5
	default:
		return theme.DefaultTheme().Size(name)
	}
}

func newPosGui(pos posApplication) {
	app := app.New()
	app.Settings().SetTheme(&customTheme{})
	window := app.NewWindow("Victim Point-of-Sale System")
	window.SetCloseIntercept(func() {
		pos.cardReader.Close()
		window.Close()
	})

	textBox := widget.NewEntry()
	window.SetContent(textBox)

	go func() {
		for {
			ccInfo, err := pos.readCard()
			if err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					textBox.Text = fmt.Sprintf("card read error: %s", err)
					textBox.Refresh()
				}
				continue
			}

			name, err := pos.parseCardHolderName(ccInfo)
			if err != nil {
				textBox.Text = fmt.Sprintf("card read error: %s", err)
				log.Printf("card read error: %s", err)
				textBox.Refresh()
			}

			textBox.Text = fmt.Sprintf("Purchase processed for:\n%s", name)
			textBox.Refresh()
			fmt.Println(ccInfo)
		}
	}()

	window.Resize(fyne.NewSize(700, 500))
	window.ShowAndRun()
}
