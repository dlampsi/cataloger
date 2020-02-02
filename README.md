# cataloger

[![Actions Status](https://github.com/dlampsi/cataloger/workflows/default/badge.svg)](https://github.com/dlampsi/cataloger/actions)

Util for interact with various catalogs systems.

Avalible catalogs types:

- Active Directory (AD)

## Install

You can download releases [here](https://github.com/dlampsi/cataloger/releases).

Example deploy script - [get_darwin_release.sh](scripts/get_darwin_release.sh) for darwin install `v0.0.1`:

```bash
./scripts/get_darwin_release.sh v0.0.1
```

### Completion

To configure your bash shell to load completions for each session add to your `~/.bashrc` or `~/.profile`:

```bash
. <(cataloger completion)
```

## Usage

Full help avalible on `-h` or `--help` flags:

```bash
cataloger -h
```

### Login

Login command, creates cataloger config and try to connect to catalog:

```bash
cataloger login
```

Or you can pass all flags to escape interactive asking:

```bash
cataloger login \
    --host ad.server.local \
    --port 636 \
    --ssl --insecure \
    --bind "CN=noname,OU=unit,DC=company,DC=com" \
    --password "fake" \
    --search-base "OU=people,DC=company,DC=com"
```

After that command cataloger creates config file in `$HOME/.cataloger.json`:

```json
{
  "host": "ad.server.local",
  "port": 636,
  "ssl": true,
  "insecure": true,
  "bind": "CN=noname,OU=unit,DC=company,DC=com",
  "password": "fake",
  "search-base": "OU=people,DC=company,DC=com",
}
```

### Custom config file

You can provide custom config file via `-c` flag:

```bash
cataloger -c ~/custom_conf.json
```

### Other examples

Search:

```bash
# Search for user entry by user sAMAccountName
cataloger search user dummyuser

# Search for user entry by user mail attribute
cataloger search user --attribute=mail dummyuser@fake.com

# Search for user entry by user sAMAccountName and display user groups
cataloger search user dummyUser -g

# Search for group entry by sAMAccountName
cataloger search group dummyGroup

# Search for group entry by 'cn' attribute
cataloger search group --attribute=cn dummyGroup dummyGroup-CN

# Search for group entry by sAMAccountName and display group direct members
cataloger search group dummyGroup -m

# Search for group entry by sAMAccountName and display group all members (include all subgroups members)
cataloger search group dummyGroup -m --nested
```

Modify:

```bash
# Add 'dummyUser' to 'dummyGroup' members
cataloger modify group members dummyGroup -a dummyUser

# Add 'dummyUser1' and 'dummyUser2' to 'dummyGroup' members
cataloger modify group members dummyGroup -a dummyUser1 -a dummyUser2

# Remove 'dummyUser' form 'dummyGroup' members
cataloger modify group members dummyGroup -d dummyUser

# Remove 'dummyUser1' and 'dummyUser2' form 'dummyGroup' members
cataloger modify group members dummyGroup -d dummyUser1 -d dummyUser2
```
