
# Invicti Enterprise Bulk Issue Updater

This tool was created for a customer request. You can update all of the issues in a scan according to the issue name with your UID and API Key.


# Invicti Enterprise Bulk Issue Updater

This tool was created for a customer request. You can update all of the issues in a scan according to the issue name with your UID and API Key.


## Usage/Examples

```bash
.\IssueUpdater.exe -h
USAGE:
    [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --uid value         User ID
   --api_key value     API Key
   --scan_id value     [Optional] Scan ID to update issues.
   --issue_name value  Name of the issue to be updated
   --update value      Update type. USAGE: --update false_positive. Available options: {false_positive, accepted_risk, fixed_unconfirmed, fixed_cant_retest} Note: You should use proper 'fixed' according to the issue type. (Retestable or not)
   --help, -h          show help
```

- uid: User ID. (Found on IE > API Settings )
- scan_id: ID of scan that you want to update issues. This parameter is optional, if you don’t set the tool will be update all of the specified issues in the user account.
- api_key: IE API Key
- issue_name: Name of the issue that you want to be updated.
- update: Update status type.  Available Options: {false_positive, accepted_risk, fixed_unconfirmed, fixed_cant_retest}


Update all of the Missing X-Frame-Options Header vulnerabilities in the user account as False Positive:
```bash
IssueUpdater.exe --uid <USER-ID>--api_key <API-KEY> --issue_name "Missing X-Frame-Options Header" --update false_positive
```
