# cataloger

Util for interact with ldap and active directory catalogs

## Install

```shell
wget https://github.com/dlampsi/cataloger/releases/...
chmod +x cataloger
mv ./cataloger /usr/local/bin/cataloger
```

## Usage

Full help avalible on `-h` or `--help` flags:
```shell
cataloger -h
```

Some examples bellow.

Get user info:

```shell
cataloger get dummyuser01 dummyuser01
cataloger get users dummyuser01 dummyuser01
```