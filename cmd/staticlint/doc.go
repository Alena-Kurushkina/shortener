// Util staticlint is a multichecker consisting of:
// • printf, shadow and structtag analyzers of golang.org/x/tools/go/analysis/passes package
// • custom ExitCheckAnalyzer
// • all analyzers of SA type from staticcheck.io package
// • S1003, ST1005, QF1010 analyzers of staticcheck.io package
// • exportloopref.Analyzer and thelper.Analyzer.
//
// Usage:
//   cmd/staticlint/main.go [path to folder for check]
//
// Example:
//   ./cmd/staticlint/main.go ./...
//
// Description of analyzers:
//
// OSEXITCHECK
//
//
// check for os.Exit call in main function of main package
//
// It returns warning if it finds os.Exit() call in main function of main package.
//
// For example:
//   package main
//
//   import (
// 	"os"
//   )
//
//   func main(){
//     os.Exit(3)
//   }
//
//
//
// PRINTF
//
// check consistency of Printf format strings and arguments
//
// The check applies to calls of the formatting functions such as
// [fmt.Printf] and [fmt.Sprintf], as well as any detected wrappers of
// those functions such as [log.Printf]. It reports a variety of
// mistakes such as syntax errors in the format string and mismatches
// (of number and type) between the verbs and their arguments.
//
// See the documentation of the fmt package for the complete set of
// format operators and their operand types.
//
// SHADOW
//
// check for possible unintended shadowing of variables
//
// This analyzer check for shadowed variables.
// A shadowed variable is a variable declared in an inner scope
// with the same name and type as a variable in an outer scope,
// and where the outer variable is mentioned after the inner one
// is declared.
//
// (This definition can be refined; the module generates too many
// false positives and is not yet enabled by default.)
//
// For example:
//
// 	func BadRead(f *os.File, buf []byte) error {
// 		var err error
// 		for {
// 			n, err := f.Read(buf) // shadows the function variable 'err'
// 			if err != nil {
// 				break // causes return of wrong value
// 			}
// 			foo(buf)
// 		}
// 		return err
// 	}
//
// STRUCTTAG
//
// check that struct field tags conform to reflect.StructTag.Get
//
// Also report certain struct tags (json, xml) used with unexported fields.
//
// SA1000
//
// invalid regular expression
//
// Available since
//     2017.1
//
//
// SA1001
//
// invalid template
//
// Available since
//     2017.1
//
//
// SA1002
//
// invalid format in time.Parse
//
// Available since
//     2017.1
//
//
// SA1003
//
// unsupported argument to functions in encoding/binary
//
// The encoding/binary package can only serialize types with known sizes.
// This precludes the use of the int and uint types, as their sizes
// differ on different architectures. Furthermore, it doesn't support
// serializing maps, channels, strings, or functions.
//
// Before Go 1.8, bool wasn't supported, either.
//
// Available since
//     2017.1
//
//
// SA1004
//
// suspiciously small untyped constant in time.Sleep
//
// The time.Sleep function takes a time.Duration as its only argument.
// Durations are expressed in nanoseconds. Thus, calling time.Sleep(1)
// will sleep for 1 nanosecond. This is a common source of bugs, as sleep
// functions in other languages often accept seconds or milliseconds.
//
// The time package provides constants such as time.Second to express
// large durations. These can be combined with arithmetic to express
// arbitrary durations, for example 5 * time.Second for 5 seconds.
//
// If you truly meant to sleep for a tiny amount of time, use
// n * time.Nanosecond to signal to Staticcheck that you did mean to sleep
// for some amount of nanoseconds.
//
// Available since
//     2017.1
//
//
// SA1005
//
// invalid first argument to exec.Command
//
// os/exec runs programs directly (using variants of the fork and exec
// system calls on Unix systems). This shouldn't be confused with running
// a command in a shell. The shell will allow for features such as input
// redirection, pipes, and general scripting. The shell is also
// responsible for splitting the user's input into a program name and its
// arguments. For example, the equivalent to
//
//     ls / /tmp
//
// would be
//
//     exec.Command("ls", "/", "/tmp")
//
// If you want to run a command in a shell, consider using something like
// the following – but be aware that not all systems, particularly
// Windows, will have a /bin/sh program:
//
//     exec.Command("/bin/sh", "-c", "ls | grep Awesome")
//
// Available since
//     2017.1
//
//
// SA1006
//
// printf with dynamic first argument and no further arguments
//
// Using fmt.Printf with a dynamic first argument can lead to unexpected
// output. The first argument is a format string, where certain character
// combinations have special meaning. If, for example, a user were to
// enter a string such as
//
//     Interest rate: 5%
//
// and you printed it with
//
//     fmt.Printf(s)
//
// it would lead to the following output:
//
//     Interest rate: 5%!(NOVERB).
//
// Similarly, forming the first parameter via string concatenation with
// user input should be avoided for the same reason. When printing user
// input, either use a variant of fmt.Print, or use the %s Printf verb
// and pass the string as an argument.
//
// Available since
//     2017.1
//
//
// SA1007
//
// invalid URL in net/url.Parse
//
// Available since
//     2017.1
//
//
// SA1008
//
// non-canonical key in http.Header map
//
// Keys in http.Header maps are canonical, meaning they follow a specific
// combination of uppercase and lowercase letters. Methods such as
// http.Header.Add and http.Header.Del convert inputs into this canonical
// form before manipulating the map.
//
// When manipulating http.Header maps directly, as opposed to using the
// provided methods, care should be taken to stick to canonical form in
// order to avoid inconsistencies. The following piece of code
// demonstrates one such inconsistency:
//
//     h := http.Header{}
//     h["etag"] = []string{"1234"}
//     h.Add("etag", "5678")
//     fmt.Println(h)
//
//     // Output:
//     // map[Etag:[5678] etag:[1234]]
//
// The easiest way of obtaining the canonical form of a key is to use
// http.CanonicalHeaderKey.
//
// Available since
//     2017.1
//
//
// SA1010
//
// (*regexp.Regexp).FindAll called with n == 0, which will always return zero results
//
// If n >= 0, the function returns at most n matches/submatches. To
// return all results, specify a negative number.
//
// Available since
//     2017.1
//
//
// SA1011
//
// various methods in the 'strings' package expect valid UTF-8, but invalid input is provided
//
// Available since
//     2017.1
//
//
// SA1012
//
// a nil context.Context is being passed to a function, consider using context.TODO instead
//
// Available since
//     2017.1
//
//
// SA1013
//
// io.Seeker.Seek is being called with the whence constant as the first argument, but it should be the second
//
// Available since
//     2017.1
//
//
// SA1014
//
// non-pointer value passed to Unmarshal or Decode
//
// Available since
//     2017.1
//
//
// SA1015
//
// using time.Tick in a way that will leak. Consider using time.NewTicker, and only use time.Tick in tests, commands and endless functions
//
// Before Go 1.23, time.Tickers had to be closed to be able to be garbage
// collected. Since time.Tick doesn't make it possible to close the underlying
// ticker, using it repeatedly would leak memory.
//
// Go 1.23 fixes this by allowing tickers to be collected even if they weren't closed.
//
// Available since
//     2017.1
//
//
// SA1016
//
// trapping a signal that cannot be trapped
//
// Not all signals can be intercepted by a process. Specifically, on
// UNIX-like systems, the syscall.SIGKILL and syscall.SIGSTOP signals are
// never passed to the process, but instead handled directly by the
// kernel. It is therefore pointless to try and handle these signals.
//
// Available since
//     2017.1
//
//
// SA1017
//
// channels used with os/signal.Notify should be buffered
//
// The os/signal package uses non-blocking channel sends when delivering
// signals. If the receiving end of the channel isn't ready and the
// channel is either unbuffered or full, the signal will be dropped. To
// avoid missing signals, the channel should be buffered and of the
// appropriate size. For a channel used for notification of just one
// signal value, a buffer of size 1 is sufficient.
//
// Available since
//     2017.1
//
//
// SA1018
//
// strings.Replace called with n == 0, which does nothing
//
// With n == 0, zero instances will be replaced. To replace all
// instances, use a negative number, or use strings.ReplaceAll.
//
// Available since
//     2017.1
//
//
// SA1019
//
// using a deprecated function, variable, constant or field
//
// Available since
//     2017.1
//
//
// SA1020
//
// using an invalid host:port pair with a net.Listen-related function
//
// Available since
//     2017.1
//
//
// SA1021
//
// using bytes.Equal to compare two net.IP
//
// A net.IP stores an IPv4 or IPv6 address as a slice of bytes. The
// length of the slice for an IPv4 address, however, can be either 4 or
// 16 bytes long, using different ways of representing IPv4 addresses. In
// order to correctly compare two net.IPs, the net.IP.Equal method should
// be used, as it takes both representations into account.
//
// Available since
//     2017.1
//
//
// SA1023
//
// modifying the buffer in an io.Writer implementation
//
// Write must not modify the slice data, even temporarily.
//
// Available since
//     2017.1
//
//
// SA1024
//
// a string cutset contains duplicate characters
//
// The strings.TrimLeft and strings.TrimRight functions take cutsets, not
// prefixes. A cutset is treated as a set of characters to remove from a
// string. For example,
//
//     strings.TrimLeft("42133word", "1234")
//
// will result in the string "word" – any characters that are 1, 2, 3 or
// 4 are cut from the left of the string.
//
// In order to remove one string from another, use strings.TrimPrefix instead.
//
// Available since
//     2017.1
//
//
// SA1025
//
// it is not possible to use (*time.Timer).Reset's return value correctly
//
// Available since
//     2019.1
//
//
// SA1026
//
// cannot marshal channels or functions
//
// Available since
//     2019.2
//
//
// SA1027
//
// atomic access to 64-bit variable must be 64-bit aligned
//
// On ARM, x86-32, and 32-bit MIPS, it is the caller's responsibility to
// arrange for 64-bit alignment of 64-bit words accessed atomically. The
// first word in a variable or in an allocated struct, array, or slice
// can be relied upon to be 64-bit aligned.
//
// You can use the structlayout tool to inspect the alignment of fields
// in a struct.
//
// Available since
//     2019.2
//
//
// SA1028
//
// sort.Slice can only be used on slices
//
// The first argument of sort.Slice must be a slice.
//
// Available since
//     2020.1
//
//
// SA1029
//
// inappropriate key in call to context.WithValue
//
// The provided key must be comparable and should not be
// of type string or any other built-in type to avoid collisions between
// packages using context. Users of WithValue should define their own
// types for keys.
//
// To avoid allocating when assigning to an interface{},
// context keys often have concrete type struct{}. Alternatively,
// exported context key variables' static type should be a pointer or
// interface.
//
// Available since
//     2020.1
//
//
// SA1030
//
// invalid argument in call to a strconv function
//
// This check validates the format, number base and bit size arguments of
// the various parsing and formatting functions in strconv.
//
// Available since
//     2021.1
//
//
// SA1031
//
// overlapping byte slices passed to an encoder
//
// In an encoding function of the form Encode(dst, src), dst and
// src were found to reference the same memory. This can result in
// src bytes being overwritten before they are read, when the encoder
// writes more than one byte per src byte.
//
// Available since
//     2024.1
//
//
// SA1032
//
// wrong order of arguments to errors.Is
//
// The first argument of the function errors.Is is the error
// that we have and the second argument is the error we're trying to match against.
// For example:
//
// 	if errors.Is(err, io.EOF) { ... }
//
// This check detects some cases where the two arguments have been swapped. It
// flags any calls where the first argument is referring to a package-level error
// variable, such as
//
// 	if errors.Is(io.EOF, err) { /* this is wrong */ }
//
// Available since
//     2024.1
//
//
// SA2000
//
// sync.WaitGroup.Add called inside the goroutine, leading to a race condition
//
// Available since
//     2017.1
//
//
// SA2001
//
// empty critical section, did you mean to defer the unlock?
//
// Empty critical sections of the kind
//
//     mu.Lock()
//     mu.Unlock()
//
// are very often a typo, and the following was intended instead:
//
//     mu.Lock()
//     defer mu.Unlock()
//
// Do note that sometimes empty critical sections can be useful, as a
// form of signaling to wait on another goroutine. Many times, there are
// simpler ways of achieving the same effect. When that isn't the case,
// the code should be amply commented to avoid confusion. Combining such
// comments with a //lint:ignore directive can be used to suppress this
// rare false positive.
//
// Available since
//     2017.1
//
//
// SA2002
//
// called testing.T.FailNow or SkipNow in a goroutine, which isn't allowed
//
// Available since
//     2017.1
//
//
// SA2003
//
// deferred Lock right after locking, likely meant to defer Unlock instead
//
// Available since
//     2017.1
//
//
// SA3000
//
// testMain doesn't call os.Exit, hiding test failures
//
// Test executables (and in turn 'go test') exit with a non-zero status
// code if any tests failed. When specifying your own TestMain function,
// it is your responsibility to arrange for this, by calling os.Exit with
// the correct code. The correct code is returned by (*testing.M).Run, so
// the usual way of implementing TestMain is to end it with
// os.Exit(m.Run()).
//
// Available since
//     2017.1
//
//
// SA3001
//
// assigning to b.N in benchmarks distorts the results
//
// The testing package dynamically sets b.N to improve the reliability of
// benchmarks and uses it in computations to determine the duration of a
// single operation. Benchmark code must not alter b.N as this would
// falsify results.
//
// Available since
//     2017.1
//
//
// SA4000
//
// binary operator has identical expressions on both sides
//
// Available since
//     2017.1
//
//
// SA4001
//
// &*x gets simplified to x, it does not copy x
//
// Available since
//     2017.1
//
//
// SA4003
//
// comparing unsigned values against negative values is pointless
//
// Available since
//     2017.1
//
//
// SA4004
//
// the loop exits unconditionally after one iteration
//
// Available since
//     2017.1
//
//
// SA4005
//
// field assignment that will never be observed. Did you mean to use a pointer receiver?
//
// Available since
//     2021.1
//
//
// SA4006
//
// a value assigned to a variable is never read before being overwritten. Forgotten error check or dead code?
//
// Available since
//     2017.1
//
//
// SA4008
//
// the variable in the loop condition never changes, are you incrementing the wrong variable?
//
// For example:
//
// 	for i := 0; i < 10; j++ { ... }
//
// This may also occur when a loop can only execute once because of unconditional
// control flow that terminates the loop. For example, when a loop body contains an
// unconditional break, return, or panic:
//
// 	func f() {
// 		panic("oops")
// 	}
// 	func g() {
// 		for i := 0; i < 10; i++ {
// 			// f unconditionally calls panic, which means "i" is
// 			// never incremented.
// 			f()
// 		}
// 	}
//
// Available since
//     2017.1
//
//
// SA4009
//
// a function argument is overwritten before its first use
//
// Available since
//     2017.1
//
//
// SA4010
//
// the result of append will never be observed anywhere
//
// Available since
//     2017.1
//
//
// SA4011
//
// break statement with no effect. Did you mean to break out of an outer loop?
//
// Available since
//     2017.1
//
//
// SA4012
//
// comparing a value against NaN even though no value is equal to NaN
//
// Available since
//     2017.1
//
//
// SA4013
//
// negating a boolean twice (!!b) is the same as writing b. This is either redundant, or a typo.
//
// Available since
//     2017.1
//
//
// SA4014
//
// an if/else if chain has repeated conditions and no side-effects; if the condition didn't match the first time, it won't match the second time, either
//
// Available since
//     2017.1
//
//
// SA4015
//
// calling functions like math.Ceil on floats converted from integers doesn't do anything useful
//
// Available since
//     2017.1
//
//
// SA4016
//
// certain bitwise operations, such as x ^ 0, do not do anything useful
//
// Available since
//     2017.1
//
//
// SA4017
//
// discarding the return values of a function without side effects, making the call pointless
//
// Available since
//     2017.1
//
//
// SA4018
//
// self-assignment of variables
//
// Available since
//     2017.1
//
//
// SA4019
//
// multiple, identical build constraints in the same file
//
// Available since
//     2017.1
//
//
// SA4020
//
// unreachable case clause in a type switch
//
// In a type switch like the following
//
//     type T struct{}
//     func (T) Read(b []byte) (int, error) { return 0, nil }
//
//     var v interface{} = T{}
//
//     switch v.(type) {
//     case io.Reader:
//         // ...
//     case T:
//         // unreachable
//     }
//
// the second case clause can never be reached because T implements
// io.Reader and case clauses are evaluated in source order.
//
// Another example:
//
//     type T struct{}
//     func (T) Read(b []byte) (int, error) { return 0, nil }
//     func (T) Close() error { return nil }
//
//     var v interface{} = T{}
//
//     switch v.(type) {
//     case io.Reader:
//         // ...
//     case io.ReadCloser:
//         // unreachable
//     }
//
// Even though T has a Close method and thus implements io.ReadCloser,
// io.Reader will always match first. The method set of io.Reader is a
// subset of io.ReadCloser. Thus it is impossible to match the second
// case without matching the first case.
//
//
// Structurally equivalent interfaces
//
// A special case of the previous example are structurally identical
// interfaces. Given these declarations
//
//     type T error
//     type V error
//
//     func doSomething() error {
//         err, ok := doAnotherThing()
//         if ok {
//             return T(err)
//         }
//
//         return U(err)
//     }
//
// the following type switch will have an unreachable case clause:
//
//     switch doSomething().(type) {
//     case T:
//         // ...
//     case V:
//         // unreachable
//     }
//
// T will always match before V because they are structurally equivalent
// and therefore doSomething()'s return value implements both.
//
// Available since
//     2019.2
//
//
// SA4021
//
// 'x = append(y)' is equivalent to 'x = y'
//
// Available since
//     2019.2
//
//
// SA4022
//
// comparing the address of a variable against nil
//
// Code such as 'if &x == nil' is meaningless, because taking the address of a variable always yields a non-nil pointer.
//
// Available since
//     2020.1
//
//
// SA4023
//
// impossible comparison of interface value with untyped nil
//
// Under the covers, interfaces are implemented as two elements, a
// type T and a value V. V is a concrete value such as an int,
// struct or pointer, never an interface itself, and has type T. For
// instance, if we store the int value 3 in an interface, the
// resulting interface value has, schematically, (T=int, V=3). The
// value V is also known as the interface's dynamic value, since a
// given interface variable might hold different values V (and
// corresponding types T) during the execution of the program.
//
// An interface value is nil only if the V and T are both
// unset, (T=nil, V is not set), In particular, a nil interface will
// always hold a nil type. If we store a nil pointer of type *int
// inside an interface value, the inner type will be *int regardless
// of the value of the pointer: (T=*int, V=nil). Such an interface
// value will therefore be non-nil even when the pointer value V
// inside is nil.
//
// This situation can be confusing, and arises when a nil value is
// stored inside an interface value such as an error return:
//
//     func returnsError() error {
//         var p *MyError = nil
//         if bad() {
//             p = ErrBad
//         }
//         return p // Will always return a non-nil error.
//     }
//
// If all goes well, the function returns a nil p, so the return
// value is an error interface value holding (T=*MyError, V=nil).
// This means that if the caller compares the returned error to nil,
// it will always look as if there was an error even if nothing bad
// happened. To return a proper nil error to the caller, the
// function must return an explicit nil:
//
//     func returnsError() error {
//         if bad() {
//             return ErrBad
//         }
//         return nil
//     }
//
// It's a good idea for functions that return errors always to use
// the error type in their signature (as we did above) rather than a
// concrete type such as *MyError, to help guarantee the error is
// created correctly. As an example, os.Open returns an error even
// though, if not nil, it's always of concrete type *os.PathError.
//
// Similar situations to those described here can arise whenever
// interfaces are used. Just keep in mind that if any concrete value
// has been stored in the interface, the interface will not be nil.
// For more information, see The Laws of
// Reflection (https://golang.org/doc/articles/laws_of_reflection.html).
//
// This text has been copied from
// https://golang.org/doc/faq#nil_error, licensed under the Creative
// Commons Attribution 3.0 License.
//
// Available since
//     2020.2
//
//
// SA4024
//
// checking for impossible return value from a builtin function
//
// Return values of the len and cap builtins cannot be negative.
//
// See https://golang.org/pkg/builtin/#len and https://golang.org/pkg/builtin/#cap.
//
// Example:
//
//     if len(slice) < 0 {
//         fmt.Println("unreachable code")
//     }
//
// Available since
//     2021.1
//
//
// SA4025
//
// integer division of literals that results in zero
//
// When dividing two integer constants, the result will
// also be an integer. Thus, a division such as 2 / 3 results in 0.
// This is true for all of the following examples:
//
// 	_ = 2 / 3
// 	const _ = 2 / 3
// 	const _ float64 = 2 / 3
// 	_ = float64(2 / 3)
//
// Staticcheck will flag such divisions if both sides of the division are
// integer literals, as it is highly unlikely that the division was
// intended to truncate to zero. Staticcheck will not flag integer
// division involving named constants, to avoid noisy positives.
//
// Available since
//     2021.1
//
//
// SA4026
//
// go constants cannot express negative zero
//
// In IEEE 754 floating point math, zero has a sign and can be positive
// or negative. This can be useful in certain numerical code.
//
// Go constants, however, cannot express negative zero. This means that
// the literals -0.0 and 0.0 have the same ideal value (zero) and
// will both represent positive zero at runtime.
//
// To explicitly and reliably create a negative zero, you can use the
// math.Copysign function: math.Copysign(0, -1).
//
// Available since
//     2021.1
//
//
// SA4027
//
// (*net/url.URL).Query returns a copy, modifying it doesn't change the URL
//
// (*net/url.URL).Query parses the current value of net/url.URL.RawQuery
// and returns it as a map of type net/url.Values. Subsequent changes to
// this map will not affect the URL unless the map gets encoded and
// assigned to the URL's RawQuery.
//
// As a consequence, the following code pattern is an expensive no-op:
// u.Query().Add(key, value).
//
// Available since
//     2021.1
//
//
// SA4028
//
// x % 1 is always zero
//
// Available since
//     2022.1
//
//
// SA4029
//
// ineffective attempt at sorting slice
//
// sort.Float64Slice, sort.IntSlice, and sort.StringSlice are
// types, not functions. Doing x = sort.StringSlice(x) does nothing,
// especially not sort any values. The correct usage is
// sort.Sort(sort.StringSlice(x)) or sort.StringSlice(x).Sort(),
// but there are more convenient helpers, namely sort.Float64s,
// sort.Ints, and sort.Strings.
//
// Available since
//     2022.1
//
//
// SA4030
//
// ineffective attempt at generating random number
//
// Functions in the math/rand package that accept upper limits, such
// as Intn, generate random numbers in the half-open interval [0,n). In
// other words, the generated numbers will be >= 0 and < n – they
// don't include n. rand.Intn(1) therefore doesn't generate 0
// or 1, it always generates 0.
//
// Available since
//     2022.1
//
//
// SA4031
//
// checking never-nil value against nil
//
// Available since
//     2022.1
//
//
// SA4032
//
// comparing runtime.GOOS or runtime.GOARCH against impossible value
//
// Available since
//     2024.1
//
//
// SA5000
//
// assignment to nil map
//
// Available since
//     2017.1
//
//
// SA5001
//
// deferring Close before checking for a possible error
//
// Available since
//     2017.1
//
//
// SA5002
//
// the empty for loop ('for {}') spins and can block the scheduler
//
// Available since
//     2017.1
//
//
// SA5003
//
// defers in infinite loops will never execute
//
// Defers are scoped to the surrounding function, not the surrounding
// block. In a function that never returns, i.e. one containing an
// infinite loop, defers will never execute.
//
// Available since
//     2017.1
//
//
// SA5004
//
// 'for { select { ...' with an empty default branch spins
//
// Available since
//     2017.1
//
//
// SA5005
//
// the finalizer references the finalized object, preventing garbage collection
//
// A finalizer is a function associated with an object that runs when the
// garbage collector is ready to collect said object, that is when the
// object is no longer referenced by anything.
//
// If the finalizer references the object, however, it will always remain
// as the final reference to that object, preventing the garbage
// collector from collecting the object. The finalizer will never run,
// and the object will never be collected, leading to a memory leak. That
// is why the finalizer should instead use its first argument to operate
// on the object. That way, the number of references can temporarily go
// to zero before the object is being passed to the finalizer.
//
// Available since
//     2017.1
//
//
// SA5007
//
// infinite recursive call
//
// A function that calls itself recursively needs to have an exit
// condition. Otherwise it will recurse forever, until the system runs
// out of memory.
//
// This issue can be caused by simple bugs such as forgetting to add an
// exit condition. It can also happen "on purpose". Some languages have
// tail call optimization which makes certain infinite recursive calls
// safe to use. Go, however, does not implement TCO, and as such a loop
// should be used instead.
//
// Available since
//     2017.1
//
//
// SA5008
//
// invalid struct tag
//
// Available since
//     2019.2
//
//
// SA5009
//
// invalid Printf call
//
// Available since
//     2019.2
//
//
// SA5010
//
// impossible type assertion
//
// Some type assertions can be statically proven to be
// impossible. This is the case when the method sets of both
// arguments of the type assertion conflict with each other, for
// example by containing the same method with different
// signatures.
//
// The Go compiler already applies this check when asserting from an
// interface value to a concrete type. If the concrete type misses
// methods from the interface, or if function signatures don't match,
// then the type assertion can never succeed.
//
// This check applies the same logic when asserting from one interface to
// another. If both interface types contain the same method but with
// different signatures, then the type assertion can never succeed,
// either.
//
// Available since
//     2020.1
//
//
// SA5011
//
// possible nil pointer dereference
//
// A pointer is being dereferenced unconditionally, while
// also being checked against nil in another place. This suggests that
// the pointer may be nil and dereferencing it may panic. This is
// commonly a result of improperly ordered code or missing return
// statements. Consider the following examples:
//
//     func fn(x *int) {
//         fmt.Println(*x)
//
//         // This nil check is equally important for the previous dereference
//         if x != nil {
//             foo(*x)
//         }
//     }
//
//     func TestFoo(t *testing.T) {
//         x := compute()
//         if x == nil {
//             t.Errorf("nil pointer received")
//         }
//
//         // t.Errorf does not abort the test, so if x is nil, the next line will panic.
//         foo(*x)
//     }
//
// Staticcheck tries to deduce which functions abort control flow.
// For example, it is aware that a function will not continue
// execution after a call to panic or log.Fatal. However, sometimes
// this detection fails, in particular in the presence of
// conditionals. Consider the following example:
//
//     func Log(msg string, level int) {
//         fmt.Println(msg)
//         if level == levelFatal {
//             os.Exit(1)
//         }
//     }
//
//     func Fatal(msg string) {
//         Log(msg, levelFatal)
//     }
//
//     func fn(x *int) {
//         if x == nil {
//             Fatal("unexpected nil pointer")
//         }
//         fmt.Println(*x)
//     }
//
// Staticcheck will flag the dereference of x, even though it is perfectly
// safe. Staticcheck is not able to deduce that a call to
// Fatal will exit the program. For the time being, the easiest
// workaround is to modify the definition of Fatal like so:
//
//     func Fatal(msg string) {
//         Log(msg, levelFatal)
//         panic("unreachable")
//     }
//
// We also hard-code functions from common logging packages such as
// logrus. Please file an issue if we're missing support for a
// popular package.
//
// Available since
//     2020.1
//
//
// SA5012
//
// passing odd-sized slice to function expecting even size
//
// Some functions that take slices as parameters expect the slices to have an even number of elements.
// Often, these functions treat elements in a slice as pairs.
// For example, strings.NewReplacer takes pairs of old and new strings,
// and calling it with an odd number of elements would be an error.
//
// Available since
//     2020.2
//
//
// SA6000
//
// using regexp.Match or related in a loop, should use regexp.Compile
//
// Available since
//     2017.1
//
//
// SA6001
//
// missing an optimization opportunity when indexing maps by byte slices
//
// Map keys must be comparable, which precludes the use of byte slices.
// This usually leads to using string keys and converting byte slices to
// strings.
//
// Normally, a conversion of a byte slice to a string needs to copy the data and
// causes allocations. The compiler, however, recognizes m[string(b)] and
// uses the data of b directly, without copying it, because it knows that
// the data can't change during the map lookup. This leads to the
// counter-intuitive situation that
//
//     k := string(b)
//     println(m[k])
//     println(m[k])
//
// will be less efficient than
//
//     println(m[string(b)])
//     println(m[string(b)])
//
// because the first version needs to copy and allocate, while the second
// one does not.
//
// For some history on this optimization, check out commit
// f5f5a8b6209f84961687d993b93ea0d397f5d5bf in the Go repository.
//
// Available since
//     2017.1
//
//
// SA6002
//
// storing non-pointer values in sync.Pool allocates memory
//
// A sync.Pool is used to avoid unnecessary allocations and reduce the
// amount of work the garbage collector has to do.
//
// When passing a value that is not a pointer to a function that accepts
// an interface, the value needs to be placed on the heap, which means an
// additional allocation. Slices are a common thing to put in sync.Pools,
// and they're structs with 3 fields (length, capacity, and a pointer to
// an array). In order to avoid the extra allocation, one should store a
// pointer to the slice instead.
//
// See the comments on https://go-review.googlesource.com/c/go/+/24371
// that discuss this problem.
//
// Available since
//     2017.1
//
//
// SA6003
//
// converting a string to a slice of runes before ranging over it
//
// You may want to loop over the runes in a string. Instead of converting
// the string to a slice of runes and looping over that, you can loop
// over the string itself. That is,
//
//     for _, r := range s {}
//
// and
//
//     for _, r := range []rune(s) {}
//
// will yield the same values. The first version, however, will be faster
// and avoid unnecessary memory allocations.
//
// Do note that if you are interested in the indices, ranging over a
// string and over a slice of runes will yield different indices. The
// first one yields byte offsets, while the second one yields indices in
// the slice of runes.
//
// Available since
//     2017.1
//
//
// SA6005
//
// inefficient string comparison with strings.ToLower or strings.ToUpper
//
// Converting two strings to the same case and comparing them like so
//
//     if strings.ToLower(s1) == strings.ToLower(s2) {
//         ...
//     }
//
// is significantly more expensive than comparing them with
// strings.EqualFold(s1, s2). This is due to memory usage as well as
// computational complexity.
//
// strings.ToLower will have to allocate memory for the new strings, as
// well as convert both strings fully, even if they differ on the very
// first byte. strings.EqualFold, on the other hand, compares the strings
// one character at a time. It doesn't need to create two intermediate
// strings and can return as soon as the first non-matching character has
// been found.
//
// For a more in-depth explanation of this issue, see
// https://blog.digitalocean.com/how-to-efficiently-compare-strings-in-go/
//
// Available since
//     2019.2
//
//
// SA6006
//
// using io.WriteString to write []byte
//
// Using io.WriteString to write a slice of bytes, as in
//
//     io.WriteString(w, string(b))
//
// is both unnecessary and inefficient. Converting from []byte to string
// has to allocate and copy the data, and we could simply use w.Write(b)
// instead.
//
// Available since
//     2024.1
//
//
// SA9001
//
// defers in range loops may not run when you expect them to
//
// Available since
//     2017.1
//
//
// SA9002
//
// using a non-octal os.FileMode that looks like it was meant to be in octal.
//
// Available since
//     2017.1
//
//
// SA9003
//
// empty body in an if or else branch
//
// Available since
//     2017.1, non-default
//
//
// SA9004
//
// only the first constant has an explicit type
//
// In a constant declaration such as the following:
//
//     const (
//         First byte = 1
//         Second     = 2
//     )
//
// the constant Second does not have the same type as the constant First.
// This construct shouldn't be confused with
//
//     const (
//         First byte = iota
//         Second
//     )
//
// where First and Second do indeed have the same type. The type is only
// passed on when no explicit value is assigned to the constant.
//
// When declaring enumerations with explicit values it is therefore
// important not to write
//
//     const (
//           EnumFirst EnumType = 1
//           EnumSecond         = 2
//           EnumThird          = 3
//     )
//
// This discrepancy in types can cause various confusing behaviors and
// bugs.
//
//
// Wrong type in variable declarations
//
// The most obvious issue with such incorrect enumerations expresses
// itself as a compile error:
//
//     package pkg
//
//     const (
//         EnumFirst  uint8 = 1
//         EnumSecond       = 2
//     )
//
//     func fn(useFirst bool) {
//         x := EnumSecond
//         if useFirst {
//             x = EnumFirst
//         }
//     }
//
// fails to compile with
//
//     ./const.go:11:5: cannot use EnumFirst (type uint8) as type int in assignment
//
//
// Losing method sets
//
// A more subtle issue occurs with types that have methods and optional
// interfaces. Consider the following:
//
//     package main
//
//     import "fmt"
//
//     type Enum int
//
//     func (e Enum) String() string {
//         return "an enum"
//     }
//
//     const (
//         EnumFirst  Enum = 1
//         EnumSecond      = 2
//     )
//
//     func main() {
//         fmt.Println(EnumFirst)
//         fmt.Println(EnumSecond)
//     }
//
// This code will output
//
//     an enum
//     2
//
// as EnumSecond has no explicit type, and thus defaults to int.
//
// Available since
//     2019.1
//
//
// SA9005
//
// trying to marshal a struct with no public fields nor custom marshaling
//
// The encoding/json and encoding/xml packages only operate on exported
// fields in structs, not unexported ones. It is usually an error to try
// to (un)marshal structs that only consist of unexported fields.
//
// This check will not flag calls involving types that define custom
// marshaling behavior, e.g. via MarshalJSON methods. It will also not
// flag empty structs.
//
// Available since
//     2019.2
//
//
// SA9006
//
// dubious bit shifting of a fixed size integer value
//
// Bit shifting a value past its size will always clear the value.
//
// For instance:
//
//     v := int8(42)
//     v >>= 8
//
// will always result in 0.
//
// This check flags bit shifting operations on fixed size integer values only.
// That is, int, uint and uintptr are never flagged to avoid potential false
// positives in somewhat exotic but valid bit twiddling tricks:
//
//     // Clear any value above 32 bits if integers are more than 32 bits.
//     func f(i int) int {
//         v := i >> 32
//         v = v << 32
//         return i-v
//     }
//
// Available since
//     2020.2
//
//
// SA9007
//
// deleting a directory that shouldn't be deleted
//
// It is virtually never correct to delete system directories such as
// /tmp or the user's home directory. However, it can be fairly easy to
// do by mistake, for example by mistakingly using os.TempDir instead
// of ioutil.TempDir, or by forgetting to add a suffix to the result
// of os.UserHomeDir.
//
// Writing
//
//     d := os.TempDir()
//     defer os.RemoveAll(d)
//
// in your unit tests will have a devastating effect on the stability of your system.
//
// This check flags attempts at deleting the following directories:
//
// - os.TempDir
// - os.UserCacheDir
// - os.UserConfigDir
// - os.UserHomeDir
//
// Available since
//     2022.1
//
//
// SA9008
//
// else branch of a type assertion is probably not reading the right value
//
// When declaring variables as part of an if statement (like in 'if
// foo := ...; foo {'), the same variables will also be in the scope of
// the else branch. This means that in the following example
//
//     if x, ok := x.(int); ok {
//         // ...
//     } else {
//         fmt.Printf("unexpected type %T", x)
//     }
//
// x in the else branch will refer to the x from x, ok
// :=; it will not refer to the x that is being type-asserted. The
// result of a failed type assertion is the zero value of the type that
// is being asserted to, so x in the else branch will always have the
// value 0 and the type int.
//
// Available since
//     2022.1
//
//
// SA9009
//
// ineffectual Go compiler directive
//
// A potential Go compiler directive was found, but is ineffectual as it begins
// with whitespace.
//
// Available since
//     2024.1
//
//
// EXPORTLOOPREF
//
// checks for pointers to enclosing loop variables
//
// THELPER
//
// report unmarked test helpers
//
//
package main
