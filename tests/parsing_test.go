package test

import (
	"github.com/SRI-CSL/gllvm/shared"
	"strings"
	"testing"
)

const input0 = `-pthread -c -Wno-unused-result -Wsign-compare -Wunreachable-code -DNDEBUG -g -fwrapv -O3 -Wall -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -flto -g -std=c99 -Wextra -Wno-unused-result -Wno-unused-parameter -Wno-missing-field-initializers -Wstrict-prototypes -Werror=implicit-function-declaration -fprofile-instr-use=code.profclangd -I./Include/internal  -I. -I./Include -D_FORTIFY_SOURCE=2 -D_FORTIFY_SOURCE=2 -fPIC -DPy_BUILD_CORE -DSOABI='"cpython-38-x86_64-linux-gnu"'	-o Python/dynload_shlib.o ./Python/dynload_shlib.c`

const input1 = `-pthread -c -Wno-unused-result -Wsign-compare -Wunreachable-code -DNDEBUG -g -fwrapv -O3 -Wall -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -flto=thin -g -std=c99 -Wextra -Wno-unused-result -Wno-unused-parameter -Wno-missing-field-initializers -Wstrict-prototypes -Werror=implicit-function-declaration -fprofile-instr-use=code.profclangd -I./Include/internal  -I. -I./Include -D_FORTIFY_SOURCE=2 -D_FORTIFY_SOURCE=2 -fPIC -DPy_BUILD_CORE -DSOABI='"cpython-38-x86_64-linux-gnu"'	-o Python/dynload_shlib.o ./Python/dynload_shlib.c`

func pt(input string, t *testing.T) {
	cmds := strings.Fields(input)
	parsed := shared.Parse(cmds)
	if !parsed.IsLTO {
		t.Errorf("Parsing of  %v FAILED %v (not LTO)\n", input, parsed)	
	}
}

func Test_parsing(t *testing.T) {
	pt(input0, t)
	pt(input1, t)
}
