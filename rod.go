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

		// no need to use leakless on aws-lambda, lambda will ensure no process leak
		Leakless(false).

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
	const timeout = 10

	// create a Rod browser instance
	var browser *rod.Browser
	u := launchInLambda().MustLaunch()
	browser = rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage()

	// Block loading any resources we don't need in headless
	// https://go-rod.github.io/#/network?id=blocking-certain-resources-from-loading
	router := page.HijackRequests()

	resources := []proto.NetworkResourceType{
		proto.NetworkResourceTypeFont,
		proto.NetworkResourceTypeImage,
		proto.NetworkResourceTypeMedia,
		proto.NetworkResourceTypeStylesheet,
	}

	router.MustAdd("*", func(ctx *rod.Hijack) {
		if slice.Contains(resources, ctx.Request.Type()) {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}

		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()

	err := rod.Try(func() {
		page.Timeout(timeout * time.Second).MustNavigate("http://www." + url).MustWaitLoad()
	})

	if err != nil {
		return "", err
	}

	// wait until the body loads
	page.MustElement("body")

	// then return the entire html
	return page.MustElement("html").MustHTML(), nil
}
