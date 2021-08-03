#!/bin/bash

echo
echo "Passed in vars"
echo "INPUT_SCAN_TYPE: $INPUT_SCAN_TYPE"
echo "INPUT_COMMENT_ON_PR: $INPUT_COMMENT_ON_PR"
echo "INPUT_TERRAFORM_WORKING_DIR: $INPUT_TERRAFORM_WORKING_DIR"
echo "INPUT_TFSEC_EXCLUDE: $INPUT_TFSEC_EXCLUDE"
echo "INPUT_TFSEC_VERSION: $INPUT_TFSEC_VERSION"
echo "INPUT_TFSEC_OUTPUT_FORMAT: $INPUT_TFSEC_OUTPUT_FORMAT"
echo "INPUT_TFSEC_OUTPUT_FILE: $INPUT_TFSEC_OUTPUT_FILE"
echo "INPUT_CHECKOV_EXCLUDE: $INPUT_CHECKOV_EXCLUDE"
echo "INPUT_TFLINT_EXCLUDE: $INPUT_TFLINT_EXCLUDE"
echo
# install tfsec from GitHub (taken from README.md)
if [[ -n "$INPUT_TFSEC_VERSION" ]]; then
  env GO111MODULE=on go install github.com/aquasecurity/tfsec/cmd/tfsec@"${INPUT_TFSEC_VERSION}"
else
  env GO111MODULE=on go get -u github.com/aquasecurity/tfsec/cmd/tfsec
fi

line_break() {
  echo
  echo "*****************************"
  echo
}

declare -i tfsec_exitcode=0
declare -i checkov_exitcode=0
declare -i tflint_exitcode=0
declare -i tfinit_exitcode=0

# Identify which Terraform folders have changes and need scanning
tf_folders_with_changes=`git diff --no-commit-id --name-only -r @^ | awk '{print $1}' | grep '.tf' | sed 's#/[^/]*$##' | uniq`
echo
echo "TF folders with changes"
echo $tf_folders_with_changes

# Get a list of all terraform folders in the repo
all_tf_folders=`find . -type f -name '*.tf' | sed 's#/[^/]*$##' | sed 's/.\///'| sort | uniq`
echo
echo "All TF folders"
echo $all_tf_folders

run_tfsec(){
  line_break
  echo "TFSEC will check the following folders:"
  echo $1
  directories=($1)
  for directory in ${directories[@]}
  do
    line_break
    echo "Running TFSEC in ${directory}"
    terraform_working_dir="/github/workspace/${directory}"
    if [[ -n "$INPUT_TFSEC_EXCLUDE" ]]; then
      echo "Excluding the following checks: ${INPUT_TFSEC_EXCLUDE}"
      /go/bin/tfsec ${terraform_working_dir} --no-colour -e "${INPUT_TFSEC_EXCLUDE}" ${INPUT_TFSEC_OUTPUT_FORMAT:+ -f "$INPUT_TFSEC_OUTPUT_FORMAT"} ${INPUT_TFSEC_OUTPUT_FILE:+ --out "$INPUT_TFSEC_OUTPUT_FILE"} 2>&1
    else
      /go/bin/tfsec ${terraform_working_dir} --no-colour ${INPUT_TFSEC_OUTPUT_FORMAT:+ -f "$INPUT_TFSEC_OUTPUT_FORMAT"} ${INPUT_TFSEC_OUTPUT_FILE:+ --out "$INPUT_TFSEC_OUTPUT_FILE"} 2>&1
    fi
    tfsec_exitcode+=$?
    echo "tfsec_exitcode=${tfsec_exitcode}"
  done
  return $tfsec_exitcode
}

run_checkov(){
  line_break
  echo "Checkov will check the following folders:"
  echo $1
  directories=($1)
  for directory in ${directories[@]}
  do
    line_break
    echo "Running Checkov in ${directory}"
    terraform_working_dir="/github/workspace/${directory}"
    if [[ -n "$INPUT_CHECKOV_EXCLUDE" ]]; then
      echo "Excluding the following checks: ${INPUT_CHECKOV_EXCLUDE}"
      checkov --quiet -d $terraform_working_dir --skip-check ${INPUT_CHECKOV_EXCLUDE} 2>&1
    else
      checkov --quiet -d $terraform_working_dir 2>&1
    fi
    checkov_exitcode+=$?
    echo "checkov_exitcode=${checkov_exitcode}"
  done
  return $checkov_exitcode
}

