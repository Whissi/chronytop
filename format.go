package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// clean strips non-printable characters so output can't inject terminal escapes.
func clean(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, s)
}

func fmtSecs(s float64, signed bool) string {
	unit, mult := "ns", 1e9
	switch a := math.Abs(s); {
	case a >= 1:
		unit, mult = "s", 1
	case a >= 1e-3:
		unit, mult = "ms", 1e3
	case a >= 1e-6:
		unit, mult = "µs", 1e6
	}
	if signed {
		return fmt.Sprintf("%+.1f%s", s*mult, unit)
	}
	return fmt.Sprintf("%.1f%s", s*mult, unit)
}

func fmtSI(n float64, unit string) string {
	switch a := math.Abs(n); {
	case a >= 1e9:
		return fmt.Sprintf("%.1fG%s", n/1e9, unit)
	case a >= 1e6:
		return fmt.Sprintf("%.1fM%s", n/1e6, unit)
	case a >= 1e3:
		return fmt.Sprintf("%.1fk%s", n/1e3, unit)
	case a >= 10 || n == 0:
		return fmt.Sprintf("%.0f%s", n, unit)
	}
	return fmt.Sprintf("%.1f%s", n, unit)
}

func fmtBytes(n float64) string {
	for _, suffix := range []string{"B", "KiB", "MiB", "GiB"} {
		if n < 1024 {
			return fmt.Sprintf("%.1f %s", n, suffix)
		}
		n /= 1024
	}
	return fmt.Sprintf("%.1f TiB", n)
}

func fmtDur(sec float64) string {
	s := int64(sec)
	d, h, m := s/86400, s%86400/3600, s%3600/60
	s %= 60
	switch {
	case d > 0:
		return fmt.Sprintf("%dd %dh %dm", d, h, m)
	case h > 0:
		return fmt.Sprintf("%dh %dm", h, m)
	case m > 0:
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// fmtAgo formats seconds as one short unit, like chronyc's LastRx column.
func fmtAgo(sec float64) string {
	s := int64(sec)
	switch {
	case s < 60:
		return fmt.Sprintf("%ds", s)
	case s < 3600:
		return fmt.Sprintf("%dm", s/60)
	case s < 86400:
		return fmt.Sprintf("%dh", s/3600)
	}
	return fmt.Sprintf("%dd", s/86400)
}

// fmtLog2 formats a log2 interval in seconds, as chronyc reports polls and client intervals.
func fmtLog2(v int) string {
	if v == unknownInterval {
		return "-"
	}
	switch sec := math.Pow(2, float64(v)); {
	case sec < 1:
		return strconv.FormatFloat(sec, 'g', 2, 64) + "s"
	case sec < 10000:
		return strconv.Itoa(int(sec)) + "s"
	default:
		return strconv.Itoa(int(sec/3600)) + "h"
	}
}

// fmtCount formats an integer with thousands separators.
func fmtCount(n float64) string {
	s := strconv.FormatInt(int64(n), 10)
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 && s[i-1] != '-' {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	return b.String()
}
