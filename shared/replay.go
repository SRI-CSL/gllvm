package shared

import (
	"bufio"
	"os"
)

const (
	recording_env_var = "GLLVM_REPLAY_LOG"
)

/*
 * Assumming GLLVM_REPLAY_LOG is set, this records the compiler call out to the file indicated by the
 * value of said environment variable. The format of the record is:
 *  <start of format>
 *  pwd
 *  compiler
 *  args[0]
 *  args[1]
 *  ...
 *  args[N]
 *
 * </start of format>
 * The end being signalled by a blank line. This format make replaying somewaht trivial.
 */
func Record(args []string, compiler string) {
	logfile := os.Getenv(recording_env_var)
	if len(logfile) > 0 {
		fp, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer fp.Close()
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fp.WriteString(dir)
		fp.WriteString("\n")
		fp.WriteString(compiler)
		fp.WriteString("\n")
		for _, arg := range args {
			fp.WriteString(arg)
			fp.WriteString("\n")
		}
		fp.WriteString("\n")
	}
}

type CompilerCall struct {
	Pwd  string
	Name string
	Args []string
}

func readCompilerCall(scanner *bufio.Scanner) (call *CompilerCall) {
	if scanner.Scan() {
		//got one line, probably an entire call too ...
		callp := new(CompilerCall) //pass in a local version later, save on mallocing.
		line := scanner.Text()
		if len(line) == 0 {
			panic("empty CompilerCall.Pwd")
		}
		callp.Pwd = line
		if !scanner.Scan() {
			panic("non-existant CompilerCall.Name")
		}
		line = scanner.Text()
		if len(line) == 0 {
			panic("empty CompilerCall.Name")
		}
		callp.Name = line

		for scanner.Scan() {
			line = scanner.Text()
			if len(line) == 0 {
				if len(callp.Args) == 0 {
					panic("empty CompilerCall.Args")
				}
				break
			}
			callp.Args = append(callp.Args, line)
		}
		call = callp
	}
	return
}

/*
 *
 */
func Replay(path string) (ok bool) {
	fp, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for {
		callp := readCompilerCall(scanner)
		if callp == nil {
			return
		}
		ok = replayCall(callp)
		if !ok {
			return
		}
	}
	ok = true
	return
}

/*
 *
 */
func replayCall(call *CompilerCall) bool {
	err := os.Chdir(call.Pwd)
	if err != nil {
		panic(err)
	}
	exitCode := Compile(call.Args, call.Name)
	if exitCode != 0 {
		return false
	}
	return true
}
