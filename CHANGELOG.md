# Changelog

This file lists the changes to Holdero repo with each version.

## 0.x.x - In Progress

### Added

* XSWD and DERO connections and calls
* Account for local DB storage
* Custom token balance support in PlaceSwap
* More error logging when downloading

### Changed

* Fyne 2.4.4
* dReams 0.12.0
* Updated card file paths
* Player handKey stored in account DB
* Update StartApp with upstream functions
* `dreams` UpdateText and SetUpdate for labels
* Remove DownloadFile and UnzipFile, pull them from `dreams`
* Combine SharedDecks widget into FaceSelect

### Fixed

* Stop sync if disconnected during
* Clear favorites list on disconnect


## 0.3.1 - January 19 2024

### Changed

* Go 1.21.5
* Fyne 2.4.3
* dReams 0.11.1
* Cleaned up `rpc` client var names

### Fixed

* Lessened object refreshes to reduce memory use


## 0.3.0 - December 23 2023

### Added

* CHANGELOG
* Pull request and templates
* `semver` versioning 
* HS gold cards
* Asset tabs with profile
* Sync screen
* Swap funcs
* Owners instructions

### Changed

* Fyne 2.4.1
* dReams 0.11.0
* Icon resources 
* Updated menu and owners controls layout
* Update dreams.AssetSelects for profile
* Rename unneeded exports
* Rename SortCardAsset to SortCardAssets
* Remove unneeded funcs exports
* Bet button text change on amount
* Condense Called into singleShot
* Player_label func image handling
* tag parma key funcs
* Confirmations to dialogs 
* implement `gnomes` and funcs
* implement `menu` ShowTxDialog and ShowConfirmDialog
* implement `rpc` PrintError, PrintLog and IsConfirmingTx

### Fixed

* Deprecated container.NewStack
* Fyne error when downloading custom
* Validator hangs