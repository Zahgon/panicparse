// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// panic crashes in various ways.
//
// It is a tool to help test pp, it is used in its unit tests.
//
// To install, run:
//
//	go install github.com/maruel/panicparse/v2/cmd/panic
//	panic -help
//	panic str |& pp
//
// Some panics require the race detector with -race:
//
//	go install -race github.com/maruel/panicparse/v2/cmd/panic
//	panic race |& pp
//
// To use with optimization (-N) and inlining (-l) disabled, build with
// -gcflags '-N -l' like:
//
//	go install -gcflags '-N -l' github.com/maruel/panicparse/v2/cmd/panic
package main

// To add a new panic stack signature, add it to types type below, keeping the
// list ordered by name. If you need utility functions, add it in the section
// below. That's it!

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/maruel/panicparse/v2/cmd/panic/internal"
	correct "github.com/maruel/panicparse/v2/cmd/panic/internal/incorrect"
	ùtf8 "github.com/maruel/panicparse/v2/cmd/panic/internal/utf8"
)

func main() {
	if len(os.Args) == 2 {
		switch n := os.Args[1]; n {
		case "-h", "-help", "--help", "help":
			usage()
			os.Exit(0)

		case "dump_commands":
			// Undocumented command to do a raw dump of the supported commands. This
			// is used by unit tests in ../../stack.
			items := make([]string, 0, len(types))
			for n := range types {
				items = append(items, n)
			}
			sort.Strings(items)
			for _, n := range items {
				fmt.Printf("%s\n", n)
			}
			os.Exit(0)

		default:
			if f, ok := types[n]; ok {
				fmt.Printf("GOTRACEBACK=%s\n", os.Getenv("GOTRACEBACK"))
				if n == "simple" {
					// Since the map lookup creates another call stack entry, add a
					// one-off "simple" panic style to test the very minimal case.
					// types["simple"].f is never called.
					panic("simple")
				}
				f.f()
				os.Exit(3)
			}
			fmt.Fprintf(stdErr, "unknown panic style %q\n", n)
			os.Exit(1)
		}
	}
	usage()
	os.Exit(1)
}

// Mocked in test.
var stdErr io.Writer = os.Stderr

// Utility functions.

func panicint(i int) { _ = "STUB: not implemented"; return }

func panicfloat64(f float64) { _ = "STUB: not implemented"; return }

func panicstr(a string) { _ = "STUB: not implemented"; return }

func panicslicestr(a []string) { _ = "STUB: not implemented"; return }

func panicArgsElided(a, b, c, d, e, f, g, h, i, j, k int) { _ = "STUB: not implemented"; return }

func recurse(i int) { _ = "STUB: not implemented"; return }

func panicRaceDisabled(name string) { _ = "STUB: not implemented"; return }

func rerunWithFastCrash() { _ = "STUB: not implemented"; return }

/* #nosec G204 */

// panicDoRaceWrite and panicDoRaceRead are extracted from panicRace() to make
// the stack trace less trivial, but in general folks will do the error with
// this code inlined.
func panicDoRaceWrite(x *int) { _ = "STUB: not implemented"; return }

func panicDoRaceRead(x *int) { _ = "STUB: not implemented"; return }

func panicRace() { _ = "STUB: not implemented"; return }

// Do two separate calls so that the 'created at' stacks are different.

//go:noinline
func panicChanStruct(x chan struct{}) {
	_ = "STUB: not implemented"

	/* TODO(maruel): This is not detected!
	   func panicRaceUnaligned() {
	   	if !raceEnabled {
	   		panicRaceDisabled("race_unaligned")
	   		return
	   	}
	   	rerunWithFastCrash()

	   	a := [8]byte{}
	   	b := (*int64)(unsafe.Pointer(&a[0]))
	   	go func() {
	   		for i := 0; ; i++ {
	   			a[4] = byte(i)
	   		}
	   	}()
	   	go func() {
	   		for {
	   			*b++
	   		}
	   	}()
	   	time.Sleep(time.Minute)
	   }
	*/return
}

//

