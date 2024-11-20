package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"

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

func bind(ctx context.Context, port int) <-chan error {
	ch := make(chan error)

	go func() {
		c := fmt.Sprintf("-L %d:127.0.0.1:%d work-pc -N", port, port)
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
	a := app.New()
	w := a.NewWindow("Developer tools")

	ctx, cancel := context.WithCancel(context.Background())

	a.Lifecycle().SetOnStopped(cancel)
	w.CenterOnScreen()

	allContents := []fyne.CanvasObject{}

	accordions := widget.NewAccordion()
	for _, port := range []int{3000, 3001, 3002, 3003, 3004, 3005, 3006} {
		btn := widget.NewButton(fmt.Sprintf("Bind %d!", port), func() {})

		btnDisconnect := widget.NewButton("Disconnect", func() {})
		btnDisconnect.Disable()

		btn.OnTapped = func() {
			ctx, cancel := context.WithCancel(ctx)
			btnDisconnect.OnTapped = cancel

			btn.Disable()
			btnDisconnect.Enable()
			ch := bind(ctx, port)

			go func(ch <-chan error, btn *widget.Button) {
				<-ch

				btn.Enable()
				btnDisconnect.Disable()
			}(ch, btn)

		}
		cont := container.NewHBox(btn, btnDisconnect)
		allContents = append(allContents, cont)
	}
	accordions.Append(widget.NewAccordionItem("trader", container.NewGridWithColumns(2, allContents...)))

	allContents = []fyne.CanvasObject{}
	for _, port := range []int{9100, 9090, 16686 /*jaeger*/} {
		btn := widget.NewButton(fmt.Sprintf("Bind %d!", port), func() {})

		btnDisconnect := widget.NewButton("Disconnect", func() {})
		btnDisconnect.Disable()

		btn.OnTapped = func() {
			ctx, cancel := context.WithCancel(ctx)
			btnDisconnect.OnTapped = cancel

			btn.Disable()
			btnDisconnect.Enable()
			ch := bind(ctx, port)

			go func(ch <-chan error, btn *widget.Button) {
				<-ch

				btn.Enable()
				btnDisconnect.Disable()
			}(ch, btn)

		}
		cont := container.NewHBox(btn, btnDisconnect)
		allContents = append(allContents, cont)
	}
	accordions.Append(widget.NewAccordionItem("observability", container.NewGridWithColumns(2, allContents...)))

	allContents = []fyne.CanvasObject{}
	for _, port := range []int{8085, 8095} {
		btn := widget.NewButton(fmt.Sprintf("Bind %d!", port), func() {})

		btnDisconnect := widget.NewButton("Disconnect", func() {})
		btnDisconnect.Disable()

		btn.OnTapped = func() {
			ctx, cancel := context.WithCancel(ctx)
			btnDisconnect.OnTapped = cancel

			btn.Disable()
			btnDisconnect.Enable()
			ch := bind(ctx, port)

			go func(ch <-chan error, btn *widget.Button) {
				<-ch

				btn.Enable()
				btnDisconnect.Disable()
			}(ch, btn)

		}
		cont := container.NewHBox(btn, btnDisconnect)
		allContents = append(allContents, cont)
	}
	accordions.Append(widget.NewAccordionItem("bidder", container.NewGridWithColumns(2, allContents...)))

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
