# Podcast local download

A simple cli tool written in go to locally download an entire podcast (including metadata) 

## Installation

    go build

## Usage

    ./goDownloadPodcast feed.rss

or specifying the feed url

    ./goDownloadPodcast https://podcast.example.com/feed.rss

the program creates a folder named **podcast** where episodes and corresponding metadata are downloaded according to the season
