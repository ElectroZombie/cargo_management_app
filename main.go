package main

import("context";"embed";"log";"os";"path/filepath";"cargo_management_app/database";"github.com/wailsapp/wails/v2";"github.com/wailsapp/wails/v2/pkg/options";"github.com/wailsapp/wails/v2/pkg/options/assetserver")
//go:embed all:frontend/dist
var assets embed.FS
type App struct{ctx context.Context;db *database.Database}
func NewApp()*App{return &App{}}
func(a *App)startup(ctx context.Context){a.ctx=ctx;path:=getDBPath();var err error;a.db,err=database.NewDatabase(path);if err!=nil{log.Fatalf("initialize database: %v",err)}}
func(a *App)shutdown(ctx context.Context){if a.db!=nil{_ = a.db.Close()}}
func getDBPath()string{home,err:=os.UserHomeDir();if err!=nil{log.Fatal(err)};dir:=filepath.Join(home,".cargo_management_app");if err:=os.MkdirAll(dir,0755);err!=nil{log.Fatal(err)};return filepath.Join(dir,"cargo_management.db")}
func main(){app:=NewApp();if err:=wails.Run(&options.App{Title:"Cargo Management App",Width:1024,Height:768,AssetServer:&assetserver.Options{Assets:assets},OnStartup:app.startup,OnShutdown:app.shutdown,Bind:[]interface{}{app}});err!=nil{log.Fatal(err)}}
