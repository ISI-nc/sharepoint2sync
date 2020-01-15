package internal

import (
	"os"
	"testing"
)

func Test_parseSharepointEntries(t *testing.T) {
	type args struct {
		fileName string
	}

	tests := []struct {
		name string
		args args
		expectedEntries int
	}{
		{name: "readfromfile", args: args{fileName:"sample.xml"}, expectedEntries:5},
		{name: "readfromfile2", args: args{fileName:"sample2.xml"}, expectedEntries:0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("read file ", tt.args.fileName)
			// Open our xmlFile
			xmlFile, err := os.Open(tt.args.fileName)
			// if we os.Open returns an error then handle it
			if err != nil {
				t.Error(err)
			}
			// defer the closing of our xmlFile so that we can parse it later on
			defer xmlFile.Close()

			t.Log("read entries")
			entries,err := ParseSharepointEntries(xmlFile)
			if err != nil {
				t.Error(err)
			}
			if len(entries) != tt.expectedEntries {
				t.Error("expected ", tt.expectedEntries, "got", len(entries))
			}
		})
	}
}
