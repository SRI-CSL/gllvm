package test

import (
	"github.com/SRI-CSL/gllvm/shared"
	"strings"
	"testing"
)

const input0 = `-pthread -c -Wno-unused-result -Wsign-compare -Wunreachable-code -DNDEBUG -g -fwrapv -O3 -Wall -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -flto -g -std=c99 -Wextra -Wno-unused-result -Wno-unused-parameter -Wno-missing-field-initializers -Wstrict-prototypes -Werror=implicit-function-declaration -fprofile-instr-use=code.profclangd -I./Include/internal  -I. -I./Include -D_FORTIFY_SOURCE=2 -D_FORTIFY_SOURCE=2 -fPIC -DPy_BUILD_CORE -DSOABI='"cpython-38-x86_64-linux-gnu"'	-o Python/dynload_shlib.o ./Python/dynload_shlib.c`

const input1 = `-pthread -c -Wno-unused-result -Wsign-compare -Wunreachable-code -DNDEBUG -g -fwrapv -O3 -Wall -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -march=x86-64 -mtune=generic -O3 -pipe -fno-plt -g -fdebug-prefix-map=/home/legend/makepkgs/python/src=/usr/src/debug -fno-semantic-interposition -flto=thin -g -std=c99 -Wextra -Wno-unused-result -Wno-unused-parameter -Wno-missing-field-initializers -Wstrict-prototypes -Werror=implicit-function-declaration -fprofile-instr-use=code.profclangd -I./Include/internal  -I. -I./Include -D_FORTIFY_SOURCE=2 -D_FORTIFY_SOURCE=2 -fPIC -DPy_BUILD_CORE -DSOABI='"cpython-38-x86_64-linux-gnu"'	-o Python/dynload_shlib.o ./Python/dynload_shlib.c`

const input2 = `-Wl,--fatal-warnings -Wl,--build-id=sha1 -fPIC -Wl,-z,noexecstack -Wl,-z,relro -Wl,-z,now -Wl,-z,defs -Wl,--as-needed -fuse-ld=lld -Wl,--icf=all -Wl,--color-diagnostics -flto=thin -Wl,--thinlto-jobs=8 -Wl,--thinlto-cache-dir=thinlto-cache -Wl,--thinlto-cache-policy,cache_size=10\%:cache_size_bytes=10g:cache_size_files=100000 -Wl,--lto-O0 -fwhole-program-vtables -Wl,--no-call-graph-profile-sort -m64 -Wl,-O2 -Wl,--gc-sections -Wl,--gdb-index -rdynamic -fsanitize=cfi-vcall -fsanitize=cfi-icall -pie -Wl,--disable-new-dtags -Wl,-O1,--sort-common,--as-needed,-z,relro,-z,now -o "./brotli" -Wl,--start-group @"./brotli.rsp"  -Wl,--end-group  -latomic -ldl -lpthread -lrt
`

const input3 = `1.c 2.c 3.c 4.c 5.c -Wl,--start-group 7.o 8.o 9.o -Wl,--end-group 10.c 11.c 12.c 13.c`

func plto(input string, t *testing.T) {
	cmds := strings.Fields(input)
	parsed := shared.Parse(cmds)
	shared.LogInfo("\n<parsed>%v</parsed>", &parsed)
	if !parsed.IsLTO {
		t.Errorf("Parsing of  %v FAILED %v (not LTO)\n", input, parsed)
	}
}

func pl(input string, t *testing.T, expected int) {
	cmds := strings.Fields(input)
	parsed := shared.Parse(cmds)
	shared.LogInfo("\n<parsed>%v</parsed>", &parsed)
	if expected != len(parsed.LinkArgs) {
		t.Errorf("Linking args %v of length %v NOT the expected length %v\n", parsed.LinkArgs, len(parsed.LinkArgs), expected)
	}
}

func Test_parsing(t *testing.T) {
	plto(input0, t)
	plto(input1, t)
	pl(input2, t, 32)
	pl(input3, t, 5)
}
