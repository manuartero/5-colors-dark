#!/bin/bash
# Script to update a kraken service definition, commit it and run kraken push

if [ $# -lt 2 ]; then
  echo "This script updates a kraken service yaml file with the provide docker tag,"
  echo "then it commits it and calls kraken to process and push it"
  echo
  echo "Syntax: ./update-service <service> <tag> [<params_for_kraken_push>]"
  exit 1
fi

service=$1
shift
tag=$1
shift

cd `dirname $0`

service_yaml=definitions-kraken/services/$service.yaml

if [ ! -r $service_yaml ]; then
  echo "ERROR: Couldn't find service definition file $service_yaml"
  exit 1
fi

num_tags=`grep "\bdocker.tag\s*:" $service_yaml | wc -l`
if [ $num_tags -eq 0 ]; then
  echo "ERROR: Couldn't find any 'docker.tag' in $service_yaml"
  exit 1
elif [ $num_tags -gt 1 ]; then
  echo "ERROR: More than one 'docker.tag' found in $service_yaml"
  exit 1
fi

echo ">>>> Updating local branch to ensure we can commit and push..."
current_branch=`git rev-parse --abbrev-ref HEAD`
git fetch -q
git merge --ff-only origin/$current_branch
if [ $? -ne 0 ]; then
  echo
  echo "ERROR: Couldn't pull correctly from the git repository. Manual intervention required"
  exit 1
fi
echo

sed -i -e "s/^\([[:space:]]*docker.tag[[:space:]]*:[[:space:]]*\"*\)\([^\"]*\)\(\"*[[:space:]]*\)$/\1$tag\3/" $service_yaml

echo ">>>> Resulting diff after applying the tag change:"
git diff $service_yaml
read ins del file <<< $(git diff --numstat $service_yaml)

if [ -z $file ]; then
  echo "ERROR: No modifications made to $service_yaml (probably same version applied)"
  exit 1
elif [ $ins -ne 1 -o $del -ne 1 ]; then
  echo "ERROR: Unexpected affected lines in $service_yaml. Please restore the file manually"
  exit 1
fi

echo
echo ">>>> Running git commit"
git commit -m "Update service $service to version $tag" $service_yaml
if [ $? -ne 0 ]; then
  echo "ERROR: Failed to commit changes. Please go on manually"
  exit 1
fi

echo
echo ">>>> Calling kraken push"
./kraken push "$@"
