package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestMember(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		//SkipByteCodeCase("LET arr = [1,2,3,4] RETURN arr[10]", BC{
		//	I(bytecode.OpLoadConst, 1, C(0)),
		//	I(bytecode.OpMove, 2, C(1)),
		//	I(bytecode.OpLoadConst, 3, C(2)),
		//	I(bytecode.OpMove, 4, C(3)),
		//	I(bytecode.OpLoadConst, 5, C(4)),
		//	I(bytecode.OpMove, 6, C(5)),
		//	I(bytecode.OpLoadArray, 7, R(2), R(4), R(6)),
		//	I(bytecode.OpMove, 0, 7),
		//	I(bytecode.OpReturn, 0, 7),
		//}),
		//Spec("LET arr = [1,2,3,4] RETURN arr[1]", 2),
		//Spec("LET arr = [1,2,3,4] LET idx = 1 RETURN arr[idx]", 2),
		//Spec(`LET obj = { foo: "bar", qaz: "wsx"} RETURN obj["qaz"]`, "wsx"),
		//Spec(fmt.Sprintf(`
		//						LET obj = { "foo": "bar", %s: "wsx"}
		//
		//						RETURN obj["qaz"]
		//					`, "`qaz`"), "wsx"),
		//Spec(fmt.Sprintf(`
		//						LET obj = { "foo": "bar", %s: "wsx"}
		//
		//						RETURN obj["let"]
		//					`, "`let`"),
		//	"wsx"),
		//Spec(`LET obj = { foo: "bar", qaz: "wsx"} LET key = "qaz" RETURN obj[key]`, "wsx"),
		//Spec(`RETURN { foo: "bar" }.foo`, "bar"),
		//Spec(`LET inexp = 1 IN {'foo': [1]}.foo
		//	LET ternaryexp = FALSE ? TRUE : {foo: TRUE}.foo
		//	RETURN inexp && ternaryexp`,
		//	true),
		//Spec(`RETURN ["bar", "foo"][0]`, "bar"),
		//Spec(`LET inexp = 1 IN [[1]][0]
		//						LET ternaryexp = FALSE ? TRUE : [TRUE][0]
		//						RETURN inexp && ternaryexp`,
		//	true),
		//Spec(`LET obj = {
		//					first: {
		//						second: {
		//							third: {
		//								fourth: {
		//									fifth: {
		//										bottom: true
		//									}
		//								}
		//							}
		//						}
		//					}
		//				}
		//
		//				RETURN obj.first.second.third.fourth.fifth.bottom`,
		//	true),
		//Spec(`LET o1 = {
		//first: {
		//  second: {
		//      ["third"]: {
		//          fourth: {
		//              fifth: {
		//                  bottom: true
		//              }
		//          }
		//      }
		//  }
		//}
		//}
		//
		//LET o2 = { prop: "third" }
		//
		//RETURN o1["first"]["second"][o2.prop]["fourth"]["fifth"].bottom`,
		//
		//	true),
		//Spec(`LET o1 = {
		//first: {
		// second: {
		//     third: {
		//         fourth: {
		//             fifth: {
		//                 bottom: true
		//             }
		//         }
		//     }
		// }
		//}
		//}
		//
		//LET o2 = { prop: "third" }
		//
		//RETURN o1.first["second"][o2.prop].fourth["fifth"]["bottom"]`,
		//
		//	true),
		//Spec(`LET obj = {
		//					attributes: {
		//						'data-index': 1
		//					}
		//				}
		//
		//				RETURN obj.attributes['data-index']`,
		//	1),
		//CaseRuntimeError(`LET obj = NONE RETURN obj.foo`),
		//CaseNil(`LET obj = NONE RETURN obj?.foo`),
		//CaseObject(`RETURN {first: {second: "third"}}.first`,
		//	map[string]any{
		//		"second": "third",
		//	}),
		//SkipCaseObject(`RETURN KEEP_KEYS({first: {second: "third"}}.first, "second")`,
		//	map[string]any{
		//		"second": "third",
		//	}),
		//CaseArray(`
		//			FOR v, k IN {f: {foo: "bar"}}.f
		//				RETURN [k, v]
		//		`,
		//	[]any{
		//		[]any{"foo", "bar"},
		//	}),
		//Spec(`RETURN FIRST([[1, 2]][0])`,
		//	1),
		//CaseArray(`RETURN [[1, 2]][0]`,
		//	[]any{1, 2}),
		//CaseArray(`
		//			FOR i IN [[1, 2]][0]
		//				RETURN i
		//		`,
		//	[]any{1, 2}),
		//Spec(`
		//			LET arr = [{ name: "Bob" }]
		//
		//			RETURN FIRST(arr).name
		//		`,
		//	"Bob"),
		//SkipByteCodeCase(`
		//			LET arr = [{ name: { first: "Bob" } }]
		//
		//			RETURN FIRST(arr)['name'].first
		//		`,
		//	BC{
		//		I(bytecode.OpLoadConst, 1, C(0)),
		//	}),
		//CaseNil(`
		//			LET obj = { foo: None }
		//
		//			RETURN obj.foo?.bar
		//		`),
	})
}
