package internal

// CompilationFrontend is the compiler assembly root.
// It owns construction and explicit wiring, but compilers must depend only on
// their named collaborators rather than the frontend itself.
type CompilationFrontend struct {
	Session         *CompilationSession
	Bindings        *BindingCompiler
	Calls           *CallResolver
	Collects        *CollectCompiler
	Dispatch        *DispatchCompiler
	Expressions     *ExprCompiler
	Literals        *LiteralCompiler
	Loops           *LoopCompiler
	Recovery        *RecoveryCompiler
	TypeFacts       *TypeFacts
	CaptureAnalyzer *CaptureAnalyzer
	Sorts           *LoopSortCompiler
	Statements      *StatementCompiler
	UDFCatalog      *UDFCatalogBuilder
	UDFs            *UDFCompiler
	Wait            *WaitCompiler
}

func NewCompilationFrontend(session *CompilationSession) *CompilationFrontend {
	front := &CompilationFrontend{
		Session:         session,
		Bindings:        NewBindingCompiler(session),
		Calls:           NewCallResolver(session),
		Expressions:     NewExprCompiler(session),
		Literals:        NewLiteralCompiler(session),
		Statements:      NewStatementCompiler(session),
		Loops:           NewLoopCompiler(session),
		Sorts:           NewLoopSortCompiler(session),
		Collects:        NewCollectCompiler(session),
		Wait:            NewWaitCompiler(session),
		Dispatch:        NewDispatchCompiler(session),
		TypeFacts:       NewTypeFacts(session),
		CaptureAnalyzer: NewCaptureAnalyzer(session),
		UDFs:            NewUDFCompiler(session),
		UDFCatalog:      NewUDFCatalogBuilder(session),
		Recovery:        NewRecoveryCompiler(session),
	}

	front.Bindings.bind(front.Expressions, front.TypeFacts)
	front.Literals.bind(front.Expressions, front.TypeFacts)
	front.CaptureAnalyzer.bind(front.Bindings)
	front.UDFCatalog.bind(front.Calls)
	front.Sorts.bind(front.Expressions, front.TypeFacts)
	front.Recovery.bind(front.Expressions, front.Literals, front.TypeFacts)
	front.Dispatch.bind(front.Expressions, front.Literals, front.Recovery, front.TypeFacts)
	front.Wait.bind(front.Expressions, front.Literals, front.Recovery, front.TypeFacts)
	front.Collects.bind(front.Bindings, front.Calls, front.Expressions, front.Recovery, front.TypeFacts)
	front.Statements.bind(front.Bindings, front.Dispatch, front.Expressions, front.Loops, front.TypeFacts, front.Wait)
	front.UDFs.bind(front.Calls, front.Expressions, front.TypeFacts, front.Recovery, front.Statements)
	front.Loops.bind(front.Bindings, front.Collects, front.Expressions, front.Literals, front.Recovery, front.Sorts, front.TypeFacts)
	front.Expressions.bind(front.Bindings, front.Calls, front.Dispatch, front.Literals, front.Loops, front.Recovery, front.TypeFacts, front.Wait)

	return front
}
