# tnote - Simple Terminal Note Taking App

`tnote` is a lightweight terminal-based note-taking application. It allows you to quickly create, view, and delete notes directly from the command line.

## Contributing

This project is made for personal use, but I'm open to feature requests and suggestions. If you have an idea or improvement, feel free to open an issue or submit a pull request.

## Features

- Create new notes with a simple command.
- View all saved notes.
- Delete notes by specifying their ID.
- Automatically organizes notes in a hidden `.tnotes` directory.
- Compatible with Windows and probablyyy linux (no testing done there).

## Usage

### Create a New Note

To create a new note, simply type:

```bash
./tnote <text>
```

Example:

```bash
./tnote "Remember to fix the bug"
```

### Create a Note with Terminal Context (Not Yet Implemented)

To create a note with terminal context (e.g., last few commands):

```bash
./tnote -n<lines> <text>
```

Example:

```bash
./tnote -n5 "Capture last 5 commands"
```

> **Note:** This feature is not yet implemented.

### View All Notes

To view all saved notes:

```bash
./tnote -v
```

### Delete a Note

To delete a note by its ID:

```bash
./tnote -d <id>
```

Example:

```bash
./tnote -d 2
```

### Show Help

To display the help message:

```bash
./tnote -help
```

## Notes Directory

All notes are saved in a hidden directory named `.tnotes` in the current working directory. this directory is automatically hidden.

## Examples

1. Create a note:

   ```bash
   ./tnote "Buy groceries"
   ```

2. View all notes:

   ```bash
   ./tnote -v
   ```

3. Delete a note:

   ```bash
   ./tnote -d 1
   ```

## Limitations

- Terminal context capture (`-n`) is not yet implemented.
- Notes are stored as plain text files in the `.tnotes` directory.

## License

This project is open-source and available under the MIT License.
