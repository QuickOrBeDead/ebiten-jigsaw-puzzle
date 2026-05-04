# Ebiten Jigsaw Puzzle

A jigsaw puzzle game built with Go and the [Ebitengine](https://github.com/hajimehoshi/ebiten) game engine.

## Features

- Select from pre-loaded images or upload your own JPEG images
- Drag-and-drop puzzle pieces with snapping mechanics
- Ghost image overlay for reference
- Original image preview toggle
- Progress tracking with percentage completion
- Move counter and elapsed time tracking
- Puzzle completion celebration with stats
- Arrange button to organize pieces

## Requirements

- Go 1.26.2 or later

## Installation

```bash
git clone https://github.com/QuickOrBeDead/ebiten-jigsaw-puzzle.git
cd ebiten-jigsaw-puzzle
```

## Usage

### Running the Game

```bash
go run .
```

### Adding Puzzle Images

Place JPEG images (`.jpg` or `.jpeg`) in the `pictures/` directory. They will automatically appear in the home scene when you start the game.

### Uploading Images

Click the "Upload Image" button on the home screen to select an image from your desktop.

## Controls

- **Left Click**: Select and drag puzzle pieces
- **Release**: Drop pieces (they will snap if close enough to matching pieces)

## UI Buttons

### Home Scene
- **Upload Image**: Open file dialog to load a custom image

### Game Scene
- **Restart**: Start a new puzzle with the same image
- **Home**: Return to the home screen
- **Image**: Toggle reference image preview
- **Ghost**: Toggle ghost image overlay
- **Arrange**: Organize puzzle pieces

## Project Structure

```
├── main.go                 # Game entry point
├── common/                 # Shared utilities (scenes, buttons, drawing, math)
├── edge/                   # Puzzle piece edge definitions
├── piece/                  # Puzzle piece geometry and grouping
├── puzzle/                 # Puzzle logic and picture handling
├── scenes/                 # Game scenes
│   ├── homeScene/         # Home/menu scene
│   └── gameScene/         # Main puzzle game scene
└── pictures/              # Pre-loaded puzzle images
```

## License

See [LICENSE](LICENSE) for details.
