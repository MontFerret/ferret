package optimization

import "fmt"

func resolvePassDependencies(passes []Pass, passIndex int, ctx *PassContext, metadataEpoch map[string]int, epoch int, resolving map[string]bool) error {
	pass := passes[passIndex]

	for _, dep := range pass.Requires() {
		if metadataIsFresh(ctx.Metadata, metadataEpoch, dep, epoch) {
			continue
		}

		depIndex := findEarlierPassIndex(passes, passIndex, dep)
		if depIndex < 0 {
			return fmt.Errorf("%w: pass %q requires %q which has not been executed", ErrMissingDependency, pass.Name(), dep)
		}

		if resolving == nil {
			resolving = make(map[string]bool)
		}

		if resolving[dep] {
			return fmt.Errorf("%w: cyclic dependency involving %q", ErrMissingDependency, dep)
		}

		resolving[dep] = true
		if err := resolvePassDependencies(passes, depIndex, ctx, metadataEpoch, epoch, resolving); err != nil {
			delete(resolving, dep)
			return err
		}
		delete(resolving, dep)

		if metadataIsFresh(ctx.Metadata, metadataEpoch, dep, epoch) {
			continue
		}

		dependency := passes[depIndex]
		passResult, err := dependency.Run(ctx)
		if err != nil {
			return fmt.Errorf("%w -> %s (dependency for %s): %w", ErrPassFailed, dependency.Name(), pass.Name(), err)
		}

		if passResult == nil {
			passResult = &PassResult{}
		}

		if passResult.Modified {
			return fmt.Errorf("%w: pass %q refreshed for %q", ErrDependencyRefreshModified, dependency.Name(), pass.Name())
		}

		ctx.Metadata[dependency.Name()] = passResult.Metadata
		metadataEpoch[dependency.Name()] = epoch
	}

	return nil
}

func metadataIsFresh(metadata map[string]any, metadataEpoch map[string]int, name string, epoch int) bool {
	if _, ok := metadata[name]; !ok {
		return false
	}

	return metadataEpoch[name] == epoch
}

func findEarlierPassIndex(passes []Pass, passIndex int, name string) int {
	for i := passIndex - 1; i >= 0; i-- {
		if passes[i].Name() == name {
			return i
		}
	}

	return -1
}
