#!/usr/bin/env bash
set -Eeuo pipefail
set -x

root=${1:?repository root is required}
work=${2:?upstream work directory is required}
mkdir -p "$work/assets" "$work/observed"

lock="$root/contracts/upstream-lock-v1.json"

while IFS=$'\t' read -r producer tag release_id immutable tag_object target_commit; do
	if [ -z "$producer" ]; then
		continue
	fi
	repository=${producer#github.com/}
	release_json="$work/observed/${repository##*/}-${tag}.json"
	gh api "repos/$repository/releases/tags/$tag" > "$release_json"
	jq -e --argjson id "$release_id" --arg tag "$tag" --argjson immutable "$immutable" \
		'.id == $id and .tag_name == $tag and .immutable == $immutable and .draft == false and .prerelease == false' "$release_json" >/dev/null
	tag_ref=$(gh api "repos/$repository/git/ref/tags/$tag")
	actual_object=$(jq -r '.object.sha' <<<"$tag_ref")
	actual_type=$(jq -r '.object.type' <<<"$tag_ref")
	tag_object_json="$work/observed/${repository##*/}-${tag}-tag.json"
	gh api "repos/$repository/git/tags/$actual_object" > "$tag_object_json"
	actual_target=$(jq -r '.object.sha' "$tag_object_json")
	actual_target_type=$(jq -r '.object.type' "$tag_object_json")
	test "$actual_object" = "$tag_object"
	test "$actual_type" = "tag"
	test "$actual_target" = "$target_commit"
	test "$actual_target_type" = "commit"

	jq -r --arg repo "$repository" --arg tag "$tag" '.assets[] | [$repo,$tag,.id,.name,.size,.digest,.browser_download_url] | @tsv' "$release_json" |
	while IFS=$'\t' read -r asset_repo asset_tag asset_id asset_name asset_size asset_digest asset_url; do
		locked=$(jq -r --arg producer "github.com/$asset_repo" --arg tag "$asset_tag" --argjson id "$asset_id" \
			'.inputs[] | select(.producer == $producer and .release == $tag) | .assets[] | select(.id == $id) | [.name,.size_bytes,.sha256] | @tsv' "$lock")
		if [ -z "$locked" ]; then
			continue
		fi
		IFS=$'\t' read -r locked_name locked_size locked_sha <<<"$locked"
		test "$asset_name" = "$locked_name"
		test "$asset_size" = "$locked_size"
		test "$asset_digest" = "sha256:$locked_sha"
		curl --fail --silent --show-error --location --output "$work/assets/$asset_name" "$asset_url"
		actual_size=$(stat -c '%s' "$work/assets/$asset_name")
		actual_sha=$(sha256sum "$work/assets/$asset_name" | awk '{print $1}')
		test "$actual_size" = "$locked_size"
		test "$actual_sha" = "$locked_sha"
	done

	gh api "repos/$repository/actions/runs/$(jq -r --arg repo "github.com/$repository" --arg tag "$tag" '.inputs[] | select(.producer == $repo and .release == $tag) | .release_run_id' "$lock")/jobs?per_page=100" > "$work/observed/${repository##*/}-${tag}-jobs.json"
done < <(jq -r '.inputs[] | [.producer,.release,.release_id,.immutable,.tag_object,.target_commit] | @tsv' "$lock")

jq -S . "$lock" > "$work/observed/upstream-lock.json"
printf '%s\n' '{"schema":"gooo/evolution-trial/upstream-fetch/v1","decision":"CLOSED","writes":{"repository":0,"upstream":0},"tag_and_asset_digests_verified":true}' > "$work/observed/fetch-report.json"
