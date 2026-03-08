#!/bin/bash
# Portfolio TUI Configuration Setup Script

set -e

echo "🎨 Portfolio TUI Configuration Setup"
echo "======================================"
echo

# Check if config.yaml already exists
if [ -f "config.yaml" ]; then
    echo "⚠️  config.yaml already exists!"
    read -p "Do you want to overwrite it? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Keeping existing config.yaml"
        exit 0
    fi
fi

# Copy example config
echo "📋 Copying config.example.yaml to config.yaml..."
cp config.example.yaml config.yaml
echo "✓ Configuration file created!"
echo

# Prompt for customization
echo "Let's personalize your portfolio!"
echo

read -p "Enter your full name [Ashik Eqbal]: " fullname
fullname=${fullname:-"Ashik Eqbal"}

read -p "Enter your username [ashikeqbal]: " username
username=${username:-"ashikeqbal"}

read -p "Enter your tagline [Portfolio]: " tagline
tagline=${tagline:-"Portfolio"}

read -p "Enter your role/bio [Developer • DevOps • Homelab Builder]: " role
role=${role:-"Developer • DevOps • Homelab Builder"}

read -p "Enter your host [localhost]: " host
host=${host:-"localhost"}

read -p "Enter your email [your.email@example.com]: " email
email=${email:-"your.email@example.com"}

read -p "Enter your phone [+1234567890]: " phone
phone=${phone:-"+1234567890"}

read -p "Enter your GitHub URL [https://github.com/yourusername]: " github
github=${github:-"https://github.com/yourusername"}

read -p "Enter your LinkedIn URL [https://linkedin.com/in/yourprofile]: " linkedin
linkedin=${linkedin:-"https://linkedin.com/in/yourprofile"}

read -p "Enter your Twitter/X URL [https://twitter.com/yourhandle]: " twitter
twitter=${twitter:-"https://twitter.com/yourhandle"}

read -p "Enter your website URL [https://yourwebsite.com]: " website
website=${website:-"https://yourwebsite.com"}

echo

# Update config.yaml with user input
echo "📝 Updating configuration..."

# Use sed to replace values (macOS/Linux compatible)
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/full_name: .*/full_name: \"$fullname\"/" config.yaml
    sed -i '' "s/username: .*/username: \"$username\"/" config.yaml
    sed -i '' "s/tagline: .*/tagline: \"$tagline\"/" config.yaml
    sed -i '' "s/role: .*/role: \"$role\"/" config.yaml
    sed -i '' "s/host: .*/host: \"$host\"/" config.yaml
    sed -i '' "s|email: .*|email: \"$email\"|" config.yaml
    sed -i '' "s|phone: .*|phone: \"$phone\"|" config.yaml
    sed -i '' "s|github: .*|github: \"$github\"|" config.yaml
    sed -i '' "s|linkedin: .*|linkedin: \"$linkedin\"|" config.yaml
    sed -i '' "s|twitter: .*|twitter: \"$twitter\"|" config.yaml
    sed -i '' "s|website: .*|website: \"$website\"|" config.yaml
else
    # Linux
    sed -i "s/full_name: .*/full_name: \"$fullname\"/" config.yaml
    sed -i "s/username: .*/username: \"$username\"/" config.yaml
    sed -i "s/tagline: .*/tagline: \"$tagline\"/" config.yaml
    sed -i "s/role: .*/role: \"$role\"/" config.yaml
    sed -i "s/host: .*/host: \"$host\"/" config.yaml
    sed -i "s|email: .*|email: \"$email\"|" config.yaml
    sed -i "s|phone: .*|phone: \"$phone\"|" config.yaml
    sed -i "s|github: .*|github: \"$github\"|" config.yaml
    sed -i "s|linkedin: .*|linkedin: \"$linkedin\"|" config.yaml
    sed -i "s|twitter: .*|twitter: \"$twitter\"|" config.yaml
    sed -i "s|website: .*|website: \"$website\"|" config.yaml
fi

echo "✅ Configuration updated successfully!"
echo
echo "Your personalized settings:"
echo "  Name:     $fullname"
echo "  Username: $username"
echo "  Tagline:  $tagline"
echo "  Role:     $role"
echo "  Host:     $host"
echo "  Email:    $email"
echo "  Phone:    $phone"
echo "  GitHub:   $github"
echo "  LinkedIn: $linkedin"
echo "  Twitter:  $twitter"
echo "  Website:  $website"
echo
echo "You can further customize config.yaml to:"
echo "  • Change the ASCII logo"
echo "  • Update app version and branding"
echo "  • Modify the sidebar title"
echo
echo "See CONFIG.md for detailed documentation."
echo
echo "🚀 Run './portfolio-tui' or 'go run .' to start!"