// types is all the supported types of panics.
//
// Keep the list sorted.
//
// TODO(maruel): Figure out a way to reliably trigger "(scan)" output:
// - disable automatic GC with runtime.SetGCPercent(-1)
// - a goroutine with a large number of items in the stack
// - large heap to make the scanning process slow enough
// - trigger a manual GC with go runtime.GC()
// - panic in the meantime
// This would still not be deterministic.
//
// TODO(maruel): Figure out a way to reliably trigger sleep output.
var types = map[string]struct {
	desc string
	f    func()
}{
	"args_elided": {
		"too many args in stack line, causing the call arguments to be elided",
		func() {
			panicArgsElided(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
		},
	},

	"chan_receive": {
		"goroutine blocked on <-c",
		func() {
			c := make(chan bool)
			go func() {
				<-c
				<-c
			}()
			c <- true
			panic(42)
		},
	},

	"chan_send": {
		"goroutine blocked on c<-",
		func() {
			c := make(chan bool)
			go func() {
				c <- true
				c <- true
			}()
			<-c
			panic(42)
		},
	},

	"chan_struct": {
		"panic with an empty chan struct{} as a parameter",
		func() {
			panicChanStruct(nil)
		},
	},

	"float": {
		"panic(4.2)",
		func() {
			panicfloat64(4.2)
		},
	},

	"goroutine_1": {
		"panic in one goroutine",
		func() {
			go func() {
				panicint(42)
			}()
			time.Sleep(time.Minute)
		},
	},

	"goroutine_100": {
		"start 100 goroutines before panicking",
		func() {
			var wg sync.WaitGroup
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					wg.Done()
					time.Sleep(time.Minute)
				}()
			}
			wg.Wait()
			panicint(42)
		},
	},

	"goroutine_dedupe_pointers": {
		"start 100 goroutines with different pointers before panicking",
		func() {
			var wg sync.WaitGroup
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func(b *int) {
					wg.Done()
					time.Sleep(time.Minute)
				}(new(int))
			}
			wg.Wait()
			panicint(42)
		},
	},

	"int": {
		"panic(42)",
		func() {
			panicint(42)
		},
	},

	"locked": {
		"thread locked goroutine via runtime.LockOSThread()",
		func() {
			runtime.LockOSThread()
			panic(42)
		},
	},

	"other": {
		"panics with other package in the call stack, with both exported and unexpected functions",
		func() {
			internal.Callback(func() {
				panic("allo")
			})
		},
	},

	"asleep": {
		"panics with 'all goroutines are asleep - deadlock'",
		func() {
			// When built with the race detector, this hangs. I suspect this is due
			// because the race detector starts a separate goroutine which causes
			// checkdead() to not trigger. See checkdead() in src/runtime/proc.go.
			// https://github.com/golang/go/issues/20588
			//
			// Repro:
			//   go install -race github.com/maruel/panicparse/v2/cmd/panic; panic asleep
			var mu sync.Mutex
			mu.Lock()
			mu.Lock()
		},
	},

	"race": {
		"cause a crash by race detector",
		panicRace,
	},

	/* TODO(maruel): This is not detected!
	"race_unaligned": {
		"cause a crash by race detector with unaligned access",
		panicRaceUnaligned,
	},
	*/

	"stack_cut_off": {
		"recursive calls with too many call lines in traceback, causing higher up calls to missing",
		func() {
			// Observed limit is 99.
			// src/runtime/runtime2.go:const _TracebackMaxFrames = 100
			recurse(100)
		},
	},

	"stack_cut_off_named": {
		"named calls with too many call lines in traceback, causing higher up calls to missing",
		func() {
			// The tool has difficulty with very deep static calls during linking.
			// See https://github.com/golang/go/issues/51814
			recurse497()
		},
	},

	"simple": {
		// This is not used for real, here for documentation.
		"skip the map for a shorter stack trace",
		func() {},
	},

	"slice_str": {
		"panic([]string{\"allo\"}) with cap=2",
		func() {
			a := make([]string, 1, 2)
			a[0] = "allo"
			panicslicestr(a)
		},
	},

	"stdlib": {
		"panics with stdlib in the call stack, with both exported and unexpected functions",
		func() {
			strings.FieldsFunc("a", func(rune) bool {
				panic("allo")
			})
		},
	},

	"stdlib_and_other": {
		"panics with both other and stdlib packages in the call stack",
		func() {
			strings.FieldsFunc("a", func(rune) bool {
				internal.Callback(func() {
					panic("allo")
				})
				return false
			})
		},
	},

	"str": {
		"panic(\"allo\")",
		func() {
			panicstr("allo")
		},
	},

	"mismatched": {
		"mismatched package and directory names",
		func() {
			correct.Panic()
		},
	},

	"utf8": {
		"non-ascii package, struct and method names",
		func() {
			s := ùtf8.Strùct{}
			s.Pànic()
		},
	},
}

func usage() { _ = "STUB: not implemented"; return }

//

