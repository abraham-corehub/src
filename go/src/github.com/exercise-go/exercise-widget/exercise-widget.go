package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/admin"
	"github.com/qor/assetfs"
	"github.com/qor/qor"
	"github.com/qor/widget"
)

var tWidgets *widget.Widgets
var tAdmin *admin.Admin

type bannerArgument struct {
	Title    string
	SubTitle string
}

func main() {
	//db := utils.TestDB()
	db, err := gorm.Open(`sqlite3`, `widget.db`)
	if err != nil {
		fmt.Println(err)
	}

	if err := db.DropTableIfExists(&widget.QorWidgetSetting{}).Error; err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&widget.QorWidgetSetting{})
	mux := http.NewServeMux()
	//Server = httptest.NewServer(mux)

	// Default implemention based on filesystem, you could overwrite with other implemention, for example bindatafs will do this.
	aFS := assetfs.AssetFS()

	// Register path to AssetFS
	aFS.RegisterPath("views")

	// Get file's content with name from path `/web/app/views`
	b, err := aFS.Asset("slider.tmpl")

	tWidgets = widget.New(&widget.Config{
		DB: db,
	})
	tWidgets.RegisterViewPath("/views")

	tAdmin = admin.New(&qor.Config{DB: db})
	tAdmin.AddResource(tWidgets)
	tAdmin.MountTo("/admin", mux)

	tWidgets.RegisterWidget(&widget.Widget{
		Name:      "Banner",
		Templates: []string{"slider"},
		Setting:   tAdmin.NewResource(&bannerArgument{}),
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			if setting != nil {
				argument := setting.(*bannerArgument)
				context.Options["Title"] = argument.Title
				context.Options["SubTitle"] = argument.SubTitle
			}
			return context
		},
	})

	tWidgets.RegisterScope(&widget.Scope{
		Name: "From Google",
		Visible: func(context *widget.Context) bool {
			if request, ok := context.Get("Request"); ok {
				_, ok := request.(*http.Request).URL.Query()["from_google"]
				return ok
			}
			return false
		},
	})

	tWidgets.RegisterWidget(&widget.Widget{
		Name:    "Slider",
		Setting: tAdmin.NewResource(&bannerArgument{}),
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			context.Body = string(b)
			return context
		},
	})

	http.ListenAndServe(":8080", mux)
}
