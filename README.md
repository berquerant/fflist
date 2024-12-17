# fflist

Search for media files and display the paths of matching files.

## Installation

``` shell
./task build
```

## Requirements

- ffprobe 7.1 https://ffmpeg.org/ffprobe.html

## Usage

``` shell
❯ fflist query --help
Search for media files and display the paths of matching files.

The QUERY is a string in the 'key=value' format.
The 'key' refers to the name of the metadata of the media file, and the 'value' is a regular expression.
If the value corresponding to 'key' matches the 'value', the file path will be output to standard output.
If there are multiple QUERY, the file paths that match all QUERY will be output.

If the QUERY contains "or" or "OR," it splits the QUERY into groups based on that.
Within a group, conditions are evaluated with AND, while between groups, they are evaluated with OR.
For example, in "name=NAME1 OR name=NAME2 artist=ARTIST" the output will include files that either meet name=NAME1 or both name=NAME2 and artist=ARTIST.

The available 'key' include the following:

- name: The name of the file
- path: The path of the file
- mode: The file permissions (in octal)
- mod_time: The last modification time of the file
- size: The file size (in bytes)
- dir: All but the last element of path
- ext: The file name extension
- basename: name but ext
- basepath: path but ext

Depending on the type of media file, the following 'key' may also be available:

- album
- artist
- composer
- genre

Note: All metadata values are interpreted as strings.

To check which 'key' are actually available, please use the 'fflist debug' command or the '--verbose' option.

Using sh 'key' allows you to execute a sh script and output the file path only if the exit status is 0.
The value of sh 'key' is the main body of the script, and environment variables are available.
The script receives the entire file metadata in jsonl format from standard input.
For example, by writing a QUERY like the following, you can output only the file paths of files that exceed 8,000,000 bytes in size.

  'sh=jq "select((.size|tonumber) > 8000000).name" -r | grep -E ".+" -q'

Using the '--config' option allows you to specify the search directory and QUERY from a file.
The file has the following format:

root:
  - ROOT1
  - ROOT2

query:
  - - name=NAME1
  - - name=NAME2
    - artist=ARTIST

or

{
  "root": ["ROOT1", "ROOT2"],
  "query": [
    ["name=NAME1"],
    ["name=NAME2", "artist=ARTIST"]
  ]
}

'query' is an array of query groups.
Nested conditions are evaluated with AND, while top-level conditions are evaluated with OR.
In the above example, it means 'name=NAME1 OR (name=NAME2 AND artist=ARTIST)'.

When the '--config' option is specified, the '--root' option and QUERY arguments are ignored.

You can use environment variables (e.g. '$VARNAME') in the file specified by the --config option, as well as in the --root option and QUERY arguments.

Exmaples:
# in ~/Music, match name
fflist query -r ~/Music 'name=NAME'
# in ~/Music, match artist and genre
fflist query -r ~/Music 'artist=ARTIST' 'genre=GENRE'
# in ~/Music, either meet name=NAME1 or both name=NAME2 and artist=ARTIST
fflist query -r ~/Music name=NAME1 OR name=NAME2 artist=ARTIST
# read paths from stdin, match name
fflist query -r - name=NAME < path.list
# create index of ~/Music
fflist query -r ~/Music --createIndex > index
# in the index, match name
fflist query -r index --readIndex 'name=NAME'

Usage:
  fflist query [QUERY...] [flags]

Flags:
  -c, --config string   Query config file
      --createIndex     Dump all metadata
  -h, --help            help for query
      --readIndex       Read metadata from root, as index file
  -r, --root strings    Root directories (default [.])
  -v, --verbose         Verbose output
  -w, --worker int      Probe worker num (default 8)

Global Flags:
      --debug          Enable debug logs
  -p, --probe string   Media analyzer command (default "ffprobe")
  -q, --quiet          Quiet logs
```