func recurse00()  { _ = "STUB: not implemented"; return }
func recurse01()  { _ = "STUB: not implemented"; return }
func recurse02()  { _ = "STUB: not implemented"; return }
func recurse03()  { _ = "STUB: not implemented"; return }
func recurse04()  { _ = "STUB: not implemented"; return }
func recurse05()  { _ = "STUB: not implemented"; return }
func recurse06()  { _ = "STUB: not implemented"; return }
func recurse07()  { _ = "STUB: not implemented"; return }
func recurse08()  { _ = "STUB: not implemented"; return }
func recurse09()  { _ = "STUB: not implemented"; return }
func recurse10()  { _ = "STUB: not implemented"; return }
func recurse11()  { _ = "STUB: not implemented"; return }
func recurse12()  { _ = "STUB: not implemented"; return }
func recurse13()  { _ = "STUB: not implemented"; return }
func recurse14()  { _ = "STUB: not implemented"; return }
func recurse15()  { _ = "STUB: not implemented"; return }
func recurse16()  { _ = "STUB: not implemented"; return }
func recurse17()  { _ = "STUB: not implemented"; return }
func recurse18()  { _ = "STUB: not implemented"; return }
func recurse19()  { _ = "STUB: not implemented"; return }
func recurse20()  { _ = "STUB: not implemented"; return }
func recurse21()  { _ = "STUB: not implemented"; return }
func recurse22()  { _ = "STUB: not implemented"; return }
func recurse23()  { _ = "STUB: not implemented"; return }
func recurse24()  { _ = "STUB: not implemented"; return }
func recurse25()  { _ = "STUB: not implemented"; return }
func recurse26()  { _ = "STUB: not implemented"; return }
func recurse27()  { _ = "STUB: not implemented"; return }
func recurse28()  { _ = "STUB: not implemented"; return }
func recurse29()  { _ = "STUB: not implemented"; return }
func recurse30()  { _ = "STUB: not implemented"; return }
func recurse31()  { _ = "STUB: not implemented"; return }
func recurse32()  { _ = "STUB: not implemented"; return }
func recurse33()  { _ = "STUB: not implemented"; return }
func recurse34()  { _ = "STUB: not implemented"; return }
func recurse35()  { _ = "STUB: not implemented"; return }
func recurse36()  { _ = "STUB: not implemented"; return }
func recurse37()  { _ = "STUB: not implemented"; return }
func recurse38()  { _ = "STUB: not implemented"; return }
func recurse39()  { _ = "STUB: not implemented"; return }
func recurse40()  { _ = "STUB: not implemented"; return }
func recurse41()  { _ = "STUB: not implemented"; return }
func recurse42()  { _ = "STUB: not implemented"; return }
func recurse43()  { _ = "STUB: not implemented"; return }
func recurse44()  { _ = "STUB: not implemented"; return }
func recurse45()  { _ = "STUB: not implemented"; return }
func recurse46()  { _ = "STUB: not implemented"; return }
func recurse47()  { _ = "STUB: not implemented"; return }
func recurse48()  { _ = "STUB: not implemented"; return }
func recurse49()  { _ = "STUB: not implemented"; return }
func recurse50()  { _ = "STUB: not implemented"; return }
func recurse51()  { _ = "STUB: not implemented"; return }
func recurse52()  { _ = "STUB: not implemented"; return }
func recurse53()  { _ = "STUB: not implemented"; return }
func recurse54()  { _ = "STUB: not implemented"; return }
func recurse55()  { _ = "STUB: not implemented"; return }
func recurse56()  { _ = "STUB: not implemented"; return }
func recurse57()  { _ = "STUB: not implemented"; return }
func recurse58()  { _ = "STUB: not implemented"; return }
func recurse59()  { _ = "STUB: not implemented"; return }
func recurse60()  { _ = "STUB: not implemented"; return }
func recurse61()  { _ = "STUB: not implemented"; return }
func recurse62()  { _ = "STUB: not implemented"; return }
func recurse63()  { _ = "STUB: not implemented"; return }
func recurse64()  { _ = "STUB: not implemented"; return }
func recurse65()  { _ = "STUB: not implemented"; return }
func recurse66()  { _ = "STUB: not implemented"; return }
func recurse67()  { _ = "STUB: not implemented"; return }
func recurse68()  { _ = "STUB: not implemented"; return }
func recurse69()  { _ = "STUB: not implemented"; return }
func recurse70()  { _ = "STUB: not implemented"; return }
func recurse71()  { _ = "STUB: not implemented"; return }
func recurse72()  { _ = "STUB: not implemented"; return }
func recurse73()  { _ = "STUB: not implemented"; return }
func recurse74()  { _ = "STUB: not implemented"; return }
func recurse75()  { _ = "STUB: not implemented"; return }
func recurse76()  { _ = "STUB: not implemented"; return }
func recurse77()  { _ = "STUB: not implemented"; return }
func recurse78()  { _ = "STUB: not implemented"; return }
func recurse79()  { _ = "STUB: not implemented"; return }
func recurse80()  { _ = "STUB: not implemented"; return }
func recurse81()  { _ = "STUB: not implemented"; return }
func recurse82()  { _ = "STUB: not implemented"; return }
func recurse83()  { _ = "STUB: not implemented"; return }
func recurse84()  { _ = "STUB: not implemented"; return }
func recurse85()  { _ = "STUB: not implemented"; return }
func recurse86()  { _ = "STUB: not implemented"; return }
func recurse87()  { _ = "STUB: not implemented"; return }
func recurse88()  { _ = "STUB: not implemented"; return }
func recurse89()  { _ = "STUB: not implemented"; return }
func recurse90()  { _ = "STUB: not implemented"; return }
func recurse91()  { _ = "STUB: not implemented"; return }
func recurse92()  { _ = "STUB: not implemented"; return }
func recurse93()  { _ = "STUB: not implemented"; return }
func recurse94()  { _ = "STUB: not implemented"; return }
func recurse95()  { _ = "STUB: not implemented"; return }
func recurse96()  { _ = "STUB: not implemented"; return }
func recurse97()  { _ = "STUB: not implemented"; return }
func recurse98()  { _ = "STUB: not implemented"; return }
func recurse99()  { _ = "STUB: not implemented"; return }
func recurse100() { _ = "STUB: not implemented"; return }

