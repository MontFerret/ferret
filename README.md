# Ferret
[![Build Status](https://travis-ci.com/MontFerret/ferret.svg?branch=master)](https://travis-ci.com/MontFerret/ferret)  

![ferret](https://raw.githubusercontent.com/MontFerret/ferret/master/assets/intro.jpg)

## What is it?
```ferret``` is a web scraping system aiming to simplify data extraction from the web for such things like machine learning and analytics.    
Having it's own declarative language, ```ferret``` abstracts away technical details and complexity of the underlying technologies, helping to focus on the data itself.    
It's extremely portable, extensible and fast.

## Show me some code
![code](https://raw.githubusercontent.com/MontFerret/ferret/master/assets/show-me.jpg)
The following example demonstrates use of dynamic pages.    
We are loading main Google Search page, type search criteria into a input box and click search button.   
The click action triggers a redirect, so we wait till it ends.   
Once the page loaded, we iterate over all search elements and assign result to a variable.   
The final for loop filters out empty elements that might be because of inaccurate use of selectors.      

```aql
LET g = DOCUMENT("https://www.google.com/", true)

INPUT(g, 'input[name="q"]', "ferret")
CLICK(g, 'input[name="btnK"]')

WAIT_NAVIGATION(g)

LET result = (
    FOR result IN ELEMENTS(g, '.g')
       RETURN {
           title: ELEMENT(result, 'h3 > a'),
           description: ELEMENT(result, '.st'),
           url: ELEMENT(result, 'cite')
       }
)

RETURN (
    FOR i IN result
    FILTER i.title != NONE
    RETURN i
)
```

## Features

* Declarative language
* Support of both static and dynamic web pages
* Embeddable
* Extensible

## Motivation
Nowadays data is everything and who owns data - owns the world.    
I have worked on multiple data-driven projects where data was an essential part of a system and I realized how cumbersome writing tons of scrapers is.    
After some time looking for a tool that would let me to not write a code, but just express what data I need, decided to come up with my own solution.    
```ferret``` project is an ambitious initiative trying to bring universal platform for writing scrapers without any hassle.    

## Inspiration
FQL (Ferret Query Language) is heavily inspired by [AQL](https://www.arangodb.com/) (ArangoDB Query Language).    
But due to the domain specifics, there are some differences in how things work.     

## WIP
Be aware, the the project is under heavy development. There is no documentation and some things may change in the final release.    
For query syntax, you may go to [ArrangoDB web site](https://docs.arangodb.com/3.3/AQL/index.html) and use AQL docs as docs for FQL - since they are identical.    


## Installation

### Prerequisites
* Go >=1.6
* GoDep
* GNU Make
* Chrome or Docker (optional)

```sh
make build
```

## Quick start

### Browserless mode

If you want to play with ```fql``` and check its syntax, you can run CLI with the following commands:
```
go run ./cmd/cli/main.go
```

```ferret``` will run in REPL mode.

```shell
Welcome to Ferret REPL
Please use `Ctrl-D` to exit this program.
>%
>LET doc = DOCUMENT('https://news.ycombinator.com/')
>FOR post IN ELEMENTS(doc, '.storylink')
>RETURN post.attributes.href
>%

```

**Note:** symbol ```%``` is used to start and end multi line queries. You also can use heredoc format.

If you want to execute a query stored in a file, just pass a file name:

```
go run ./cmd/cli/main.go ./docs/examples/hackernews.fql
```

```
cat ./docs/examples/hackernews.fql | go run ./cmd/cli/main.go 
```

```
go run ./cmd/cli/main.go < ./docs/examples/hackernews.fql
```


### Browser mode

By default, ``ferret`` loads HTML pages via http protocol, because it's faster.    
But nowadays, there are more and more websites rendered with JavaScript, and therefore, this 'old school' approach does not really work.    
For such cases, you may fetch documents using Chrome or Chromium via Chrome DevTools protocol (aka CDP).    
First, you need to make sure that you launched Chrome with ```remote-debugging-port=9222``` flag.    
Second, you need to pass the address to ```ferret``` CLI.    

```
go run ./cmd/cli/main.go --cdp http://127.0.0.1:9222
```

**NOTE:** By default, ```ferret``` will try to use this local address as a default one, so it makes sense to explicitly pass the parameter only in case of either different port number or remote address.    

Alternatively, you can tell CLI to launch Chrome for you.

```shell
go run ./cmd/cli/main.go --cdp-launch
```

**NOTE:** Launch command is currently broken on MacOS.

Once ```ferret``` knows how to communicate with Chrome, you can use a function ```DOCUMENT(url, isDynamic)``` with ```true``` boolean value for dynamic pages:

```shell
Welcome to Ferret REPL
Please use `exit` or `Ctrl-D` to exit this program.
>%
>LET doc = DOCUMENT('https://soundcloud.com/charts/top', true)
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
>LET doc = DOCUMENT("https://github.com/", true)
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
	"github.com/MontFerret/ferret/pkg/compiler"
	"os"
)

type Topic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func main() {
	topics, err := getTopTenTrendingTopics()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, topic := range topics {
		fmt.Println(fmt.Sprintf("%s: %s %s", topic.Name, topic.Description, topic.Url))
	}
}

func getTopTenTrendingTopics() ([]*Topic, error) {
	query := `
		LET doc = DOCUMENT("https://github.com/topics")

		FOR el IN ELEMENTS(doc, ".py-4.border-bottom")
			LIMIT 10
			LET url = ELEMENT(el, "a")
			LET name = ELEMENT(el, ".f3")
			LET desc = ELEMENT(el, ".f5")

			RETURN {
				name: TRIM(name.innerText),
				description: TRIM(desc.innerText),
				url: "https://github.com" + url.attributes.href
			}
	`

	comp := compiler.New()

	program, err := comp.Compile(query)

	if err != nil {
		return nil, err
	}

	out, err := program.Run(context.Background())

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

## Extensibility

That said, ```ferret``` is a very modular system which also allows not only embed it, but extend its standard library.

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"os"
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
		err := core.ValidateArgs(args, 1)

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

		return str.Concat(values.NewString("_ferret")).ToUpper(), nil
	}

	query := `
		FOR el IN ["foo", "bar", "qaz"]
			// conventionally all functions are registered in upper case
			RETURN TRANSFORM(el)
	`

	comp := compiler.New()
	comp.RegisterFunction("transform", transform)

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

On top of that, you can completely turn off standard library, by passing the following option:

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
