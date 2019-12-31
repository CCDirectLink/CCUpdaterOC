package main

import (
	"github.com/CCDirectLink/CCUpdaterUI/design"
	"github.com/CCDirectLink/CCUpdaterUI/frenyard"
	"github.com/CCDirectLink/CCUpdaterUI/frenyard/framework"
	"github.com/CCDirectLink/CCUpdaterUI/frenyard/integration"
)

func (app *upApplication) ShowWaiter(text string, a func(func(string)), b func()) {
	label := framework.NewUILabelPtr(integration.NewTextTypeChunk("", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{})
	app.Teleport(design.LayoutDocument(design.Header{
		Title: text,
	}, label, false))
	go func () {
		a(func (text string) {
			app.upQueued <- func () {
				label.SetText(integration.NewTextTypeChunk(text, design.GlobalFont))
			}
		})
		app.upQueued <- b
	}()
}

func (app *upApplication) MessageBox(title string, text string, b func()) {
	app.Teleport(design.LayoutDocument(design.Header{
		Title: title,
	}, design.LayoutMsgbox(text, b), true))
}