func recurse101() { _ = "STUB: not implemented"; return }
func recurse102() { _ = "STUB: not implemented"; return }
func recurse103() { _ = "STUB: not implemented"; return }
func recurse104() { _ = "STUB: not implemented"; return }
func recurse105() { _ = "STUB: not implemented"; return }
func recurse106() { _ = "STUB: not implemented"; return }
func recurse107() { _ = "STUB: not implemented"; return }
func recurse108() { _ = "STUB: not implemented"; return }
func recurse109() { _ = "STUB: not implemented"; return }
func recurse110() { _ = "STUB: not implemented"; return }
func recurse111() { _ = "STUB: not implemented"; return }
func recurse112() { _ = "STUB: not implemented"; return }
func recurse113() { _ = "STUB: not implemented"; return }
func recurse114() { _ = "STUB: not implemented"; return }
func recurse115() { _ = "STUB: not implemented"; return }
func recurse116() { _ = "STUB: not implemented"; return }
func recurse117() { _ = "STUB: not implemented"; return }
func recurse118() { _ = "STUB: not implemented"; return }
func recurse119() { _ = "STUB: not implemented"; return }
func recurse120() { _ = "STUB: not implemented"; return }
func recurse121() { _ = "STUB: not implemented"; return }
func recurse122() { _ = "STUB: not implemented"; return }
func recurse123() { _ = "STUB: not implemented"; return }
func recurse124() { _ = "STUB: not implemented"; return }
func recurse125() { _ = "STUB: not implemented"; return }
func recurse126() { _ = "STUB: not implemented"; return }
func recurse127() { _ = "STUB: not implemented"; return }
func recurse128() { _ = "STUB: not implemented"; return }
func recurse129() { _ = "STUB: not implemented"; return }
func recurse130() { _ = "STUB: not implemented"; return }
func recurse131() { _ = "STUB: not implemented"; return }
func recurse132() { _ = "STUB: not implemented"; return }
func recurse133() { _ = "STUB: not implemented"; return }
func recurse134() { _ = "STUB: not implemented"; return }
func recurse135() { _ = "STUB: not implemented"; return }
func recurse136() { _ = "STUB: not implemented"; return }
func recurse137() { _ = "STUB: not implemented"; return }
func recurse138() { _ = "STUB: not implemented"; return }
func recurse139() { _ = "STUB: not implemented"; return }
func recurse140() { _ = "STUB: not implemented"; return }
func recurse141() { _ = "STUB: not implemented"; return }
func recurse142() { _ = "STUB: not implemented"; return }
func recurse143() { _ = "STUB: not implemented"; return }
func recurse144() { _ = "STUB: not implemented"; return }
func recurse145() { _ = "STUB: not implemented"; return }
func recurse146() { _ = "STUB: not implemented"; return }
func recurse147() { _ = "STUB: not implemented"; return }
func recurse148() { _ = "STUB: not implemented"; return }
func recurse149() { _ = "STUB: not implemented"; return }
func recurse150() { _ = "STUB: not implemented"; return }
func recurse151() { _ = "STUB: not implemented"; return }
func recurse152() { _ = "STUB: not implemented"; return }
func recurse153() { _ = "STUB: not implemented"; return }
func recurse154() { _ = "STUB: not implemented"; return }
func recurse155() { _ = "STUB: not implemented"; return }
func recurse156() { _ = "STUB: not implemented"; return }
func recurse157() { _ = "STUB: not implemented"; return }
func recurse158() { _ = "STUB: not implemented"; return }
func recurse159() { _ = "STUB: not implemented"; return }
func recurse160() { _ = "STUB: not implemented"; return }
func recurse161() { _ = "STUB: not implemented"; return }
func recurse162() { _ = "STUB: not implemented"; return }
func recurse163() { _ = "STUB: not implemented"; return }
func recurse164() { _ = "STUB: not implemented"; return }
func recurse165() { _ = "STUB: not implemented"; return }
func recurse166() { _ = "STUB: not implemented"; return }
func recurse167() { _ = "STUB: not implemented"; return }
func recurse168() { _ = "STUB: not implemented"; return }
func recurse169() { _ = "STUB: not implemented"; return }
func recurse170() { _ = "STUB: not implemented"; return }
func recurse171() { _ = "STUB: not implemented"; return }
func recurse172() { _ = "STUB: not implemented"; return }
func recurse173() { _ = "STUB: not implemented"; return }
func recurse174() { _ = "STUB: not implemented"; return }
func recurse175() { _ = "STUB: not implemented"; return }
func recurse176() { _ = "STUB: not implemented"; return }
func recurse177() { _ = "STUB: not implemented"; return }
func recurse178() { _ = "STUB: not implemented"; return }
func recurse179() { _ = "STUB: not implemented"; return }
func recurse180() { _ = "STUB: not implemented"; return }
func recurse181() { _ = "STUB: not implemented"; return }
func recurse182() { _ = "STUB: not implemented"; return }
func recurse183() { _ = "STUB: not implemented"; return }
func recurse184() { _ = "STUB: not implemented"; return }
func recurse185() { _ = "STUB: not implemented"; return }
func recurse186() { _ = "STUB: not implemented"; return }
func recurse187() { _ = "STUB: not implemented"; return }
func recurse188() { _ = "STUB: not implemented"; return }
func recurse189() { _ = "STUB: not implemented"; return }
func recurse190() { _ = "STUB: not implemented"; return }
func recurse191() { _ = "STUB: not implemented"; return }
func recurse192() { _ = "STUB: not implemented"; return }
func recurse193() { _ = "STUB: not implemented"; return }
func recurse194() { _ = "STUB: not implemented"; return }
func recurse195() { _ = "STUB: not implemented"; return }
func recurse196() { _ = "STUB: not implemented"; return }
func recurse197() { _ = "STUB: not implemented"; return }
func recurse198() { _ = "STUB: not implemented"; return }
func recurse199() { _ = "STUB: not implemented"; return }

