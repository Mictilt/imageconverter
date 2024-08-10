package main

import (
	"fmt"
	"os/exec"
)

func processFile(inputPath string, outputPath string) error {
	// Ensure Ghostscript is installed

	var gsPath string

	_, err := exec.LookPath("gs")
	if err != nil {
		_, err = exec.LookPath("gswin64c")
		if err != nil {
			return fmt.Errorf("ghostscript is not installed or not in PATH: %v", err)
		} else {
			gsPath = "gswin64c"
		}
	} else {
		gsPath = "gs"
	}

	// Prepare Ghostscript command
	cmd := exec.Command(gsPath,
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS=/default",
		"-dQUIET",
		"-dDetectDuplicateImages",
		"-dCompressFonts=true",
		"-dSubsetFonts=true",
		"-dCompressPages=true",
		"-dEmbedAllFonts=true",
		"-dMaxInlineImageSize=4000",
		"-dDownsampleColorImages=true",
		"-dColorImageResolution=150",
		"-dDownsampleGrayImages=true",
		"-dGrayImageResolution=150",
		"-dDownsampleMonoImages=true",
		"-dMonoImageResolution=150",
		"-o", outputPath,
		inputPath)

	// Run Ghostscript command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ghostscript error: %v\nOutput: %s", err, output)
	}

	fmt.Println("PDF processed successfully.")
	return nil
}
