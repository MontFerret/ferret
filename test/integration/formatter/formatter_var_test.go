package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterVarBindings(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
	var    foo =      10
	 foo    =   foo +   1
	 return foo
		`, `var foo = 10
foo = foo + 1
return foo`+"\n\n"),
		S(`
	var total=10
	total+=1
	total-=2
	total*=3
	total/=3
	return total
		`, `var total = 10
total += 1
total -= 2
total *= 3
total /= 3
return total`+"\n\n"),
		S(`
	func   run( ){
	var total= 1
	 total   =total+2
 return total
}
return run()
		`, `func run() {
    var total = 1
    total = total + 2
    return total
}
return run()`+"\n\n"),
		S(`
	func run(){
	var total=10
	total+=1
	total-=2
	total*=3
	total/=3
	return total
	}
	return run()
		`, `func run() {
    var total = 10
    total += 1
    total -= 2
    total *= 3
    total /= 3
    return total
}
return run()`+"\n\n"),
		S(`
let    STEP =  10
return STEP
`, `let STEP = 10
return STEP`+"\n\n"),
		S(`
for item in [ 1, 2 ]
var current = item
current=current+1
return current
`, `for item in [1, 2]
    var current = item
    current = current + 1
    return current`+"\n\n"),
		S(`
let item={ deprecated: true, keep: true}
delete    item.deprecated
return item
`, `let item = { deprecated: true, keep: true }
delete item.deprecated
return item`+"\n\n"),
		S(`
func clean(payload){
delete  payload["debug"]
return payload
}
return clean({})
`, `func clean(payload) {
    delete payload["debug"]
    return payload
}
return clean({})`+"\n\n"),
		S(`
for item in [{ stale: true }]
delete item.stale
return item
`, `for item in [
    { stale: true }
]
    delete item.stale
    return item`+"\n\n"),
	})
}
