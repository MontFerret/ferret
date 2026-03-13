package frame

// CallStack manages structural call frames for UDF execution.
type (
	CallStack struct {
		frames []CallFrame
	}

	TraceEntry struct {
		FnName     string
		CallSitePC int
		FnID       int
	}
)

func NewCallStack() CallStack {
	return CallStack{}
}

// Len returns the number of active frames.
func (s *CallStack) Len() int {
	return len(s.frames)
}

// Push adds a new frame to the stack.
func (s *CallStack) Push(frame CallFrame) {
	s.frames = append(s.frames, frame)
}

// Pop removes the top frame from the stack.
func (s *CallStack) Pop() (CallFrame, bool) {
	if len(s.frames) == 0 {
		return CallFrame{}, false
	}

	top := len(s.frames) - 1
	frame := s.frames[top]
	s.frames[top] = CallFrame{}
	s.frames = s.frames[:top]
	return frame, true
}

// Top returns the top frame without removing it.
func (s *CallStack) Top() *CallFrame {
	if len(s.frames) == 0 {
		return nil
	}

	return &s.frames[len(s.frames)-1]
}

// NearestRecoveryBoundary returns the index of the nearest protected frame.
func (s *CallStack) NearestRecoveryBoundary() int {
	for i := len(s.frames) - 1; i >= 0; i-- {
		if s.frames[i].RecoveryBoundary {
			return i
		}
	}

	return -1
}

// SetTopFnID updates the top frame's function id when present.
func (s *CallStack) SetTopFnID(fnID int) bool {
	if len(s.frames) == 0 {
		return false
	}

	s.frames[len(s.frames)-1].FnID = fnID
	return true
}

// SetTopCall updates call metadata of the top frame when present.
func (s *CallStack) SetTopCall(fnID int, fnName string, callSitePC int) bool {
	if len(s.frames) == 0 {
		return false
	}

	top := &s.frames[len(s.frames)-1]
	top.FnID = fnID
	top.FnName = fnName
	top.CallSitePC = callSitePC
	top.HasCallSite = true

	return true
}

// TraceEntries returns caller trace entries from nearest to farthest frame.
func (s *CallStack) TraceEntries() []TraceEntry {
	if len(s.frames) == 0 {
		return nil
	}

	traces := make([]TraceEntry, 0, len(s.frames))

	for i := len(s.frames) - 1; i >= 0; i-- {
		frame := s.frames[i]
		if !frame.HasCallSite {
			continue
		}

		traces = append(traces, TraceEntry{
			CallSitePC: frame.CallSitePC,
			FnID:       frame.FnID,
			FnName:     frame.FnName,
		})
	}

	if len(traces) == 0 {
		return nil
	}

	return traces
}

// CallSitePCs returns caller PCs from nearest to farthest frame.
func (s *CallStack) CallSitePCs() []int {
	traces := s.TraceEntries()
	if len(traces) == 0 {
		return nil
	}

	pcs := make([]int, 0, len(traces))
	for _, trace := range traces {
		pcs = append(pcs, trace.CallSitePC)
	}

	return pcs
}
