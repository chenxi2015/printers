package winprinters

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
)

// AddCustomPaperSize 添加自定义纸张规格
//
// # Reference C# code
//
// - https://www.cnblogs.com/datacool/p/datacool_windowsapi_printerhelper.html
// - https://github.com/vanloc0301/ecouponsprinter/blob/master/ECouponsPrinter/ECouponsPrinter/Printer.cs
//
//goland:noinspection GoSnakeCaseUsage
func AddCustomPaperSize(printerName, paperName string, widthMM, heightMM, leftMM, topMM uint32) (err error) {
	const (
		PRINTER_ACCESS_USE        uint32 = 0x00000008
		PRINTER_ACCESS_ADMINISTER uint32 = 0x00000004
	)
	defaults := &PrinterDefaults{
		pDevMode:      new(DevMode),
		DesiredAccess: PRINTER_ACCESS_ADMINISTER | PRINTER_ACCESS_USE,
	}

	var p *Printer
	if p, err = OpenWithDefaults(printerName, defaults); err != nil {
		return
	}
	defer func() {
		_ = p.Close()
	}()
	pFormName, _ := windows.UTF16FromString(paperName)
	paper := &pFormName[0]
	_ = DeleteForm(p.h, paper) // 删除已存在的同名自定义纸张大小

	pageSize := SIZE{
		Width:  widthMM * 1000,
		Height: heightMM * 1000,
	}
	formInfo := FORM_INFO_1{
		Flags: 0,
		pName: paper,
		Size:  pageSize,
		ImageableArea: Rect{
			Left:   leftMM * 1000,
			Right:  pageSize.Width,
			Top:    topMM * 1000,
			Bottom: pageSize.Height,
		},
	}
	if err = AddForm(p.h, 1, &formInfo); err != nil {
		err = fmt.Errorf("向打印机 [%s] 添加自定义纸张大小 [%s] 失败！错误：%s", printerName, paperName, err.Error())
	}

	pDeviceName, _ := windows.UTF16PtrFromString(printerName)
	printerInfo := &PRINTER_INFO_9{}
	err = DocumentProperties(0, p.h, pDeviceName, printerInfo.pDevMode, printerInfo.pDevMode, DM_MODIFY|DM_COPY)
	if err != nil {
		return // 无法为打印机设定打印方向
	}

	if printerInfo, err = p.GetPrinter9(); err != nil {
		return // 调用 GetPrinter 方法失败，无法获取 PRINTER_INFO_9 结构
	}
	var devMode *DevMode
	if devMode, err = p.DocumentPropertiesGet(printerName); err != nil {
		return // 无法获取 DevMode 结构
	}
	devMode.dmFields = DM_FORMNAME
	devMode.dmFormName = *paper
	printerInfo.pDevMode = devMode

	// windows.ERROR_INVALID_PARAMETER
	if err = p.SetPrinter9(printerInfo); err != nil {
		return // 调用 SetPrinter 方法失败，无法进行打印机设置
	}
	return
}

// DeleteCustomPaperSize 删除自定义纸张规格
func DeleteCustomPaperSize(printerName, paperName string) (err error) {
	var p *Printer
	if p, err = Open(printerName); err != nil {
		return
	}
	defer func() {
		_ = p.Close()
	}()

	pFormName, _ := windows.UTF16FromString(paperName)
	pName := &(pFormName)[0]
	if err = DeleteForm(p.h, pName); errors.Is(err, windows.ERROR_INVALID_FORM_NAME) {
		err = nil
	}
	return
}
