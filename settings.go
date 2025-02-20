package main

import (
	_ "embed"
	"encoding/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type appSettings struct {
	Mappings map[string][]settingGroup `json:"mappings"`
}

type settingGroup struct {
	Source string
	Target string
	Name   *string
}

//go:embed default-settings.json
var defaultSettings string

func openSettingsWindow(a fyne.App) {
	prefs := a.Preferences()
	// Create a new window for settings
	w := a.NewWindow("Settings")

	// A multi-line text entry to display or edit the JSON
	jsonEntry := widget.NewMultiLineEntry()
	settings := prefs.String("settings")
	if settings == "" {
		settings = defaultSettings
	}

	jsonEntry.Text = settings

	// A button to parse the JSON from the text entry and save/use it
	saveButton := widget.NewButton("Save", func() {
		rawJSON := jsonEntry.Text

		var config appSettings
		if err := json.Unmarshal([]byte(rawJSON), &config); err != nil {
			dialog.ShowError(err, w)
			return
		}

		prefs.SetString("settings", rawJSON)
		w.Close()
	})

	// Lay out our widgets in a vertical box,
	// with load/save buttons in a horizontal box below
	content := container.NewVBox(
		widget.NewLabel("Enter the JSON configuration:"),
		jsonEntry,
		container.NewHBox(saveButton),
		widget.NewLabel("Changes will be visible after application restart"),
	)

	w.SetContent(content)
	jsonEntry.Resize(fyne.Size{
		Width:  400,
		Height: 600,
	})
	w.Show()
}
