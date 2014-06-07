package main

/*
help for Windows: http://www.microsoft.com/resources/documentation/windows/xp/all/proddocs/en-us/shutdown.mspx
build: go build "-ldflags=-H windowsgui"

Create syso using rsrc (https://github.com/akavel/rsrc)
rsrc [-manifest FILE.exe.manifest] [-ico FILE.ico[,FILE2.ico...]] -o FILE.syso

*/

%SystemRoot%\system32;%SystemRoot%;%SystemRoot%\System32\Wbem;%SYSTEMROOT%\System32\WindowsPowerShell\v1.0\;C:\Dropbox\short;C:\Dropbox\d-language\dmd\windows\bin;C:\Dropbox\d-language\dm\bin;D:\Go\bin;D:\golib\bin;D:\MinGW32\bin;C:\Dropbox\software;E:\Applications\TortoiseHg\;C:\Program Files (x86)\QuickTime\QTSystem\;C:\Program Files (x86)\Common Files\Adobe\AGL;D:\msysgit\bin;C:\Program Files (x86)\GNU\Claws Mail\pub

import (
	"flag"
	"fmt"
	"github.com/AllenDang/gform"
	"github.com/AllenDang/w32"
	jww "github.com/spf13/jwalterweatherman"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var space = 10
var height = 20

var mw *gform.Form
var beforeHour, beforeMin, afterHour, afterMin *gform.Edit
var btnSave, btnHide, btnPost, btnHelp, btnClose *gform.PushButton
var lblTip *gform.Label

var closeTimes = 0

func btnHide_onclick(arg *gform.EventArg) {
	mw.Hide()
}

func btnClose_onclick(arg *gform.EventArg) {
	closeTimes++
	if closeTimes > 2 {
		opt.Show = false
		mw.Close()
		os.Exit(0)
	}
}

func btnPost_onclick(arg *gform.EventArg) {
	if minuteOfNow() < minuteOfBefore() {
		opt.BeforeMin = opt.BeforeMin - 5
		if opt.BeforeMin < 0 {
			opt.BeforeHour--
			opt.BeforeMin = opt.BeforeMin + 60
		}
		jww.INFO.Println("Posted 5 minutes")
	} else if minuteOfNow() > minuteOfAfter() {
		opt.AfterMin = opt.AfterMin + 5
		if opt.AfterMin > 60 {
			opt.AfterHour++
			opt.AfterMin = opt.AfterMin - 60
		}
		jww.INFO.Println("Posted 5 minutes")
	}
	btnPost.SetEnabled(false)
	initValue()
}

func btnHelp_onclick(arg *gform.EventArg) {
	w32.ShellExecute(w32.HWND(0), "open", home, "", "", 3)
	setTip("Open: " + home)
}

func btnSave_onclick(arg *gform.EventArg) {
	i, err := strconv.Atoi(beforeHour.Caption())
	if err != nil {
		setTip("Invalid: " + beforeHour.Caption())
		return
	}
	opt.BeforeHour = i

	i, err = strconv.Atoi(beforeMin.Caption())
	if err != nil {
		setTip("Invalid value: " + beforeMin.Caption())
		return
	}
	opt.BeforeMin = i

	i, err = strconv.Atoi(afterHour.Caption())
	if err != nil {
		setTip("Invalid value: " + afterHour.Caption())
		return
	}
	opt.AfterHour = i

	i, err = strconv.Atoi(afterMin.Caption())
	if err != nil {
		setTip("Invalid value: " + afterMin.Caption())
		return
	}
	opt.AfterMin = i

	f := &flag.Flag{Name: "beforehour", Value: myflag{Int: opt.BeforeHour}}
	conf.Set("", f)

	f = &flag.Flag{Name: "beforemin", Value: myflag{Int: opt.BeforeMin}}
	conf.Set("", f)

	f = &flag.Flag{Name: "afterhour", Value: myflag{Int: opt.AfterHour}}
	conf.Set("", f)

	f = &flag.Flag{Name: "aftermin", Value: myflag{Int: opt.AfterMin}}
	conf.Set("", f)

	setTip("Saved")
}

func initValue() {
	jww.INFO.Println("Init Value in")
	text := strconv.Itoa(opt.BeforeHour)
	beforeHour.SetCaption(text)

	text = strconv.Itoa(opt.BeforeMin)
	beforeMin.SetCaption(text)

	text = strconv.Itoa(opt.AfterHour)
	afterHour.SetCaption(text)

	text = strconv.Itoa(opt.AfterMin)
	afterMin.SetCaption(text)
	jww.INFO.Println("Init Value Out")
}

func initGUI() {
	jww.INFO.Println("initGUI in")
	gform.Init()

	mw = gform.NewForm(nil)
	mw.EnableMaxButton(false)
	mw.EnableMinButton(false)

	//mw.SetPos(300, 100)
	mw.Center()
	mw.SetSize(350, 200)
	mw.SetCaption(name + " v" + version)

	lbl := gform.NewLabel(mw)
	lbl.SetCaption("Shutdown before Oclock:")
	lbl.SetPos(space+space/2, space+space/2)
	lbl.SetSize(170, height)

	beforeHour = gform.NewEdit(mw)
	beforeHour.SetCaption("5")
	w, h := lbl.Size()
	x, y := lbl.Pos()
	beforeHour.SetPos(x+w+space, y)
	beforeHour.SetSize(40, height)

	w, h = beforeHour.Size()
	x, y = beforeHour.Pos()
	lbl = gform.NewLabel(mw)
	lbl.SetCaption(":")
	lbl.SetPos(x+w+3, y)
	lbl.SetSize(4, height)

	beforeMin = gform.NewEdit(mw)
	beforeMin.SetCaption("0")
	w, h = lbl.Size()
	x, y = lbl.Pos()
	beforeMin.SetPos(x+w+3, y)
	beforeMin.SetSize(40, height)

	// line 2
	lbl = gform.NewLabel(mw)
	lbl.SetCaption("Shutdown after Oclock:")
	lbl.SetPos(space+space/2, y+h+space)
	lbl.SetSize(170, height)

	afterHour = gform.NewEdit(mw)
	afterHour.SetCaption("23")
	w, h = lbl.Size()
	x, y = lbl.Pos()
	afterHour.SetPos(x+w+space, y)
	afterHour.SetSize(40, height)

	w, h = afterHour.Size()
	x, y = afterHour.Pos()
	lbl = gform.NewLabel(mw)
	lbl.SetCaption(":")
	lbl.SetPos(x+w+3, y)
	lbl.SetSize(4, height)

	afterMin = gform.NewEdit(mw)
	afterMin.SetCaption("0")
	w, h = lbl.Size()
	x, y = lbl.Pos()
	afterMin.SetPos(x+w+3, y)
	afterMin.SetSize(40, height)

	// line 3
	btnSave = gform.NewPushButton(mw)
	btnSave.SetPos(space+space/2, y+h+space+4)
	btnSave.SetCaption("&Save")
	btnSave.SetSize(45, 25)
	btnSave.OnLBUp().Bind(btnSave_onclick)

	w, h = btnSave.Size()
	x, y = btnSave.Pos()
	btnHide = gform.NewPushButton(mw)
	btnHide.SetPos(x+w+space, y)
	btnHide.SetCaption("&Hide")
	btnHide.SetSize(45, 25)
	btnHide.OnLBUp().Bind(btnHide_onclick)

	w, h = btnHide.Size()
	x, y = btnHide.Pos()
	btnPost = gform.NewPushButton(mw)
	btnPost.SetPos(x+w+space, y)
	btnPost.SetCaption("&Post 5 Min")
	btnPost.SetSize(75, 25)
	btnPost.SetEnabled(false)
	btnPost.OnLBUp().Bind(btnPost_onclick)

	w, h = btnPost.Size()
	x, y = btnPost.Pos()
	btnHelp = gform.NewPushButton(mw)
	btnHelp.SetPos(x+w+space, y)
	btnHelp.SetSize(45, 25)
	btnHelp.SetCaption("Hel&p")
	btnHelp.OnLBUp().Bind(btnHelp_onclick)

	w, h = btnHelp.Size()
	x, y = btnHelp.Pos()
	btnClose = gform.NewPushButton(mw)
	btnClose.SetPos(x+w+space, y)
	btnClose.SetSize(45, 25)
	btnClose.SetCaption("&Close")
	btnClose.OnLBUp().Bind(btnClose_onclick)

	// Line 5
	lblTip = gform.NewLabel(mw)
	lblTip.SetPos(space, y+h+space)
	lblTip.SetSize(300, 25)
	lblTip.SetCaption("Tips: 0:0 to disable shutdown")
	mw.Show()

	initValue()

	gform.RunMainLoop()

	jww.INFO.Println("initGUI out")
}

func setTip(text string) {
	var h, m, s = time.Now().Clock()
	str := fmt.Sprintf(" (%v:%v:%v)", h, m, s)
	lblTip.SetCaption(text + str)
}

func showGUI() {
	mw.Show()
}

func closeGUI() {
	mw.Close()
}

func enablePost() {
	btnPost.SetEnabled(true)
}

func shutdown() {
	cmd := exec.Command("shutdown", "-s", "-f", "-c", name, "-t", "60")
	cmd.Run()
	//jww.INFO.Println("Shutdown...")
}
