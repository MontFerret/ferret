package engine

import (
	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
)

// Option configures an Engine during construction.
type Option = gooptions.Option[config]

// WithOptimizationLevel configures optimization for normal query compilation.
// OptimizationFull is used by default; debug compilation always uses OptimizationNone.
func WithOptimizationLevel(level OptimizationLevel) Option {
	return gooptions.New(
		func(config *config, level OptimizationLevel) {
			config.optimizationLevel = level
		},
	).
		Named("optimization level").
		Validators(
			gooptions.OneOf(
				OptimizationNone,
				OptimizationBasic,
				OptimizationFull,
			),
		).
		Value(level).
		Build()
}

// WithProgramLoader sets a custom artifact loader for Engine.Load.
func WithProgramLoader(loader *artifact.Loader) Option {
	return gooptions.New(func(opts *config, loader *artifact.Loader) {
		opts.programLoader = loader
	}).
		Value(loader).
		Named("program loader").
		Validators(gooptions.NotNilPtr[artifact.Loader]()).
		Build()
}

// WithMaxActiveSessions sets an engine-wide limit on concurrently active sessions.
//
// This limit applies to Session objects created from any plan compiled by the
// engine. When the limit is reached, Plan.NewSession blocks until another
// session is closed or the provided context is canceled.
//
// Use this when you want to put a global cap on query execution concurrency and
// the host-side resources that come with it, such as CPU, memory, network
// traffic, or downstream service pressure.
//
// This is different from WithMaxIdleVMsPerPlan and WithMaxVMsPerPlan:
// WithMaxActiveSessions controls how many sessions may be running or checked
// out at once across the engine, while the VM options control how each
// individual plan manages its VM pool.
//
// A value of 0 disables the limit.
func WithMaxActiveSessions(n int) Option {
	return gooptions.New(func(opts *config, n int) {
		opts.maxActiveSessions = n
	}).
		Value(n).
		Named("max active sessions").
		Validators(gooptions.NonNegative[int]()).
		Build()
}

// WithMaxIdleVMsPerPlan sets how many closed-session VMs each plan keeps warm
// for reuse after they become idle.
//
// This is a retention setting, not a concurrency limit. It only controls how
// many unused VMs remain cached in a plan's pool after sessions are closed.
// When the idle cache is full, additional returned VMs are closed instead of
// retained.
//
// Use this when the same compiled plan is executed repeatedly and you want to
// trade some steady-state memory for faster session creation by reusing already
// initialized VMs.
//
// This is different from WithMaxVMsPerPlan:
// WithMaxIdleVMsPerPlan controls how many unused VMs stay cached, while
// WithMaxVMsPerPlan controls the maximum total number of VMs the plan may own
// at all, including both idle and currently borrowed VMs.
//
// A value of 0 disables idle retention for the plan.
func WithMaxIdleVMsPerPlan(n int) Option {
	return gooptions.New(func(opts *config, n int) {
		opts.maxIdleVMsPerPlan = n
	}).
		Value(n).
		Named("max idle VMs per plan").
		Validators(gooptions.NonNegative[int]()).
		Build()
}

// WithMaxVMsPerPlan sets a hard per-plan limit on the total number of VMs the
// plan's pool may own at one time.
//
// The total includes both idle VMs kept in the pool and VMs currently borrowed
// by active sessions created from that plan. When the limit is reached and no
// idle VM is available to reuse, session creation fails with vm.ErrPoolExhausted.
//
// Use this when you need a strict upper bound on the memory or resource
// footprint of a single hot plan, even if that plan is under heavy concurrent
// load.
//
// This is different from WithMaxActiveSessions:
// WithMaxVMsPerPlan limits VM ownership for one plan, while
// WithMaxActiveSessions limits active session concurrency across the entire
// engine.
//
// This is also different from WithMaxIdleVMsPerPlan:
// WithMaxVMsPerPlan is a hard cap, while WithMaxIdleVMsPerPlan only decides how
// many unused VMs are retained after demand drops.
//
// A value of 0 means the plan may create as many VMs as needed, subject only to
// other limits such as WithMaxActiveSessions.
func WithMaxVMsPerPlan(n int) Option {
	return gooptions.New(func(opts *config, n int) {
		opts.maxVMsPerPlan = n
	}).
		Value(n).
		Named("max VMs per plan").
		Validators(gooptions.NonNegative[int]()).
		Build()
}
