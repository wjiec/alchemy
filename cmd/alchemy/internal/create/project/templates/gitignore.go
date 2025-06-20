package templates

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

type GitIgnore struct{}

func (g *GitIgnore) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (g *GitIgnore) Path() string { return ".gitignore" }
func (g *GitIgnore) Body() string { return gitIgnoreTemplate }

const gitIgnoreTemplate = `# General
.DS_Store
.AppleDouble
.LSOverride

# Icon must end with two \r
Icon

# Thumbnails
._*

# Files that might appear in the root of a volume
.DocumentRevisions-V100
.fseventsd
.Spotlight-V100
.TemporaryItems
.Trashes
.VolumeIcon.icns
.com.apple.timemachine.donotpresent

# Directories potentially created on remote AFP share
.AppleDB
.AppleDesktop
Network Trash Folder
Temporary Items
.apdisk

# If you prefer the allow list template instead of the deny list, see community template:
# https://github.com/github/gitignore/blob/main/community/Golang/Go.AllowList.gitignore
#
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with "go test -c"
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out
`