run_tflint(){
  line_break
  echo "tflint will check the following folders:"
  echo $1
  directories=($1)
  for directory in ${directories[@]}
  do
    line_break
    echo "Running tflint in ${directory}"
    terraform_working_dir="/github/workspace/${directory}"
    if [[ "${directory}" != *"templates"* ]]
    then
      if [[ -n "$INPUT_TFLINT_EXCLUDE" ]]; then
        echo "Excluding the following checks: ${INPUT_TFLINT_EXCLUDE}"
        tflint --disable-rule="${INPUT_TFLINT_EXCLUDE}" ${terraform_working_dir} 2>&1
      else
        tflint ${terraform_working_dir} 2>&1
      fi
    else 
      echo "Skipping folder as path name contains *templates*"
    fi
    tflint_exitcode+=$?
    echo "tflint_exitcode=${tflint_exitcode}"
  done
  return $tflint_exitcode
}

case ${INPUT_SCAN_TYPE} in

  full)
    line_break
    echo "Starting full scan"
    TFSEC_OUTPUT=$(run_tfsec "${all_tf_folders}")
    tfsec_exitcode=$?
    wait
    CHECKOV_OUTPUT=$(run_checkov "${all_tf_folders}")
    checkov_exitcode=$?
    wait
    TFLINT_OUTPUT=$(run_tflint "${all_tf_folders}")
    tflint_exitcode=$?
    wait
    ;;

  changed)
    line_break
    echo "Starting scan of changed folders"
    TFSEC_OUTPUT=$(run_tfsec "${tf_folders_with_changes}")
    tfsec_exitcode=$?
    wait
    CHECKOV_OUTPUT=$(run_checkov "${tf_folders_with_changes}")
    checkov_exitcode=$?
    wait
    TFLINT_OUTPUT=$(run_tflint "${tf_folders_with_changes}")
    tflint_exitcode=$?
    wait
    ;;
  *)
    line_break
    echo "Starting single folder scan"
    TFSEC_OUTPUT=$(run_tfsec "${INPUT_TERRAFORM_WORKING_DIR}")
    tfsec_exitcode=$?
    wait
    CHECKOV_OUTPUT=$(run_checkov "${INPUT_TERRAFORM_WORKING_DIR}")
    checkov_exitcode=$?
    wait
    TFLINT_OUTPUT=$(run_tflint "${INPUT_TERRAFORM_WORKING_DIR}")
    tflint_exitcode=$?
    wait
    ;;
esac

if [ $tfsec_exitcode -eq 0 ]; then
  TFSEC_STATUS="Success"
else
  TFSEC_STATUS="Failed"
fi

if [ $checkov_exitcode -eq 0 ]; then
  CHECKOV_STATUS="Success"
else
  CHECKOV_STATUS="Failed"
fi

if [ $tflint_exitcode -eq 0 ]; then
  TFLINT_STATUS="Success"
else
  TFLINT_STATUS="Failed"
fi

# Print output.
line_break
echo "${TFSEC_OUTPUT}"
echo "${CHECKOV_OUTPUT}"
echo "${TFLINT_OUTPUT}"

# Comment on the pull request if necessary.
if [ "${INPUT_COMMENT_ON_PR}" == "1" ] || [ "${INPUT_COMMENT_ON_PR}" == "true" ]; then
  COMMENT=1
else
  COMMENT=0
fi

if [ "${GITHUB_EVENT_NAME}" == "pull_request" ] && [ -n "${GITHUB_TOKEN}" ] && [ "${COMMENT}" == "1" ] ; then
    COMMENT="#### \`TFSEC Scan\` ${TFSEC_STATUS}
<details><summary>Show Output</summary>

\`\`\`hcl
${TFSEC_OUTPUT}
\`\`\`

</details>

#### \`Checkov Scan\` ${CHECKOV_STATUS}
<details><summary>Show Output</summary>

\`\`\`hcl
${CHECKOV_OUTPUT}
\`\`\`

</details>

#### \`CTFLint Scan\` ${TFLINT_STATUS}
<details><summary>Show Output</summary>

\`\`\`hcl
${TFLINT_OUTPUT}
\`\`\`

</details>"

  PAYLOAD=$(echo "${COMMENT}" | jq -R --slurp '{body: .}')
  URL=$(jq -r .pull_request.comments_url "${GITHUB_EVENT_PATH}")
  echo "${PAYLOAD}" | curl -s -S -H "Authorization: token ${GITHUB_TOKEN}" --header "Content-Type: application/json" --data @- "${URL}" > /dev/null
fi

line_break
echo "Total of TFSEC exit codes: $tfsec_exitcode"
echo "Total of Checkov exit codes: $checkov_exitcode"
echo "Total of tflint exit codes: $tflint_exitcode"

if [ $tfsec_exitcode -gt 0 ] || [ $checkov_exitcode -gt 0 ] || [ $tflint_exitcode -gt 0 ];then
  echo "Exiting with error(s)"  
  exit 1
else
  echo "Exiting with no error"  
  exit 0
fi
