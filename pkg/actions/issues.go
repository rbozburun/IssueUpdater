package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//var wg sync.WaitGroup

// Issue item for requesting to update
type SingleIssueReqItem struct {
	IssueID string `json:"IssueId"`
	State   string `json:"State"`
}

type ScanResultIssue []struct {
	IssueUrl string `json:"IssueUrl"`
	Title    string `json:"Title"`
	Type     string `json:"Type"`
	Url      string `json:"Url"`
}

type AllIssuesIssue struct {
	AssigneeName                   string `json:"AssigneeName"`
	FirstSeenDate                  string `json:"FirstSeenDate"`
	ID                             string `json:"Id"`
	IsAddressed                    bool   `json:"IsAddressed"`
	IsDetectedByShark              bool   `json:"IsDetectedByShark"`
	IsPresent                      bool   `json:"IsPresent"`
	LastSeenDate                   string `json:"LastSeenDate"`
	StateFixedConfirmedDate        string `json:"StateFixedConfirmedDate"`
	UpdatedDate                    string `json:"UpdatedDate"`
	Severity                       string `json:"Severity"`
	State                          string `json:"State"`
	Title                          string `json:"Title"`
	URL                            string `json:"Url"`
	LatestVulnerabilityIsConfirmed bool   `json:"LatestVulnerabilityIsConfirmed"`
	WebsiteID                      string `json:"WebsiteId"`
	WebsiteName                    string `json:"WebsiteName"`
	WebsiteRootURL                 string `json:"WebsiteRootUrl"`
}

type AllIssuesApiResult struct {
	FirstItemOnPage     int              `json:"FirstItemOnPage"`
	HasNextPage         bool             `json:"HasNextPage"`
	HasPreviousPage     bool             `json:"HasPreviousPage"`
	HasPreviIsFirstPage bool             `json:"HasPreviIsFirstPage"`
	IsLastPage          bool             `json:"IsLastPage"`
	LastItemOnPage      int              `json:"LastItemOnPage"`
	List                []AllIssuesIssue `json:"List"`
	PageCount           int              `json:"PageCount"`
	PageNumber          int              `json:"PageNumber"`
	PageSize            int              `json:"PageSize"`
	TotalItemCount      int              `json:"TotalItemCount"`
}

