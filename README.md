# cataloger

[![Actions Status](https://github.com/dlampsi/cataloger/workflows/default/badge.svg)](https://github.com/dlampsi/cataloger/actions)

Util for interact with various catalogs systems.

Avalible catalogs types:

- Active Directory (AD)

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

### Login

Login command, creates cataloger config and try to connect to catalog.

```bash
cataloger login \
    --host ad.server.local \
    --port 636 \
    --ssl --insecure \
    --bind-dn "CN=noname,OU=unit,DC=company,DC=com" \
    --bind-pass "fake" \
    --base "OU=people,DC=company,DC=com"
```

After that command cataloger creates config file in `$HOME/.cataloger/config.json`.

Example:

```json
{
  "auth": {
    "bind_dn": "CN=noname,OU=unit,DC=company,DC=com",
    "bind_pass": "ZmFrZQ=="
  },
  "params": {
    "group_search_base": "",
    "search_base": "OU=people,DC=company,DC=com",
    "user_search_base": ""
  },
  "server": {
    "host": "ad.server.local",
    "insecure": true,
    "port": 636,
    "ssl": true
  }
}
```

**Note**: Password storing in base64 encoding.

### Custom config file

You can provide custom config file via `-c` flag:

```bash
cataloger get dummyuser01 -c ~/custom_conf.json
```

### Other examples

Some usage examples bellow:

```bash
# Get one user info
cataloger get dummyuser01
cataloger get dummyuser01 dummygroup01

# Get one group info
cataloger get groups dummygroup01
cataloger get groups dummygroup01 dummygroup02

# Add members to group
cataloger mod group dummygroup01 -a dummyuser01
cataloger mod group dummygroup01 -a dummyuser01 -a dummyuser02

# Delete members from group
cataloger mod group dummygroup01 -d dummyuser01
cataloger mod group dummygroup01 -d dummyuser01 -d dummyuser02
```
