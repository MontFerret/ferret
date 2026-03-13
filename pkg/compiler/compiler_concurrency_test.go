package compiler

import (
	"fmt"
	"maps"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestCompilerCompileConcurrentSharedCompiler(t *testing.T) {
	t.Parallel()

	compiler := New(WithOptimizationLevel(O1))
	workers := maxInt(8, runtime.GOMAXPROCS(0)*2)
	iterations := 80

	validQueries := []string{
		"RETURN 1",
		"LET payload = { value: 42, tags: [1, 2, 3] }\nRETURN payload.value",
		"FOR item IN [1, 2, 3, 4] RETURN item * 2",
		"LET data = [{ name: 'a' }, { name: 'b' }]\nFOR row IN data RETURN row.name",
	}

	t.Run("shared_sources", func(t *testing.T) {
		sources := make([]*file.Source, 0, len(validQueries))
		for i, query := range validQueries {
			sources = append(sources, file.NewSource(fmt.Sprintf("shared_%d", i), query))
		}

		runConcurrentCompileWorkers(t, workers, iterations, func(worker, iter int) error {
			source := sources[(worker+iter)%len(sources)]

			program, err := compiler.Compile(source)
			if err != nil {
				return fmt.Errorf("compile failed: %w", err)
			}

			if program == nil {
				return fmt.Errorf("program is nil")
			}

			if program.Source != source {
				return fmt.Errorf("unexpected source pointer in output program")
			}

			if len(program.Bytecode) == 0 {
				return fmt.Errorf("compiled program has no bytecode")
			}

			return nil
		})
	})

	t.Run("fresh_sources", func(t *testing.T) {
		runConcurrentCompileWorkers(t, workers, iterations, func(worker, iter int) error {
			query := validQueries[(worker+iter)%len(validQueries)]
			source := file.NewSource(fmt.Sprintf("fresh_%d_%d", worker, iter), query)

			program, err := compiler.Compile(source)
			if err != nil {
				return fmt.Errorf("compile failed: %w", err)
			}

			if program == nil {
				return fmt.Errorf("program is nil")
			}

			if program.Source != source {
				return fmt.Errorf("unexpected source pointer in output program")
			}

			if len(program.Bytecode) == 0 {
				return fmt.Errorf("compiled program has no bytecode")
			}

			return nil
		})
	})

	udfQueries := []struct {
		expectedHost map[string]int
		name         string
		query        string
		expectedUDFs int
	}{
		{
			name: "udf_alias_function_name",
			query: `
USE FOO::TEST_FN AS FN
FUNC WRAP(v) => FN(v)
RETURN WRAP(1)
`,
			expectedHost: map[string]int{"FOO::TEST_FN": 1},
			expectedUDFs: 1,
		},
		{
			name: "udf_namespace_alias",
			query: `
USE BAR AS B
FUNC RUN(v) => B::OTHER_FN(v, v + 1)
RETURN RUN(2)
`,
			expectedHost: map[string]int{"BAR::OTHER_FN": 2},
			expectedUDFs: 1,
		},
		{
			name: "udf_case_distinct_hosts",
			query: `
FUNC wrap() => Foo(1) + foo(2)
RETURN wrap()
`,
			expectedHost: map[string]int{"Foo": 1, "foo": 1},
			expectedUDFs: 1,
		},
	}

	t.Run("udf_isolation_shared_sources", func(t *testing.T) {
		type sharedSourceCase struct {
			source *file.Source
			spec   struct {
				expectedHost map[string]int
				expectedUDFs int
			}
		}

		sources := make([]sharedSourceCase, 0, len(udfQueries))

		for i, query := range udfQueries {
			item := sharedSourceCase{
				source: file.NewSource(fmt.Sprintf("udf_shared_%d", i), query.query),
			}

			item.spec.expectedHost = query.expectedHost
			item.spec.expectedUDFs = query.expectedUDFs
			sources = append(sources, item)
		}

		runConcurrentCompileWorkers(t, workers, iterations, func(worker, iter int) error {
			entry := sources[(worker+iter)%len(sources)]

			program, err := compiler.Compile(entry.source)
			if err != nil {
				return fmt.Errorf("compile failed: %w", err)
			}

			if err := assertCompiledProgram(program, entry.source); err != nil {
				return err
			}

			return assertUDFMetadata(program, entry.spec.expectedHost, entry.spec.expectedUDFs)
		})
	})

	t.Run("udf_isolation_fresh_sources", func(t *testing.T) {
		runConcurrentCompileWorkers(t, workers, iterations, func(worker, iter int) error {
			query := udfQueries[(worker+iter)%len(udfQueries)]
			source := file.NewSource(fmt.Sprintf("udf_fresh_%d_%d", worker, iter), query.query)

			program, err := compiler.Compile(source)
			if err != nil {
				return fmt.Errorf("compile failed: %w", err)
			}

			if err := assertCompiledProgram(program, source); err != nil {
				return err
			}

			return assertUDFMetadata(program, query.expectedHost, query.expectedUDFs)
		})
	})
}

func TestCompilerCompileConcurrentInvalidQueries(t *testing.T) {
	t.Parallel()

	compiler := New(WithOptimizationLevel(O1))
	workers := maxInt(8, runtime.GOMAXPROCS(0)*2)
	iterations := 80

	invalidQueries := []string{
		"",
		"RETURN",
		"LET x =\nRETURN x",
		"FOR i IN RETURN i",
	}

	t.Run("invalid_sources", func(t *testing.T) {
		runConcurrentCompileWorkers(t, workers, iterations, func(worker, iter int) error {
			query := invalidQueries[(worker+iter)%len(invalidQueries)]
			source := file.NewSource(fmt.Sprintf("invalid_%d_%d", worker, iter), query)

			program, err := compiler.Compile(source)
			if err == nil {
				return fmt.Errorf("expected compilation error")
			}

			if program != nil {
				return fmt.Errorf("expected nil program for invalid query")
			}

			message := strings.ToLower(err.Error())
			if strings.Contains(message, "unhandled panic") || strings.Contains(message, "goroutine ") {
				return fmt.Errorf("invalid query returned panic-like error: %s", err.Error())
			}

			return nil
		})
	})
}

func runConcurrentCompileWorkers(t *testing.T, workers, iterations int, fn func(worker, iter int) error) {
	t.Helper()

	errs := make(chan error, workers*iterations)

	var wg sync.WaitGroup

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			defer func() {
				if recovered := recover(); recovered != nil {
					errs <- fmt.Errorf("worker %d panicked: %v", workerID, recovered)
				}
			}()

			for iter := 0; iter < iterations; iter++ {
				if err := fn(workerID, iter); err != nil {
					errs <- fmt.Errorf("worker %d iteration %d: %w", workerID, iter, err)
				}
			}
		}(worker)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func assertCompiledProgram(program *bytecode.Program, source *file.Source) error {
	if program == nil {
		return fmt.Errorf("program is nil")
	}

	if program.Source != source {
		return fmt.Errorf("unexpected source pointer in output program")
	}

	if len(program.Bytecode) == 0 {
		return fmt.Errorf("compiled program has no bytecode")
	}

	return nil
}

func assertUDFMetadata(program *bytecode.Program, expectedHost map[string]int, expectedUDFs int) error {
	host := program.Functions.Host
	if !maps.Equal(host, expectedHost) {
		return fmt.Errorf("unexpected host metadata: got %v, expected %v", host, expectedHost)
	}

	if got := len(program.Functions.UserDefined); got != expectedUDFs {
		return fmt.Errorf("unexpected udf count: got %d, expected %d", got, expectedUDFs)
	}

	return nil
}
