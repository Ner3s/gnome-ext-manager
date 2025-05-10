# GNOME Extension Manager

A command-line tool for backup and restoration of GNOME Shell extensions.

[PT-BR](README_pt-BR.md)

## Features

- **Backup Extensions**: Create a backup of all your enabled GNOME extensions along with their settings
- **Restore Extensions**: Restore previously backed up extensions and their configurations
- **Progress Indicators**: Visual progress bars to show operation status
- **Simple Interface**: Easy-to-use command-line interface

## Requirements

- GNOME Shell
- Go programming language
- `gnome-extensions` command-line tool
- `dconf` command-line tool

## Installation

Clone this repository and build the application:

```bash
curl -L https://github.com/ner3s/gnome-ext-manager/releases/latest/download/gnome-ext-manager -o /usr/local/bin/gnome-ext-manager
chmod +x /usr/local/bin/gnome-ext-manager
```

## Usage

Run the application:

```bash
gnome-ext-manager
```

### Backup Your Extensions

Select option 1 from the main menu and enter a filename for your backup (or press Enter for the default name).

The backup process will:
1. Detect your GNOME Shell version
2. List all enabled extensions
3. Save extension data to a JSON file
4. Export extension settings using dconf
5. Create a zip archive containing all the backup data

### Restore Your Extensions

Select option 2 from the main menu and provide the path to your backup zip file.

The restore process will:
1. Extract the zip archive
2. Install each extension from extensions.gnome.org
3. Enable all previously enabled extensions
4. Restore extension settings using dconf

## How It Works

The backup functionality:
- Uses `gnome-extensions list --enabled` to identify installed extensions
- Fetches extension metadata from extensions.gnome.org
- Uses `dconf dump` to export extension settings
- Packages everything into a zip file

The restore functionality:
- Extracts the backup archive
- Downloads extensions from extensions.gnome.org
- Installs them with `gnome-extensions install`
- Enables each extension with `gnome-extensions enable`
- Restores settings with `dconf load`

## License

This project is licensed under the terms of the LICENSE file included in this repository.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
