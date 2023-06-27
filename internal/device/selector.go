// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package device

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/intel/ixl-go/internal/config"
)

type matcher interface {
	match(wq *config.WorkQueue) bool
}

type matchAll struct{}

func (m matchAll) match(wq *config.WorkQueue) bool {
	return true
}

type matchDevice struct {
	device int
}

func (m *matchDevice) match(wq *config.WorkQueue) bool {
	return wq.Device.ID == uint64(m.device)
}

type matchWQ struct {
	device int
	wq     int
}

func (m *matchWQ) match(wq *config.WorkQueue) bool {
	return wq.Device.ID == uint64(m.device) && wq.ID == m.wq
}

type matchRange struct {
	from matchWQ
	to   matchWQ
}

func (m *matchRange) match(wq *config.WorkQueue) bool {
	if wq.Device.ID > uint64(m.from.device) && wq.Device.ID < uint64(m.to.device) {
		return true
	}
	if wq.Device.ID == uint64(m.from.device) {
		if m.from.device != m.to.device {
			return wq.ID >= m.from.wq
		} else {
			return wq.ID >= m.from.wq && wq.ID <= m.to.wq
		}
	}
	if wq.Device.ID == uint64(m.to.device) {
		return wq.ID <= m.to.wq
	}
	return false
}

type notMatch struct {
	m matcher
}

func (m *notMatch) match(wq *config.WorkQueue) bool {
	return !m.m.match(wq)
}

type orMatch struct {
	x matcher
	y matcher
}

func (m *orMatch) match(wq *config.WorkQueue) bool {
	return m.x.match(wq) || m.y.match(wq)
}

type andMatch struct {
	x matcher
	y matcher
}

func (m *andMatch) match(wq *config.WorkQueue) bool {
	return m.x.match(wq) && m.y.match(wq)
}

type orMatchers []matcher

func (o orMatchers) match(wq *config.WorkQueue) bool {
	for _, m := range o {
		if m.match(wq) {
			return true
		}
	}
	return false
}

func getMatcher(selector string) (matcher, error) {
	m, err := parseSelector(selector)
	if err != nil {
		return nil, err
	}
	return orMatchers(m), nil
}

func parseSelector(selector string) (matchers []matcher, err error) {
	strs := strings.Split(selector, ",")
	temp := []string{}
	for _, v := range strs {
		result := strings.TrimSpace(v)
		if result != "" {
			temp = append(temp, result)
		}
	}

	for _, m := range temp {
		matcher, err := parseExpr([]rune(m))
		if err != nil {
			return nil, err
		}
		matchers = append(matchers, matcher)
	}
	return
}

func parseExpr(runes []rune) (matcher, error) {
	const (
		stateInit = 0
		stateWord = 1
	)
	matchers := []matcher{}
	operators := []rune{}
	word := []rune{}
	state := stateInit
	popOp := func() {
		for len(operators) != 0 {
			op := operators[len(operators)-1]
			switch op {
			case '!':
				if len(matchers) == 0 {
					return
				}
				matchers[len(matchers)-1] = &notMatch{matchers[len(matchers)-1]}
				operators = operators[:len(operators)-1]
			case '&':
				if len(matchers) < 2 {
					return
				}
				matchers[len(matchers)-2] = &andMatch{matchers[len(matchers)-1], matchers[len(matchers)-2]}
				matchers = matchers[:len(matchers)-1]
				operators = operators[:len(operators)-1]
			case '|':
				if len(matchers) < 2 {
					return
				}
				matchers[len(matchers)-2] = &orMatch{matchers[len(matchers)-1], matchers[len(matchers)-2]}
				matchers = matchers[:len(matchers)-1]
				operators = operators[:len(operators)-1]
			}
		}
	}
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch state {
		case stateInit:
			switch r {
			case '*':
				matchers = append(matchers, matchAll{})
				popOp()
				continue
			case '(':
				if i == len(runes)-1 {
					return nil, errBadExpr
				}
				j := findRightParenthesis(runes[i+1:])
				if j == -1 {
					return nil, errBadExpr
				}
				matcher, err := parseExpr(runes[i+1 : i+1+j])
				if err != nil {
					return nil, err
				}
				matchers = append(matchers, matcher)
				i = i + j + 1
				popOp()
				continue
			case '&', '!', '|':
				operators = append(operators, r)
				continue
			case ' ':
				continue
			}
			if r >= '0' && r <= '9' {
				state = stateWord
				i--
				continue
			}
			return nil, errBadExpr
		case stateWord:
			if (r >= '0' && r <= '9') || r == '.' || r == '*' || r == ' ' || r == '~' {
				word = append(word, r)
				continue
			}

			matcher, err := parseWord(string(word))
			if err != nil {
				return nil, err
			}

			matchers = append(matchers, matcher)
			word = word[:0]
			i--
			state = stateInit
			popOp()

			continue
		}
	}
	if state == stateWord {
		matcher, err := parseWord(string(word))
		if err != nil {
			return nil, err
		}
		matchers = append(matchers, matcher)
	}
	if len(matchers) == 0 {
		return matchAll{}, nil
	}

	return matchers[0], nil
}

func parseWord(word string) (matcher, error) {
	ranges := strings.Split(word, "~")
	switch len(ranges) {
	case 1:
		return parseWQMatcher(word)
	case 2:
		left, err := parseWQMatcher(ranges[0])
		if err != nil {
			return nil, err
		}
		right, err := parseWQMatcher(ranges[1])
		if err != nil {
			return nil, err
		}
		lw, ok := left.(*matchWQ)
		if !ok {
			return nil, fmt.Errorf("[%w]unknown expr format: %s", errBadExpr, word)
		}
		lr, ok := right.(*matchWQ)
		if !ok {
			return nil, fmt.Errorf("[%w]unknown expr format: %s", errBadExpr, word)
		}
		return &matchRange{
			from: *lw,
			to:   *lr,
		}, nil
	default:
		return nil, fmt.Errorf("[%w]unknown expr format: %s", errBadExpr, word)
	}
}

func parseWQMatcher(word string) (matcher, error) {
	word = strings.TrimSpace(word)
	wqname := strings.Split(word, ".")
	switch len(wqname) {
	case 1:
		deviceID, err := strconv.Atoi(word)
		if err != nil {
			return nil, fmt.Errorf("[%w]unknown expr format: %s", errBadExpr, word)
		}
		return &matchDevice{
			device: deviceID,
		}, nil
	case 2:
		deviceID, err := strconv.Atoi(wqname[0])
		if err != nil {
			return nil, fmt.Errorf("[%w]unknown expr format: %s", errBadExpr, word)
		}

		wqID, err := strconv.Atoi(wqname[1])
		if err != nil {
			if wqname[1] == "*" {
				return &matchDevice{
					device: deviceID,
				}, nil
			}
			return nil, fmt.Errorf("[%w]unknown expr format: %s", errBadExpr, word)
		}
		return &matchWQ{
			device: deviceID,
			wq:     wqID,
		}, nil
	default:
		return nil, fmt.Errorf("[%w]unknown expr format: %s", errBadExpr, word)
	}
}

var errBadExpr = errors.New("bad expr")

func findRightParenthesis(runes []rune) int {
	l := 1
	for i, r := range runes {
		switch r {
		case '(':
			l++
		case ')':
			l--
			if l == 0 {
				return i
			}
		}
	}
	return -1
}
