# VHS demo tape for drim - One-click deployment tool for Kaneo
# https://github.com/charmbracelet/vhs

Output demo.gif

# Settings
Set FontSize 28
Set Width 2000
Set Height 1200
Set Padding 20
Set Margin 20
Set MarginFill "#1a1b26"
Set BorderRadius 8
Set WindowBar Colorful
Set WindowBarSize 60
Set Theme "TokyoNight"
Set TypingSpeed 50ms
Set PlaybackSpeed 1.0

# Requirements
Require drim

# Start recording
Hide
Type "clear"
Enter
Show

Sleep 500ms

# Show the main setup command
Type "drim setup"
Enter
Sleep 1s

# Answer the reverse proxy prompt
Type "y"
Enter
Sleep 1s

# Enter domain
Type "kaneo.example.com"
Enter
Sleep 3s

# Wait for setup to complete - check the whole screen for the completion message
Sleep 12s
