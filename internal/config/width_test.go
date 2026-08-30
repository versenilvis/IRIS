package config

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestWidthAcceptsBothForms(t *testing.T) {
	cases := []struct {
		toml string
		want Width
	}{
		{`max-width = 80`, Width{n: 80}},
		{`max-width = 0`, Width{}},
		{`max-width = "80%"`, Width{n: 80, percent: true}},
		{`max-width = "100%"`, Width{n: 100, percent: true}},
		{`max-width = " 60 % "`, Width{n: 60, percent: true}},
		{`max-width = "80"`, Width{n: 80}},
	}
	for _, c := range cases {
		var got struct {
			MaxWidth Width `toml:"max-width"`
		}
		if _, err := toml.Decode(c.toml, &got); err != nil {
			t.Errorf("decode %q: %v", c.toml, err)
			continue
		}
		if got.MaxWidth != c.want {
			t.Errorf("decode %q = %+v; want %+v", c.toml, got.MaxWidth, c.want)
		}
	}
}

func TestWidthRejectsNonsense(t *testing.T) {
	for _, in := range []string{`max-width = "wide"`, `max-width = "%"`, `max-width = "8 0%"`, `max-width = true`} {
		var got struct {
			MaxWidth Width `toml:"max-width"`
		}
		if _, err := toml.Decode(in, &got); err == nil {
			t.Errorf("decode %q: want an error, got %+v", in, got.MaxWidth)
		}
	}
}

func TestWidthResolve(t *testing.T) {
	cases := []struct {
		w    Width
		term int
		want int
	}{
		{Width{n: 80}, 200, 80},              // absolute ignores the terminal
		{Width{n: 80, percent: true}, 200, 160},
		{Width{n: 50, percent: true}, 81, 40}, // truncates rather than rounds up
		{Width{n: 80, percent: true}, 0, 0},   // unknown terminal falls through
		{Width{}, 200, 0},                     // unset falls through
	}
	for _, c := range cases {
		if got := c.w.Resolve(c.term); got != c.want {
			t.Errorf("Width%+v.Resolve(%d) = %d; want %d", c.w, c.term, got, c.want)
		}
	}
}

func TestWidthValidate(t *testing.T) {
	ok := []Width{{}, {n: 80}, {n: 1, percent: true}, {n: 100, percent: true}}
	for _, w := range ok {
		if err := w.validate(); err != nil {
			t.Errorf("Width%+v.validate() = %v; want nil", w, err)
		}
	}
	bad := []Width{{n: -1}, {n: 0, percent: true}, {n: 101, percent: true}}
	for _, w := range bad {
		if err := w.validate(); err == nil {
			t.Errorf("Width%+v.validate() = nil; want an error", w)
		}
	}
}

// iris config show round-trips the config through the encoder, so both forms
// have to come back out the way they went in.
func TestWidthRoundTrips(t *testing.T) {
	for _, in := range []string{`max-width = 80`, `max-width = "80%"`} {
		var cfg struct {
			MaxWidth Width `toml:"max-width"`
		}
		if _, err := toml.Decode(in, &cfg); err != nil {
			t.Fatal(err)
		}
		var out []byte
		b, err := cfg.MaxWidth.MarshalTOML()
		if err != nil {
			t.Fatal(err)
		}
		out = append([]byte("max-width = "), b...)
		if string(out) != in {
			t.Errorf("round trip of %q = %q", in, out)
		}
	}
}
