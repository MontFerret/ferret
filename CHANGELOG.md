## Changelog

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
