package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func bindProxy(ctx context.Context, port int, host string) <-chan error {
	ch := make(chan error)

	go func() {

		c := fmt.Sprintf("-D %d -q -N -C %s", port, host)
		cmd := exec.CommandContext(ctx, "ssh", strings.Split(c, " ")...)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		if err := cmd.Run(); err != nil {
			log.Print(err)
			ch <- err
			return
		}

		ch <- nil
	}()

	return ch
}

func bind(ctx context.Context, source, targetPort string) <-chan error {
	ch := make(chan error)

	go func() {
		c := fmt.Sprintf("-L %s:127.0.0.1:%s work-pc -N", targetPort, source)
		fmt.Println(c)
		cmd := exec.CommandContext(ctx, "ssh", strings.Split(c, " ")...)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		if err := cmd.Run(); err != nil {
			log.Print(err)
			ch <- err
			return
		}

		ch <- nil
	}()

	return ch
}

func main() {
	a := app.NewWithID("developer-tools")
	w := a.NewWindow("Developer tools")

	optionsMenu := fyne.NewMenu("Options",
		fyne.NewMenuItem("Settings", func() {
			openSettingsWindow(a)
		}),
	)

	mainMenu := fyne.NewMainMenu(optionsMenu)
	w.SetMainMenu(mainMenu)

	ctx, cancel := context.WithCancel(context.Background())

	a.Lifecycle().SetOnStopped(cancel)
	w.CenterOnScreen()

	settings := a.Preferences().String("settings")
	if settings == "" {
		settings = defaultSettings
	}
	var config appSettings
	if err := json.Unmarshal([]byte(settings), &config); err != nil {
		dialog.ShowError(err, w)
	}

	accordions := widget.NewAccordion()

	for name, mappings := range config.Mappings {
		allContents := []fyne.CanvasObject{}

		for _, mapping := range mappings {
			btnName := fmt.Sprintf("Bind %s!", mapping.Target)
			if mapping.Name != nil {
				btnName = *mapping.Name
			}

			btn := widget.NewButton(btnName, func() {})

			btnDisconnect := widget.NewButton("Disconnect", func() {})
			btnDisconnect.Disable()

			btn.OnTapped = func() {
				ctx, cancel := context.WithCancel(ctx)
				btnDisconnect.OnTapped = cancel

				btn.Disable()
				btnDisconnect.Enable()
				ch := bind(ctx, mapping.Source, mapping.Target)

				go func(ch <-chan error, btn *widget.Button) {
					<-ch

					btn.Enable()
					btnDisconnect.Disable()
				}(ch, btn)

			}
			cont := container.NewHBox(btn, btnDisconnect)

			allContents = append(allContents, cont)
		}
		accordions.Append(widget.NewAccordionItem(name, container.NewGridWithColumns(2, allContents...)))
	}

	w.SetContent(container.NewGridWithColumns(1,
		accordions,
		bindProxyWidget(ctx),
	))

	w.ShowAndRun()
}

func bindProxyWidget(ctx context.Context) fyne.CanvasObject {
	port := 1080

	btn := widget.NewButton("Bind proxy", func() {})

	btnDisconnect := widget.NewButton("Disconnect", func() {})
	btnDisconnect.Disable()

	host := widget.NewEntry()
	host.SetPlaceHolder("work-pc")
	host.Text = "work-pc"

	btn.OnTapped = func() {
		ctx, cancel := context.WithCancel(ctx)
		btnDisconnect.OnTapped = cancel

		btn.Disable()
		btnDisconnect.Enable()
		ch := bindProxy(ctx, port, host.Text)

		go func(ch <-chan error, btn *widget.Button) {
			<-ch

			btn.Enable()
			btnDisconnect.Disable()
		}(ch, btn)

	}
	return container.NewVBox(host, btn, btnDisconnect)
}
