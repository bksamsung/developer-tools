package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func bindProxy(ctx context.Context, port int, host string) <-chan error {
	ch := make(chan error)

	go func() {
		
		c := fmt.Sprintf("-D %d -q -N -C %s", port, host)
		cmd := exec.CommandContext(ctx, "ssh", strings.Split(c, " ")...)

		if err := cmd.Run(); err != nil {
			log.Print(err)
			ch <- err
			return
		}

		ch <- nil
	}()

	return ch
}

func bind(ctx context.Context, port int) <-chan error {
	ch := make(chan error)

	go func() {
		c := fmt.Sprintf("-L %d:127.0.0.1:%d work-pc -N", port, port)
		cmd := exec.CommandContext(ctx, "ssh", strings.Split(c, " ")...)

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
	a := app.New()
	w := a.NewWindow("Developer tools")

	ctx, cancel := context.WithCancel(context.Background())

	a.Lifecycle().SetOnStopped(cancel)
	w.CenterOnScreen()

	hello := widget.NewLabel("Select ports you want to bind")

	allContents := []fyne.CanvasObject{hello, bindProxyWidget(ctx)}

	for _, port := range []int{3000, 3001, 3002, 3003, 3004, 3005, 16686 /*jaeger*/} {
		btn := widget.NewButton("Bind it!", func() {})
		errLabel := widget.NewLabel("")

		btnDisconnect := widget.NewButton("Disconnect", func() {})
		btnDisconnect.Disable()

		btn.OnTapped = func() {
			ctx, cancel := context.WithCancel(ctx)
			btnDisconnect.OnTapped = cancel

			btn.Disable()
			btnDisconnect.Enable()
			ch := bind(ctx, port)

			go func(ch <-chan error, btn *widget.Button) {
				err := <-ch

				btn.Enable()
				btnDisconnect.Disable()
				errLabel.Text = ""
				if err != nil {
					errLabel.Text = err.Error()
				}
			}(ch, btn)

		}
		cont := container.NewHBox(widget.NewLabel(fmt.Sprintf("Port %d", port)), btn, btnDisconnect, errLabel)
		allContents = append(allContents, cont)
	}

	w.SetContent(container.NewVBox(
		allContents...,
	))

	w.ShowAndRun()
}

func bindProxyWidget(ctx context.Context) fyne.CanvasObject {
	port := 1080

	btn := widget.NewButton("Bind proxy (port 1080)!", func() {})
	errLabel := widget.NewLabel("")

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
			err := <-ch

			btn.Enable()
			btnDisconnect.Disable()
			errLabel.Text = ""
			if err != nil {
				errLabel.Text = err.Error()
			}
		}(ch, btn)

	}
	return container.NewHBox(widget.NewLabel(fmt.Sprintf("Port %d", port)), host, btn, btnDisconnect, errLabel)
}