func recurse200() { _ = "STUB: not implemented"; return }
func recurse201() { _ = "STUB: not implemented"; return }
func recurse202() { _ = "STUB: not implemented"; return }
func recurse203() { _ = "STUB: not implemented"; return }
func recurse204() { _ = "STUB: not implemented"; return }
func recurse205() { _ = "STUB: not implemented"; return }
func recurse206() { _ = "STUB: not implemented"; return }
func recurse207() { _ = "STUB: not implemented"; return }
func recurse208() { _ = "STUB: not implemented"; return }
func recurse209() { _ = "STUB: not implemented"; return }
func recurse210() { _ = "STUB: not implemented"; return }
func recurse211() { _ = "STUB: not implemented"; return }
func recurse212() { _ = "STUB: not implemented"; return }
func recurse213() { _ = "STUB: not implemented"; return }
func recurse214() { _ = "STUB: not implemented"; return }
func recurse215() { _ = "STUB: not implemented"; return }
func recurse216() { _ = "STUB: not implemented"; return }
func recurse217() { _ = "STUB: not implemented"; return }
func recurse218() { _ = "STUB: not implemented"; return }
func recurse219() { _ = "STUB: not implemented"; return }
func recurse220() { _ = "STUB: not implemented"; return }
func recurse221() { _ = "STUB: not implemented"; return }
func recurse222() { _ = "STUB: not implemented"; return }
func recurse223() { _ = "STUB: not implemented"; return }
func recurse224() { _ = "STUB: not implemented"; return }
func recurse225() { _ = "STUB: not implemented"; return }
func recurse226() { _ = "STUB: not implemented"; return }
func recurse227() { _ = "STUB: not implemented"; return }
func recurse228() { _ = "STUB: not implemented"; return }
func recurse229() { _ = "STUB: not implemented"; return }
func recurse230() { _ = "STUB: not implemented"; return }
func recurse231() { _ = "STUB: not implemented"; return }
func recurse232() { _ = "STUB: not implemented"; return }
func recurse233() { _ = "STUB: not implemented"; return }
func recurse234() { _ = "STUB: not implemented"; return }
func recurse235() { _ = "STUB: not implemented"; return }
func recurse236() { _ = "STUB: not implemented"; return }
func recurse237() { _ = "STUB: not implemented"; return }
func recurse238() { _ = "STUB: not implemented"; return }
func recurse239() { _ = "STUB: not implemented"; return }
func recurse240() { _ = "STUB: not implemented"; return }
func recurse241() { _ = "STUB: not implemented"; return }
func recurse242() { _ = "STUB: not implemented"; return }
func recurse243() { _ = "STUB: not implemented"; return }
func recurse244() { _ = "STUB: not implemented"; return }
func recurse245() { _ = "STUB: not implemented"; return }
func recurse246() { _ = "STUB: not implemented"; return }
func recurse247() { _ = "STUB: not implemented"; return }
func recurse248() { _ = "STUB: not implemented"; return }
func recurse249() { _ = "STUB: not implemented"; return }
func recurse250() { _ = "STUB: not implemented"; return }
func recurse251() { _ = "STUB: not implemented"; return }
func recurse252() { _ = "STUB: not implemented"; return }
func recurse253() { _ = "STUB: not implemented"; return }
func recurse254() { _ = "STUB: not implemented"; return }
func recurse255() { _ = "STUB: not implemented"; return }
func recurse256() { _ = "STUB: not implemented"; return }
func recurse257() { _ = "STUB: not implemented"; return }
func recurse258() { _ = "STUB: not implemented"; return }
func recurse259() { _ = "STUB: not implemented"; return }
func recurse260() { _ = "STUB: not implemented"; return }
func recurse261() { _ = "STUB: not implemented"; return }
func recurse262() { _ = "STUB: not implemented"; return }
func recurse263() { _ = "STUB: not implemented"; return }
func recurse264() { _ = "STUB: not implemented"; return }
func recurse265() { _ = "STUB: not implemented"; return }
func recurse266() { _ = "STUB: not implemented"; return }
func recurse267() { _ = "STUB: not implemented"; return }
func recurse268() { _ = "STUB: not implemented"; return }
func recurse269() { _ = "STUB: not implemented"; return }
func recurse270() { _ = "STUB: not implemented"; return }
func recurse271() { _ = "STUB: not implemented"; return }
func recurse272() { _ = "STUB: not implemented"; return }
func recurse273() { _ = "STUB: not implemented"; return }
func recurse274() { _ = "STUB: not implemented"; return }
func recurse275() { _ = "STUB: not implemented"; return }
func recurse276() { _ = "STUB: not implemented"; return }
func recurse277() { _ = "STUB: not implemented"; return }
func recurse278() { _ = "STUB: not implemented"; return }
func recurse279() { _ = "STUB: not implemented"; return }
func recurse280() { _ = "STUB: not implemented"; return }
func recurse281() { _ = "STUB: not implemented"; return }
func recurse282() { _ = "STUB: not implemented"; return }
func recurse283() { _ = "STUB: not implemented"; return }
func recurse284() { _ = "STUB: not implemented"; return }
func recurse285() { _ = "STUB: not implemented"; return }
func recurse286() { _ = "STUB: not implemented"; return }
func recurse287() { _ = "STUB: not implemented"; return }
func recurse288() { _ = "STUB: not implemented"; return }
func recurse289() { _ = "STUB: not implemented"; return }
func recurse290() { _ = "STUB: not implemented"; return }
func recurse291() { _ = "STUB: not implemented"; return }
func recurse292() { _ = "STUB: not implemented"; return }
func recurse293() { _ = "STUB: not implemented"; return }
func recurse294() { _ = "STUB: not implemented"; return }
func recurse295() { _ = "STUB: not implemented"; return }
func recurse296() { _ = "STUB: not implemented"; return }
func recurse297() { _ = "STUB: not implemented"; return }
func recurse298() { _ = "STUB: not implemented"; return }
func recurse299() { _ = "STUB: not implemented"; return }

