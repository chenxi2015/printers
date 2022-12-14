// Package printers
package printers

import "testing"

func TestAddCustomPaperSize(t *testing.T) {
	type args struct {
		PrinterName string
		PaperName   string
		WidthMM     uint32
		HeightMM    uint32
		LeftMM      uint32
		TopMM       uint32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"HP-100x200mm", args{`\\192.168.227.136\HP LaserJet 1020`, "_Custom.100x200mm", 100, 200, 0, 0}, false},
		{"Canon-100x200mm", args{`Canon Generic Plus UFR II`, "_Custom.100x200mm", 100, 200, 0, 0}, false},
		{"XPS-100x200mm", args{`Microsoft XPS Document Writer`, "_Custom.100x200mm", 100, 200, 0, 0}, false},
		{"PDF-50x100mm", args{`Canon Generic Plus UFR II`, "_Custom.50x100mm", 50, 100, 0, 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddCustomPaperSize(tt.args.PrinterName, tt.args.PaperName, tt.args.WidthMM, tt.args.HeightMM, tt.args.LeftMM, tt.args.TopMM); (err != nil) != tt.wantErr {
				t.Errorf("AddCustomPaperSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteCustomPaperSize(t *testing.T) {
	type args struct {
		printerName string
		paperName   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Canon-100x200mm", args{`Canon Generic Plus UFR II`, "_Custom.100x200mm"}, false},
		{"Canon-100x200mm", args{`Canon Generic Plus UFR II`, "_Custom.50x100mm"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteCustomPaperSize(tt.args.printerName, tt.args.paperName); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCustomPaperSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
