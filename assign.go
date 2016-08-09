package main
import "fmt"
import "github.com/google/go-github/github"
import "golang.org/x/oauth2"
import "math/rand"


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
    //error checking: if len(possibleAssignees) == 0
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

    if len(unassignedIssues) == 0 {
        fmt.Println("You have no unassignedIssues")
        return
    }

    for issue := range unassignedIssues {
        assignRandomly(org, *curRepo.Name, client, issue, possibleAssignees)
    }
}

func assignRandomly(org string, repo string, client *github.Client, issue *github.Issue, possibleAssignees []*github.User) {
    teamSize := len(possibleAssignees)
    rand.Seed(int64(teamSize))
    assignedDev := possibleAssignees[rand.Intn(teamSize)]
    iss, _, _ := client.Issues.AddAssignees(org, repo, *issue.Number, []string{*assignedDev.Login})
    if len(iss.Assignees) == 0 {
        fmt.Println(fmt.Sprintf("Unable to assign %s to issue number %d", *assignedDev.Login, *iss.Number))
    }
}