package duration

import "testing"

var durationTests = []struct {
	str string
	d   Duration
}{
	{"0s", 0},
	{"1ns", 1 * Nanosecond},
	{"1.1µs", 1100 * Nanosecond},
	{"2.2ms", 2200 * Microsecond},
	{"3.3s", 3300 * Millisecond},
	{"4m5s", 4*Minute + 5*Second},
	{"4m5.001s", 4*Minute + 5001*Millisecond},
	{"5h6m7.001s", 5*Hour + 6*Minute + 7001*Millisecond},
	{"2d3h4m", 2*Day + 3*Hour + 4*Minute},
	{"6w3d0h", 6*Week + 3*Day},
	{"8m0.000000001s", 8*Minute + 1*Nanosecond},
	{"15250w1d23h", 1<<63 - 1},
	{"-15250w1d23h", -1 << 63},
}

func TestDurationString(t *testing.T) {
	for _, tt := range durationTests {
		if str := tt.d.String(); str != tt.str {
			t.Errorf("Duration(%d).String() = %s, want %s", int64(tt.d), str, tt.str)
		} else {
			t.Logf("Duration(%d).String() = %s", int64(tt.d), str)
		}
		if tt.d > 0 {
			if str := (-tt.d).String(); str != "-"+tt.str {
				t.Errorf("Duration(%d).String() = %s, want %s", int64(-tt.d), str, "-"+tt.str)
			} else {
				t.Logf("Duration(%d).String() = %s", int64(-tt.d), str)
			}
		}
	}
}

var parseDurationTests = []struct {
	in   string
	ok   bool
	want Duration
}{
	// simple
	{"0", true, 0},
	{"5s", true, 5 * Second},
	{"30s", true, 30 * Second},
	{"1478s", true, 1478 * Second},
	// sign
	{"-5s", true, -5 * Second},
	{"+5s", true, 5 * Second},
	{"-0", true, 0},
	{"+0", true, 0},
	// decimal
	{"5.0s", true, 5 * Second},
	{"5.6s", true, 5*Second + 600*Millisecond},
	{"5.s", true, 5 * Second},
	{".5s", true, 500 * Millisecond},
	{"1.0s", true, 1 * Second},
	{"1.00s", true, 1 * Second},
	{"1.004s", true, 1*Second + 4*Millisecond},
	{"1.0040s", true, 1*Second + 4*Millisecond},
	{"100.00100s", true, 100*Second + 1*Millisecond},
	// different units
	{"10ns", true, 10 * Nanosecond},
	{"11us", true, 11 * Microsecond},
	{"12µs", true, 12 * Microsecond}, // U+00B5
	{"12μs", true, 12 * Microsecond}, // U+03BC
	{"13ms", true, 13 * Millisecond},
	{"14s", true, 14 * Second},
	{"15m", true, 15 * Minute},
	{"16h", true, 16 * Hour},
	{"12d", true, 12 * Day},
	{"3w", true, 3 * Week},
	// composite durations
	{"3h30m", true, 3*Hour + 30*Minute},
	{"10.5s4m", true, 4*Minute + 10*Second + 500*Millisecond},
	{"-2m3.4s", true, -(2*Minute + 3*Second + 400*Millisecond)},
	{"1h2m3s4ms5us6ns", true, 1*Hour + 2*Minute + 3*Second + 4*Millisecond + 5*Microsecond + 6*Nanosecond},
	{"39h9m14.425s", true, 39*Hour + 9*Minute + 14*Second + 425*Millisecond},
	{"2w3d12h", true, 2*Week + 3*Day + 12*Hour},
	// large value
	{"52763797000ns", true, 52763797000 * Nanosecond},
	// more than 9 digits after decimal point, see https://golang.org/issue/6617
	{"0.3333333333333333333h", true, 20 * Minute},
	// 9007199254740993 = 1<<53+1 cannot be stored precisely in a float64
	{"9007199254740993ns", true, (1<<53 + 1) * Nanosecond},
	// largest duration that can be represented by int64 in nanoseconds
	{"9223372036854775807ns", true, (1<<63 - 1) * Nanosecond},
	{"9223372036854775.807us", true, (1<<63 - 1) * Nanosecond},
	{"9223372036s854ms775us807ns", true, (1<<63 - 1) * Nanosecond},
	// large negative value
	{"-9223372036854775807ns", true, -1<<63 + 1*Nanosecond},

	// errors
	{"", false, 0},
	{"3", false, 0},
	{"-", false, 0},
	{"s", false, 0},
	{".", false, 0},
	{"-.", false, 0},
	{".s", false, 0},
	{"+.s", false, 0},
	{"3000000h", false, 0},                  // overflow
	{"9223372036854775808ns", false, 0},     // overflow
	{"9223372036854775.808us", false, 0},    // overflow
	{"9223372036854ms775us808ns", false, 0}, // overflow
	// largest negative value of type int64 in nanoseconds should fail
	// see https://go-review.googlesource.com/#/c/2461/
	{"-9223372036854775808ns", false, 0},
}

func TestParseDuration(t *testing.T) {
	for _, tc := range parseDurationTests {
		d, err := ParseDuration(tc.in)
		if tc.ok && (err != nil || d != tc.want) {
			t.Errorf("ParseDuration(%q) = %v, %v, want %v, nil", tc.in, d, err, tc.want)
		} else if !tc.ok && err == nil {
			t.Errorf("ParseDuration(%q) = _, nil, want _, non-nil", tc.in)
		} else {
			t.Logf("ParseDuration(%q) = %v", tc.in, d)
		}
	}
}
