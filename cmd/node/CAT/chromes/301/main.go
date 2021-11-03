package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"io/ioutil"
	"log"
)

func Pdrintf(format string, v ...interface{}) {
	fmt.Println("start......")
	fmt.Printf(format, v...)
}

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	chromedp.WithBrowserOption(
		chromedp.WithConsolef(Pdrintf),
		chromedp.WithBrowserLogf(Pdrintf),
		chromedp.WithBrowserDebugf(Pdrintf),
	)
	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel1()
	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(Pdrintf))
	defer cancel()

	defer cancel()
	// run
	var b1 []byte

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventExceptionThrown:
			fmt.Printf("Event Exception, console time > %s \n", ev.Timestamp.Time())
			fmt.Printf("\tException Type > %s \n", ev.ExceptionDetails.Exception.Type.String())
			fmt.Printf("\tException Description > %s \n", ev.ExceptionDetails.Exception.Description)
			fmt.Printf("\tException Stacktrace Text > %s \n", ev.ExceptionDetails.Exception.ClassName)
		case *runtime.StackTrace:
			fmt.Printf("Stack Trace, console type > %s \n", ev.Description)
			for _, frames := range ev.CallFrames {
				fmt.Printf("Frame line # %s\n", frames.LineNumber)
			}
		case *runtime.EventConsoleAPICalled:
			fmt.Printf("Event Console API Called, console type > %s call:\n", ev.Type)
			for _, arg := range ev.Args {
				fmt.Printf("%s - %s\n", arg.Type, arg.Value)
			}
		case *cdproto.Message:
		case *network.EventLoadingFailed:
			fmt.Println("浏览器加载失败的日志...........")
			fmt.Println(ev.ErrorText)
			json, _ := ev.MarshalJSON()
			fmt.Println(string(json))
			fmt.Println("浏览器加载失败的日志...........")

		case *network.EventLoadingFinished:
			fmt.Println("浏览器加载完成...........")
			json, _ := ev.MarshalJSON()
			fmt.Println(string(json))
			fmt.Println("浏览器加载完成...........")

		default:
			//log.Println("type: ", reflect.TypeOf(ev), "ev data: ", ev)
		}
	})
	if err := chromedp.Run(ctx,
		// emulate iPhone 7 landscape
		chromedp.Emulate(device.IPhone8Plus),
		chromedp.Navigate(`https://baidu.com`),
		chromedp.CaptureScreenshot(&b1),
	); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("baidu_IPhone8Plus.png", b1, 0777); err != nil {
		log.Fatal(err)
	}

}
