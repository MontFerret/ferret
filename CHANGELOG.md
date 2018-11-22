## Changelog

### 0.5.1

#### Fixes
- Unable to change a page load timeout [#186](https://github.com/MontFerret/ferret/pull/186).
- ``RETURN doc`` returns an empty string [#187](https://github.com/MontFerret/ferret/pull/187).
- Unable to pass an HTML Node without a selector to ``INNER_TEXT`` and ``INNER_HTML`` [#187](https://github.com/MontFerret/ferret/pull/187).
- ``doc.innerText`` returns an error [#187](https://github.com/MontFerret/ferret/pull/187).

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
