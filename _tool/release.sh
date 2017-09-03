#/bin/bash -eu
SCRIPT_DIR=$(cd $(dirname $0);pwd)
REGEX='v[0-9].[0-9].[0-9]'

onerror()
{
    status=$?
    script=$0
    line=$1
    shift

    args=
    for i in "$@"; do 
        args+="\"$i\" "
    done

    echo ""
    echo "------------------------------------------------------------"
    echo "Error occured on $script [Line $line]: Status $status"
    echo ""
    echo "PID: $$"
    echo "User: $USER"
    echo "Current directory: $PWD"
    echo "Command line: $script $args"
    echo "------------------------------------------------------------"
    echo ""
}

begintrap()
{
    set -e
    trap 'onerror $LINENO "$@"' ERR
}

validate()
{
    TAG_NAME=$1
    echo ${TAG_NAME}
    echo ${REGEX}
    if [[ ${TAG_NAME} =~ ${REGEX} ]]; then
        echo "Input tag name: ${TAG_NAME}"    
        return 0
    else
        echo 'Invalid tag name.' >&2
        return 2
    fi
}

yes_or_no()
{
    PS3="Continue? "
    while true; do
        echo "Type 1 or 2."
        select answer in yes no; do
            case $answer in
                yes)
                    echo -e "tyeped yes.\n"
                    return 0
                    ;;
                no)
                    echo -e "tyeped no.\n"
                    return 1
                    ;;
                *)
                    echo -e "cannot understand your answer.\n"
                    ;;
            esac
        done
    done
}

release()
{
    BRANCH_NAME="bump-version-${TAG_NAME//./_}"
    git checkout -b ${BRANCH_NAME}

    # Replace version
    echo ${TAG_NAME} > ${SCRIPT_DIR}/version
    sed -i '' -e "s/${REGEX}/${TAG_NAME}/g" ${SCRIPT_DIR}/../mog.go

    # generate CHANGELOG.md
    github-changes -o szyn -r mog --only-pulls -n ${TAG_NAME}
    git add ${SCRIPT_DIR}/../CHANGELOG.md
    git commit -m "Update CHANGELOG.md"
    
    git add .
    git commit -m "Bump version to ${TAG_NAME}"
    git tag -a ${TAG_NAME} -m "Bump version to ${TAG_NAME}"
}

begintrap
validate $1
yes_or_no
release