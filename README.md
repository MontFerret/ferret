# Ferret
<p align="center">
	<a href="https://goreportcard.com/report/github.com/MontFerret/ferret">
		<img alt="Go Report Status" src="https://goreportcard.com/badge/github.com/MontFerret/ferret">
	</a>
	<a href="https://github.com/MontFerret/ferret/actions">
		<img alt="Build Status" src="https://github.com/MontFerret/ferret/workflows/build/badge.svg">
	</a>
	<a href="https://codecov.io/gh/MontFerret/ferret">
		<img src="https://codecov.io/gh/MontFerret/ferret/branch/master/graph/badge.svg" />
	</a>
	<a href="https://discord.gg/kzet32U">
		<img alt="Discord Chat" src="https://img.shields.io/discord/501533080880676864.svg">
	</a>
	<a href="https://github.com/MontFerret/ferret/releases">
		<img alt="Ferret release" src="https://img.shields.io/github/release/MontFerret/ferret.svg">
	</a>
	<a href="http://opensource.org/licenses/MIT">
		<img alt="MIT License" src="http://img.shields.io/badge/license-MIT-brightgreen.svg">
	</a>
</p>

![ferret](https://raw.githubusercontent.com/MontFerret/ferret/master/assets/intro.jpg)

## What is it?
```ferret``` is a web scraping system. It aims to simplify data extraction from the web for UI testing, machine learning, analytics and more.    
```ferret``` allows users to focus on the data. It abstracts away the technical details and complexity of underlying technologies using its own declarative language. 
It is extremely portable, extensible and fast.

[Read the introductory blog post about Ferret here!](https://medium.com/@ziflex/say-hello-to-ferret-a-modern-web-scraping-tool-5c9cc85ba183)

### Features

* Declarative language
* Support of both static and dynamic web pages
* Embeddable
* Extensible

### Show me some code
The following example demonstrates the use of dynamic pages.    
We load the main Google Search page, type search criteria into an input box and then click a search button.   
The click action triggers a redirect, so we wait until its end.   
Once the page gets loaded, we iterate over all elements in search results and assign the output to a variable.   
The final for loop filters out empty elements that might be because of inaccurate use of selectors.      

```aql
LET google = DOCUMENT("https://www.google.com/", {
    driver: "cdp",
    userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.87 Safari/537.36"
})

INPUT(google, 'input[name="q"]', "ferret")
CLICK(google, 'input[name="btnK"]')

WAIT_NAVIGATION(google)

FOR result IN ELEMENTS(google, '.g')
    // filter out extra elements like videos and 'People also ask'
    FILTER TRIM(result.attributes.class) == 'g'
    RETURN {
        title: INNER_TEXT(result, 'h3'),
        description: INNER_TEXT(result, '.st'),
        url: INNER_TEXT(result, 'cite')
    }
```

More examples you can find [here](./examples)


### Motivation
Nowadays data is everything and who owns data - owns the world.    
I have worked on multiple data-driven projects where data was an essential part of a system and I realized how cumbersome writing tons of scrapers is.    
After some time looking for a tool that would let me to not write a code, but just express what data I need, decided to come up with my own solution.    
```ferret``` project is an ambitious initiative trying to bring the universal platform for writing scrapers without any hassle.    

### Inspiration
FQL (Ferret Query Language) is heavily inspired by [AQL](https://www.arangodb.com/) (ArangoDB Query Language).    
But due to the domain specifics, there are some differences in syntax and how things work.     


## Installation

### Binary
You can download latest binaries from [here](https://github.com/MontFerret/ferret/releases).

### Source code
#### Production
* Go >=1.11
* Chrome or Docker

#### Development
* GNU Make
* ANTLR4 >=4.8


```sh
go get github.com/MontFerret/ferret
```

### Environment

In order to use all Ferret features, you will need to have Chrome either installed locally or running in Docker.
For ease of use we recommend to run [Chromium inside a Docker container](https://github.com/MontFerret/chromium):

```sh
docker pull montferret/chromium
docker run -d -p 9222:9222 montferret/chromium
```

But if you want to see what's happening during query execution, just start your Chrome with remote debugging port:

```sh
chrome.exe --remote-debugging-port=9222
```

## Quick start

### Browserless mode

If you want to play with ```fql``` and check its syntax, you can run CLI with the following commands:
```
ferret
```

```ferret``` will run in REPL mode.

```shell
Welcome to Ferret REPL
Please use `exit` or `Ctrl-D` to exit this program.
>%
>LET doc = DOCUMENT('https://news.ycombinator.com/')
>FOR post IN ELEMENTS(doc, '.storylink')
>RETURN post.attributes.href
>%

```

**Note:** symbol ```%``` is used to start and end multi-line queries. You also can use the heredoc format.

If you want to execute a query stored in a file, just pass a file name:

```
ferret ./docs/examples/static-page.fql
```

```
cat ./docs/examples/static-page.fql | ferret
```

```
ferret < ./docs/examples/static-page.fql
```


### Browser mode

By default, ``ferret`` loads HTML pages via HTTP protocol, because it's faster.    
But nowadays, there are more and more websites rendered with JavaScript, and therefore, this 'old school' approach does not really work.    
For such cases, you may fetch documents using Chrome or Chromium via Chrome DevTools protocol (aka CDP).    
First, you need to make sure that you launched Chrome with ```remote-debugging-port=9222``` flag.    
Second, you need to pass the address to ```ferret``` CLI.    

```
ferret --cdp http://127.0.0.1:9222
```

**NOTE:** By default, ```ferret``` will try to use this local address as a default one, so it makes sense to explicitly pass the parameter only in case of either different port number or remote address.    

Alternatively, you can tell CLI to launch Chrome for you.

```shell
ferret --cdp-launch
```

**NOTE:** Launch command is currently broken on MacOS.

Once ```ferret``` knows how to communicate with Chrome, you can use a function ```DOCUMENT(url, isDynamic)``` with ```true``` boolean value for dynamic pages:

```shell
Welcome to Ferret REPL
Please use `exit` or `Ctrl-D` to exit this program.
>%
>LET doc = DOCUMENT('https://soundcloud.com/charts/top', { driver: "cdp" })
>WAIT_ELEMENT(doc, '.chartTrack__details', 5000)
>LET tracks = ELEMENTS(doc, '.chartTrack__details')
>FOR track IN tracks
>    LET username = ELEMENT(track, '.chartTrack__username')
>    LET title = ELEMENT(track, '.chartTrack__title')
>    RETURN {
>       artist: username.innerText,
>        track: title.innerText
>    }
>%
```

```shell
Welcome to Ferret REPL
Please use `exit` or `Ctrl-D` to exit this program.
>%
>LET doc = DOCUMENT("https://github.com/", { driver: "cdp" })
>LET btn = ELEMENT(doc, ".HeaderMenu a")

>CLICK(btn)
>WAIT_NAVIGATION(doc)
>WAIT_ELEMENT(doc, '.IconNav')

>FOR el IN ELEMENTS(doc, '.IconNav a')
>    RETURN TRIM(el.innerText)
>%
```

### Embedded mode

```ferret``` is a very modular system and therefore, can be easily be embedded into your Go application.

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
)

type Topic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func main() {
	topics, err := getTopTenTrendingTopics()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, topic := range topics {
		fmt.Println(fmt.Sprintf("%s: %s %s", topic.Name, topic.Description, topic.URL))
	}
}

func getTopTenTrendingTopics() ([]*Topic, error) {
	query := `
		LET doc = DOCUMENT("https://github.com/topics")

		FOR el IN ELEMENTS(doc, ".py-4.border-bottom")
			LIMIT 10
			LET url = ELEMENT(el, "a")
			LET name = ELEMENT(el, ".f3")
			LET description = ELEMENT(el, ".f5")

			RETURN {
				name: TRIM(name.innerText),
				description: TRIM(description.innerText),
				url: "https://github.com" + url.attributes.href
			}
	`

	comp := compiler.New()

	program, err := comp.Compile(query)

	if err != nil {
		return nil, err
	}

	// create a root context
	ctx := context.Background()

	// enable HTML drivers
	// by default, Ferret Runtime does not know about any HTML drivers
	// all HTML manipulations are done via functions from standard library
	// that assume that at least one driver is available
	ctx = drivers.WithContext(ctx, cdp.NewDriver())
	ctx = drivers.WithContext(ctx, http.NewDriver(), drivers.AsDefault())

	out, err := program.Run(ctx)

	if err != nil {
		return nil, err
	}

	res := make([]*Topic, 0, 10)

	err = json.Unmarshal(out, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

```

## Extras

### Extensibility

That said, ```ferret``` is a very modular system which also allows not only embed it, but extend its standard library.

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func main() {
	strs, err := getStrings()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, str := range strs {
		fmt.Println(str)
	}
}

func getStrings() ([]string, error) {
	// function implements is a type of a function that ferret supports as a runtime function
	transform := func(ctx context.Context, args ...core.Value) (core.Value, error) {
		// it's just a helper function which helps to validate a number of passed args
		err := core.ValidateArgs(args, 1, 1)

		if err != nil {
			// it's recommended to return built-in None type, instead of nil
			return values.None, err
		}

		// this is another helper functions allowing to do type validation
		err = core.ValidateType(args[0], core.StringType)

		if err != nil {
			return values.None, err
		}

		// cast to built-in string type
		str := args[0].(values.String)

		return values.NewString(strings.ToUpper(str.String() + "_ferret")), nil
	}

	query := `
		FOR el IN ["foo", "bar", "qaz"]
			// conventionally all functions are registered in upper case
			RETURN TRANSFORM(el)
	`

	comp := compiler.New()

	if err := comp.RegisterFunction("transform", transform); err != nil {
		return nil, err
	}

	program, err := comp.Compile(query)

	if err != nil {
		return nil, err
	}

	out, err := program.Run(context.Background())

	if err != nil {
		return nil, err
	}

	res := make([]string, 0, 3)

	err = json.Unmarshal(out, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
```

On top of that, you can completely turn off the standard library, bypassing the following option:

```go
comp := compiler.New(compiler.WithoutStdlib())
```

And after that, you can easily provide your own implementation of functions from standard library.    

If you don't need a particular set of functions from standard library, you can turn off the entire ```stdlib``` and register separate packages from that:    

```go
package main

import (
    "github.com/MontFerret/ferret/pkg/compiler"
    "github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func main() {
    comp := compiler.New(compiler.WithoutStdlib())

    comp.RegisterFunctions(strings.NewLib())
}
```

### Proxy

By default, Ferret does not use any proxies. Partially, due to inability to force Chrome/Chromium (or any other Chrome Devtools Protocol compatible browser) to use a particular proxy. It should be done during a browser launch.

But you can pass an address of a proxy server you want to use for static pages.

#### CLI

```sh
ferret --proxy=http://localhost:8888 my-query.fql
```

#### Code

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
	
    "github.com/MontFerret/ferret/pkg/compiler"
    "github.com/MontFerret/ferret/pkg/drivers"
    "github.com/MontFerret/ferret/pkg/drivers/http"
)

func run(q string) ([]byte, error) {
    proxy := "http://localhost:8888"
    comp := compiler.New()
    program := comp.MustCompile(q)

    // create a root context
    ctx := context.Background()

    // we inform the driver what proxy to use
    ctx = drivers.WithContext(ctx, http.NewDriver(http.WithProxy(proxy)), drivers.AsDefault())

    return program.Run(ctx)
}

```

### Cookies

#### Get, Set, Delete
For more precise work, you can set/get/delete cookies manually before and after loading the page:

```
LET doc = DOCUMENT("https://www.google.com", {
    driver: "cdp",
    cookies: [
         {
             name: "foo",
             value: "bar"
         }
    ]
})

COOKIE_SET(doc, { name: "baz", value: "qaz"}, { name: "daz", value: "gag" })
COOKIE_DEL(doc, "foo")

LET c = COOKIE_GET(doc, "baz")

FOR cookie IN doc.cookies
    RETURN cookie.name
	
```

#### Access previously-set cookies (non-incognito mode)

By default, ``CDP`` driver execute each query in an incognito mode in order to avoid any collisions related to some persisted cookies from previous queries.   
However, sometimes it might not be a desirable behavior and a query needs to be executed within a Chrome tab with earlier persisted cookies.   
In order to do that, we need to inform the driver to execute all queries in regular tabs. Here is how to do that:

##### CLI

```sh
ferret --cdp-keep-cookies my-query.fql
```

##### Code

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
)

func run(q string) ([]byte, error) {
	comp := compiler.New()
	program := comp.MustCompile(q)

	// create a root context
	ctx := context.Background()

	// we inform the driver to keep cookies between queries
	ctx = drivers.WithContext(
		ctx,
		cdp.NewDriver(cdp.WithKeepCookies()),
		drivers.AsDefault(),
	)

	return program.Run(ctx)
}
```

##### Query
```
LET doc = DOCUMENT("https://www.google.com", {
    driver: "cdp",
    keepCookies: true
})
```

### File System

#### Write
```
USE IO::FS

LET favicon = DOWNLOAD("https://www.google.com/favicon.ico")

RETURN WRITE("google.favicon.ico", favicon)
```

#### Read
```
USE IO::FS

LET urls_data = READ("urls.json")
LET urls = JSON_PARSE(urls_data)

FOR url IN urls
    RETURN DOCUMENT(url)
```

## Contributors

Thanks to everyone who contributed.
<a href="https://github.com/MontFerret/ferret/graphs/contributors"><img src="https://opencollective.com/ferret/contributors.svg?width=890&button=false" /></a>

## Financial support

Support this project by becoming a sponsor. Your logo will show up here with a link to your website.    

<p>
		    <a href="https://opencollective.com/ferret#section-contributors" target="_blank">
	<img src="https://opencollective.com/ferret/sponsors.svg?width=890&button=false" />
	</a>
</p>
<p>
	    <a href="https://opencollective.com/ferret#section-contributors" target="_blank">
	<img src="https://opencollective.com/ferret/backers.svg?width=890&button=false" />
	</a>
</p>

<p align="center">
    <a href="https://opencollective.com/ferret/donate" target="_blank">
        <img src="https://opencollective.com/ferret/donate/button@2x.png?color=blue" width="300" />
    </a>
</p>
