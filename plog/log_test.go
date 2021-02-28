package plog

import (
	"fmt"
	"os"
	"testing"

	"github.com/CharlesPu/flamingo/perror"
)

func TestLog(t *testing.T) {
	Logf(LevelInfo, "hahah")
	Infof("haha")
	Infof("%+v", perror.New("ahaha").OriginError())
	fmt.Fprintf(os.Stderr, "hahah")
	fmt.Fprintf(os.Stderr, "1111")
}
