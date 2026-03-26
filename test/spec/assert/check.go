package assert

type (
	Unary func(actual any) error

	Binary func(actual, expected any) error
)

func ComposeUnary(checks ...Unary) Unary {
	return func(actual any) error {
		for _, check := range checks {
			if err := check(actual); err != nil {
				return err
			}
		}

		return nil
	}
}

func ComposeBinary(checks ...Binary) Binary {
	return func(actual, expected any) error {
		for _, check := range checks {
			if err := check(actual, expected); err != nil {
				return err
			}
		}

		return nil
	}
}
