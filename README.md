# BitBin
BitBin is a Work-In-Progress aimed to be a more lightweight but mostly compatible alternative to 
[ByteBin](https://github.com/lucko/bytebin/).
Due to it being written in Go instead of Java, it comes with a by-default lower resource footprint.

> [!NOTE]
> This project is not meant as a drop-in replacement for ByteBin. While the goal is to maintain user-side compatibility, 
> it will still differ from ByteBin itself in many aspects.

---
### Configuration
Configuration can be done by using environment variables. Each one must be prefixed with `BITBIN_`.
A list of all possible environment variables and their functions:

| Variable Name            | Default Value                                 | Description                                                                                                                                                     |
|--------------------------|-----------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `BITBIN_PORT`            | 8080                                          | The port this application will listen on.                                                                                                                       |
| `BITBIN_DIRECTORY`       | `/data/bitbin`                                | The storage directory for the DB/files                                                                                                                          |
| `BITBIN_LINK_TITLE`      | `source`                                      | Title for the link shown on the right side of the navbar.                                                                                                       |
| `BITBIN_LINK`            | `https://git.gamecrash.dev/game.crash/BitBin` | Address of the link on the right side of the navbar.                                                                                                            |
| `BITBIN_PASTE_TTL`       | 43200                                         | The lifetime of a paste in minutes; will not expire if value is 0. After existing for the given time, it will be deleted.  Defaults to 30 days (43200 minutes). |
| `BITBIN_MAX_UPLOAD_SIZE` | 10485760                                      | The maximum upload size of a paste in MiB. Defaults to 10MiB.                                                                                                   |
