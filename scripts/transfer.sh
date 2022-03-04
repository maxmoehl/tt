#!/bin/zsh

for i in {0..254}; do
  data=$(jq -cM ".[${i}]" ./data.json | sed -e 's/ID/id/' -e 's/Start/start/' -e 's/Stop/stop/' -e 's/Project/project/' -e 's/Task/task/' -e 's/Tags/tags/')
  if [[ $(echo ${data} | jq '.task') == '""' ]]; then
    data=$(echo "${data}" | jq -cM 'del(.task)')
  fi
  echo "INSERT INTO timers VALUES('$(echo "${data}" | jq -r '.id')','${data}'); --${i}"
done