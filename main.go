package main

import (
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"

	"github.com/therecipe/qt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/qml"
	"github.com/therecipe/qt/widgets"
)

const defQSS = `#QLineEdit { background-color: rgb(255, 0, 0) }
#QGroupBox { border: 1px solid white }
QComboBox { color: blue }`

func main() {

	//create qt widgets application
	qApp := widgets.NewQApplication(len(os.Args), os.Args)
	switch runtime.GOARCH {
	case "js", "wasm":
		if qt.Global.Call("eval", "window.location.hash.search('windows')").Int() != -1 {
			widgets.QApplication_SetStyle2("windows")
		}
	}

	//window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Widgets Playground")
	window.SetCentralWidget(widgets.NewQWidget(nil, 0))
	windowLayout := widgets.NewQHBoxLayout2(window.CentralWidget())

	//left layout
	lVBox := widgets.NewQVBoxLayout()

	//tab widget setup
	tabWidget := widgets.NewQTabWidget(nil)
	lVBox.AddWidget(tabWidget, 0, 0)

	switch runtime.GOARCH {
	case "js", "wasm":
		widgets.QApplication_Desktop().ConnectResized(func(screen int) {
			if widgets.QApplication_Desktop().Screen(screen).Width() < 1024 ||
				widgets.QApplication_Desktop().Screen(screen).Height() < 768 {

				tabWidget.SetStyleSheet("QTabBar::scroller { width: 105px }")
				tabWidget.SetTabPosition(widgets.QTabWidget__West)
			} else {
				tabWidget.SetStyleSheet("")
				tabWidget.SetTabPosition(widgets.QTabWidget__North)
			}
		})
		widgets.QApplication_Desktop().Resized(0)
	}

	//create textedits
	var textEdits []*widgets.QTextEdit
	var textEditQSS *widgets.QPlainTextEdit

	for i := 0; i <= 9; i++ {
		te := newTextEdit(fileForIndex(i))
		te.SetReadOnly(i >= 4)
		textEdits = append(textEdits, te)

		if i == 0 {
			cWidget := widgets.NewQWidget(nil, 0)
			vBox := widgets.NewQVBoxLayout2(cWidget)

			vBox.AddWidget(te, 5, 0)

			textEditQSS = newTextEditQSS()
			vBox.AddWidget(textEditQSS, 1, 0)

			tabWidget.AddTab(cWidget, fileForIndex(i))
		} else {
			tabWidget.AddTab(te, fileForIndex(i))
		}
	}

	//add info tab
	tabWidget.AddTab(newInfoWidget(), "info")

	//left sub layout and right layout
	lSubHBox := widgets.NewQHBoxLayout()
	rVBox := widgets.NewQVBoxLayout2(nil)

	var runButton *widgets.QPushButton
	var liveToggle *widgets.QPushButton

	//reset button
	resetButton := widgets.NewQPushButton2("reset", nil)
	resetButton.SetIcon(resetButton.Style().StandardIcon(widgets.QStyle__SP_DialogResetButton, nil, nil))
	lSubHBox.AddWidget(resetButton, 1, 0)
	resetButton.ConnectClicked(func(bool) {
		file := core.NewQFile2(":/qml/" + fileForIndex(tabWidget.CurrentIndex()))
		if file.Open(core.QIODevice__ReadOnly) {
			defer file.Close()
			textEdits[tabWidget.CurrentIndex()].Document().SetPlainText(file.ReadAll().ConstData())
		}
		if tabWidget.CurrentIndex() == 0 {
			textEditQSS.SetPlainText(defQSS)
		}
		if !liveToggle.IsChecked() {
			runButton.Click()
		}
	})

	var engine *qml.QJSEngine

	//run button
	runButton = widgets.NewQPushButton2("run", nil)
	runButton.SetIcon(runButton.Style().StandardIcon(widgets.QStyle__SP_MediaPlay, nil, nil))
	lSubHBox.AddWidget(runButton, 2, 0)
	runButton.ConnectClicked(func(bool) {

		//delete old widget
		for rVBox.Count() > 0 {
			ci := rVBox.TakeAt(0)
			ci.Widget().Hide()
			ci.Widget().DeleteLater()
			ci.DestroyQLayoutItem()
		}

		//add new widget
		var cWidget *widgets.QWidget
		switch i := tabWidget.CurrentIndex(); i {
		case 0, 1, 2, 3:
			switch runtime.GOARCH {
			case "js": //TODO: browser js api support for wasm
				cWidget = widgets.NewQWidgetFromPointer(unsafe.Pointer(js.Global.Call("eval", textEdits[tabWidget.CurrentIndex()].ToPlainText()).Unsafe()))
			default:
				if engine == nil {
					engine = qml.NewQJSEngine()
				}
				if engine.GlobalObject().Property("setInterval").IsUndefined() {
					engine.NewGoType("setInterval", func(f func(), msec int) {
						t := core.NewQTimer(nil)
						t.ConnectTimeout(f)
						t.Start(msec)
					})
				}
				cWidget = widgets.NewQWidgetFromPointer(unsafe.Pointer(uintptr(engine.Evaluate(textEdits[tabWidget.CurrentIndex()].ToPlainText(), "", 0).ToVariant().ToULongLong(nil))))
			}
			if i == 0 {
				cWidget.SetStyleSheet(textEditQSS.ToPlainText())
			}
		case 4:
			cWidget = newListView()
		case 5:
			cWidget = newTableView()
		case 6:
			cWidget = newTreeView()
		case 7:
			cWidget = newTextEditExample(qApp)
		case 8:
			cWidget = new_xkcdWidget()
		case 9, 10:
			cWidget = widgets.NewQWidget(nil, 0)
		}
		rVBox.AddWidget(cWidget, 0, 0)
	})

	//live toggle
	liveToggle = widgets.NewQPushButton2("live", nil)
	liveToggle.SetIcon(liveToggle.Style().StandardIcon(widgets.QStyle__SP_BrowserReload, nil, nil))
	liveToggle.ConnectClicked(func(c bool) {
		if c {
			liveToggle.SetIcon(liveToggle.Style().StandardIcon(widgets.QStyle__SP_BrowserStop, nil, nil))
		} else {
			liveToggle.SetIcon(liveToggle.Style().StandardIcon(widgets.QStyle__SP_BrowserReload, nil, nil))
		}
	})
	liveToggle.SetCheckable(true)
	lSubHBox.AddWidget(liveToggle, 1, 0)

	var timer *time.Timer
	triggerUpdate := func() {
		if liveToggle.IsChecked() {
			if timer == nil {
				timer = time.AfterFunc(250*time.Millisecond, func() {
					runButton.Click()
					timer = nil
				})
			} else {
				timer.Reset(250 * time.Millisecond)
			}
		}
	}

	//connect text changed signals
	for _, textEdit := range textEdits {
		textEdit.ConnectTextChanged(triggerUpdate)
	}
	textEditQSS.ConnectTextChanged(triggerUpdate)

	//connect tab changed signal
	tabWidget.ConnectCurrentChanged(func(i int) {
		if i <= 9 && textEdits[i].Document().IsEmpty() {
			file := core.NewQFile2(":/qml/" + fileForIndex(i))
			if file.Open(core.QIODevice__ReadOnly) {
				defer file.Close()
				liveToogleWasChecked := liveToggle.IsChecked()
				liveToggle.SetChecked(false)
				textEdits[i].Document().SetPlainText(file.ReadAll().ConstData())
				liveToggle.SetChecked(liveToogleWasChecked)
			}
			textEdits[i].MoveCursor(gui.QTextCursor__Start, gui.QTextCursor__MoveAnchor)
			textEdits[i].EnsureCursorVisible()
		}

		liveToggle.SetEnabled(i <= 3)
		resetButton.SetEnabled(i <= 3)

		switch i {
		case 9, 10:
			resetButton.SetVisible(false)
			runButton.SetVisible(false)
			liveToggle.SetVisible(false)

			rVBox.ItemAt(0).Widget().SetVisible(false)
			windowLayout.SetStretchFactor2(rVBox, 0)

		default:
			resetButton.SetVisible(true)
			runButton.SetVisible(true)
			liveToggle.SetVisible(true)

			runButton.Clicked(false)

			rVBox.ItemAt(0).Widget().SetVisible(true)
			windowLayout.SetStretchFactor2(rVBox, 2)
		}
	})
	tabWidget.CurrentChanged(0)

	//add sub layout to lVBox
	lVBox.AddLayout(lSubHBox, 0)

	//add main layouts to window layout
	windowLayout.AddLayout(lVBox, 3)
	windowLayout.AddLayout(rVBox, 2)

	//show window and exec
	switch runtime.GOARCH {
	case "js", "wasm":
		window.ShowFullScreen()
	default:
		window.Show()
	}
	widgets.QApplication_Exec()
}

