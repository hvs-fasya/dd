package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const (
	inFileName    = "in_file"
	outFileName   = "out_file"
	execFileName  = "go-dd"
	wrongFileName = "wrong"
	dirName       = "testdir"
)

var (
	inFileBytes = []byte("test text file\nsecond line")
)

func TestMain(m *testing.M) {
	//init test source file
	err := ioutil.WriteFile(inFileName, inFileBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			clear()
		}
	}()
	//init empty dir
	os.Mkdir(dirName, 0777)
	makeCmd := exec.Command("make")
	err = makeCmd.Run()
	if err != nil {
		log.Fatalf("failed to makeCmd binary for %s: %v", execFileName, err)
	}
	m.Run()
	clear()
	os.Exit(0)
}

func TestMainExecute(t *testing.T) {
	type Args struct {
		Infile  string
		OutFile string
	}
	tests := []struct {
		name           string
		args           []string
		wantOutContent string
		wantOutput     string
		outExist       bool
	}{
		{name: "simple", args: []string{`--from=` + inFileName, `--to=` + outFileName}, wantOutContent: string(inFileBytes), wantOutput: fmt.Sprintf("TOTAL bytes copied: %d", len(inFileBytes)), outExist: true},
		{name: "with_offset", args: []string{`--from=` + inFileName, `--to=` + outFileName, fmt.Sprintf(`--offset=%d`, len(inFileBytes)-2)}, wantOutContent: string(inFileBytes[len(inFileBytes)-2:]), wantOutput: "", outExist: true},
		{name: "with_limit", args: []string{`--from=` + inFileName, `--to=` + outFileName, `--limit=2`}, wantOutContent: string(inFileBytes[:2]), wantOutput: "TOTAL bytes copied: 2", outExist: true},
		{name: "no_in_file", args: []string{`--from=` + wrongFileName, `--to=` + outFileName}, wantOutContent: "", wantOutput: "Error: open " + wrongFileName, outExist: false},
		{name: "in_file_is_dir", args: []string{`--from=` + dirName, `--to=` + outFileName}, wantOutContent: "", wantOutput: "Error: " + ErrNotRegularFile.Error(), outExist: false},
		{name: "bad_offset", args: []string{`--from=` + inFileName, `--to=` + outFileName, fmt.Sprintf(`--offset=%d`, len(inFileBytes)+1)}, wantOutContent: "", wantOutput: "Error: " + ErrLimitOffset.Error(), outExist: false},
		{name: "bad_limit", args: []string{`--from=` + inFileName, `--to=` + outFileName, fmt.Sprintf(`--limit=%d`, len(inFileBytes)+1)}, wantOutContent: "", wantOutput: "Error: " + ErrLimitOffset.Error(), outExist: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(`./`+execFileName, tt.args...)
			output, _ := cmd.CombinedOutput()
			exists := fileExists(outFileName)
			if exists != tt.outExist {
				t.Errorf("out file exists = %v, want %v", exists, tt.outExist)
			}
			if !strings.Contains(string(output), tt.wantOutput) {
				t.Errorf("\noutput:\n %v\n, want contain:\n %v", string(output), tt.wantOutput)
			}
			outContent := getOutFileContent()
			if string(outContent) != tt.wantOutContent {
				t.Errorf("\nout_file content:\n %v\n, want out_file content:\n %v", string(outContent), tt.wantOutContent)
			}
			os.Remove(outFileName)
		})
	}
}

func clear() {
	os.Remove(inFileName)
	os.Remove(outFileName)
	os.Remove(execFileName)
	os.Remove(dirName)
}

func fileExists(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
}

func getOutFileContent() (content []byte) {
	content, _ = ioutil.ReadFile(outFileName) //no errors handling in order to leave content empty if error
	return content
}
