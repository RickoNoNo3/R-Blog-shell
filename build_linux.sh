# make built folder
echo "Make built folder"
mkdir built
built_folder=built/`cat version`_linux_`date +%Y%m%d_%H%M%S`
mkdir $built_folder

# go build
echo "Run go build"
go build -o "$built_folder/blog_shell.run" .

# copy files to built folder
echo "Copy config.json"
cp config.json $built_folder/

echo "Build done. Check $built_folder."
