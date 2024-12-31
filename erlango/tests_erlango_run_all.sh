for TESTNAME in $(grep -riI "[ ]*func[ ]*Test" | awk -F "(" '{print $1}' | awk '{print $2}'); do 
  echo $TESTNAME; 
  go test -run $TESTNAME *.go

done
