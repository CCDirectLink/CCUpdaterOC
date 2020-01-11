package main

import (
	"os"
	"time"

	"github.com/CCDirectLink/CCUpdaterUI/frenyard/framework"
	"github.com/CCDirectLink/ccmu/pkg"
)

type transactionOperation int

const (
	opInstall transactionOperation = iota
	opUninstall
	opUpdate
)

type transaction map[pkg.Package]transactionOperation

// PerformTransaction performs a transaction, showing UI for it as well.
func (app *upApplication) PerformTransaction(back framework.ButtonBehavior, tx transaction) {
	// It begins...
	log := "-- Log started at " + time.Now().Format(time.RFC1123) + " (ccmodupdater.log) --"
	app.ShowWaiter("Working...", func(progress func(string)) {
		for p, op := range tx {
			info, _ := p.Info()

			switch op {
			case opInstall:
				log += "\nInstalling " + info.NiceName + "..."
				progress(log)
				p.Install()
			case opUninstall:
				log += "\nRemoving " + info.NiceName + "..."
				progress(log)
				p.Uninstall()
			case opUpdate:
				log += "\nUpgrading " + info.NiceName + "..."
				progress(log)
				p.Update()
			}
		}
	}, func() {
		cfgFile, err := os.OpenFile("ccmodupdater.log", os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err == nil {
			// Oh well
			cfgFile.WriteString(log + "\n")
		}
		cfgFile.Close()
		app.GSInstant()
		app.MessageBox("Report", log, back)
	})
}
