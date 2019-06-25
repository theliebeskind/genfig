#!/bin/bash

exit_script() {
    trap - SIGINT SIGTERM # clear the trap
    kill -- -$$ # Sends SIGTERM to child/sub processes
}

trap exit_script SIGINT SIGTERM

watch_qtpl() {
  go generate
  while true; do
    fswatch -1 ./**/*.qtpl
    go generate
  done
}

watch_qtpl