package drivers

type (
	ScreenSize struct {
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
		Cookies     []HTTPCookie
		Header      HTTPHeader
		Screen      *ScreenSize
	}
)
