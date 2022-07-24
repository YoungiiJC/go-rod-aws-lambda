package main

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
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
	const timeout = 10

	// create a Rod browser instance
	var browser *rod.Browser
	u := launchInLambda().MustLaunch()
	browser = rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	// visit the page
	var page *rod.Page
	err := rod.Try(func() {
		page = browser.Timeout(timeout * time.Second).MustPage("http://www." + url)
	})

	if err != nil {
		return "", err
	}

	// return the entire html
	return page.MustElement("html").MustHTML(), nil
}
