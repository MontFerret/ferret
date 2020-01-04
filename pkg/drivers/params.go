package drivers

type (
	Viewport struct {
		Height      int
		Width       int
		ScaleFactor float64
		Mobile      bool
		Landscape   bool
	}

	Params struct {
		URL         string
		UserAgent   string
		KeepCookies bool
		Cookies     HTTPCookies
		Headers     HTTPHeaders
		Viewport    *Viewport
	}

	ParseParams struct {
		Content     []byte
		KeepCookies bool
		Cookies     HTTPCookies
		Headers     HTTPHeaders
		Viewport    *Viewport
	}
)
