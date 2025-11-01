# CBR to CBZ Converter

A command-line tool to convert `.cbr` (Comic Book RAR) files to `.cbz` (Comic Book ZIP) files with maximum compression.

## Features

- **Batch Processing** - Processes a batch of files in a directory automatically
- **Maximum Compression** - Uses the max compression rate for zips automatically
- **Progress Tracking** - Shows conversion progress and file size comparisons
- **Smart Skipping** - Avoids overwriting existing CBZ files
- **Optional Cleanup** - Prompts to delete original CBR files after conversion

## Download

Download the latest release for your platform from the [Releases](https://github.com/driftywinds/cbr2cbz/releases) page. I have only personally tested the Windows .exe file, but all the other should work in theory. Supported platforms are Windows, MacOS, Linux (amd64, arm64)

## Quick Start

### Windows
1. Download `cbr2cbz.exe` from [Releases](https://github.com/driftywinds/cbr2cbz/releases) or compile it yourself
2. Place it in the folder containing your `.cbr` files
3. Double-click `cbr2cbz.exe` or run it from Command Prompt
4. Follow the prompts

### Linux/macOS

1. Download the binary from [Releases](https://github.com/driftywinds/cbr2cbz/releases) or compile it yourself
2. Make it executable:
   ```chmod +x cbr2cbz-linux  # or cbr2cbz-macos```
3. Place it in the folder with your `.cbr` files
4. Run it:
   ```./cbr2cbz-linux  # or ./cbr2cbz-macos```

## Usage

Simply place the executable in a directory containing `.cbr` files and run it:

```bash
cbr2cbz.exe  # Windows
./cbr2cbz    # Linux/macOS
```

The tool will:
1. Scan the current directory for all `.cbr` files
2. Convert each file to `.cbz` format with maximum compression
3. Display progress and file size comparisons
4. Provide a summary of successful/failed conversions
5. Ask if you want to delete the original `.cbr` files

### Example Output

```
CBR to CBZ Converter
============================================================
Found 3 CBR file(s) to convert
============================================================

Processing: Amazing Spider-Man 001.cbr
  → Extracting...
  → Creating CBZ with 24 files...
  ✓ Created: Amazing Spider-Man 001.cbz
    Original: 45.32 MB
    New: 42.18 MB
    Change: -6.9%

Processing: Batman 001.cbr
  → Extracting...
  → Creating CBZ with 28 files...
  ✓ Created: Batman 001.cbz
    Original: 52.14 MB
    New: 48.92 MB
    Change: -6.2%

============================================================
Conversion Summary:
  ✓ Successful: 2
  ⚠ Skipped: 0
  ✗ Failed: 0
  Total: 2

============================================================
Delete original .cbr files? (yes/no):
```

## Building from Source

### Prerequisites
- [Go 1.16 or later](https://go.dev/dl/)

### Build Instructions

1. Clone the repository:
   ```bash
   git clone https://github.com/driftywinds/cbr2cbz.git
   cd cbr2cbz
   ```

2. Initialize Go module and install dependencies:
   ```bash
   go mod init cbr2cbz
   go get github.com/nwaples/rardecode
   ```

3. Build for your platform:

   **Windows:**
   ```bash
   go build -ldflags="-s -w" -o cbr2cbz.exe main.go
   ```

   **Linux:**
   ```bash
   go build -ldflags="-s -w" -o cbr2cbz-linux main.go
   ```

   **macOS:**
   ```bash
   go build -ldflags="-s -w" -o cbr2cbz-macos main.go
   ```

4. **Cross-compile** for other platforms:
   ```bash
   # Windows from Linux/Mac
   GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o cbr2cbz.exe main.go
   
   # Linux from Windows/Mac
   GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o cbr2cbz-linux main.go
   
   # macOS from Windows/Linux
   GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o cbr2cbz-macos main.go
   ```

The `-ldflags="-s -w"` flags reduce the executable size by stripping debug information.

## FAQ

**Q: What's the difference between CBR and CBZ?**  
A: CBR files use RAR compression, while CBZ files use ZIP compression. CBZ files are more widely supported and don't require proprietary software to open.

**Q: Will this damage my original files?**  
A: No. The tool only deletes original files if you explicitly choose "yes" when prompted. Until then, both CBR and CBZ files coexist.

**Q: Does this work with password-protected CBR files?**  
A: No, password-protected RAR files are not supported.

**Q: Why is the CBZ file sometimes larger than the CBR?**  
A: RAR typically achieves better compression than ZIP. However, CBZ files offer better compatibility and are easier to work with.

**Q: Can I convert CBZ back to CBR?**  
A: This tool only converts CBR → CBZ. For the reverse, you'd need a different tool.

## Dependencies

The project uses the following Go library:
- [nwaples/rardecode](https://github.com/nwaples/rardecode) - RAR archive reading

---

Made with ❤️ for comic book enthusiasts
