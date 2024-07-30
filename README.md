# Convoke: A game server in go
> Open for name suggestions

An game server written in go, using rethinkDB and websockets


## Features:
| Feature | Status |
| --- | --- |
| Lobbies     | 📆 |
| Matchmaking | 📆 |
| P2P Logic   | 🔨 |
| Player Data | 📆 |
| WebUI       | 🔨 |
| [Godot Addon](https://github.com/Cyteon/convoke-godot) | 🔨 |
| Basic rooms | 🔨 |

✅ - Done | 🔨 - In Progress | 📆 - Planned | ❌ - Not Planned


## Run Locally

Clone the project

```bash
  git clone https://github.com/Cyteon/convoke
```

Go to the project directory

```bash
  cd convoke
```

Rename `config.example.yaml` to `config.yaml` and populate the values

Start the server

```bash
  go run .
```


## Todo:
- [x]  Basic player authentication
- [ ]  Create a lobby
- [ ]  Join a lobby
- [ ]  Chat in a lobby
- [ ]  Game logic
