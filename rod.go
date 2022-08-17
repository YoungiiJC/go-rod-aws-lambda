package main

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/ssoroka/slice"
)

func launchInLambda() *launcher.Launcher {
	return launcher.New().
		// where lambda runtime stores chromium
		Bin("/opt/chromium").

		// recommended flags to run in serverless environments
		// see https://github.com/alixaxel/chrome-aws-lambda/blob/master/source/index.ts
		Set("allow-running-insecure-content").
		Set("autoplay-policy", "user-gesture-required").
		Set("disable-component-update").
		Set("disable-domain-reliability").
		Set("disable-features", "AudioServiceOutOfProcess", "IsolateOrigins", "site-per-process").
		Set("disable-print-preview").
		Set("disable-setuid-sandbox").
		Set("disable-site-isolation-trials").
		Set("disable-speech-api").
		Set("disable-web-security").
		Set("disk-cache-size", "33554432").
		Set("enable-features", "SharedArrayBuffer").
		Set("hide-scrollbars").
		Set("ignore-gpu-blocklist").
		Set("in-process-gpu").
		Set("mute-audio").
		Set("no-default-browser-check").
		Set("no-pings").
		Set("no-sandbox").
		Set("no-zygote").
		Set("single-process").
		Set("use-gl", "swiftshader").
		Set("window-size", "1920", "1080")
}

func getPageHTML(url string) (string, error) {
	// If Rod fails, it needs to correctly timeout before the timeout we set as the lambda fn's timeout
	// this ensures that the browser instance is properly killed and cleaned up
	//
	// these timeouts should collectively be less than the timeout we set for the lambda
	const navigateTimeout = 5 * time.Second
	const navigationTimeout = 5 * time.Second
	const requestIdleTimeout = 10 * time.Second
	const htmlTimeout = 5 * time.Second

	var html string

	err := rod.Try(func() {
		// instantiate the chromium launcher
		launcher := launchInLambda()
		defer launcher.Cleanup()
		defer launcher.Kill()
		u := launcher.MustLaunch()

		// create a browser instance
		browser := rod.New().ControlURL(u).MustConnect()

		// open a page
		page := browser.MustPage()

		// Block loading any resources we don't need in headless
		// https://go-rod.github.io/#/network?id=blocking-certain-resources-from-loading
		router := page.HijackRequests()

		resources := []proto.NetworkResourceType{
			proto.NetworkResourceTypeFont,
			proto.NetworkResourceTypeImage,
			proto.NetworkResourceTypeMedia,
			proto.NetworkResourceTypeStylesheet,
			proto.NetworkResourceTypeWebSocket, // we don't need websockets to fetch html
		}

		router.MustAdd("*", func(ctx *rod.Hijack) {
			if slice.Contains(resources, ctx.Request.Type()) {
				ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
				return
			}

			ctx.ContinueRequest(&proto.FetchContinueRequest{})
		})

		go router.Run()

		// go to the url
		page.Timeout(navigateTimeout).MustNavigate(url)

		// follow any redirects
		// https://github.com/go-rod/rod/issues/640#issuecomment-1171941374
		waitNavigation := page.Timeout(navigationTimeout).MustWaitNavigation()
		waitNavigation()

		// wait until requests stop firing so we can get
		// any html rendered by js scripts or iframes
		waitRequestIdle := page.Timeout(requestIdleTimeout).MustWaitRequestIdle()
		waitRequestIdle()

		// return the html
		html = page.Timeout(htmlTimeout).MustElement("html").MustHTML()
	})

	if err != nil {
		return "", err
	}

	return html, nil
}
