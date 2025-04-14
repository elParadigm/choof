
# Choof

![Choof Demo](choof.gif)

Choof is a lightweight and intuitive CLI tool for viewing and managing file information efficiently. It provides a visually appealing and structured way to inspect file details, delete files, rename them, and open them directly from the terminal.

## Features
- **Detailed File Info:** Displays metadata such as size, creation/modification dates, permissions, and hashes (SHA256, MD5).
- **File Operations:** Delete, rename, and open files with simple keybindings.
- **A permission modification tool**
- **User-Friendly Interface:** Clean, structured output with colors and icons.
- **Keyboard Shortcuts:**
  - `Ctrl + c / q` - Close the program
  - `d` - Delete file (with confirmation)
  - `Ctrl + r` - Rename file
  - `Ctrl + o` - Open file
  - `Ctrl + p` - Permission changer
  - `p` - Toggle permission explainer
  - `p` - Toggle size




## Installation
```sh
# Clone the repository
git clone https://github.com/elParadigm/choof.git
cd choof

# Build and run
go build -o choof
./choof <filename>
```

## Usage
```sh
choof /path/to/file
```
Navigate the interface using the provided keybindings to manage your files efficiently.

## Contributing
Pull requests are welcome! If you have suggestions or find issues, feel free to open an issue on the repository.

