package main
import "fmt"
import "github.com/google/go-github/github"
import "golang.org/x/oauth2"
// import "math"


func main(){

    org := "sourcegraph-beta"
    userRepo := "sourcegraph-desktop"

    ts := oauth2.StaticTokenSource(
        &oauth2.Token{ AccessToken: "68ed6eb0012b7c189133900507ac2b5967a147c9"},
    )
    tc := oauth2.NewClient(oauth2.NoContext, ts)

    client := github.NewClient(tc)


    repos, resp, _ := client.Repositories.ListByOrg("sourcegraph-beta", &github.RepositoryListByOrgOptions{})

    fmt.Printf(resp.Status) //need error checking

    var curRepo *github.Repository
    for _,repo := range repos {
        fmt.Println(*repo.Name)
        if *repo.Name == userRepo {
            curRepo = repo
        }
    }

    fmt.Println(*curRepo.Name) //need error checking

    if (curRepo == nil) {return}
    possibleAssignees, _, _ := client.Issues.ListAssignees(org, *curRepo.Name, &github.ListOptions{})
    for _, assignee := range possibleAssignees {
        fmt.Println(*assignee.Login)
    }

    allIssues, _, _ := client.Issues.ListByRepo(org, *curRepo.Name, &github.IssueListByRepoOptions{})
    unassignedIssues := make(map[*github.Issue]bool)
    assignedIssues := make(map[*github.Issue]bool)

    for _, issue := range allIssues {
        if len(issue.Assignees) == 0 {
            unassignedIssues[issue] = true
        } else {
            assignedIssues[issue] = true
        }
    }

    for issue := range unassignedIssues {
        fmt.Println(issue.String())
    }
}
