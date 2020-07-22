# Git FuzzyFind Commit Message (fcm)

fcm provides a Git commit message template fuzzy search feature.

![description](https://raw.githubusercontent.com/wiki/wataboru/git-fuzzy-find-commit-message/images/fcm_description.gif)

## Use as CLI

### Install

```
$ go get github.com/wataboru/git-fuzzy-find-commit-message
```

Download from [releases](https://github.com/wataboru/git-fuzzy-find-commit-message/releases)

### Usage

1. Run
```
$ fcm
```

2. Enter commit message and Choose templates
```
 null TextEditorComponent::domNode during visibility check           ┌────────────────────────────────────────────────────────────────────┐
  Avoid infinite recursion when bad values are passed to tz aware .. │  Add build script                                                  │
  Avoid distinct if a subquery has already materialized              │                                                                    │
  Add -enable-experimental-nested-generic-types frontend flag        │                                                                    │
  Add support for activating and deactivating package-specific key.. │                                                                    │
  Add a helper method mayHaveOpenedArchetypeOperands to SILInstruc.. │                                                                    │
  Add a typealias to avoid a build ordering dependency between pro.. │                                                                    │
  Add --main-process flag to run specs in the main process           │                                                                    │
  Add support for closure contexts to readMetadataFromInstance()     │                                                                    │
  Add a basic test for opening an editor in largeFileMode if >= 2M.. │                                                                    │
  Add support for allocators that require tensors with zero          │                                                                    │
  Add "event" parameter for "click" handler of MenuItem              │                                                                    │
  Add a design-decisions section to the CONTRIBUTING guide           │                                                                    │
  Add Throws flag and ThrowsLoc to AbstractFunctionDecl              │                                                                    │
  Add TODO about blinkFeatures -> enableBlinkFeatures                │                                                                    │
  Add validation test for projecting existentials                    │                                                                    │
  Add failing spec for Menu.buildFromTemplate                        │                                                                    │
  Add SkUserConfig.h with blank SkDebugf macro                       │                                                                    │
  Add support for launching HTML files directly                      │                                                                    │
  Add documentation for --proxy-bypass-list                          │                                                                    │
  Add assertions for no available bookmark                           │                                                                    │
  Add assert for role with app name in label                         │                                                                    │
  Add convenience API for demangling                                 │                                                                    │
  Add specs for moveSelectionLeft()                                  │                                                                    │
  Add comment about map key/values                                   │                                                                    │
  Add TypeLowering::hasFixedSize()                                   │                                                                    │
  Add File > Exit menu on Windows                                    │                                                                    │
  Add tests for pending pane items                                   │                                                                    │
  Add missing period in comment                                      │                                                                    │
  Add docs for app.getLocale()                                       │                                                                    │
  Add asserts for properties                                         │                                                                    │
  Add style.less examples                                            │                                                                    │
  Add overflow scrolling                                             │                                                                    │
  Add npm start script                                               │                                                                    │
  Add missing return                                                 │                                                                    │
> Add build script                                                   │                                                                    │
  38/125                                                             │                                                                    │
> Add
```

3. Edit message and save
```
Add Makefile and build script

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch master
# Your branch is up to date with 'origin/master'.
#
# Changes to be committed:
#       modified:   Makefile
#       modified:   .github/workflows/release.yml
#
```

4. Commited
```
[master ******] Add Makefile and build script
 2 files changed, 114 insertions(+), 8 deletions(-)
 rewrite LICENSE (79%)
```

### Version

```
$ fcm -v
```

## Generate files

This app generates the following files.
- ~/.fcm  
  Message Template. You can add your own additions to increase the number of Fuzzy Find candidates.
- ~/.fcm_history  
  Each time you commit using fcm, the history is added to this page. The history is also a candidate for a Fuzzy Find.

### Format

- `~/.fcm` or `~/.fcm_history`
```
# This line is CommentOut, Not use.
FuzzyFind candidate1
FuzzyFind candidate2
FuzzyFind candidate3
```

## Use as a libary

- https://github.com/ktr0731/go-fuzzyfinder
- https://github.com/tcnksm/ghr

## Refer to the following for implementation

- https://github.com/skanehira/fk  
I referred to CI and Build as a whole.
- https://anond.hatelabo.jp/20160725092419  
I got all the default message templates from here.
