package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	"github.com/marcusolsson/tui-go"
)

type stack []float64

// RPMEngine is the engine for a RPN calculator
type RPMEngine struct {
	stack        stack
	catalog      map[string]interface{}
	mode         string
	haserror     bool
	precision    int
	history      []string
	historyindex int
}

// Unit constants
const (
	RAD = "RAD"
	DEG = "DEG"
)

func (s stack) Push(v float64) stack {
	return append(s, v)
}

func (s stack) Len() int {
	return len(s)
}

func (s stack) Pop() (stack, float64) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func (s stack) Dup() stack {
	l := len(s)
	return append(s, s[l-1])
}

func (s stack) Clear() stack {
	return []float64{}
}

func main() {
	engine := NewRPMEngine()

	theme := tui.NewTheme()
	normal := tui.Style{Bg: tui.ColorWhite, Fg: tui.ColorBlack}
	theme.SetStyle("normal", normal)
	theme.SetStyle("label.warning", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorRed})
	theme.SetStyle("label.unit", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorCyan})

	stack := tui.NewVBox()

	stackScroll := tui.NewScrollArea(stack)
	//stackScroll.SetAutoscrollToBottom(true)

	stackBox := tui.NewVBox(stackScroll)
	stackBox.SetBorder(true)

	inputBox := tui.NewEntry()
	inputBox.SetFocused(true)

	inputOuter := tui.NewHBox(inputBox)
	inputOuter.SetBorder(true)
	inputOuter.SetSizePolicy(tui.Expanding, tui.Maximum)

	angleLabel := tui.NewLabel("RAD")
	errorLabel := tui.NewLabel(" ")
	precisionLabel := tui.NewLabel(" ")
	angleLabel.SetStyleName("unit")
	precisionLabel.SetStyleName("unit")

	infoBox := tui.NewHBox(angleLabel)
	infoBox.SetBorder(true)

	precisionBox := tui.NewHBox(precisionLabel)
	precisionBox.SetBorder(true)

	errorBox := tui.NewHBox(errorLabel)
	errorBox.SetBorder(true)

	mainBox := tui.NewHBox(inputOuter, infoBox, precisionBox, errorBox)
	mainBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	root := tui.NewVBox(stackBox, mainBox)
	ui := tui.New(root)
	ui.SetTheme(theme)

	// inputBox.SetBorder(true)
	inputBox.OnSubmit(func(e *tui.Entry) {
		if "quit" == e.Text() {
			ui.Quit()
		}

		analyseInput(engine, e.Text())

		// Repaint angle units
		angleLabel.SetText(engine.mode)

		// Repaint precision box if needed
		if engine.precision != -1 {
			precisionLabel.SetText(strconv.Itoa(engine.precision))
		}
		// Repaint error cell
		errorLabel.SetText(" ")
		if engine.haserror {
			errorLabel.SetText("E")
			errorLabel.SetStyleName("warning")

			engine.haserror = false
		}

		// Empty widget...
		for i := stack.Length() - 1; i >= 0; i-- {
			stack.Remove(i)
		}

		// ...and repaint
		for _, val := range engine.stack {
			stack.Append(tui.NewHBox(
				tui.NewLabel(strconv.FormatFloat(val, 'f', -1, 64)),
				tui.NewSpacer(),
			))
		}
		inputBox.SetText("")
	})

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Up", func() {
		if len(engine.history) == 0 {
			return
		}
		inputBox.SetText(engine.history[engine.historyindex])
		if engine.historyindex--; engine.historyindex < 0 {
			engine.historyindex = 0
		}
	})
	ui.SetKeybinding("Down", func() {
		if len(engine.history) == 0 {
			return
		}
		if engine.historyindex++; engine.historyindex > (len(engine.history) - 1) {
			engine.historyindex = len(engine.history) - 1
		}
		inputBox.SetText(engine.history[engine.historyindex])
	})

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

func analyseInput(engine *RPMEngine, input string) error {
	re := regexp.MustCompile("([0-9\\.]+)|([a-z_#]+)|([+-\\/\\*])")
	tokens := re.FindAllString(input, -1)

	// No tokens ? Just Dup !
	if len(tokens) == 0 {
		engine.Compute("dup")
		return nil
	}

	for _, tok := range tokens {
		// Skip any token separator (, or ' ')
		if tok == " " || tok == "," {
			continue
		}

		// Saving input for "repeat" feature
		// we can not do append(engine.history, tokens...) since
		// we would get unfiltered " " or "," entries
		engine.history = append(engine.history, tok)

		number, err := strconv.ParseFloat(tok, 64)

		if err != nil {
			// This is not a number
			// Call Compute to the rescue
			engine.Compute(tok)
		} else {
			engine.Push(number)
		}
	}

	engine.historyindex = len(engine.history) - 1
	return nil
}

func convertRad2Deg(engine RPMEngine, x float64) float64 {
	if engine.mode == DEG {
		return x * 180 / math.Pi
	}

	return x
}

func convertDeg2Rad(engine RPMEngine, x float64) float64 {
	if engine.mode == DEG {
		return x * math.Pi / 180
	}

	return x
}

