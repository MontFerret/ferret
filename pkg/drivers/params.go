package drivers

type (
	ResourceFilter struct {
		URL  string
		Type string
	}

	StatusCodeFilter struct {
		URL  string
		Code int
	}

	Ignore struct {
		Resources   []ResourceFilter
		StatusCodes []StatusCodeFilter
	}

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
		Cookies     *HTTPCookies
		Headers     *HTTPHeaders
		Viewport    *Viewport
		Charset     string
		Ignore      *Ignore
	}

	ParseParams struct {
		Content     []byte
		KeepCookies bool
		Cookies     *HTTPCookies
		Headers     *HTTPHeaders
		Viewport    *Viewport
	}
)
