module benchmark

go 1.24.0

toolchain go1.24.1

require (
	github.com/kelindar/goap v0.0.0-20231112144204-e9595370b8d7
	goapai v1.0.0
)

require (
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	golang.org/x/sys v0.14.0 // indirect
)

replace goapai => ./../.