// NewRPMEngine returns a RPMEngine with a default function catalog
func NewRPMEngine() *RPMEngine {
	engine := &RPMEngine{}
	engine.mode = RAD
	engine.haserror = false
	engine.precision = -1
	engine.catalog = map[string]interface{}{
		"+": func(x, y float64) float64 { return x + y },
		"-": func(x, y float64) float64 { return x - y },
		"*": func(x, y float64) float64 { return x * y },
		"/": func(x, y float64) float64 { return x / y },
		"^": func(x, y float64) float64 { return math.Pow(x, y) },
		"%": func(x, y float64) float64 {
			my, _ := math.Modf(y)
			mx, _ := math.Modf(x)
			return math.Mod(mx, my)
		},
		"pow":  func(x, y float64) float64 { return math.Pow(x, y) },
		"sqrt": func(x float64) float64 { return math.Sqrt(x) },
		// Trig
		"sin":  func(x float64) float64 { return math.Sin(convertDeg2Rad(*engine, x)) },
		"cos":  func(x float64) float64 { return math.Cos(convertDeg2Rad(*engine, x)) },
		"tan":  func(x float64) float64 { return math.Tan(convertDeg2Rad(*engine, x)) },
		"asin": func(x float64) float64 { return convertRad2Deg(*engine, math.Asin(x)) },
		"acos": func(x float64) float64 { return convertRad2Deg(*engine, math.Acos(x)) },
		"atan": func(x float64) float64 { return convertRad2Deg(*engine, math.Atan(x)) },
		// Precision functions
		"abs":       func(x float64) float64 { return math.Abs(x) },
		"ceil":      func(x float64) float64 { return math.Ceil(x) },
		"floor":     func(x float64) float64 { return math.Floor(x) },
		"round":     func(x float64) float64 { return math.Round(x) },
		"trunc":     func(x, y float64) float64 { return math.Round(x*math.Pow(10, y)) / math.Pow(10, y) },
		"precision": func(x float64) { engine.precision = int(x) },
		"#":         func(x float64) { engine.precision = int(x) },
		// Mode
		"rad": func() { engine.mode = RAD },
		"deg": func() { engine.mode = DEG },
		// Stack ops
		"dup": func() { engine.Dup() },
		"drop": func() {
			if len(engine.stack) > 0 {
				_ = engine.Pop()
			} else {
				engine.haserror = true
			}
		},
		"clear": func() { engine.stack = engine.stack.Clear() },
		"swap": func() {
			if len(engine.stack) < 2 {
				engine.haserror = true
				return
			}
			a, b := engine.Pop(), engine.Pop()
			engine.Push(a)
			engine.Push(b)
		},
		// Constants
		"pi":  func() { engine.Push(math.Pi) },
		"phi": func() { engine.Push(math.Phi) },
	}
	return engine
}

// PushNaked pushes a value to the internal stack without handling precision
func (e *RPMEngine) PushNaked(v float64) {
	e.stack = e.stack.Push(v)
}

// Dup duplicates last value on stack
func (e *RPMEngine) Dup() {
	if e.stack.Len() > 0 {
		n := e.Pop()
		e.PushNaked(n)
		e.Push(n)
	}
}

// Push a value to the internal stack
func (e *RPMEngine) Push(v float64) {
	p := float64(e.precision)
	if p != -1 {
		v = math.Round(v*math.Pow(10, p)) / math.Pow(10, p)
	}

	e.stack = e.stack.Push(v)
}

// Pop a value from the internal stack
func (e *RPMEngine) Pop() float64 {
	var v float64

	e.stack, v = e.stack.Pop()

	return v
}

// Compute an operation
// If the operation results in a value, push it onto the stack
func (e *RPMEngine) Compute(operation string) error {

	opFunc, ok := e.catalog[operation]

	if !ok {
		match, _ := regexp.MatchString("[a-z]{3}_[a-z]{3}", operation)

		if !match {
			e.haserror = true
			return fmt.Errorf("Operation %s not found", operation)
		}
		return e.ExchangeRate(operation)
	}

	method := reflect.ValueOf(opFunc)
	numOperands := method.Type().NumIn()
	if e.stack.Len() < numOperands {
		e.haserror = true
		// return fmt.Errorf("Too few operands for requested operation %s", operation)
	}

	operands := make([]reflect.Value, numOperands)
	for i := 0; i < numOperands; i++ {
		operands[numOperands-i-1] = reflect.ValueOf(e.Pop())
	}

	results := method.Call(operands)
	if len(results) == 1 {
		result := results[0].Float()
		// p := float64(e.precision)
		// if p != -1 {
		// 	result = math.Round(result*math.Pow(10, p)) / math.Pow(10, p)
		// }
		e.Push(result)
	}

	return nil
}

//ExchangeRate returns exchange rate for currencies specified
// as src_dst
func (e *RPMEngine) ExchangeRate(currencies string) error {
	resp, err := http.Get("http://free.currencyconverterapi.com/api/v5/convert?q=" + currencies + "&compact=y")

	if err != nil {
		e.haserror = true
		return fmt.Errorf("Currencies conversion unavailable")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// Extract floats from string
	re := regexp.MustCompile("([0-9\\.]+)")
	tokens := re.FindAllString(string(body), -1)

	// No tokens ?
	if len(tokens) == 0 {
		e.haserror = true
		return nil
	}

	res, _ := strconv.ParseFloat(tokens[0], 64)
	e.Push(res)

	return nil
}