func IssueActions(uid string, api_key string, scan_id string, issue_name string, update_type string) {
	if scan_id != "" {
		//Get all issues for specified website
		issue_ids := GetIssueIds(uid, api_key, scan_id, issue_name)
		fmt.Printf("[+] Issue IDs that will be updated: %s\n", issue_ids)
		l.Printf("[+] Issue IDs that will be updated: %s", issue_ids)

		//Sync for issue tasks
		//wg.Add(len(issue_ids)) -> Deactivated because API cannot reply when the script running fast..

		fmt.Println("[i] Update phase starting...")
		for i := 0; i < len(issue_ids); i++ {
			UpdateIssue(uid, api_key, scan_id, issue_name, issue_ids[i], update_type)
		}
		// wg.Wait()  -> Deactivated because API cannot reply when the script running fast..

	} else {
		// If scan ID is empty, update all of the issues in the account
		fmt.Printf("[i] Processing for all of the %s issues\n", issue_name)
		l.Printf("[i] Processing for all of the %s issues", issue_name)

		//Get all issues for specified account
		issues, total_item_count := GetAllIssueIds(uid, api_key, issue_name)

		fmt.Printf("[i] %d issue found.\n", len(issues))
		l.Printf("[i] %d issue found.", len(issues))

		number_of_updated := 0
		var updated_issue_urls []string
		for i := 0; i < len(issues); i++ {
			//Append issue to updated issue url list
			issue_url := "https://www.netsparkercloud.com/issues/detail/" + issues[i]
			updated_issue_urls = append(updated_issue_urls, issue_url)

			fmt.Printf("%d. Issue ID: %s\n", i+1, issues[i])
			UpdateIssue(uid, api_key, scan_id, issue_name, issues[i], update_type)
			number_of_updated += 1
		}

		fmt.Println()
		fmt.Println("-------- SUMMARY --------")
		fmt.Printf("Number of all issues in the account: %d\n", total_item_count)
		fmt.Printf("Number of %s issues in the account: %d\n", issue_name, len(issues))
		fmt.Printf("Number of updated issues: %d\n", number_of_updated)
		fmt.Println("\nUpdated issue URLs:\n-------------------------")
		for i := 0; i < len(updated_issue_urls); i++ {
			fmt.Println(updated_issue_urls[i])
		}

		l.Println("-------- SUMMARY --------")
		l.Printf("Number of all issues in the account: %d", total_item_count)
		l.Printf("Number of %s issues in the account: %d", issue_name, len(issues))
		l.Printf("Number of updated issues: %d", number_of_updated)

	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func UpdateIssue(uid string, api_key string, scan_id string, issue_name string, issue_id string, update_type string) {

	//defer wg.Done()

	//Update given issue as FP
	var state string
	if update_type == "false_positive" {
		state = "FalsePositive"

	} else if update_type == "accepted_risk" {
		state = "AcceptedRisk"

	} else if update_type == "fixed_unconfirmed" {
		state = "FixedUnconfirmed "

	} else if update_type == "fixed_cant_retest" {
		state = "FixedCantRetest "

	} else {
		fmt.Println("You should provide a valid update_type!")
		return
	}

	data := SingleIssueReqItem{
		IssueID: issue_id,
		State:   state,
	}

	singleIssueReqItemBytes, err := json.Marshal(data)
	checkErr(err)

	body := bytes.NewReader(singleIssueReqItemBytes)

	req, err := http.NewRequest("POST", "https://www.netsparkercloud.com/api/1.0/issues/update", body)
	checkErr(err)

	req.SetBasicAuth(uid, api_key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	fmt.Printf("[+] Issue with ID: %s marked as %s\n", issue_id, state)
	l.Printf("[+] Issue with ID: %s marked as %s", issue_id, state)
	checkErr(err)
	defer resp.Body.Close()
}

func SendRequest(uid string, api_key string, endpoint string) []byte {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(uid, api_key)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// Get response
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

//Return target issue IDs that will be updated and number of all issues
func GetAllIssueIds(uid string, api_key string, issue_name string) ([]string, int) {
	fmt.Printf("[i] Getting %s issues for the user id: %s\n", issue_name, uid)
	l.Printf("[i] Getting %s issues for the user id: %s", issue_name, uid)

	endpoint := "https://www.netsparkercloud.com/api/1.0/issues/allissues?pageSize=200"
	page_count := 0
	current_page := 0
	total_item_count := 0

	//Send request
	body := SendRequest(uid, api_key, endpoint)
	var allIssuesApiResult AllIssuesApiResult

	//Unmarshal body data to to allIssuesApiResult
	json.Unmarshal(body, &allIssuesApiResult)

	//Set the page icount and current page
	page_count = allIssuesApiResult.PageCount
	current_page = allIssuesApiResult.PageNumber
	total_item_count = allIssuesApiResult.TotalItemCount

	fmt.Printf("[i] Page count fetched: %d\n", page_count)
	l.Printf("[i] Page count fetched: %d", page_count)
	fmt.Printf("[i] Current page: %d\n\n", current_page)
	l.Printf("[i] Current page: %d", current_page)

	var issue_object_list []AllIssuesIssue

	fmt.Println("All API results fetching via pagination...")
	//Request to AllIssues endpoint until current page is the last page (including the last one)
	for current_page <= page_count {
		url := endpoint + "&page=" + strconv.Itoa(current_page)
		fmt.Printf("[i] Sending request to %s\n", url)
		l.Printf("[i] Sending request to %s\n", url)
		body := SendRequest(uid, api_key, url)

		//Unmarshal body data to to allIssuesApiResult
		err := json.Unmarshal(body, &allIssuesApiResult)
		if err != nil {
			//Handle error if needed
			fmt.Println("[!] Cannot unmarshall to allIssuesApiResult")
			l.Println("[!] Cannot unmarshall to allIssuesApiResult")
		}

		//Add each object from allIssuesApiResult.List to the issue_object_list
		issue_object_list = append(issue_object_list, allIssuesApiResult.List...)

		//Move to the next page
		current_page = current_page + 1
	}

	//Prepare Issue ID list to return
	var target_issue_ids []string
	for i := 0; i < len(issue_object_list); i++ {
		if issue_object_list[i].Title == issue_name {
			issue_id := issue_object_list[i].ID
			target_issue_ids = append(target_issue_ids, issue_id)
		}
	}

	if len(target_issue_ids) <= 0 {
		fmt.Println("[!] Could not found any issue with that name. Check the issue name is correct or not. Or there is no issue.")
		panic("Warn")
	}

	return target_issue_ids, total_item_count

}

func GetIssueIds(uid string, api_key string, scan_id string, issue_name string) []string {
	//Return all issues for the specified scan
	fmt.Printf("[i] Getting %s issues for the scan %s...\n", issue_name, scan_id)
	l.Printf("[i] Getting %s issues for the scan %s...\n", issue_name, scan_id)

	endpoint := "https://www.netsparkercloud.com/api/1.0/scans/result/" + scan_id

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(uid, api_key)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Get response
	body, _ := ioutil.ReadAll(resp.Body)

	var scanResultIssue ScanResultIssue

	//Unmarshal body data to to scanResultIssue
	json.Unmarshal(body, &scanResultIssue)

	//Get all issues for specified scan
	issues := scanResultIssue

	var target_issue_ids []string

	for i := 0; i < len(issues); i++ {
		if issues[i].Title == issue_name {
			issue_id := strings.Split(issues[i].IssueUrl, "/")[5]
			target_issue_ids = append(target_issue_ids, issue_id)
		}
	}

	if len(target_issue_ids) < 1 {
		fmt.Println("[!] Fetched issue names and provided issue name not matched. Please check the issue_name.")
		l.Printf("[!] Fetched issue name and provided issue name not matched. Provided issue name: %s", issue_name)
		panic("Exiting...")
	}

	fmt.Println("[+] All issues fetched.")
	l.Println("[+] All issues fetched.")

	return target_issue_ids

}
