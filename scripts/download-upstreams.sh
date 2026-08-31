#!/usr/bin/env bash
set -Eeuo pipefail

root=${1:?repository root is required}
work=${2:?upstream work directory is required}
mkdir -p "$work/assets" "$work/observed"

lock="$root/contracts/upstream-lock-v1.json"

input_count=$(jq '.inputs | length' "$lock")
for ((input_index = 0; input_index < input_count; input_index++)); do
	producer=$(jq -r ".inputs[$input_index].producer" "$lock")
	tag=$(jq -r ".inputs[$input_index].release" "$lock")
	release_id=$(jq -r ".inputs[$input_index].release_id" "$lock")
	immutable=$(jq -r ".inputs[$input_index].immutable" "$lock")
	tag_object=$(jq -r ".inputs[$input_index].tag_object" "$lock")
	target_commit=$(jq -r ".inputs[$input_index].target_commit" "$lock")
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

	asset_count=$(jq ".inputs[$input_index].assets | length" "$lock")
	for ((asset_index = 0; asset_index < asset_count; asset_index++)); do
		asset_id=$(jq -r ".inputs[$input_index].assets[$asset_index].id" "$lock")
		locked_name=$(jq -r ".inputs[$input_index].assets[$asset_index].name" "$lock")
		locked_size=$(jq -r ".inputs[$input_index].assets[$asset_index].size_bytes" "$lock")
		locked_sha=$(jq -r ".inputs[$input_index].assets[$asset_index].sha256" "$lock")
		asset_url=$(jq -er --argjson id "$asset_id" '.assets[] | select(.id == $id) | .browser_download_url' "$release_json")
		asset_name=$(jq -er --argjson id "$asset_id" '.assets[] | select(.id == $id) | .name' "$release_json")
		asset_size=$(jq -er --argjson id "$asset_id" '.assets[] | select(.id == $id) | .size' "$release_json")
		asset_digest=$(jq -er --argjson id "$asset_id" '.assets[] | select(.id == $id) | .digest' "$release_json")
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
done

jq -S . "$lock" > "$work/observed/upstream-lock.json"
printf '%s\n' '{"schema":"gooo/evolution-trial/upstream-fetch/v1","decision":"CLOSED","writes":{"repository":0,"upstream":0},"tag_and_asset_digests_verified":true}' > "$work/observed/fetch-report.json"
