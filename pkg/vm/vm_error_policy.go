package vm

// handleProtectedError applies protected-frame unwinding policy.
func (vm *VM) handleProtectedError(err error) error {
	if err == nil {
		return nil
	}

	if vm.unwindToProtected() {
		return nil
	}

	return err
}

// handleError applies catch-table then protected-frame error policy.
func (vm *VM) handleError(err error) error {
	return vm.handleErrorWithCatch(err, nil)
}

// handleErrorWithCatch applies catch-table then protected-frame error policy
// and allows a catch-specific fallback assignment/action.
func (vm *VM) handleErrorWithCatch(err error, onCatch func()) error {
	if err == nil {
		return nil
	}

	if catch, ok := vm.tryCatch(vm.pc); ok {
		if onCatch != nil {
			onCatch()
		}

		if catch[2] >= 0 {
			vm.pc = catch[2]
		}

		return nil
	}

	return vm.handleProtectedError(err)
}
