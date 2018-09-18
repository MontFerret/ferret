# Ferret
> Web scraping query language

## Motivation
Nowadays data is everything and who owns data - owns the world.
I have worked on multiple data-driven projects where data was an essential part of the system where I realized how cumbersome writing tons of scrapers is.
I was looking for some kind of tool that would let me to not write a code, but just express what data I need.
Unfortunately, I didn't find anything, and therefore decided to create one.
```Ferret``` project is an ambitious initiative to bring universal platform for writing scrapers without any hassle.

## Inspiration
FQL (Ferret Query Language) is heavily inspired by [AQL](https://www.arangodb.com/) (ArangoDB Query Language).
But due to domain specifics, there are some differences in how things work.

## WIP
Be aware, the the project is under heavy development. There is no documentation and some things may change in the final release.
For query syntax, you may go to [ArrangoDB web site](https://docs.arangodb.com/3.3/AQL/index.html) and use AQL docs as a docs for FQL - since they are identical.

## Quick stark

### Browserless mode

If you want to play with ```fql``` and check its syntax, run CLI with the following commands:
```
go run ./cmd/cli/main.go

```

```ferret``` will run REPL.

```shell
Welcome to Ferret REPL
Please use `Ctrl-D` to exit this program.
>LET doc = DOCUMENT('https://news.ycombinator.com/')\
>FOR post IN ELEMENTS(doc, '.storylink')\
>RETURN post.attributes.href

```

**Note:** blackslash is used for multiline queries.

If you want to execute a query store in a file, just type a file name

```
go run ./cmd/cli/main.go ./docs/examples/hackernews.fql
```


### Browser mode

By default, ``ferret`` loads HTML pages via http protocol since it's faster.
But nowadays, there are more and more websites rendered with JavaScript, and therefore, this 'old school' approach does not really work.
For this case, you may fetch documents using Chrome or Chromium via Chrome DevTools protocol (aka CDP).

```shell
go run ./cmd/cli/main.go --cdp-launch
```

**Note:** Launch command is currently broken on MacOS.

Alternatively, you may open Chrome manually with ```remote-debugging-port=9222``` arguments and bass the address to ``ferret``:

```
./bin/ferret --cdp http://127.0.0.1:9222
```

In this case, you can use function ```DOCUMENT(url, isJsRendered)``` with ```true``` for loading JS rendered pages:

```shell
Welcome to Ferret REPL
Please use `exit` or `Ctrl-D` to exit this program.
>LET doc = DOCUMENT('https://soundcloud.com/charts/top', true)
>SLEEP(2000) // WAIT WHEN THE PAGE GETS RENDERED
>LET tracks = ELEMENTS(doc, '.chartTrack__details')
>LOG("found", LENGTH(tracks), "tracks")
>FOR track IN tracks
>    LET username = ELEMENT(track, '.chartTrack__username')
>    LET title = ELEMENT(track, '.chartTrack__title')
>    RETURN {
>       artist: username.innerText,
>        track: title.innerText
>    }
```