func recurse300() { _ = "STUB: not implemented"; return }
func recurse301() { _ = "STUB: not implemented"; return }
func recurse302() { _ = "STUB: not implemented"; return }
func recurse303() { _ = "STUB: not implemented"; return }
func recurse304() { _ = "STUB: not implemented"; return }
func recurse305() { _ = "STUB: not implemented"; return }
func recurse306() { _ = "STUB: not implemented"; return }
func recurse307() { _ = "STUB: not implemented"; return }
func recurse308() { _ = "STUB: not implemented"; return }
func recurse309() { _ = "STUB: not implemented"; return }
func recurse310() { _ = "STUB: not implemented"; return }
func recurse311() { _ = "STUB: not implemented"; return }
func recurse312() { _ = "STUB: not implemented"; return }
func recurse313() { _ = "STUB: not implemented"; return }
func recurse314() { _ = "STUB: not implemented"; return }
func recurse315() { _ = "STUB: not implemented"; return }
func recurse316() { _ = "STUB: not implemented"; return }
func recurse317() { _ = "STUB: not implemented"; return }
func recurse318() { _ = "STUB: not implemented"; return }
func recurse319() { _ = "STUB: not implemented"; return }
func recurse320() { _ = "STUB: not implemented"; return }
func recurse321() { _ = "STUB: not implemented"; return }
func recurse322() { _ = "STUB: not implemented"; return }
func recurse323() { _ = "STUB: not implemented"; return }
func recurse324() { _ = "STUB: not implemented"; return }
func recurse325() { _ = "STUB: not implemented"; return }
func recurse326() { _ = "STUB: not implemented"; return }
func recurse327() { _ = "STUB: not implemented"; return }
func recurse328() { _ = "STUB: not implemented"; return }
func recurse329() { _ = "STUB: not implemented"; return }
func recurse330() { _ = "STUB: not implemented"; return }
func recurse331() { _ = "STUB: not implemented"; return }
func recurse332() { _ = "STUB: not implemented"; return }
func recurse333() { _ = "STUB: not implemented"; return }
func recurse334() { _ = "STUB: not implemented"; return }
func recurse335() { _ = "STUB: not implemented"; return }
func recurse336() { _ = "STUB: not implemented"; return }
func recurse337() { _ = "STUB: not implemented"; return }
func recurse338() { _ = "STUB: not implemented"; return }
func recurse339() { _ = "STUB: not implemented"; return }
func recurse340() { _ = "STUB: not implemented"; return }
func recurse341() { _ = "STUB: not implemented"; return }
func recurse342() { _ = "STUB: not implemented"; return }
func recurse343() { _ = "STUB: not implemented"; return }
func recurse344() { _ = "STUB: not implemented"; return }
func recurse345() { _ = "STUB: not implemented"; return }
func recurse346() { _ = "STUB: not implemented"; return }
func recurse347() { _ = "STUB: not implemented"; return }
func recurse348() { _ = "STUB: not implemented"; return }
func recurse349() { _ = "STUB: not implemented"; return }
func recurse350() { _ = "STUB: not implemented"; return }
func recurse351() { _ = "STUB: not implemented"; return }
func recurse352() { _ = "STUB: not implemented"; return }
func recurse353() { _ = "STUB: not implemented"; return }
func recurse354() { _ = "STUB: not implemented"; return }
func recurse355() { _ = "STUB: not implemented"; return }
func recurse356() { _ = "STUB: not implemented"; return }
func recurse357() { _ = "STUB: not implemented"; return }
func recurse358() { _ = "STUB: not implemented"; return }
func recurse359() { _ = "STUB: not implemented"; return }
func recurse360() { _ = "STUB: not implemented"; return }
func recurse361() { _ = "STUB: not implemented"; return }
func recurse362() { _ = "STUB: not implemented"; return }
func recurse363() { _ = "STUB: not implemented"; return }
func recurse364() { _ = "STUB: not implemented"; return }
func recurse365() { _ = "STUB: not implemented"; return }
func recurse366() { _ = "STUB: not implemented"; return }
func recurse367() { _ = "STUB: not implemented"; return }
func recurse368() { _ = "STUB: not implemented"; return }
func recurse369() { _ = "STUB: not implemented"; return }
func recurse370() { _ = "STUB: not implemented"; return }
func recurse371() { _ = "STUB: not implemented"; return }
func recurse372() { _ = "STUB: not implemented"; return }
func recurse373() { _ = "STUB: not implemented"; return }
func recurse374() { _ = "STUB: not implemented"; return }
func recurse375() { _ = "STUB: not implemented"; return }
func recurse376() { _ = "STUB: not implemented"; return }
func recurse377() { _ = "STUB: not implemented"; return }
func recurse378() { _ = "STUB: not implemented"; return }
func recurse379() { _ = "STUB: not implemented"; return }
func recurse380() { _ = "STUB: not implemented"; return }
func recurse381() { _ = "STUB: not implemented"; return }
func recurse382() { _ = "STUB: not implemented"; return }
func recurse383() { _ = "STUB: not implemented"; return }
func recurse384() { _ = "STUB: not implemented"; return }
func recurse385() { _ = "STUB: not implemented"; return }
func recurse386() { _ = "STUB: not implemented"; return }
func recurse387() { _ = "STUB: not implemented"; return }
func recurse388() { _ = "STUB: not implemented"; return }
func recurse389() { _ = "STUB: not implemented"; return }
func recurse390() { _ = "STUB: not implemented"; return }
func recurse391() { _ = "STUB: not implemented"; return }
func recurse392() { _ = "STUB: not implemented"; return }
func recurse393() { _ = "STUB: not implemented"; return }
func recurse394() { _ = "STUB: not implemented"; return }
func recurse395() { _ = "STUB: not implemented"; return }
func recurse396() { _ = "STUB: not implemented"; return }
func recurse397() { _ = "STUB: not implemented"; return }
func recurse398() { _ = "STUB: not implemented"; return }
func recurse399() { _ = "STUB: not implemented"; return }

