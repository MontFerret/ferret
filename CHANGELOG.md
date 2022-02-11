## Changelog

### 0.16.6

### Fixed
- Invalid XPath evaluation in HTTP driver [#725](https://github.com/MontFerret/ferret/pull/725)
- Go routines leakage [#726](https://github.com/MontFerret/ferret/pull/726)

### Updated
- Small tweak in FQL Parser for FQL Formatter [#723](https://github.com/MontFerret/ferret/pull/723)

### 0.16.5

### Fixed
- Query fails if an element is not found (regression) [#722](https://github.com/MontFerret/ferret/pull/722)

### Updated
- Small tweak in FQL Parser for FQL Formatter [#723](https://github.com/MontFerret/ferret/pull/723)

### 0.16.4

### Fixed
- Fixed inability to parse custom date formats with DATE function [#720](https://github.com/MontFerret/ferret/pull/720)

### 0.16.3

### Fixed
- Panic during XPath execution by HTTP driver [#715](https://github.com/MontFerret/ferret/pull/715)

### 0.16.2

### Fixed
- Unable to use dynamic values in LIMIT clause [#706](https://github.com/MontFerret/ferret/pull/706)
- HTTP driver does not allow to override header values [#707](https://github.com/MontFerret/ferret/pull/707), [#709](https://github.com/MontFerret/ferret/pull/709)
- Cleaned up deps [#710](https://github.com/MontFerret/ferret/pull/710), [#711](https://github.com/MontFerret/ferret/pull/711)

### 0.16.1

### Fixed
- Logical precedence in ternary operator condition [#704](https://github.com/MontFerret/ferret/pull/704)

### 0.16.0

### Added
- New ``WAITFOR EVENT`` syntax [#590](https://github.com/MontFerret/ferret/pull/590)
- Support of optional chaining [#634](https://github.com/MontFerret/ferret/pull/634)
- Tracing to CDP driver [#648](https://github.com/MontFerret/ferret/pull/648)
- Support of errors suppression in function calls [#652](https://github.com/MontFerret/ferret/pull/652)
- Support of error suppression in inline expressions [#671](https://github.com/MontFerret/ferret/pull/671)
- Support of XPath selectors throughout drivers API [#657](https://github.com/MontFerret/ferret/pull/657)
- Zero-allocation in ``FOR`` loops returning ``NONE`` [#673](https://github.com/MontFerret/ferret/pull/673)
- Ignorable ``_`` variable [#673](https://github.com/MontFerret/ferret/pull/673)

### Changed
- Updated Root API [#622](https://github.com/MontFerret/ferret/pull/622)
- Increased websocket maximum buffer size in CDP driver [#648](https://github.com/MontFerret/ferret/pull/648)

### Fixed
- ``values.Parse`` does not parse int64 [#621](https://github.com/MontFerret/ferret/pull/621)
- CPU leakage [#635](https://github.com/MontFerret/ferret/pull/635), [90792bc](https://github.com/MontFerret/ferret/pull/648/commits/90792bcf3cd0b95075988aafe5d1c5072f2985bc)
- Nil pointer exception [#636](https://github.com/MontFerret/ferret/pull/636)
- Use of deprecated CDP API methods [#637](https://github.com/MontFerret/ferret/pull/637)
- HTTP driver makes multiple requests [#642](https://github.com/MontFerret/ferret/pull/642)
- Log level is ignored [aeb1247](https://github.com/MontFerret/ferret/commit/aeb1247ab34c3107b66ef76cfedde0c747904889)

#### Dependencies 
- Upgraded github.com/mafredri/cdp
- Upgraded github.com/antchfx/xpath
- Upgraded github.com/PuerkitoBio/goquery
- Upgraded github.com/rs/zerolog

#### Other
- Dropped support of Go 1.13 [0cb7623](https://github.com/MontFerret/ferret/commit/0cb7623a7fca00cc044ba3b62822b78557d96f2f)
- New Eval API in CDP driver [#651](https://github.com/MontFerret/ferret/pull/651), [#658](https://github.com/MontFerret/ferret/pull/658)

### 0.15.0
#### Added
- Support of document charset in HTTP driver [#609](https://github.com/MontFerret/ferret/pull/609)
- ``Walk`` method to FQL Parser [80c278e](https://github.com/MontFerret/ferret/commit/80c278ec6c783e29a8df12865da8208d1c148c65)
- Possibility to send keyboard events like 'Enter' or 'Shift' [#618](https://github.com/MontFerret/ferret/pull/618)

#### Changed
- Moved CLI to a separate repository [#608](https://github.com/MontFerret/ferret/pull/608)

#### Fixed
- Passing headers and cookies to HTTP driver [#614](https://github.com/MontFerret/ferret/pull/614)
- Reading property of anyonymous object [#616](https://github.com/MontFerret/ferret/pull/616)
- Clearing input text containing special characteers [#619](https://github.com/MontFerret/ferret/pull/619)

### 0.14.1
#### Fixed
- Parsing HTTP headers and cookies [#598](https://github.com/MontFerret/ferret/pull/598)
- Parsing cookie expiration datetime [#602](https://github.com/MontFerret/ferret/pull/602)

### 0.14.0
#### Added
- Support of History API [#584](https://github.com/MontFerret/ferret/pull/584)
- Support of custom http transport in HTTP driver [#586](https://github.com/MontFerret/ferret/pull/586)
- ``LIKE`` operator [#591](https://github.com/MontFerret/ferret/pull/591)
- Support of ignoring page resources [#592](https://github.com/MontFerret/ferret/pull/592)
- Support of handling non-200 status codes in HTTP driver [#593](https://github.com/MontFerret/ferret/pull/593)
- ``DOCUMENT_EXISTS`` function [#594](https://github.com/MontFerret/ferret/pull/594)

#### Fixed
- ``RAND(0,100)`` always same result [#579](https://github.com/MontFerret/ferret/pull/579)
- Element.children always returns empty array [#580](https://github.com/MontFerret/ferret/pull/580)
- Passing parameters with a nested nil structure leads to panic [#587](https://github.com/MontFerret/ferret/pull/587)

### 0.13.0
#### Added
- ``WHILE`` loop and ``ATTR_QUERY`` function [#567](https://github.com/MontFerret/ferret/pull/567)
- Support of Element.nextElementSibling and Element.previousElement [#569](https://github.com/MontFerret/ferret/pull/569)
- Support of Element.getParentElement [#571](https://github.com/MontFerret/ferret/pull/571)
- Support of computed styles [#570](https://github.com/MontFerret/ferret/pull/570)

#### Fixed
- HTML escaping [#573](https://github.com/MontFerret/ferret/pull/573)

#### Updated 
- Upgraded CDP client [#536](https://github.com/MontFerret/ferret/pull/563)
- Upgraded GoQuery [#562](https://github.com/MontFerret/ferret/pull/562)
- Upgraded XPath [#572](https://github.com/MontFerret/ferret/pull/572)

### 0.12.1
#### Fixed
- Missing regexp FILTER operator [#558](https://github.com/MontFerret/ferret/pull/558)
- Open tabs on page load error [#564](https://github.com/MontFerret/ferret/pull/564)
- Docs for WAIT_NAVIGATION [#557](https://github.com/MontFerret/ferret/pull/557)

### 0.12.0
#### Added
- iFrame navigation handling [#535](https://github.com/MontFerret/ferret/pull/535)
- ``FRAMES`` function for fast frame lookup [#535](https://github.com/MontFerret/ferret/pull/535)
- Assertion library [#526](https://github.com/MontFerret/ferret/pull/526)

#### Changed
- Removed property caching and tracking [#531](https://github.com/MontFerret/ferret/pull/531)
- Updated dependencies [#528](https://github.com/MontFerret/ferret/pull/528), [#525](https://github.com/MontFerret/ferret/pull/525)
- ``IO::FS::WRITE`` accepts any type as a file content [#544](https://github.com/MontFerret/ferret/pull/544)
- Print errors on stderr [#539](https://github.com/MontFerret/ferret/pull/539)

#### Fixed
- ``WAIT`` does not respect cancellation signal [#524](https://github.com/MontFerret/ferret/pull/524)
- Missed ``DATE_COMPARE`` [#537](https://github.com/MontFerret/ferret/pull/537)
- Spelling [#534](https://github.com/MontFerret/ferret/pull/534)
- ``SCREENSHOT`` param type check [#545](https://github.com/MontFerret/ferret/pull/545)
- Wrong base for int formatter [e283722](https://github.com/MontFerret/ferret/commit/e283722d37f392f755ace2a42232c0d4b37d1838)

### 0.11.1
#### Fixed
- Fixed use of unquoted scroll options [#521](https://github.com/MontFerret/ferret/pull/521)
- Upgraded ANTLR version [#517](https://github.com/MontFerret/ferret/pull/517)


### 0.11.0
#### Added
- USE statement. [#470](https://github.com/MontFerret/ferret/pull/470)
- Scroll options. [#471](https://github.com/MontFerret/ferret/pull/471)
- Functions for working with file paths. [#505](https://github.com/MontFerret/ferret/pull/505)
- Fuzzer. [#501](https://github.com/MontFerret/ferret/pull/501)

## Updated
- ``DECODED_URI_COMPONENT`` decodes unicode symbols now. [#499](https://github.com/MontFerret/ferret/pull/499) 
- Dependencies. [87265cf](https://github.com/MontFerret/ferret/commit/87265cf470c4b614d144706020729dd453620a0c)

# Fixed
- ``RAND`` always returns same result . [#484](https://github.com/MontFerret/ferret/pull/484)
- ``RAND`` does not work on Windows. [#497](https://github.com/MontFerret/ferret/pull/497)
- ``IO::FS::WRITE`` does not add read permissions. [#494](https://github.com/MontFerret/ferret/pull/494)
- Unable to use keywords in namespaces. [#481](https://github.com/MontFerret/ferret/pull/481)

### 0.10.2
#### Updated
- Updated dependencies. [#466](https://github.com/MontFerret/ferret/pull/466) [#467](https://github.com/MontFerret/ferret/pull/467)

### 0.10.1
#### Fixed
- Added string functions with correct names. [#461](https://github.com/MontFerret/ferret/pull/461)
- Added missed datetime library. [#462](https://github.com/MontFerret/ferret/pull/462)

### 0.10.0
#### Added
- Response information to drivers. [#391](https://github.com/MontFerret/ferret/pull/391), [#450](https://github.com/MontFerret/ferret/pull/450)
- Compilation check whether parameter values are provided. [#396](https://github.com/MontFerret/ferret/pull/396)
- Allowed HTTP response codes to HTTP driver. [#398](https://github.com/MontFerret/ferret/pull/398)
- IO functions to standard library. [#403](https://github.com/MontFerret/ferret/pull/403), [#405](https://github.com/MontFerret/ferret/pull/405), [#452](https://github.com/MontFerret/ferret/pull/452)
- Compilation check whether a variable name is unique. [#416](https://github.com/MontFerret/ferret/pull/416)
- Loading HTML page into memory. Supported by all drivers. [#413](https://github.com/MontFerret/ferret/pull/434)

#### Fixed
- Fixes in HTTP driver. [#390](https://github.com/MontFerret/ferret/pull/390)
- Inability to handle redirects correctly. [#432](https://github.com/MontFerret/ferret/pull/432)
- XPath selector gives faulty output. [#435](https://github.com/MontFerret/ferret/pull/435)
- Typos in README and comments. [#446](https://github.com/MontFerret/ferret/pull/446)
- ``PAGINATION`` fails during redirects. [#448](https://github.com/MontFerret/ferret/pull/448)

#### Changed
- Made FQL keywords case insensitive. [#393](https://github.com/MontFerret/ferret/pull/393)
- Performance boost in EventBroker. [#402](https://github.com/MontFerret/ferret/pull/402), [#407](https://github.com/MontFerret/ferret/pull/407), [#408](https://github.com/MontFerret/ferret/pull/408)
- Updated dependencies.


### 0.9.0
#### Added
- ``INPUT_CLEAR`` function to clear input's value. [#366](https://github.com/MontFerret/ferret/pull/366)
- Support of tick for string literals. [#367](https://github.com/MontFerret/ferret/pull/367)
- Support of default headers and cookies. [#372](https://github.com/MontFerret/ferret/pull/372)
- Support of use of params in dot notation. [#378](https://github.com/MontFerret/ferret/pull/378)
- Optional count param to ``CLICK`` function. [#377](https://github.com/MontFerret/ferret/pull/377)
- ``BLUR`` function. [#379](https://github.com/MontFerret/ferret/pull/379)

#### Fixed
- Tabs don't get closed on page load error. [#359](https://github.com/MontFerret/ferret/pull/359)
- ``CLICK`` function does not allow to use element with a selector. [#355](https://github.com/MontFerret/ferret/pull/355)
- Unable to use member expression right after a function call. [#368](https://github.com/MontFerret/ferret/pull/368)

#### Changed
- Updated zerolog. [#352](https://github.com/MontFerret/ferret/pull/352)
- Runtime ``Object`` and ``Array`` values implement ``core.Getter`` interface. [#353](https://github.com/MontFerret/ferret/pull/353)
- Externalized default timeout values. [#371](https://github.com/MontFerret/ferret/pull/371) 
- Refactored ``drivers.HTMLDocument`` and ``drivers.HTMLElement`` interfaces. [#376](https://github.com/MontFerret/ferret/pull/376), [#375](https://github.com/MontFerret/ferret/pull/375)

### 0.8.3
#### Fixed
- Unable to click by selector using an element.

### 0.8.2
#### Fixed
- Scrolling position is not centered. [#343](https://github.com/MontFerret/ferret/pull/343)
- Unable to set custom logger fields. [#346](https://github.com/MontFerret/ferret/pull/346)
- Fixed ``INNER_HTML``, ``INNER_TEXT``, ``INNER_HTML_SET``, ``INNER_TEXT_SET`` functions. [#347](https://github.com/MontFerret/ferret/pull/347)
- Unable to set custom headers. [#348](https://github.com/MontFerret/ferret/pull/348)

### 0.8.1
#### Fixed
- Added existence check to ``CLICK`` and ``CLICK_ALL`` functions. [#341](https://github.com/MontFerret/ferret/pull/341)
- Added a check whether an element is in the viewport before scrolling. [#342](https://github.com/MontFerret/ferret/pull/342)

### 0.8.0
#### Added
- Delay randomization for inputs. [#283](https://github.com/MontFerret/ferret/pull/283)
- Namespace support. [#269](https://github.com/MontFerret/ferret/pull/296)
- iframe support. [#315](https://github.com/MontFerret/ferret/pull/315)
- Better emulation of user interaction. [#316](https://github.com/MontFerret/ferret/pull/316), [#331](https://github.com/MontFerret/ferret/pull/331)
- ``ESCAPE_HTML``, ``UNESCAPE_HTML`` and ``DECODE_URI_COMPONENT`` functions. [#318](https://github.com/MontFerret/ferret/pull/318)
- XPath support. [#322](https://github.com/MontFerret/ferret/pull/322)
- Regular expression operator. [#326](https://github.com/MontFerret/ferret/pull/326)
- ``INNER_HTML_SET`` and ``INNER_TEXT_SET`` functions. [#329](https://github.com/MontFerret/ferret/pull/329)
- Possibility to set viewport size. [#334](https://github.com/MontFerret/ferret/pull/334)
- ``FOCUS`` function. [#340](https://github.com/MontFerret/ferret/pull/340)

#### Changed
- ``RAND`` accepts optional upper and lower limits. [#271](https://github.com/MontFerret/ferret/pull/271)
- Updated CDP definitions. [#328](https://github.com/MontFerret/ferret/pull/328) 
- Logic of iterator termination. [#330](https://github.com/MontFerret/ferret/pull/330)

#### Fixed
- Order of arguments in ``SCROLL`` function. [#269](https://github.com/MontFerret/ferret/pull/269)
- The command line parameter "--param" does not support colon. [#282](https://github.com/MontFerret/ferret/pull/282)
- Race condition during ``WAIT_NAVIGATION`` call. [#281](https://github.com/MontFerret/ferret/pull/281)
- Arithmetic operators. [#298](https://github.com/MontFerret/ferret/pull/298)
- Missed UA setting for HTTP driver. [#318](https://github.com/MontFerret/ferret/pull/318)
- Improper math operator used in calculating page load timeout. [#319](https://github.com/MontFerret/ferret/pull/319)
- Wrong function names in README. [#321](https://github.com/MontFerret/ferret/pull/321)
- JSON serialization for HTTPHeader type. [#323](https://github.com/MontFerret/ferret/pull/323)

### 0.7.0
#### Added
- Autocomplete to CLI [#219](https://github.com/MontFerret/ferret/pull/219).
- New mouse functions - ``MOUSE(x, y)`` and ``SCROLL(x, y)`` [#237](https://github.com/MontFerret/ferret/pull/237).
- ``WAIT_NO_ELEMENT``, ``WAIT_NO_CLASS`` and ``WAIT_NO_CLASS_ALL`` functions [#249](https://github.com/MontFerret/ferret/pull/249).
- Computed ``HTMLElement.style`` property [#255](https://github.com/MontFerret/ferret/pull/255).
- ``ATTR_GET``, ``ATTR_SET``, ``ATTR_REMOVE``, ``STYLE_GET``, ``STYLE_SET`` and ``STYLE_REMOVE`` functions [#255](https://github.com/MontFerret/ferret/pull/255).
- ``WAIT_STYLE``, ``WAIT_NO_STYLE``, ``WAIT_STYLE_ALL`` and ``WAIT_NO_STYLE_ALL`` functions [#256](https://github.com/MontFerret/ferret/pull/260).
- Cookies support. Now a document can be loaded with preset cookies. Also, HTMLDocument has ``.cookies`` property.
In order to manipulate with cookies, ``COOKIE_DEL``, ``COOKIE_SET`` AND ``COOKIE_GET`` functions were added [#242](https://github.com/MontFerret/ferret/pull/242).

```
LET doc = DOCUMENT(url, {
    driver: "cdp",
    cookies: [{
        name: "x-e2e",
        value: "test"
    }, {
        name: "x-e2e-2",
        value: "test2"
    }]
})
```

#### Changed
- Renamed ParseTYPEP to MustParseTYPE [#231](https://github.com/MontFerret/ferret/pull/231).
- Added context to all HTML object [#235](https://github.com/MontFerret/ferret/pull/235).

#### Fixed
- Click events are not cancellable [#222](https://github.com/MontFerret/ferret/pull/222).
- Name collision [#223](https://github.com/MontFerret/ferret/pull/223).
- Invalid return in FQL Compiler constructor [#227](https://github.com/MontFerret/ferret/pull/227).
- Incorrect string length computation [#238](https://github.com/MontFerret/ferret/pull/238).
- Access to HTML object properties via dot notation [#239](https://github.com/MontFerret/ferret/pull/239).
- Graceful process termination [#240](https://github.com/MontFerret/ferret/pull/240).
- Browser launcher for macOS [#246](https://github.com/MontFerret/ferret/pull/246). 

#### Breaking changes
- New runtime type system [#232](https://github.com/MontFerret/ferret/pull/232).
- Moved and renamed ``collections.IterableCollection`` and ```collections.CollectionIterator``` interfaces.
Now they are in ``core`` package and called ``Iterable`` and ``Iterator`` [1af8b37](https://github.com/MontFerret/ferret/commit/f8e061cc8034fd4cfa4ce2a094276d50137a4b98).
- Renamed ``collections.Collection`` interface to ``collections.Measurable`` [1af8b37](https://github.com/MontFerret/ferret/commit/f8e061cc8034fd4cfa4ce2a094276d50137a4b98).
- Moved html interfaces from ``runtime/values`` package into ``drivers`` package [#234](https://github.com/MontFerret/ferret/pull/234).
- Changed drivers initialization. Replaced old ``drivers.WithDynamic`` and ``drivers.WithStatic`` methods with a new ``drivers.WithContext`` method with optional parameter ``drivers.AsDefault()`` [#234](https://github.com/MontFerret/ferret/pull/234).
- New document load params [#234](https://github.com/MontFerret/ferret/pull/234).
```
LET doc = DOCUMENT(url, {
    driver: "cdp"
})
```


### 0.6.0
#### Added
- Added support for ```context.Done()``` to interrupt an execution [#201](https://github.com/MontFerret/ferret/pull/201).
- Added support for custom HTML drivers [#209](https://github.com/MontFerret/ferret/pull/209).
- Added support for dot notation access and assignments for custom types [#214](https://github.com/MontFerret/ferret/pull/214/commits/0ea36e511540e569ef53b8748301512b6d8a046b)
- Added ```ELEMENT_EXISTS(doc, selector) -> Boolean``` function [#210](https://github.com/MontFerret/ferret/pull/210).
```
LET exists = ELEMENT_EXISTS(doc, ".nav")
```
- Added ```PageLoadParams``` to ```DOCUMENT``` function [#214](https://github.com/MontFerret/ferret/pull/214/commits/3434323cd08ca3186e90cb5ab1faa26e28a28709).
```
LET doc = DOCUMENT("https://www.google.com/", {
    dynamic: true,
    timeout: 10000
})
```
 
#### Fixed
- Math operators precedence [#202](https://github.com/MontFerret/ferret/pull/202).
- Memory leak in ```DOWNLOAD``` function [#213](https://github.com/MontFerret/ferret/pull/213).

#### Breaking change
- **(Embedded)** Removed builtin drivers initialization in Program [#198](https://github.com/MontFerret/ferret/pull/198).
The initialization must be done via context manually.

### 0.5.2
#### Fixed
- Does not close browser tab when fails to load a page [#193](https://github.com/MontFerret/ferret/pull/193).
- ```HTMLElement.value``` does not return actual value [#195](https://github.com/MontFerret/ferret/pull/195)
- Compiles a query with duplicate variable in FOR statement [#196](https://github.com/MontFerret/ferret/pull/196)
- Default CDP address [#197](https://github.com/MontFerret/ferret/pull/197).  

### 0.5.1
#### Fixed
- Unable to change a page load timeout [#186](https://github.com/MontFerret/ferret/pull/186).
- ``RETURN doc`` returns an empty string [#187](https://github.com/MontFerret/ferret/pull/187).
- Unable to pass an HTML Node without a selector to ``INNER_TEXT`` and ``INNER_HTML`` [#187](https://github.com/MontFerret/ferret/pull/187).
- ``doc.innerText`` returns an error [#187](https://github.com/MontFerret/ferret/pull/187).
- Panics when ``WAIT_CLASS`` does not receive all required arguments [#192](https://github.com/MontFerret/ferret/pull/192).

### 0.5.0
#### Added
- ``FMT`` function [#151](https://github.com/MontFerret/ferret/pull/151).
- DateTime functions [#152](https://github.com/MontFerret/ferret/pull/152), [#153](https://github.com/MontFerret/ferret/pull/153), [#154](https://github.com/MontFerret/ferret/pull/154), [#156](https://github.com/MontFerret/ferret/pull/156), [#157](https://github.com/MontFerret/ferret/pull/157), [#165](https://github.com/MontFerret/ferret/pull/165), [#175](https://github.com/MontFerret/ferret/pull/175), [#182](https://github.com/MontFerret/ferret/pull/182).
- ``PAGINATION`` function [#173](https://github.com/MontFerret/ferret/pull/173).
- ``SCROLL_TOP``, ``SCROLL_BOTTOM`` and ``SCROLL_ELEMENT`` functions [#174](https://github.com/MontFerret/ferret/pull/174).
- ``HOVER`` function [#178](https://github.com/MontFerret/ferret/pull/178).
- Panic recovery mechanism [#158](https://github.com/MontFerret/ferret/pull/158).

#### Fixed
- Unable to define variables and make function calls before FILTER, SORT and etc statements [#148](https://github.com/MontFerret/ferret/pull/148).
- Unable to use params in LIMIT clause [#173](https://github.com/MontFerret/ferret/pull/173).
- ```RIGHT``` should return substr counting from right rather than left [#164](https://github.com/MontFerret/ferret/pull/164).
- ``INNER_HTML`` returns outer HTML instead for dynamic elements [#170](https://github.com/MontFerret/ferret/pull/170).
- ``INNER_TEXT`` returns HTML instead from dynamic elements [#170](https://github.com/MontFerret/ferret/pull/170).

#### Breaking change:
- Name collision between ```math``` and ```utils``` packages in standard library. Renamed ```LOG``` to ```PRINT``` [#162](https://github.com/MontFerret/ferret/pull/162).

### 0.4.0
#### Added
- ``COLLECT`` keyword [#141](https://github.com/MontFerret/ferret/pull/141)
- ``VALUES`` function [#128](https://github.com/MontFerret/ferret/pull/128) 
- ``MERGE_RECURSIVE`` function [#140](https://github.com/MontFerret/ferret/pull/140) 

#### Fixed
- Unable to use string literals as object properties [commit](https://github.com/MontFerret/ferret/commit/685c5872aaed42852ce32e7ab8b69b1a269185be)

### 0.3.0

#### Added
- ``FROM_BASE64`` function [commit](https://github.com/MontFerret/ferret/commit/5db8df55db46336927ca32ab096569fa09df58d3)
- Support for multi line strings [commit](https://github.com/MontFerret/ferret/commit/cf70088fd84fa0e02887c0f34298793b98f96073)
- ``DOWNLOAD`` function [commit](https://github.com/MontFerret/ferret/commit/dd13878f80f340c4727d3ad5a6a70859dd958b92)
- Binary expressions [commit](https://github.com/MontFerret/ferret/commit/e5ca63bcdb83418b40792bc65bf83f58a0cb1b4e)

#### Fixed
- ``KEEP`` function does not perform deep cloning [commit](https://github.com/MontFerret/ferret/commit/0f3128e8428cd3dc5377a2ead3134c1ae14cc9a0)
- WaitForNavigation callback can get called more than once [commit](https://github.com/MontFerret/ferret/commit/1d6a23fa967643a737cd052234d480052d3ec2d9)
- Concurrent map iteration and map write  [commit](https://github.com/MontFerret/ferret/commit/1d6a23fa967643a737cd052234d480052d3ec2d9)

#### Breaking changes
- Renamed ``.innerHtml`` to ``.innerHTML`` [commit](https://github.com/MontFerret/ferret/commit/393980029976405d9e432faadd407e964c995fd4)

### 0.2.0

#### Added
- Numeric functions [commit](https://github.com/MontFerret/ferret/commit/5f94b77a39709846a922a3bf421f81e78c2b0c7e)
- ``PDF`` function [commit](https://github.com/MontFerret/ferret/commit/2417be3f9da6db49dcee5ac6f061cc66142fbef5)
- ``ZIP`` function [commit](https://github.com/MontFerret/ferret/commit/5d0d9ec5374d42b0e882436955666c737d9dab0c)
- ``MERGE`` function [commit](https://github.com/MontFerret/ferret/commit/446ce3ead5812fe105726bae16196fb7ce4a7185)
