package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nwaples/rardecode"
)

type ConversionStats struct {
	Successful int
	Skipped    int
	Failed     int
	Total      int
}

func main() {
	fmt.Println("CBR to CBZ Converter")
	fmt.Println(strings.Repeat("=", 60))

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Find all .cbr files
	cbrFiles, err := filepath.Glob(filepath.Join(currentDir, "*.cbr"))
	if err != nil {
		fmt.Printf("Error searching for .cbr files: %v\n", err)
		return
	}

	if len(cbrFiles) == 0 {
		fmt.Println("No .cbr files found in current directory.")
		return
	}

	fmt.Printf("Found %d CBR file(s) to convert\n", len(cbrFiles))
	fmt.Println(strings.Repeat("=", 60))

	// Process each file
	stats := ConversionStats{Total: len(cbrFiles)}

	for _, cbrPath := range cbrFiles {
		result := convertCBRtoCBZ(cbrPath)
		switch result {
		case "success":
			stats.Successful++
		case "skipped":
			stats.Skipped++
		case "failed":
			stats.Failed++
		}
	}

	// Print summary
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Conversion Summary:")
	fmt.Printf("  ✓ Successful: %d\n", stats.Successful)
	fmt.Printf("  ⚠ Skipped: %d\n", stats.Skipped)
	fmt.Printf("  ✗ Failed: %d\n", stats.Failed)
	fmt.Printf("  Total: %d\n", stats.Total)

	// Ask about deletion
	if stats.Successful > 0 {
		fmt.Println()
		fmt.Println(strings.Repeat("=", 60))
		fmt.Print("Delete original .cbr files? (yes/no): ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "yes" || response == "y" {
			fmt.Println("\nDeleting original files...")
			for _, cbrPath := range cbrFiles {
				cbzPath := strings.TrimSuffix(cbrPath, ".cbr") + ".cbz"
				if _, err := os.Stat(cbzPath); err == nil {
					if err := os.Remove(cbrPath); err != nil {
						fmt.Printf("  ✗ Failed to delete %s: %v\n", filepath.Base(cbrPath), err)
					} else {
						fmt.Printf("  ✓ Deleted: %s\n", filepath.Base(cbrPath))
					}
				}
			}
			fmt.Println("\nDone!")
		} else {
			fmt.Println("\nOriginal files kept.")
		}
	}
}

func convertCBRtoCBZ(cbrPath string) string {
	cbzPath := strings.TrimSuffix(cbrPath, ".cbr") + ".cbz"
	baseName := filepath.Base(cbrPath)

	fmt.Printf("\nProcessing: %s\n", baseName)

	// Check if output file already exists
	if _, err := os.Stat(cbzPath); err == nil {
		fmt.Printf("  ⚠ Skipping: %s already exists\n", filepath.Base(cbzPath))
		return "skipped"
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "cbr2cbz-*")
	if err != nil {
		fmt.Printf("  ✗ Error creating temp directory: %v\n", err)
		return "failed"
	}
	defer os.RemoveAll(tempDir)

	// Extract RAR file
	fmt.Println("  → Extracting...")
	extractedFiles, err := extractRAR(cbrPath, tempDir)
	if err != nil {
		fmt.Printf("  ✗ Error extracting: %v\n", err)
		return "failed"
	}

	if len(extractedFiles) == 0 {
		fmt.Println("  ✗ Error: No files found in archive")
		return "failed"
	}

	fmt.Printf("  → Creating CBZ with %d files...\n", len(extractedFiles))

	// Create CBZ file
	err = createCBZ(cbzPath, tempDir, extractedFiles)
	if err != nil {
		fmt.Printf("  ✗ Error creating CBZ: %v\n", err)
		os.Remove(cbzPath)
		return "failed"
	}

	// Get file sizes
	originalInfo, _ := os.Stat(cbrPath)
	newInfo, _ := os.Stat(cbzPath)

	originalSize := float64(originalInfo.Size()) / 1024 / 1024
	newSize := float64(newInfo.Size()) / 1024 / 1024
	ratio := (1 - float64(newInfo.Size())/float64(originalInfo.Size())) * 100

	fmt.Printf("  ✓ Created: %s\n", filepath.Base(cbzPath))
	fmt.Printf("    Original: %.2f MB\n", originalSize)
	fmt.Printf("    New: %.2f MB\n", newSize)
	fmt.Printf("    Change: %+.1f%%\n", ratio)

	return "success"
}

func extractRAR(rarPath, destDir string) ([]string, error) {
	var extractedFiles []string

	// Open RAR file
	r, err := rardecode.OpenReader(rarPath, "")
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Extract each file
	for {
		header, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Skip directories and hidden files
		if header.IsDir || strings.HasPrefix(filepath.Base(header.Name), ".") {
			continue
		}

		// Create destination path
		destPath := filepath.Join(destDir, header.Name)

		// Create directories if needed
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return nil, err
		}

		// Create file
		outFile, err := os.Create(destPath)
		if err != nil {
			return nil, err
		}

		// Copy content
		_, err = io.Copy(outFile, r)
		outFile.Close()
		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, header.Name)
	}

	return extractedFiles, nil
}

func createCBZ(cbzPath, sourceDir string, files []string) error {
	// Create ZIP file with best compression
	outFile, err := os.Create(cbzPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Sort files for consistent ordering
	sort.Strings(files)

	// Add each file to ZIP
	for _, file := range files {
		// Create ZIP entry with best compression
		header := &zip.FileHeader{
			Name:   file,
			Method: zip.Deflate,
		}
		header.SetMode(0644)

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Open source file
		sourcePath := filepath.Join(sourceDir, file)
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			return err
		}

		// Copy to ZIP
		_, err = io.Copy(writer, sourceFile)
		sourceFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}