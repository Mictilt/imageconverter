package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func installGhostscript() {
	url := "https://github.com/ArtifexSoftware/ghostpdl-downloads/releases/download/gs10031/gs10031w64.exe"

	// Download the Ghostscript installer
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading Ghostscript:", err)
		return
	}
	defer resp.Body.Close()

	// Create the installer file
	installerPath := filepath.Join(os.TempDir(), "ghostscript_installer.exe")
	installer, err := os.Create(installerPath)
	if err != nil {
		fmt.Println("Error creating installer file:", err)
		return
	}

	// Write the body to file
	_, err = io.Copy(installer, resp.Body)
	if err != nil {
		fmt.Println("Error writing installer to file:", err)
		installer.Close()
		return
	}
	installer.Close() // Ensure the file is closed before executing

	// Run the installer (opens the GUI installer)
	cmd := exec.Command(installerPath)
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error running Ghostscript installer:", err)
		return
	}

	fmt.Println("Ghostscript installer started. Please complete the installation manually.")

	// Wait for the installer to complete before proceeding
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error while waiting for Ghostscript installer to finish:", err)
		return
	}

	// Optionally, remove the installer after installation
	err = os.Remove(installerPath)
	if err != nil {
		fmt.Println("Error deleting Ghostscript installer:", err)
	}

	ghostscriptPath := `C:\Program Files\gs\gs10.03.1\bin` // Adjust this path if needed
	cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command",
    "[Environment]::SetEnvironmentVariable('Path', " +
    "[Environment]::GetEnvironmentVariable('Path', [EnvironmentVariableTarget]::Machine) + ';' + '" + ghostscriptPath + "', " +
    "[EnvironmentVariableTarget]::Machine)")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error adding to PATH: %v\nOutput: %s\n", err, output)
		return
	}
	fmt.Printf("PATH update output: %s\n", output)
}

func installLibvips() {
	// URL of libvips zip file
	url := "https://github.com/libvips/build-win64-mxe/releases/download/v8.12.0/vips-dev-w64-web-8.12.0.zip"

	// Download the zip file
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading libvips:", err)
		return
	}
	defer resp.Body.Close()

	// Create the zip file
	zipFile, err := os.Create("libvips.zip")
	if err != nil {
		fmt.Println("Error creating zip file:", err)
		return
	}

	// Write the body to file
	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		fmt.Println("Error writing zip to file:", err)
		zipFile.Close()
		return
	}
	zipFile.Close() // Ensure the file is closed before extracting

	// Extract the zip file to C://Program Files/libvips
	extractPath := `C:\Program Files\libvips`
	err = unzip("libvips.zip", extractPath)
	if err != nil {
		fmt.Println("Error extracting zip file:", err)
		return
	}

	// Delete the zip file
	err = os.Remove("libvips.zip")
	if err != nil {
		fmt.Println("Error deleting zip file:", err)
		return
	}

	// Add libvips to PATH
	vipsBinPath := filepath.Join(extractPath, "vips-dev-8.12", "bin")
	println(vipsBinPath)

	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command",
    "[Environment]::SetEnvironmentVariable('Path', " +
    "[Environment]::GetEnvironmentVariable('Path', [EnvironmentVariableTarget]::Machine) + ';' + '" + vipsBinPath + "', " +
    "[EnvironmentVariableTarget]::Machine)")
	println(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error adding to PATH: %v\nOutput: %s\n", err, output)
		return
	}
	fmt.Printf("PATH update output: %s\n", output)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Create directories if needed
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		// Ensure directory exists
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Create file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close() // Close any opened file before returning
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close() // Ensure the file is closed after writing
		rc.Close()      // Ensure the file is closed after reading

		if err != nil {
			return err
		}
	}
	return nil
}
