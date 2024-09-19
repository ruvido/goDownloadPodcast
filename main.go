package main

import (
    "log"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strings"
    "regexp"

    "github.com/mmcdole/gofeed"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
    "github.com/JohannesKaufmann/html-to-markdown"
)

var downloadFiles  = false
var createMetadata = true

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please provide the RSS file path OR a valid RSS url")
        return
    }

    input := os.Args[1]

    // Parse the RSS feed from URL or file
    var feed *gofeed.Feed
    var err error

    if isValidURL(input) {
        fmt.Println("Reading RSS from URL:", input)
        feed, err = parseRSSFromURL(input)
    } else {
        fmt.Println("Reading RSS from file:", input)
        feed, err = parseRSSFromFile(input)
    }

    if err != nil {
        fmt.Println("Error reading RSS:", err)
        return
    }
    // file, err := os.Open(os.Args[1])
    // if err != nil {
    //     fmt.Println("Error opening file:", err)
    //     return
    // }
    // defer file.Close()

    // fp := gofeed.NewParser()
    // feed, err := fp.Parse(file)
    // if err != nil {
    //     fmt.Println("Error parsing feed:", err)
    //     return
    // }

    // Directory setup
    podcastDir := "podcast"
    os.MkdirAll(podcastDir, os.ModePerm)

    for _, item := range feed.Items {
        // Process only season 4 episodes
        if item.ITunesExt.Season == "4" {
            // Format season and episode with leading zeros
            epType := item.ITunesExt.EpisodeType
            seasonNumber := fmt.Sprintf("%02d", parseSeasonNumber(item.ITunesExt.Season))
            episodeNumber := fmt.Sprintf("%02d", parseEpisodeNumber(item.ITunesExt.Episode))
            if epType == "bonus" {
                episodeNumber = epType
            }

            season := fmt.Sprintf("s%s", seasonNumber)
            episode := fmt.Sprintf("e%s", episodeNumber)
            title := item.Title
            slug := slugify(title)
            epDate := item.PublishedParsed.Format("2006-01-02")
            alias := fmt.Sprintf("/%s-%s", season, episode)
            if epType == "bonus" {
                alias = ""
            }
            filename := fmt.Sprintf("%s-%s-%s.mp3", season, episode, slug)
            // Debug print
            log.Printf("Season: %s, Episode: %s| %s %s | %s\n", seasonNumber, episodeNumber, epDate, epType, title)

            episodesDir := filepath.Join(podcastDir, "episodes", fmt.Sprintf("season-%s", season))
            contentDir := filepath.Join(podcastDir, "content", fmt.Sprintf("season-%s", season))

            filepathMp3 := filepath.Join(episodesDir, filename)

            if downloadFiles {
                os.MkdirAll(episodesDir, os.ModePerm)
                fmt.Printf("Downloading %s...\n", filename)
                if err := downloadFile(item.Enclosures[0].URL, filepathMp3); err != nil {
                    fmt.Printf("Error downloading %s: %v\n", filename, err)
                    continue
                }
            }

            if createMetadata {
                os.MkdirAll(contentDir, os.ModePerm)
                description, err := htmlToMarkdown(item.Description)
                if err != nil {
                    fmt.Printf("Error converting description to Markdown")
                }

                // Create metadata file
                metadata := fmt.Sprintf(`---
title:    "%s"
season:   "%s"
number:   "%s"
date:     "%s"
file:     "episodes/%s/%s"
length:   "%s"
duration: "%s"
guid:     "%s"
aliases:  ["%s"]
slug:     "%s"
---
%s
                `, title, seasonNumber, episodeNumber, item.PublishedParsed.Format("2006-01-02"), season, filename, item.Enclosures[0].Length, item.ITunesExt.Duration, item.GUID, alias, slug, description)

                filepathMd := filepath.Join(contentDir, fmt.Sprintf("%s-%s.md", alias, slug))
                fmt.Printf("Writing metadata to %s\n", filepathMd)
                if err := os.WriteFile(filepathMd, []byte(metadata), 0644); err != nil {
                    fmt.Printf("Error writing metadata for %s: %v\n", filename, err)
                }
            }
        }
    }
}

func parseSeasonNumber(season string) int {
    var seasonNumber int
    fmt.Sscanf(season, "%d", &seasonNumber)
    return seasonNumber
}

func parseEpisodeNumber(episode string) int {
    var episodeNumber int
    fmt.Sscanf(episode, "%d", &episodeNumber)
    return episodeNumber
}

// func sanitizeTitle(title string) string {
//     // Normalize and remove special characters from the title
//     t := transform.Chain(norm.NFD, transform.RemoveFunc(isNotASCII), norm.NFC)
//     sanitized, _, _ := transform.String(t, title)
//     return sanitized
// }

func isNotASCII(r rune) bool {
    return r > 127
}

func downloadFile(url string, filepath string) error {

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Set the custom User-Agent (mimic wget)
	req.Header.Set("User-Agent", "Wget/1.21.1")

	// Create an HTTP client
	client := &http.Client{}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	// Create the file on disk
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the response body to the file
	_, err = io.Copy(out, resp.Body)
	return err
}

func parseRSSFromFile(path string) (*gofeed.Feed, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    fp := gofeed.NewParser()
    return fp.Parse(file)
}


func parseRSSFromURL(rssURL string) (*gofeed.Feed, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", rssURL, nil)
    if err != nil {
        return nil, err
    }

    // Set a User-Agent to mimic a normal browser or wget
    req.Header.Set("User-Agent", "Wget/1.20.3 (linux-gnu)")

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Read the body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Uncomment this to debug the body
    // fmt.Println("Response Body:\n", string(body))

    //
    // fmt.Println("Response Body:\n", string(body)) // Print the response body for debugging

    // Re-parse the body using a new reader
    fp := gofeed.NewParser()
    return fp.Parse(strings.NewReader(string(body)))
}

func isValidURL(str string) bool {
    u, err := url.Parse(str)
    return err == nil && u.Scheme != "" && u.Host != ""
}

func slugify(title string) string {
	// Convert to lowercase
	title = strings.ToLower(title)

   // Normalize and remove non-ASCII characters (like accents)
    t := transform.Chain(norm.NFD, transform.RemoveFunc(isNotASCII), norm.NFC)
    normalized, _, _ := transform.String(t, title)

	// Replace spaces and special characters with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug := re.ReplaceAllString(normalized, "-")

	// Trim hyphens from the start and end
	slug = strings.Trim(slug, "-")

	return slug
}


func htmlToMarkdown(html string) (string, error) {
    converter := md.NewConverter("", true, nil)
    markdown, err := converter.ConvertString(html)
    if err != nil {
        return "", err
    }
    return markdown, nil
}
