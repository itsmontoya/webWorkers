```
go test --bench=BenchmarkPara. --count 4 --benchtime 4s
PASS
BenchmarkParaWWBasic-4            100000             59968 ns/op           14694 B/op        108 allocs/op
BenchmarkParaWWBasic-4            100000             61405 ns/op           14696 B/op        109 allocs/op
BenchmarkParaWWBasic-4            100000             59727 ns/op           14694 B/op        108 allocs/op
BenchmarkParaWWBasic-4            100000             59902 ns/op           14695 B/op        108 allocs/op
BenchmarkParaStdlibBasic-4        200000            128350 ns/op            6032 B/op         61 allocs/op
BenchmarkParaStdlibBasic-4         10000            662088 ns/op            5815 B/op         60 allocs/op
BenchmarkParaStdlibBasic-4         20000            738417 ns/op            6152 B/op         62 allocs/op
BenchmarkParaStdlibBasic-4         10000            436768 ns/op            5326 B/op         58 allocs/op
```
