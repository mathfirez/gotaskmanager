package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Task struct {
	name        string
	description string
	text        string
	status      string
	lastChanged string
}

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("Task Manager")

	toggleEdit := false
	toggleNew := false
	statuses := []string{"Open", "Closed"}

	// adding Edit shortcut
	ctrlE := &desktop.CustomShortcut{KeyName: fyne.KeyE, Modifier: fyne.KeyModifierControl}

	// ading New shorcut
	ctrlW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierControl}

	// adding Save current text shortcut
	ctrlS := &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}

	// adding Delete current task shortcut
	ctrlD := &desktop.CustomShortcut{KeyName: fyne.KeyD, Modifier: fyne.KeyModifierControl}

	// adding Delete current task shortcut
	lightMode := false
	ctrlL := &desktop.CustomShortcut{KeyName: fyne.KeyL, Modifier: fyne.KeyModifierControl}

	ctrlOne := &desktop.CustomShortcut{KeyName: fyne.Key1, Modifier: fyne.KeyModifierControl}
	ctrlTwo := &desktop.CustomShortcut{KeyName: fyne.Key2, Modifier: fyne.KeyModifierControl}
	ctrlThree := &desktop.CustomShortcut{KeyName: fyne.Key3, Modifier: fyne.KeyModifierControl}
	ctrlFour := &desktop.CustomShortcut{KeyName: fyne.Key4, Modifier: fyne.KeyModifierControl}

	//adding to canvas and defining function for the ctrl E shortcut

	file, err := os.Open("tasks.csv")
	if err != nil {
		os.Create("tasks.csv")
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	// Initializing list
	var data = []string{}

	// Intialing selected task to be handled later
	var currTask = new(Task)
	currTask.name = ""
	currTask.description = ""
	currTask.text = ""
	currTask.status = ""
	currTask.lastChanged = ""

	var currId int
	currId = -1

	fmt.Println(records)

	// Passing names to the sidebar
	for i := 0; i < len(records); i++ {
		fmt.Println(records[i][0])
		data = append(data, records[i][0])
	}

	taskList := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i])
		})

	headerContent := "# Select a task or create a new one (Ctrl+W)"
	taskHeader := widget.NewRichTextFromMarkdown(" ")
	taskHeader.ParseMarkdown(headerContent)

	notes := widget.NewMultiLineEntry()
	notes.SetText(currTask.text)

	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		// Deal on how to save the textbox to the file
		// Must change commas in user text to "," or find other way to handle it
		fmt.Println("Pressionado")
		currTask.text = notes.Text
		records[currId][2] = currTask.text
		records[currId][4] = time.Now().String()[0:19]
		fmt.Println(records[currId])

		fileW, e := os.OpenFile("tasks.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

		writer := csv.NewWriter(fileW)
		e = writer.WriteAll(records)
		if e != nil {
			fmt.Println(e)
		}
		fmt.Println(time.Now().String()[0:19])
	})

	bottomContent := container.NewVBox(saveBtn)

	taskList.OnSelected = func(id widget.ListItemID) {
		fmt.Println(id)
		currId = id
		currTask.name = records[currId][0]
		currTask.description = records[currId][1]
		currTask.text = records[currId][2]
		currTask.status = records[currId][3]
		currTask.lastChanged = records[currId][4]
		headerContent = "# " + currTask.name + "\n## " + currTask.description + "\n ### Status: " + currTask.status + " - Changed on: " + currTask.lastChanged
		fmt.Println(headerContent)
		taskHeader.ParseMarkdown(headerContent)
		// Setting the text in the editor for the current task
		notes.SetText(currTask.text)
		content := container.NewBorder(taskHeader, bottomContent, taskList, nil, notes)
		w.SetContent(content)
	}

	// Setting edit window
	editName := widget.NewEntry()
	editDescription := widget.NewEntry()
	editStatus := widget.NewSelect(statuses, func(s string) {
		fmt.Println(s)
	})

	saveEditBtn := widget.NewButtonWithIcon("Save Task", theme.DocumentSaveIcon(), func() {
		// Deal on how to save the textbox to the file
		// Must change commas in user text to "," or find other way to handle it
		fmt.Println("Pressionado")
		//editName.Text = currTask.name
		//editDescription.Text = currTask.description
		//editStatus.Selected = currTask.status
		records[currId][0] = editName.Text
		records[currId][1] = editDescription.Text
		records[currId][3] = editStatus.Selected
		records[currId][4] = time.Now().String()[0:19]
		fmt.Println(records[currId])

		fileW, e := os.OpenFile("tasks.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

		writer := csv.NewWriter(fileW)
		e = writer.WriteAll(records)
		if e != nil {
			fmt.Println(e)
		}

		// Redefining the sidebar so it can pickup name changes
		file, err = os.Open("tasks.csv")
		if err != nil {
			os.Create("tasks.csv")
		}
		// Clearing data array so it can be populated again
		data = []string{}
		reader = csv.NewReader(file)
		records, _ = reader.ReadAll()
		for i := 0; i < len(records); i++ {
			fmt.Println(records[i][0])
			data = append(data, records[i][0])
		}

		// Now must refresh top header and go back to main window
		currTask.name = records[currId][0]
		currTask.description = records[currId][1]
		currTask.status = records[currId][3]
		currTask.lastChanged = records[currId][4]
		headerContent = "# " + currTask.name + "\n## " + currTask.description + "\n ### Status: " + currTask.status + " - Changed on: " + currTask.lastChanged
		taskHeader.ParseMarkdown(headerContent)
		taskHeader.Refresh()
		taskList.Refresh()
		content := container.NewBorder(taskHeader, bottomContent, taskList, nil, notes)
		w.SetContent(content)
		toggleEdit = false

	})

	editContent := container.NewVBox(widget.NewLabel("Task Name"), editName, widget.NewLabel("Task Description"), editDescription, widget.NewLabel("Task Status"), editStatus, saveEditBtn)

	w.Canvas().AddShortcut(ctrlE, func(shortcut fyne.Shortcut) {
		// Flip between true and false to toggle edit mdoe or not
		fmt.Println("Ctrl+E")
		if currId == -1 {
			dialog.ShowInformation("Error", "Select a task before editing.", w)
			return
		}
		toggleEdit = !toggleEdit
		if toggleEdit {
			toggleNew = false
			editName.Text = currTask.name
			editDescription.Text = currTask.description
			editStatus.Selected = currTask.status
			editContent = container.NewVBox(widget.NewLabel("Task Name"), editName, widget.NewLabel("Task Description"), editDescription, widget.NewLabel("Task Status"), editStatus, saveEditBtn)
			editContent.Refresh()
			content := container.NewBorder(taskHeader, bottomContent, taskList, nil, editContent)
			w.SetContent(content)
		} else {
			content := container.NewBorder(taskHeader, bottomContent, taskList, nil, notes)
			w.SetContent(content)
		}
	})

	// Setting New task window
	newName := widget.NewEntry()
	newDescription := widget.NewEntry()
	newStatus := widget.NewSelect(statuses, func(s string) {
		fmt.Println(s)
	})

	saveNewBtn := widget.NewButtonWithIcon("Add New Task", theme.DocumentSaveIcon(), func() {
		// Deal on how to save the textbox to the file
		// Must change commas in user text to "," or find other way to handle it
		fmt.Println("Novo - Pressionado")

		// SELECT A NEW ID MAX + 1 AND APPEND TO THE RECORDS SLICE
		if newName.Text == "" {
			newName.Text = "(new task)"
		}

		if newStatus.Selected == "" {
			newStatus.Selected = statuses[0]
		}

		newRecord := []string{newName.Text, newDescription.Text, "Insert notes here...", newStatus.Selected, time.Now().String()[0:19]}

		records = append(records, newRecord)

		fileW, e := os.OpenFile("tasks.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

		writer := csv.NewWriter(fileW)
		e = writer.WriteAll(records)
		if e != nil {
			fmt.Println(e)
		}
		// Redefining the sidebar
		file, err = os.Open("tasks.csv")
		if err != nil {
			os.Create("tasks.csv")
		}
		// Clearing data array so it can be populated again
		data = []string{}
		reader = csv.NewReader(file)
		records, _ = reader.ReadAll()
		for i := 0; i < len(records); i++ {
			fmt.Println(records[i][0])
			data = append(data, records[i][0])
		}
		// Now must refresh the sidebar and go back to main window

		taskList.Refresh()
		content := container.NewBorder(taskHeader, bottomContent, taskList, nil, notes)
		content.Refresh()
		w.SetContent(content)
		toggleNew = false

	})

	newContent := container.NewVBox(widget.NewLabel("New Task Name"), newName, widget.NewLabel("New Task Description"), newDescription, widget.NewLabel("New Task Status"), newStatus, saveNewBtn)

	w.Canvas().AddShortcut(ctrlW, func(shortcut fyne.Shortcut) {
		// Flip between true and false to toggle edit mdoe or not
		fmt.Println("Ctrl+W")
		toggleNew = !toggleNew
		if toggleNew {
			toggleEdit = false
			content := container.NewBorder(taskHeader, bottomContent, taskList, nil, newContent)
			w.SetContent(content)
		} else {
			content := container.NewBorder(taskHeader, bottomContent, taskList, nil, notes)
			w.SetContent(content)
		}
	})

	w.Canvas().AddShortcut(ctrlS, func(shortcut fyne.Shortcut) {
		// Flip between true and false to toggle mode or not
		fmt.Println("Ctrl+S")
		// Deal on how to save the textbox to the file
		// Must change commas in user text to "," or find other way to handle it
		fmt.Println("Pressionado")
		currTask.text = notes.Text
		records[currId][2] = currTask.text
		records[currId][4] = time.Now().String()[0:19]
		fmt.Println(records[currId])

		fileW, e := os.OpenFile("tasks.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

		writer := csv.NewWriter(fileW)
		e = writer.WriteAll(records)
		if e != nil {
			fmt.Println(e)
		}

		// Updating Header
		fmt.Println(time.Now().String()[0:19])
		currTask.lastChanged = records[currId][4]
		headerContent = "# " + currTask.name + "\n## " + currTask.description + "\n ### Status: " + currTask.status + " - Changed on: " + currTask.lastChanged
		taskHeader.ParseMarkdown(headerContent)
		taskHeader.Refresh()
		content := container.NewBorder(taskHeader, bottomContent, taskList, nil, notes)
		w.SetContent(content)
	})

	w.Canvas().AddShortcut(ctrlD, func(shortcut fyne.Shortcut) {
		fmt.Println("Ctrl+D")

		if toggleEdit == true || toggleNew == true {
			return
		}

		if currId == -1 {
			return
		}
		records = append(records[:currId], records[currId+1:]...)

		fileW, e := os.OpenFile("tasks.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

		writer := csv.NewWriter(fileW)
		e = writer.WriteAll(records)
		if e != nil {
			fmt.Println(e)
		}
		// Redefining the sidebar
		file, err = os.Open("tasks.csv")
		if err != nil {
			os.Create("tasks.csv")
		}
		// Clearing data array so it can be populated again
		data = []string{}
		reader = csv.NewReader(file)
		records, _ = reader.ReadAll()
		for i := 0; i < len(records); i++ {
			fmt.Println(records[i][0])
			data = append(data, records[i][0])
		}

		// Reseting list
		currId = -1
		taskList.UnselectAll()

		// Now must refresh the sidebar and go back to main window
		taskList.Refresh()
		headerContent = "# Select a task or create a new one (Ctrl+W)"
		taskHeader.ParseMarkdown(headerContent)
		taskHeader.Refresh()
		content := container.NewBorder(taskHeader, bottomContent, taskList, nil, notes)
		content.Refresh()
		w.SetContent(content)
	})

	w.Canvas().AddShortcut(ctrlL, func(shortcut fyne.Shortcut) {
		fmt.Println("Ctrl+L")
		if lightMode == true {
			a.Settings().SetTheme(theme.DarkTheme())
			lightMode = false
			return
		}
		a.Settings().SetTheme(theme.LightTheme())
		lightMode = true
	})

	// Numeric shortcuts
	// Ctrl + 1
	w.Canvas().AddShortcut(ctrlOne, func(shortcut fyne.Shortcut) {
		fmt.Println("Ctrl+1")
		taskList.Select(0)
	})

	// Ctrl + 2
	w.Canvas().AddShortcut(ctrlTwo, func(shortcut fyne.Shortcut) {
		fmt.Println("Ctrl+2")
		taskList.Select(1)
	})

	// Ctrl + 3
	w.Canvas().AddShortcut(ctrlThree, func(shortcut fyne.Shortcut) {
		fmt.Println("Ctrl+3")
		taskList.Select(2)
	})

	// Ctrl + 4
	w.Canvas().AddShortcut(ctrlFour, func(shortcut fyne.Shortcut) {
		fmt.Println("Ctrl+4")
		taskList.Select(3)
	})

	content := container.NewBorder(taskHeader, bottomContent, taskList, nil, notes)

	w.SetContent(content)

	w.CenterOnScreen()
	w.Resize(fyne.NewSize(1280, 720))
	w.ShowAndRun()
}
