# Mitril

A simple CLI tool written in Go to locally download an entire podcast, including metadata. It organizes downloaded episodes into folders within a `download` directory, where:
- `audio/` contains the MP3 audio files.
- `content/` stores the episode metadata as Markdown pages.

These Markdown files are designed for integration with static site generators, such as Hugo, to create a self-hosted podcast website. The tool supports podcasts with seasons, but won't throw an error if the RSS feed does not include season information.

**Mitril** is primarily intended for migration from a hosted podcast service to a self-hosted solution, powered by Hugo.


## Features

- Downloads podcast episodes from RSS feeds.
- Stores episodes as MP3s in a structured folder system.
- Automatically generates Markdown metadata for each episode.
- Supports season and episode organization.
- Easy to use, simple command-line interface.

## Installation

To build the tool:

```bash
git clone https://github.com/ruvido/mitril.git
cd mitril
go build
```
in case of errors, check your go version (or remove the go version tag from the go.mod file)

## Usage

Get the help by running ``./mitril`` with no flags

```
Usage: mitril [options] <RSS_FILE_OR_URL>
Options:
  -d    Download audio files
  -m    Create metadata files
  -p    Prefix for web audio files
```
Run the tool by providing an RSS feed URL:

    ./mitril https://podcast.example.com/feed.rss

Alternatively, specify a RSS feed file:

    ./mitril feed.rss

without ``-d`` or ``-m`` a dry run will be executed. 
For a complete download, run:

    ./mitril -d -m https://podcast.example.com/feed.rss

## Folder Structure

Once the episodes are downloaded, the folder structure will look like this:

```plaintext
download/
├── audio/
│   ├── season-01/
│   │   ├── s01-e01-episode-title.mp3
│   │   ├── s01-e02-episode-title.mp3
│   └── season-02/
│       └── s02-e01-episode-title.mp3
└── content/
    └── podcast/
        ├── season-01/
        │   ├── s01-e01-episode-title.md
        │   ├── s01-e02-episode-title.md
        └── season-02/
            └── s02-e01-episode-title.md
```

Each episode will have a corresponding .md metadata file in the content/ folder, ready for use in Hugo or another static site generator.

## Example Metadata File

```
---
title:    "Episode Title"
season:   "01"
number:   "01"
date:     "2024-09-20"
audio:    "/season-01/s01-e01-episode-title.mp3"
length:   "12345678"
duration: "3512"
guid:     "abc123"
aliases:  ["/s01e01"]
slug:     "episode-title"
---
This is a brief description of the episode.
```

## Prefix for web audio files

Using ``-p`` adds a prefix to the filepath in the audio tag of the markdown files, e.g.

    ./mitril -d -m -p "https://storage.ex.com" https://podcast.example.com/feed.rss

modifies the yaml tag to:

    audio:    "https://storage.ex.com/season-01/s01-e01-episode-title.mp3"

this is useful when serving podcast pages and audio files from different servers.

## Contribution & Support

Feel free to submit issues if you encounter any problems, or suggest features via GitHub issues.

Happy podcast migrating!


