/**
 * @type {import('semantic-release').GlobalConfig}
 */
export default {
    plugins: [
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        "@semantic-release/github",
        [
            "@semantic-release/exec",
            {
                publishCmd:
                    'echo "${nextRelease.notes}" > /tmp/release-notes.md; goreleaser release --release-notes /tmp/release-notes.md --clean',
            },
        ],
    ],
};
