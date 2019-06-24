#!/bin/bash
set -x
echo package engine                                           > engine/overrideversion.go
echo import \"fmt\"                                           >> engine/overrideversion.go
echo "// This is an autogenerated for GoLand builds"          >> engine/overrideversion.go
echo "// It should not exist during regular build procedures" >> engine/overrideversion.go
echo func init\(\) \{                                         >> engine/overrideversion.go
echo    if Build !\= \"\" \{ return \}                       >> engine/overrideversion.go
echo    fmt.Println\(\"Build and version come from overrideversion.go\"\)   >> engine/overrideversion.go
echo    Build = \"GoLand-`git rev-parse HEAD``git status | grep -Eo modified |sort -u`\"               >> engine/overrideversion.go
echo    FactomdVersion = \"GoLand-`cat VERSION``git status | grep -Eo modified |sort -u`\"              >> engine/overrideversion.go
echo \}                                                       >> engine/overrideversion.go
