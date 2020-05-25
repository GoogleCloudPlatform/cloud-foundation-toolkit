import logging
import os
import argparse
from github import Github, GithubException

"""Creates new comment or updates existing comment on a PR using a PAT token"""

def create_update_comment(token, org, repo_name, pr_number, comment_body):
    """Creates or updates existing comment on a PR"""
    # auth with GH token
    gh = Github(token)
    # ref to repo
    github_repo = gh.get_repo(f'{org}/{repo_name}')
    # get current user id of token owner
    current_bot_user = gh.get_user().id
    logging.info(f'Current bot user: {current_bot_user}')
    # try to get the PR
    try:
        pr = github_repo.get_pull(pr_number)
    except GithubException as e:
        logging.info(f'Unable to get PR: {pr_number} in {org}/{repo_name}')
        raise
    # get all comments in the pr (doesnt capture review comments)
    pr_comments = pr.get_issue_comments()
    # check if bot has already commented
    existing_comments = [
        pr_comment for pr_comment in pr_comments
        if pr_comment.user.id == current_bot_user
        ]
    if not existing_comments:
        # add a comment
        comment = pr.create_issue_comment(comment_body)
        logging.info(f'Added new comment: {comment}')
    else:
        # edit existing comment
        existing_comments[0].edit(comment_body)
        logging.info(f'Edited existing comment: {existing_comments[0]}')

def parse_args():
    parser = argparse.ArgumentParser(description='Add/edit comments to PRs')
    parser.add_argument('-o', '--org', default='terraform-google-modules', help='Github organization, defaults to cft modules repo', action='store')
    parser.add_argument('-r', '--repo', help='Github repo name', action='store')
    parser.add_argument('-p', '--pr', help='Github PR number to add/edit comment', type=int, action='store')
    parser.add_argument('-c', '--comment', help='Comment body to create or update ', action='store')
    return parser.parse_args()

if __name__ == "__main__":
    # setup logging
    logging.basicConfig(level=logging.INFO)
    # check if GITHUB_PAT_TOKEN token is set in env
    if os.environ.get('GITHUB_PAT_TOKEN') is None:
        raise RuntimeError('Unable to find GITHUB_PAT_TOKEN token in env')
    gh_token = os.environ.get('GITHUB_PAT_TOKEN')
    # parse args
    parser = parse_args()
    create_update_comment(gh_token, parser.org, parser.repo, parser.pr, parser.comment)
