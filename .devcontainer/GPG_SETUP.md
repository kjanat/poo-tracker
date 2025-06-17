# GPG Commit Signing in DevContainer

This guide explains how to set up GPG commit signing in the DevContainer environment.

## 🔐 How GPG Integration Works

The DevContainer automatically mounts your host GPG configuration and keys:

- **Host GPG Directory**: `~/.gnupg` (Linux/macOS) or `%USERPROFILE%\.gnupg` (Windows)
- **Container Mount**: `/home/kjanat/.gnupg`
- **Git Config**: Your host `.gitconfig` is also mounted

## 📋 Prerequisites

### On Your Host System

1. **Install GPG**:

   ```bash
   # Ubuntu/Debian
   sudo apt install gnupg
   
   # macOS
   brew install gnupg
   
   # Windows (use Git Bash or WSL)
   # GPG is included with Git for Windows
   ```

2. **Generate a GPG Key** (if you don't have one):

   ```bash
   gpg --full-generate-key
   ```

   Choose:

   - Key type: RSA and RSA (default)
   - Key size: 4096 bits
   - Expiration: 0 (doesn't expire) or set a date
   - Enter your name and email (use your GitHub email)

3. **Configure Git for Signing**:

   ```bash
   # List your GPG keys to get the key ID
   gpg --list-secret-keys --keyid-format LONG
   
   # Copy the key ID (after "sec rsa4096/")
   # Example: if you see "sec rsa4096/ABC123DEF456", the key ID is "ABC123DEF456"
   
   # Configure git with your key
   git config --global user.signingkey YOUR_KEY_ID
   git config --global commit.gpgsign true
   git config --global tag.gpgsign true
   ```

4. **Add GPG Key to GitHub**:

   ```bash
   # Export your public key
   gpg --armor --export YOUR_KEY_ID
   
   # Copy the output and add it to GitHub:
   # GitHub → Settings → SSH and GPG keys → New GPG key
   ```

## 🚀 DevContainer Auto-Configuration

The DevContainer automatically:

1. **Mounts GPG Directory**: Your `~/.gnupg` is available in the container
2. **Sets Environment Variables**: `GPG_TTY` for proper terminal interaction
3. **Configures Git**: Auto-detects your GPG keys and configures signing
4. **Fixes Permissions**: Ensures proper GPG directory permissions

## 🛠️ Manual Configuration (if needed)

If auto-configuration doesn't work, run these commands in the DevContainer:

```bash
# Fix GPG permissions
chmod 700 ~/.gnupg
chmod 600 ~/.gnupg/*

# List available keys
gpg --list-secret-keys --keyid-format LONG

# Configure git signing
git config --global user.signingkey YOUR_KEY_ID
git config --global commit.gpgsign true

# Test signing
echo "test" | gpg --clearsign
```

## 🧪 Testing GPG Signing

1. **Make a test commit**:

   ```bash
   echo "GPG test" > test-gpg.txt
   git add test-gpg.txt
   git commit -m "Test GPG signing"
   ```

2. **Verify the signature**:

   ```bash
   git log --show-signature -1
   ```

   You should see: `gpg: Good signature from "Your Name <your.email@example.com>"`

## 🔧 Troubleshooting

### GPG Agent Issues

If you get "gpg: signing failed: No such file or directory":

```bash
# Kill existing GPG agent
gpgconf --kill gpg-agent

# Restart GPG agent
gpg-agent --daemon

# Test again
echo "test" | gpg --clearsign
```

### Permission Issues

```bash
# Fix all GPG permissions
chmod 700 ~/.gnupg
find ~/.gnupg -type f -exec chmod 600 {} \;
find ~/.gnupg -type d -exec chmod 700 {} \;
```

### Key Not Found

```bash
# Check if keys are available
gpg --list-secret-keys

# If empty, your host keys aren't mounted properly
# Check that your host ~/.gnupg directory exists and has keys
```

### Windows-Specific Issues

If using Windows, make sure:

1. **Git for Windows** is installed (includes GPG)
2. **GPG keys exist** in `%USERPROFILE%\.gnupg`
3. **Use Git Bash** or WSL for key generation
4. **File paths** use forward slashes in DevContainer

## 🎯 Best Practices

1. **Use a Separate Signing Subkey**: Generate a dedicated signing subkey
2. **Set Expiration Dates**: Use expiring keys for better security
3. **Backup Your Keys**: Export and securely store your private keys
4. **Use Hardware Tokens**: Consider hardware security keys for production

## 📚 Additional Resources

- [GitHub GPG Documentation](https://docs.github.com/en/authentication/managing-commit-signature-verification)
- [Git GPG Signing](https://git-scm.com/book/en/v2/Git-Tools-Signing-Your-Work)
- [GPG Best Practices](https://riseup.net/en/security/message-security/openpgp/best-practices)

The DevContainer handles most of the complexity for you - your host GPG setup should work seamlessly inside the container! 🔐✨
