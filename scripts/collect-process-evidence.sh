#!/usr/bin/env bash
set -Eeuo pipefail

repository=${1:?repository is required}
phase=${2:?phase is required}
output=${3:?output is required}

commits=$(gh api "repos/$repository/commits?sha=main&per_page=100")
main_head=$(jq -r '.[0].sha' <<<"$commits")
bootstrap_root=$(jq -r '[.[] | select((.parents | length) == 0) | .sha] | last // empty' <<<"$commits")
root_parent_count=$(jq --arg root "$bootstrap_root" -r '[.[] | select(.sha == $root) | (.parents | length)] | .[0] // -1' <<<"$commits")
post_bootstrap_direct_main=0
current="$main_head"
while [ "$current" != "$bootstrap_root" ]; do
	commit=$(gh api "repos/$repository/commits/$current")
	parent_count=$(jq '.parents | length' <<<"$commit")
	if [ "$parent_count" -eq 1 ]; then
		post_bootstrap_direct_main=$((post_bootstrap_direct_main + 1))
	fi
	current=$(jq -r '.parents[0].sha // empty' <<<"$commit")
	if [ -z "$current" ]; then
		echo "main first-parent chain did not reach bootstrap root" >&2
		exit 1
	fi
done
pull_requests=$(gh pr list --repo "$repository" --state all --limit 20 --json number,url,state,isDraft,mergedAt,baseRefName,headRefName,headRefOid,mergeCommit)

jq -n \
	--arg schema "gooo/evolution-trial/process-evidence/v1" \
	--arg repository "$repository" --arg phase "$phase" --arg main_head "$main_head" \
	--arg bootstrap_root "$bootstrap_root" --argjson root_parent_count "$root_parent_count" \
	--argjson post_bootstrap_direct_main "$post_bootstrap_direct_main" \
	--argjson pull_requests "$pull_requests" \
	'{schema:$schema,repository:$repository,phase:$phase,main_head:$main_head,bootstrap_root:$bootstrap_root,bootstrap_root_parent_count:$root_parent_count,bootstrap_direct_main:1,post_bootstrap_direct_main:$post_bootstrap_direct_main,pull_requests:$pull_requests,repository_writes_by_experiment:0,upstream_writes_by_experiment:0,local_test_executions:0,utility_global_core:{state:"UNKNOWN",status:"NOT_MADE"},exact:($root_parent_count == 0 and $post_bootstrap_direct_main == 0),next_operation:(if $phase == "release" then "VERIFY_IMMUTABLE_RELEASE_AND_ASSET_DIGESTS" else "MERGE_ONE_FEATURE_PR_THEN_PUBLISH_IMMUTABLE_RELEASE" end)}' > "$output"
