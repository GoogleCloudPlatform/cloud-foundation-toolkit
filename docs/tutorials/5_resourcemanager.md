#  Installation

## Introduction

<walkthrough-tutorial-duration duration="10"></walkthrough-tutorial-duration>

This tutorial explains how to set up a hierarchy for managing resources in GCP.

In conjunction with Cloud IAM, the Resource Manager (RM) service provides the foundation for security governance in GCP.

Resource manager organizes GCP resources hierarchically.

The tutorial assume you have Organizations Admin privileges to create folders.

<!-- TODO: add overview/diagram -->

## Open Organization Resource Manager

To get started, open the [Manage resources](https://console.cloud.google.com/cloud-resource-manager) page under the IAM section in the menu.

<walkthrough-menu-navigation sectionid="IAM_ADMIN_SECTION"></walkthrough-menu-navigation>

## Create a folder

Ensure that your organization domain is shown in the top left of the page.
If not, click the <walkthrough-spotlight-pointer cssSelector="cfc-purview-picker-org">dropdown</walkthrough-spotlight-pointer> to select it.

### Create folder

Click the <walkthrough-spotlight-pointer cssSelector="#create-folder-button">Create Folder</walkthrough-spotlight-pointer> button at the top of the page.

### Name folder

Give your folder a <walkthrough-spotlight-pointer cssSelector="label[for='folderName']">name</walkthrough-spotlight-pointer> by typing it in form.

Keep in mind these restrictions when creating the folder:
- The name must be 30 characters or less.
- The name must be distinct from all other folders that share its parent.

**For your first folder, you might create it with the name `sandbox`.**

After entering the name, click <walkthrough-spotlight-pointer cssSelector="#createFolder">Create</walkthrough-spotlight-pointer>.

## Create additional folders
To complete your folder hierarchy, continue the previous step for each folder you wish to create.

For example, you can start by creating folders for each environment:

- `sandbox`
- `dev`
- `test`
- `production`

### Nesting Folders
You might want to create additional folders to represent your organization hierarchy (such as folders for each business unit).

To create a nested folder, change the **Destination** in the folder creation form to the parent folder.

## Move projects into folders
Now that you have folders created, you can begin [moving your existing projects into folders](https://cloud.google.com/resource-manager/docs/creating-managing-folders#moving_a_project_into_a_folder).

### Select projects
From the **Manage resources** page, select the <walkthrough-spotlight-pointer cssSelector=".mat-pseudo-checkbox:not(.cfctest-table-select-all-checkbox)">checkbox</walkthrough-spotlight-pointer>
next to each project you would like to move into a folder.


### Move projects
Once you have selected the projects, click the <walkthrough-spotlight-pointer cssSelector="#move-button">Move</walkthrough-spotlight-pointer>
button at the top of the page.

In the dialogue, click the folder you would like to move the projects into then click Select.

Repeat this process for all projects you would like to move.
