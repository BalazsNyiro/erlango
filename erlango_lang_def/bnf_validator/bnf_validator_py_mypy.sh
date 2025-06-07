#!/usr/bin/env bash
clear

if [ "$(which mypy > /dev/null; echo $?)" -eq "1" ]; then
  echo "please install mypy: 'pip install mypy'"
else
  # mypy src/p0_base/prg_general_config_and_state.py
  # mypy src/p1_pixels/img_10_pixels.py


  for FILE_PY in *.py; do
    echo
    echo "file to check with mypy: $FILE_PY"
    mypy $FILE_PY  --check-untyped-defs
  done

fi

