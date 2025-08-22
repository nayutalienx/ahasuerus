# Ahasuerus

A 2D platformer game with time rewind mechanics built using Go and Raylib.

## Description

**English:** Ahasuerus is a unique platformer featuring time manipulation gameplay. Players can rewind time to solve puzzles, avoid obstacles, and navigate challenging levels. The game combines traditional platforming mechanics with innovative time-rewind features that affect both gameplay and audio.

**Русский:** Платформер с возможностью отмотки времени.

## Demo

https://github.com/user-attachments/assets/6db9ef27-c6f1-435f-9a38-d499aa989b19

## Features

- **Time Rewind Mechanics**: Hold Left Shift to rewind time and undo your actions
- **Dynamic Audio**: Audio plays in reverse during time rewind for immersive experience
- **Multiple Scenes**: Menu, gameplay, and level editor modes
- **Collision Detection**: Sophisticated collision system with polygon-based hitboxes
- **Particle Effects**: Visual effects and lighting system
- **Shader Support**: Custom GLSL shaders for visual enhancements
- **Configurable Graphics**: Support for multiple resolutions (720p, 1080p, 1440p)
- **Level Editor**: Built-in level creation and editing tools

## Technology Stack

- **Language**: Go 1.19+
- **Graphics**: Raylib-go (Go bindings for Raylib)
- **Audio**: Beep library for audio processing with MP3 support
- **Dependencies**: 
  - Raylib for rendering and input
  - UUID for unique identifiers
  - Particle system for visual effects

## Controls

- **Arrow Keys**: Movement (left/right)
- **Space**: Jump
- **Left Shift**: Time rewind
- **F1**: Toggle edit mode
- **Mouse**: Camera control and interaction

## Building and Running

### Prerequisites

- Go 1.19 or later
- Raylib dependencies for your platform

### Build Instructions

```bash
# Clone the repository
git clone https://github.com/nayutalienx/ahasuerus.git
cd ahasuerus

# Download dependencies
go mod tidy

# Build the game
go build -o ahasuerus

# Run the game
./ahasuerus
```

## Project Structure

```
ahasuerus/
├── audio/          # Audio processing and effects
├── collision/      # Collision detection system
├── config/         # Game configuration
├── container/      # Object management containers
├── controls/       # Input handling
├── data/           # Level data and game content
├── game/           # Core game logic
├── models/         # Game objects (player, NPCs, etc.)
├── particle/       # Particle effects system
├── repository/     # Data persistence layer
├── resources/      # Assets (textures, shaders, audio)
├── scene/          # Scene management (menu, game, editor)
└── main.go         # Application entry point
```

## Game Mechanics

### Time Rewind System

The game's core feature is its time rewind mechanism:
- Hold **Left Shift** to activate time rewind
- All game elements reverse their state
- Audio plays in reverse during rewind
- Collision detection works in both directions
- Rewind speed can vary based on gameplay context

### Player Abilities

- **Movement**: Standard left/right movement with arrow keys
- **Jumping**: Space bar for jumping with gravity physics
- **Animation System**: Multiple animation states (idle, running, jumping)
- **Collision Response**: Dynamic hitbox system for precise collision detection

### Visual Features

- **Shader Effects**: Custom fragment shaders for visual enhancement
- **Particle Systems**: Dynamic particle effects throughout the game
- **Lighting**: Dynamic lighting system with multiple light sources
- **SDF Text Rendering**: Smooth text rendering using signed distance fields

## Development

The game uses a modular architecture with clear separation of concerns:

- **Scene Management**: Different scenes for menu, gameplay, and editing
- **Resource Management**: Efficient loading and caching of assets
- **Object System**: Unified interface for all game objects
- **Event System**: Decoupled communication between game systems

## Configuration

Game settings can be modified in `config/config.go`:
- Resolution settings (720p/1080p/1440p)
- Frame rate configuration
- Graphics quality options

## License

This project is open source. See the repository for license details.

