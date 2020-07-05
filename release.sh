FILES=dist/*
mkdir -p "release"
for f in $FILES
do
  echo "Processing $f file..."
  cp "$f" "gvm"
  if [[ $f == *"windows"* ]]
  then
    zip "$f.zip" "gvm"
    mv "$f.zip" ./release
  else
    tar -zcvf "$f.tar.gz" "gvm"
    mv "$f.tar.gz" ./release
  fi
  rm -rf "$f"
done