//
// ELF Section related code

package shared

import (
	"debug/elf"
	"debug/macho"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// type Sectionator struct {
// 	Filename string
// }

// func NewSectionator(filename string) *Sectionator {
// 	this := &Sectionator{
// 		Filename: filename,
// 	}
// 	return this
// }

// data2file takes byte contents and writes them to a file
func data2file(data []byte) (filepath *os.File, err error) {

	tmpFile, err := ioutil.TempFile("", "data2file")
	if err != nil {
		return nil, fmt.Errorf("data2file(): Unable to create temp file")
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(data)
	if err != nil {
		return nil, fmt.Errorf("Unable to write data to %s", tmpFile.Name())
	}
	return tmpFile, nil
}

// SectionWrite writes data to sectionName of an elf or mach file
func SectionWrite(filename string, data []byte, segmentName string, sectionName string) (err error) {
	// We can only attach a bitcode path to certain file types
	extension := filepath.Ext(filename)
	switch extension {
	case
		".o",
		".lo",
		".os",
		".So",
		".po":
	default:
		return fmt.Errorf("Extension %s not supported", extension)
	}

	// Sanity checks
	if len(segmentName) > 0 && runtime.GOOS != osDARWIN {
		return fmt.Errorf("Segment name requires Mac OS")
	}

	if len(segmentName) == 0 && runtime.GOOS != osLINUX {
		return fmt.Errorf("Empty segment is only compatible with Liux")
	}

	tmpFile, err := data2file(data)
	if err != nil {
		return fmt.Errorf("Unable to create temp file")
	}
	defer os.Remove(tmpFile.Name())

	var attachCmd string
	var attachCmdArgs []string
	if runtime.GOOS == osDARWIN {
		if len(LLVMLd) > 0 {
			attachCmd = LLVMLd
		} else {
			attachCmd = "ld"
		}
		attachCmdArgs = []string{"-r", "-keep_private_externs", filename, "-sectcreate", segmentName, sectionName, tmpFile.Name()}
		attachCmdArgs = append(attachCmdArgs, "-o", filename)

	} else {
		if len(LLVMObjcopy) > 0 {
			attachCmd = LLVMObjcopy
		} else {
			attachCmd = "objcopy"
		}
		attachCmdArgs = []string{"--add-section", sectionName + "=" + tmpFile.Name()}
		attachCmdArgs = append(attachCmdArgs, filename)
	}

	// Run the attach command and ignore errors
	_, nerr := execCmd(attachCmd, attachCmdArgs, "")
	if nerr != nil {
		return fmt.Errorf("SectionWrite: %v %v failed because %v", attachCmd, attachCmdArgs, nerr)
	}

	return nil
}

// SectionRead reads elf or mach sectionName from filename
func SectionRead(filename string, sectionName string) (data []byte, err error) {
	sectionName = platformizeSectionName(sectionName)

	switch platform := runtime.GOOS; platform {
	case osFREEBSD, osLINUX:
		objFile, err := elf.Open(filename)
		if err != nil {
			return nil, err
		}
		section := objFile.Section(sectionName)
		if section == nil {
			return nil, err
		}
		data, err = section.Data()
		if err != nil {
			return nil, err
		}
	case osDARWIN:
		objFile, err := macho.Open(filename)
		if err != nil {
			return nil, err
		}
		section := objFile.Section(sectionName)
		if section != nil {
			return nil, err
		}
		data, err = section.Data()
		if err != nil {
			return nil, err
		}
	}

	err = nil
	return data, err
}

func platformizeSectionName(name string) string {
	switch runtime.GOOS {
	case osDARWIN:
		return "__" + name
	case osLINUX, osFREEBSD:
		return "." + name
	default:
		return "<ERROR>"
	}
}
