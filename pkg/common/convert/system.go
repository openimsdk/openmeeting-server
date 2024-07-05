package convert

// ByInt32 Define a type ByInt32 to implement the sort.Interface
type ByInt32 []int32

// Len returns the length of the array
func (a ByInt32) Len() int {
	return len(a)
}

// Less compares two elements in the array and determines their order
func (a ByInt32) Less(i, j int) bool {
	return a[i] < a[j]
}

// Swap swaps two elements in the array
func (a ByInt32) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
