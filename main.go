package main

import "C"

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"

	"database/sql"

	"fyne.io/fyne/driver/mobile"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-sql-driver/mysql"
)

type products struct {
	id          int
	name        string
	description string
	price       float32
}

func existeError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

var path = "C:/Users/dany0/OneDrive/Documentos/img"

func readDatabase() {

	//Verifica que el archivo existe
	var _, err = os.Stat(path)
	//Crea el archivo si no existe
	if os.IsNotExist(err) {
		var file, err = os.Create(path + "/prueba.json")
		if existeError(err) {
			return
		}
		defer file.Close()
	}
	fmt.Println("File Created Successfully", path)

}

func crearArchivo(date string, nameFile string) {
	f, err := os.Create(path + nameFile + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("old falcon\n" + date)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("Hecho datos")

}

type Databases struct {
	Databases []Database `json:"databases"`
}

type Database struct {
	Host              string `json:"host"`
	Usuario           string `json:"user"`
	Pass              string `json:"password"`
	NombreBaseDeDatos string `json:"nameDatabase"`
}

type Products struct {
	Id                  int
	Nombre, Description string
	Price               float32
}

func obtenerContactos() ([]Products, error) {
	contactos := []Products{}
	db, err := obtenerBaseDeDatos()

	if err != nil {
		return nil, err
	}

	defer db.Close()
	filas, err := db.Query("select idproducts, name, DESCRIPTION, price FROM local.products")

	if err != nil {
		return nil, err
	}

	defer filas.Close()

	var c Products

	for filas.Next() {
		err = filas.Scan(&c.Id, &c.Nombre, &c.Description, &c.Price)
		if err != nil {
			return nil, err
		}
		contactos = append(contactos, c)
	}
	return contactos, nil
}

type numericalEntry struct{ widget.Entry }

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

func obtenerBaseDeDatos() (db *sql.DB, e error) {

	readDatabase()

	var operatingSystem string = "C:/Users/dany0/OneDrive/Documentos/img/prueba.json"

	/*

		operatingSystem = "C:/Users/dany0/OneDrive/Documentos/img/prueba.json"
		if runtime.GOOS == "windows" {
			operatingSystem = "C:/Users/dany0/OneDrive/Documentos/img"
		} else {
			operatingSystem = "/home/daniel/Documentos/data/prueba.json"
		}

	*/

	jsonFile, _ := os.Open(operatingSystem)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var databases Databases

	json.Unmarshal(byteValue, &databases)

	var host, usuario, pass, nombreBaseDeDatos string

	host = databases.Databases[0].Host
	usuario = databases.Databases[0].Usuario
	pass = databases.Databases[0].Pass
	nombreBaseDeDatos = databases.Databases[0].NombreBaseDeDatos

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, nombreBaseDeDatos))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func gridProduct() *fyne.Container {
	contactos, _ := obtenerContactos()

	gridContainer := container.NewGridWithColumns(2)

	for i := 0; i < 7; i++ {
		grid := container.NewGridWithColumns(
			1,
			widget.NewCard(
				contactos[i].Nombre,
				"Tama: ",
				nil,
			),
		)

		top := grid
		middle := newNumericalEntry()
		middle.SetText("si")
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
	return gridContainer
}

type optionCar struct {
	Seccion string
	Lavado  string
}

func newOptionCar(seccion string) *optionCar {
	option := optionCar{Seccion: seccion}
	return &option
}

func typeCar() fyne.Widget {
	return widget.NewSelect([]string{"Moto", "Taxi chico", "Taxi Grande", "Colegial", "Seda", "SUV", "Pick Up", "Busito"}, func(value string) {
		seccion := newOptionCar(value)
		fmt.Println("seccion seleccion:", seccion.Seccion)
		crearArchivo(value, "/dataTypeCar")
	})
}
func LavadoSelect() fyne.Widget {
	LavadoSelect := widget.NewSelect([]string{"Espuma", "Sin espuma"}, func(value string) {
		fmt.Println("seccion seleccion:", value)
		crearArchivo(value, "/dataLavado")
	})
	return LavadoSelect
}

func containerOption() *fyne.Container {
	seccionBox := container.NewVBox(widget.NewLabel("Sección"), layout.NewSpacer(), typeCar())

	LavadoBox := container.NewVBox(widget.NewLabel("Lavado"), layout.NewSpacer(), LavadoSelect())
	return container.NewHBox(LavadoBox, seccionBox)
}

func form(w fyne.Window) fyne.Widget {
	entryPlaca := widget.NewEntry()
	entryModel := widget.NewEntry()
	textArea := widget.NewMultiLineEntry()

	return &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Placa", Widget: entryPlaca},
			{Text: "Modelo", Widget: entryModel},
			{Text: "Detalles", Widget: textArea},
			{Text: "Opciones", Widget: containerOption()},
			{Text: "Accesorio", Widget: gridProduct()},
		},
		SubmitText: "Siguiente",
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", entryPlaca.Text)
			log.Println("multiline:", textArea.Text)
			w.Close()
		},
	}
}

func scroll(w fyne.Window) fyne.Widget {
	return container.NewVScroll(form(w))
}

func tabsSecond(w fyne.Window) fyne.Widget {
	return container.NewAppTabs(
		container.NewTabItem("Factura nueva", scroll(w)),
		container.NewTabItem("Tab 1", widget.NewLabel("Hello 1")),
	)
}

func tabs(w fyne.Window) fyne.Widget {

	text1 := canvas.NewText("Hello", color.White)
	text2 := canvas.NewText("There", color.White)
	text3 := canvas.NewText("(right)", color.White)
	content := container.New(layout.NewHBoxLayout(), text1, text2, layout.NewSpacer(), text3)

	text4 := canvas.NewText("centered", color.White)
	centered := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), text4, layout.NewSpacer())

	tabs := container.NewAppTabs(
		container.NewTabItem("Factura nueva", tabsSecond(w)),
		container.NewTabItem("Tab 1", container.New(layout.NewVBoxLayout(), content, centered)),
		container.NewTabItem("Tab 1", widget.NewLabel("Hello 2")),
		container.NewTabItem("Tab 1", widget.NewLabel("Hello 3")),
		container.NewTabItem("Cerrar", widget.NewButton("Cerrar aplicación", func() { w.Close() })),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	return tabs
}

func main() {

	print(runtime.GOOS)
	myApp := app.New()
	w := myApp.NewWindow("Title")

	w.SetContent(tabs(w))
	//w.Resize(fyne.NewSize(float32(920), float32(620)))

	w.SetFullScreen(true)

	w.CenterOnScreen()
	w.ShowAndRun()

}
