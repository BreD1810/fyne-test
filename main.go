package main

import (
	"encoding/json"
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const DataFile = "data.json"

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("List")

	loadedData := loadData()

	data := binding.NewStringList()
	data.Set(loadedData)

	defer saveData(data)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject { // New item
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) { // Update item
			o.(*widget.Label).Bind(i.(binding.String))
		})

	list.OnSelected = func(id widget.ListItemID) { // Update/delete items
		list.Unselect(id)
		d, _ := data.GetValue(id)
		w := myApp.NewWindow("Edit Data")

		itemName := widget.NewEntry()
		itemName.Text = d

		updateDataButton := widget.NewButton("Update", func() {
			data.SetValue(id, itemName.Text)
			w.Close()
		})

		cancelButton := widget.NewButton("Cancel", func() {
			w.Close()
		})

		deleteDataButton := widget.NewButton("Delete", func() {
			newData := []string{}
			oldData, _ := data.Get()
			for i, item := range oldData {
				if i != id {
					newData = append(newData, item)
				}
			}
			data.Set(newData)
			w.Close()
		})

		w.SetContent(container.New(layout.NewVBoxLayout(), itemName, updateDataButton, deleteDataButton, cancelButton))
		w.Resize(fyne.NewSize(400, 600))
		w.CenterOnScreen()
		w.Show()
	}

	addButton := widget.NewButton("Add", func() {
		w := myApp.NewWindow("Add Data")
		itemName := widget.NewEntry()

		addDataButton := widget.NewButton("Add", func() {
			data.Append(itemName.Text)
			w.Close()
		})

		cancelButton := widget.NewButton("Cancel", func() {
			w.Close()
		})

		w.SetContent(container.New(layout.NewVBoxLayout(), itemName, addDataButton, cancelButton))
		w.Resize(fyne.NewSize(400, 200))
		w.CenterOnScreen()
		w.Show()
	})

	exitButton := widget.NewButton("Quit", func() {
		myWindow.Close()
	})

	myWindow.SetContent(container.NewBorder(nil, container.New(layout.NewVBoxLayout(), addButton, exitButton), nil, nil, list))
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.SetMaster() // Should stop the app if this is closed
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun() // Runs the application
}

// loadData loads items from data.json
func loadData() []string {
	input, _ := ioutil.ReadFile(DataFile)
	data := []string{}
	json.Unmarshal(input, &data)
	return data
}

// saveData saves data to data.json
func saveData(data binding.StringList) {
	d, _ := data.Get()
	jsonData, _ := json.Marshal(d)
	ioutil.WriteFile(DataFile, jsonData, 0644)
}