func recurse400() { _ = "STUB: not implemented"; return }
func recurse401() { _ = "STUB: not implemented"; return }
func recurse402() { _ = "STUB: not implemented"; return }
func recurse403() { _ = "STUB: not implemented"; return }
func recurse404() { _ = "STUB: not implemented"; return }
func recurse405() { _ = "STUB: not implemented"; return }
func recurse406() { _ = "STUB: not implemented"; return }
func recurse407() { _ = "STUB: not implemented"; return }
func recurse408() { _ = "STUB: not implemented"; return }
func recurse409() { _ = "STUB: not implemented"; return }
func recurse410() { _ = "STUB: not implemented"; return }
func recurse411() { _ = "STUB: not implemented"; return }
func recurse412() { _ = "STUB: not implemented"; return }
func recurse413() { _ = "STUB: not implemented"; return }
func recurse414() { _ = "STUB: not implemented"; return }
func recurse415() { _ = "STUB: not implemented"; return }
func recurse416() { _ = "STUB: not implemented"; return }
func recurse417() { _ = "STUB: not implemented"; return }
func recurse418() { _ = "STUB: not implemented"; return }
func recurse419() { _ = "STUB: not implemented"; return }
func recurse420() { _ = "STUB: not implemented"; return }
func recurse421() { _ = "STUB: not implemented"; return }
func recurse422() { _ = "STUB: not implemented"; return }
func recurse423() { _ = "STUB: not implemented"; return }
func recurse424() { _ = "STUB: not implemented"; return }
func recurse425() { _ = "STUB: not implemented"; return }
func recurse426() { _ = "STUB: not implemented"; return }
func recurse427() { _ = "STUB: not implemented"; return }
func recurse428() { _ = "STUB: not implemented"; return }
func recurse429() { _ = "STUB: not implemented"; return }
func recurse430() { _ = "STUB: not implemented"; return }
func recurse431() { _ = "STUB: not implemented"; return }
func recurse432() { _ = "STUB: not implemented"; return }
func recurse433() { _ = "STUB: not implemented"; return }
func recurse434() { _ = "STUB: not implemented"; return }
func recurse435() { _ = "STUB: not implemented"; return }
func recurse436() { _ = "STUB: not implemented"; return }
func recurse437() { _ = "STUB: not implemented"; return }
func recurse438() { _ = "STUB: not implemented"; return }
func recurse439() { _ = "STUB: not implemented"; return }
func recurse440() { _ = "STUB: not implemented"; return }
func recurse441() { _ = "STUB: not implemented"; return }
func recurse442() { _ = "STUB: not implemented"; return }
func recurse443() { _ = "STUB: not implemented"; return }
func recurse444() { _ = "STUB: not implemented"; return }
func recurse445() { _ = "STUB: not implemented"; return }
func recurse446() { _ = "STUB: not implemented"; return }
func recurse447() { _ = "STUB: not implemented"; return }
func recurse448() { _ = "STUB: not implemented"; return }
func recurse449() { _ = "STUB: not implemented"; return }
func recurse450() { _ = "STUB: not implemented"; return }
func recurse451() { _ = "STUB: not implemented"; return }
func recurse452() { _ = "STUB: not implemented"; return }
func recurse453() { _ = "STUB: not implemented"; return }
func recurse454() { _ = "STUB: not implemented"; return }
func recurse455() { _ = "STUB: not implemented"; return }
func recurse456() { _ = "STUB: not implemented"; return }
func recurse457() { _ = "STUB: not implemented"; return }
func recurse458() { _ = "STUB: not implemented"; return }
func recurse459() { _ = "STUB: not implemented"; return }
func recurse460() { _ = "STUB: not implemented"; return }
func recurse461() { _ = "STUB: not implemented"; return }
func recurse462() { _ = "STUB: not implemented"; return }
func recurse463() { _ = "STUB: not implemented"; return }
func recurse464() { _ = "STUB: not implemented"; return }
func recurse465() { _ = "STUB: not implemented"; return }
func recurse466() { _ = "STUB: not implemented"; return }
func recurse467() { _ = "STUB: not implemented"; return }
func recurse468() { _ = "STUB: not implemented"; return }
func recurse469() { _ = "STUB: not implemented"; return }
func recurse470() { _ = "STUB: not implemented"; return }
func recurse471() { _ = "STUB: not implemented"; return }
func recurse472() { _ = "STUB: not implemented"; return }
func recurse473() { _ = "STUB: not implemented"; return }
func recurse474() { _ = "STUB: not implemented"; return }
func recurse475() { _ = "STUB: not implemented"; return }
func recurse476() { _ = "STUB: not implemented"; return }
func recurse477() { _ = "STUB: not implemented"; return }
func recurse478() { _ = "STUB: not implemented"; return }
func recurse479() { _ = "STUB: not implemented"; return }
func recurse480() { _ = "STUB: not implemented"; return }
func recurse481() { _ = "STUB: not implemented"; return }
func recurse482() { _ = "STUB: not implemented"; return }
func recurse483() { _ = "STUB: not implemented"; return }
func recurse484() { _ = "STUB: not implemented"; return }
func recurse485() { _ = "STUB: not implemented"; return }
func recurse486() { _ = "STUB: not implemented"; return }
func recurse487() { _ = "STUB: not implemented"; return }
func recurse488() { _ = "STUB: not implemented"; return }
func recurse489() { _ = "STUB: not implemented"; return }
func recurse490() { _ = "STUB: not implemented"; return }
func recurse491() { _ = "STUB: not implemented"; return }
func recurse492() { _ = "STUB: not implemented"; return }
func recurse493() { _ = "STUB: not implemented"; return }
func recurse494() { _ = "STUB: not implemented"; return }
func recurse495() { _ = "STUB: not implemented"; return }
func recurse496() { _ = "STUB: not implemented"; return }
func recurse497() { _ = "STUB: not implemented"; return }
