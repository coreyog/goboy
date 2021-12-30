package goboy

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/wacul/ptr"
)

func TestOpCodeCompleteness(t *testing.T) {
	for k, v := range unprefixed {
		if v.Operation == nil {
			t.Fail()
			t.Logf("%s: %s", formatInst(k, nil), opNames[v.Code])
		}
	}

	for k, v := range cb {
		if v.Operation == nil {
			t.Fail()
			t.Logf("%s: %s", formatInst(k, ptr.String("cb")), opNames[v.Code])
		}
	}
}

func formatInst(inst byte, prefix *string) (output string) {
	i := strconv.FormatInt(int64(inst), 2)       // to binary
	i = strings.Repeat("0", 8-len(i)) + i        // pad to 8 bits
	output = i[0:2] + "_" + i[2:5] + "_" + i[5:] // split into X/Y/Z parts

	if prefix != nil {
		output = fmt.Sprintf("%s_%s", *prefix, output)
	}

	return output
}
