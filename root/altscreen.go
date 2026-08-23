package root

import "bytes"

var altScreenTransitions = []struct {
	seq   []byte
	enter bool
}{
	{[]byte("\x1b[?1049h"), true},
	{[]byte("\x1b[?1049l"), false},
	{[]byte("\x1b[?1047h"), true},
	{[]byte("\x1b[?1047l"), false},
	{[]byte("\x1b[?47h"), true},
	{[]byte("\x1b[?47l"), false},
}

// altScreenCarryLen is one byte short of the longest sequence above, so a
// switch split across two PTY reads is still seen whole.
const altScreenCarryLen = 7

// lastAltScreenTransition reports the final alternate screen switch in chunk.
//
// The last one wins rather than the first: fish probes the terminal at startup
// by entering and leaving the alternate screen inside a single write, and
// treating that as "a full screen app took over" suppressed the overlay for the
// rest of the session.
func lastAltScreenTransition(chunk []byte) (enter bool, found bool) {
	last := -1
	for _, t := range altScreenTransitions {
		if i := bytes.LastIndex(chunk, t.seq); i > last {
			last, enter, found = i, t.enter, true
		}
	}
	return enter, found
}

// scanAltScreen reports the final alternate screen switch in chunk, also
// looking at the boundary with the previous read so a sequence split across two
// reads is not missed.
func scanAltScreen(carry, chunk []byte) (enter bool, found bool) {
	if len(carry) > 0 {
		head := chunk
		if len(head) > altScreenCarryLen {
			head = head[:altScreenCarryLen]
		}
		joined := make([]byte, 0, len(carry)+len(head))
		joined = append(append(joined, carry...), head...)
		if e, ok := lastAltScreenTransition(joined); ok {
			enter, found = e, true
		}
	}
	// anything wholly inside chunk comes after the boundary, so it wins
	if e, ok := lastAltScreenTransition(chunk); ok {
		enter, found = e, true
	}
	return enter, found
}

func keepAltScreenCarry(chunk []byte) []byte {
	if len(chunk) > altScreenCarryLen {
		chunk = chunk[len(chunk)-altScreenCarryLen:]
	}
	return append([]byte(nil), chunk...)
}
