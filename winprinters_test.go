// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package winprinters

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPrinter(t *testing.T) {
	name, err := GetDefault()
	if err != nil {
		t.Fatalf("GetDefault failed: %v", err)
	}
	t.Logf("GetDefault: %s", name)

	p, err := Open(name)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer func(p *Printer) {
		_ = p.Close()
	}(p)

	err = p.StartDocument("my document", "RAW")
	if err != nil {
		t.Fatalf("StartDocument failed: %v", err)
	}
	defer func(p *Printer) {
		_ = p.EndDocument()
	}(p)
	err = p.StartPage()
	if err != nil {
		t.Fatalf("StartPage failed: %v", err)
	}
	_, _ = fmt.Fprintf(p, "Hello %q\n", name)
	err = p.EndPage()
	if err != nil {
		t.Fatalf("EndPage failed: %v", err)
	}
}

func TestReadNames(t *testing.T) {
	names, err := ReadNames()
	if err != nil {
		t.Fatalf("ReadNames failed: %v", err)
	}
	jsonNames, _ := json.MarshalIndent(names, "", "  ")
	t.Logf("ReadNames: %s", jsonNames)

	name, err := GetDefault()
	if err != nil {
		t.Fatalf("GetDefault failed: %v", err)
	}
	t.Logf("GetDefault: %s", name)

	// make sure default printer is listed
	for _, v := range names {
		if v == name {
			return
		}
	}
	t.Fatalf("Default printed %q is not listed amongst printers returned by ReadNames %q", name, names)
}

func TestDriverInfo(t *testing.T) {
	name, err := GetDefault()
	if err != nil {
		t.Fatalf("GetDefault failed: %v", err)
	}
	t.Logf("GetDefault: %s", name)

	p, err := Open(name)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer func(p *Printer) {
		_ = p.Close()
	}(p)

	di, err := p.DriverInfo()
	if err != nil {
		t.Fatalf("DriverInfo failed: %v", err)
	}
	t.Logf("%+v", di)
}

func TestJobs(t *testing.T) {
	names, err := ReadNames()
	if err != nil {
		t.Fatalf("ReadNames failed: %v", err)
	}
	for _, name := range names {
		t.Log("Printer Name:", name)
		p, err := Open(name)
		if err != nil {
			t.Fatalf("Open failed: %v", err)
		}

		pj, err := p.Jobs()
		if err != nil {
			closePrinter(p)
			t.Fatalf("Jobs failed: %v", err)
		}
		if len(pj) > 0 {
			t.Log("Print Jobs:", len(pj))
			for _, j := range pj {
				b, err := json.MarshalIndent(j, "", "   ")
				if err == nil && len(b) > 0 {
					t.Log(string(b))
				}
			}
		}
		closePrinter(p)
	}
}

func closePrinter(p *Printer) {
	if p == nil {
		return
	}
	_ = p.Close()
}

func TestSetDefault(t *testing.T) {
	type args struct {
		printer string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"PDF", args{"Microsoft Print To PDF"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetDefault(tt.args.printer); (err != nil) != tt.wantErr {
				t.Errorf("SetDefault() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrinter_Forms(t *testing.T) {
	tests := []struct {
		name    string
		printer string
		wantErr bool
	}{
		{"Fax", "Fax", false},
		{"PDF", "Microsoft Print To PDF", false},
		{"XPS", "Microsoft XPS Document Writer", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Open(tt.printer)
			if err != nil {
				t.Errorf("Open got error: %v", err)
				return
			}

			var gotForms []FormInfo
			if gotForms, err = p.Forms(); (err != nil) != tt.wantErr {
				t.Errorf("Forms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, form := range gotForms {
				//formBytes, _ := json.Marshal(form)
				//t.Logf("%d. %s", i+1, formBytes)
				t.Logf("%d. %#v", i+1, form)
			}
		})
	}
}
