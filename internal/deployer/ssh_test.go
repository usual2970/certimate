package deployer

import (
	"os"
	"path"
	"testing"
)

func TestPath(t *testing.T) {
	dir := path.Dir("./a/b/c")
	os.MkdirAll(dir, 0o755)
}
