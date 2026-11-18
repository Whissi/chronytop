package main

import "testing"

func TestFormat(t *testing.T) {
	cases := []struct{ got, want string }{
		{fmtSecs(0.0000131, true), "+13.1µs"},
		{fmtSecs(-0.0016, true), "-1.6ms"},
		{fmtSecs(0.000000844, false), "844.0ns"},
		{fmtSI(0, ""), "0"},
		{fmtSI(2.5, ""), "2.5"},
		{fmtSI(1500, "bps"), "1.5kbps"},
		{fmtBytes(1536), "1.5 KiB"},
		{fmtDur(3725), "1h 2m"},
		{fmtAgo(1016), "16m"},
		{fmtLog2(-3), "0.12s"},
		{fmtLog2(10), "1024s"},
		{fmtLog2(unknownInterval), "-"},
		{fmtCount(887837), "887,837"},
		{fmtCount(-1234), "-1,234"},
		{clean("bad\x1b[2Jname"), "bad[2Jname"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("got %q, want %q", c.got, c.want)
		}
	}
}

func TestParse(t *testing.T) {
	tr := parseTracking(parseCSV("0A630001,10.99.0.1,2,1763434000.0,-0.000013,0.000015,0.0000113," +
		"14.635,0.001,0.385,0.0133,0.0001486,16.0,Normal\n"))
	if tr == nil || tr.Name != "10.99.0.1" || tr.Leap != "Normal" || tr.Freq != 14.635 {
		t.Errorf("tracking: %+v", tr)
	}
	if parseTracking(parseCSV("too,short\n")) != nil {
		t.Error("short tracking row should be rejected")
	}

	src := parseSources(parseCSV("^,*,10.99.0.1,1,4,377,22,0.0003108,0.0003108,0.0073\n"))
	if len(src) != 1 || src[0].State != "*" || src[0].Poll != 4 || src[0].Reach != "377" {
		t.Errorf("sources: %+v", src)
	}

	cl := parseClients(parseCSV("203.0.113.5,2100,1400,-2,127,3,0,0,127,0\n"))
	if len(cl) != 1 || cl[0].Dropped != 1400 || cl[0].Interval != -2 {
		t.Errorf("clients: %+v", cl)
	}

	ss := parseServerStats(parseCSV("887837,13218,628,0,0,0,0,0,42,43,30350\n"))
	if len(ss) != 11 || ss[1] != 13218 {
		t.Errorf("serverstats: %v", ss)
	}
}
