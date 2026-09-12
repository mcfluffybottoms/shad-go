//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type stack struct {
	items []int
}

func (s *stack) Push(item int) {
	s.items = append(s.items, item)
}

func (s *stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *stack) Peek() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	}
	return s.items[len(s.items)-1], true
}

func (s *stack) Pop() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	}
	i := len(s.items) - 1
	item := s.items[i]
	s.items = s.items[:i]
	return item, true
}

type Evaluator struct {
	s          stack
	operations map[string][]string
}

func (e *Evaluator) Put(i int) {
	e.s.Push(i)
}

func (e *Evaluator) Dup() error {
	item, ok := e.s.Peek()
	if !ok {
		return errors.New("DUP - stack is empty")
	} else {
		e.s.Push(item)
		return nil
	}
}

func (e *Evaluator) Over() error {
	if len(e.s.items) < 2 {
		return errors.New("Over - stack is too small")
	} else {
		i := e.s.items[len(e.s.items)-2]
		e.Put(i)
		return nil
	}
}

func (e *Evaluator) Drop() error {
	if e.s.IsEmpty() {
		return errors.New("DROP - stack is empty")
	} else {
		e.s.Pop()
		return nil
	}
}

func (e *Evaluator) Swap() error {
	if len(e.s.items) < 2 {
		return errors.New("SWAP - stack is too small")
	} else {
		i1 := len(e.s.items) - 1
		i2 := len(e.s.items) - 2
		e.s.items[i1], e.s.items[i2] = e.s.items[i2], e.s.items[i1]
		return nil
	}
}

func (e *Evaluator) Math(op string) error {
	if len(e.s.items) < 2 {
		return errors.New("MATH - stack is too small")
	} else {
		i1 := e.s.items[len(e.s.items)-1]
		i2 := e.s.items[len(e.s.items)-2]
		e.s.Pop()
		e.s.Pop()
		var i int
		switch op {
		case "+":
			i = i2 + i1
		case "-":
			i = i2 - i1
		case "*":
			i = i2 * i1
		case "/":
			if i1 == 0 {
				return errors.New("MATH - division by zero")
			}
			i = i2 / i1
		default:
			return errors.New("MATH - operand not supported")
		}
		e.s.Push(i)
		return nil
	}
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		operations: map[string][]string{
			"drop": {"drop"},
			"swap": {"swap"},
			"dup":  {"dup"},
			"over": {"over"},
			"+":    {"+"},
			"-":    {"-"},
			"/":    {"/"},
			"*":    {"*"},
		},
	}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) wordFormation(words []string) (bool, error) {
	if words[0] == ":" && len(words) > 2 && words[len(words)-1] == ";" {
		if len(words) < 3 {
			return true, errors.New("need at least function init")
		} else {
			newOp := strings.ToLower(words[1])
			if _, er := strconv.Atoi(newOp); er == nil {
				return true, errors.New("name cant be a number")
			}
			newOps := make([]string, 0)
			for i := 2; i < len(words)-1; i++ {
				op := strings.ToLower(words[i])
				_, er := strconv.Atoi(op)
				if ops, ok := e.operations[op]; !ok && er != nil {
					return true, errors.New("error while making new command")
				} else if er == nil {
					newOps = append(newOps, op)
				} else {
					newOps = append(newOps, ops...)
				}
			}
			e.operations[newOp] = newOps
		}
		return true, nil
	}
	return false, nil
}

func (e *Evaluator) applyCustom(comm string) (bool, error) {
	if ops, ok := e.operations[comm]; ok {
		for _, op := range ops {
			err := e.processBasic(op)
			if err != nil {
				return true, err
			}
		}
		return true, nil
	}
	return false, nil
}

func (e *Evaluator) processBasic(comm string) (err error) {
	switch strings.ToLower(comm) {
	case "dup":
		err = e.Dup()
	case "over":
		err = e.Over()
	case "drop":
		err = e.Drop()
	case "swap":
		err = e.Swap()
	case "+", "-", "*", "/":
		err = e.Math(comm)
	default:
		d, er := strconv.Atoi(comm)
		if er == nil {
			e.Put(d)
		} else {
			err = errors.New("Operation not found")
		}
	}
	return err
}

func (e *Evaluator) Process(row string) ([]int, error) {
	words := strings.Split(row, " ")

	done, err := e.wordFormation(words)
	if done {
		return e.s.items, err
	}

	for _, comm := range words {
		comm = strings.ToLower(comm)
		done, err = e.applyCustom(comm)
		if done && err != nil {
			return e.s.items, err
		} else if done && err == nil {
			continue
		} else {
			err = e.processBasic(comm)
		}

		if err != nil {
			return e.s.items, err
		}
	}

	return e.s.items, err
}
