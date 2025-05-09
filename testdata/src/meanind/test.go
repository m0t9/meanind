package meanind

func meaningfulIndexIteration() {
	// Key & value used, no worries
	for meaningfulName, item := range []int{} {
		_ = meaningfulName
		_ = item
	}

	// Ranging over non-array / slice
	for meaningfulName := range 5 {
		_ = meaningfulName
	}
	for meaningfulName := range map[string]string{} {
		_ = meaningfulName
	}

	for x := range []int{} { // want `for-range loop with confusing iteratee name detected. Is "x" name suitable for array/slice index?`
		_ = x
	}

	// Well known names.
	for idx := range []int{} {
		_ = idx
	}
	for jdx := range []int{} {
		_ = jdx
	}
	for jdx := range []int{} {
		_ = jdx
	}

	for ind := range []int{} {
		_ = ind
	}

	// Well-known names modifications
	for idxOfArr := range []int{} {
		_ = idxOfArr
	}
	for jdxOfArr := range []int{} {
		_ = jdxOfArr
	}
	for kdxOfArr := range []int{} {
		_ = kdxOfArr
	}

	for arrIdx := range []int{} {
		_ = arrIdx
	}
	for arrJdx := range []int{} {
		_ = arrJdx
	}
	for arrKdx := range []int{} {
		_ = arrKdx
	}

	for indOfArr := range []int{} {
		_ = indOfArr
	}

	for arrInd := range []int{} {
		_ = arrInd
	}

	// Works on arrays
	for indOfArr := range [5]int{0, 1, 2, 3, 4} {
		_ = indOfArr
	}

	for app := range [5]int{0, 1, 2, 3, 4} { // want `for-range loop with confusing iteratee name detected. Is "app" name suitable for array/slice index?`
		_ = app
	}

	for arrInd := range []int{} {
		_ = arrInd
	}

	// Does not report cases when iteratee is used as index.
	arr := []int{0, 1, 2}
	for x := range arr {
		_ = arr[x]
	}

	// Does not report if it participates in index expression.
	for x := range arr {
		_ = arr[x+1-1]
	}

	// Report cases when used, but with shadowed variables.
	for x := range arr { // want `for-range loop with confusing iteratee name detected. Is "x" name suitable for array/slice index?`
		if true {
			// Used in index but in shadowed instance.
			arr := []int{1, 0}
			_ = arr[x]
		}
	}

	// Report cases where access to some other array / slice.
	for x := range arr { // want `for-range loop with confusing iteratee name detected. Is "x" name suitable for array/slice index?`
		t := func() []int {
			return []int{1, 2, 3}
		}

		_ = t()[x]
	}

	for x := range arr { // want `for-range loop with confusing iteratee name detected. Is "x" name suitable for array/slice index?`
		idx := 0
		_ = arr[idx]
		_ = x
	}
}
