package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestMember(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase("LET arr = [1,2,3,4] RETURN arr[10]", BC{
			I(vm.OpLoadConst, 1, C(0)),
			I(vm.OpMove, 2, C(1)),
			I(vm.OpLoadConst, 3, C(2)),
			I(vm.OpMove, 4, C(3)),
			I(vm.OpLoadConst, 5, C(4)),
			I(vm.OpMove, 6, C(5)),
			I(vm.OpLoadArray, 7, R(2), R(4), R(6)),
			I(vm.OpMove, 0, 7),
			I(vm.OpReturn, 0, 7),
		}),
		//Case("LET arr = [1,2,3,4] RETURN arr[1]", 2),
		//Case("LET arr = [1,2,3,4] LET idx = 1 RETURN arr[idx]", 2),
		//Case(`LET obj = { foo: "bar", qaz: "wsx"} RETURN obj["qaz"]`, "wsx"),
		//Case(fmt.Sprintf(`
		//						LET obj = { "foo": "bar", %s: "wsx"}
		//
		//						RETURN obj["qaz"]
		//					`, "`qaz`"), "wsx"),
		//Case(fmt.Sprintf(`
		//						LET obj = { "foo": "bar", %s: "wsx"}
		//
		//						RETURN obj["let"]
		//					`, "`let`"),
		//	"wsx"),
		//Case(`LET obj = { foo: "bar", qaz: "wsx"} LET key = "qaz" RETURN obj[key]`, "wsx"),
		//Case(`RETURN { foo: "bar" }.foo`, "bar"),
		//Case(`LET inexp = 1 IN {'foo': [1]}.foo
		//	LET ternaryexp = FALSE ? TRUE : {foo: TRUE}.foo
		//	RETURN inexp && ternaryexp`,
		//	true),
		//Case(`RETURN ["bar", "foo"][0]`, "bar"),
		//Case(`LET inexp = 1 IN [[1]][0]
		//						LET ternaryexp = FALSE ? TRUE : [TRUE][0]
		//						RETURN inexp && ternaryexp`,
		//	true),
		//Case(`LET obj = {
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
		//Case(`LET o1 = {
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
		//Case(`LET o1 = {
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
		//Case(`LET obj = {
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
		//Case(`RETURN FIRST([[1, 2]][0])`,
		//	1),
		//CaseArray(`RETURN [[1, 2]][0]`,
		//	[]any{1, 2}),
		//CaseArray(`
		//			FOR i IN [[1, 2]][0]
		//				RETURN i
		//		`,
		//	[]any{1, 2}),
		//Case(`
		//			LET arr = [{ name: "Bob" }]
		//
		//			RETURN FIRST(arr).name
		//		`,
		//	"Bob"),
		ByteCodeCase(`
					LET arr = [{ name: { first: "Bob" } }]
		
					RETURN FIRST(arr)['name'].first
				`,
			BC{
				I(vm.OpLoadConst, 1, C(0)),
			}),
		//CaseNil(`
		//			LET obj = { foo: None }
		//
		//			RETURN obj.foo?.bar
		//		`),
	})
}
