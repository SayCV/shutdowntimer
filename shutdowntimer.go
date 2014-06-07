package main

/*
	Shutdown Timer for Windows
	date: 2014-06-06
	Auther: Donald Tu
	Contact: support@software-download.name
*/

import (
	"flag"
	"github.com/rakyll/globalconf"
	jww "github.com/spf13/jwalterweatherman"
	"runtime"
	"strconv"
	"time"
)

var name = "Shutdown Timer"
var version = "201406"
var home = "http://software-download.name/2012/windows-command-line-shutdown/"

var timerDelayed = false

type myflag struct {
	Int int
}

func (m myflag) String() string {
	return strconv.Itoa(m.Int)
}

func (m myflag) Set(str string) (err error) {
	m.Int, err = strconv.Atoi(str)
	return err
}

type Opts struct {
	Show bool

	BeforeHour int
	BeforeMin  int
	AfterHour  int
	AfterMin   int
}

var opt = &Opts{}

func init() {
	//jww.UseTempLogFile("YourAppName")
	//jww.SetLogThreshold(jww.LevelDebug)
	//    LevelDebug LevelInfo LevelWarn
	//jww.SetStdoutThreshold(jww.LevelWarn)

	jww.INFO.Println("Init in")

	flag.BoolVar(&opt.Show, "show", false, "true to disable GUI window")
	flag.IntVar(&opt.BeforeHour, "beforehour", 5, "Shutdown if current time is before specified oclock")
	flag.IntVar(&opt.BeforeMin, "beforemin", 0, "Shutdown if current time is before specified minute time")
	flag.IntVar(&opt.AfterHour, "afterhour", 22, "Shutdown if current time is after specified oclock")
	flag.IntVar(&opt.AfterMin, "aftermin", 0, "Shutdown if current time is after specified minute time")
	jww.INFO.Println("Init out")

}

var conf *globalconf.GlobalConf

func testShutdown() {
	for {
		jww.INFO.Println("for in testShutdown() ")
		if isShutdown() {
			jww.INFO.Println("isShutdown() true")
			if !timerDelayed && opt.Show {
				timerDelayed = true
				showGUI()
				enablePost()
				setTip("Warning: Windows will shutdow now")
				time.Sleep(time.Minute)
				continue
			}

			if opt.Show {
				closeGUI()
				opt.Show = false
			}

			shutdown()
			break
		}
		//runtime.Gosched()
		time.Sleep(time.Minute)
	}
}

func main() {
	conf, _ = globalconf.New(name)
	flag.Parse()
	if flag.NFlag() == 0 {
		opt.Show = true
	}
	conf.ParseAll() //parseAll will set flags from ini, so we need test NFlag before this

	runtime.GOMAXPROCS(runtime.NumCPU())

	go testShutdown()

	if opt.Show == true {
		jww.INFO.Println("Start initGUI")

		initGUI()
		jww.INFO.Println("GUI started")
	} else {
	}
}

func minuteOfNow() int {
	return time.Now().Hour()*60 + time.Now().Minute()
}

func minuteOfBefore() int {
	return opt.BeforeHour*60 + opt.BeforeMin
}

func minuteOfAfter() int {
	return opt.AfterHour*60 + opt.AfterMin
}

func isShutdown() bool {
	if (minuteOfBefore() != 0 && minuteOfNow() < minuteOfBefore()) ||
		(minuteOfAfter() != 0 && minuteOfNow() > minuteOfAfter()) {
		return true
	}
	return false
}
