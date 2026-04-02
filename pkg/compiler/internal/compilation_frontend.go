package internal

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

	front.Bindings.front = front
	front.Collects.front = front
	front.Dispatch.front = front
	front.Expressions.front = front
	front.Literals.front = front
	front.Loops.front = front
	front.Recovery.front = front
	front.CaptureAnalyzer.front = front
	front.Sorts.front = front
	front.Statements.front = front
	front.UDFCatalog.front = front
	front.UDFs.front = front
	front.Wait.front = front

	return front
}
