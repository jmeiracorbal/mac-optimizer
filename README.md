# Mac Optimizer 🚀⚡

A comprehensive macOS optimization and cleanup tool that improves system performance through intelligent cache management, service optimization, and resource cleanup.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![macOS](https://img.shields.io/badge/macOS-Monterey%2B-blue.svg)](https://www.apple.com/macos/)
[![Homebrew](https://img.shields.io/badge/Homebrew-Compatible-orange.svg)](https://brew.sh/)
[![Shell](https://img.shields.io/badge/Shell-Bash-green.svg)](https://www.gnu.org/software/bash/)

## 🚀 Quick Start

### Install via Homebrew (Recommended)
```bash
brew tap jmeiracorbal/tools
brew install mac-optimizer
mac-optimizer
```

### Manual Installation
```bash
# Download and run directly
curl -O https://raw.githubusercontent.com/jmeiracorbal/mac-optimizer/main/mac-optimizer.sh
chmod +x mac-optimizer.sh
./mac-optimizer.sh
```

## ✨ Features

- **🧹 Smart Cache Cleanup**: Removes system and application caches while preserving important data
- **🔍 Spotlight Optimization**: Rebuilds search index for improved performance (with confirmation)
- **🔐 Permission Repair**: Fixes system permissions and directory access rights
- **⚡ Service Management**: Restarts problematic system services (Finder, Dock, etc.)
- **💾 Memory Liberation**: Frees up RAM and cleans memory pressure
- **📝 Log Cleanup**: Removes old system and application logs
- **🌐 Browser Cache Cleanup**: Clears browser caches and temporary data
- **📁 Downloads Management**: Cleans old files from Downloads folder
- **🔒 Safe Operations**: Dry-run mode and user confirmations for destructive actions
- **📊 System Statistics**: Before/after performance metrics

## 🛠️ Usage

### Basic Usage
```bash
# Complete system optimization (interactive)
mac-optimizer

# Preview what would be done
mac-optimizer --dry-run

# Silent automated optimization
mac-optimizer --force --quiet
```

### Module Selection
```bash
# Only clean caches and optimize Spotlight
mac-optimizer cache spotlight

# Clean browsers and downloads only
mac-optimizer browsers downloads

# Run specific optimizations
mac-optimizer memory services permissions
```

### Command Line Options
```bash
mac-optimizer [OPTIONS] [MODULES]

OPTIONS:
  -h, --help          Show help message
  -v, --version       Show version information
  -d, --dry-run       Preview changes without executing
  -q, --quiet         Run in silent mode
  -f, --force         Skip confirmations
  --log-file FILE     Custom log file location

MODULES:
  cache              Clean system and application caches
  spotlight          Optimize and reindex Spotlight
  permissions        Repair system permissions
  services           Restart problematic services
  memory             Free memory and clean processes
  logs               Clean old system logs
  browsers           Clean browser caches and data
  downloads          Clean old downloads (with confirmation)
  all                Run all modules (default)
```

## 📊 Performance Impact

### Typical Space Recovery
- **Cache Cleanup**: 500MB - 2GB+ freed
- **Browser Data**: 100MB - 1GB+ freed
- **Old Downloads**: Varies by user habits
- **System Logs**: 50MB - 500MB freed

### System Improvements
- **Memory Usage**: 5-15% improvement
- **Spotlight Speed**: Faster search indexing
- **App Launch Times**: Reduced cache overhead
- **System Responsiveness**: Cleaner service states

## 🔒 Safety Features

### Confirmation Prompts
Operations that require user confirmation:
- **Spotlight Optimization**: Requires admin permissions, takes time
- **Browser Cleanup**: May require re-login to websites
- **Downloads Cleanup**: Prevents accidental file deletion

### Backup & Recovery
- **Automatic Backups**: Critical files backed up before changes
- **Detailed Logging**: Complete operation history in `~/.mac-optimizer.log`
- **Dry-Run Mode**: Preview all changes before execution
- **Non-Destructive**: Preserves directory structures and important data

## 📱 Example Output

```
=== Mac Optimizer v1.0.0 ===
Starting system optimization...

=== System Statistics ===
Free Memory: 53%
Disk Usage: 84%
Uptime: 4 days

🔄 Starting cache cleanup...
ℹ️  Cleaning /Users/user/Library/Caches (size: 720M)
✅ Cache cleaned: /Users/user/Library/Caches (now: 0B)
ℹ️  Cleaning DNS cache
✅ Cache cleanup completed

🔄 Spotlight optimization requested...
⚠️  Spotlight optimization will:
   • Require administrator permissions
   • Temporarily disable indexing
   • Rebuild the search index (takes time)
   • May affect search functionality temporarily
Continue with Spotlight optimization? (y/N): y

🔄 Freeing system memory...
ℹ️  Memory purge completed
✅ Memory freed

🎉 Optimization completed!
📋 Log saved at: ~/.mac-optimizer.log

=== System Statistics ===
Free Memory: 67% (+14%)
Disk Usage: 82% (-2%)
Uptime: 4 days
```

## 🧹 Cleanup Modules

### Cache Cleanup (`cache`)
- **User Caches**: `~/Library/Caches/*`
- **Application Support**: Crash reports, temporary files
- **WebKit Data**: Browser engine caches
- **System Caches**: `/Library/Caches/*` (requires admin)
- **DNS Cache**: System DNS resolver cache

### Spotlight Optimization (`spotlight`)
- **Index Rebuild**: Complete search index reconstruction
- **Performance Boost**: Faster file and content searching
- **Background Process**: Continues optimization after completion
- **Admin Required**: Needs administrator permissions

### Browser Cleanup (`browsers`)
- **Chrome**: Cache, application cache, GPU cache
- **Firefox**: Profile caches, startup cache
- **Safari**: Browser caches, local storage, databases
- **WebKit**: Shared browser engine data

### Downloads Management (`downloads`)
- **Age-Based**: Only files older than 30 days
- **Size Reporting**: Shows total space that will be freed
- **File Preview**: Lists sample files before deletion
- **User Confirmation**: Prevents accidental data loss

## 🔧 Technical Details

### System Requirements
- **macOS**: Monterey (12.0) or later
- **Architecture**: Intel and Apple Silicon supported
- **Permissions**: Some operations require administrator access
- **Dependencies**: Standard macOS utilities (no external dependencies)

### File Locations
- **Log File**: `~/.mac-optimizer.log`
- **Backup Location**: Same directory with `.backup` suffix
- **Target Directories**: Standard macOS cache and data locations

### Performance Considerations
- **CPU Usage**: Minimal during operation
- **Disk I/O**: Moderate during cache cleanup
- **Memory Usage**: Low script footprint
- **Network**: DNS cache flush may briefly affect connectivity

## 🍺 Homebrew Installation

### Add the Tap
```bash
brew tap jmeiracorbal/tools
```

### Install the Package
```bash
brew install mac-optimizer
```

### Update to Latest Version
```bash
brew update
brew upgrade mac-optimizer
```

### Uninstall
```bash
brew uninstall mac-optimizer
brew untap jmeiracorbal/tools  # Optional: remove the tap
```

## 🔄 Automation & Scripting

### Cron Job Example
```bash
# Run optimization weekly on Sundays at 3 AM
0 3 * * 0 /usr/local/bin/mac-optimizer --force --quiet
```

### CI/CD Integration
```bash
# Pre-deployment cleanup
mac-optimizer cache memory --force
```

### Custom Workflows
```bash
# Development environment reset
mac-optimizer cache browsers downloads --force

# Performance troubleshooting
mac-optimizer --dry-run > system-analysis.log
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup
```bash
git clone https://github.com/jmeiracorbal/mac-optimizer.git
cd mac-optimizer
chmod +x mac-optimizer.sh
./mac-optimizer.sh --dry-run
```

### Testing
```bash
# Test all modules in dry-run mode
./mac-optimizer.sh --dry-run all

# Test specific functionality
./mac-optimizer.sh --dry-run cache spotlight
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ FAQ

**Q: Is it safe to run mac-optimizer regularly?**  
A: Yes! The script is designed for regular use with built-in safety mechanisms and confirmations for potentially disruptive operations.

**Q: Will this break my system or applications?**  
A: No. The script only cleans temporary files, caches, and logs. It doesn't modify system files or application binaries.

**Q: How often should I run this?**  
A: Weekly or bi-weekly for regular users. Daily for heavy users or development environments.

**Q: What if Spotlight indexing takes too long?**  
A: Spotlight rebuilding happens in the background. You can continue using your Mac normally while it completes.

**Q: Can I run this on company/managed Macs?**  
A: Check with your IT department first, as some operations require administrator permissions.

**Q: Does this work with external drives?**  
A: The script focuses on system optimization. External drives are not affected.

## ⚠️ Important Notes

- **Backup First**: While safe, consider backing up important data before first use
- **Admin Rights**: Some optimizations require administrator permissions
- **Browser Sessions**: Browser cleanup may require re-logging into websites
- **Spotlight Time**: Index rebuilding can take 30 minutes to several hours depending on data volume
- **System Restart**: Consider restarting after major optimizations for best results

## 📈 Performance Monitoring

Monitor the effectiveness using the built-in statistics:

```bash
# Before optimization
mac-optimizer --dry-run

# Run optimization
mac-optimizer

# Check improvements in the summary
```

## 🔗 Related Projects

- [Cloud Storage Symlinks](https://github.com/jmeiracorbal/cloud-storage-symlinks) - Organize cloud storage access
- [Homebrew](https://brew.sh/) - Package manager for macOS

---

**Made with ❤️ for the macOS community**

*Keep your Mac running at peak performance!* 🚀
