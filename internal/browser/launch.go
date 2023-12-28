package browser

import (
	"bahamut/internal/config"
	"strings"

	"github.com/playwright-community/playwright-go"
)

func Launch(c *config.ConfigBrowser) (playwright.Browser, playwright.BrowserContext, error) {
	runOpts := playwright.RunOptions{
		SkipInstallBrowsers: c.SkipInstallation,
		Browsers:            []string{c.Type},
	}
	err := playwright.Install(&runOpts)
	if err != nil {
		return nil, nil, err
	}
	pw, err := playwright.Run(&runOpts)
	if err != nil {
		return nil, nil, err
	}

	launchOpts := playwright.BrowserTypeLaunchOptions{
		Args:     c.Args,
		Headless: &c.Headless,
	}
	if c.Driver == "chromium" {
		launchOpts.Channel = &c.Type
	}
	browser, err := getBrowserDriver(pw, c.Driver).Launch(launchOpts)
	if err != nil {
		return nil, nil, err
	}
	context, err := browser.NewContext()
	if err != nil {
		return nil, nil, err
	}
	return browser, context, nil
}

func getBrowserDriver(pw *playwright.Playwright, driver string) playwright.BrowserType {
	switch strings.ToLower(driver) {
	case "firefox":
		return pw.Firefox
	case "webkit":
		return pw.WebKit
	case "chromium":
		return pw.Chromium
	}
	return pw.Chromium
}
