package frame

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// CallStack manages call frames and a register pool for UDF execution.
type CallStack struct {
	frames []CallFrame
	pool   Pool
}

type TraceEntry struct {
	FnName     string
	CallSitePC int
	FnID       int
}

// Init initializes the underlying register pool.
func (s *CallStack) Init(maxPoolSize int) {
	s.pool.Init(maxPoolSize)
}

// Reset clears the call stack while leaving the pool intact.
func (s *CallStack) Reset() {
	if len(s.frames) == 0 {
		return
	}

	s.frames = s.frames[:0]
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

	frame := s.frames[len(s.frames)-1]
	s.frames = s.frames[:len(s.frames)-1]
	return frame, true
}

// Top returns the top frame without removing it.
func (s *CallStack) Top() *CallFrame {
	if len(s.frames) == 0 {
		return nil
	}

	return &s.frames[len(s.frames)-1]
}

// AcquireRegisters returns a register window of the requested size.
func (s *CallStack) AcquireRegisters(size int) []runtime.Value {
	return s.pool.Get(size)
}

// ReleaseRegisters releases a register window back into the pool.
func (s *CallStack) ReleaseRegisters(reg []runtime.Value) {
	s.pool.Put(reg)
}

// ReturnToCaller unwinds one frame, restores caller registers, and applies retVal.
func (s *CallStack) ReturnToCaller(active []runtime.Value, retVal runtime.Value) ([]runtime.Value, int, bool) {
	frame, ok := s.Pop()
	if !ok {
		return nil, 0, false
	}

	// Restore caller registers and write the return value.
	s.pool.Put(active)
	registers := frame.Registers
	registers[frame.ReturnDest] = retVal

	return registers, frame.ReturnPC, true
}

// UnwindToProtectedFrame drops frames until a protected frame is reached.
func (s *CallStack) UnwindToProtectedFrame(active []runtime.Value) ([]runtime.Value, int, bool) {
	for i := len(s.frames) - 1; i >= 0; i-- {
		if !s.frames[i].Protected {
			continue
		}

		frame := s.frames[i]
		for j := i + 1; j < len(s.frames); j++ {
			s.pool.Put(s.frames[j].Registers)
		}

		// Reclaim registers above the protected frame and reset its return dest.
		s.frames = s.frames[:i]
		s.pool.Put(active)
		registers := frame.Registers
		if frame.ReturnDest.IsRegister() {
			registers[frame.ReturnDest] = runtime.None
		}

		return registers, frame.ReturnPC, true
	}

	return nil, 0, false
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
