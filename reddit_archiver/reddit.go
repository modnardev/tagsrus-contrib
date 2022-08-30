package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	snoo "github.com/modnardev/go-reddit"
)

var tagsrusHost = "http://192.168.1.100:8080"

// var tagsrusHost = "http://localhost:8080"

var credentials = snoo.Credentials{}

func init() {
	credentials.ID = os.Getenv("REDDIT_ID")
	credentials.Secret = os.Getenv("REDDIT_Secret")
	credentials.Username = os.Getenv("REDDIT_Username")
	credentials.Password = os.Getenv("REDDIT_Password")
}

func getUpvoted(client *snoo.Client, c chan *snoo.Post) error {
	after := ""

	for {
		opts := snoo.ListUserOverviewOptions{
			ListOptions: snoo.ListOptions{
				After: after,
				Limit: 100,
			},
			Time: "all",
			Sort: "new",
		}

		posts, _, err := client.User.Upvoted(context.Background(), &opts)
		if err != nil {
			return err
		}

		if len(posts) == 0 {
			break
		}

		for _, post := range posts {
			c <- post
		}

		after = posts[len(posts)-1].FullID
	}

	return nil
}

func getSaved(client *snoo.Client, c chan *snoo.Post) error {
	after := ""

	for {
		opts := snoo.ListUserOverviewOptions{
			ListOptions: snoo.ListOptions{
				After: after,
				Limit: 100,
			},
			Time: "all",
			Sort: "new",
		}

		posts, _, _, err := client.User.Saved(context.Background(), &opts)
		if err != nil {
			return err
		}

		if len(posts) == 0 {
			break
		}

		for _, post := range posts {
			c <- post
		}

		after = posts[len(posts)-1].FullID
	}

	return nil
}

func main() {
	posts := make(chan *snoo.Post)
	go func() {
		browserExtURL, _ := url.Parse(tagsrusHost + "/api/import")

		for post := range posts {
			q := browserExtURL.Query()
			q.Set("link", fmt.Sprintf("https://reddit.com%v", post.Permalink))

			browserExtURL.RawQuery = q.Encode()
			req, err := http.NewRequestWithContext(context.Background(), "GET", browserExtURL.String(), nil)
			if err != nil {
				fmt.Println("error:", err)
			}

			resp, err := (&http.Client{}).Do(req)
			if err != nil {
				fmt.Println("error:", err)
			}
			resp.Body.Close()

			fmt.Printf("[%v] https://reddit.com%v\n", post.Created.Format("2006-01-02"), post.Permalink)
		}
	}()

	httpClient := &http.Client{}
	client, e := snoo.NewClient(credentials, snoo.WithUserAgent("tagsrus-archiver/0.1.0"), snoo.WithHTTPClient(httpClient))
	if e != nil {
		panic(e)
	}

	e = getUpvoted(client, posts)
	if e != nil {
		panic(e)
	}

	e = getSaved(client, posts)
	if e != nil {
		panic(e)
	}

	log.Println("done!")
}
