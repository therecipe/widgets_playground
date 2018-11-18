// +build !ios,!android

package main

import (
	"fmt"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/printsupport"
	"github.com/therecipe/qt/widgets"
)

func (t *TextEdit) filePrintPreview() {
	if len(os.Getenv("QT_NO_PRINTER")) == 0 && len(os.Getenv("QT_NO_PRINTDIALOG")) == 0 {
		var (
			printer = printsupport.NewQPrinter(printsupport.QPrinter__HighResolution)
			preview = printsupport.NewQPrintPreviewDialog(printer, t, 0)
		)
		preview.ConnectPaintRequested(t.printPreview)
		preview.Exec()
	}
}

func (t *TextEdit) printPreview(printer *printsupport.QPrinter) {
	if len(os.Getenv("QT_NO_PRINTER")) == 0 {
		t.textEdit.Print(printer)
	}
}

func (t *TextEdit) filePrintPdf() {
	if len(os.Getenv("QT_NO_PRINTER")) == 0 {
		var fileDialog = widgets.NewQFileDialog2(t, "Export PDF", "", "")
		fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptSave)
		fileDialog.SetMimeTypeFilters([]string{"application/pdf"})
		fileDialog.SetDefaultSuffix("pdf")
		if fileDialog.Exec() != int(widgets.QDialog__Accepted) {
			return
		}

		var (
			fileName = fileDialog.SelectedFiles()[0]
			printer  = printsupport.NewQPrinter(printsupport.QPrinter__HighResolution)
		)
		printer.SetOutputFormat(printsupport.QPrinter__PdfFormat)
		printer.SetOutputFileName(fileName)
		t.textEdit.Document().Print(printer)
		t.StatusBar().ShowMessage(fmt.Sprintf("Exported %v", core.QDir_ToNativeSeparators(fileName)), 0)
	}
}

func (t *TextEdit) filePrint() {
	if len(os.Getenv("QT_NO_PRINTER")) == 0 && len(os.Getenv("QT_NO_PRINTDIALOG")) == 0 {
		var (
			printer = printsupport.NewQPrinter(printsupport.QPrinter__HighResolution)
			dlg     = printsupport.NewQPrintDialog(printer, t)
		)
		if t.textEdit.TextCursor().HasSelection() {
			dlg.SetOption(printsupport.QAbstractPrintDialog__PrintSelection, true)
		}
		dlg.SetWindowTitle("Print Document")
		if dlg.Exec() == int(widgets.QDialog__Accepted) {
			t.textEdit.Print(printer)
		}
		printer.DestroyQPrinter()
	}
}
