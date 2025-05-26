package runtime

import (
	"context"
)

type (
	SortLess = func(ctx context.Context, a, b Int) (Boolean, error)
	SortSwap = func(ctx context.Context, a, b Int) error
)

// sortInsertionSort sorts data[a:b] using insertion Sort.
func sortInsertionSort(ctx context.Context, less SortLess, swap SortSwap, a, b Int) error {
	for i := a + 1; i < b; i++ {
		for j := i; j > a; j-- {
			isLess, err := less(ctx, j, j-1)

			if err != nil {
				return err
			}

			if !isLess {
				break
			}

			if err := swap(ctx, j, j-1); err != nil {
				return err
			}
		}
	}

	return nil
}

func sortSwapRange(ctx context.Context, swap SortSwap, a, b, n Int) error {
	for i := ZeroInt; i < n; i++ {
		if err := swap(ctx, a+i, b+i); err != nil {
			return err
		}
	}

	return nil
}

func stableSort(ctx context.Context, less SortLess, swap SortSwap, n Int) error {
	blockSize := Int(20) // must be > 0
	a, b := ZeroInt, blockSize
	for b <= n {
		if err := sortInsertionSort(ctx, less, swap, a, b); err != nil {
			return err
		}

		a = b
		b += blockSize
	}

	if err := sortInsertionSort(ctx, less, swap, a, n); err != nil {
		return err
	}

	for blockSize < n {
		a, b = 0, 2*blockSize

		for b <= n {
			if err := sortSymMerge(ctx, less, swap, a, a+blockSize, b); err != nil {
				return err
			}

			a = b
			b += 2 * blockSize
		}

		if m := a + blockSize; m < n {
			if err := sortSymMerge(ctx, less, swap, a, m, n); err != nil {
				return err
			}
		}

		blockSize *= 2
	}

	return nil
}

func sortSymMerge(ctx context.Context, less SortLess, swap SortSwap, a, m, b Int) error {
	// Avoid unnecessary recursions of symMerge
	// by direct insertion of data[a] into data[m:b]
	// if data[a:m] only contains one element.
	if m-a == 1 {
		// Use binary search to find the lowest index i
		// such that data[i] >= data[a] for m <= i < b.
		// Exit the search loop with i == b in case no such index exists.
		i := m
		j := b
		for i < j {
			h := Int(uint(i+j) >> 1)

			isLess, err := less(ctx, h, a)

			if err != nil {
				return err
			}

			if isLess {
				i = h + 1
			} else {
				j = h
			}
		}
		// Swap values until data[a] reaches the position before i.
		for k := a; k < i-1; k++ {
			if err := swap(ctx, k, k+1); err != nil {
				return err
			}
		}

		return nil
	}

	// Avoid unnecessary recursions of symMerge
	// by direct insertion of data[m] into data[a:m]
	// if data[m:b] only contains one element.
	if b-m == 1 {
		// Use binary search to find the lowest index i
		// such that data[i] > data[m] for a <= i < m.
		// Exit the search loop with i == m in case no such index exists.
		i := a
		j := m
		for i < j {
			h := Int(uint(i+j) >> 1)

			isLess, err := less(ctx, m, h)

			if err != nil {
				return err
			}

			if !isLess {
				i = h + 1
			} else {
				j = h
			}
		}
		// Swap values until data[m] reaches the position i.
		for k := m; k > i; k-- {
			if err := swap(ctx, k, k-1); err != nil {
				return err
			}
		}

		return nil
	}

	mid := Int(uint(a+b) >> 1)
	n := mid + m
	var start, r Int

	if m > mid {
		start = n - b
		r = mid
	} else {
		start = a
		r = m
	}
	p := n - 1

	for start < r {
		c := Int(uint(start+r) >> 1)

		isLess, err := less(ctx, p-c, c)

		if err != nil {
			return err
		}

		if !isLess {
			start = c + 1
		} else {
			r = c
		}
	}

	end := n - start
	if start < m && m < end {
		if err := sortRotate(ctx, swap, start, m, end); err != nil {
			return err
		}
	}
	if a < start && start < mid {
		if err := sortSymMerge(ctx, less, swap, a, start, mid); err != nil {
			return err
		}
	}
	if mid < end && end < b {
		if err := sortSymMerge(ctx, less, swap, mid, end, b); err != nil {
			return err
		}
	}

	return nil
}

func sortRotate(ctx context.Context, swap SortSwap, a, m, b Int) error {
	i := m - a
	j := b - m

	for i != j {
		if i > j {
			if err := sortSwapRange(ctx, swap, m-i, m, j); err != nil {
				return err
			}

			i -= j
		} else {
			if err := sortSwapRange(ctx, swap, m-i, m+j-i, i); err != nil {
				return err
			}

			j -= i
		}
	}
	// i == j
	return sortSwapRange(ctx, swap, m-i, m, i)
}
