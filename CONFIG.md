# Portfolio TUI - Configuration Guide

This application supports personalized configuration through a YAML file.

## Setup

1. Copy the example configuration:
   ```bash
   cp config.example.yaml config.yaml
   ```

2. Edit `config.yaml` with your personal information

3. Run the application - it will automatically load your config

## Configuration Priority

The app loads configuration in this order:
1. Config file (`config.yaml`, `config.yml`, `.portfolio-tui.yaml`, `.portfolio-tui.yml`)
2. Built-in defaults

## Customization Options

### App Information
- `app.name`: Application name displayed in neofetch
- `app.version`: Version number
- `app.runtime`: Runtime/framework name

### Owner Information
- `owner.full_name`: Your full name (e.g., "John Doe")
- `owner.username`: Username for display (e.g., "johndoe")
- `owner.tagline`: Short tagline shown in headers (e.g., "Portfolio")
- `owner.role`: Role/bio shown on home screen (e.g., "Developer • Designer • Creator")

### Network
- `network.host`: Host identifier (defaults to "localhost")

### Contact
- `contact.email`: Contact email displayed in contact fallback
- `contact.phone`: Contact phone displayed in contact fallback

### Social
- `social.github`: GitHub profile URL
- `social.linkedin`: LinkedIn profile URL
- `social.twitter`: X/Twitter profile URL
- `social.website`: Personal website URL
- `social.youtube`: YouTube channel URL
- `social.instagram`: Instagram profile URL

### Branding
- `branding.sidebar_title`: Title shown in sidebar (e.g., "PORTFOLIO")
- `branding.ascii_logo`: Main ASCII logo (newline-separated)
- `branding.sidebar_logo`: Compact sidebar logo (newline-separated)

### SSH
- `ssh.host`: Interface the app listens on inside the process or container
- `ssh.port`: SSH port used when running `portfolio-tui serve`
- `ssh.host_key_path`: Path to the SSH host key file

## ASCII Logo Format

Logos should be defined as strings with `\n` for line breaks and `\\` for backslashes:

```yaml
ascii_logo: "   /\\\n  /  \\\n /_  _\\"
```

## Example Configuration

See `config.example.yaml` for a complete example with all available options.

## Environment Variables (.env)

`.env` is used for Sanity connection values and SSH runtime overrides:
- `SANITY_PROJECT_ID`
- `SANITY_DATASET`
- `SANITY_API_VERSION`
- `PORTFOLIO_SSH_HOST`
- `PORTFOLIO_SSH_PORT`
- `PORTFOLIO_SSH_KEY`
- `PORTFOLIO_BIND_IP`

`config.yaml` controls the app's profile and can also define SSH defaults under the `ssh:` section. Environment variables override those defaults at runtime, which is especially useful for Docker and Docker Compose. If the SSH host key file does not exist at `PORTFOLIO_SSH_KEY`, the server generates it automatically on startup.
