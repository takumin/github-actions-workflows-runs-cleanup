package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func main() {
	// Set Flags
	var token, owner, repos, workflow string
	flag.StringVar(&token, "token", "", "github token")
	flag.StringVar(&owner, "owner", "", "github owner")
	flag.StringVar(&repos, "repos", "", "github repository")
	flag.StringVar(&workflow, "workflow", "", "github actions workflow name")
	flag.Parse()

	// Check Flags
	if token == "" || owner == "" || repos == "" || workflow == "" {
		log.Fatalln("require [token|owner|repos|workflow] flag.")
	}

	// Authentication
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// GitHub Client
	client := github.NewClient(tc)

	// Get Workflows
	workflows, _, err := client.Actions.ListWorkflows(ctx, owner, repos, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get Workflow ID
	var id int64
	for _, v := range workflows.Workflows {
		if *v.Name != workflow {
			continue
		}
		id = v.GetID()
	}

	// Set Pagination
	opt := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// Get Workflow Run IDs
	var ids []int64
	for {
		runs, resp, err := client.Actions.ListWorkflowRunsByID(ctx, owner, repos, id, opt)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range runs.WorkflowRuns {
			ids = append(ids, v.GetID())
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// Delete Workflow Runs
	for _, v := range ids {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%d", owner, repos, v)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Set Header
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

		// Delete Request
		resp, err := tc.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()

		// Logger
		log.Println("delete workflow run:", url)
	}
}
