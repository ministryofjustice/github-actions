from datetime import datetime

import pytest
import semantic_version
from filter_advisories import (
    ValidationError,
    advisory_to_slack_block,
    filter_advisories,
    format_slack_output,
    parse_vulnerabilities,
)


@pytest.fixture
def advisories():
    return [
        {
            "ghsa_id": "GHSA-xxxx-xxxx-xxxx",
            "severity": "medium",
            "published_at": "2023-06-01T00:00:00Z",
            "summary": "This is a test advisory",
            "html_url": "https://github.com/owner/target-repo/security/advisories/GHSA-xxxx-xxxx-xxxx",
            "vulnerabilities": [
                {
                    "vulnerable_version_range": "< 0.12.0",
                    "patched_versions": ">= 0.12.0",
                }
            ],
        },
        {
            "ghsa_id": "GHSA-yyyy-yyyy-yyyy",
            "severity": "high",
            "published_at": "2023-07-01T00:00:00Z",
            "summary": "This is another test advisory",
            "html_url": "https://github.com/owner/target-repo/security/advisories/GHSA-yyyy-yyyy-yyyy",
            "vulnerabilities": [
                {
                    "vulnerable_version_range": "ALL",
                    "patched_versions": "",
                }
            ],
        },
        {
            "ghsa_id": "GHSA-zzzz-zzzz-zzzz",
            "severity": "critical",
            "published_at": "2023-08-01T00:00:00Z",
            "summary": "This is yet another test advisory",
            "html_url": "https://github.com/owner/target-repo/security/advisories/GHSA-zzzz-zzzz-zzzz",
            "vulnerabilities": [
                {
                    "vulnerable_version_range": "<= 0.10.1",
                    "patched_versions": "",
                },
                {
                    "vulnerable_version_range": "= 0.11.1",
                    "patched_versions": "",
                },
            ],
        },
    ]


class TestFilterAdvisories:

    @pytest.mark.parametrize(
        "minimal_version,last_run_datetime,expected_count,expected_ids",
        [
            (
                "0.11.0",
                "2023-05-31T00:00:00+00:00",
                2,
                ["GHSA-xxxx-xxxx-xxxx", "GHSA-yyyy-yyyy-yyyy"],
            ),
            (
                "0.12.0",
                "2023-05-31T00:00:00+00:00",
                1,
                ["GHSA-yyyy-yyyy-yyyy"],
            ),
            (
                "0.11.1",
                "2023-05-31T00:00:00+00:00",
                3,
                ["GHSA-xxxx-xxxx-xxxx", "GHSA-yyyy-yyyy-yyyy", "GHSA-zzzz-zzzz-zzzz"],
            ),
            (
                "0.11.1",
                "2023-07-02T00:00:00+00:00",
                1,
                ["GHSA-zzzz-zzzz-zzzz"],
            ),
        ],
    )
    def test_filter_advisories(
        self,
        advisories,
        minimal_version,
        last_run_datetime,
        expected_count,
        expected_ids,
    ):
        minimal_version = semantic_version.Version(minimal_version)
        last_run_datetime = datetime.fromisoformat(last_run_datetime)

        filtered = filter_advisories(advisories, minimal_version, last_run_datetime)
        assert len(filtered) == expected_count
        assert [adv["ghsa_id"] for adv in filtered] == expected_ids

    def test_filter_advisories_validation_error(self):
        minimal_version = semantic_version.Version("0.11.0")
        last_run_datetime = datetime.fromisoformat("2023-06-01T00:00:00+00:00")
        broken_advisories = [{"ghsa_id": "GHSA-zzzz-zzzz-zzzz", "published_at": ""}]

        with pytest.raises(ValidationError):
            filter_advisories(broken_advisories, minimal_version, last_run_datetime)


