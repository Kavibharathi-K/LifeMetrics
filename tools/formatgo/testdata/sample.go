package testdata

func Inner(a int, b int) int {
	return a + b
}

func Outer(a int, b int) int {
	return Inner(a, b) + Inner(b, a)
}

func Example() int {
	return Outer(Inner(1, 2), Inner(3, 4))
}
