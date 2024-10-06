package main

func caller2() {
	panic("panic")
}

func caller1() {
	caller2()
}

func main() {
}
