# cataloger

[![GoDoc](https://godoc.org/github.com/dlampsi/cataloger?status.svg)](https://godoc.org/github.com/dlampsi/cataloger) [![Actions Status](https://github.com/dlampsi/cataloger/workflows/badge.svg)](https://github.com/dlampsi/cataloger/actions)

Util for interact with ldap and active directory catalogs

## Install

```bash
wget https://github.com/dlampsi/cataloger/releases/...
chmod +x cataloger
mv ./cataloger /usr/local/bin/cataloger
```

## Usage

Full help avalible on `-h` or `--help` flags:
```bash
cataloger -h
```

Some examples bellow:
```bash
# Get info
cataloger get dummyuser01 dummygroup01
# Get users infor
cataloger get users dummyuser01 dummyuser02
# Get groups info
cataloger get groups dummygroup01 dummygroup02

```