class TestParseVulnerabilities:

    @pytest.mark.parametrize(
        "minimal_version,advisory_id,expected_ids",
        [
            ("0.11.0", "GHSA-xxxx-xxxx-xxxx", ["GHSA-xxxx-xxxx-xxxx"]),
            ("0.12.0", "GHSA-xxxx-xxxx-xxxx", []),
            ("0.11.0", "GHSA-yyyy-yyyy-yyyy", ["GHSA-yyyy-yyyy-yyyy"]),
            ("0.11.0", "GHSA-zzzz-zzzz-zzzz", []),
            ("0.11.1", "GHSA-zzzz-zzzz-zzzz", ["GHSA-zzzz-zzzz-zzzz"]),
            ("0.10.0", "GHSA-zzzz-zzzz-zzzz", ["GHSA-zzzz-zzzz-zzzz"]),
        ],
    )
    def test_parse_vulnerabilities(
        self, advisories, minimal_version, advisory_id, expected_ids
    ):
        minimal_version = semantic_version.Version(minimal_version)
        filtered = []

        advisory = next(adv for adv in advisories if adv["ghsa_id"] == advisory_id)
        parse_vulnerabilities(advisory, filtered, minimal_version)

        assert len(filtered) == len(expected_ids)
        assert [adv["ghsa_id"] for adv in filtered] == expected_ids

    @pytest.mark.parametrize(
        "advisory",
        [
            {
                "ghsa_id": "GHSA-zzzz-zzzz-zzzz",
                "severity": "low",
                "published_at": "2023-06-03T00:00:00Z",
                "summary": "Missing vulnerable version range",
                "html_url": "https://github.com/owner/target-repo/security/advisories/GHSA-zzzz-zzzz-zzzz",
                "vulnerabilities": [{}],
            }
        ],
    )
    def test_parse_vulnerabilities_validation_error(self, advisory):
        minimal_version = semantic_version.Version("0.11.0")
        filtered = []

        with pytest.raises(ValidationError):
            parse_vulnerabilities(advisory, filtered, minimal_version)

    def test_parse_vulnerabilities_invalid_version(self):
        minimal_version = semantic_version.Version("0.11.0")
        filtered = []

        advisory = {
            "ghsa_id": "GHSA-aaaa-aaaa-aaaa",
            "severity": "low",
            "published_at": "2023-06-04T00:00:00Z",
            "summary": "Invalid version range",
            "html_url": "https://github.com/owner/target-repo/security/advisories/GHSA-aaaa-aaaa-aaaa",
            "vulnerabilities": [{"vulnerable_version_range": "invalid_version"}],
        }
        parse_vulnerabilities(advisory, filtered, minimal_version)
        assert len(filtered) == 1
        assert filtered[0]["ghsa_id"] == advisory["ghsa_id"]


class TestAdvisoryToSlackBlock:

    @pytest.mark.parametrize(
        "advisory_id,expected_high_severity",
        [
            ("GHSA-xxxx-xxxx-xxxx", False),
            ("GHSA-yyyy-yyyy-yyyy", True),
        ],
    )
    def test_advisory_to_slack_block(
        self, advisories, advisory_id, expected_high_severity
    ):
        advisory = next(adv for adv in advisories if adv["ghsa_id"] == advisory_id)
        block, high_severity = advisory_to_slack_block(advisory)

        assert block["type"] == "section"
        assert block["text"]["type"] == "mrkdwn"
        assert f"*ID:* {advisory['ghsa_id']}" in block["text"]["text"]
        assert high_severity == expected_high_severity


class TestFormatSlackOutput:

    @pytest.mark.parametrize(
        "filtered_advisory_ids,expected_block_count,contains_high_severity",
        [
            (
                ["GHSA-xxxx-xxxx-xxxx", "GHSA-yyyy-yyyy-yyyy"],
                5,  # 5 blocks: 1x warning, 2x div, 2x advisory
                True,
            ),
            (["GHSA-xxxx-xxxx-xxxx"], 1, False),  # 1 blocks: 1x advisory
            (
                ["GHSA-yyyy-yyyy-yyyy"],
                3,  # 3 blocks: 1x warning, 1x div, 1x advisory
                True,
            ),
        ],
    )
    def test_format_slack_output(
        self,
        advisories,
        filtered_advisory_ids,
        expected_block_count,
        contains_high_severity,
    ):
        filtered_advisories = [
            adv for adv in advisories if adv["ghsa_id"] in filtered_advisory_ids
        ]
        result = format_slack_output(filtered_advisories)
        assert len(result["blocks"]) == expected_block_count

        if contains_high_severity:
            assert result["blocks"][0]["type"] == "section"
            assert (
                ":warning: contains _high severity_ advisory :warning:"
                in result["blocks"][0]["text"]["text"]
            )
        else:
            assert result["blocks"][0]["type"] == "section"
            assert (
                f"*ID:* {filtered_advisory_ids[0]}"
                in result["blocks"][0]["text"]["text"]
            )
