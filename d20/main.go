/*
 * Day DD of AoC 2023
 *
 * Idea: Text
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package main

import (
	"aoc23/tools"
	"flag"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"
	"time"
)

var TESTMODE = true
var nt = flag.Bool("nt", false, "exec no test mode")
var inputfile = flag.String("f", "input.txt", "name of input file")
var part1only = flag.Bool("1", false, "run part 1 only")
var part2only = flag.Bool("2", false, "run part 2 only")

func main() {
	flag.Parse()
	if *nt {
		TESTMODE = false
	}
	log.SetPrefix("  ")
	log.SetFlags(0)

	if !*part2only {
		part01()
	}
	if !*part1only {
		part02()
	}
}

var testinput2 = `broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> output`

var testinput1 = `broadcaster -> a, b, c
%a -> b
%b -> c
%c -> inv
&inv -> a`

var testinput = testinput2

var cntLow int
var cntHigh int

const (
	low  int = 0
	high int = 1
)

type pulse struct {
	sender   string
	value    int // low, high
	receiver module
}

func (p pulse) String() string {
	return fmt.Sprintf("(%v -%v-> %v)", p.sender, p.value, p.receiver.getName())
}

type pulseQueue struct {
	values []pulse
}

func (pq pulseQueue) String() string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%v: ", len(pq.values)))
	for _, p := range pq.values {
		buf.WriteString(fmt.Sprintf("%v ", p))
	}
	buf.WriteString("]")
	return buf.String()
}

func (pq *pulseQueue) addPulses(sender string, pulseVal int, receivers []module) {
	for _, r := range receivers {
		p := pulse{sender, pulseVal, r}
		(*pq).values = append((*pq).values, p)
	}
}

type module interface {
	getName() string
	receive(string, int)
	addReceiver(r module)
	getReceivers() []module
	reset()
}

type base struct {
	name      string
	receivers []module
	queue     *pulseQueue
}

type flipflop struct {
	base
	on bool
}

func mkFlipflop(name string, pq *pulseQueue) *flipflop {
	ff := flipflop{
		base{name, make([]module, 0), pq},
		false,
	}
	return &ff
}

func (ff *flipflop) getName() string {
	return ff.name
}

func (ff *flipflop) reset() {
	ff.on = false
}

func (ff *flipflop) addReceiver(r module) {
	ff.receivers = append(ff.receivers, r)
}

func (ff *flipflop) getReceivers() []module {
	return ff.receivers
}

func (ff *flipflop) receive(name string, pulse int) {
	if pulse == low {
		ff.on = !ff.on
		newP := low
		if ff.on {
			newP = high
		}
		ff.queue.addPulses(ff.name, newP, ff.receivers)
	}
}
func (ff *flipflop) String() string {
	return fmt.Sprintf("%v: %v", ff.name, ff.on)
}

type conjunction struct {
	base
	lastPulse map[string]int
}

func mkConjunction(name string, pq *pulseQueue) *conjunction {
	c := conjunction{
		base{name, make([]module, 0), pq},
		make(map[string]int),
	}
	return &c
}

func (c *conjunction) getName() string {
	return c.name
}

func (c *conjunction) reset() {
	for k := range c.lastPulse {
		c.lastPulse[k] = 0
	}
}

func (c *conjunction) receive(name string, pulse int) {
	c.lastPulse[name] = pulse
	newP := low
	for _, v := range c.lastPulse {
		if v == low {
			newP = high
		}
	}
	c.queue.addPulses(c.name, newP, c.receivers)
}

func (c *conjunction) addReceiver(r module) {
	c.receivers = append(c.receivers, r)
}

func (c *conjunction) getReceivers() []module {
	return c.receivers
}
func (c *conjunction) String() string {
	return fmt.Sprintf("%v: %v", c.name, c.lastPulse)
}

type broadcaster struct {
	base
}

func (b *broadcaster) getName() string {
	return b.name
}

func (b *broadcaster) reset() {
}

func mkBroadcaster(name string, pq *pulseQueue) *broadcaster {
	bc := broadcaster{
		base{name, make([]module, 0), pq},
	}
	return &bc
}
func (b *broadcaster) addReceiver(r module) {
	b.receivers = append(b.receivers, r)
}

func (b *broadcaster) getReceivers() []module {
	return b.receivers
}
func (b *broadcaster) receive(name string, pulseVal int) {
	b.queue.addPulses(b.name, pulseVal, b.receivers)
}
func (b *broadcaster) String() string {
	return fmt.Sprintf("%v", b.name)
}

func readModules(pq *pulseQueue, printit bool) map[string]module {

	lines := getInput()

	allModules := make(map[string]module)

	modReceivers := make(map[string][]string)
	allReceivers := make(map[string]bool)

	for _, line := range lines {
		// log.Printf("Processing line %v with len %v\n", cnt, len(line))
		parts := strings.Split(line, " -> ")

		mod := parts[0]
		rcvs := strings.Split(parts[1], ", ")
		var mname = mod[1:]
		// fmt.Printf("%v: Read '%v' with: '%v'\n", cnt, mod, rcvs)
		if mod[0] == 'b' {
			mname = mod
			m := mkBroadcaster(mname, pq)
			allModules[m.name] = m
		} else if parts[0][0] == '%' {
			m := mkFlipflop(mname, pq)
			allModules[m.name] = m
		} else if parts[0][0] == '&' {
			m := mkConjunction(mname, pq)
			allModules[m.name] = m
		}
		modReceivers[mname] = rcvs
		for _, r := range rcvs {
			allReceivers[r] = true
		}
	}
	// ensure there is a module for all receivers
	for k := range allReceivers {
		if _, ok := allModules[k]; !ok {
			// the receiver dos not exist as a module - lets us add it as pure output
			allModules[k] = mkFlipflop(k, pq)
		}
	}
	// now link actual receivers to the modules
	for k, v := range allModules {
		if mrl, ok := modReceivers[k]; ok {
			for _, r := range mrl {
				rcv := allModules[r]
				v.addReceiver(rcv)
				if w, ok := rcv.(*conjunction); ok {
					w.lastPulse[v.getName()] = low
				}
			}
		}
	}

	if printit {
		// just as info
		for _, v := range allModules {
			r := v.getReceivers()
			fmt.Printf("Mod '%v' has %v receivers: %v\n", v, len(r), r)
		}
	}

	return allModules
}

func (pq *pulseQueue) step(name string) string {
	// take first elem of queue
	p := pq.values[0]
	pq.values = pq.values[1:]
	// and send
	r := p.receiver
	if p.value == high {
		cntHigh++
	} else {
		cntLow++
	}
	r.receive(p.sender, p.value)
	if name != "" && p.value == high && r.getName() == name {
		return p.sender
	}
	return ""
}

func part01() {
	startTime := time.Now()

	pq := pulseQueue{}
	allModules := readModules(&pq, false)

	bc := allModules["broadcaster"]
	cnt := 0
	for i := 0; i < 1000; i++ {
		cnt++
		pq.addPulses("button", low, []module{bc})
		for len(pq.values) > 0 {
			_ = pq.step("") // send empty string, as we are not interested in checking for registers
		}
	}
	total := cntLow * cntHigh
	elapsed := time.Since(startTime)
	fmt.Printf("Buttons: %v, High %v, Low: %v\n", cnt, cntHigh, cntLow)
	fmt.Printf("Result part 01 (%v): %v\n\n", elapsed, total)

}

func findSenders(nm string, allModules map[string]module) []module {
	senders := []module{}
	for _, m := range allModules {
		rcvs := m.getReceivers()
		for _, r := range rcvs {
			if r.getName() == nm {
				senders = append(senders, m)
			}
		}
	}
	return senders
}

func findTargets(nm string, allModules map[string]module) []module {
	prev := findSenders(nm, allModules)[0]
	targets := findSenders(prev.getName(), allModules)
	return targets
}

func part02() {
	startTime := time.Now()
	cnt := 0

	pq := pulseQueue{}
	allModules := readModules(&pq, false)

	sender := findSenders("rx", allModules)[0]
	targets := findTargets("rx", allModules)
	// fmt.Printf("Sender: %v, Targets %v\n\n", sender, targets)

	tgCounter := make(map[string]int)
	bc := allModules["broadcaster"]

	for len(tgCounter) != len(targets) {
		cnt++
		pq.addPulses("button", low, []module{bc})
		for len(pq.values) > 0 {
			s := pq.step(sender.getName())
			if s != "" {
				if _, ok := tgCounter[s]; !ok {
					tgCounter[s] = cnt
				}
			}
		}
	}

	vals := slices.Collect(maps.Values(tgCounter))
	lcm := tools.LCM(vals[0], vals[1], vals[2:]...)
	elapsed := time.Since(startTime)
	fmt.Printf("Result part 02 (%v): %v (%v buttons)\n\n", elapsed, lcm, cnt)
}

func getInput(inputs ...string) []string {
	if TESTMODE {
		input := testinput
		if len(inputs) > 0 {
			input = inputs[0]
		}
		return tools.ReadInputString(input)
	} else {
		return tools.ReadInputFile(*inputfile)
	}

}
