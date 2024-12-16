# Vilo: A Lightweight Version Control CLI Tool

`vilo` is a simple version control system written in Go, designed to handle file versioning, encryption, and compression using AES and Gzip.

---

## Features
1. **Initialize**: Create the required `.vilo` directory structure.
2. **Add**: Stage files to the staging area.
3. **Commit**: Encrypt, compress, and save staged files into the `.vilo/objects` directory.
4. **Secure Storage**: Files are stored securely with AES encryption.
5. **Custom Commit Messages**: Tag your commits with meaningful messages.

---

## Installation
### Prerequisites
- Go (Golang) version 1.16 or higher

### Clone the Repository
```bash
$ git clone https://github.com/ammargit93/vilo.git
$ cd vilo
```

### Build the CLI
Run the following command to build the executable:
```bash
$ go build -o vilo.exe main.go utils.go commands.go
```

This will create the `vilo` executable in the current directory.

---

## Usage
### Initialize Vilo
Initialize a new repository in the current directory:
```bash
$ ./vilo init
```
This creates a `.vilo` directory containing:
- `.vilo/HEAD`: Stores commit hashes
- `.vilo/stage.json`: Tracks staged files
- `.vilo/objects/`: Stores committed files

### Add Files to Staging Area
Add files to the staging area using the `add` command:
```bash
$ ./vilo add --files file1.txt,file2.txt
```
**Output**:
```
/path/to/file1.txt Staged for commit
/path/to/file2.txt Staged for commit
```

### Commit Files
Commit files with a custom message:
```bash
$ ./vilo commit --Message "Initial commit"
```
**Output**:
```
Commit successful!
```

Committed files are encrypted, compressed, and saved into `.vilo/objects/<hash>/`.

### View Staging Area
You can check the staged files by viewing `.vilo/stage.json`.

---

## Internals
### Encryption & Compression
Files are:
1. **Compressed** using Gzip to save space.
2. **Encrypted** using AES (CFB mode) for secure storage.

### Directory Structure
After a commit, the structure looks like this:
```
.vilo/
├── HEAD
├── objects/
│   └── <commit-hash>/
│       ├── file1.txt.enc
│       └── file2.txt.enc
└── stage.json
```

---

## Future Improvements
- Add `log` command to view commit history.
- Implement `checkout` to restore previous versions.
- Add diff functionality to compare file changes.

---

## Contributing
Feel free to fork this repository and submit pull requests!

---

## License
This project is licensed under the MIT License. See `LICENSE` for details.

---

## Author
**Ammar Ansari**

For any questions or suggestions.

