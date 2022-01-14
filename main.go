package main

/*
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

unsigned short *get_screen_size(void)
{
    static unsigned short size[2];
    char *array[8];
    char screen_size[64];
    char* token = NULL;

    FILE *cmd = popen("xdpyinfo | awk '/dimensions/ {print $2}'", "r");

    if (!cmd)
        return 0;

    while (fgets(screen_size, sizeof(screen_size), cmd) != NULL);
    pclose(cmd);

    token = strtok(screen_size, "x\n");

    if (!token)
        return 0;

    for (unsigned short i = 0; token != NULL; ++i) {
        array[i] = token;
        token = strtok(NULL, "x\n");
    }
    size[0] = atoi(array[0]);
    size[1] = atoi(array[1]);
    size[2] = -1;

    return size;
}


int width()
{
    unsigned short *size = get_screen_size();

    return size[0];
}

int height(){
	unsigned short *size = get_screen_size();

	return size[1];
}



*/
import "C"

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"fyne.io/fyne/driver/mobile"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type numericalEntry struct {
	widget.Entry
}

func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *numericalEntry) TypedRune(r rune) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', ',':
		e.Entry.TypedRune(r)
	}
}

func (e *numericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *numericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func form(w fyne.Window) fyne.Widget {
	entry := widget.NewEntry()
	textArea := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Si", Widget: entry}},
		SubmitText: "Siguiente",
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", entry.Text)
			log.Println("multiline:", textArea.Text)
			w.Close()
		},
	}

	// we can also append items
	form.Append("Description", textArea)

	provinceSelect := widget.NewSelect([]string{"anhui", "zhejiang", "shanghai"}, func(value string) {
		fmt.Println("province:", value)
	})
	provinceBox := container.NewVBox(widget.NewLabel("Province"), layout.NewSpacer(), provinceSelect)

	form.Append("hola", provinceBox)

	grid := container.NewGridWithColumns(1)
	grid.Add(widget.NewCard(
		"This is my title",
		"prueba 2",
		nil,
	))

	gridContainer := container.NewGridWithColumns(3)

	for i := 0; i < 9; i++ {

		item1 := widget.NewAccordionItem("A",
			container.NewVBox(
				widget.NewLabel("A for Apple A for Apple A for Apple A for Apple"),
				widget.NewLabel("A for Apple A for Apple A for Apple A for Apple")))
		ac := widget.NewAccordion(item1)
		top := ac
		middle := newNumericalEntry()
		left := widget.NewButton("Borrar", func() {
			log.Println("tapped")
		})
		right := widget.NewButton("Agregar", nil)

		buttonTitle := "Disable"

		changeButton := func() {
			// here could be your logic
			// how to disable/enable button
			if right.Text == "Disable" {
				buttonTitle = "Enable"
				//button.Disable()
			}
			right.SetText(buttonTitle)
			right.Refresh()
		}
		right.OnTapped = changeButton

		content := container.New(layout.NewBorderLayout(top, nil, left, right),
			top, left, right, middle)
		gridContainer.Add(content)
	}

	form.Append("Items", gridContainer)
	return form
}

func scroll(w fyne.Window) fyne.Widget {
	return container.NewVScroll(form(w))
}

func tabs(w fyne.Window) fyne.Widget {

	text1 := canvas.NewText("Hello", color.White)
	text2 := canvas.NewText("There", color.White)
	text3 := canvas.NewText("(right)", color.White)
	content := container.New(layout.NewHBoxLayout(), text1, text2, layout.NewSpacer(), text3)

	text4 := canvas.NewText("centered", color.White)
	centered := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), text4, layout.NewSpacer())

	tabs := container.NewAppTabs(
		container.NewTabItem("Factura nueva", scroll(w)),
		container.NewTabItem("Tab 1", container.New(layout.NewVBoxLayout(), content, centered)),
		container.NewTabItem("Tab 1", widget.NewLabel("Hello")),
		container.NewTabItem("Tab 1", widget.NewLabel("Hello")),
		container.NewTabItem("Tab 1", widget.NewLabel("Hello")),
		container.NewTabItem("Tab 1", widget.NewLabel("Hello")),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	return tabs
}

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Title")

	w.SetContent(tabs(w))
	w.Resize(fyne.NewSize(float32(C.width()), float32(C.height())))

	w.CenterOnScreen()
	w.ShowAndRun()
}
