//source: https://github.com/therecipe/qt/tree/master/internal/examples/widgets/xkcd

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/therecipe/qt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var data_struct struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

func new_xkcdWidget() *widgets.QWidget {

	widget := widgets.NewQWidget(nil, 0)

	layout := widgets.NewQFormLayout(widget)
	layout.SetFieldGrowthPolicy(widgets.QFormLayout__AllNonFixedFieldsGrow)

	widgetmap := make(map[string]*widgets.QWidget)
	for i := 0; i < reflect.TypeOf(data_struct).NumField(); i++ {
		name := reflect.TypeOf(data_struct).Field(i).Tag.Get("json")

		if name != "img" {
			widgetmap[name] = widgets.NewQLineEdit(nil).QWidget_PTR()

			layout.AddRow3(name, widgetmap[name])
		} else {
			widgetmap[name] = widgets.NewQLineEdit(nil).QWidget_PTR()

			label := widgets.NewQLabel(nil, 0)
			label.ConnectLinkActivated(func(link string) {
				switch runtime.GOARCH {
				case "js", "wasm":
					qt.Global.Call("eval", "window.open('"+link+"', '_blank')")
				default:
					gui.QDesktopServices_OpenUrl(core.NewQUrl3(link, 0))
				}
			})
			widgetmap[name+"_label"] = label.QWidget_PTR()

			layout.AddRow3(name, widgetmap[name])
			layout.AddRow3(name+"_label", widgetmap[name+"_label"])
		}
	}

	button := widgets.NewQPushButton2("random xkcd", nil)
	layout.AddWidget(button)
	button.ConnectClicked(func(bool) {
		go func() {
			rand.Seed(time.Now().UnixNano())

			url := fmt.Sprintf("https://xkcd.com/%v/info.0.json", rand.Intn(614))
			switch runtime.GOARCH {
			case "js", "wasm":
				url = "https://yacdn.org/proxy/" + url
			}

			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			data, _ := ioutil.ReadAll(resp.Body)

			json.Unmarshal(data, &data_struct)

			for i := 0; i < reflect.TypeOf(data_struct).NumField(); i++ {
				name := reflect.TypeOf(data_struct).Field(i).Tag.Get("json")

				if name != "img" {
					switch reflect.ValueOf(data_struct).Field(i).Kind() {
					case reflect.String:
						widgets.NewQLineEditFromPointer(widgetmap[name].Pointer()).SetText(reflect.ValueOf(data_struct).Field(i).String())
					case reflect.Int:
						widgets.NewQLineEditFromPointer(widgetmap[name].Pointer()).SetText(strconv.Itoa(int(reflect.ValueOf(data_struct).Field(i).Int())))
					}
				} else {
					url := reflect.ValueOf(data_struct).Field(i).String()

					widgets.NewQLineEditFromPointer(widgetmap[name].Pointer()).SetText(url)
					widgets.NewQLabelFromPointer(widgetmap[name+"_label"].Pointer()).SetText(fmt.Sprintf("<a href=\"%[1]v\">%[1]v</a>", url))
				}
			}
		}()
	})

	return widget
}
