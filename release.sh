FILES=dist/*
mkdir -p "release"
for f in $FILES
do
  echo "Processing $f file..."
  if [[ $f == *"windows"* ]]
  then
    cp "$f" gvm.exe
    zip "$f.zip" gvm.exe
    mv "$f.zip" ./release
    rm -rf gvm.exe
  else
    cp "$f" "dist/gvm"
    tar -C dist -zcvf "$f.tar.gz" gvm
    mv "$f.tar.gz" ./release
    rm -rf dist/gvm
  fi
  rm -rf "$f"
done