func newDefaultFont() *gui.QFont {
	font := gui.NewQFont()
	font.SetPointSize(10)
	return font
}

func newTextEdit(n string) *widgets.QTextEdit {
	//code syntax highlighter
	doc := gui.NewQTextDocument(nil)
	doc.SetDefaultFont(newDefaultFont())

	//code text area
	textEdit := widgets.NewQTextEdit(nil)
	textEdit.SetDocument(doc)
	textEdit.SetTabStopDistance(textEdit.TabStopDistance() / 3)
	NewGolangHighlighter(textEdit.Document())
	return textEdit
}

func newTextEditQSS() *widgets.QPlainTextEdit {
	textEditQSS := widgets.NewQPlainTextEdit(nil)
	textEditQSS.SetFont(newDefaultFont())
	textEditQSS.SetPlainText(defQSS)
	return textEditQSS
}

func fileForIndex(i int) string {
	switch i {
	case 0:
		return "line_edits.js"
	case 1:
		return "graphics_scene.js"
	case 2:
		return "pixel_editor.js"
	case 3:
		return "basic.js"
	case 4:
		return "list_view.go"
	case 5:
		return "table_view.go"
	case 6:
		return "tree_view.go"
	case 7:
		return "text_edit.go"
	case 8:
		return "xkcd.go"
	case 9:
		return "source.go"
	}
	return ""
}

