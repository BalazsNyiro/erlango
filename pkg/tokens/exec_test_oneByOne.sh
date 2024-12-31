for testName in $(grep -rniI "func *Test" . | awk '{print $2}' | awk -F '(' '{print $1}' ); do
  echo
  echo ===== $testName =====
  go test -v -run $testName
  read -p "press ENTER to the next (prev was: $testName)" 
done