func newInfoWidget() *widgets.QWidget {

	infoWidget := widgets.NewQWidget(nil, 0)
	repoLayout := widgets.NewQVBoxLayout2(infoWidget)

	//label
	label := widgets.NewQLabel2(`This playground is a Proof of Concept showcase to illustrate that it's possible to create SPAs entirely in Go and/or JavaScript by utilizing <a href="https://www.qt.io">Qt</a> and <a href="https://github.com/therecipe/qt">therecipe/qt</a><br>
<br>
This showcase isn't meant to show an aesthetically pleasing application but rather an functional application using the Qt Widgets module.<br>
There will be another playground to showcase more modern looking QML applications in the future.<br>
It will show amongst other things examples with Material and Metro design as well as examples with heavy 2 and 3D animations.<br>
<br>
However, this showcase is still experimental and might bug out from time to time, there are also at least two known issues making the WebAssembly binary (or generated JavaScript code) unnecessary large.<br>
It should be possible to reduce the size at least about 20% or more in the future and also improve the general performance.<br>
<br>
This playground can also be compiled for all major operating systems (desktop and mobile) and also for the full WebAssembly target without any changes.<br>
The JavaScript examples won't work though, but if there is enough interest then there is probably a way to make the JavaScript api work on the desktop and mobile targets as well.<br>
<br>
The current version of this showcase did use gopherjs to transpile the Go code and expose the Qt api to JavaScript, later versions will use the Go WebAssembly target for both tasks.<br>
The Qt code is always compiled to WebAssembly using emscripten for both the "js" and the "wasm" target.<br>
<br>
Special thanks to <b>@neelance</b> for creating <a href="https://github.com/gopherjs/gopherjs">gopherjs</a> and also for the work on the WebAssembly target for Go, <b>@lpotter</b> for working on <a href="https://wiki.qt.io/Qt_for_WebAssembly">Qt for WebAssembly</a>, <b>@kripken</b> for <a href="https://github.com/kripken/emscripten">emscripten</a>, <b>@5k3105</b> for the original code of some of the examples and porting over the syntax highlighter to Go, and also <b>@egonelbre</b> for the <a href="https://github.com/egonelbre/gophers">gopher image</a>.<br>`, nil, 0)
	label.SetTextInteractionFlags(core.Qt__LinksAccessibleByMouse)
	label.ConnectLinkActivated(func(link string) {
		switch runtime.GOARCH {
		case "js", "wasm":
			qt.Global.Call("eval", "window.open('"+link+"', '_blank')")
		default:
			gui.QDesktopServices_OpenUrl(core.NewQUrl3(link, 0))
		}
	})
	label.SetWordWrap(true)
	repoLayout.AddWidget(label, 0, 0)

	//button
	button := widgets.NewQPushButton2("open repo in new tab", nil)
	button.ConnectClicked(func(bool) {
		link := "https://github.com/therecipe/widgets_playground"
		switch runtime.GOARCH {
		case "js", "wasm":
			qt.Global.Call("eval", "window.open('"+link+"', '_blank')")
		default:
			gui.QDesktopServices_OpenUrl(core.NewQUrl3(link, 0))
		}
	})
	repoLayout.AddWidget(button, 0, 0)

	return infoWidget
}
